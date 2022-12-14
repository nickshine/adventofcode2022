package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

//go:embed example.txt
var exampleInput string

//go:embed input.txt
var input string

type comparison int

const (
	ORDERED comparison = iota
	UNORDERED
	EQUAL
)

func parsePairs(in string) [][]any {
	pairs := [][]any{}

	in = strings.TrimSpace(in)
	rawPairs := strings.Split(in, "\n\n")

	for _, p := range rawPairs {
		pair := strings.Split(p, "\n")
		var left, right any
		json.Unmarshal([]byte(pair[0]), &left)
		json.Unmarshal([]byte(pair[1]), &right)

		pairs = append(pairs, []any{left, right})
	}

	return pairs
}

func bothLists(left, right any) bool {
	_, lok := left.([]any)
	_, rok := right.([]any)
	return lok && rok
}

func bothDigits(left, right any) bool {
	_, lok := left.(float64)
	_, rok := right.(float64)
	return lok && rok
}

func compare(left, right any) comparison {

	switch {
	case bothLists(left, right):
		l, r := left.([]any), right.([]any)
		for i, v := range l {
			if i >= len(r) {
				return UNORDERED
			}
			cmp := compare(v, r[i])
			if cmp != EQUAL {
				return cmp
			}
		}

		if len(l) < len(r) {
			return ORDERED
		} else if len(l) > len(r) {
			return UNORDERED
		}
		return EQUAL
	case bothDigits(left, right):
		l, r := left.(float64), right.(float64)

		if l < r {
			return ORDERED
		} else if l > r {
			return UNORDERED
		}
		return EQUAL
	default: // mixed types
		if l, ok := left.([]any); ok { // right is digit
			r := right.(float64)
			return compare(l, []any{r})
		}

		l := left.(float64)
		r := right.([]any)
		return compare([]any{l}, r)
	}

}

func part1(in string) int {
	pairs := parsePairs(in)
	var sum int

	for i, p := range pairs {
		if compare(p[0], p[1]) == ORDERED {
			sum += i + 1
		}
	}

	return sum
}

func part2(in string) int {
	pairs := parsePairs(in)
	p2 := []any{[]any{float64(2)}}
	p6 := []any{[]any{float64(6)}}

	var packets []any
	for _, p := range pairs {
		for _, pp := range p {
			packets = append(packets, pp)
		}
	}
	packets = append(packets, p2)
	packets = append(packets, p6)

	sort.Slice(packets, func(i, j int) bool {
		return compare(packets[i], packets[j]) == ORDERED
	})

	key := 1
	for i, p := range packets {
		if reflect.DeepEqual(p, p2) || reflect.DeepEqual(p, p6) {
			key *= (i + 1)

		}

	}

	return key
}

func main() {
	fmt.Printf("Part 1 example: %d\n", part1(exampleInput))
	fmt.Printf("Part 1: %d\n", part1(input))
	fmt.Printf("Part 2 example: %d\n", part2(exampleInput))
	fmt.Printf("Part 2: %d\n", part2(input))
}
