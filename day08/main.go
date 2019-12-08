package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const width = 25
const height = 6
const pixels = width * height

// layers are 25*6 = 150 pixels
type layer struct {
	data  [pixels]int32
	zeros int
	ones  int
	twos  int
}

func (l *layer) Summary() {
	for _, v := range l.data {
		switch v {
		case 0:
			l.zeros++
		case 1:
			l.ones++
		case 2:
			l.twos++
		}
	}
}

func (l layer) Print() {
	for i, val := range l.data {
		fmt.Printf("%v", val)
		if (i+1)%width == 0 {
			fmt.Printf("\n")
		}
	}
}

func (l *layer) Calculate(layers []layer) {
	for i := range l.data {
		l.data[i] = pixelValue(i, layers)
	}
}

func pixelValue(i int, layers []layer) int32 {
	for _, l := range layers {
		if l.data[i] != 2 {
			return l.data[i]
		}
	}
	return 0
}

func main() {
	inputText := ""
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		inputText = scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}

	layers := make([]layer, len(inputText)/pixels)
	layerIdx := 0
	charIdx := 0
	for _, ch := range inputText {
		layers[layerIdx].data[charIdx] = int32(ch - '0')
		if charIdx++; charIdx == pixels {
			layerIdx++
			charIdx = 0
		}
	}

	minIdx := 0
	for i, l := range layers {
		l.Summary()
		if l.zeros < layers[minIdx].zeros {
			minIdx = i
		}
		layers[i] = l
	}
	log.Printf("p1 checksum: %v", layers[minIdx].ones*layers[minIdx].twos)
	log.Printf("p2 final message")
	final := new(layer)
	final.Calculate(layers)
	final.Print()
}
