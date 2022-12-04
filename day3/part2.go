package main

import (
	"fmt"
	"os"
	"strings"
)

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// groups chucks lines in to groups of 3
func groups(lines []string) [][]string {

	var groups [][]string
	for {

		if len(lines) < 3 {
			break
		}

		groups = append(groups, lines[0:3])
		lines = lines[3:]
	}

	return groups
}

// set dedups a strings chars
func set(s string) string {
	seen := make(map[rune]struct{})
	for _, r := range s {
		seen[r] = struct{}{}
	}

	var result []rune
	for r := range seen {
		result = append(result, r)
	}

	return string(result)
}

func sharedItem(group []string) string {
	seen := make(map[rune]int)
	for _, rucksack := range group {
		for _, r := range set(rucksack) {
			seen[r] += 1
		}
	}

	for r, v := range seen {
		if v == 3 {
			return string(r)
		}
	}

	return ""
}

func main() {
	data, _ := os.ReadFile(os.Args[1])
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")

	total := 0
	for _, group := range groups(lines) {

		shared := sharedItem(group)
		priority := strings.Index(chars, shared) + 1
		total += priority

	}

	fmt.Println(total)
}
