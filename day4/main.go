package main

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

// parsePair return low and hi for each in pair
func parsePair(in string) (ll int, lh int, rl int, rh int) {
	pair := strings.Split(in, ",")
	left := strings.Split(pair[0], "-")
	right := strings.Split(pair[1], "-")

	ll, _ = strconv.Atoi(left[0])
	lh, _ = strconv.Atoi(left[1])
	rl, _ = strconv.Atoi(right[0])
	rh, _ = strconv.Atoi(right[1])

	return
}

func contained(in string) bool {

	ll, lh, rl, rh := parsePair(in)

	// left contains right, or right contains left
	return ll >= rl && lh <= rh || rl >= ll && rh <= lh
}

func overlapped(in string) bool {

	ll, lh, rl, rh := parsePair(in)

	// if left high overlaps right, or right high overlaps left
	return lh >= rl && lh <= rh || rh >= ll && rh <= lh
}

func part1() int {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var count int
	for _, l := range lines {
		if contained(l) {
			count++
		}
	}

	return count
}

func part2() int {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var count int
	for _, l := range lines {
		if overlapped(l) {
			count++
		}
	}

	return count
}

func main() {
	fmt.Println(part1())
	fmt.Println(part2())
}
