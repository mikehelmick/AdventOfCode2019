package computer

import (
	"log"
	"time"
)

// Output is a single piece of output from the emulator. If done is false
// then it is an integer output from the machine. If done is true, the
// machine has terminated.
type Output struct {
	Val  int64
	Done bool
}

// Emulator representing a single instance of the intcode computer.
type Emulator struct {
	mem          []int64
	pc           int64
	relativeBase int64

	nonBlockingInput bool

	input  chan int64
	output chan Output
}

const positionMode = 0
const immediateMode = 1
const relativeMode = 2

// NewEmulator initializes a new intcode computer
func NewEmulator(program []int64, input chan int64, output chan Output, nonBlockingInput bool) *Emulator {
	emu := new(Emulator)
	emu.mem = make([]int64, len(program)*10)
	copy(emu.mem, program)
	emu.input = input
	emu.output = output
	emu.nonBlockingInput = nonBlockingInput
	return emu
}

// Debug prints currents state to the log.
func (c *Emulator) Debug() {
	log.Printf("PC: %05d mem: %v base: %v", c.pc, c.mem[c.pc], c.relativeBase)
}

func (c *Emulator) getArgAddr(offset int64, part int64) int64 {
	switch c.mem[c.pc] / part % 10 {
	case immediateMode:
		return c.pc + offset
	case relativeMode:
		return c.relativeBase + c.mem[c.pc+offset]
	}
	return c.mem[c.pc+offset]
}

func (c *Emulator) getP1Addr() int64 {
	return c.getArgAddr(1, 100)
}

func (c *Emulator) getP2Addr() int64 {
	return c.getArgAddr(2, 1000)
}

func (c *Emulator) getP3Addr() int64 {
	return c.getArgAddr(3, 10000)
}

// Execute runs the program sent in on initializion until completion.
func (c *Emulator) Execute() {
	for {
		opcode := c.mem[c.pc] % 100
		increase := int64(4) // default increase
		//log.Printf("PC %4d : %05d (%d,%d,%d)", c.pc, c.mem[c.pc], c.mem[c.pc+1], c.mem[c.pc+2], c.mem[c.pc+3])

		switch opcode {
		case 1:
			c.mem[c.getP3Addr()] = c.mem[c.getP1Addr()] + c.mem[c.getP2Addr()]
		case 2:
			c.mem[c.getP3Addr()] = c.mem[c.getP1Addr()] * c.mem[c.getP2Addr()]
		case 3:
			if c.nonBlockingInput {
				select {
				case val := <-c.input:
					c.mem[c.getP1Addr()] = val
				case <-time.After(100 * time.Millisecond):
					c.mem[c.getP1Addr()] = -1
				}
			} else {
				c.mem[c.getP1Addr()] = <-c.input
			}
			increase = 2
		case 4:
			c.output <- Output{c.mem[c.getP1Addr()], false}
			increase = 2
		case 5:
			p1 := c.mem[c.getP1Addr()]
			if p1 != 0 {
				c.pc = c.mem[c.getP2Addr()]
				increase = 0
			} else {
				increase = 3
			}
		case 6:
			p1 := c.mem[c.getP1Addr()]
			if p1 == 0 {
				c.pc = c.mem[c.getP2Addr()]
				increase = 0
			} else {
				increase = 3
			}
		case 7:
			p1 := c.mem[c.getP1Addr()]
			p2 := c.mem[c.getP2Addr()]
			if p1 < p2 {
				c.mem[c.getP3Addr()] = 1
			} else {
				c.mem[c.getP3Addr()] = 0
			}
		case 8:
			p1 := c.mem[c.getP1Addr()]
			p2 := c.mem[c.getP2Addr()]
			if p1 == p2 {
				c.mem[c.getP3Addr()] = 1
			} else {
				c.mem[c.getP3Addr()] = 0
			}
		case 9:
			// relative base adjustment.
			c.relativeBase += c.mem[c.getP1Addr()]
			increase = 2
		case 99:
			//log.Printf("END\n")
			c.output <- Output{0, true}
			return
		default:
			log.Fatalf("invalid input. data: %v opcode: %v pos: %v\n", c.mem[c.pc], opcode, c.pc)
		}
		c.pc += increase
	}
}
