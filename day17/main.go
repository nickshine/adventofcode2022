package main

import (
	_ "embed"
	"fmt"
	"strings"

	"log"
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

func part1(in string, count int) int {
	jets := parseJets(in)
	var maxHeight, numStopped, shapeType int

	var rocks []*rock

	jetIdx := 0
	for numStopped < count {

		r := newRock(2, maxHeight+3, shapes[shapeType])
		// log.Printf("new rock: %s", r)

		// run current rock sequence until stopped
		for !r.stopped {
			push(r, rocks, jets[jetIdx])
			fall(r, rocks)
			jetIdx = (jetIdx + 1) % len(jets)
		}

		if r.y+r.height > maxHeight {
			maxHeight = r.y + r.height
			// log.Printf("new maxHeight: %d", maxHeight)
		}

		rocks = append(rocks, r)
		numStopped++
		shapeType = (shapeType + 1) % 5

	}

	for i, rock := range rocks {
		log.Printf("%d: %s", i, rock)
	}

	return maxHeight
}

func main() {
	// fmt.Printf("Part 1 example: %d\n", part1(exampleInput, 2022))
	fmt.Printf("Part 1: %d\n", part1(input, 2022))
	// fmt.Printf("Part 2 example: %d\n", part2(exampleInput))
	// fmt.Printf("Part 2: %d\n", part2(input))
}
