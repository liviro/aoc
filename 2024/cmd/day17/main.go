package main

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/liviro/aoc/2024/internal/parse"
)

type machine struct {
	a, b, c int
	program []int
	pointer int
	output  []int
}

func (m machine) copy() machine {
	nm := machine{
		a: m.a, b: m.b, c: m.c, pointer: m.pointer,
		program: []int{},
		output:  []int{},
	}
	nm.program = append(nm.program, m.program...)
	nm.output = append(nm.output, m.output...)
	return nm
}

func (m machine) combo() int {
	op := m.program[m.pointer+1]
	if op < 4 {
		return op
	}
	if op == 4 {
		return m.a
	}
	if op == 5 {
		return m.b
	}
	if op == 6 {
		return m.c
	}
	panic("Invalid combo operand!")
}

func (m machine) literal() int {
	return m.program[m.pointer+1]
}

func (m *machine) adv() {
	num := m.a
	den := math.Pow(2, float64(m.combo()))
	m.a = int(math.Trunc(float64(num) / den))
	m.pointer += 2
}

func (m *machine) bxl() {
	m.b = m.b ^ m.literal()
	m.pointer += 2
}

func (m *machine) bst() {
	m.b = m.combo() % 8
	m.pointer += 2
}

func (m *machine) jnz() {
	if m.a == 0 {
		m.pointer += 2
		return
	}
	m.pointer = m.literal()
}

func (m *machine) bxc() {
	m.b = m.b ^ m.c
	m.pointer += 2
}

func (m *machine) out() {
	m.output = append(m.output, m.combo()%8)
	m.pointer += 2
}

func (m *machine) bdv() {
	num := m.a
	den := math.Pow(2, float64(m.combo()))
	m.b = int(math.Trunc(float64(num) / den))
	m.pointer += 2
}

func (m *machine) cdv() {
	num := m.a
	den := math.Pow(2, float64(m.combo()))
	m.c = int(math.Trunc(float64(num) / den))
	m.pointer += 2
}

func (m *machine) run() {
	for {
		if m.pointer >= len(m.program) {
			return
		}
		switch {
		case m.program[m.pointer] == 0:
			m.adv()
		case m.program[m.pointer] == 1:
			m.bxl()
		case m.program[m.pointer] == 2:
			m.bst()
		case m.program[m.pointer] == 3:
			m.jnz()
		case m.program[m.pointer] == 4:
			m.bxc()
		case m.program[m.pointer] == 5:
			m.out()
		case m.program[m.pointer] == 6:
			m.bdv()
		case m.program[m.pointer] == 7:
			m.cdv()
		}
	}
}

func (m machine) printOutput() string {
	var b strings.Builder
	for _, v := range m.output {
		b.WriteString(fmt.Sprintf("%d,", v))
	}
	return b.String()[:b.Len()-1]
}

func extractMachine(name string) machine {
	raw, err := os.ReadFile(name)
	if err != nil {
		panic("Cannot open file!")
	}
	sections := strings.Split(strings.TrimSpace(string(raw)), "\n\n")
	m := machine{program: []int{}}
	fmt.Sscanf(sections[0], "Register A: %d\nRegister B: %d\nRegister C: %d", &m.a, &m.b, &m.c)
	for _, i := range strings.Split(strings.TrimPrefix(sections[1], "Program: "), ",") {
		m.program = append(m.program, parse.MustInt(i))
	}
	return m
}

func to10bit(a []int) int {
	r := 0
	for i, v := range a {
		r += v * int(math.Pow(8, float64(len(a)-(i+1))))
	}
	return r
}

func part2(m machine) int {
	soFar := []int{}
	t := 0
	for {
		n := -1
	D:
		for ; t < 8; t++ {
			mc := m.copy()
			mc.a = to10bit(append(soFar, t))
			mc.run()

			for j := 1; j <= len(soFar)+1; j++ {
				if m.program[len(m.program)-j] != mc.output[len(mc.output)-j] {
					continue D
				}
			}
			n = t
			break D
		}
		if n != -1 {
			soFar = append(soFar, n)
			t = 0
		} else {
			// Need to roll back
		R:
			for {
				last := soFar[len(soFar)-1]
				soFar = soFar[:len(soFar)-1]
				if last != 7 {
					t = last + 1
					break R
				}
			}
		}
		if len(soFar) == len(m.program) {
			break
		}
	}
	return to10bit(soFar)
}

func main() {
	t := time.Now()
	machine := extractMachine(os.Args[1])
	m1 := machine.copy()
	m1.run()
	fmt.Printf("Part 1: %s\n", m1.printOutput())
	fmt.Printf("Part 2: %d\n", part2(machine.copy()))
	fmt.Printf("Time elapsed: %s\n", time.Since(t))
}
