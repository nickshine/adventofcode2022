package main

import (
	_ "embed"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

//go:embed example.txt
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

func part1() int {
	rawMonkeys := strings.Split(strings.TrimSpace(input), "\n\n")
	monkeys := readMonkeys(rawMonkeys)

	rounds := 20

	for n := 0; n < rounds; n++ {

		for _, m := range monkeys {
			for _, item := range m.items {
				m.count++
				wl := m.op(item) / 3
				if m.test(wl) {
					// fmt.Printf("Monkey %d throwing item %d to %d\n", i, wl, m.tID)
					throw(monkeys[m.tID], wl)
				} else {
					// fmt.Printf("Monkey %d throwing item %d to %d\n", i, wl, m.fID)
					throw(monkeys[m.fID], wl)
				}
			}
			// clear items
			m.items = []int{}
		}

	}

	counts := []int{}
	for i, m := range monkeys {
		counts = append(counts, m.count)
		fmt.Printf("monkey %d: %v, count: %d\n", i, m.items, m.count)
	}

	sort.Ints(counts)

	fmt.Printf("Max1: %d, max2: %d\n", counts[len(counts)-1], counts[len(counts)-2])

	return counts[len(counts)-1] * counts[len(counts)-2]
}

func main() {
	fmt.Printf("Part 1: %d\n", part1())
}
