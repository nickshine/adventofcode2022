package main

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

//go:embed example.txt
var exampleInput string

//go:embed input.txt
var input string

const decryptionKey = 811589153

type node struct {
	prev, next *node
	value      int
}

type list struct {
	head, tail *node
}

func (l *list) insert(value int) *node {

	n := &node{value: value}

	if l.head == nil {
		l.head = n
		l.tail = n
		return n
	}

	p := l.head
	for p.next != l.tail.next { // l.tail.next == l.head in circular list
		p = p.next
	}
	n.prev = p
	p.next = n
	l.tail = n

	l.head.prev = l.tail // update circular refs
	l.tail.next = l.head

	return n
}

// insertNode inserts b after a
func (l *list) insertNode(a, b *node) {
	next := a.next
	a.next = b
	b.prev = a
	b.next = next
	next.prev = b

	if a == l.tail {
		l.tail = b
	}
}

// removeNode removes n from the list.
//
// It is assumed that n exists in the list.
func (l *list) removeNode(n *node) {
	n.prev.next = n.next
	n.next.prev = n.prev

	if n == l.head {
		l.head = n.next
	} else if n == l.tail {
		l.tail = n.prev
	}
}

func (l *list) display() {
	p := l.head
	for {

		if p.next == l.head {
			fmt.Printf("-> %dt ", p.value)
			break
		}

		if p == l.head {
			fmt.Printf("-> %dh ", p.value)
		} else {
			fmt.Printf("-> %d ", p.value)

		}

		p = p.next
	}

	fmt.Println()
}

func parseInput(in string, part2 bool) ([]*node, *list) {
	lines := strings.Split(strings.TrimSpace(in), "\n")
	list := &list{}
	var nodes []*node

	for _, l := range lines {
		v, err := strconv.Atoi(l)
		if err != nil {
			panic("invalid input")
		}

		if part2 {
			v *= decryptionKey
		}
		node := list.insert(v)
		nodes = append(nodes, node)
	}

	return nodes, list
}

// mix returns the zero node
func mix(nodes []*node, l *list) *node {
	var zero *node
	for _, n := range nodes {
		p := n
		switch {
		case n.value > 0:
			l.removeNode(n)
			for i := 0; i < n.value; i++ {
				p = p.next
			}
		case n.value < 0:
			l.removeNode(n)
			for i := 0; i >= n.value; i-- {
				p = p.prev
			}
		case n.value == 0:
			zero = n
			continue
		default:
			panic("invalid input")
		}

		l.insertNode(p, n)
	}

	return zero
}

func mix2(nodes []*node, l *list) *node {
	var zero *node
	length := len(nodes)
	for i := 0; i < 10; i++ {
		for _, n := range nodes {
			p := n
			switch {
			case n.value == 0:
				zero = n
				continue
			default:
				l.removeNode(n)
				if n.value > 0 {
					for i := 0; i < n.value%(length-1); i++ {
						p = p.next
					}
				} else if n.value < 0 {
					for i := 0; i >= n.value%(length-1); i-- {
						p = p.prev
					}
				}
			}

			l.insertNode(p, n)
		}
	}

	return zero
}

func part1(in string) int {
	nodes, list := parseInput(in, false)

	zero := mix(nodes, list)

	sum := 0
	p := zero
	for i := 0; i < 3; i++ {
		for j := 0; j < 1000; j++ {
			p = p.next
		}
		sum += p.value
	}

	return sum
}

func part2(in string) int {
	nodes, list := parseInput(in, true)

	zero := mix2(nodes, list)

	sum := 0
	p := zero
	for i := 0; i < 3; i++ {
		for j := 0; j < 1000; j++ {
			p = p.next
		}
		sum += p.value
	}

	return sum
}

func main() {
	fmt.Printf("Part 1 example: %d\n", part1(exampleInput))
	fmt.Printf("Part 1: %d\n", part1(input))
	fmt.Printf("Part 2 example: %d\n", part2(exampleInput))
	fmt.Printf("Part 2: %d\n", part2(input))
}
