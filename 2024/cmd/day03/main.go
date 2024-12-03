package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/liviro/aoc/2024/internal/parse"
)

func extractMem(name string) (string, error) {
	fp, err := os.Open(name)
	if err != nil {
		return "", err
	}
	defer fp.Close()
	s := bufio.NewScanner(fp)
	mem := ""
	for s.Scan() {
		mem += s.Text()
		mem += "\n"
	}
	return mem, nil
}

// Extracts the raw instructions from the input.
// Index 0 is the full match.
// If a multiplier instruction is detected, indices 1 & 2 will have multiplicands.
// If enabling instruction is detected, index 3 is non-empty.
// If disabling instruction is detected, index 4 is non-empty.
func rawInstructions(mem string) [][]string {
	re := regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)|(do\(\))|(don\'t\(\))`)
	return re.FindAllStringSubmatch(mem, -1)
}

func res(mem string, withToggle bool) int {
	raw := rawInstructions(mem)
	sum := 0
	on := true
	for _, inst := range raw {
		// Turn on
		if inst[3] != "" {
			on = true
		}
		// Turn off
		if inst[4] != "" {
			on = false
		}
		// Add to sum, if on
		if (on || !withToggle) && inst[1] != "" {
			sum += parse.MustInt(inst[1]) * parse.MustInt(inst[2])
		}
	}
	return sum
}

func main() {
	t := time.Now()
	mem, err := extractMem(os.Args[1])
	if err != nil {
		fmt.Printf("extractMem: %v", err)
	}
	fmt.Printf("Part 1: %d\n", res(mem, false))
	fmt.Printf("Part 2: %d\n", res(mem, true))
	fmt.Printf("Time elapsed: %s\n", time.Since(t))
}
