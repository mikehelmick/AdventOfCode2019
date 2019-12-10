package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
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

func (p *point) distance(o point) float64 {
	return math.Abs(math.Sqrt(math.Pow(float64(o.x-p.x), 2) + math.Pow(float64(o.y-p.y), 2)))
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
// and b sits between a and c
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

	if f[p.y][p.x] == 0 {
		return
	}
	defer func() { f[p.y][p.x]-- }()

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
			if f[candidate.y][candidate.x] == 0 {
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
				f[p.y][p.x]++
			}
		case <-time.After(1 * time.Second):
			return
		}
	}
}

func findAngles(f [][]int, p point, w, h int, angles map[point]int32) {
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if f[y][x] != 0 && !(p.x == x && p.y == y) {
				angle := math.Atan2(float64(y-p.y), float64(x-p.x))*(180.0/math.Pi) + 90.0
				if angle < 0 {
					angle += 360
				}
				var err error
				angle, err = strconv.ParseFloat(fmt.Sprintf("%.2f", angle), 64)
				if err != nil {
					panic("failed float trimming")
				}
				angles[point{x, y}] = int32(angle * 100)
			}
		}
	}
}

type spoint struct {
	p     point
	angle int32
	dist  float64
}

func buildSpoints(p point, angles map[point]int32, res []spoint) {
	idx := 0
	for k, v := range angles {
		res[idx] = spoint{k, v, p.distance(k)}
		idx++
	}
}

type byAngleDist []spoint

func (a byAngleDist) Len() int {
	return len(a)
}

func (a byAngleDist) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a byAngleDist) Less(i, j int) bool {
	if a[i].angle < a[j].angle {
		return true
	} else if a[i].angle == a[j].angle {
		if a[i].dist < a[j].dist {
			return true
		}
	}
	return false
}

func destroy(pts []spoint, destC chan spoint) {
	lastAngle := int32(-1)
	sendAny := false
	for i, pt := range pts {
		if pt.p.x >= 0 && pt.angle != lastAngle {
			lastAngle = pt.angle
			destC <- pt
			pts[i] = spoint{point{-1, -1}, -1, -1}
			sendAny = true
		}
	}
	if sendAny {
		destroy(pts, destC)
	} else {
		close(destC)
	}
}

func main() {
	// Build the field where 1 is astroid, 0 is nothing.
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

	// In parallel, launch a BFS search for visible astroids from each astroid.
	done := make(chan bool, 100)
	defer close(done)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			go search(field, point{x, y}, width, height, done)
		}
	}
	// wait for the serches to be done.
	for i := 0; i < width*height; i++ {
		<-done
	}

	// Calculate the answer to part 1, most visible astroids.
	max := 0
	var maxp point
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			fmt.Printf("%4d ", field[y][x])
			if field[y][x] > max {
				max = field[y][x]
				maxp = point{x, y}
			}
		}
		fmt.Printf("\n")
	}
	log.Printf("Max: %v @ %v", max, maxp)

	// Sweep the laser until the laser doesn't sweep anymore.
	angles := make(map[point]int32)
	findAngles(field, maxp, width, height, angles)

	/* // Debugging - print the angles
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if ang, ok := angles[point{x, y}]; !ok {
				if maxp.x == x && maxp.y == y {
					fmt.Printf("    XX  ")
				} else {
					fmt.Printf("        ")
				}
			} else {
				fmt.Printf("%6d ", ang)
			}
		}
		fmt.Printf("\n")
	}
	*/

	//angOrd = make(chan point, 500)
	spoints := make([]spoint, len(angles))
	buildSpoints(maxp, angles, spoints)
	sort.Sort(byAngleDist(spoints))
	//log.Printf("%v", spoints)

	destC := make(chan spoint)
	go destroy(spoints, destC)
	order := 1
	for pt := range destC {
		log.Printf("O: %4d  -- {%3d,%3d} @ %6d  -- Ans: %5d", order, pt.p.x, pt.p.y, pt.angle, (pt.p.x*100 + pt.p.y))
		order++
	}
}
