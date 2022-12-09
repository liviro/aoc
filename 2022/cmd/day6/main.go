package main

import (
	"fmt"
	"os"
)

func charsDistinct(s string) bool {
	m := make(map[rune]struct{})
	for _, r := range s {
		if _, ok := m[r]; ok {
			return false
		}
		m[r] = struct{}{}
	}
	return true
}

func markerIndex(buf string, n int) int {
	for i := n; i < len(buf); i++ {
		if charsDistinct(buf[i-n : i]) {
			return i
		}
	}
	return -1
}

func main() {
	in := "input.txt"
	buf, err := os.ReadFile(in)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Input ingestion went wrong: ", err)
		os.Exit(1)
	}
	fmt.Println("Part 1:", markerIndex(string(buf), 4))
	fmt.Println("Part 2:", markerIndex(string(buf), 14))
}
