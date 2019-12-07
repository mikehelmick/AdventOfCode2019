package main

import (
	"log"
	"strconv"
	"strings"
)

func getP1(data []int, pos int) int {
	if data[pos]/100%10 == 1 {
		return data[pos+1]
	}
	return data[data[pos+1]]
}

func getP2(data []int, pos int) int {
	if data[pos]/1000%10 == 1 {
		return data[pos+2]
	}
	return data[data[pos+2]]
}

func process(data []int, pos int, inC chan int, outC chan int) {
	opcode := data[pos] % 100
	//log.Printf("pc: %v :: %v\n", pos, data[pos])
	increase := 4 // default increase

	if opcode == 99 {
		//log.Printf("END\n")
		return
	} else if opcode == 1 {
		p1 := getP1(data, pos)
		p2 := getP2(data, pos)
		//log.Printf("Addition: %v -- %v + %v", data[pos], p1, p2)
		data[data[pos+3]] = p1 + p2
	} else if opcode == 2 {
		p1 := getP1(data, pos)
		p2 := getP2(data, pos)
		data[data[pos+3]] = p1 * p2
		//log.Printf("Multiplication: %v -- %v * %v", data[pos], p1, p2)
	} else if opcode == 3 {
		input := <-inC
		data[data[pos+1]] = input
		increase = 2
	} else if opcode == 4 {
		val := getP1(data, pos)
		log.Printf("Output value %v\n", val)
		outC <- val
		increase = 2
	} else if opcode == 5 {
		p1 := getP1(data, pos)
		if p1 != 0 {
			pos = getP2(data, pos)
			increase = 0
		} else {
			increase = 3
		}
	} else if opcode == 6 {
		p1 := getP1(data, pos)
		if p1 == 0 {
			pos = getP2(data, pos)
			increase = 0
		} else {
			increase = 3
		}
	} else if opcode == 7 {
		p1 := getP1(data, pos)
		p2 := getP2(data, pos)
		if p1 < p2 {
			data[data[pos+3]] = 1
		} else {
			data[data[pos+3]] = 0
		}
	} else if opcode == 8 {
		p1 := getP1(data, pos)
		p2 := getP2(data, pos)
		if p1 == p2 {
			data[data[pos+3]] = 1
		} else {
			data[data[pos+3]] = 0
		}
	} else {
		log.Fatalf("invalid input. data: %v opcode: %v pos: %v\n", data[pos], opcode, pos)
	}
	process(data, pos+increase, inC, outC)
}

// Generate all permutations of an array, send them to an output channel
func permutations(arr []int, ch chan []int) {
	var helper func([]int, int)

	helper = func(arr []int, n int) {
		if n == 1 {
			tmp := make([]int, len(arr))
			copy(tmp, arr)
			ch <- tmp
		} else {
			for i := 0; i < n; i++ {
				helper(arr, n-1)
				if n%2 == 1 {
					tmp := arr[i]
					arr[i] = arr[n-1]
					arr[n-1] = tmp
				} else {
					tmp := arr[0]
					arr[0] = arr[n-1]
					arr[n-1] = tmp
				}
			}
		}
	}
	helper(arr, len(arr))
}

func main() {
	input := "3,8,1001,8,10,8,105,1,0,0,21,34,43,60,81,94,175,256,337,418,99999,3,9,101,2,9,9,102,4,9,9,4,9,99,3,9,102,2,9,9,4,9,99,3,9,102,4,9,9,1001,9,4,9,102,3,9,9,4,9,99,3,9,102,4,9,9,1001,9,2,9,1002,9,3,9,101,4,9,9,4,9,99,3,9,1001,9,4,9,102,2,9,9,4,9,99,3,9,102,2,9,9,4,9,3,9,1001,9,2,9,4,9,3,9,102,2,9,9,4,9,3,9,101,1,9,9,4,9,3,9,1002,9,2,9,4,9,3,9,1001,9,2,9,4,9,3,9,101,2,9,9,4,9,3,9,101,1,9,9,4,9,3,9,102,2,9,9,4,9,3,9,1002,9,2,9,4,9,99,3,9,101,2,9,9,4,9,3,9,1001,9,1,9,4,9,3,9,101,1,9,9,4,9,3,9,102,2,9,9,4,9,3,9,1002,9,2,9,4,9,3,9,1001,9,1,9,4,9,3,9,1001,9,2,9,4,9,3,9,101,1,9,9,4,9,3,9,1001,9,1,9,4,9,3,9,1001,9,2,9,4,9,99,3,9,1002,9,2,9,4,9,3,9,102,2,9,9,4,9,3,9,1002,9,2,9,4,9,3,9,102,2,9,9,4,9,3,9,1001,9,1,9,4,9,3,9,1002,9,2,9,4,9,3,9,1001,9,2,9,4,9,3,9,101,1,9,9,4,9,3,9,101,2,9,9,4,9,3,9,1001,9,2,9,4,9,99,3,9,101,2,9,9,4,9,3,9,101,1,9,9,4,9,3,9,101,1,9,9,4,9,3,9,101,1,9,9,4,9,3,9,1001,9,1,9,4,9,3,9,1001,9,2,9,4,9,3,9,1001,9,2,9,4,9,3,9,102,2,9,9,4,9,3,9,1001,9,2,9,4,9,3,9,1002,9,2,9,4,9,99,3,9,1001,9,1,9,4,9,3,9,1001,9,1,9,4,9,3,9,1002,9,2,9,4,9,3,9,1001,9,1,9,4,9,3,9,1002,9,2,9,4,9,3,9,102,2,9,9,4,9,3,9,1001,9,2,9,4,9,3,9,1001,9,2,9,4,9,3,9,101,2,9,9,4,9,3,9,102,2,9,9,4,9,99"
	dataS := strings.Split(input, ",")

	var data = []int{}
	for _, elem := range dataS {
		i, err := strconv.ParseInt(elem, 10, 64)
		if err != nil {
			log.Fatal("unable to parse int")
		}
		data = append(data, int(i))
	}

	inCh := make(chan int, 2)
	outCh := make(chan int)
	permCh := make(chan []int, 100)
	maxThurst := 0
	a := []int{0, 1, 2, 3, 4}
	go func() {
		permutations(a, permCh)
		permCh <- []int{0}
	}()

	tmp := make([]int, len(data))
	for {
		perm := <-permCh
		log.Printf("Received: %v", perm)
		if len(perm) == 1 {
			log.Printf("Max Thrust: %v", maxThurst)
			break
		}

		thrust := 0 // initial parm to machine 1
		for i := 0; i < 5; i++ {
			log.Printf("IN: %v %v", perm[i], thrust)
			inCh <- perm[i]
			inCh <- thrust
			copy(tmp, data)
			go process(tmp, 0, inCh, outCh)
			thrust = <-outCh
		}
		//log.Printf("Max thrust %v from input %v", thrust, perm)
		if thrust > maxThurst {
			maxThurst = thrust
		}
	}
}
