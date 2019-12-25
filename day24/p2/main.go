package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
)

// Pos is a position in a single level.
type Pos struct {
	r, c int
}

var offsets = []Pos{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

func (p Pos) add(o Pos) Pos {
	return Pos{p.r + o.r, p.c + o.c}
}

// Within the bounds of this current level, 2,2 isn't valid at any level.
func (p Pos) isValid() bool {
	return p.r >= 0 && p.r <= 4 && p.c >= 0 && p.c <= 4 &&
		!(p.r == 2 && p.c == 2)
}

// Cord represents an adjacent position, including level
type Cord struct {
	level int
	pos   Pos
}

// Add a pos offset, keep current level.
func (c Cord) add(o Pos) Cord {
	return Cord{c.level, Pos{c.pos.r + o.r, c.pos.c + o.c}}
}

func (l Levels) print() {
	var keys []int
	for k := range l {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, depth := range keys {
		fmt.Printf("Depth #%d\n", depth)
		for r := 0; r < 5; r++ {
			for c := 0; c < 5; c++ {
				if r == 2 && c == 2 {
					fmt.Print("?")
				} else {
					if l[depth][r][c] == 1 {
						fmt.Print("#")
					} else {
						fmt.Print(".")
					}
				}
			}
			fmt.Println("")
		}
		fmt.Println("")
	}
}

// Levels represents the map of levels to grids.
type Levels map[int][][]int

func (l Levels) aliveCount(level int) int {
	count := 0
	for r := 0; r < 5; r++ {
		for c := 0; c < 5; c++ {
			if l[level][r][c] == 1 {
				count++
			}
		}
	}
	return count
}

func (l Levels) initLevel(level int) {
	if _, ok := l[level]; ok {
		return
	}
	l[level] = make([][]int, 5)
	for r := 0; r < 5; r++ {
		l[level][r] = make([]int, 5)
	}
}

func (l Levels) minMaxDepth() (min int, max int) {
	min, max = 0, 0
	for k := range l {
		if k < min {
			min = k
		}
		if k > max {
			max = k
		}
	}
	return
}

func loadLine(row int, levels Levels, l int, line string) {
	for i, rune := range line {
		ch := string(rune)
		if ch == "#" {
			levels[l][row][i] = 1
		}
	}
}

func adjacentSpaces(c Cord) []Cord {
	adj := make([]Cord, 0, 8)
	for _, o := range offsets {
		if can := c.add(o); can.pos.isValid() {
			adj = append(adj, can)
		}
	}
	// Special conditions of adjacent spaces. For spaces along outer edge,
	// add spaces from 1 level higher.
	if c.pos.r == 0 {
		adj = append(adj, Cord{c.level + 1, Pos{1, 2}})
	}
	if c.pos.r == 4 {
		adj = append(adj, Cord{c.level + 1, Pos{3, 2}})
	}
	if c.pos.c == 0 {
		adj = append(adj, Cord{c.level + 1, Pos{2, 1}})
	}
	if c.pos.c == 4 {
		adj = append(adj, Cord{c.level + 1, Pos{2, 3}})
	}
	// For spaces along the inner edge, add spaces 1 level lower.
	if c.pos.r == 1 && c.pos.c == 2 {
		// Whole top row of 1 depth lower.
		for col := 0; col < 5; col++ {
			adj = append(adj, Cord{c.level - 1, Pos{0, col}})
		}
	}
	if c.pos.r == 2 && c.pos.c == 3 {
		// Whole right col of 1 depth lower
		for row := 0; row < 5; row++ {
			adj = append(adj, Cord{c.level - 1, Pos{row, 4}})
		}
	}
	if c.pos.r == 3 && c.pos.c == 2 {
		// whole bottom row of 1 depth lower
		for col := 0; col < 5; col++ {
			adj = append(adj, Cord{c.level - 1, Pos{4, col}})
		}
	}
	if c.pos.r == 2 && c.pos.c == 1 {
		// Whole left col of 1 depth lower
		for row := 0; row < 5; row++ {
			adj = append(adj, Cord{c.level - 1, Pos{row, 0}})
		}
	}

	//if c.level == 0 {
	//	log.Printf("DEBUG: Adjacent %v --> %v", c, adj)
	//}

	return adj
}

func (l Levels) adjacentCount(c Cord) int {
	adj := 0
	for _, p := range adjacentSpaces(c) {
		if l.getValue(p) == 1 {
			adj++
		}
	}
	return adj
}

// Does this particular level have anything alive along the inner square.
// Indicates the need to recurse lower.
func (l Levels) aliveInner(lvl int) bool {
	return l[lvl][2][1] == 1 || l[lvl][1][2] == 1 ||
		l[lvl][2][3] == 1 || l[lvl][3][2] == 1
}

func (l Levels) aliveOuter(lvl int) bool {
	for c := 0; c < 5; c++ {
		if l[lvl][0][c] == 1 || l[lvl][4][c] == 1 || l[lvl][c][0] == 1 || l[lvl][c][4] == 1 {
			return true
		}
	}
	return false
}

func (l Levels) getValue(c Cord) int {
	if grid, ok := l[c.level]; ok {
		return grid[c.pos.r][c.pos.c]
	}
	return 0
}

func (l Levels) setValue(c Cord, val int) {
	l[c.level][c.pos.r][c.pos.c] = val
}

func generation(l Levels) Levels {
	// Setup next generation.
	min, max := l.minMaxDepth()
	//log.Printf("Old Gen min,max = (%v, %v)", min, max)
	if l.aliveInner(min) {
		// will have a new lower level in next generation.
		min--
	}
	if l.aliveOuter(max) {
		// new higher level in next generation (likely)
		max++
	}
	var nextLevels Levels = make(map[int][][]int)
	for i := min; i <= max; i++ {
		nextLevels.initLevel(i)
	}
	// Next levels now contains everything needed.
	//log.Printf("New Gen min,max = (%v, %v)", min, max)

	for lvl := min; lvl <= max; lvl++ {
		for r := 0; r < 5; r++ {
			for c := 0; c < 5; c++ {
				if r == 2 && c == 2 {
					continue
				}
				thisCord := Cord{lvl, Pos{r, c}}
				adj := l.adjacentCount(thisCord)
				curValue := l.getValue(thisCord)
				nextLevels[lvl][r][c] = curValue

				if lvl == -1 && r == 4 && c == 4 {
					log.Printf("DEBUG %v %v %v", thisCord, adj, curValue)
				}

				if curValue == 1 && adj != 1 {
					nextLevels.setValue(thisCord, 0)
				} else if curValue == 0 && (adj == 2 || adj == 1) {
					nextLevels.setValue(thisCord, 1)
				}
			}
		}
	}

	return nextLevels
}

func (l Levels) countAlive() int {
	alive := 0
	for lvl := range l {
		for r := 0; r < 5; r++ {
			for c := 0; c < 5; c++ {
				if l[lvl][r][c] == 1 {
					alive++
				}
			}
		}
	}
	return alive
}

func main() {
	var levels Levels = make(map[int][][]int)
	levels.initLevel(0)

	scanner := bufio.NewScanner(os.Stdin)
	for row := 0; scanner.Scan() && row < 5; row++ {
		loadLine(row, levels, 0, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
	fmt.Printf("-- generation 0 --\n")
	levels.print()

	log.Printf("%v", levels)
	for gen := 1; gen <= 200; gen++ {
		levels = generation(levels)
		fmt.Printf("-- generation %v --\n", gen)
		levels.print()
	}

	log.Printf("Alive Count: %v", levels.countAlive())
}
