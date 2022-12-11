package main

import (
	_ "embed"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

type apply func(l, r int) int
type op func(n int) int

type monkey struct {
	items          []int
	count, divisor int
	op             op
	test           func(n int) bool
	tID, fID       int
}

func parseOperation(left, operator, right string) (out op) {
	var fn apply
	var leftOld, rightOld bool

	l, err := strconv.Atoi(left)
	if err != nil {
		leftOld = true
	}

	r, err := strconv.Atoi(right)
	if err != nil {
		rightOld = true
	}

	switch operator {
	case "*":
		fn = func(l, r int) int {
			return l * r
		}

	case "+":
		fn = func(l, r int) int {
			return l + r
		}
	case "-":
		fn = func(l, r int) int {
			return l - r
		}
	}

	switch {
	case leftOld && rightOld:
		out = func(n int) int {
			return fn(n, n)
		}
	case leftOld:
		out = func(n int) int {
			return fn(n, r)
		}
	case rightOld:
		out = func(n int) int {
			return fn(l, n)
		}
	default:
		out = func(n int) int {
			return fn(l, r)
		}
	}

	return out

}

func readMonkeys(in []string) []*monkey {
	monkeys := []*monkey{}

	for _, m := range in {
		monkey := &monkey{}
		for _, l := range strings.Split(strings.TrimSpace(m), "\n") {
			parts := strings.Fields(l)
			switch parts[0] {
			case "Monkey":
			case "Starting":
				for _, s := range parts[2:] {
					item, err := strconv.Atoi(strings.Trim(s, ","))
					if err != nil {
						panic(err)
					}
					monkey.items = append(monkey.items, item)
				}
			case "Operation:":
				left := parts[len(parts)-3]
				operator := parts[len(parts)-2]
				right := parts[len(parts)-1]

				monkey.op = parseOperation(left, operator, right)
			case "Test:":
				if len(parts) != 4 {
					panic("invalid test")
				}
				divisor, err := strconv.Atoi(parts[3])
				if err != nil {
					panic(err)
				}

				monkey.divisor = divisor
				monkey.test = func(n int) bool {
					return n%divisor == 0
				}

			default:
				if len(parts) != 6 {
					panic("invalid ifs")
				}
				id, err := strconv.Atoi(parts[5])
				if err != nil {
					panic(err)
				}

				switch parts[1] {
				case "true:":
					monkey.tID = id
				case "false:":
					monkey.fID = id
				}
			}
		}
		monkeys = append(monkeys, monkey)
	}

	return monkeys
}

func throw(m *monkey, item int) {
	m.items = append(m.items, item)
}

func findCommon(monkeys []*monkey) int {
	common := 1
	for _, m := range monkeys {
		common *= m.divisor
	}

	return common
}

func doRound(monkeys []*monkey, manage func(int) int) {

	for _, m := range monkeys {
		for _, item := range m.items {
			m.count++
			wl := manage(m.op(item))
			if m.test(wl) {
				throw(monkeys[m.tID], wl)
			} else {
				throw(monkeys[m.fID], wl)
			}
		}
		// clear items
		m.items = []int{}
	}
}

func maxLevel(monkeys []*monkey) int {
	counts := []int{}
	for _, m := range monkeys {
		counts = append(counts, m.count)
	}

	sort.Ints(counts)
	return counts[len(counts)-1] * counts[len(counts)-2]
}

func part1(monkeys []*monkey) int {
	const rounds = 20
	for n := 0; n < rounds; n++ {
		doRound(monkeys, func(item int) int {
			return item / 3
		})
	}

	return maxLevel(monkeys)
}

func part2(monkeys []*monkey) int {
	const rounds = 10000
	divisor := findCommon(monkeys)

	for n := 0; n < rounds; n++ {
		doRound(monkeys, func(v int) int {
			return v % divisor
		})
	}

	return maxLevel(monkeys)
}

func main() {
	in := strings.Split(strings.TrimSpace(input), "\n\n")
	fmt.Printf("Part 1: %d\n", part1(readMonkeys(in)))
	fmt.Printf("Part 2: %d\n", part2(readMonkeys(in)))
}
