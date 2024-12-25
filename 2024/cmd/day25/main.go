package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func extractSchematics(name string) ([][5]int, [][5]int) {
	raw, err := os.ReadFile(name)
	if err != nil {
		panic("Cannot open file!")
	}
	sections := strings.Split(strings.TrimSpace(string(raw)), "\n\n")
	locks := [][5]int{}
	keys := [][5]int{}
	for _, s := range sections {
		if strings.HasPrefix(s, "#") {
			locks = append(locks, parseLock(s))
		} else {
			keys = append(keys, parseKey(s))
		}
	}
	return locks, keys
}

func parseLock(raw string) [5]int {
	l := [5]int{}
	rows := strings.Split(raw, "\n")
	for i := 1; i < len(rows); i++ {
		for j, v := range strings.Split(rows[i], "") {
			if v == "#" {
				l[j]++
			}
		}
	}
	return l
}

func parseKey(raw string) [5]int {
	l := [5]int{}
	rows := strings.Split(raw, "\n")
	for i := len(rows) - 2; i > 0; i-- {
		for j, v := range strings.Split(rows[i], "") {
			if v == "#" {
				l[j]++
			}
		}
	}
	return l
}

func hasOverlap(lock, key [5]int) bool {
	for i := range lock {
		if lock[i]+key[i] > 5 {
			return true
		}
	}
	return false
}

func part1(locks, keys [][5]int) int {
	c := 0
	for _, l := range locks {
		for _, k := range keys {
			if !hasOverlap(l, k) {
				c++
			}
		}
	}
	return c
}

func main() {
	t := time.Now()
	locks, keys := extractSchematics(os.Args[1])
	fmt.Printf("Part 1: %d\n", part1(locks, keys))
	fmt.Printf("Part2: %d\n", 0)
	fmt.Printf("Time elapsed: %s\n", time.Since(t))
}
