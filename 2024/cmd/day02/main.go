package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/liviro/aoc/2024/internal/parse"
)

func extractReports(name string) ([][]int, error) {
	fp, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	s := bufio.NewScanner(fp)
	var rs [][]int
	for s.Scan() {
		raw := strings.Split(s.Text(), " ")
		r := []int{}
		for _, v := range raw {
			r = append(r, parse.MustInt(v))
		}
		rs = append(rs, r)
	}
	return rs, nil
}

func safe(report []int) bool {
	var incr bool
	if report[0] < report[1] {
		incr = true
	}
	safe := true
	for i := range report {
		if i == len(report)-1 {
			break
		}
		d := report[i] - report[i+1]
		if incr && (d > -1 || d < -3) {
			safe = false
		}
		if !incr && (d < 1 || d > 3) {
			safe = false
		}
	}
	return safe
}

func countSafe(reports [][]int) int {
	s := 0
	for _, r := range reports {
		if safe(r) {
			s++
		}
	}
	return s
}

func countDampenedSafe(reports [][]int) int {
	s := 0
	for _, r := range reports {
		dos := dampenedOptions(r)
		for _, do := range dos {
			if safe(do) {
				s++
				break
			}
		}
	}
	return s
}

func dampenedOptions(report []int) [][]int {
	rs := [][]int{}
	rs = append(rs, report)
	for i := range report {
		opt := []int{}
		for j, v := range report {
			if i != j {
				opt = append(opt, v)
			}
		}
		rs = append(rs, opt)
	}
	return rs
}

func main() {
	t := time.Now()
	reports, err := extractReports(os.Args[1])
	if err != nil {
		fmt.Printf("extractReports: %v", err)
	}
	fmt.Printf("Part 1: %d\n", countSafe(reports))
	fmt.Printf("Part 2: %d\n", countDampenedSafe(reports))
	fmt.Printf("Time elapsed: %s\n", time.Since(t))
}
