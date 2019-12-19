package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/mikehelmick/adventofcode2019/computer"
)

type pos struct {
	x, y int64
}

type prog []int64

func (data prog) testCord(x, y int64) bool {
	inC := make(chan int64, 5)
	defer close(inC)
	outC := make(chan computer.Output, 50)
	defer close(outC)
	emulator := computer.NewEmulator(data, inC, outC)
	go emulator.Execute()

	inC <- x
	inC <- y

	o := <-outC
	<-outC
	return o.Val == 1
}

func main() {
	input := "109,424,203,1,21102,1,11,0,1106,0,282,21101,18,0,0,1105,1,259,2102,1,1,221,203,1,21101,0,31,0,1105,1,282,21101,0,38,0,1105,1,259,20101,0,23,2,21202,1,1,3,21101,0,1,1,21101,57,0,0,1105,1,303,1202,1,1,222,21002,221,1,3,21001,221,0,2,21102,1,259,1,21101,0,80,0,1105,1,225,21101,0,175,2,21102,1,91,0,1106,0,303,2101,0,1,223,21001,222,0,4,21102,259,1,3,21101,225,0,2,21102,1,225,1,21102,1,118,0,1105,1,225,21002,222,1,3,21101,70,0,2,21101,0,133,0,1105,1,303,21202,1,-1,1,22001,223,1,1,21102,1,148,0,1105,1,259,2102,1,1,223,21002,221,1,4,21002,222,1,3,21102,24,1,2,1001,132,-2,224,1002,224,2,224,1001,224,3,224,1002,132,-1,132,1,224,132,224,21001,224,1,1,21101,195,0,0,105,1,109,20207,1,223,2,21002,23,1,1,21101,0,-1,3,21102,1,214,0,1106,0,303,22101,1,1,1,204,1,99,0,0,0,0,109,5,2102,1,-4,249,21202,-3,1,1,22102,1,-2,2,21201,-1,0,3,21101,0,250,0,1106,0,225,21201,1,0,-4,109,-5,2105,1,0,109,3,22107,0,-2,-1,21202,-1,2,-1,21201,-1,-1,-1,22202,-1,-2,-2,109,-3,2105,1,0,109,3,21207,-2,0,-1,1206,-1,294,104,0,99,21202,-2,1,-2,109,-3,2106,0,0,109,5,22207,-3,-4,-1,1206,-1,346,22201,-4,-3,-4,21202,-3,-1,-1,22201,-4,-1,2,21202,2,-1,-1,22201,-4,-1,1,22101,0,-2,3,21101,343,0,0,1105,1,303,1105,1,415,22207,-2,-3,-1,1206,-1,387,22201,-3,-2,-3,21202,-2,-1,-1,22201,-3,-1,3,21202,3,-1,-1,22201,-3,-1,2,21201,-4,0,1,21101,0,384,0,1105,1,303,1105,1,415,21202,-4,-1,-4,22201,-4,-3,-4,22202,-3,-2,-2,22202,-2,-4,-4,22202,-3,-2,-3,21202,-4,-1,-2,22201,-3,-2,1,21201,1,0,-4,109,-5,2106,0,0"
	dataS := strings.Split(input, ",")

	var data prog = make([]int64, len(dataS)*10)
	for idx, elem := range dataS {
		i, err := strconv.ParseInt(elem, 10, 64)
		if err != nil {
			log.Fatal("unable to parse int")
		}
		data[idx] = i
	}

	x := int64(0)
	y := int64(101)

	for {
		// Find the first value X for this Y
		for !data.testCord(x, y) {
			x++
		}
		//fmt.Printf("Row: %v, strating @ %v\n", y, x)
		// See if the box would fit if the lower left corner was at x,y
		if data.testCord(x, y-99) && data.testCord(x+99, y-99) {
			fmt.Printf("{%v, %v} fits :: %v\n", x, y-99, x*10000+(y-99))
			break
		}
		// move on to next row
		y++
	}

}
