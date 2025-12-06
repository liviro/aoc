package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/liviro/aoc/2025/internal/parse"
)

type problem struct {
	operator string
	numbers  []int
}

func (p problem) solve() int {
	var r int
	if p.operator == "*" {
		r = 1
		for _, v := range p.numbers {
			r *= v
		}
	}
	if p.operator == "+" {
		for _, v := range p.numbers {
			r += v
		}
	}
	return r
}

func extractProblemsPt1(name string) []problem {
	fp, err := os.Open(name)
	if err != nil {
		panic("Unable to open file")
	}
	defer fp.Close()
	s := bufio.NewScanner(fp)
	var ps []problem
	for s.Scan() {
		raw := slices.DeleteFunc(strings.Split(s.Text(), " "), func(s string) bool {
			return s == ""
		})
		if len(ps) == 0 {
			ps = make([]problem, len(raw))
		}
		for i, v := range raw {
			if v == "+" {
				ps[i].operator = "+"
				continue
			}
			if v == "*" {
				ps[i].operator = "*"
				continue
			}

			ps[i].numbers = append(ps[i].numbers, parse.MustInt(v))
		}
	}
	return ps
}

func extractProblemsPt2(name string) []problem {
	fp, err := os.Open(name)
	if err != nil {
		panic("Unable to open file")
	}
	defer fp.Close()
	s := bufio.NewScanner(fp)
	var rawLines []string
	for s.Scan() {
		rawLines = append(rawLines, s.Text())
	}
	var ps []problem
	p := problem{}
	for i := len(rawLines[0]) - 1; i >= 0; i-- {
		// If all-blank column, move onto next problem
		allBlank := true
		for y := 0; y < len(rawLines); y++ {
			if rawLines[y][i] != ' ' {
				allBlank = false
			}
		}
		if allBlank {
			ps = append(ps, p)
			p = problem{}
			continue
		}

		var rn strings.Builder
		for j := 0; j < len(rawLines)-1; j++ {
			if rawLines[j][i] != ' ' {
				rn.WriteByte(rawLines[j][i])
			}
		}
		p.numbers = append(p.numbers, parse.MustInt(rn.String()))

		if rawLines[len(rawLines)-1][i] != ' ' {
			p.operator = string(rawLines[len(rawLines)-1][i])
		}
	}
	ps = append(ps, p)
	return ps
}

func grandTotal(ps []problem) int {
	s := 0
	for _, p := range ps {
		s += p.solve()
	}
	return s
}

func main() {
	t := time.Now()
	ps1 := extractProblemsPt1(os.Args[1])
	fmt.Printf("Part 1: %d\n", grandTotal(ps1))
	ps2 := extractProblemsPt2(os.Args[1])
	fmt.Printf("Part 2: %d\n", grandTotal(ps2))
	fmt.Printf("Time elapsed: %s\n", time.Since(t))
}
