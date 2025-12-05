package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/liviro/aoc/2025/internal/parse"
)

type fresh struct{ min, max int }

func (f fresh) String() string {
	return fmt.Sprintf("[%d, %d]", f.min, f.max)
}

func parseFresh(raw string) []fresh {
	var fs []fresh
	for _, r := range strings.Split(raw, "\n") {
		ps := strings.Split(r, "-")
		fs = append(fs, fresh{
			min: parse.MustInt(ps[0]),
			max: parse.MustInt(ps[1]),
		})
	}
	return fs
}

func parseIngredients(raw string) []int {
	var is []int
	for r := range strings.SplitSeq(raw, "\n") {
		is = append(is, parse.MustInt(r))
	}
	return is
}

func extractDatabase(name string) ([]fresh, []int) {
	raw, err := os.ReadFile(name)
	if err != nil {
		panic("Unable to open file")
	}
	sections := strings.Split(strings.TrimSpace(string(raw)), "\n\n")
	return parseFresh(sections[0]), parseIngredients(sections[1])
}

func isFresh(freshRanges []fresh, i int) bool {
	for _, r := range freshRanges {
		if i >= r.min && i <= r.max {
			return true
		}
	}
	return false
}

func part1(freshRanges []fresh, ingredients []int) int {
	c := 0
	for _, i := range ingredients {
		if isFresh(freshRanges, i) {
			c++
		}
	}
	return c
}

func unify(freshRanges []fresh) []fresh {
	var nf []fresh
	for _, r := range freshRanges {
		if len(nf) == 0 {
			nf = append(nf, r)
			continue
		}
		for i, c := range nf {
			// Non-overlapping, smaller than ranges so far
			if r.max < c.min {
				nf = slices.Insert(nf, i, r)
				break
			}
			// r is fully enclosed by c
			if r.min >= c.min && r.max <= c.max {
			}
			// Overlapping, r has only lower elements
			if r.min < c.min && r.max >= c.min {
				nf[i].min = r.min
			}
			// Overlapping, r has only higher elements
			if r.min <= c.max && r.max > c.max {
				nf[i].max = r.max
			}
			// c is fully enclosed by r
			if r.min < c.min && r.max > c.max {
				nf[i] = r
			}
		}
		// r bigger than all ranges so far
		if r.min > nf[len(nf)-1].max {
			nf = append(nf, r)
		}
	}
	return nf
}

func part2(freshRanges []fresh) int {
	fr := freshRanges
	for {
		oldLength := len(fr)
		fr = unify(fr)
		if oldLength == len(fr) {
			break
		}
	}
	c := 0
	for _, r := range fr {
		c += r.max - r.min + 1
	}
	return c
}

func main() {
	t := time.Now()
	fs, is := extractDatabase(os.Args[1])
	fmt.Printf("Part 1: %d\n", part1(fs, is))
	fmt.Printf("Part 2: %d\n", part2(fs))
	fmt.Printf("Time elapsed: %s\n", time.Since(t))
}
