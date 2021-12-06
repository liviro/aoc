package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type fish struct {
	ticker int32
}

// parseFish returns a new fish, as parsed from the input string of its tick.
func parseFish(raw string) (fish, error) {
	n, err := strconv.Atoi(raw)
	if err != nil {
		return fish{}, err
	}
	return fish{int32(n)}, nil
}

// tick returns a ticked copy of the input fish, and, if applicable, a pointer to a newly spawned fish.
// New fish are spawned if the passed in fish is at tick 0.
func tick(f fish) (oldF fish, newF *fish) {
	switch f.ticker {
	case 0:
		return fish{6}, &fish{8}
	default:
		return fish{f.ticker - 1}, nil
	}
}

type school map[fish]int64

// extractSchool extracts the input school of fish in the given file name.
func extractSchool(fileName string) (school, error) {
	raw, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	rs := strings.Split(strings.TrimSpace(string(raw)), ",")
	sch := make(school)
	for _, s := range rs {
		f, err := parseFish(s)
		if err != nil {
			return nil, err
		}
		sch[f] += 1
	}
	return sch, nil
}

// size computes the number of fishes in the given school.
func (s school) size() int64 {
	var c int64
	for _, v := range s {
		c += v
	}
	return c
}

// evolve returns the evolved (by one tick) school of the given input.
func evolve(s school) school {
	ns := make(school)
	for k, v := range s {
		of, nf := tick(k)
		ns[of] += v
		if nf != nil {
			ns[*nf] += v
		}
	}
	return ns
}

// laterSchoolSize returns the size of the input school after the given number of ticks.
func laterSchoolSize(s school, d int) int64 {
	ns := s
	for i := 0; i < d; i++ {
		ns = evolve(ns)
	}
	return ns.size()
}

func main() {
	s, err := extractSchool("input.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Input ingestion went wrong: ", err)
		os.Exit(1)
	}
	fmt.Println("Part 1:", laterSchoolSize(s, 80))
	fmt.Println("Part 2:", laterSchoolSize(s, 256))
}
