package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/mikehelmick/adventofcode2019/computer"
)

func permutations(arr []int64, ch chan []int64) {
	var helper func([]int64, int64)

	helper = func(arr []int64, n int64) {
		if n == 1 {
			tmp := make([]int64, len(arr))
			copy(tmp, arr)
			ch <- tmp
		} else {
			for i := int64(0); i < n; i++ {
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
	helper(arr, int64(len(arr)))
}

func link(input chan int64, output chan computer.Output, doneC chan bool) {
	for {
		o := <-output
		if o.Done {
			break
		}
		input <- o.Val
	}
	doneC <- true
}

func startMachine(data []int64, input chan int64, output chan computer.Output) {
	c := computer.NewEmulator(data, input, output, false)
	c.Execute()
}

func main() {
	input := "3,8,1001,8,10,8,105,1,0,0,21,34,43,60,81,94,175,256,337,418,99999,3,9,101,2,9,9,102,4,9,9,4,9,99,3,9,102,2,9,9,4,9,99,3,9,102,4,9,9,1001,9,4,9,102,3,9,9,4,9,99,3,9,102,4,9,9,1001,9,2,9,1002,9,3,9,101,4,9,9,4,9,99,3,9,1001,9,4,9,102,2,9,9,4,9,99,3,9,102,2,9,9,4,9,3,9,1001,9,2,9,4,9,3,9,102,2,9,9,4,9,3,9,101,1,9,9,4,9,3,9,1002,9,2,9,4,9,3,9,1001,9,2,9,4,9,3,9,101,2,9,9,4,9,3,9,101,1,9,9,4,9,3,9,102,2,9,9,4,9,3,9,1002,9,2,9,4,9,99,3,9,101,2,9,9,4,9,3,9,1001,9,1,9,4,9,3,9,101,1,9,9,4,9,3,9,102,2,9,9,4,9,3,9,1002,9,2,9,4,9,3,9,1001,9,1,9,4,9,3,9,1001,9,2,9,4,9,3,9,101,1,9,9,4,9,3,9,1001,9,1,9,4,9,3,9,1001,9,2,9,4,9,99,3,9,1002,9,2,9,4,9,3,9,102,2,9,9,4,9,3,9,1002,9,2,9,4,9,3,9,102,2,9,9,4,9,3,9,1001,9,1,9,4,9,3,9,1002,9,2,9,4,9,3,9,1001,9,2,9,4,9,3,9,101,1,9,9,4,9,3,9,101,2,9,9,4,9,3,9,1001,9,2,9,4,9,99,3,9,101,2,9,9,4,9,3,9,101,1,9,9,4,9,3,9,101,1,9,9,4,9,3,9,101,1,9,9,4,9,3,9,1001,9,1,9,4,9,3,9,1001,9,2,9,4,9,3,9,1001,9,2,9,4,9,3,9,102,2,9,9,4,9,3,9,1001,9,2,9,4,9,3,9,1002,9,2,9,4,9,99,3,9,1001,9,1,9,4,9,3,9,1001,9,1,9,4,9,3,9,1002,9,2,9,4,9,3,9,1001,9,1,9,4,9,3,9,1002,9,2,9,4,9,3,9,102,2,9,9,4,9,3,9,1001,9,2,9,4,9,3,9,1001,9,2,9,4,9,3,9,101,2,9,9,4,9,3,9,102,2,9,9,4,9,99"
	dataS := strings.Split(input, ",")

	// Convert string input to int array.
	var data = []int64{}
	for _, elem := range dataS {
		i, err := strconv.ParseInt(elem, 10, 64)
		if err != nil {
			log.Fatal("unable to parse int")
		}
		data = append(data, i)
	}

	// Span a function to generare all permutations of the phase sequences.
	permCh := make(chan []int64, 1000)
	a := []int64{9, 8, 7, 6, 5}
	go func() {
		permutations(a, permCh)
		permCh <- []int64{0}
	}()

	// Read back the channel of phase permutations and run each one.
	maxThurst := int64(0)
	for {
		perm := <-permCh
		log.Printf("Received: %v", perm)
		if len(perm) == 1 {
			log.Printf("Max Thrust: %v", maxThurst)
			break
		}

		// separate copy of the program for each machine & input channel.
		channels := make([]chan int64, 5)
		outChannels := make([]chan computer.Output, 5)
		for i := 0; i < 5; i++ {
			channels[i] = make(chan int64, 20)
			channels[i] <- perm[i]
			if i == 0 {
				// input to 1st machine.
				channels[i] <- 0
			}
			outChannels[i] = make(chan computer.Output, 20)
		}

		// single done output channel. We'll just wait for 5 dones.
		doneC := make(chan bool, 5)

		// start and link the parallel machines
		go startMachine(data, channels[0], outChannels[0])
		go link(channels[1], outChannels[0], doneC)
		go startMachine(data, channels[1], outChannels[1])
		go link(channels[2], outChannels[1], doneC)
		go startMachine(data, channels[2], outChannels[2])
		go link(channels[3], outChannels[2], doneC)
		go startMachine(data, channels[3], outChannels[3])
		go link(channels[4], outChannels[3], doneC)
		go startMachine(data, channels[4], outChannels[4])
		go link(channels[0], outChannels[4], doneC)

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
