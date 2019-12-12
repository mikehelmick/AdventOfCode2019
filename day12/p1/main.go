package main

import "log"

type pos struct {
	x, y, z int32
}

type moon struct {
	pos [3]int32
	vel [3]int32
}

func updateGravity(a, b *moon) {
	for i := 0; i < len(b.pos); i++ {
		if a.pos[i] < b.pos[i] {
			a.vel[i]++
			b.vel[i]--
		} else if a.pos[i] > b.pos[i] {
			a.vel[i]--
			b.vel[i]++
		}
	}
}

func (m moon) toPos() pos {
	return pos{m.pos[0], m.pos[1], m.pos[2]}
}

func (m moon) toVel() pos {
	return pos{m.vel[0], m.vel[1], m.vel[2]}
}

func (m *moon) move() {
	for i := 0; i < len(m.pos); i++ {
		m.pos[i] += m.vel[i]
	}
}

func abs(a int32) int32 {
	if a < 0 {
		return a * -1
	}
	return a
}

func (m moon) energy() int32 {
	return (abs(m.pos[0]) + abs(m.pos[1]) + abs(m.pos[2])) * (abs(m.vel[0]) + abs(m.vel[1]) + abs(m.vel[2]))
}

func (m moon) print() {
	log.Printf("pos=<x=%3d, y=%3d, z=%3d>, vel=<x=%3d, y=%3d, z=%3d>", m.pos[0], m.pos[1], m.pos[2], m.vel[0], m.vel[1], m.vel[2])
}

func (m moon) init(x, y, z int32) moon {
	m.pos[0] = x
	m.pos[1] = y
	m.pos[2] = z
	return m
}

func printAll(moons []moon, s int) {
	log.Printf("After %v steps", s)
	for _, m := range moons {
		m.print()
	}
	log.Printf(" ")
}

func energy(moons []moon) int32 {
	e := int32(0)
	for _, m := range moons {
		e += m.energy()
	}
	return e
}

func axisEquals(init, moons []moon, pos int) bool {
	res := true
	for i := 0; i < len(init) && res; i++ {
		res = init[i].pos[pos] == moons[i].pos[pos] &&
			init[i].vel[pos] == moons[i].vel[pos]
	}
	return res
}

func gcd(a, b int64) int64 {
	for b != 0 {
		tmp := b
		b = a % b
		a = tmp
	}
	return a
}

func lcm(a, b int64, rest ...int64) int64 {
	result := a * b / gcd(a, b)

	for i := 0; i < len(rest); i++ {
		result = lcm(result, rest[i])
	}

	return result
}

func main() {
	moons := make([]moon, 4)
	moons[0] = moons[0].init(1, 4, 4)
	moons[1] = moons[1].init(-4, -1, 19)
	moons[2] = moons[2].init(-15, -14, 12)
	moons[3] = moons[3].init(-17, 1, 10)
	//moons[0] = moons[0].init(-1, 0, 2)
	//moons[1] = moons[1].init(2, -10, -7)
	//moons[2] = moons[2].init(4, -8, 8)
	//moons[3] = moons[3].init(3, 5, -1)

	iState := make([]moon, 4)
	copy(iState, moons)

	printAll(moons, 0)
	for s := 1; s < 1001; s++ {
		for a := 0; a < len(moons)-1; a++ {
			for b := a + 1; b < len(moons); b++ {
				updateGravity(&moons[a], &moons[b])
			}
		}
		for i := range moons {
			(&moons[i]).move()
		}

		printAll(moons, s)
	}

	energy := int32(0)
	for _, m := range moons {
		energy += m.energy()
	}
	log.Printf("Energy %v", energy)
}
