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

	var data = []int{}
	for _, elem := range dataS {
		i, err := strconv.ParseInt(elem, 10, 64)
		if err != nil {
			log.Fatal("unable to parse int")
		}
		data = append(data, int(i))
	}

	permCh := make(chan []int, 100)
	maxThurst := 0
	a := []int{9, 8, 7, 6, 5}
	go func() {
		permutations(a, permCh)
		permCh <- []int{0}
	}()

	for {
		perm := <-permCh
		log.Printf("Received: %v", perm)
		if len(perm) == 1 {
			log.Printf("Max Thrust: %v", maxThurst)
			break
		}

		// programming competition style...
		// separate copy of the program for each machine.
		p1 := make([]int, len(data))
		copy(p1, data)
		p2 := make([]int, len(data))
		copy(p2, data)
		p3 := make([]int, len(data))
		copy(p3, data)
		p4 := make([]int, len(data))
		copy(p4, data)
		p5 := make([]int, len(data))
		copy(p5, data)

		// input channel for each machine.
		ch1 := make(chan int, 20)
		ch1 <- perm[0]
		ch1 <- 0 // initial input
		ch2 := make(chan int, 20)
		ch2 <- perm[1]
		ch3 := make(chan int, 20)
		ch3 <- perm[2]
		ch4 := make(chan int, 20)
		ch4 <- perm[3]
		ch5 := make(chan int, 20)
		ch5 <- perm[4]

		// single done output channel. We'll just wait for 5 dones.
		doneC := make(chan bool, 5)

		// start and link the parallel machines
		go process(p1, 0, ch1, ch2, doneC)
		go process(p2, 0, ch2, ch3, doneC)
		go process(p3, 0, ch3, ch4, doneC)
		go process(p4, 0, ch4, ch5, doneC)
		go process(p5, 0, ch5, ch1, doneC)

		for i := 0; i < 5; i++ {
			<-doneC
		}
		// p5 will have written it's last output to p1's channel.
		// p1 should have already shut down, fingers crossed.
		thrust := <-ch1
		if thrust > maxThurst {
			maxThurst = thrust
		}
	}
}
