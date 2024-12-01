package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/liviro/aoc/2024/internal/parse"
)

func extractLists(name string) ([]int, []int, error) {
	fp, err := os.Open(name)
	if err != nil {
		return nil, nil, err
	}
	defer fp.Close()
	s := bufio.NewScanner(fp)
	var l1, l2 []int
	for s.Scan() {
		r := strings.Split(s.Text(), "   ")
		l1 = append(l1, parse.MustInt(r[0]))
		l2 = append(l2, parse.MustInt(r[1]))
	}
	slices.Sort(l1)
	slices.Sort(l2)
	return l1, l2, nil
}

func distance(l1, l2 []int) int {
	res := 0
	for i, a := range l1 {
		res += int(math.Abs(float64(a - l2[i])))
	}
	return res
}

func similarity(l1, l2 []int) int {
	res := 0
	d := make(map[int]int)
	for _, a := range l2 {
		d[a] += 1
	}
	for _, b := range l1 {
		if c, ok := d[b]; ok {
			res += b * c
		}
	}
	return res
}

func main() {
	t := time.Now()
	l1, l2, err := extractLists(os.Args[1])
	if err != nil {
		fmt.Printf("extractLists: %v", err)
		return
	}
	fmt.Printf("Part 1: %d\n", distance(l1, l2))
	fmt.Printf("Part 2: %d\n", similarity(l1, l2))
	fmt.Printf("Time elapsed: %s\n", time.Since(t))
}
