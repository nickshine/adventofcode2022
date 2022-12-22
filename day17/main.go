package main

import (
	"bytes"
	"crypto/sha1"
	_ "embed"
	"fmt"
	"log"
	"strings"
)

//go:embed example.txt
var exampleInput string

//go:embed input.txt
var input string

const (
	maxWidth = 7
)

var (
	shapes = []shape{
		{
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{1, 1, 1, 1},
		},
		{
			{0, 0, 0, 0},
			{0, 1, 0, 0},
			{1, 1, 1, 0},
			{0, 1, 0, 0},
		},
		{
			{0, 0, 0, 0},
			{0, 0, 1, 0},
			{0, 0, 1, 0},
			{1, 1, 1, 0},
		},
		{
			{1, 0, 0, 0},
			{1, 0, 0, 0},
			{1, 0, 0, 0},
			{1, 0, 0, 0},
		},
		{
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{1, 1, 0, 0},
			{1, 1, 0, 0},
		},
	}
)

type shape [4][4]int

type rock struct {
	x, y          int // the coordinate of bottom left edge of a 4x4 grid of units
	width, height int
	s             shape // a 4x4 grid where 1's indicate rock and 0's indicate empty space
	jetIndex      int
	stopped       bool
}

func newRock(x, y int, s shape) *rock {
	r := &rock{x: x, y: y, s: s}

	// set width and height
	maxx, maxy := 0, 4
	// for each column
	for y := 0; y < len(r.s); y++ {
		hasY := false
		for x := 0; x < len(r.s); x++ {
			if r.s[y][x] == 0 {
				continue
			}

			hasY = true
			if x > maxx {
				maxx = x
			}
		}
		if !hasY {
			maxy--
		}

	}

	r.width = maxx + 1
	r.height = maxy

	return r

}

func (r *rock) String() string {
	return fmt.Sprintf("{x:%d,y:%d,w:%d,h:%d}", r.x, r.y, r.width, r.height)
}

func (r *rock) overlapping(t *rock) bool {

	if t.y+t.height <= r.y {
		// log.Printf("%s lower than current %s, skipping", t, r)
		return false
	}

	// ts left or right of current cannot block
	if t.x+t.width <= r.x || t.x >= r.x+r.width {
		// log.Printf("%s left or right of current %s, skipping", t, r)
		return false
	}

	for y := 0; y < len(r.s); y++ {
		for x := 0; x < len(r.s); x++ {
			if r.s[y][x] == 0 { // empty space can overlap
				continue
			}
			// is r.s[y][x] overlapping any of t.s?
			for ty := 0; ty < len(t.s); ty++ {
				for tx := 0; tx < len(t.s); tx++ {
					if t.s[ty][tx] == 0 {
						continue
					}

					if r.x+x != t.x+tx {
						continue
					}

					if r.y+len(r.s)-1-y != t.y+len(t.s)-1-ty {
						continue
					}

					// log.Printf("%s overlaps %s at %d,%d", r, t, r.x+x, r.y+len(r.s)-1-y)
					return true
				}
			}
		}
	}
	return false
}

func (r *rock) draw() {
	for _, v := range r.s {
		for _, vv := range v {
			if vv == 1 {
				fmt.Printf("#")
			} else {
				fmt.Printf(" ")
			}
		}
		fmt.Println()
	}
}

func parseJets(in string) []int {
	s := strings.TrimSpace(in)
	jets := make([]int, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] == '<' {
			jets[i] = -1
		} else {
			jets[i] = 1
		}
	}

	return jets
}

func push(r *rock, rocks []*rock, push int) bool {
	x := r.x + push
	// cannot push past left or right wall
	if x < 0 || x+r.width > maxWidth {
		return false
	}

	curx := r.x
	r.x += push
	for _, rock := range rocks {
		if r.overlapping(rock) {

			r.x = curx // put r.x back
			return false
		}
	}

	// log.Printf("push %d, to %s", push, r)
	return true
}

func fall(r *rock, rocks []*rock) bool {
	// cannot fall below floor
	if r.y == 0 {
		// log.Printf("rock %s overlapping floor, stopping", r)
		r.stopped = true
		return false
	}

	curY := r.y
	r.y--
	for _, rock := range rocks {
		if r.overlapping(rock) {

			r.y = curY // put r.y back
			// log.Printf("rock %s overlapping %s, stopping", r, rock)
			r.stopped = true
			return false
		}
	}

	// log.Printf("fall one to %s", r)
	return true
}

