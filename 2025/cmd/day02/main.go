package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/liviro/aoc/2025/internal/parse"
)

type interval struct {
	start, end int
}

func extractRanges(name string) []interval {
	fp, err := os.Open(name)
	if err != nil {
		panic("Unable to open file")
	}
	defer fp.Close()
	s := bufio.NewScanner(fp)
	res := []interval{}
	for s.Scan() {
		for _, r := range strings.Split(s.Text(), ",") {
			start, end, found := strings.Cut(r, "-")
			if !found {
				panic("Not a proper interval!")
			}
			res = append(res, interval{
				start: parse.MustInt(start),
				end:   parse.MustInt(end),
			})
		}
	}
	return res
}

func isInvalidPt1(n int) bool {
	s := fmt.Sprint(n)
	if len(s)%2 != 0 {
		return false
	}
	a := parse.MustInt(s[:len(s)/2])
	b := parse.MustInt(s[len(s)/2:])
	return a == b
}

func isInvalidPt2(n int) bool {
	s := fmt.Sprint(n)
L:
	for l := 2; l <= len(s); l++ {
		if len(s)%l != 0 {
			continue
		}
		x := parse.MustInt(s[:len(s)/l])
		rest := s[len(s)/l:]
		for {
			if len(rest) == 0 {
				return true
			}
			y := parse.MustInt(rest[:len(s)/l])
			if x != y {
				continue L
			}
			rest = rest[len(s)/l:]
		}
	}
	return false
}

func invalidSums(its []interval) (int, int) {
	p1, p2 := 0, 0
	for _, it := range its {
		for a := it.start; a <= it.end; a++ {
			if isInvalidPt1(a) {
				p1 += a
			}
			if isInvalidPt2(a) {
				p2 += a
			}
		}
	}
	return p1, p2
}

func main() {
	t := time.Now()
	its := extractRanges(os.Args[1])
	p1, p2 := invalidSums(its)
	fmt.Printf("Part 1: %d\n", p1)
	fmt.Printf("Part 2: %d\n", p2)
	fmt.Printf("Time elapsed: %s\n", time.Since(t))
}
