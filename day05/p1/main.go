package main

import (
	"log"
	"strconv"
	"strings"
)

// Get the param value based on the mode for P1
func getP1(data []int, pos int) int {
	if data[pos]/100%10 == 1 {
		return data[pos+1]
	}
	return data[data[pos+1]]
}

// Get the param value based on the mode for P1
func getP2(data []int, pos int) int {
	if data[pos]/1000%10 == 1 {
		return data[pos+2]
	}
	return data[data[pos+2]]
}

func process(data []int, pos int) {
	opcode := data[pos] % 100
	//log.Printf("pc: %v :: %v\n", pos, data[pos])
	increase := 4

	if opcode == 99 {
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
		input := 1 // by definition
		data[data[pos+1]] = input
		increase = 2
	} else if opcode == 4 {
		val := getP1(data, pos)
		log.Printf("Output value %v\n", val)
		increase = 2
	} else {
		log.Fatalf("invalid input. data: %v opcode: %v pos: %v\n", data[pos], opcode, pos)
	}
	process(data, pos+increase)
}

func main() {
	input := "3,225,1,225,6,6,1100,1,238,225,104,0,1102,27,28,225,1,113,14,224,1001,224,-34,224,4,224,102,8,223,223,101,7,224,224,1,224,223,223,1102,52,34,224,101,-1768,224,224,4,224,1002,223,8,223,101,6,224,224,1,223,224,223,1002,187,14,224,1001,224,-126,224,4,224,102,8,223,223,101,2,224,224,1,224,223,223,1102,54,74,225,1101,75,66,225,101,20,161,224,101,-54,224,224,4,224,1002,223,8,223,1001,224,7,224,1,224,223,223,1101,6,30,225,2,88,84,224,101,-4884,224,224,4,224,1002,223,8,223,101,2,224,224,1,224,223,223,1001,214,55,224,1001,224,-89,224,4,224,102,8,223,223,1001,224,4,224,1,224,223,223,1101,34,69,225,1101,45,67,224,101,-112,224,224,4,224,102,8,223,223,1001,224,2,224,1,223,224,223,1102,9,81,225,102,81,218,224,101,-7290,224,224,4,224,1002,223,8,223,101,5,224,224,1,223,224,223,1101,84,34,225,1102,94,90,225,4,223,99,0,0,0,677,0,0,0,0,0,0,0,0,0,0,0,1105,0,99999,1105,227,247,1105,1,99999,1005,227,99999,1005,0,256,1105,1,99999,1106,227,99999,1106,0,265,1105,1,99999,1006,0,99999,1006,227,274,1105,1,99999,1105,1,280,1105,1,99999,1,225,225,225,1101,294,0,0,105,1,0,1105,1,99999,1106,0,300,1105,1,99999,1,225,225,225,1101,314,0,0,106,0,0,1105,1,99999,1007,677,677,224,102,2,223,223,1005,224,329,101,1,223,223,1108,226,677,224,1002,223,2,223,1005,224,344,101,1,223,223,1008,677,677,224,102,2,223,223,1005,224,359,101,1,223,223,8,226,677,224,1002,223,2,223,1006,224,374,101,1,223,223,108,226,677,224,1002,223,2,223,1006,224,389,1001,223,1,223,1107,226,677,224,102,2,223,223,1005,224,404,1001,223,1,223,7,226,677,224,1002,223,2,223,1005,224,419,101,1,223,223,1107,677,226,224,102,2,223,223,1006,224,434,1001,223,1,223,1107,226,226,224,1002,223,2,223,1006,224,449,101,1,223,223,1108,226,226,224,1002,223,2,223,1005,224,464,101,1,223,223,8,677,226,224,102,2,223,223,1005,224,479,101,1,223,223,8,226,226,224,1002,223,2,223,1006,224,494,1001,223,1,223,1007,226,677,224,1002,223,2,223,1006,224,509,1001,223,1,223,108,226,226,224,1002,223,2,223,1006,224,524,1001,223,1,223,1108,677,226,224,102,2,223,223,1006,224,539,101,1,223,223,1008,677,226,224,102,2,223,223,1006,224,554,101,1,223,223,107,226,677,224,1002,223,2,223,1006,224,569,101,1,223,223,107,677,677,224,102,2,223,223,1006,224,584,101,1,223,223,7,677,226,224,102,2,223,223,1005,224,599,101,1,223,223,1008,226,226,224,1002,223,2,223,1005,224,614,1001,223,1,223,107,226,226,224,1002,223,2,223,1005,224,629,101,1,223,223,7,226,226,224,102,2,223,223,1006,224,644,1001,223,1,223,1007,226,226,224,102,2,223,223,1006,224,659,101,1,223,223,108,677,677,224,102,2,223,223,1005,224,674,1001,223,1,223,4,223,99,226"
	dataS := strings.Split(input, ",")

	var data = []int{}
	for _, elem := range dataS {
		i, err := strconv.ParseInt(elem, 10, 32)
		if err != nil {
			log.Fatal("unable to parse int")
		}
		data = append(data, int(i))
	}

	process(data, 0)
	//log.Printf("Answer %v \n", data)
}
