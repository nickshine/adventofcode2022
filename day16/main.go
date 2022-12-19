package main

import (
	_ "embed"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

//go:embed example.txt
var exampleInput string

//go:embed input.txt
var input string

var inputRE = regexp.MustCompile(`^Valve ([A-Z]{2}) has flow rate=(\d+); tunnels? leads? to valves? ([A-Z, ]+)$`)

type node struct {
	ID string
	v  int
	i  int
}

func (n node) String() string {
	return fmt.Sprintf("%d:%s:%d", n.i, n.ID, n.v)
}

type edge struct {
	from, to string
	weight   int
}

type graph struct {
	nodes map[string]node
	edges map[string]map[edge]struct{}
}

func newGraph() *graph {
	return &graph{
		nodes: make(map[string]node),
		edges: make(map[string]map[edge]struct{}),
	}
}

func (g *graph) addNode(n node) {
	g.nodes[n.ID] = n
}

func (g *graph) addEdge(uID, vID string, weight int) {
	if _, ok := g.edges[uID]; !ok {
		g.edges[uID] = make(map[edge]struct{})
	}

	g.edges[uID][edge{uID, vID, weight}] = struct{}{}

	// undirected
	if _, ok := g.edges[vID]; !ok {
		g.edges[vID] = make(map[edge]struct{})
	}
	g.edges[vID][edge{vID, uID, weight}] = struct{}{}
}

func (g *graph) display() {
	for id, n := range g.nodes {
		fmt.Printf("Node: %s (%d)\n", id, n.v)
		for e := range g.edges[id] {
			fmt.Printf("  %s->%s\n", e.from, e.to)
		}
	}
}
func readGraph(in string) *graph {
	g := newGraph()
	lines := strings.Split(strings.Trim(in, "\n"), "\n")
	for i, l := range lines {
		res := inputRE.FindStringSubmatch(l)
		if len(res) != 4 {
			panic("invalid input")
		}

		valve := res[1]
		rate, err := strconv.Atoi(res[2])
		if err != nil {
			panic("invalid input")
		}
		adj := strings.Split(res[3], ", ")
		g.addNode(node{ID: valve, v: rate, i: i})
		for _, n := range adj {
			g.addEdge(valve, n, 1) // each edge is 1 weight, or 1 minute cost
		}

	}

	return g
}

func deleteNode(nodes map[string]node, k string) map[string]node {

	out := make(map[string]node, len(nodes))

	for k, v := range nodes {
		out[k] = v
	}

	delete(out, k)
	return out
}

// allShortest is an implementation of Floyd-Warshall algorithm to find
// shortest paths between each pair of nodes.
func (g *graph) allShortest() map[string]map[string]int {
	// the distance from every pair of nodes
	dist := make(map[string]map[string]int, len(g.nodes))

	// set all dists to max, and selfs to 0
	for k := range g.nodes {
		dist[k] = make(map[string]int, len(g.nodes))
		for n := range g.nodes {
			dist[k][n] = math.MaxInt32
		}
		dist[k][k] = 0 // node to itself is 0
	}

	// set distances for all adjacent nodes
	for id, adj := range g.edges {
		for edge := range adj {
			dist[id][edge.to] = edge.weight
		}
	}

	for k, dmk := range dist { // for each node k, distance map dmk
		for _, dmi := range dist { // for each distance map dmi
			for n, d := range dmi { // for each node n, distance d
				if distance := dmi[k] + dmk[n]; d > distance {
					dmi[n] = distance
				}
			}
		}
	}

	return dist

}

func release(nodes map[string]node, distances map[string]map[string]int, node string, time, pressure, flow, limit int) int {
	// assume no more valves will open at first
	max := pressure + (limit-time)*flow

	for id, n := range nodes {
		if n.v == 0 { // don't bother calculating if node has zero flow
			continue
		}
		cost := int(distances[node][id]) + 1 //shortest path from node to id node, + 1 to open
		if time+cost >= limit {
			continue
		}
		t := time + cost
		p := pressure + cost*flow
		f := flow + n.v
		result := release(deleteNode(nodes, id), distances, id, t, p, f, limit)
		if result > max {
			max = result
		}
	}

	return max
}

func visit(g *graph, src string, opened, time, released int, distances map[string]map[string]int, state map[int]int) {
	if time <= 0 {
		return
	}

	if r, ok := state[opened]; !ok || released > r {
		state[opened] = released
	}

	for id, n := range g.nodes {
		if n.v == 0 { // don't bother calculating if node has zero flow
			continue
		}
		if opened&(1<<n.i) != 0 {
			continue
		}

		cost := (time - distances[src][id] - 1) // minus one to open
		score := n.v * cost
		visit(g, id, opened|(1<<n.i), cost, released+score, distances, state)
	}
}

func part1(in string) int {
	g := readGraph(in)
	distances := g.allShortest()
	// max := release(g.nodes, distances, "AA", 0, 0, 0, 30)

	state := make(map[int]int)
	visit(g, "AA", 0, 30, 0, distances, state)
	max := 0
	for _, v := range state {
		if v > max {
			max = v
		}
	}
	return max
}

func part2(in string) int {
	g := readGraph(in)
	distances := g.allShortest()
	state := make(map[int]int)

	max := 0
	visit(g, "AA", 0, 26, 0, distances, state)
	for p1, v1 := range state {
		for p2, v2 := range state {
			if (p1 & p2) == 0 {
				if v1+v2 > max {
					max = v1 + v2
				}
			}
		}

	}

	return max
}

func main() {
	fmt.Printf("Part 1 example: %d\n", part1(exampleInput))
	fmt.Printf("Part 1: %d\n", part1(input))
	fmt.Printf("Part 2 example: %d\n", part2(exampleInput))
	fmt.Printf("Part 2: %d\n", part2(input))
}
