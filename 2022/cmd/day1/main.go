package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func parseCarryTotal(raw string) int {
	rawCals := strings.Split(raw, "\n")
	t := 0
	for _, c := range rawCals {
		n, err := strconv.Atoi(c)
		if err != nil {
			panic("Could not parse int")
		}
		t += n
	}
	return t
}

func extractInput(fileName string) ([]int, error) {
	raw, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	rc := strings.Split(strings.TrimSpace(string(raw)), "\n\n")
	var cs []int
	for _, c := range rc {
		cs = append(cs, parseCarryTotal(c))
	}
	return cs, nil
}

func main() {
	cs, err := extractInput("input.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Input ingestion went wrong: ", err)
		os.Exit(1)
	}

	sort.Ints(cs)
	fmt.Println("Part 1:", cs[len(cs)-1])
	fmt.Println("Part 2:", cs[len(cs)-1]+cs[len(cs)-2]+cs[len(cs)-3])
}
