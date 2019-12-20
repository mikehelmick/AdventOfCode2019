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
	name            string
	inside, outside pos
}

type maze struct {
	grid    map[pos]string
	warp    map[string]portal
	warpPos map[pos]string
	w, h    int
}

func (m *maze) isInsideWarp(p pos) bool {
	if pName, ok := m.warpPos[p]; ok {
		return m.warp[pName].inside == p
	}
	return false
}

func (m *maze) isOutsideWarp(p pos) bool {
	if pName, ok := m.warpPos[p]; ok {
		return m.warp[pName].outside == p
	}
	return false
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

func (m *maze) updatePortal(name string, p pos, inside bool) {
	if _, ok := m.warp[name]; !ok {
		m.warp[name] = portal{name, pos{}, pos{}}
	}
	if pr := m.warp[name]; name == "AA" || name == "ZZ" || !inside {
		pr.outside = p
		m.warp[name] = pr
	} else {
		pr.inside = p
		m.warp[name] = pr
	}
	if name != "AA" && name != "ZZ" {
		m.warpPos[p] = name
	}
	m.grid[p] = "8"
}

// Get warp position - either down one level (-1) or up one level (1)
func (m *maze) getWarpPos(p pos) (pos, int) {
	for _, v := range m.warp {
		if v.inside == p {
			// descend 1 level
			return v.inside, -1
		} else if v.outside == p {
			// ascend 1 level
			return v.outside, 1
		}
	}
	log.Fatalf("Unable to find pair for warp point %v", p)
	panic("I just can't.")
}

func (m *maze) isPosOutside(p pos) bool {
	return (p.x == 2 || p.x == m.w-2) ||
		(p.y == 2 || p.y == m.h-3)
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
						portalPos := pos{x - 1, y}
						m.updatePortal(pName, portalPos, !m.isPosOutside(portalPos))
					} else {
						portalPos := pos{x + 2, y}
						m.updatePortal(pName, portalPos, !m.isPosOutside(portalPos))
					}
				} else if m.isChar(pos{x, y + 1}) {
					pName := fmt.Sprintf("%s%s", m.grid[p], m.grid[pos{x, y + 1}])
					// Update portal. Position of the portal is either above the first or
					// below the second.
					if m.isEmpty(pos{x, y - 1}) {
						portalPos := pos{x, y - 1}
						m.updatePortal(pName, portalPos, !m.isPosOutside(portalPos))
					} else {
						if pName == "OA" {
							log.Printf(" \\-> Above the maze")
						}
						portalPos := pos{x, y + 2}
						m.updatePortal(pName, portalPos, !m.isPosOutside(portalPos))
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

type posAndLevel struct {
	p     pos
	level int
}

type posAndLevelAndWarp struct {
	pLevel posAndLevel
	warpOk bool
}

func (p posAndLevelAndWarp) add(o pos) posAndLevelAndWarp {
	return posAndLevelAndWarp{posAndLevel{pos{p.pLevel.p.x + o.x, p.pLevel.p.y + o.y}, p.pLevel.level}, true}
}

func enqueueNeighbors(p posAndLevelAndWarp, m *maze, visited map[posAndLevel]bool, q []posAndLevelAndWarp) []posAndLevelAndWarp {
	if p.warpOk {
		// if you're on a portal, go to the other end.
		if pName, ok := m.warpPos[p.pLevel.p]; ok {
			if m.isInsideWarp(p.pLevel.p) {
				newPos := posAndLevelAndWarp{posAndLevel{m.warp[pName].outside, p.pLevel.level + 1}, false}
				visited[newPos.pLevel] = true
				return append(q, newPos)
			}
			if p.pLevel.level > 0 {
				newPos := posAndLevelAndWarp{posAndLevel{m.warp[pName].inside, p.pLevel.level - 1}, false}
				visited[newPos.pLevel] = true
				return append(q, newPos)
			}
		}
	}

	for _, offset := range osets {
		if neighbor := p.add(offset); neighbor.pLevel.p.isValid(m) && !visited[neighbor.pLevel] {
			if m.grid[neighbor.pLevel.p] == "." || m.grid[neighbor.pLevel.p] == "8" {
				if m.grid[neighbor.pLevel.p] == "8" && m.warpPos[neighbor.pLevel.p] == "ZZ" && neighbor.pLevel.level != 0 {
					continue
				} else if m.grid[neighbor.pLevel.p] == "8" && neighbor.pLevel.level == 0 && m.warpPos[neighbor.pLevel.p] != "ZZ" && m.isOutsideWarp(neighbor.pLevel.p) {
					continue
				}
				visited[neighbor.pLevel] = true
				q = append(q, neighbor)
			}
		}
	}
	return q
}

type visitedMap map[pos]bool

func (m *maze) search(origin, goal pos) int {
	q := make([]posAndLevelAndWarp, 0, 100)
	q = append(q, posAndLevelAndWarp{posAndLevel{origin, 0}, false})
	dist := 0

	actualGoal := posAndLevel{goal, 0}

	visited := make(map[posAndLevel]bool)

	for len(q) > 0 {
		nextQ := make([]posAndLevelAndWarp, 0, 100)

		//log.Printf("Wave %v front: %v", dist, q)

		for _, can := range q {
			if can.pLevel == actualGoal {
				return dist
			}
			nextQ = enqueueNeighbors(can, m, visited, nextQ)
		}
		dist++
		q = nextQ
	}

	return dist
}

func (m *maze) printPortals() {
	fmt.Printf("Dim: %v, %v\n", m.w, m.h)
	for k, v := range m.warp {
		fmt.Printf("portal %v\n", k)
		fmt.Printf(" \\-> inside %v\n", v.inside)
		fmt.Printf(" \\-> outside %v\n", v.outside)
	}
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
	maze.printPortals()
	origin := maze.warp["AA"].outside
	log.Printf("Origin: %v", origin)
	goal := maze.warp["ZZ"].outside
	log.Printf("Goal: %v", goal)
	maze.print()

	steps := maze.search(origin, goal)
	fmt.Printf("Answer: %v steps\n", steps)
}
