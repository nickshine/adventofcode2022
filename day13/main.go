package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
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

// func compareOld(left, right string) int {
//
// 	fmt.Printf("------------------------nleft: %s, right: %s\n", left, right)
//
// 	if len(left) == 0 || len(right) == 0 {
// 		if len(left) <= len(right) {
// 			return ORDERED
// 		} else {
// 			return UNORDERED
// 		}
// 	}
//
// 	cl, cr := rune(left[0]), rune(right[0])
//
// 	// fmt.Printf("cl: %c, cr: %c\n", cl, cr)
// 	switch {
// 	case bothLists(left, right):
// 		fmt.Println("splitting left and right lists")
// 		ll, rl := split(left), split(right)
//
// 		fmt.Printf("l split: %v, len: %d\n", ll, len(ll))
// 		fmt.Printf("r split: %v, len: %d\n", rl, len(rl))
// 		fmt.Printf("len(ll): %d, len(rl): %d\n", len(ll), len(rl))
// 		for i := range ll {
// 			if i >= len(rl) { // right ran out of items
// 				return UNORDERED
// 			}
// 			cmp := compare(ll[i], rl[i])
// 			if cmp != CONTINUE {
// 				return cmp
// 			}
// 		}
//
// 		if len(ll) < len(rl) {
// 			return ORDERED
// 		} else if len(ll) > len(rl) {
// 			return UNORDERED
// 		}
//
// 		return CONTINUE
//
// 	case bothDigits(cl, cr):
// 		dl, dr := cl-'0', cr-'0'
// 		fmt.Printf("comparing digits: %d,%d\n", dl, dr)
// 		if dl > dr {
// 			fmt.Printf("dl > dr: %d,%d\n", dl, dr)
// 			return UNORDERED
// 		} else if dl < dr {
// 			fmt.Printf("dl < dr: %d,%d\n", dl, dr)
// 			return ORDERED
// 		} else { // same
// 			return compare(left[1:], right[1:])
// 		}
// 	default: // one is digit
// 		if unicode.IsDigit(cl) {
// 			return compare(fmt.Sprintf("[%c]%s", cl, left[1:]), right)
// 		} else {
// 			return compare(left, fmt.Sprintf("[%c]%s", cr, right[1:]))
// 		}
//
// 	}
//
// 	// return false
// }

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

func main() {
	fmt.Printf("Part 1 example: %d\n", part1(exampleInput))
	fmt.Printf("Part 1: %d\n", part1(input))
	// fmt.Printf("Part 2 example: %d\n", part2(exampleInput))
	// fmt.Printf("Part 2: %d\n", part2(input))
}
