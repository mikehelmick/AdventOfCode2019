package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func loadLine(row int, values [][]int64, line string) {
	for i, rune := range line {
		ch := string(rune)
		if ch == "#" {
			values[row][i] = 1
		}
	}
}

func bioDiversity(values [][]int64) int64 {
	d := int64(0)
	d += values[0][0] * 1
	d += values[0][1] * 2
	d += values[0][2] * 4
	d += values[0][3] * 8
	d += values[0][4] * 16

	d += values[1][0] * 32
	d += values[1][1] * 64
	d += values[1][2] * 128
	d += values[1][3] * 256
	d += values[1][4] * 512

	d += values[2][0] * 1024
	d += values[2][1] * 2048
	d += values[2][2] * 4096
	d += values[2][3] * 8192
	d += values[2][4] * 16384

	d += values[3][0] * 32768
	d += values[3][1] * 65536
	d += values[3][2] * 131072
	d += values[3][3] * 262144
	d += values[3][4] * 524288

	d += values[4][0] * 1048576
	d += values[4][1] * 2097152
	d += values[4][2] * 4194304
	d += values[4][3] * 8388608
	d += values[4][4] * 16777216

	return d
}

func signature(values [][]int64) string {
	s := ""
	for r := 0; r < 5; r++ {
		for c := 0; c < 5; c++ {
			s = fmt.Sprintf("%s%d", s, values[r][c])
		}
	}
	return s
}

type pos struct {
	r, c int
}

var offsets = []pos{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

func (p pos) add(o pos) pos {
	return pos{p.r + o.r, p.c + o.c}
}

func (p pos) isValid() bool {
	return p.r >= 0 && p.r <= 4 && p.c >= 0 && p.c <= 4
}

func adjacentSpaces(p pos) []pos {
	adj := make([]pos, 0, 8)
	for _, o := range offsets {
		if can := p.add(o); can.isValid() {
			adj = append(adj, can)
		}
	}
	return adj
}

func adjacentCount(values [][]int64, p pos) int {
	adj := 0
	for _, p := range adjacentSpaces(p) {
		if values[p.r][p.c] == 1 {
			adj++
		}
	}
	return adj
}

func generation(values [][]int64) [][]int64 {
	next := make([][]int64, 5)
	for i := 0; i < 5; i++ {
		next[i] = make([]int64, 5)
	}

	for r := 0; r < 5; r++ {
		for c := 0; c < 5; c++ {
			adj := adjacentCount(values, pos{r, c})
			next[r][c] = values[r][c]
			if values[r][c] == 1 && adj != 1 {
				next[r][c] = 0
			} else if values[r][c] == 0 && (adj == 2 || adj == 1) {
				next[r][c] = 1
			}
		}
	}

	return next
}

func print(values [][]int64) {
	for r := 0; r < 5; r++ {
		for c := 0; c < 5; c++ {
			if values[r][c] == 1 {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println("")
	}
}

func main() {
	values := make([][]int64, 5)

	scanner := bufio.NewScanner(os.Stdin)
	for row := 0; scanner.Scan() && row < 5; row++ {
		values[row] = make([]int64, 5)
		loadLine(row, values, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}

	fmt.Printf("-- generation 0 --\n")
	print(values)

	cache := make(map[string]bool)
	cache[signature(values)] = true
	gen := 1
	done := false
	for !done {
		values = generation(values)
		sig := signature(values)

		fmt.Printf("-- generation %v --\n", gen)
		print(values)
		gen++

		if v, ok := cache[sig]; ok && v {
			// match
			log.Printf("P1 Answer %v", bioDiversity(values))
			done = true
		} else {
			cache[sig] = true
		}
	}
}
