package main

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

//go:embed example.txt
var exampleInput string

//go:embed input.txt
var input string

type facing int

const (
	RIGHT facing = iota
	DOWN
	LEFT
	UP
)

const (
	_ = iota
	turnR
	_
	turnL
)

type step struct {
	turn     int
	distance int
}

type grid [][]rune

type translateFunc func(x, y int, f facing) (int, int, facing)

type region struct {
	xmin, xmax int
	ymin, ymax int
	translate  translateFunc
}

// getRegion returns the region that x,y is in from the provided array of regions.
// xmin, and ymin are included in the range, while xmax and ymax are excluded.
func getRegion(regions []*region, x, y int) *region {

	for _, r := range regions {
		if x >= r.xmin && x < r.xmax && y >= r.ymin && y < r.ymax {
			return r
		}
	}

	return nil
}

func turn(current facing, turn int) facing {
	return facing((int(current) + turn) % 4)
}

func parseInput(in string) (grid, []step) {
	parts := strings.Split(in, "\n\n")
	rows := strings.Split(parts[0], "\n")
	path := strings.TrimSpace(parts[1])

	grid := make([][]rune, len(rows))

	maxX := 0
	for _, row := range rows {
		if len(row) > maxX {
			maxX = len(row)
		}
	}

	for i, row := range rows {
		grid[i] = make([]rune, maxX)
		if len(row) < maxX {
			// prefill rows that are shorter than max
			for j := 0; j < maxX; j++ {
				grid[i][j] = ' '
			}
		}
		for j, c := range row {
			grid[i][j] = c
		}
	}

	var steps []step
	curTurn := turnR
	for i, j := 0, 0; j <= len(path); {
		if j == len(path) || unicode.IsLetter(rune(path[j])) {
			distance, err := strconv.Atoi(path[i:j])
			if err != nil {
				panic(err)
			}

			steps = append(steps, step{turn: curTurn, distance: distance})
		}

		if j == len(path) {
			break
		} else if unicode.IsDigit(rune(path[j])) {
			j++
		} else if unicode.IsLetter(rune(path[j])) {
			switch path[j] {
			case 'L':
				curTurn = turnL
			case 'R':
				curTurn = turnR
			}

			j++
			i = j
		}

	}

	return grid, steps

}

func display(g grid) {
	for _, row := range g {
		for _, col := range row {
			fmt.Printf("%c", col)
		}
		fmt.Println()
	}
}

func (g grid) findStart() (int, int) {

	row := g[0]
	for i, c := range row {
		if c == ' ' {
			continue
		} else {
			return 0, i
		}
	}

	return 0, 0
}

// dxy returns the delta x and y to move in the proper direction.
func (f facing) dxy() (int, int) {
	var dy, dx int
	switch f {
	case RIGHT:
		dy, dx = 0, 1
	case DOWN:
		dy, dx = 1, 0
	case LEFT:
		dy, dx = 0, -1
	case UP:
		dy, dx = -1, 0
	}

	return dx, dy
}

func (g grid) move(x, y, dx, dy int) (int, int, bool) {

	nextY := (y + dy) % len(g)
	if nextY < 0 {
		nextY = len(g) + nextY
	}
	nextX := (x + dx) % len(g[0])
	if nextX < 0 {
		nextX = len(g[0]) + nextX
	}
	switch g[nextY][nextX] {
	case '.':
		return nextX, nextY, true
	case '#':
		return x, y, false
	default:
		return g.move(nextX, nextY, dx, dy)
	}
}

func (g grid) cubeMove(x, y int, f facing, r *region) (int, int, facing, bool) {
	dx, dy := f.dxy()

	// get current location reg
	var nextX, nextY int
	var nextF facing

	switch {
	case x+dx >= r.xmax, x+dx < r.xmin, y+dy >= r.ymax, y+dy < r.ymin:
		nextX, nextY, nextF = r.translate(x, y, f)
	default:
		nextX, nextY, nextF = x+dx, y+dy, f
	}

	switch g[nextY][nextX] {
	case '.':
		return nextX, nextY, nextF, true
	case '#':
		return x, y, f, false
	default:
		panic("unexpected condition")
	}
}

func part1(in string) int {
	grid, steps := parseInput(in)
	y, x := grid.findStart()
	f := UP

	for _, step := range steps {
		f = turn(f, step.turn)
		dx, dy := f.dxy()
		for i := 0; i < step.distance; i++ {
			nextX, nextY, ok := grid.move(x, y, dx, dy)
			if !ok {
				break
			}
			x, y = nextX, nextY
		}

	}

	return 1000*(y+1) + 4*(x+1) + int(f)
}

func part2(in string, regions []*region) int {
	grid, steps := parseInput(in)
	y, x := grid.findStart()
	f := UP // start up so first turn will end in RIGHT facing

	for _, step := range steps {
		f = turn(f, step.turn)

		for i := 0; i < step.distance; i++ {
			r := getRegion(regions, x, y)
			nextX, nextY, nextF, ok := grid.cubeMove(x, y, f, r)
			if !ok {
				break
			}
			x, y, f = nextX, nextY, nextF
		}

	}

	return 1000*(y+1) + 4*(x+1) + int(f)

}

