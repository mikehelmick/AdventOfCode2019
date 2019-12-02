package main

import (
	"log"
	"strconv"
	"strings"
)

func process(data []int64, pos int) {
	if data[pos] == 99 {
		return
	} else if data[pos] == 1 {
		data[data[pos+3]] = data[data[pos+1]] + data[data[pos+2]]
	} else if data[pos] == 2 {
		data[data[pos+3]] = data[data[pos+1]] * data[data[pos+2]]
	} else {
		log.Fatal("invalid input")
	}
	process(data, pos+4)
}

func main() {
	input := "1,0,0,3,1,1,2,3,1,3,4,3,1,5,0,3,2,1,10,19,1,9,19,23,1,13,23,27,1,5,27,31,2,31,6,35,1,35,5,39,1,9,39,43,1,43,5,47,1,47,5,51,2,10,51,55,1,5,55,59,1,59,5,63,2,63,9,67,1,67,5,71,2,9,71,75,1,75,5,79,1,10,79,83,1,83,10,87,1,10,87,91,1,6,91,95,2,95,6,99,2,99,9,103,1,103,6,107,1,13,107,111,1,13,111,115,2,115,9,119,1,119,6,123,2,9,123,127,1,127,5,131,1,131,5,135,1,135,5,139,2,10,139,143,2,143,10,147,1,147,5,151,1,151,2,155,1,155,13,0,99,2,14,0,0"
	dataS := strings.Split(input, ",")

	var data = []int64{}
	for _, elem := range dataS {
		i, err := strconv.ParseInt(elem, 10, 64)
		if err != nil {
			log.Fatal("unable to parse int")
		}
		data = append(data, i)
	}

	tmp := make([]int64, len(data))

	for i := int64(0); i <= 99; i++ {
		for j := int64(0); j <= 99; j++ {
			copy(tmp, data)
			tmp[1] = i
			tmp[2] = j
			process(tmp, 0)

			if tmp[0] == 19690720 {
				log.Printf("Answer %v and %v \n", i, j)
				ans := 100*i + j
				log.Printf("Answer %v \n", ans)
			}
		}
	}
}
