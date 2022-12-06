package main

import (
	_ "embed"
	"fmt"
	"strings"
)

//go:embed input.txt
var input string

func scan(size int) int {
	in := strings.TrimSpace(input)
	l, r := 0, 1
	seen := make(map[byte]int, size)
	seen[in[l]] = l

	for {
		if r == len(in) {
			return -1
		} else if r-l == size {
			break
		}

		if idx, ok := seen[in[r]]; ok && idx >= l {
			l = idx + 1
		}

		seen[in[r]] = r
		r++
	}

	return r
}

func main() {
	fmt.Printf("part 1: %d\n", scan(4))
	fmt.Printf("part 2: %d\n", scan(14))
}
