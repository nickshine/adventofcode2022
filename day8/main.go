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
	// out of bounds
	if y < 0 || x < 0 || x >= len(grid) || y >= len(grid[0]) {
		fmt.Printf("Out of bounds: grid[%d][%d]\n", x, y)
		return true
	}

	// fmt.Printf("grid[%d][%d]: %d < %d (tree) ? %t\n", x, y, grid[x][y], tree, grid[x][y] < tree)
	return grid[x][y] < tree && isVisible(grid, tree, x+dx, y+dy, dx, dy)
}

func isPerimeter(grid [][]int, x, y int) bool {

	return x == 0 || x == len(grid)-1 || y == 0 || y == len(grid[x])-1
}

func part1(grid [][]int) int {

	fmt.Println(grid)
	fmt.Printf("perimeter size: %d\n", calcPerimeter(grid))

	visible := 0

	fmt.Printf("len grid: %d x %d\n", len(grid), len(grid[0]))

	for x, row := range grid {
		for y, tree := range row {
			if isPerimeter(grid, x, y) {
				// fmt.Printf("Perimeter: grid[%d][%d]: %d\n", x, y, grid[x][y])
				visible++
				continue
			}

			// fmt.Printf("Non perimeter: grid[%d][%d]: %d, checking left\n", x, y, grid[x][y])
			if isVisible(grid, tree, x, y-1, 0, -1) || // left
				isVisible(grid, tree, x, y+1, 0, 1) || // right
				isVisible(grid, tree, x-1, y, -1, 0) || // up
				isVisible(grid, tree, x+1, y, 1, 0) { // down
				// fmt.Println("is visible")
				visible++
			}
		}
	}

	return visible
}

func main() {
	fmt.Printf("part 1: %d\n", part1(readGrid(input)))
	// fmt.Printf("part 2: %d\n", part2())
}
