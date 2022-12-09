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

func move(h, t *knot) {

	xd := h.x - t.x
	yd := h.y - t.y

	switch {
	case abs(xd) < 2 && abs(yd) < 2: // t equal or adjacent to h
		return
	case xd == 0 && yd != 0: // t moves U or D
		if yd > 1 {
			t.y++
		} else {
			t.y--
		}
	case yd == 0 && xd != 0: // t moves L or R
		if xd > 1 {
			t.x++
		} else {
			t.x--
		}

	default: // t moves diagonal
		if yd > 0 {
			t.y++
		} else {
			t.y--
		}
		if xd > 0 {
			t.x++
		} else {
			t.x--
		}

	}
}

func travel(rope []*knot, dx, dy, n int, seen map[knot]bool) {

	h, t := rope[0], rope[len(rope)-1]
	for m := 0; m < n; m++ {
		h.x, h.y = h.x+dx, h.y+dy

		for i := 1; i < len(rope); i++ {
			move(rope[i-1], rope[i])
		}
		seen[*t] = true
	}
}

func run(in string, size int) int {
	lines := strings.Split(in, "\n")

	rope := make([]*knot, size)
	for i := 0; i < len(rope); i++ {
		rope[i] = &knot{0, 0}
	}
	seen := map[knot]bool{{0, 0}: true}

	for _, l := range lines {
		parts := strings.Fields(l)
		dir := parts[0]
		n, err := strconv.Atoi(parts[1])
		if err != nil {
			panic(err)
		}

		switch dir {
		case "L":
			travel(rope, -1, 0, n, seen)
		case "R":
			travel(rope, 1, 0, n, seen)
		case "U":
			travel(rope, 0, 1, n, seen)
		case "D":
			travel(rope, 0, -1, n, seen)
		default:
			panic("invalid input")
		}
	}

	return len(seen)

}

func main() {
	fmt.Printf("part 1: %d\n", run(strings.TrimSpace(input), 2))
	fmt.Printf("part 2: %d\n", run(strings.TrimSpace(input), 10))
}