func run(jets []int, count int) (int, []*rock) {
	var maxHeight, numStopped, shapeType int

	var rocks []*rock

	jetIdx := -1
	for numStopped < count {

		r := newRock(2, maxHeight+3, shapes[shapeType])
		// log.Printf("new rock: %s", r)

		// run current rock sequence until stopped
		for {
			jetIdx = (jetIdx + 1) % len(jets)
			push(r, rocks, jets[jetIdx])
			fall(r, rocks)
			if r.stopped {
				break
			}
		}

		r.jetIndex = jetIdx

		if r.y+r.height > maxHeight {
			maxHeight = r.y + r.height
			// log.Printf("new maxHeight: %d", maxHeight)
		}

		rocks = append(rocks, r)
		numStopped++
		shapeType = (shapeType + 1) % 5

	}

	return maxHeight, rocks

}

func nRows(rocks []*rock, ystart, count int) [][]int {
	rows := make([][]int, count)
	for i := 0; i < len(rows); i++ {
		rows[i] = make([]int, maxWidth)
	}

	for _, r := range rocks {
		if r.y >= ystart+len(rows) {
			// log.Printf("r.y out of bounds: r.y:%d, ystart:%d", r.y, ystart)
			continue
		}

		if r.y+r.height <= ystart {
			// log.Printf("r.y out of bounds: r.y+height:%d, ystart:%d", r.y+r.height, ystart)
			continue
		}

		for y := 0; y < len(r.s); y++ {
			for x := 0; x < len(r.s); x++ {
				if r.s[y][x] == 0 {
					continue
				}

				yy := r.y + len(r.s) - 1 - y
				if yy-ystart >= len(rows) || yy-ystart < 0 {
					// log.Printf("yy out of bounds with ystart: yy:%d, ystart:%d", yy, ystart)

					continue
				}

				rows[yy-ystart][r.x+x] = 1
				// fmt.Fprintf(&buf, "%d:", r.jetIndex)
			}
		}
	}

	return rows

}

// shaNRows creates the chamber rows from the returned rocks, then returns a sha of them.
//
// Since part1 was done without using a grid, the rows are created from the rock data.
func shaNRows(rocks []*rock, ystart, count int) string {
	var buf bytes.Buffer

	rows := nRows(rocks, ystart, count)

	for i := 0; i < len(rows); i++ {
		for j := 0; j < len(rows[i]); j++ {
			fmt.Fprintf(&buf, "%d", rows[i][j])
		}
	}
	sha := fmt.Sprintf("%x", sha1.Sum(buf.Bytes()))

	return sha
}

func shaNCols(rocks []*rock, ystart, xoffset, count int) string {
	var buf bytes.Buffer

	rows := nRows(rocks, ystart, count)

	for i := 0; i < len(rows); i++ {
		fmt.Fprintf(&buf, "%d", rows[i][xoffset])
	}

	sha := fmt.Sprintf("%x", sha1.Sum(buf.Bytes()))

	return sha
}

func part1(in string, count int) int {
	jets := parseJets(in)
	maxHeight, _ := run(jets, count)
	return maxHeight
}

func part2(in string) int {
	jets := parseJets(in)

	log.Printf("len jets: %d", len(jets))

	const limit = 10000
	const rowCount = 25

	maxHeight, rocks := run(jets, limit)
	log.Printf("%d,%d", maxHeight, len(rocks))
	sha := shaNCols(rocks, 0, 0, rowCount) // get sha of first rowCount rows
	log.Printf(sha)

	// take sha of first n rocks, and look for that same sha again to find cycle
	var cycles []int
	for y := rowCount; y < maxHeight; y++ {
		shab := shaNCols(rocks, y, 0, rowCount)
		if sha == shab {
			log.Printf("same sha at %d", y)
			cycles = append(cycles, y)
		}
	}

	var cycle int
	for i := 1; i < len(cycles); i++ {
		cycle = cycles[i] - cycles[i-1]

		log.Printf("cycle: %d", cycle)
	}

	maxHeight, _ = run(jets, 1740)
	log.Printf("Max height: %d", maxHeight)

	rockCycle := 1740

	const max = 1000000000000
	factor := max / rockCycle
	log.Printf("factor: %d", factor)
	remainder := max % rockCycle
	log.Printf("remainder: %d", remainder)

	remainderHeight, _ := run(jets, remainder)
	log.Printf("remainder height: %d", remainderHeight)

	remainderHeightPlusOneCycle, _ := run(jets, remainder+rockCycle)
	log.Printf("remainder +one cycle height: %d", remainderHeightPlusOneCycle)

	total := factor*cycle + remainderHeight
	return total
}

func main() {
	// fmt.Printf("Part 1 example: %d\n", part1(exampleInput, 2022))
	fmt.Printf("Part 1: %d\n", part1(input, 2022))
	// fmt.Printf("Part 2 example: %d\n", part2(exampleInput))
	fmt.Printf("Part 2: %d\n", part2(input))
}
