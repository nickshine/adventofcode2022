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
		fmt.Printf("%03d ", i)
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

		switch {
		case prev.x == cur.x: // vertical path
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

func setFloor(grid [][]*point, y int) {
	for x := 0; x < len(grid); x++ {
		grid[y][x] = &point{x, y, '#'}
	}
}

func parsePaths(in string, size, offset int) [][]*point {
	paths := strings.Split(strings.TrimSpace(in), "\n")

	grid := make([][]*point, size)
	for i := 0; i < size; i++ {
		grid[i] = make([]*point, size)
		for j := 0; j < len(grid[i]); j++ {
			grid[i][j] = &point{j + offset, i, '.'}
		}
	}

	maxY := 0
	for _, path := range paths {
		points := pathToPoints(path)
		for _, p := range points {
			if p.y > maxY {
				maxY = p.y
			}
			setPoint(grid, p, offset)
		}
	}

	setFloor(grid, maxY+2)

	return grid

}

func fall(grid [][]*point, p *point, offset int) {
	x, y := p.x-offset, p.y

	d := grid[y+1][x]
	dl := grid[y+1][x-1]
	dr := grid[y+1][x+1]

	switch {
	case p.isBlocked():
		fmt.Println("Start blocked")
	case d.isAir():
		fall(grid, d, offset)
	case d.isBlocked():
		if !dl.isBlocked() {
			fall(grid, dl, offset)
		} else if !dr.isBlocked() {
			fall(grid, dr, offset)
		} else {
			p.v = 'o' // at rest
		}
	default:
		panic("unexpected condition")
	}
}

func run(in string, size, offset, cycles int) {
	grid := parsePaths(in, size, offset)

	start := &point{500, 0, '+'}
	setPoint(grid, start, offset)

	for i := 0; i < cycles; i++ {
		fall(grid, start, offset)
	}
	display(grid)
}

func main() {
	run(exampleInput, 35, 485, 94)
	run(input, 500, 250, 26461)
}
