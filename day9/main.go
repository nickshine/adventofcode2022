package main

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type knot struct {
	x int
	y int
}

func (s *knot) equal(k *knot) bool {
	return s.x == k.x && s.y == k.y
}

func (s *knot) adjacent(k *knot) bool {
	dx := abs(k.x - s.x)
	dy := abs(k.y - s.y)

	return dx == 1 && dy == 1 ||
		dx == 1 && dy == 0 ||
		dx == 0 && dy == 1
}

func travel(h, t *knot, dx, dy, n int, seen map[knot]int) {

	for i := 0; i < n; i++ {

		h.x, h.y = h.x+dx, h.y+dy
		// fmt.Printf("H: %d,%d\n", h.x, h.y)
		if t.adjacent(h) || t.equal(h) {
			// t doesn't move
			continue
		}

		// if not adjacent/equal, t moves in same direction as head (dx,dy)
		t.x, t.y = t.x+dx, t.y+dy

		// if h travels L or R && h.y != t.y, t needs to move U or D to h
		if dx != 0 && h.y != t.y {
			t.y = h.y
		}

		// if h travels U or D && h.x != t.x, t needs to move L or R to h
		if dy != 0 && h.x != t.x {
			t.x = h.x
		}

		seen[*t]++
	}
}

func part1(in string) int {
	lines := strings.Split(in, "\n")

	seen := map[knot]int{}

	h, t := &knot{0, 0}, &knot{0, 0}
	seen[*t]++

	for _, l := range lines {
		parts := strings.Fields(l)
		dir := parts[0]
		n, err := strconv.Atoi(parts[1])
		if err != nil {
			panic(err)
		}

		// cases in which T doesn't move
		// - if H coord == T coord (H overlapped T)
		// - if H coord stays within surrounding T coords

		// case in which T moves diagonally
		// - if H and T aren't in same row or column and not touching, T will move diagonally

		switch dir {
		case "L":
			travel(h, t, -1, 0, n, seen)
		case "R":
			travel(h, t, 1, 0, n, seen)
		case "U":
			travel(h, t, 0, 1, n, seen)
		case "D":
			travel(h, t, 0, -1, n, seen)
		default:
			panic("invalid input")
		}

	}

	return len(seen)
}

func main() {
	fmt.Printf("part 1: %d\n", part1(strings.TrimSpace(input)))
	// fmt.Printf("part 2: %d\n", part2(readGrid(input)))
}
