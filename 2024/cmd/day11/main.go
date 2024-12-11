package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/liviro/aoc/2024/internal/parse"
)

func extractStones(name string) []int {
	fp, err := os.Open(name)
	if err != nil {
		panic("Unable to open file")
	}
	defer fp.Close()
	s := bufio.NewScanner(fp)
	res := []int{}
	for s.Scan() {
		for _, r := range strings.Split(s.Text(), " ") {
			res = append(res, parse.MustInt(r))
		}
	}
	return res
}

// In a memoized way, return the number of stones that the input stone S
// will yield after the given number of blinks.
// Memoization is a map of stone -> (map of blink # -> resulting stones' count).
func evolveStones(s int, blinks int, memo map[int]map[int]int) int {
	// Base case: no blinking, no change in stone.
	if blinks == 0 {
		return 1
	}

	// Memoization table lookup
	_, ok := memo[s]
	if ok {
		if s, ok := memo[s][blinks]; ok {
			return s
		}
	}
	// Initialize the map if no entry for stone is not found.
	if !ok {
		memo[s] = map[int]int{}
	}

	// 0 -> 1
	if s == 0 {
		next := evolveStones(1, blinks-1, memo)
		memo[s][blinks] = next
		return memo[s][blinks]
	}
	// Even digits -> split into 2 stones
	if len(fmt.Sprintf("%d", s))%2 == 0 {
		orig := fmt.Sprintf("%d", s)
		s1 := parse.MustInt(orig[:len(orig)/2])
		s2 := parse.MustInt(orig[len(orig)/2:])
		c1 := evolveStones(s1, blinks-1, memo)
		c2 := evolveStones(s2, blinks-1, memo)
		memo[s][blinks] = c1 + c2
		return memo[s][blinks]
	}
	// x -> s*2024
	next := evolveStones(s*2024, blinks-1, memo)
	memo[s][blinks] = next
	return memo[s][blinks]
}

func countStones(stones []int, blinks int) int {
	memo := map[int]map[int]int{}
	sum := 0
	for _, s := range stones {
		sum += evolveStones(s, blinks, memo)
	}
	return sum
}

func main() {
	t := time.Now()
	stones := extractStones(os.Args[1])
	fmt.Printf("Part 1: %d\n", countStones(stones, 25))
	fmt.Printf("Part 2: %d\n", countStones(stones, 75))
	fmt.Printf("Time elapsed: %s\n", time.Since(t))
}
