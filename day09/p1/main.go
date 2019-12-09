package main

import (
	"log"
	"strconv"
	"strings"
)

var (
	relativeBase = int64(0)
)

func getArgAddr(data []int64, pos int64, offset int64, part int64) int64 {
	switch data[pos] / part % 10 {
	case 1:
		return pos + offset
	case 2:
		return relativeBase + data[pos+offset]
	}
	return data[pos+offset]
}

func getP1Addr(data []int64, pos int64) int64 {
	return getArgAddr(data, pos, 1, 100)
}

func getP2Addr(data []int64, pos int64) int64 {
	return getArgAddr(data, pos, 2, 1000)
}

func getP3Addr(data []int64, pos int64) int64 {
	return getArgAddr(data, pos, 3, 10000)
}

func process(data []int64, pos int64, inC chan int64, outC chan int64, doneC chan bool) {
	opcode := data[pos] % 100
	increase := int64(4) // default increase
	//log.Printf("PC %4d : %05d (%d,%d,%d)", pos, data[pos], data[pos+1], data[pos+2], data[pos+3])

	if opcode == 99 {
		log.Printf("END\n")
		doneC <- true
		return
	} else if opcode == 1 {
		p1 := data[getP1Addr(data, pos)]
		p2 := data[getP2Addr(data, pos)]
		data[getP3Addr(data, pos)] = p1 + p2
	} else if opcode == 2 {
		p1 := data[getP1Addr(data, pos)]
		p2 := data[getP2Addr(data, pos)]
		data[getP3Addr(data, pos)] = p1 * p2
	} else if opcode == 3 {
		data[getP1Addr(data, pos)] = <-inC
		increase = 2
	} else if opcode == 4 {
		outC <- data[getP1Addr(data, pos)]
		increase = 2
	} else if opcode == 5 {
		p1 := data[getP1Addr(data, pos)]
		if p1 != 0 {
			pos = data[getP2Addr(data, pos)]
			increase = 0
		} else {
			increase = 3
		}
	} else if opcode == 6 {
		p1 := data[getP1Addr(data, pos)]
		if p1 == 0 {
			pos = data[getP2Addr(data, pos)]
			increase = 0
		} else {
			increase = 3
		}
	} else if opcode == 7 {
		p1 := data[getP1Addr(data, pos)]
		p2 := data[getP2Addr(data, pos)]
		if p1 < p2 {
			data[getP3Addr(data, pos)] = 1
		} else {
			data[getP3Addr(data, pos)] = 0
		}
	} else if opcode == 8 {
		p1 := data[getP1Addr(data, pos)]
		p2 := data[getP2Addr(data, pos)]
		if p1 == p2 {
			data[getP3Addr(data, pos)] = 1
		} else {
			data[getP3Addr(data, pos)] = 0
		}
	} else if opcode == 9 {
		// relative base adjustment.
		relativeBase += data[getP1Addr(data, pos)]
		increase = 2
	} else {
		log.Fatalf("invalid input. data: %v opcode: %v pos: %v\n", data[pos], opcode, pos)
	}
	process(data, pos+increase, inC, outC, doneC)
}

