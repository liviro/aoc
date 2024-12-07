package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/liviro/aoc/2024/internal/parse"
)

type equation struct {
	testVal int
	nums    []int
}

func extractEquations(name string) []equation {
	fp, err := os.Open(name)
	if err != nil {
		panic("Unable to open file")
	}
	defer fp.Close()
	s := bufio.NewScanner(fp)
	eqs := []equation{}
	for s.Scan() {
		r := strings.Split(s.Text(), ": ")
		e := equation{
			testVal: parse.MustInt(r[0]),
			nums:    []int{},
		}
		for _, n := range strings.Split(r[1], " ") {
			e.nums = append(e.nums, parse.MustInt(n))
		}
		eqs = append(eqs, e)
	}
	return eqs
}

func concat(a, b int) int {
	return parse.MustInt(fmt.Sprintf("%d%d", a, b))
}

func digits(a int) int {
	return int(math.Ceil((math.Log10(float64(a) + 1))))
}

func possible(e equation, withConcat bool) bool {
	if len(e.nums) == 2 {
		if e.nums[0]+e.nums[1] == e.testVal ||
			e.nums[0]*e.nums[1] == e.testVal ||
			(withConcat && concat(e.nums[0], e.nums[1]) == e.testVal) {
			return true
		}
		return false
	}
	last := e.nums[len(e.nums)-1]
	if e.testVal%last == 0 {
		if possible(equation{
			testVal: e.testVal / last,
			nums:    append([]int(nil), e.nums[:len(e.nums)-1]...),
		}, withConcat) {
			return true
		}
	}
	if e.testVal-last >= 0 {
		if possible(equation{
			testVal: e.testVal - last,
			nums:    append([]int(nil), e.nums[:len(e.nums)-1]...),
		}, withConcat) {
			return true
		}
	}
	if withConcat && ((e.testVal-last)%int(math.Pow10(digits(last))) == 0) {
		if possible(equation{
			testVal: (e.testVal - last) / int(math.Pow10(digits(last))),
			nums:    append([]int(nil), e.nums[:len(e.nums)-1]...),
		}, true) {
			return true
		}
	}
	return false
}

func sumPossible(eqs []equation, withConcat bool) int {
	s := 0
	for _, e := range eqs {
		if possible(e, withConcat) {
			s += e.testVal
		}
	}
	return s
}

func main() {
	t := time.Now()
	eqs := extractEquations(os.Args[1])
	fmt.Printf("Part 1: %d\n", sumPossible(eqs, false))
	fmt.Printf("Part 2: %d\n", sumPossible(eqs, true))
	fmt.Printf("Time elapsed: %s\n", time.Since(t))
}
