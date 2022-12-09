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

func travel2(rope []*knot, dx, dy, n int, seen map[knot]int) {

	h, t := rope[0], rope[len(rope)-1]
	for m := 0; m < n; m++ {
		h.x, h.y = h.x+dx, h.y+dy

		for i := 1; i < len(rope); i++ {
			leader, follower := rope[i-1], rope[i]

			xd := leader.x - follower.x
			yd := leader.y - follower.y

			if yd == 2 && xd >= 1 || yd >= 1 && xd == 2 { // follower moves U and R
				follower.x, follower.y = follower.x+1, follower.y+1
			} else if yd >= 1 && xd == -2 || yd == 2 && xd <= -1 { // follower moves U and L
				follower.x, follower.y = follower.x-1, follower.y+1
			} else if yd == -2 && xd <= -1 || yd <= -1 && xd == -2 { // follower moves D and L
				follower.x, follower.y = follower.x-1, follower.y-1
			} else if yd == -2 && xd >= 1 || yd <= -1 && xd == 2 { // follower moves D and R
				follower.x, follower.y = follower.x+1, follower.y-1
			} else if xd == 0 && yd > 1 { // follower moves U
				follower.y++
			} else if xd == 0 && yd < -1 { // follower moves D
				follower.y--
			} else if xd > 1 && yd == 0 { // follower moves R
				follower.x++
			} else if xd < -1 && yd == 0 { // follower moves L
				follower.x--
			}
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

func part2(in string) int {
	lines := strings.Split(in, "\n")

	rope := make([]*knot, 10)
	for i := 0; i < len(rope); i++ {
		rope[i] = &knot{0, 0}
	}
	seen := map[knot]int{}
	seen[knot{0, 0}]++

	for _, l := range lines {
		parts := strings.Fields(l)
		dir := parts[0]
		n, err := strconv.Atoi(parts[1])
		if err != nil {
			panic(err)
		}

		switch dir {
		case "L":
			travel2(rope, -1, 0, n, seen)
		case "R":
			travel2(rope, 1, 0, n, seen)
		case "U":
			travel2(rope, 0, 1, n, seen)
		case "D":
			travel2(rope, 0, -1, n, seen)
		default:
			panic("invalid input")
		}
	}

	return len(seen)

}

func main() {
	fmt.Printf("part 1: %d\n", part1(strings.TrimSpace(input)))
	fmt.Printf("part 2: %d\n", part2(strings.TrimSpace(input)))
}
