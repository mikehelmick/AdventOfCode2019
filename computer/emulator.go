package computer

import "log"

// Emulator representing a single instance of the intcode computer.
type Emulator struct {
	mem          []int64
	pc           int64
	relativeBase int64

	input  chan int64
	output chan int64
	done   chan bool
}

// NewEmulator initializes a new intcode computer
func NewEmulator(program []int64, input, output chan int64, done chan bool) *Emulator {
	emu := new(Emulator)
	emu.mem = make([]int64, len(program)*10)
	copy(emu.mem, program)
	emu.input = input
	emu.output = output
	emu.done = done
	return emu
}

// Debug prints currents state to the log.
func (c *Emulator) Debug() {
	log.Printf("PC: %05d mem: %v base: %v", c.pc, c.mem[c.pc], c.relativeBase)
}

func (c *Emulator) getArgAddr(offset int64, part int64) int64 {
	switch c.mem[c.pc] / part % 10 {
	case 1:
		return c.pc + offset
	case 2:
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
		//log.Printf("PC %4d : %05d (%d,%d,%d)", pos, c.mem[pos], c.mem[pos+1], c.mem[pos+2], c.mem[pos+3])

		if opcode == 99 {
			log.Printf("END\n")
			c.done <- true
			return
		} else if opcode == 1 {
			c.mem[c.getP3Addr()] = c.mem[c.getP1Addr()] + c.mem[c.getP2Addr()]
		} else if opcode == 2 {
			c.mem[c.getP3Addr()] = c.mem[c.getP1Addr()] * c.mem[c.getP2Addr()]
		} else if opcode == 3 {
			c.mem[c.getP1Addr()] = <-c.input
			increase = 2
		} else if opcode == 4 {
			c.output <- c.mem[c.getP1Addr()]
			increase = 2
		} else if opcode == 5 {
			p1 := c.mem[c.getP1Addr()]
			if p1 != 0 {
				c.pc = c.mem[c.getP2Addr()]
				increase = 0
			} else {
				increase = 3
			}
		} else if opcode == 6 {
			p1 := c.mem[c.getP1Addr()]
			if p1 == 0 {
				c.pc = c.mem[c.getP2Addr()]
				increase = 0
			} else {
				increase = 3
			}
		} else if opcode == 7 {
			p1 := c.mem[c.getP1Addr()]
			p2 := c.mem[c.getP2Addr()]
			if p1 < p2 {
				c.mem[c.getP3Addr()] = 1
			} else {
				c.mem[c.getP3Addr()] = 0
			}
		} else if opcode == 8 {
			p1 := c.mem[c.getP1Addr()]
			p2 := c.mem[c.getP2Addr()]
			if p1 == p2 {
				c.mem[c.getP3Addr()] = 1
			} else {
				c.mem[c.getP3Addr()] = 0
			}
		} else if opcode == 9 {
			// relative base adjustment.
			c.relativeBase += c.mem[c.getP1Addr()]
			increase = 2
		} else {
			log.Fatalf("invalid input. data: %v opcode: %v pos: %v\n", c.mem[c.pc], opcode, c.pc)
		}
		c.pc += increase
	}
}
