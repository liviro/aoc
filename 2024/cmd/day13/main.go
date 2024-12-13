package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type button struct {
	x, y int
}

type machine struct {
	a, b button
	x, y int
}

func extractMachines(name string) []machine {
	raw, err := os.ReadFile(name)
	if err != nil {
		panic("Cannot open file!")
	}
	rawMachines := strings.Split(strings.TrimSpace(string(raw)), "\n\n")
	ms := []machine{}
	for _, rm := range rawMachines {
		ms = append(ms, parseMachine(rm))
	}
	return ms
}

func parseMachine(raw string) machine {
	m := machine{
		a: button{},
		b: button{},
	}
	rows := strings.Split(raw, "\n")
	fmt.Sscanf(rows[0], "Button A: X+%d, Y+%d", &m.a.x, &m.a.y)
	fmt.Sscanf(rows[1], "Button B: X+%d, Y+%d", &m.b.x, &m.b.y)
	fmt.Sscanf(rows[2], "Prize: X=%d, Y=%d", &m.x, &m.y)
	return m
}

func fixMachines(ms []machine) []machine {
	fms := []machine{}
	for _, m := range ms {
		fms = append(fms, machine{
			a: m.a, b: m.b,
			x: 10000000000000 + m.x,
			y: 10000000000000 + m.y,
		})
	}
	return fms
}

func tokens(m machine) int {
	dn := m.a.x*m.b.y - m.a.y*m.b.x
	in := (m.x*m.b.y - m.y*m.b.x)
	jn := (m.y*m.a.x - m.x*m.a.y)
	if in%dn != 0 || jn%dn != 0 {
		return -1
	}
	if in*dn <= 0 || jn*dn <= 0 {
		return -1
	}
	return 3*(in/dn) + (jn / dn)
}

func countTokens(ms []machine) int {
	s := 0
	for _, m := range ms {
		mt := tokens(m)
		if mt != -1 {
			s += mt
		}
	}
	return s
}

func main() {
	t := time.Now()
	ms := extractMachines(os.Args[1])
	fmt.Printf("Part 1: %d\n", countTokens(ms))
	fmt.Printf("Part 2: %d\n", countTokens(fixMachines(ms)))
	fmt.Printf("Time elapsed: %s\n", time.Since(t))
}