func main() {
	fmt.Printf("Part 1 example: %d\n", part1(exampleInput))
	fmt.Printf("Part 1: %d\n", part1(input))
	fmt.Printf("Part 2 example: %d\n", part2(exampleInput, exampleRegions()))
	fmt.Printf("Part 2: %d\n", part2(input, inputRegions()))
}

func exampleRegions() []*region {
	r1 := &region{xmin: 8, xmax: 12, ymin: 0, ymax: 4}
	r2 := &region{xmin: 0, xmax: 4, ymin: 4, ymax: 8}
	r3 := &region{xmin: 4, xmax: 8, ymin: 4, ymax: 8}
	r4 := &region{xmin: 8, xmax: 12, ymin: 4, ymax: 8}
	r5 := &region{xmin: 8, xmax: 12, ymin: 8, ymax: 12}
	r6 := &region{xmin: 12, xmax: 16, ymin: 8, ymax: 12}
	r1.translate = func(x, y int, f facing) (int, int, facing) {
		var nextX, nextY int
		var nextF facing
		switch f {
		case RIGHT:
			nextX = r6.xmax - 1
			nextY = r6.ymax - 1 - (y - r1.ymin)
			nextF = LEFT
		case DOWN:
			nextX = x
			nextY = r4.ymin
			nextF = DOWN
		case LEFT:
			nextX = r3.xmin + (x - r1.xmin)
			nextY = r3.ymin
			nextF = DOWN
		case UP:
			nextX = r2.xmax - 1 - (x - r1.xmin)
			nextY = r2.ymin
			nextF = DOWN
		}

		return nextX, nextY, nextF
	}
	r2.translate = func(x, y int, f facing) (int, int, facing) {
		var nextX, nextY int
		var nextF facing
		switch f {
		case RIGHT:
			nextX = r3.xmin
			nextY = y
			nextF = RIGHT
		case DOWN:
			nextX = r5.xmax - 1 - (x - r2.xmin)
			nextY = r5.ymax - 1
			nextF = UP
		case LEFT:
			nextX = r6.xmax - 1 - (y - r2.ymin)
			nextY = r6.ymax - 1
			nextF = UP
		case UP:
			nextX = r1.xmax - 1 - (x - r2.xmin)
			nextY = r1.ymin
			nextF = DOWN
		}

		return nextX, nextY, nextF
	}
	r3.translate = func(x, y int, f facing) (int, int, facing) {
		var nextX, nextY int
		var nextF facing
		switch f {
		case RIGHT:
			nextX = r4.xmin
			nextY = y
			nextF = RIGHT
		case DOWN:
			nextX = r5.xmin
			nextY = r5.ymax - 1 - (x - r3.xmin)
			nextF = RIGHT
		case LEFT:
			nextX = r2.xmax - 1
			nextY = y
			nextF = LEFT
		case UP:
			nextX = r1.xmin
			nextY = r1.ymin + (x - r3.xmin)
			nextF = RIGHT
		}

		return nextX, nextY, nextF
	}
	r4.translate = func(x, y int, f facing) (int, int, facing) {
		var nextX, nextY int
		var nextF facing
		switch f {
		case RIGHT:
			nextX = r6.xmax - 1 - (y - r4.ymin)
			nextY = r6.ymin
			nextF = DOWN
		case DOWN:
			nextX = x
			nextY = r5.ymin
			nextF = DOWN
		case LEFT:
			nextX = r3.xmax - 1
			nextY = y
			nextF = LEFT
		case UP:
			nextX = x
			nextY = r1.ymax - 1
			nextF = UP
		}

		return nextX, nextY, nextF
	}
	r5.translate = func(x, y int, f facing) (int, int, facing) {
		var nextX, nextY int
		var nextF facing
		switch f {
		case RIGHT:
			nextX = r6.xmin
			nextY = y
			nextF = RIGHT
		case DOWN:
			nextX = r2.xmax - 1 - (x - r5.xmin)
			nextY = r2.ymax - 1
			nextF = UP
		case LEFT:
			nextX = r3.xmax - 1 - (y - r5.ymin)
			nextY = r3.ymax - 1
			nextF = UP
		case UP:
			nextX = x
			nextY = r4.ymax - 1
			nextF = UP
		}

		return nextX, nextY, nextF
	}

	r6.translate = func(x, y int, f facing) (int, int, facing) {
		var nextX, nextY int
		var nextF facing
		switch f {
		case RIGHT:
			nextX = r1.xmax - 1
			nextY = y
			nextF = LEFT
		case DOWN:
			nextX = r2.xmin
			nextY = r2.ymax - 1 - (x - r6.xmin)
			nextF = RIGHT
		case LEFT:
			nextX = r5.xmax - 1
			nextY = y
			nextF = LEFT
		case UP:
			nextX = r4.xmax - 1
			nextY = r4.ymax - 1 - (x - r6.xmin)
			nextF = LEFT
		}

		return nextX, nextY, nextF
	}

	return []*region{r1, r2, r3, r4, r5, r6}

}

