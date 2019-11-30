package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

func main() {

	var total int64

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		str := scanner.Text()
		val, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			log.Fatalf("err %v", err)
		}
		total += val
	}
	log.Printf("Answer %v \n", total)

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}
