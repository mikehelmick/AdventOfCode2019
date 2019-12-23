package main

import (
	"bufio"
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

func (m *maze) isDoor(p pos) bool {
	return isDoor(m.grid[p])
}

func (m *maze) isKey(p pos) bool {
	return isKey(m.grid[p])
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
	p       pos
	d       int
	reqKeys []string
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

func getNeighbors(p pos, m *maze) []pos {
	neighbors := make([]pos, 0, 4)
	for _, offset := range osets {
		if neighbor := p.add(offset); neighbor.isValid(m) {
			if m.grid[neighbor] != "#" {
				neighbors = append(neighbors, neighbor)
			}
		}
	}
	return neighbors
}

func copyDoors(doors []string) []string {
	cpy := make([]string, len(doors))
	copy(cpy, doors)
	return cpy
}

func (m *maze) distanceDFS(p pos, visited map[pos]bool, doorsInWay []string, curDist int, distance map[pos]dist) {
	visited[p] = true
	curDist++
	for _, neighbor := range getNeighbors(p, m) {
		if visited[neighbor] {
			continue
		}
		if m.isDoor(neighbor) || m.isKey(neighbor) {
			// record this door or key
			if m.isDoor(neighbor) {
				doorsInWay = append(copyDoors(doorsInWay), m.grid[neighbor])
			}
			distance[neighbor] = dist{neighbor, curDist, copyDoors(doorsInWay)}
			m.distanceDFS(neighbor, visited, doorsInWay, curDist, distance)
		} else {
			// traverse this node
			m.distanceDFS(neighbor, visited, doorsInWay, curDist, distance)
		}
	}

	delete(visited, p)
}

func (m *maze) countDistanceToOthers(p pos) map[pos]dist {
	log.Printf("Calculating adjency for %v @ %v", m.grid[p], p)
	visited := make(map[pos]bool)
	//enqueueNeighbors(p, m, visited, q)

	distance := make(map[pos]dist)
	doorsInWay := make([]string, 0, 100)
	m.distanceDFS(p, visited, doorsInWay, 0, distance)
	log.Printf("Adjacency for %v @ %v: %v", p, m.grid[p], distance)

	/*
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
	*/
	return distance
}

func (m *maze) printCosts(cost map[pos]int) {
	for k, v := range cost {
		fmt.Printf("%v @ %v -> %v\n", m.grid[k], k, v)
	}
}

// INFINITY is just a big number
const INFINITY = 32000000

/*
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
*/

func hasKey(m *maze, heldKeys map[string]bool, p pos) bool {
	return heldKeys[m.grid[p]]
}

func pickupKey(m *maze, heldKeys map[string]bool, p pos) {
	heldKeys[strings.ToUpper(m.grid[p])] = true
}

func dropKey(m *maze, heldKeys map[string]bool, p pos) {
	delete(heldKeys, strings.ToUpper(m.grid[p]))
}

func hasNecessaryKeys(heldKeys map[string]bool, doors []string) bool {
	for _, key := range doors {
		if !heldKeys[key] {
			return false
		}
	}
	return true
}

func (m *maze) printStack(stack []pos) {
	s := ""
	for _, p := range stack {
		s = fmt.Sprintf("%s -> %s", s, m.grid[p])
	}
	log.Printf("%s\n", s)
}

func (m *maze) runDfs(stack []pos, visited map[pos]bool, heldKeys map[string]bool, target int, distTo map[pos]map[pos]dist, dist int, toBeat int) int {
	if dist > toBeat {
		return toBeat
	}
	p := stack[len(stack)-1]
	if len(heldKeys) == target {
		//log.Printf("Searching from %v -> %v is last key at dist %v", p, m.grid[p], dist)
		return dist
	}
	visited[p] = true
	log.Printf("%v :: Searching from %v -> %v, dist:%v", stack, p, m.grid[p], dist)
	//log.Printf("Adjacency: %v", distTo[p])
	result := toBeat

	waitFor := 0
	resCh := make(chan int)

	for _, dst := range distTo[p] {
		if visited[dst.p] {
			continue
		}
		log.Printf(" \\-> Considering %v -> %v, haveKeys %v", dst, m.grid[dst.p], heldKeys)
		if m.isDoor(dst.p) || m.isKey(dst.p) {
			if !hasNecessaryKeys(heldKeys, dst.reqKeys) {
				//log.Printf(" \\-> Don't have necessary key for %v -> %v", dst.p, m.grid[dst.p])
			} else {
				newStack := make([]pos, len(stack)+1)
				copy(newStack, stack)
				newStack[len(stack)] = dst.p

				if dist == 0 {
					newVisited := make(map[pos]bool)
					for k, v := range visited {
						newVisited[k] = v
					}

					newKeys := make(map[string]bool)
					for k, v := range heldKeys {
						newKeys[k] = v
					}
					waitFor++
					go func() {
						resCh <- m.runDfs(newStack, newVisited, newKeys, target, distTo, dist+dst.d, toBeat)
					}()
				} else {
					if m.isKey(dst.p) {
						pickupKey(m, heldKeys, dst.p)
					}

					stop := make(chan int)
					go func() {
						stop <- m.runDfs(newStack, visited, heldKeys, target, distTo, dist+dst.d, toBeat)
					}()
					res := <-stop
					//} else {
					//	m.runDfs(newStack, visited, heldKeys, target, distTo, dist+dst.d, toBeat)
					//}

					if res < result {
						if len(heldKeys) == target {
							log.Printf("New shortest path: %v", res)
							//m.printStack(stack)
						}
						result = res
						toBeat = res
					}
					if m.isKey(dst.p) {
						dropKey(m, heldKeys, dst.p)
					}
				}
			}
		}
	}
	for waitFor > 0 {
		res := <-resCh
		log.Printf("New candidate: %v", res)
		if res < result {
			result = res
		}
		waitFor--
	}

	delete(visited, p)
	return result
}

// Solve using DFS
func (m *maze) solve2(p pos, doors, keys map[pos]string) int {
	visited := make(map[pos]bool)
	heldKeys := make(map[string]bool)
	distTo := make(map[pos]map[pos]dist)
	for k, v := range m.grid {
		if v == "@" || isDoor(v) || isKey(v) {
			distTo[k] = m.countDistanceToOthers(k)
		}
	}
	log.Printf("----------------------------------")
	stack := make([]pos, 1)
	stack[0] = p
	return m.runDfs(stack, visited, heldKeys, len(keys), distTo, 0, 7902)
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
	steps := maze.solve2(*origin, doors, keys)
	fmt.Printf("Steps needed %v\n", steps)
}
