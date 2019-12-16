package main

import (
	"fmt"
	"log"
	"strconv"
)

func toArray(s string) []int64 {
	res := make([]int64, len(s))
	for i := range s {
		val, err := strconv.ParseInt(s[i:i+1], 10, 64)
		if err != nil {
			log.Fatal("unable to parse int")
		}
		res[i] = val
	}
	return res
}

var base = []int64{0, 1, 0, -1}

func pattern(pos int) []int64 {
	patt := make([]int64, pos*4)
	for bi, val := range base {
		for i := 0; i < pos; i++ {
			patt[bi*pos+i] = val
		}
	}
	return patt
}

func runPhase(in []int64) []int64 {
	out := make([]int64, len(in))

	// i is position being calculated
	for i := range in {
		pattern := pattern(i + 1)
		//log.Printf("idx: %v pattern: %v", i, pattern)
		runner := 1

		oVal := int64(0)
		for _, val := range in {
			//log.Printf("  %v * %v", pattern[runner], val)
			oVal = oVal + (pattern[runner] * val)
			runner++
			if runner >= len(pattern) {
				runner = 0
			}
		}
		if oVal < 0 {
			oVal = -oVal
		}
		if oVal >= 10 {
			oVal = oVal % 10
		}
		out[i] = oVal
		//log.Printf("idx: %v newVal: %v", i, oVal)
	}

	return out
}

func main() {
	inputS := "59754835304279095723667830764559994207668723615273907123832849523285892960990393495763064170399328763959561728553125232713663009161639789035331160605704223863754174835946381029543455581717775283582638013183215312822018348826709095340993876483418084566769957325454646682224309983510781204738662326823284208246064957584474684120465225052336374823382738788573365821572559301715471129142028462682986045997614184200503304763967364026464055684787169501819241361777789595715281841253470186857857671012867285957360755646446993278909888646724963166642032217322712337954157163771552371824741783496515778370667935574438315692768492954716331430001072240959235708"

	input := toArray(inputS)

	log.Printf("%v", input)
	for i := 0; i < 100; i++ {
		input = runPhase(input)
	}
	log.Printf("%v", input)

	for _, v := range input {
		fmt.Printf("%d", v)
	}
	fmt.Println("")
	fmt.Printf("Answer: %v\n", input[0:8])

}