func inputRegions() []*region {
	r1 := &region{xmin: 50, xmax: 100, ymin: 0, ymax: 50}
	r2 := &region{xmin: 100, xmax: 150, ymin: 0, ymax: 50}
	r3 := &region{xmin: 50, xmax: 100, ymin: 50, ymax: 100}
	r4 := &region{xmin: 0, xmax: 50, ymin: 100, ymax: 150}
	r5 := &region{xmin: 50, xmax: 100, ymin: 100, ymax: 150}
	r6 := &region{xmin: 0, xmax: 50, ymin: 150, ymax: 200}
	r1.translate = func(x, y int, f facing) (int, int, facing) {
		var nextX, nextY int
		var nextF facing
		switch f {
		case RIGHT:
			nextX = r2.xmin
			nextY = y
			nextF = RIGHT
		case DOWN:
			nextX = x
			nextY = r3.ymin
			nextF = DOWN
		case LEFT:
			nextX = r4.xmin
			nextY = r4.ymax - 1 - (y - r1.ymin)
			nextF = RIGHT
		case UP:
			nextX = r6.xmin
			nextY = r6.ymin + (x - r1.xmin)
			nextF = RIGHT
		}

		return nextX, nextY, nextF
	}
	r2.translate = func(x, y int, f facing) (int, int, facing) {
		var nextX, nextY int
		var nextF facing
		switch f {
		case RIGHT:
			nextX = r5.xmax - 1
			nextY = r5.ymax - 1 - (y - r2.ymin)
			nextF = LEFT
		case DOWN:
			nextX = r3.xmax - 1
			nextY = r3.ymin + (x - r2.xmin)
			nextF = LEFT
		case LEFT:
			nextX = r1.xmax - 1
			nextY = y
			nextF = LEFT
		case UP:
			nextX = r6.xmin + (x - r2.xmin)
			nextY = r6.ymax - 1
			nextF = UP
		}

		return nextX, nextY, nextF
	}
	r3.translate = func(x, y int, f facing) (int, int, facing) {
		var nextX, nextY int
		var nextF facing
		switch f {
		case RIGHT:
			nextX = r2.xmin + (y - r3.ymin)
			nextY = r2.ymax - 1
			nextF = UP
		case DOWN:
			nextX = x
			nextY = r5.ymin
			nextF = DOWN
		case LEFT:
			nextX = r4.xmin + (y - r3.ymin)
			nextY = r4.ymin
			nextF = DOWN
		case UP:
			nextX = x
			nextY = r1.ymax - 1
			nextF = UP
		}

		return nextX, nextY, nextF
	}
	r4.translate = func(x, y int, f facing) (int, int, facing) {
		var nextX, nextY int
		var nextF facing
		switch f {
		case RIGHT:
			nextX = r5.xmin
			nextY = y
			nextF = RIGHT
		case DOWN:
			nextX = x
			nextY = r6.ymin
			nextF = DOWN
		case LEFT:
			nextX = r1.xmin
			nextY = r1.ymax - 1 - (y - r4.ymin)
			nextF = RIGHT
		case UP:
			nextX = r3.xmin
			nextY = r3.ymin + (x - r4.xmin)
			nextF = RIGHT
		}

		return nextX, nextY, nextF
	}
	r5.translate = func(x, y int, f facing) (int, int, facing) {
		var nextX, nextY int
		var nextF facing
		switch f {
		case RIGHT:
			nextX = r2.xmax - 1
			nextY = r2.ymax - 1 - (y - r5.ymin)
			nextF = LEFT
		case DOWN:
			nextX = r6.xmax - 1
			nextY = r6.ymin + (x - r5.xmin)
			nextF = LEFT
		case LEFT:
			nextX = r4.xmax - 1
			nextY = y
			nextF = LEFT
		case UP:
			nextX = x
			nextY = r3.ymax - 1
			nextF = UP
		}

		return nextX, nextY, nextF
	}

	r6.translate = func(x, y int, f facing) (int, int, facing) {
		var nextX, nextY int
		var nextF facing
		switch f {
		case RIGHT:
			nextX = r5.xmin + (y - r6.ymin)
			nextY = r5.ymax - 1
			nextF = UP
		case DOWN:
			nextX = r2.xmin + (x - r6.xmin)
			nextY = r2.ymin
			nextF = DOWN
		case LEFT:
			nextX = r1.xmin + (y - r6.ymin)
			nextY = r1.ymin
			nextF = DOWN
		case UP:
			nextX = x
			nextY = r4.ymax - 1
			nextF = UP
		}

		return nextX, nextY, nextF
	}

	return []*region{r1, r2, r3, r4, r5, r6}
}
