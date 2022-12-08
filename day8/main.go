package main

import (
	_ "embed"
	"fmt"
	"strings"
)

//go:embed input.txt
var input string

func readGrid(in string) [][]int {
	rows := strings.Split(strings.TrimSpace(in), "\n")
	grid := make([][]int, len(rows))

	for i, row := range rows {
		grid[i] = make([]int, len(row))
		for j, col := range row {
			grid[i][j] = int(col - '0')
		}
	}

	return grid
}

func calcPerimeter(grid [][]int) int {
	p1 := len(grid)
	p2 := len(grid[0])

	return 2*(p1+p2) - 4
}

func isVisible(grid [][]int, tree, x, y, dx, dy int) bool {
	if y < 0 || x < 0 || x >= len(grid) || y >= len(grid[0]) {
		return true
	}

	return grid[x][y] < tree && isVisible(grid, tree, x+dx, y+dy, dx, dy)
}

func isPerimeter(grid [][]int, x, y int) bool {

	return x == 0 || x == len(grid)-1 || y == 0 || y == len(grid[x])-1
}

func part1(grid [][]int) int {
	visible := 0

	for x, row := range grid {
		for y, tree := range row {
			if isPerimeter(grid, x, y) {
				visible++
				continue
			}

			if isVisible(grid, tree, x, y-1, 0, -1) || // left
				isVisible(grid, tree, x, y+1, 0, 1) || // right
				isVisible(grid, tree, x-1, y, -1, 0) || // up
				isVisible(grid, tree, x+1, y, 1, 0) { // down
				visible++
			}
		}
	}

	return visible
}

func score(grid [][]int, tree, x, y, dx, dy int) int {

	if y < 0 || x < 0 || x >= len(grid) || y >= len(grid[0]) {
		return 0
	}

	if grid[x][y] >= tree {
		return 1
	}

	return 1 + score(grid, tree, x+dx, y+dy, dx, dy)
}

func part2(grid [][]int) int {
	max := 0
	for x, row := range grid {
		for y, tree := range row {
			score := score(grid, tree, x, y-1, 0, -1) *
				score(grid, tree, x, y+1, 0, 1) *
				score(grid, tree, x-1, y, -1, 0) *
				score(grid, tree, x+1, y, 1, 0)
			if score > max {
				max = score
			}
		}
	}

	return max
}

func main() {
	fmt.Printf("part 1: %d\n", part1(readGrid(input)))
	fmt.Printf("part 2: %d\n", part2(readGrid(input)))
}
