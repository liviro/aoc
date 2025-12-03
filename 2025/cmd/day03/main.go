package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/liviro/aoc/2025/internal/parse"
)

func extractGrid(name string) [][]int {
	fp, err := os.Open(name)
	if err != nil {
		panic("Unable to open file")
	}
	defer fp.Close()
	s := bufio.NewScanner(fp)
	grid := [][]int{}
	for s.Scan() {
		bank := []int{}
		for _, b := range strings.Split(s.Text(), "") {
			bank = append(bank, parse.MustInt(b))
		}
		grid = append(grid, bank)
	}
	return grid
}

func maxJoltage(row []int, digits int) int {
	res := 0
	idx := -1
	for d := digits - 1; d >= 0; d-- {
		x := 0
		for i := idx + 1; i < len(row)-d; i++ {
			if row[i] > x {
				x = row[i]
				idx = i
			}
		}
		res += x * int(math.Pow10(d))
	}
	return res
}

func joltages(grid [][]int) (int, int) {
	s1 := 0
	s2 := 0
	for _, b := range grid {
		s1 += maxJoltage(b, 2)
		s2 += maxJoltage(b, 12)
	}
	return s1, s2
}

func main() {
	t := time.Now()
	grid := extractGrid(os.Args[1])
	p1, p2 := joltages(grid)
	fmt.Printf("Part 1: %d\n", p1)
	fmt.Printf("Part 2: %d\n", p2)
	fmt.Printf("Time elapsed: %s\n", time.Since(t))
}
