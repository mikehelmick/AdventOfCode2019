package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

type point struct {
	x, y int
}

func (p point) add(d point) point {
	return point{p.x + d.x, p.y + d.y}
}

func (p *point) isValid(w, h int) bool {
	return p.x >= 0 && p.y >= 0 && p.x < w && p.y < h
}

var osets = []point{{-1, 0}, {-1, -1}, {0, -1}, {-1, 1}, {1, 0}, {1, 1}, {0, 1}, {-1, 1}}

func enqueueNeighbors(p point, visited map[point]bool, w, h int, q chan point) {
	for _, offset := range osets {
		if neighbor := p.add(offset); neighbor.isValid(w, h) && !visited[neighbor] {
			visited[neighbor] = true
			q <- neighbor
		}
	}
}

// Are a b and c coliniear w/ b and c on the same vector
func colinear(a, b, c point) bool {
	val := a.x*(b.y-c.y) +
		b.x*(c.y-a.y) +
		c.x*(a.y-b.y)
	if val != 0 {
		return false
	}

	return (b.x <= a.x && b.y <= a.y && c.x <= b.x && c.y <= b.y) ||
		(b.x >= a.x && b.y <= a.y && c.x >= b.x && c.y <= b.y) ||
		(b.x >= a.x && b.y >= a.y && c.x >= b.x && c.y >= b.y) ||
		(b.x <= a.x && b.y >= a.y && c.x <= b.x && c.y >= b.y)
}

func search(f [][]int, p point, w, h int, done chan bool) {
	defer func() { done <- true }()

	if f[p.x][p.y] == 0 {
		return
	}
	defer func() { f[p.x][p.y]-- }()

	seen := make(map[point]bool)
	visited := make(map[point]bool)
	visited[p] = true

	log.Printf("Starting search %v", p)

	q := make(chan point, w*h)
	enqueueNeighbors(p, visited, w, h, q)
	for {
		select {
		case candidate := <-q:
			enqueueNeighbors(candidate, visited, w, h, q)
			if f[candidate.x][candidate.y] == 0 {
				// no data here, move on.
				//log.Printf("no data")
				continue
			}
			//log.Printf("candidate: %v", candidate)
			// Is there a point in the seen set that is colinear w/ this point and the origin.
			visible := true
			for s := range seen {
				//log.Printf("  checking seen %v", s)
				if colinear(p, s, candidate) {
					//log.Printf("  %v is blocked by %v", candidate, s)
					visible = false
					break
				}
			}
			seen[candidate] = true
			if visible {
				//log.Printf("  %v is visible by %v", candidate, p)
				f[p.x][p.y]++
			}
		case <-time.After(1 * time.Second):
			return
		}
	}
}

func main() {
	field := make([][]int, 0, 50)

	row := 0
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		field = append(field, make([]int, len(text)))
		for i, ch := range text {
			if ch == '#' {
				field[row][i] = 1
			}
		}
		row++
	}
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}

	width := len(field[0])
	height := len(field)
	log.Printf("w,h = %v,%v", width, height)
	log.Printf("Loaded %v", field)

	done := make(chan bool, 100)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			go search(field, point{x, y}, width, height, done)
		}
	}
	//go search(field, point{2, 4}, height, width, done)
	//<-done
	for i := 0; i < width*height; i++ {
		<-done
	}

	max := 0
	var maxP point
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			fmt.Printf("%4d ", field[x][y])
			if field[x][y] > max {
				max = field[x][y]
				maxP = point{x, y}
			}
		}
		fmt.Printf("\n")
	}
	log.Printf("Max: %v @ %v", max, maxP)
}
