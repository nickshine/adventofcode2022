package main

import (
	_ "embed"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

//go:embed example.txt
var exampleInput string

//go:embed input.txt
var input string

var inputRE = regexp.MustCompile(`^Blueprint (\d+): Each ore robot costs (\d+) ore. Each clay robot costs (\d+) ore. Each obsidian robot costs (\d+) ore and (\d+) clay. Each geode robot costs (\d+) ore and (\d+) obsidian.$`)

type blueprint struct {
	id                int
	oreOreCost        int
	clayOreCost       int
	obsidianOreCost   int
	obsidianClayCost  int
	geodeOreCost      int
	geodeObsidianCost int
	maxOre            int
	maxClay           int
	maxObsidian       int
}

func parseBlueprints(in string) []blueprint {
	var blueprints []blueprint
	lines := strings.Split(strings.TrimSpace(in), "\n")
	for _, l := range lines {
		res := inputRE.FindStringSubmatch(l)
		if len(res) != 8 {
			panic("invalid input")
		}

		id, _ := strconv.Atoi(res[1])
		oreOreCost, _ := strconv.Atoi(res[2])
		clayOreCost, _ := strconv.Atoi(res[3])
		obsidianOreCost, _ := strconv.Atoi(res[4])
		obsidianClayCost, _ := strconv.Atoi(res[5])
		geodeOreCost, _ := strconv.Atoi(res[6])
		geodeObsidianCost, _ := strconv.Atoi(res[7])

		maxOre := oreOreCost
		maxClay := obsidianClayCost
		maxObsidian := geodeObsidianCost

		if clayOreCost > maxOre {
			maxOre = clayOreCost
		}
		if obsidianOreCost > maxOre {
			maxOre = obsidianOreCost
		}
		if geodeOreCost > maxOre {
			maxOre = geodeOreCost
		}

		blueprints = append(blueprints, blueprint{
			id,
			oreOreCost,
			clayOreCost,
			obsidianOreCost,
			obsidianClayCost,
			geodeOreCost,
			geodeObsidianCost,
			maxOre,
			maxClay,
			maxObsidian,
		})
	}

	return blueprints
}

const (
	maxTime = 24
)

type state struct {
	oreRobots      int
	clayRobots     int
	obsidianRobots int
	geodeRobots    int
	ore            int
	clay           int
	obsidian       int
	geode          int
	best           map[int]int // best geode for a time
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func run(b blueprint, time int, s state) int {
	if time > maxTime {
		return s.geode
	}

	ore, clay, obsidian := s.ore, s.clay, s.obsidian
	// log.Printf("ore, clay obsidian: %d,%d,%d, time: %d", ore, clay, obsidian, time)

	// each robot collects
	s.ore += s.oreRobots
	s.clay += s.clayRobots
	s.obsidian += s.obsidianRobots
	s.geode += s.geodeRobots

	// don't continue searching if there is already a time with more geodes
	if s.geode < s.best[time] {
		return s.geode
	}

	s.best[time] = max(s.geode, s.best[time])

	var options []state
	// geode - if geode robot can be built, force build it over other options
	if ore >= b.geodeOreCost && obsidian >= b.geodeObsidianCost {
		opt := s
		opt.ore -= b.geodeOreCost
		opt.obsidian -= b.geodeObsidianCost
		opt.geodeRobots++
		// log.Printf("geode spend, time %d, option: %#v", time, opt)
		options = append(options, opt)
	} else {
		// ore
		if ore >= b.oreOreCost && s.oreRobots < b.maxOre {
			opt := s
			opt.ore -= b.oreOreCost
			opt.oreRobots++
			// log.Printf("ore spend, time %d, option: %#v", time, opt)
			options = append(options, opt)
		}
		// clay
		if ore >= b.clayOreCost && s.clayRobots < b.maxClay {
			opt := s
			opt.ore -= b.clayOreCost
			opt.clayRobots++
			// log.Printf("clay spend, time %d, option: %#v", time, opt)
			options = append(options, opt)
		}
		// obsidian
		if ore >= b.obsidianOreCost && clay >= b.obsidianClayCost && s.obsidianRobots < b.maxObsidian {
			opt := s
			opt.ore -= b.obsidianOreCost
			opt.clay -= b.obsidianClayCost
			opt.obsidianRobots++
			// log.Printf("obsidian spend, time %d, option: %#v", time, opt)
			options = append(options, opt)
		}

	}
	maxGeode := 0
	for _, opt := range options {
		numGeode := run(b, time+1, opt)
		if numGeode > maxGeode {
			maxGeode = numGeode
		}
	}

	// log.Printf("time: %d, state: %#v", time, s)
	return max(maxGeode, run(b, time+1, s))

}

func part1(in string) int {
	blueprints := parseBlueprints(in)

	total := 0
	for _, b := range blueprints {
		// log.Printf("%+v", b)
		max := run(b, 1, state{oreRobots: 1, best: map[int]int{}})
		total += b.id * max
	}

	return total
}

func main() {
	fmt.Printf("Part 1 example: %d\n", part1(exampleInput))
	fmt.Printf("Part 1: %d\n", part1(input))
	// fmt.Printf("Part 2 example: %d\n", part2(exampleInput))
	// fmt.Printf("Part 2: %d\n", part2(input))
}
