package main

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

type knot struct {
	x int
	y int
}

func checkCycle(cycle, x int, cycles map[int]int) {
	switch cycle {
	case 20, 60, 100, 140, 180, 220:
		cycles[cycle] = x
	}
}

func display(cycle, x int) {

	pos := cycle % 40

	// if pos == 0 {
	// 	fmt.Printf("\n")
	// }

	switch pos {
	case x - 1, x, x + 1:
		fmt.Printf("#")
	default:
		fmt.Printf(".")
	}

	if pos == 39 {
		fmt.Printf("\n")
	}

}

func main() {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	cycles := map[int]int{}

	x, cycle := 1, 1
	display(0, x)

	for _, l := range lines {
		parts := strings.Fields(l)

		display(cycle, x)
		cycle++
		checkCycle(cycle, x, cycles)

		switch parts[0] {
		case "noop":
			// noop
		case "addx":
			val, err := strconv.Atoi(parts[1])
			if err != nil {
				panic(err)
			}

			x += val
			display(cycle, x)
			cycle++
			checkCycle(cycle, x, cycles)
		default:
			panic("bad instruction")
		}
	}

	fmt.Println()

	var sum int
	for c, x := range cycles {
		sum += c * x
	}

	fmt.Printf("Part 1: %d\n", sum)
}
