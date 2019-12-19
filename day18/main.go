package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"os"
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

func loadLine(maze *maze, line string, o *pos, keys, doors map[string]pos) {
	for i, rune := range line {
		ch := string(rune)
		pos := pos{i, maze.h}
		maze.grid[pos] = ch
		if ch == "@" {
			o.x = i
			o.y = maze.h
		} else if isKey(ch) {
			keys[ch] = pos
		} else if isDoor(ch) {
			doors[ch] = pos
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
				continue
			}
			nQ = enqueueNeighbors(cand, m, visited, nQ)
		}
		q = nQ
		rounds++
	}
	return distance
}

func (m *maze) solve(p pos, keys map[string]pos) int {
	cost := make(map[pos]int)

	// Calculate the reachable distance between each point.
	distTo := make(map[pos][]dist)
	for k, v := range m.grid {
		if v == "@" || isDoor(v) || isKey(v) {
			distTo[k] = m.countDistanceToOthers(k)
			cost[k] = 32000000
		}
	}
	cost[p] = 0

	q := make(priorityQueue, 0, len(keys))
	heap.Init(&q)
	heap.Push(&q, &dist{p, 0})

	prev := make(map[pos]pos)

	for len(q) > 0 {
		n := heap.Pop(&q).(*dist)
		log.Printf("Processing %v", n)
		for _, d := range distTo[n.p] {
			alt := n.d + d.d
			log.Printf(" n: %v alt: %v", d, alt)
			if alt < cost[d.p] {
				cost[d.p] = alt
				prev[d.p] = d.p

				if idx := q.index(d.p); idx >= 0 {
					heap.Remove(&q, idx)
				}
				heap.Push(&q, &dist{d.p, alt})
			}
		}
	}

	log.Printf("%v", cost)

	totCost := 0
	for _, v := range cost {
		totCost += v
	}

	return totCost
}

func main() {
	maze := &maze{}
	maze.grid = make(map[pos]string)

	origin := &pos{}
	keys := make(map[string]pos)
	doors := make(map[string]pos)

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
	fmt.Printf("doors: %v\n", doors)
	fmt.Println("---solving---")
	steps := maze.solve(*origin, keys)
	fmt.Printf("Steps needed %v\n", steps)
}
