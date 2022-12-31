package main

import (
	_ "embed"
	"fmt"
	"log"
	"strings"
)

//go:embed example.txt
var exampleInput string

//go:embed input.txt
var input string

var snafuMap = map[byte]int{'2': 2, '1': 1, '0': 0, '-': -1, '=': -2}
var decMap = map[int]byte{2: '2', 1: '1', 0: '0', -1: '-', -2: '='}

func toDecimal(snafu string) int {
	d := 0
	for i, place := len(snafu)-1, 1; i >= 0; i, place = i-1, place*5 {
		d += snafuMap[snafu[i]] * place
	}
	return d
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func toSnafu(d int) string {
	var b strings.Builder

	// find start place
	var min, max int
	place := 1
	max = 2*place + place/2

	for abs(d) > max {
		place *= 5
		max = 2*place + place/2
	}

	// determine digit for each place to 0
	for ; place > 0; place /= 5 {
		sign := 1
		if d < 0 {
			sign = -1
		}

		min = place - place/2   // place - place/2
		max = 2*place + place/2 // 2*place + place/2

		switch {
		case abs(d) > max-place: // digit is 2 or -2
			d -= 2 * place * sign
			b.WriteByte(decMap[2*sign])
		case abs(d) >= min: // digit is 1 or -1
			d -= place * sign
			b.WriteByte(decMap[sign])
		default: // digit is 0
			b.WriteByte(decMap[0])
		}
	}

	return b.String()
}

func toSnafuBetter(d int) string {

	var s string

	const base = 5

	for d > 0 {
		c := byte('0')
		switch d % base {
		case 0:
			c = decMap[0]
		case 1:
			c = decMap[1]
		case 2:
			c = decMap[2]
		case 3:
			c = decMap[-2]
			d += base
		case 4:
			c = decMap[-1]
			d += base
		default:
			panic("invalid")
		}

		s = string(c) + s
		d /= base
	}

	return s
}

func part1(in string) string {
	snafus := strings.Fields(strings.TrimSpace(in))

	total := 0
	for _, s := range snafus {
		total += toDecimal(s)
	}

	log.Printf("total: %d", total)
	log.Printf("snafu traditional: %s", toSnafuBetter(total))
	return toSnafu(total)
}

func main() {
	fmt.Printf("Part 1 example: %s\n", part1(exampleInput))
	fmt.Printf("Part 1: %s\n", part1(input))
}
