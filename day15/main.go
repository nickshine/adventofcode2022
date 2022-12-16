package main

import (
	_ "embed"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

//go:embed example.txt
var exampleInput string

//go:embed input.txt
var input string

const limit = 4000000

type point struct {
	x, y int
}

type sensor point

type span struct {
	l, r int
}

func (p point) String() string {
	return fmt.Sprintf("{%d,%d}", p.x, p.y)
}

func (s sensor) String() string {
	return fmt.Sprintf("sensor{%d,%d}", s.x, s.y)
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func (s sensor) distance(b point) int {
	return abs(s.x-b.x) + abs(s.y-b.y)
}

func parseSensors(in string) map[sensor]point {
	sensors := map[sensor]point{} // sensor to closest beacon
	lines := strings.Split(strings.TrimSpace(in), "\n")

	extract := func(s string) int {
		st := strings.Trim(s, ",:")
		parts := strings.Split(st, "=")
		if len(parts) != 2 {
			panic("invalid input")
		}

		d, err := strconv.Atoi(parts[1])
		if err != nil {
			panic(fmt.Errorf("invalid input: %w", err))
		}

		return d
	}

	for _, l := range lines {
		parts := strings.Fields(l)
		if len(parts) != 10 {
			panic("invalid input")
		}
		sx := extract(parts[2]) // x=2,
		sy := extract(parts[3]) // y=18:
		bx := extract(parts[8]) // x=-2,
		by := extract(parts[9]) // y=15

		sensors[sensor{sx, sy}] = point{bx, by}
	}

	return sensors
}

func part1(in string, y int) int {
	sensors := parseSensors(in)

	positions := map[point]struct{}{}
	for s, b := range sensors {
		d := s.distance(b)
		ymin, ymax := s.y-d, s.y+d
		if y < ymin || y > ymax { // if s's field does not overlap y, skip
			continue
		}
		// record all the positions where s's field touch row y
		radius := d - abs(y-s.y)
		for i := s.x - radius; i <= s.x+radius; i++ {
			positions[point{i, y}] = struct{}{}

		}
		if b.y == y { // if beacon on row y, delete it from positions
			delete(positions, point{b.x, b.y})
		}
	}

	return len(positions)

}

func part2(in string, min, max int) int {
	sensors := parseSensors(in)

	spans := map[int][]span{}
	for s, b := range sensors {
		d := s.distance(b)
		// for each sensor, calculate a list of x-ranges (spans) that it covers, and add to others for same y
		for y := s.y - d; y <= s.y+d; y++ {
			radius := d - abs(y-s.y)
			if _, ok := spans[y]; !ok {
				spans[y] = []span{}
			}
			spans[y] = append(spans[y], span{s.x - radius, s.x + radius})
		}
	}

	for y, xspans := range spans {
		if y < 0 || y > max {
			continue
		}

		sort.Slice(xspans, func(i, j int) bool {
			return xspans[i].l < xspans[j].l
		})

		cmax := xspans[0].r
		for _, span := range xspans {
			if cmax >= span.l-1 { // no gaps between spans
				if span.r >= cmax {
					cmax = span.r
				}
			} else { // non-covered point, aka hidden beacon

				return (cmax+1)*limit + y
			}
		}

	}

	return -1
}

func main() {
	fmt.Printf("Part 1 example: %d\n", part1(exampleInput, 10))
	fmt.Printf("Part 1: %d\n", part1(input, 2000000))
	fmt.Printf("Part 2 example: %d\n", part2(exampleInput, 0, 20))
	fmt.Printf("Part 2: %d\n", part2(input, 0, limit))
}

func part2Slow(in string, min, max int) int {
	sensors := parseSensors(in)
	positions := map[point]struct{}{}

	filtered := map[sensor]point{}

	for s, b := range sensors {
		d := s.distance(b)
		// if s's field does not overlap distress range
		ymin, ymax := s.y-d, s.y+d
		xmin, xmax := s.x-d, s.x+d

		// if ymin > max
		if ymin > max || ymax < min || xmin > max || xmax < min { // if s's field does not overlap distress range, skip
			continue
		}

		filtered[s] = b
	}

	for y := min; y <= max; y++ {
		for x := min; x <= max; x++ {
			// check if point within distance of s
			for s, b := range filtered {
				p := point{x, y}
				if s.distance(p) <= s.distance(b) {
					positions[p] = struct{}{}
				}
			}
		}
	}

	var xp, yp int
	for y := min; y <= max; y++ {
		for x := min; x <= max; x++ {
			if _, ok := positions[point{x, y}]; !ok {
				xp, yp = x, y
			}
		}
	}

	return xp*limit + yp
}
