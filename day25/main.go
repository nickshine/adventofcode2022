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
	max = place * 5 / 2 // 2*place + place/2

	for abs(d) > max {
		place *= 5
		max = place * 5 / 2
	}

	// determine digit for each place to 0
	for ; place > 0; place /= 5 {
		sign := 1
		if d < 0 {
			sign = -1
		}

		min = place * 1 / 2 // place - place/2
		max = place * 5 / 2 // 2*place + place/2

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

func part1(in string) string {
	snafus := strings.Fields(strings.TrimSpace(in))

	total := 0
	for _, s := range snafus {
		total += toDecimal(s)
	}

	log.Printf("total: %d", total)
	return toSnafu(total)
}

func main() {
	fmt.Printf("Part 1 example: %s\n", part1(exampleInput))
	fmt.Printf("Part 1: %s\n", part1(input))
}

// rules:

// 5^n == 2*(2*5^n-1 + 2*5^n-2 ... + 2*5^n-n) + 1

// max for 5^n == 5^n + (5^n)/2
//     max 5^4 == 5^4 + (5^4)/2
//     max 625 == 625 + 312 == 937
// snafu: 1222 == 1*5^4 + 2*5^3 + 2*5^2 + 2*5^1 + 2*5^0

// min for 5^n == 5^n - (5^n)/2
//     min 5^4 == 5^4 - 5^4/2
//     min 625 == 625 - 312 == 313
//                625 - 2*125 - 2*25 - 2*5 - 2*1

// if 5^n - 5^n/2 > d
// if 5^4 - 5^4/2 > d
// if 625 - 312 > d
// if 313 > d

// if 5^n + 5^n/2 < d
// if 5^4 + 5^4/2 < d
// if 625 + 312 < d
// if 937 < d
// if 2*5^n + 5^n/2 < d
// if 625+937 < d
// if 1562 < d

// 5^5 place:
// if d within 1*5^n - 5^n/2 && 2*5^n+5^n/2
// if d within 1*5^5 - 5^5/2 && 2*5^5+5^5/2
// if d within  3125-1562 && 6250+1562
// if d within  1563 && 7812

// 5^4 place:
// if d within 1*5^n - 5^n/2 && 2*5^n+5^n/2
// if d within 1*5^4 - 5^4/2 && 2*5^4+5^4/2
// if d within  625-312 && 1250+312
// if d within  313 && 1562

// 5^3 place:
// if d within 5^3 - 5^3/2 && 2*5^3+5^3/2
// if d within 125 - 62 && 250+62
// if d within 63 && 312
// 1 or 2?
// if > 1*5^3+5^3/2, 2
// if > 1*125+62
// if > 187, 2
// if > max-5^3 == 312-125 == 187, 2

// 5^2 place:
// if d within 5^2 - 5^2/2 && 2*5^2+5^2/2
// if d within 25 - 12 && 50+12
// if d within 13 && 62

// 5^1 place:
// if d within 5^1 - 5^1/2 && 2*5^1+5^1/2
// if d within 5 - 2 && 10+2
// if d within 3 && 12

// 5^0 place:
// if d within 5^0 - 5^0/2 && 2*5^0+5^0/2
// if d within 1 - 0 && 2+0
// if d within 1 && 2
