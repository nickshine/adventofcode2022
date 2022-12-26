package main

import (
	_ "embed"
	"fmt"
	"log"
	"strconv"
	"strings"
)

//go:embed example.txt
var exampleInput string

//go:embed input.txt
var input string

type job func(a, b int) int
type monkey struct {
	name     string
	monkeys  [2]string
	lOp, rOp *int
	job      job
}

func (m *monkey) String() string {
	l, r := "nil", "nil"

	if m.lOp != nil {
		l = fmt.Sprintf("%d", *m.lOp)
	}
	if m.rOp != nil {
		r = fmt.Sprintf("%d", *m.rOp)
	}

	return fmt.Sprintf("{name:%s, lOp:%s, rOp:%s}", m.name, l, r)
}

func (m *monkey) resolve(simple map[string]int) bool {
	if len(m.monkeys) != 2 {
		panic("invalid setup")
	}

	if v, ok := simple[m.monkeys[0]]; ok {
		if m.lOp == nil {
			m.lOp = &v
		}
	}

	if v, ok := simple[m.monkeys[1]]; ok {
		m.rOp = &v
	}

	if m.lOp != nil && m.rOp != nil {
		return true
	}

	return false
}

func (m *monkey) doJob() (int, error) {
	if m.lOp == nil || m.rOp == nil {
		return 0, fmt.Errorf("must resolve operands before doing job")
	}

	return m.job(*m.lOp, *m.rOp), nil
}

func parseInput(in string, part2 bool) (map[string]int, []*monkey) {
	lines := strings.Split(strings.TrimSpace(in), "\n")

	simpleMonkeys := map[string]int{}
	var complexMonkeys []*monkey

	for _, l := range lines {
		parts := strings.Fields(l)
		if len(parts) == 2 {
			name := strings.TrimRight(parts[0], ":")
			num, err := strconv.Atoi(parts[1])
			if err != nil {
				panic("invalid input")
			}
			simpleMonkeys[name] = num
		} else if len(parts) == 4 {
			name := strings.TrimRight(parts[0], ":")
			var job job
			if part2 {
				if name == "root" {
					job = func(a, b int) int {
						log.Printf("a should match b for root: %d,%d", a, b)
						return a - b
					}
					complexMonkeys = append(complexMonkeys, &monkey{name: name, monkeys: [2]string{parts[1], parts[3]}, job: job})
					continue
				}
			}
			switch parts[2] {
			case "+":
				job = func(a, b int) int { return a + b }
			case "-":
				job = func(a, b int) int { return a - b }
			case "*":
				job = func(a, b int) int { return a * b }
			case "/":
				job = func(a, b int) int { return a / b }
			default:
				panic("invalid input")
			}

			complexMonkeys = append(complexMonkeys, &monkey{name: name, monkeys: [2]string{parts[1], parts[3]}, job: job})
		} else {
			panic("invalid input")
		}
	}

	return simpleMonkeys, complexMonkeys
}

func copyMonkeys(simple map[string]int, complex []*monkey) (map[string]int, []*monkey) {
	newSimple := make(map[string]int, len(simple))
	for k, v := range simple {
		newSimple[k] = v
	}
	newComplex := make([]*monkey, len(complex))
	for i, v := range complex {
		m := *v
		newComplex[i] = &m
	}

	return newSimple, newComplex
}

func part1(in string) int {
	simpleMonkeys, complexMonkeys := parseInput(in, false)

	for len(complexMonkeys) != 0 {
		monkey := complexMonkeys[0]
		complexMonkeys = complexMonkeys[1:] // dequeue
		if monkey.resolve(simpleMonkeys) {
			res, err := monkey.doJob()
			if err != nil {
				panic(err)
			}
			simpleMonkeys[monkey.name] = res
		} else {
			complexMonkeys = append(complexMonkeys, monkey) // enqueue
		}
	}

	return simpleMonkeys["root"]
}

func part2(in string, start, end int) int {
	simple, complex := parseInput(in, true)

	for i := start; i < end; i++ {
		simpleMonkeys, complexMonkeys := copyMonkeys(simple, complex)

		simpleMonkeys["humn"] = i

		for len(complexMonkeys) != 0 {
			monkey := complexMonkeys[0]
			complexMonkeys = complexMonkeys[1:] // dequeue
			if monkey.resolve(simpleMonkeys) {
				res, err := monkey.doJob()
				if err != nil {
					panic(err)
				}
				simpleMonkeys[monkey.name] = res
			} else {
				complexMonkeys = append(complexMonkeys, monkey) // enqueue
			}
		}

		if simpleMonkeys["root"] == 0 { // operands match
			return simpleMonkeys["humn"]
		}

	}
	return -1
}

func main() {
	fmt.Printf("Part 1 example: %d\n", part1(exampleInput))
	fmt.Printf("Part 1: %d\n", part1(input))
	fmt.Printf("Part 2 example: %d\n", part2(exampleInput, 300, 302))
	fmt.Printf("Part 2: %d\n", part2(input, 3099532690000, 3099532700000))
}
