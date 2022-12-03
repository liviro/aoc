package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func findShared(sack string) rune {
	m := make(map[rune]struct{})
	for _, b := range sack[:len(sack)/2] {
		m[b] = struct{}{}
	}
	for _, b := range sack[len(sack)/2:] {
		if _, ok := m[b]; ok {
			return b
		}
	}
	panic("Nothing shared!")
}

func findBadge(sacks []string) rune {
	m := make(map[rune]struct{})
	for _, b := range sacks[0] {
		m[b] = struct{}{}
	}
	for k, _ := range m {
		if strings.ContainsRune(sacks[1], k) && strings.ContainsRune(sacks[2], k) {
			return k
		}
	}
	panic("Expected a badge item!")
}

func prio(r rune) int {
	if 'A' <= r && r <= 'Z' {
		return int(r-'A') + 27
	}
	return int(r-'a') + 1
}

func part1(sacks []string) int {
	res := 0
	for _, s := range sacks {
		res += prio(findShared(s))
	}
	return res
}

func part2(sacks []string) int {
	res := 0
	for i := 0; i < len(sacks); i += 3 {
		res += prio(findBadge(sacks[i : i+3]))
	}
	return res
}

func extractSacks(fileName string) ([]string, error) {
	fp, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	s := bufio.NewScanner(fp)
	var sacks []string
	for s.Scan() {
		sacks = append(sacks, s.Text())
	}
	return sacks, nil
}

func main() {
	in := "input.txt"
	sacks, err := extractSacks(in)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Input ingestion went wrong: ", err)
		os.Exit(1)
	}
	fmt.Println("Part 1:", part1(sacks))
	fmt.Println("Part 2:", part2(sacks))
}
