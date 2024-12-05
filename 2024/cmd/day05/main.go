package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/liviro/aoc/2024/internal/parse"
)

type rule struct {
	before, after int
}

func parseRules(raw string) []rule {
	rs := []rule{}
	for _, r := range strings.Split(raw, "\n") {
		ps := strings.Split(r, "|")
		rs = append(rs, rule{
			before: parse.MustInt(ps[0]),
			after:  parse.MustInt(ps[1]),
		})
	}
	return rs
}

func parseUpdates(raw string) [][]int {
	us := [][]int{}
	for _, r := range strings.Split(raw, "\n") {
		ps := strings.Split(r, ",")
		u := []int{}
		for _, p := range ps {
			u = append(u, parse.MustInt(p))
		}
		us = append(us, u)
	}
	return us
}

func extractRulesAndUpdates(name string) ([]rule, [][]int, error) {
	raw, err := os.ReadFile(name)
	if err != nil {
		return nil, nil, err
	}
	sections := strings.Split(strings.TrimSpace(string(raw)), "\n\n")
	return parseRules(sections[0]), parseUpdates(sections[1]), nil
}

// We can safely swap if we see a rule being broken: all previous
// indexes i will still be before (so same rule) as the following
// ones.
func maybeCorrect(update []int, rules []rule) bool {
	hadUpdate := false
	for i := range update {
		for j := i + 1; j < len(update); j++ {
			if slices.Contains(rules, rule{update[j], update[i]}) {
				update[i], update[j] = update[j], update[i]
				hadUpdate = true
			}
		}
	}
	return hadUpdate
}

func sums(updates [][]int, rules []rule) (int, int) {
	corrects := 0
	wrongs := 0
	for _, u := range updates {
		wasUpdated := maybeCorrect(u, rules)
		if wasUpdated {
			wrongs += u[len(u)/2]
		} else {
			corrects += u[len(u)/2]
		}
	}
	return corrects, wrongs
}

func main() {
	t := time.Now()
	rs, us, err := extractRulesAndUpdates(os.Args[1])
	if err != nil {
		fmt.Printf("extractRulesAndUpdates: %v", err)
	}
	correctSums, wrongSums := sums(us, rs)
	fmt.Printf("Part 1: %d\n", correctSums)
	fmt.Printf("Part 2: %d\n", wrongSums)
	fmt.Printf("Time elapsed: %s\n", time.Since(t))
}
