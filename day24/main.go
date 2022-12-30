package main

import (
	_ "embed"
	"fmt"
	"strings"
)

//go:embed example.txt
var exampleInput string

//go:embed input.txt
var input string

type point struct {
	x, y   int
	isWall bool
}

type grid [][]point

type state struct {
	p    point
	time int
}

func parseInput(in string) (grid, map[point][]rune) {
	rows := strings.Split(strings.TrimSpace(in), "\n")

	grid := make([][]point, len(rows))
	blizzards := map[point][]rune{}

	for y, row := range rows {
		grid[y] = make([]point, len(row))
		for x, c := range row {
			p := point{x: x, y: y}
			switch c {
			case '#':
				p.isWall = true
			case '^', 'v', '<', '>':
				blizzards[p] = append(blizzards[p], c)
			}
			grid[y][x] = p

		}

	}

	return grid, blizzards
}

func (g grid) display(blizzards map[point][]rune) {
	for _, row := range g {
		for _, p := range row {
			switch {
			case p.isWall:
				fmt.Printf("#")
			case len(blizzards[p]) == 1:
				fmt.Printf("%c", blizzards[p][0])
			case len(blizzards[p]) == 2:
				fmt.Printf("2")
			default:
				fmt.Printf(".")
			}
		}
		fmt.Println()
	}
}

func (g grid) nextBlizzards(blizzards map[point][]rune) map[point][]rune {
	nextBlizzards := map[point][]rune{}

	for p, b := range blizzards {
		x, y := p.x, p.y
		for _, blizzard := range b {
			var next point
			switch blizzard {
			case '^':
				next = g[y-1][x]
				if next.isWall {
					next = g[len(g)-2][x]
				}
			case 'v':
				next = g[y+1][x]
				if next.isWall {
					next = g[1][x]
				}
			case '<':
				next = g[y][x-1]
				if next.isWall {
					next = g[y][len(g[y])-2]
				}
			case '>':
				next = g[y][x+1]
				if next.isWall {
					next = g[y][1]
				}
			}

			nextBlizzards[next] = append(nextBlizzards[next], blizzard)
		}
	}

	return nextBlizzards
}

func (g grid) allBlizzards(blizzards map[point][]rune, maxTime int) map[int]map[point][]rune {

	allBlizzards := map[int]map[point][]rune{}
	for i := 0; i < maxTime; i++ {
		allBlizzards[i] = blizzards
		blizzards = g.nextBlizzards(blizzards)
	}

	return allBlizzards
}

func bfs(time int, g grid, start, end point, blizzards map[int]map[point][]rune) int {
	queue := []state{{start, time}}
	seen := map[state]struct{}{queue[0]: {}}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		x, y := cur.p.x, cur.p.y
		choices := []point{
			{x: x, y: y - 1}, // up
			{x: x, y: y + 1}, // down
			{x: x - 1, y: y}, // left
			{x: x + 1, y: y}, // right
			cur.p,            // wait
		}

		for _, p := range choices {
			if p.y < 0 || p.y >= len(g) {
				continue
			}
			p = g[p.y][p.x]
			next := state{p, cur.time + 1}

			if next.p == end {
				return next.time
			}

			// if current position and time has already been seen, then no need to re-add choices
			if _, ok := seen[next]; ok {
				continue
			}

			if next.p.isWall {
				continue
			}

			if _, ok := blizzards[next.time][next.p]; ok {
				continue
			}

			seen[next] = struct{}{}
			queue = append(queue, next)
		}

	}
	return -1
}

func part1(in string, maxTime int) int {
	grid, blizzards := parseInput(in)
	allBlizzards := grid.allBlizzards(blizzards, maxTime)

	start := point{x: 1, y: 0}
	end := point{x: len(grid[0]) - 2, y: len(grid) - 1}

	return bfs(0, grid, start, end, allBlizzards)
}

func part2(in string, maxTime int) int {
	grid, blizzards := parseInput(in)
	allBlizzards := grid.allBlizzards(blizzards, maxTime)

	start := point{x: 1, y: 0}
	end := point{x: len(grid[0]) - 2, y: len(grid) - 1}

	return bfs(bfs(bfs(0, grid, start, end, allBlizzards), grid, end, start, allBlizzards), grid, start, end, allBlizzards)
}

func main() {
	fmt.Printf("Part 1 example: %d\n", part1(exampleInput, 30))
	fmt.Printf("Part 1: %d\n", part1(input, 300))
	fmt.Printf("Part 2 example: %d\n", part2(exampleInput, 60))
	fmt.Printf("Part 2: %d\n", part2(input, 800))
}
