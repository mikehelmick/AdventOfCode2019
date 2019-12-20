package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"unicode/utf8"
)

type pos struct {
	x, y int
}

func (p pos) add(o pos) pos {
	return pos{p.x + o.x, p.y + o.y}
}

type portal struct {
	name string
	a, b pos
}

type maze struct {
	grid    map[pos]string
	warp    map[string]portal
	warpPos map[pos]string
	w, h    int
}

func (m *maze) print() {
	for y := 0; y < m.h; y++ {
		for x := 0; x <= m.w; x++ {
			if s, ok := m.grid[pos{x, y}]; ok {
				fmt.Printf("%s", s)
			} else {
				fmt.Printf(" ")
			}
		}
		fmt.Println()
	}
}

func (m *maze) isChar(p pos) bool {
	if s, ok := m.grid[p]; ok {
		r, _ := utf8.DecodeRuneInString(s[0:])
		return r >= 'A' && r <= 'Z'
	}
	return false
}

func (m *maze) isEmpty(p pos) bool {
	if s, ok := m.grid[p]; ok {
		return s == "."
	}
	return false
}

func (m *maze) updatePortal(name string, p pos) {
	if pr, ok := m.warp[name]; ok {
		pr.b = p
		m.warp[name] = pr
	} else {
		m.warp[name] = portal{name, p, pos{}}
	}
	if name != "AA" && name != "ZZ" {
		m.warpPos[p] = name
	}
	m.grid[p] = "8"
}

func (m *maze) getWarpPos(p pos) pos {
	for _, v := range m.warp {
		if v.a == p {
			return v.b
		} else if v.b == p {
			return v.a
		}
	}
	log.Fatalf("Unable to find pair for warp point %v", p)
	return pos{0, 0}
}

func (m *maze) loadPortals() {

	for y := 0; y < m.h; y++ {
		for x := 0; x <= m.w; x++ {
			p := pos{x, y}
			if m.isChar(p) {
				// Either the char to the right or the char below must be a rune
				if m.isChar(pos{x + 1, y}) {
					pName := fmt.Sprintf("%s%s", m.grid[p], m.grid[pos{x + 1, y}])
					// Updat portal. Position of the portal could be to the left or the right
					// of this.
					if m.isEmpty(pos{x - 1, y}) {
						m.updatePortal(pName, pos{x - 1, y})
					} else {
						m.updatePortal(pName, pos{x + 2, y})
					}
				} else if m.isChar(pos{x, y + 1}) {
					pName := fmt.Sprintf("%s%s", m.grid[p], m.grid[pos{x, y + 1}])
					// Update portal. Position of the portal is either above the first or
					// below the second.
					if m.isEmpty(pos{x, y - 1}) {
						m.updatePortal(pName, pos{x, y - 1})
					} else {
						m.updatePortal(pName, pos{x, y + 2})
					}
				}
				// Else - we've already processed this portal
			}
		}
	}
}

func loadLine(maze *maze, line string) {
	for i, rune := range line {
		ch := string(rune)
		pos := pos{i, maze.h}
		maze.grid[pos] = ch
		if i > maze.w {
			maze.w = i

		}
	}
	maze.h = maze.h + 1
}

func (p *pos) isValid(m *maze) bool {
	return p.x >= 0 && p.y >= 0 && p.x < m.w && p.y < m.h
}

var osets = []pos{{-1, 0}, {0, -1}, {1, 0}, {0, 1}}

type queue []pos

func enqueueNeighbors(p pos, m *maze, visited map[pos]bool, q []pos) []pos {
	// if you're on a portal, go to the other end.
	if pName, ok := m.warpPos[p]; ok {
		next := m.getWarpPos(p)
		q = append(q, next)
		// destroy the portal so that it doesnt get used again.
		delete(m.warpPos, next)
		delete(m.warpPos, p)
		delete(m.warp, pName)
		return q
	}

	for _, offset := range osets {
		if neighbor := p.add(offset); neighbor.isValid(m) && !visited[neighbor] {
			if m.grid[neighbor] == "." || m.grid[neighbor] == "8" {
				visited[neighbor] = true
				q = append(q, neighbor)
			}
		}
	}
	return q
}

func (m *maze) search(origin, goal pos) int {
	q := make([]pos, 0, 100)
	q = append(q, origin)
	dist := 0

	visited := make(map[pos]bool)
	visited[origin] = true

	for len(q) > 0 {
		nextQ := make([]pos, 0, 100)

		log.Printf("Wave %v front: %v", dist, q)

		for _, can := range q {
			if can == goal {
				return dist
			}
			nextQ = enqueueNeighbors(can, m, visited, nextQ)
		}
		dist++
		q = nextQ
	}

	return dist
}

func main() {
	maze := &maze{}
	maze.grid = make(map[pos]string)
	maze.warp = make(map[string]portal)
	maze.warpPos = make(map[pos]string)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		loadLine(maze, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
	maze.loadPortals()
	log.Printf("Portals: %v", maze.warp)
	origin := maze.warp["AA"].a
	log.Printf("Origin: %v", origin)
	goal := maze.warp["ZZ"].a
	log.Printf("Goal: %v", goal)
	maze.print()

	steps := maze.search(origin, goal)
	fmt.Printf("Answer: %v steps\n", steps)
}
