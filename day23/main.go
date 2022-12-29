package main

import (
	_ "embed"
	"fmt"
	"math"
	"strings"
)

//go:embed example.txt
var exampleInput string

//go:embed input.txt
var input string

type direction int

const (
	NORTH direction = iota
	SOUTH
	WEST
	EAST
)

type elf struct {
	firstDirection direction
}

type position struct {
	x, y int
	elf  *elf
}

type grid [][]position

func parseInput(in string, size int) grid {
	rows := strings.Split(strings.TrimSpace(in), "\n")

	g := make([][]position, size)

	for y := range g {
		g[y] = make([]position, size)
		for x := range g[y] {
			g[y][x] = position{x: x, y: y}
		}
	}

	// place input in center of larger grid
	start := size/2 - len(rows)/2 - 1

	for i, row := range rows {
		for j, c := range row {
			if c == '#' {
				p := g[start+i][start+j]
				p.elf = &elf{NORTH}
				g[start+i][start+j] = p
			}
		}
	}

	return g

}

// crop determines the smallest rectangle within the grid that contains all elfs.
func (g grid) crop() (x, y, width, height int) {

	var ymax, xmax int
	xmin := math.MaxInt32
	ymin := xmin

	for y, row := range g {
		for x, p := range row {
			if p.elf == nil {
				continue
			}
			if y < ymin {
				ymin = y
			}
			if x < xmin {
				xmin = x
			}
			if y > ymax {
				ymax = y
			}
			if x > xmax {
				xmax = x
			}
		}
	}

	return xmin, ymin, xmax + 1 - xmin, ymax + 1 - ymin
}

func (g grid) countGround(x, y, width, height int) int {
	sum := 0
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			p := g[y+i][x+j]
			if p.elf == nil {
				sum++
			}
		}
	}

	return sum
}

func (g grid) display() {
	fmt.Printf("    %s%s%s\n", strings.Repeat(" ", 10), strings.Repeat("1", 10), strings.Repeat("2", 10))
	fmt.Printf("    %s\n", strings.Repeat("0123456789", 3))
	for i, row := range g {
		fmt.Printf("%3d ", i)
		for _, p := range row {
			if p.elf != nil {
				fmt.Printf("#")
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Println()
	}
}

func (g grid) allClear(p position) bool {
	x, y := p.x, p.y

	nw := g[y-1][x-1]
	n := g[y-1][x]
	ne := g[y-1][x+1]
	e := g[y][x+1]
	se := g[y+1][x+1]
	s := g[y+1][x]
	sw := g[y+1][x-1]
	w := g[y][x-1]

	return nw.elf == nil &&
		n.elf == nil &&
		ne.elf == nil &&
		e.elf == nil &&
		se.elf == nil &&
		s.elf == nil &&
		sw.elf == nil &&
		w.elf == nil
}

func (g grid) nextPosition(p position) (position, bool) {
	if p.elf == nil {
		return position{}, false
	}

	x, y := p.x, p.y

	var p1, p2, p3 position
	var nextPosition position
	var hasAvailableDirection bool

	d := p.elf.firstDirection
	for i := 0; i < 4; i++ {
		switch d {
		case NORTH:
			p1 = g[y-1][x-1] // nw
			p2 = g[y-1][x]   // n
			p3 = g[y-1][x+1] // ne
		case SOUTH:
			p1 = g[y+1][x-1] // sw
			p2 = g[y+1][x]   // s
			p3 = g[y+1][x+1] // se
		case WEST:
			p1 = g[y-1][x-1] // nw
			p2 = g[y][x-1]   // w
			p3 = g[y+1][x-1] // sw
		case EAST:
			p1 = g[y-1][x+1] // ne
			p2 = g[y][x+1]   // e
			p3 = g[y+1][x+1] // se
		}

		if p1.elf == nil && p2.elf == nil && p3.elf == nil {
			nextPosition = p2
			hasAvailableDirection = true
			break
		}

		d = (d + 1) % 4
	}

	p.elf.firstDirection = (p.elf.firstDirection + 1) % 4

	if hasAvailableDirection {
		return nextPosition, true
	}

	return position{}, false
}

func proposedPositions(g grid) map[position][]position {
	proposedPositions := map[position][]position{}
	for _, row := range g {
		for _, p := range row {
			if p.elf == nil {
				continue
			}

			if g.allClear(p) {
				p.elf.firstDirection = (p.elf.firstDirection + 1) % 4
				continue
			}

			if nextP, ok := g.nextPosition(p); ok {
				proposedPositions[nextP] = append(proposedPositions[nextP], p)
			}
		}
	}

	return proposedPositions
}

func (g grid) move(to, from position) {
	to.elf = from.elf
	g[to.y][to.x] = to
	from.elf = nil
	g[from.y][from.x] = from
}

func (g grid) moveAll(proposed map[position][]position) {
	for p, v := range proposed {
		// if more than one elf proposed this position, skip
		if len(v) > 1 {
			continue
		}
		g.move(p, v[0])
	}
}

func part1(in string, size int) int {
	grid := parseInput(in, size)

	rounds := 10

	for round := 0; round < rounds; round++ {
		proposedPositions := proposedPositions(grid)
		grid.moveAll(proposedPositions)
	}

	x, y, width, height := grid.crop()
	total := grid.countGround(x, y, width, height)
	return total
}

func part2(in string, size int) int {
	grid := parseInput(in, size)

	rounds := 1
	for {
		proposedPositions := proposedPositions(grid)
		if len(proposedPositions) == 0 {
			break
		}

		grid.moveAll(proposedPositions)
		rounds++
	}

	return rounds
}

func main() {
	fmt.Printf("Part 1 example: %d\n", part1(exampleInput, 30))
	fmt.Printf("Part 1: %d\n", part1(input, 150))
	fmt.Printf("Part 2 example: %d\n", part2(exampleInput, 30))
	fmt.Printf("Part 2: %d\n", part2(input, 200))
}
