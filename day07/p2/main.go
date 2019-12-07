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

func process(data []int, pos int, inC chan int, outC chan int, doneC chan bool) {
	opcode := data[pos] % 100
	//log.Printf("pc: %v :: %v\n", pos, data[pos])
	increase := 4 // default increase

	if opcode == 99 {
		log.Printf("END\n")
		doneC <- true
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
		log.Printf("Read Input %v", input)
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
	process(data, pos+increase, inC, outC, doneC)
}

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

	// Convert string input to int array.
	var data = []int{}
	for _, elem := range dataS {
		i, err := strconv.ParseInt(elem, 10, 64)
		if err != nil {
			log.Fatal("unable to parse int")
		}
		data = append(data, int(i))
	}

	// Span a function to generare all permutations of the phase sequences.
	permCh := make(chan []int, 1000)
	a := []int{9, 8, 7, 6, 5}
	go func() {
		permutations(a, permCh)
		permCh <- []int{0}
	}()

	// Read back the channel of phase permutations and run each one.
	maxThurst := 0
	for {
		perm := <-permCh
		log.Printf("Received: %v", perm)
		if len(perm) == 1 {
			log.Printf("Max Thrust: %v", maxThurst)
			break
		}

		// separate copy of the program for each machine & input channel.
		programs := make([][]int, 5)
		channels := make([]chan int, 5)
		for i := 0; i < 5; i++ {
			programs[i] = make([]int, len(data))
			copy(programs[i], data)
			channels[i] = make(chan int, 20)
			channels[i] <- perm[i]
			if i == 0 {
				// input to 1st machine.
				channels[i] <- 0
			}
		}

		// single done output channel. We'll just wait for 5 dones.
		doneC := make(chan bool, 5)

		// start and link the parallel machines
		go process(programs[0], 0, channels[0], channels[1], doneC)
		go process(programs[1], 0, channels[1], channels[2], doneC)
		go process(programs[2], 0, channels[2], channels[3], doneC)
		go process(programs[3], 0, channels[3], channels[4], doneC)
		go process(programs[4], 0, channels[4], channels[0], doneC)

		for i := 0; i < 5; i++ {
			<-doneC
		}
		// programs[4] will have written it's last output to p1's channel.
		// p1 should have already shut down, fingers crossed.
		thrust := <-channels[0]
		if thrust > maxThurst {
			maxThurst = thrust
		}
	}
}
