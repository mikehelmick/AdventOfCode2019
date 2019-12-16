package main

import (
	"fmt"
	"log"
	"strconv"
)

const multiplier = 10000

// Convert the input string to an int array and copy it multiplier times.
func toArray(s string) []int64 {
	res := make([]int64, len(s)*multiplier)

	starter := make([]int64, len(s))
	for i := range s {
		val, err := strconv.ParseInt(s[i:i+1], 10, 64)
		if err != nil {
			log.Fatal("unable to parse int")
		}
		starter[i] = val
	}

	for i := 0; i < multiplier; i++ {
		copy(res[i*len(s):], starter)
	}

	return res
}

// Sum a slice
func sum(in []int64) int64 {
	sum := int64(0)
	for _, val := range in {
		sum += val
	}
	return sum
}

func runPhase(in []int64, start int64) []int64 {
	out := make([]int64, len(in))

	// Mask is "1" for start to end because of size of offset.
	// Add up start
	// As you go, each elemnt is the sum - the sum of the values before it
	sum := sum(in[start:])
	rem := int64(0)

	limit := int64(len(in))
	for i := start; i < limit; i++ {
		out[i] = (sum - rem) % 10
		rem += in[i]
	}
	return out
}

func main() {
	inputS := "59754835304279095723667830764559994207668723615273907123832849523285892960990393495763064170399328763959561728553125232713663009161639789035331160605704223863754174835946381029543455581717775283582638013183215312822018348826709095340993876483418084566769957325454646682224309983510781204738662326823284208246064957584474684120465225052336374823382738788573365821572559301715471129142028462682986045997614184200503304763967364026464055684787169501819241361777789595715281841253470186857857671012867285957360755646446993278909888646724963166642032217322712337954157163771552371824741783496515778370667935574438315692768492954716331430001072240959235708"

	input := toArray(inputS)

	toSkip, err := strconv.ParseInt(inputS[0:7], 10, 64)
	if err != nil {
		panic("Unable to parse first 7 chars.")
	}

	log.Printf("input size: %v startingAt: %v", len(input), toSkip)
	for i := 0; i < 100; i++ {
		input = runPhase(input, toSkip)
		log.Printf("done with phase %v", i+1)
	}

	//for _, v := range input {
	//	fmt.Printf("%d", v)
	//}
	fmt.Println("")
	log.Printf("Offset %v", toSkip)

	log.Printf("Answer: %v", input[toSkip:toSkip+8])
}
