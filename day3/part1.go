package main

import (
	"fmt"
	"os"
	"strings"
)

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func sharedItem(rucksack string) string {

	seen := make(map[rune]struct{})
	for _, r := range rucksack[0 : len(rucksack)/2] {
		seen[r] = struct{}{}
	}

	for _, r := range rucksack[len(rucksack)/2:] {
		if _, ok := seen[r]; ok {
			return string(r)
		}
	}

	return ""
}

func main() {
	data, _ := os.ReadFile(os.Args[1])
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")

	total := 0
	for _, rucksack := range lines {

		shared := sharedItem(rucksack)

		priority := strings.Index(chars, shared) + 1
		total += priority
	}

	fmt.Println(total)
}
