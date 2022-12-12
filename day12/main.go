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

type square struct {
	x, y      int
	elevation rune
	score     int
	visited   bool
}

func (s *square) String() string {
	return fmt.Sprintf("Elevation: %c, score: %d\n", s.elevation, s.score)
}

type heightMap [][]*square

func (h heightMap) String() string {
	var sb strings.Builder
	for _, row := range [][]*square(h) {
		for _, col := range row {
			sb.WriteRune(col.elevation)
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func readMap(in string) heightMap {
	lines := strings.Split(strings.TrimSpace(in), "\n")
	hmap := heightMap(make([][]*square, len(lines)))

	for i, l := range lines {
		hmap[i] = make([]*square, len(l))

		for j, r := range l {
			hmap[i][j] = &square{x: i, y: j, elevation: r}
		}
	}

	return hmap
}

func visit(m heightMap, x, y int, prev *square, score int) {

	if x < 0 || x >= len(m) || y < 0 || y >= len(m[x]) { // out of bounds
		return
	}

	s := m[x][y]

	if s.visited && s.score <= score { // shorter path already exists
		return
	}

	if s.elevation-prev.elevation > 1 { // too steep to climb
		return
	}

	s.visited = true
	s.score = score

	visit(m, x-1, y, s, score+1)
	visit(m, x, y+1, s, score+1)
	visit(m, x+1, y, s, score+1)
	visit(m, x, y-1, s, score+1)
}

func part1(in string) int {
	m := readMap(in)
	var start, end *square

	for _, row := range m {
		for _, square := range row {
			switch square.elevation {
			case 'S':
				start = square
				start.elevation = 'a'
			case 'E':
				end = square
				end.elevation = 'z'
			}
		}
	}

	visit(m, start.x, start.y, start, 0)

	return end.score
}

func part2(in string) int {
	m := readMap(in)
	var start, end *square
	var startSquares []*square

	for _, row := range m {
		for _, square := range row {
			switch square.elevation {
			case 'S':
				start = square
				start.elevation = 'a'
				startSquares = append(startSquares, start)
			case 'E':
				end = square
				end.elevation = 'z'
			case 'a':
				startSquares = append(startSquares, square)
			}
		}
	}

	min := math.MaxInt
	for _, s := range startSquares {
		visit(m, s.x, s.y, s, 0)
		if end.score < min {
			min = end.score
		}
	}

	return min
}

func main() {
	fmt.Printf("Part 1 example: %d\n", part1(exampleInput))
	fmt.Printf("Part 1: %d\n", part1(input))
	fmt.Printf("Part 2 example: %d\n", part2(exampleInput))
	fmt.Printf("Part 2: %d\n", part2(input))
}
