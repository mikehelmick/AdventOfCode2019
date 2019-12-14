package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type quant struct {
	amt  int64
	elem string
}

func (q quant) print() {
	fmt.Printf("%d %v", q.amt, q.elem)
}

type reaction struct {
	in  []quant
	out quant
}

func newReaction() reaction {
	var r reaction
	r.in = make([]quant, 0, 20)
	return r
}

func (r *reaction) addInput(i quant) {
	r.in = append(r.in, i)
}

func (r reaction) print() {
	for i, q := range r.in {
		q.print()
		if i < len(r.in)-1 {
			fmt.Print(", ")
		}
	}
	fmt.Print(" => ")
	r.out.print()
	fmt.Printf("\n")
}

func parseQuant(s string) quant {
	log.Printf("parsing quantity: '%s'", s)
	parts := strings.Split(strings.TrimSpace(s), " ")
	amt, err := strconv.ParseInt(strings.TrimSpace(parts[0]), 10, 64)
	if err != nil {
		log.Fatalf("Can't parse integer: '%v'", strings.TrimSpace(parts[0]))
	}
	elem := strings.TrimSpace(parts[1])
	return quant{amt, elem}
}

func parseReaction(s string) reaction {
	log.Printf("parsing line '%v'", s)
	r := newReaction()
	halves := strings.Split(s, "=>")
	if len(halves) != 2 {
		log.Fatalf("Invalid reaction format, %s", s)
	}
	for _, qs := range strings.Split(strings.TrimSpace(halves[0]), ",") {
		r.addInput(parseQuant(qs))
	}
	r.out = parseQuant(halves[1])
	return r
}

func main() {
	react := make([]reaction, 0, 50)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		react = append(react, parseReaction(line))
	}
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}

	for _, r := range react {
		r.print()
	}
}