func main() {
	input := "1102,34463338,34463338,63,1007,63,34463338,63,1005,63,53,1102,3,1,1000,109,988,209,12,9,1000,209,6,209,3,203,0,1008,1000,1,63,1005,63,65,1008,1000,2,63,1005,63,902,1008,1000,0,63,1005,63,58,4,25,104,0,99,4,0,104,0,99,4,17,104,0,99,0,0,1101,309,0,1024,1101,0,24,1002,1102,388,1,1029,1102,1,21,1019,1101,0,33,1015,1102,1,304,1025,1101,344,0,1027,1101,25,0,1003,1102,1,1,1021,1101,29,0,1012,1101,0,23,1005,1102,1,32,1007,1102,38,1,1000,1101,30,0,1016,1102,1,347,1026,1101,0,26,1010,1101,0,39,1004,1102,1,36,1011,1101,0,393,1028,1101,0,37,1013,1101,0,35,1008,1101,34,0,1001,1101,0,495,1022,1102,1,28,1018,1101,0,0,1020,1102,1,22,1006,1101,488,0,1023,1102,31,1,1009,1102,1,20,1017,1101,0,27,1014,109,10,21102,40,1,4,1008,1014,37,63,1005,63,205,1001,64,1,64,1106,0,207,4,187,1002,64,2,64,109,-18,1207,8,37,63,1005,63,227,1001,64,1,64,1106,0,229,4,213,1002,64,2,64,109,17,1207,-7,25,63,1005,63,247,4,235,1106,0,251,1001,64,1,64,1002,64,2,64,109,-8,1202,6,1,63,1008,63,29,63,1005,63,275,1001,64,1,64,1106,0,277,4,257,1002,64,2,64,109,25,1205,-6,293,1001,64,1,64,1105,1,295,4,283,1002,64,2,64,109,-4,2105,1,2,4,301,1106,0,313,1001,64,1,64,1002,64,2,64,109,-9,1208,-4,31,63,1005,63,335,4,319,1001,64,1,64,1105,1,335,1002,64,2,64,109,16,2106,0,-2,1106,0,353,4,341,1001,64,1,64,1002,64,2,64,109,-13,2102,1,-8,63,1008,63,38,63,1005,63,373,1105,1,379,4,359,1001,64,1,64,1002,64,2,64,109,9,2106,0,3,4,385,1105,1,397,1001,64,1,64,1002,64,2,64,109,-11,21107,41,42,0,1005,1014,415,4,403,1106,0,419,1001,64,1,64,1002,64,2,64,109,14,1206,-7,431,1106,0,437,4,425,1001,64,1,64,1002,64,2,64,109,-23,2107,37,-5,63,1005,63,455,4,443,1105,1,459,1001,64,1,64,1002,64,2,64,109,10,21107,42,41,-2,1005,1013,475,1105,1,481,4,465,1001,64,1,64,1002,64,2,64,2105,1,8,1001,64,1,64,1106,0,497,4,485,1002,64,2,64,109,-6,21108,43,41,8,1005,1017,517,1001,64,1,64,1106,0,519,4,503,1002,64,2,64,109,5,2101,0,-9,63,1008,63,23,63,1005,63,541,4,525,1106,0,545,1001,64,1,64,1002,64,2,64,109,-13,1201,5,0,63,1008,63,20,63,1005,63,565,1105,1,571,4,551,1001,64,1,64,1002,64,2,64,109,16,1205,4,589,4,577,1001,64,1,64,1106,0,589,1002,64,2,64,109,-16,1202,4,1,63,1008,63,23,63,1005,63,615,4,595,1001,64,1,64,1106,0,615,1002,64,2,64,109,1,2101,0,6,63,1008,63,33,63,1005,63,639,1001,64,1,64,1105,1,641,4,621,1002,64,2,64,109,8,21101,44,0,8,1008,1018,44,63,1005,63,667,4,647,1001,64,1,64,1105,1,667,1002,64,2,64,109,-7,1201,1,0,63,1008,63,39,63,1005,63,689,4,673,1106,0,693,1001,64,1,64,1002,64,2,64,109,7,2102,1,-8,63,1008,63,24,63,1005,63,715,4,699,1105,1,719,1001,64,1,64,1002,64,2,64,109,5,2108,34,-7,63,1005,63,739,1001,64,1,64,1105,1,741,4,725,1002,64,2,64,109,-22,2108,25,10,63,1005,63,763,4,747,1001,64,1,64,1106,0,763,1002,64,2,64,109,31,1206,-4,781,4,769,1001,64,1,64,1105,1,781,1002,64,2,64,109,-10,21101,45,0,5,1008,1019,47,63,1005,63,805,1001,64,1,64,1105,1,807,4,787,1002,64,2,64,109,2,21108,46,46,-3,1005,1013,825,4,813,1106,0,829,1001,64,1,64,1002,64,2,64,109,-22,2107,40,10,63,1005,63,845,1105,1,851,4,835,1001,64,1,64,1002,64,2,64,109,17,1208,-7,36,63,1005,63,871,1001,64,1,64,1105,1,873,4,857,1002,64,2,64,109,16,21102,47,1,-9,1008,1018,47,63,1005,63,899,4,879,1001,64,1,64,1106,0,899,4,64,99,21102,1,27,1,21101,0,913,0,1105,1,920,21201,1,39657,1,204,1,99,109,3,1207,-2,3,63,1005,63,962,21201,-2,-1,1,21102,1,940,0,1105,1,920,21201,1,0,-1,21201,-2,-3,1,21101,955,0,0,1105,1,920,22201,1,-1,-2,1106,0,966,21202,-2,1,-2,109,-3,2105,1,0"
	dataS := strings.Split(input, ",")

	// Convert string input to int slice w/ extra capacity.
	var data = make([]int64, len(dataS)*10)
	for idx, elem := range dataS {
		i, err := strconv.ParseInt(elem, 10, 64)
		if err != nil {
			log.Fatal("unable to parse int")
		}
		data[idx] = i
	}

	inC := make(chan int64, 5)
	outC := make(chan int64, 50)
	doneC := make(chan bool, 1)

	inC <- 1
	go func() {
		process(data, 0, inC, outC, doneC)
		close(outC)
	}()

	for x := range outC {
		log.Printf("Output: %v", x)
	}
	<-doneC
}
