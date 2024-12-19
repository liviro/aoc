package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func extractTowels(name string) ([]string, []string) {
	raw, err := os.ReadFile(name)
	if err != nil {
		panic("Cannot open file!")
	}
	sections := strings.Split(strings.TrimSpace(string(raw)), "\n\n")
	patterns := strings.Split(sections[0], ", ")
	designs := strings.Split(sections[1], "\n")
	return patterns, designs
}

func possibilities(design string, patterns []string, memo map[string]int) int {
	if poss, ok := memo[design]; ok {
		return poss
	}
	c := 0
	for _, p := range patterns {
		if strings.HasPrefix(design, p) {
			if len(design) == len(p) {
				c++
				continue
			}
			nd := strings.Clone(design)
			c += possibilities(nd[len(p):], patterns, memo)
		}
	}
	memo[design] = c
	return c
}

func solve(designs []string, patterns []string) (int, int) {
	p1 := 0
	p2 := 0
	memo := map[string]int{}
	for _, d := range designs {
		if ps := possibilities(d, patterns, memo); ps > 0 {
			p1++
			p2 += ps
		}
	}
	return p1, p2
}

func main() {
	t := time.Now()
	patterns, designs := extractTowels(os.Args[1])
	part1, part2 := solve(designs, patterns)
	fmt.Printf("Part 1: %d\n", part1)
	fmt.Printf("Part 2: %d\n", part2)
	fmt.Printf("Time elapsed: %s\n", time.Since(t))
}
