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

func oreForFuel(react map[string]reaction, fuel int64) int64 {
	// Amount of an element on hand.
	bank := make(map[string]int64)
	ore := int64(0)
	need := make(map[string]int64)
	need["FUEL"] = fuel

	//for len(need) > 0 {
	for len(need) > 0 { //round := 0; round < 3; round++ {
		//log.Printf("NEEDS: %v", need)
		//log.Printf("Ore: %v | Bank: %v", ore, bank)
		// k is target element, v is number needed of that element.
		nextRnd := make(map[string]int64)
		for k, v := range need {
			//log.Printf("need to produce %v of %v", v, k)
			if banked, ok := bank[k]; ok {
				//log.Printf(" There is %v of %v in the bank", banked, k)
				if banked >= v {
					v = 0
					bank[k] = banked - v
				} else {
					v -= banked
					delete(bank, k)
				}
			}
			// Able to completly satisfy this need from the bank. yay.
			if v == 0 {
				continue
			}

			r := react[k]
			//log.Printf("Reaction %v", r)
			times := v / r.out.amt
			if v%r.out.amt > 0 {
				times++
			}

			if len(r.in) == 1 && r.in[0].elem == "ORE" {
				inQ := r.in[0]
				ore += (inQ.amt * times)
			} else {
				// Add needs for the input to the reaction
				for _, quant := range r.in {
					if cur, ok := nextRnd[quant.elem]; ok {
						nextRnd[quant.elem] = cur + quant.amt*times
					} else {
						nextRnd[quant.elem] = quant.amt * times
					}
				}
			}
			v -= times * r.out.amt
			// Bank the escess V of K
			if v < 0 {
				if banked, ok := bank[k]; ok {
					bank[k] = banked + -v
				} else {
					bank[k] = -v
				}
			}
		}
		need = nextRnd
		//log.Printf("Ore: %v | Bank %v", ore, bank)
		//log.Printf("----------")
	}
	return ore
}

// Binary search over a range to find the fuel produced
func fuelForOre(react map[string]reaction, ore int64, high int64) int64 {
	sRange := int64(0)
	eRange := high
	lastResult := int64(0)
	guess := int64(1)
	for rnd := 0; rnd < 1000; rnd++ {
		guess = (eRange-sRange)/2 + sRange
		if guess == lastResult {
			break
		}
		o := oreForFuel(react, guess)
		if o == ore {
			return guess
		}

		if o > ore {
			eRange = guess
		} else {
			sRange = guess
		}
		lastResult = guess
		log.Printf("%v Fuel used %v ore", guess, o)
	}
	return guess
}

func main() {
	react := make(map[string]reaction)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		r := parseReaction(line)
		if _, ok := react[r.out.elem]; ok {
			log.Fatalf("more than 1 way to produce %v", r.out.elem)
			return
		}
		react[r.out.elem] = r
	}
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}

	for _, r := range react {
		r.print()
	}
	ore := oreForFuel(react, 1)
	log.Printf("1 fuel Consumed %v ORE", ore)
	log.Printf("----")

	maxOre := int64(1000000000000)
	highRange := maxOre / ore * 2
	fuel := fuelForOre(react, maxOre, highRange)
	log.Printf("%v ORE produces %v fuel", maxOre, fuel)
}
