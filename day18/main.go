package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"os"
	"strings"
)

type pos struct {
	x, y int
}

func (p pos) add(o pos) pos {
	return pos{p.x + o.x, p.y + o.y}
}

type maze struct {
	grid map[pos]string
	w, h int
}

func (m *maze) print() {
	for y := 0; y < m.h; y++ {
		for x := 0; x <= m.w; x++ {
			fmt.Printf("%s", m.grid[pos{x, y}])
		}
		fmt.Println()
	}
}

func isDoor(s string) bool {
	return s[0] >= 'A' && s[0] <= 'Z'
}

func isKey(s string) bool {
	return s[0] >= 'a' && s[0] <= 'z'
}

func loadLine(maze *maze, line string, o *pos, keys, doors map[pos]string) {
	for i, rune := range line {
		ch := string(rune)
		pos := pos{i, maze.h}
		maze.grid[pos] = ch
		if ch == "@" {
			o.x = i
			o.y = maze.h
		} else if isKey(ch) {
			keys[pos] = ch
		} else if isDoor(ch) {
			doors[pos] = ch
		}

		maze.w = i
	}
	maze.h = maze.h + 1
}

func (p *pos) isValid(m *maze) bool {
	return p.x >= 0 && p.y >= 0 && p.x < m.w && p.y < m.h
}

var osets = []pos{{-1, 0}, {0, -1}, {1, 0}, {0, 1}}

type queue []pos

func enqueueNeighbors(p pos, m *maze, visited map[pos]bool, q []pos) []pos {
	for _, offset := range osets {
		if neighbor := p.add(offset); neighbor.isValid(m) && !visited[neighbor] {
			if m.grid[neighbor] != "#" {
				visited[neighbor] = true
				q = append(q, neighbor)
			}
		}
	}
	return q
}

type dist struct {
	p pos
	d int
}

type priorityQueue []*dist

func (pq priorityQueue) Len() int           { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool { return pq[i].d < pq[j].d }
func (pq priorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }

func (pq *priorityQueue) Push(x interface{}) {
	item := x.(*dist)
	*pq = append(*pq, item)
}

func (pq priorityQueue) index(p pos) int {
	for i, d := range pq {
		if d.p == p {
			return i
		}
	}
	return -1
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*pq = old[0 : n-1]
	return item
}

func (m *maze) countDistanceToOthers(p pos) []dist {
	visited := make(map[pos]bool)
	visited[p] = true
	q := make([]pos, 0, 100)
	q = enqueueNeighbors(p, m, visited, q)

	rounds := 1

	distance := make([]dist, 0, 100)

	for len(q) > 0 {
		//log.Printf("Queue: %v", q)
		nQ := make([]pos, 0, 100)
		for _, cand := range q {
			val := m.grid[cand]
			if isKey(val) || isDoor(val) {
				distance = append(distance, dist{cand, rounds})
			}
			nQ = enqueueNeighbors(cand, m, visited, nQ)
		}
		q = nQ
		rounds++
	}
	return distance
}

func (m *maze) printCosts(cost map[pos]int) {
	for k, v := range cost {
		fmt.Printf("%v @ %v -> %v\n", m.grid[k], k, v)
	}
}

// INFINITY is just a big number
const INFINITY = 32000000

func (m *maze) solve(p pos, doors, keys map[pos]string) int {
	cost := make(map[pos]int)

	// Calculate the reachable distance between each point.
	distTo := make(map[pos][]dist)
	for k, v := range m.grid {
		if v == "@" || isDoor(v) || isKey(v) {
			distTo[k] = m.countDistanceToOthers(k)
			cost[k] = INFINITY
		}
	}
	log.Printf("----DISTANCE----")
	for k, v := range distTo {
		log.Printf("%v -> %v\n", k, v)
	}
	cost[p] = 0

	foundKeys := make(map[pos]map[string]bool)
	for k := range distTo {
		foundKeys[k] = make(map[string]bool)
	}

	q := make(priorityQueue, 0, len(keys))
	heap.Init(&q)
	heap.Push(&q, &dist{p, 0})

	prev := make(map[pos]pos)
	for len(q) > 0 {
		m.printCosts(cost)
		n := heap.Pop(&q).(*dist)

		if cost[n.p] < n.d {
			continue
		}

		if k, ok := keys[n.p]; ok {
			foundKeys[n.p][strings.ToUpper(k)] = true
		}
		if len(foundKeys) == len(keys) {
			return cost[n.p]
		}

		log.Printf("Processing %v -> %v", n, m.grid[n.p])
		for _, d := range distTo[n.p] {
			log.Printf(" \\-> Checking %v -> %v", d.p, m.grid[d.p])
			log.Printf(" \\-> Held keys %v", foundKeys[n.p])
			if isDoor(m.grid[d.p]) && !foundKeys[n.p][m.grid[d.p]] {
				log.Printf(" \\-> missing key.")
				continue
			}
			for k := range foundKeys[n.p] {
				foundKeys[d.p][k] = true
			}
			if isKey(m.grid[d.p]) {
				foundKeys[d.p][strings.ToUpper(m.grid[d.p])] = true
			}

			alt := n.d + d.d
			log.Printf(" n: %v alt: %v", d, alt)
			if alt < cost[d.p] {
				cost[d.p] = alt
				prev[d.p] = d.p
				log.Printf(" \\-> Adding search from %v -> %v", d.p, m.grid[d.p])
				heap.Push(&q, &dist{d.p, alt})
			}
		}
	}

	log.Printf("%v", cost)

	maxCost := 0
	for k, v := range cost {
		if isKey(m.grid[k]) {
			maxCost += v
		}
	}

	return maxCost
}

type record struct {
	key   pos
	doors []pos
}

func (m *maze) solve2(p pos, doors, keys map[pos]string) int {

}

func main() {
	maze := &maze{}
	maze.grid = make(map[pos]string)

	origin := &pos{}
	keys := make(map[pos]string)
	doors := make(map[pos]string)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		loadLine(maze, scanner.Text(), origin, keys, doors)
	}
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
	maze.print()
	fmt.Printf("origin: %v\n", origin)
	fmt.Printf("keys (%v): %v\n", len(keys), keys)
	fmt.Printf("doors (%v): %v\n", len(doors), doors)
	fmt.Println("---solving---")
	steps := maze.solve(*origin, doors, keys)
	fmt.Printf("Steps needed %v\n", steps)
}
