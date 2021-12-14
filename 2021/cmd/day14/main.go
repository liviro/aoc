package main

import (
	"fmt"
	"math"
	"os"
	"strings"
)

// polymer represents a polymer, with its rules and a count of the existing pairs, as well as the last element.
// The last element is used in counting the number of times an element appears.
type polymer struct {
	rules map[string][]string
	last  rune
	pairs map[string]int
}

// extractRules parses a raw string representation of rules into a map of pairs to their post-step pairs.
func extractRules(raw string) map[string][]string {
	rs := make(map[string][]string)
	for _, r := range strings.Split(raw, "\n") {
		pts := strings.Split(r, " -> ")
		next := []string{string(pts[0][0]) + pts[1], pts[1] + string(pts[0][1])}
		rs[pts[0]] = next
	}
	return rs
}

// extractPairs extracts the polymer pairs from a string representation of a polymer template.
func extractPairs(template string) map[string]int {
	ps := make(map[string]int)
	for i := 0; i < len(template)-1; i++ {
		ps[template[i:i+2]] += 1
	}
	return ps
}

// extractPolymer extracts a polymer from the given file.
func extractPolymer(fileName string) (polymer, error) {
	raw, err := os.ReadFile(fileName)
	if err != nil {
		return polymer{}, err
	}

	blocks := strings.Split(strings.TrimSpace(string(raw)), "\n\n")

	return polymer{
		rules: extractRules(blocks[1]),
		last:  rune(blocks[0][len(blocks[0])-1]),
		pairs: extractPairs(blocks[0]),
	}, nil
}

// step evolves the given polymer one step.
func (pm *polymer) step() {
	np := make(map[string]int)
	for p, c := range pm.pairs {
		rs := pm.rules[p]
		np[rs[0]] += c
		np[rs[1]] += c
	}
	pm.pairs = np
}

// frequencyDelta returns the difference between the frequencies of the most and least common elements of the polymer.
func (pm *polymer) frequencyDelta() int {
	elemCts := make(map[rune]int)
	for p, c := range pm.pairs {
		elemCts[rune(p[0])] += c
	}
	// The above loop only increments the counts by the first element in the known pairs. This leaves
	// out the last element, which is accounted for here.
	elemCts[pm.last] += 1
	min := math.MaxInt64
	max := 0
	for _, c := range elemCts {
		if c > max {
			max = c
		}
		if c < min {
			min = c
		}
	}
	return max - min
}

// frequencyDeltaAfterSteps returns the frequency delta after the polymer has been evolved the given number of // steps.
func (pm polymer) frequencyDeltaAfterSteps(steps int) int {
	for i := 0; i < steps; i++ {
		pm.step()
	}
	return pm.frequencyDelta()
}

func main() {
	p, err := extractPolymer("input.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Input ingestion went wrong: ", err)
		os.Exit(1)
	}

	fmt.Println("Part 1:", p.frequencyDeltaAfterSteps(10))
	fmt.Println("Part 2:", p.frequencyDeltaAfterSteps(40))
}
