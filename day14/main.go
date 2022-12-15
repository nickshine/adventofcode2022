package main

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

const (
	gridSize = 180
	xOffset  = 400
	maxFall  = 5000
)

//go:embed example.txt
var exampleInput string

//go:embed input.txt
var input string

type point struct {
	x, y int
	v    rune
}

func (p *point) String() string {
	return fmt.Sprintf("%c", p.v)
}

func (p *point) isAir() bool {
	return p.v == '.'
}

func (p *point) isRock() bool {
	return p.v == '#'

}
func (p *point) isSand() bool {
	return p.v == 'o'
}

func (p *point) isBlocked() bool {
	return p.isRock() || p.isSand()
}

func display(grid [][]*point) {
	for i := 0; i < len(grid); i++ {
		fmt.Printf("%d ", i)
		for j := 0; j < len(grid[i]); j++ {
			fmt.Printf("%s", grid[i][j])
		}
		fmt.Println()
	}
}

func toInt(c []string) (int, int) {
	if len(c) != 2 {
		panic("invalid input")
	}

	x, err := strconv.Atoi(c[0])
	if err != nil {
		panic("invalid input")
	}
	y, err := strconv.Atoi(c[1])
	if err != nil {
		panic("invalid input")
	}

	return x, y
}

func pathToPoints(path string) []*point {
	parts := strings.Split(path, " -> ")
	var coords []point
	for _, p := range parts {
		coord := strings.Split(p, ",")
		x, y := toInt(coord)
		coords = append(coords, point{x, y, '#'})
	}

	var points []*point
	for i := 1; i < len(coords); i++ {
		prev := coords[i-1]
		cur := coords[i]

		// fmt.Printf("cur: %#v, prev: %#v\n", cur, prev)

		switch {
		case prev.x == cur.x: // vertical path
			// fmt.Printf("prev.x == cur.x, %d == %d\n", prev.x, cur.x)
			x := prev.x
			var start, end int
			if prev.y < cur.y {
				start, end = prev.y, cur.y
			} else if prev.y > cur.y {
				start, end = cur.y, prev.y
			} else {
				panic("unexpected x and y both equal")
			}

			for y := start; y <= end; y++ {
				points = append(points, &point{x, y, '#'})
			}
		case prev.y == cur.y: // horizontal path
			// fmt.Printf("prev.y == cur.y, %d == %d\n", prev.y, cur.y)
			y := prev.y
			var start, end int
			if prev.x < cur.x {
				start, end = prev.x, cur.x
			} else if prev.x > cur.x {
				start, end = cur.x, prev.x
			} else {
				panic("unexpected x and y both equal")
			}

			for x := start; x <= end; x++ {
				points = append(points, &point{x, y, '#'})
			}
		default:
			panic("path is not vertical or horizontal")

		}
	}

	return points
}

func setPoint(grid [][]*point, p *point, offset int) {
	grid[p.y][p.x-offset] = p
}

func parsePaths(in string) [][]*point {
	paths := strings.Split(strings.TrimSpace(in), "\n")

	grid := make([][]*point, gridSize)
	for i := 0; i < gridSize; i++ {
		grid[i] = make([]*point, gridSize)
		for j := 0; j < len(grid[i]); j++ {
			grid[i][j] = &point{j + xOffset, i, '.'}
		}
	}

	for _, path := range paths {
		points := pathToPoints(path)
		for _, p := range points {
			// fmt.Printf("x: %d, y: %d\n", p.x, p.y)
			setPoint(grid, p, xOffset)
		}
	}

	return grid

}

func fall(grid [][]*point, p *point, offset int, count int) {

	x, y := p.x-offset, p.y

	d := grid[y+1][x]
	dl := grid[y+1][x-1]
	dr := grid[y+1][x+1]

	switch {
	case d.isAir():
		fall(grid, d, offset, count+1)
	case d.isBlocked():
		if !dl.isBlocked() {
			fall(grid, dl, offset, count+1)
		} else if !dr.isBlocked() {
			fall(grid, dr, offset, count+1)
		} else {
			if count == maxFall {
				panic("STOP")
			}
			p.v = 'o' // at rest
			// fmt.Printf("at rest: %d,%d\n", p.x, p.y)
		}
	default:
		panic("unexpected condition")
	}
}

func run(in string) {
	grid := parsePaths(in)

	start := &point{500, 0, '+'}
	setPoint(grid, start, xOffset)

	for i := 0; i < 888; i++ {
		fall(grid, start, xOffset, 0)
	}
	display(grid)
}

func main() {
	run(exampleInput)
	// run(input)
}
