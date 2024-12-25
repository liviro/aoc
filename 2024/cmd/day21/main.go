package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/liviro/aoc/2024/internal/parse"
)

type coord struct{ x, y int }

var keypad = map[string]coord{
	"7": {0, 0}, "8": {1, 0}, "9": {2, 0},
	"4": {0, 1}, "5": {1, 1}, "6": {2, 1},
	"1": {0, 2}, "2": {1, 2}, "3": {2, 2},
	"G": {0, 3}, "0": {1, 3}, "A": {2, 3},
}

var dirpad = map[string]coord{
	"G": {0, 0}, "^": {1, 0}, "A": {2, 0},
	"<": {0, 1}, "v": {1, 1}, ">": {2, 1},
}

func step(pad map[string]coord, from, to string) string {
	s := pad[from]
	e := pad[to]
	var vert, horiz string
	if s.x-e.x > 0 {
		horiz = strings.Repeat("<", s.x-e.x)
	} else {
		horiz = strings.Repeat(">", e.x-s.x)
	}
	if s.y-e.y > 0 {
		vert = strings.Repeat("^", s.y-e.y)
	} else {
		vert = strings.Repeat("v", e.y-s.y)
	}
	// If moving right and gap isn't in the way: go vertical first
	if s.x-e.x < 0 && !(pad["G"].x == s.x && pad["G"].y == e.y) {
		return vert + horiz + "A"
	}
	// If gap isn't in the way: go horizontal first
	if !(pad["G"].x == e.x && pad["G"].y == s.y) {
		return horiz + vert + "A"
	}
	// Else, can definitely go vertical first.
	return vert + horiz + "A"
}

func codeToSequence(code string, numDirpads int, memo map[int]map[string]int) int {
	// Robot 1: num to dir
	var r1b strings.Builder
	csp := strings.Split("A"+code, "")
	for i := 0; i < len(csp)-1; i++ {
		r1b.WriteString(step(keypad, csp[i], csp[i+1]))
	}
	l := presses(r1b.String(), numDirpads, memo)
	return l
}

func presses(seq string, depth int, memo map[int]map[string]int) int {
	dm, ok := memo[depth]
	if !ok {
		memo[depth] = map[string]int{}
	}
	if p, ok := dm[seq]; ok {
		return p
	}

	if depth == 1 {
		l := 0
		sp := strings.Split("A"+seq, "")
		for i := 0; i < len(sp)-1; i++ {
			l += len(step(dirpad, sp[i], sp[i+1]))
		}
		memo[depth][seq] = l
		return l
	}
	sp := strings.Split("A"+seq, "")
	l := 0
	for i := 0; i < len(sp)-1; i++ {
		l += presses(step(dirpad, sp[i], sp[i+1]), depth-1, memo)
	}
	memo[depth][seq] = l
	return l
}

func complexity(codes []string, numDirpads int) int {
	s := 0
	memo := map[int]map[string]int{}
	for _, c := range codes {
		cts := codeToSequence(c, numDirpads, memo)
		s += cts * parse.MustInt(c[:3])
	}
	return s
}

func extractCodes(name string) []string {
	fp, err := os.Open(name)
	if err != nil {
		panic("Unable to open file")
	}
	defer fp.Close()
	s := bufio.NewScanner(fp)
	cs := []string{}
	for s.Scan() {
		cs = append(cs, s.Text())
	}
	return cs
}

func main() {
	t := time.Now()
	codes := extractCodes(os.Args[1])
	fmt.Printf("Part 1: %d\n", complexity(codes, 2))
	fmt.Printf("Part 2: %d\n", complexity(codes, 25))
	fmt.Printf("Time elapsed: %s\n", time.Since(t))
}
