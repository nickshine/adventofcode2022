package main

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

//go:embed example.txt
var exampleInput string

//go:embed input.txt
var input string

const (
	Empty = iota
	Filled
)

type grid [][][]int

func parse(in string) (positions [][3]int, maxX int, maxY int, maxZ int) {
	lines := strings.Split(strings.TrimSpace(in), "\n")

	for _, l := range lines {
		xyz := strings.Split(l, ",")
		if len(xyz) != 3 {
			panic("invalid input")
		}
		x, err := strconv.Atoi(xyz[0])
		if err != nil {
			panic("invalid input")
		}
		if x > maxX {
			maxX = x
		}

		y, err := strconv.Atoi(xyz[1])
		if err != nil {
			panic("invalid input")
		}
		if y > maxY {
			maxY = y
		}
		z, err := strconv.Atoi(xyz[2])
		if err != nil {
			panic("invalid input")
		}
		if z > maxZ {
			maxZ = z
		}

		positions = append(positions, [3]int{x, y, z})
	}

	return positions, maxX + 1, maxY + 1, maxZ + 1
}

// newGrid creates a 3d slice of ints - [z][y][x]
func newGrid(maxX, maxY, maxZ int) grid {
	grid := make([][][]int, maxZ)
	for i := 0; i < maxZ; i++ {
		grid[i] = make([][]int, maxY)
		for j := 0; j < maxY; j++ {
			grid[i][j] = make([]int, maxX)
		}
	}

	return grid
}

func (g grid) get(x, y, z int) (int, bool) {

	if z < 0 || y < 0 || x < 0 {
		return 0, false
	}

	if z >= len(g) || y >= len(g[z]) || x >= len(g[z][y]) {
		return 0, false
	}

	return g[z][y][x], true
}

// insert inserts the given cube's xyz and returns the surface area gained or lost.
// It is assumed the grid is large enough to hold the given xyz.
func (g grid) insert(x, y, z int) int {
	surfaceAreaDelta := 6 // 6 sides to a cube

	if g[z][y][x] != Empty {
		panic("invalid input")
	}

	// check all 6 sides of cube for adjacent:
	// z-axis
	if c, ok := g.get(x, y, z+1); ok {
		if c == Filled {
			surfaceAreaDelta -= 2 // current cube loses 1, and the adjacent loses 1
		}
	}
	if c, ok := g.get(x, y, z-1); ok {
		if c == Filled {
			surfaceAreaDelta -= 2
		}
	}

	// y-axis
	if c, ok := g.get(x, y+1, z); ok {
		if c == Filled {
			surfaceAreaDelta -= 2
		}
	}
	if c, ok := g.get(x, y-1, z); ok {
		if c == Filled {
			surfaceAreaDelta -= 2
		}
	}

	// x-axis
	if c, ok := g.get(x+1, y, z); ok {
		if c == Filled {
			surfaceAreaDelta -= 2
		}
	}
	if c, ok := g.get(x-1, y, z); ok {
		if c == Filled {
			surfaceAreaDelta -= 2
		}
	}

	g[z][y][x] = Filled

	return surfaceAreaDelta
}

func part1(in string) int {
	positions, maxX, maxY, maxZ := parse(in)
	grid := newGrid(maxX, maxY, maxZ)
	surfaceArea := 0

	for _, p := range positions {
		surfaceArea += grid.insert(p[0], p[1], p[2])

	}

	return surfaceArea
}

func main() {
	fmt.Printf("Part 1 example: %d\n", part1(exampleInput))
	fmt.Printf("Part 1: %d\n", part1(input))
	// fmt.Printf("Part 2 example: %d\n", part2(exampleInput))
	// fmt.Printf("Part 2: %d\n", part2(input))
}
