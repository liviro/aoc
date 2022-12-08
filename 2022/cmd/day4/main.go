package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type pairs struct {
	oneStart int
	oneEnd   int
	twoStart int
	twoEnd   int
}

func mustParseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic("Could not parse int!")
	}
	return i
}

// Input format: a-b,c-d
func parsePair(r string) pairs {
	ps := strings.Split(r, ",")
	one := strings.Split(ps[0], "-")
	two := strings.Split(ps[1], "-")
	return pairs{
		oneStart: mustParseInt(one[0]),
		oneEnd:   mustParseInt(one[1]),
		twoStart: mustParseInt(two[0]),
		twoEnd:   mustParseInt(two[1]),
	}
}

func extractPairs(fileName string) ([]pairs, error) {
	fp, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	s := bufio.NewScanner(fp)
	var ps []pairs
	for s.Scan() {
		ps = append(ps, parsePair(s.Text()))
	}
	return ps, nil
}

func (p pairs) isFullyContained() bool {
	return p.oneStart <= p.twoStart && p.oneEnd >= p.twoEnd || // One fully contains two
		p.twoStart <= p.oneStart && p.twoEnd >= p.oneEnd // Two fully contins one
}

func (p pairs) hasOverlap() bool {
	return !(p.oneEnd < p.twoStart || p.twoEnd < p.oneStart)
}

func part1(ps []pairs) int {
	c := 0
	for _, p := range ps {
		if p.isFullyContained() {
			c++
		}
	}
	return c
}

func part2(ps []pairs) int {
	c := 0
	for _, p := range ps {
		if p.hasOverlap() {
			c++
		}
	}
	return c
}

func main() {
	in := "input.txt"
	ps, err := extractPairs(in)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Input ingestion went wrong: ", err)
		os.Exit(1)
	}
	fmt.Println("Part 1:", part1(ps))
	fmt.Println("Part 2:", part2(ps))
}
