package main

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

type stack struct {
	v []string
}

func (s *stack) pop(n int) []string {
	items := s.v[len(s.v)-n:]
	s.v = s.v[:len(s.v)-n]
	return items
}

func (s *stack) push(items ...string) {
	s.v = append(s.v, items...)
}

// returns a slice of stacks
func parseStacks(in string) []stack {
	in = strings.TrimRight(in, " ")

	// total stack count is last char of input
	stackCount, err := strconv.Atoi(string(in[len(in)-1]))
	if err != nil {
		panic(err)
	}

	stacks := make([]stack, stackCount)
	lines := strings.Split(in, "\n")
	// start from bottom of stacks
	for i := len(lines) - 2; i >= 0; i-- {
		l := lines[i]
		// scan each crate
		for j, idx := 1, 0; j < len(l); j, idx = j+4, idx+1 {
			crate := string(l[j])
			if crate != " " {
				stacks[idx].push(crate)
			}
		}
	}

	return stacks
}

func parseSteps(in string) [][]int {
	lines := strings.Split(in, "\n")
	steps := make([][]int, len(lines))
	for i, line := range lines {
		l := strings.Split(line, " ")
		s1, _ := strconv.Atoi(l[1])
		s2, _ := strconv.Atoi(l[3])
		s3, _ := strconv.Atoi(l[5])
		steps[i] = []int{s1, s2, s3}

	}
	return steps
}

func tops(s []stack) string {
	var top []string
	for _, stack := range s {
		top = append(top, stack.pop(1)...)
	}
	return strings.Join(top, "")
}

func part1() string {
	lines := strings.Split(strings.TrimRight(input, "\n"), "\n\n")
	stacks := parseStacks(lines[0])
	steps := parseSteps(lines[1])

	for _, step := range steps {
		moves, from, to := step[0], step[1], step[2]

		for i := 0; i < moves; i++ {
			crate := stacks[from-1].pop(1)
			stacks[to-1].push(crate...)
		}
	}

	return tops(stacks)
}

func part2() string {
	lines := strings.Split(strings.TrimRight(input, "\n"), "\n\n")
	stacks := parseStacks(lines[0])
	steps := parseSteps(lines[1])

	for _, step := range steps {
		moves, from, to := step[0], step[1], step[2]

		crates := stacks[from-1].pop(moves)
		stacks[to-1].push(crates...)
	}

	return tops(stacks)
}

func main() {
	fmt.Println(part1())
	fmt.Println(part2())
}
