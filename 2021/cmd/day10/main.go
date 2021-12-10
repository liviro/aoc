package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// extractLines returns the lines in the input file.
func extractLines(fileName string) ([]string, error) {
	fp, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	var ls []string
	s := bufio.NewScanner(fp)
	s.Split(bufio.ScanLines)
	for s.Scan() {
		ls = append(ls, s.Text())
	}
	return ls, s.Err()
}

// isOpener returns whether the rune is an opening chunk.
func isOpener(r rune) bool {
	switch r {
	case '[', '(', '{', '<':
		return true
	default:
		return false
	}
}

// isOpener returns whether the rune is a closing chunk.
func isCloser(r rune) bool {
	switch r {
	case ']', ')', '}', '>':
		return true
	default:
		return false
	}
}

// matchingCloser returns the matching closer rune for the given opener.
func matchingCloser(r rune) rune {
	switch r {
	case '(':
		return ')'
	case '[':
		return ']'
	case '{':
		return '}'
	case '<':
		return '>'
	default:
		panic("invalid opener chunk rune!")
	}
}

// corruptedScore returns the score of the first corrupted rune in the given line.
// If a line is not corrupted, the score returned is 0.
func corruptedScore(l string) int {
	var chunkStack []rune
	for _, c := range l {
		if isOpener(c) {
			chunkStack = append([]rune{c}, chunkStack...)
		} else if c == matchingCloser(chunkStack[0]) {
			chunkStack = chunkStack[1:]
		} else {
			switch c {
			case ')':
				return 3
			case ']':
				return 57
			case '}':
				return 1197
			case '>':
				return 25137
			}
		}
	}
	return 0
}

// corruptedScoreSum computes the sum of the corrupted scores of all given lines.
func corruptedScoreSum(ls []string) int {
	s := 0
	for _, l := range ls {
		s += corruptedScore(l)
	}
	return s
}

// incompleteScore returns the score of completing the incomplete given line.
func incompleteScore(l string) int {
	var chunkStack []rune
	for _, c := range l {
		if isOpener(c) {
			chunkStack = append([]rune{c}, chunkStack...)
		} else if c == matchingCloser(chunkStack[0]) {
			chunkStack = chunkStack[1:]
		} else {
			panic("Got corrupted line!")
		}
	}

	s := 0
	for _, c := range chunkStack {
		s *= 5
		switch c {
		case '(':
			s += 1
		case '[':
			s += 2
		case '{':
			s += 3
		case '<':
			s += 4
		default:
			panic("invalid chunk rune!")
		}
	}
	return s
}

// completionWinner returns the winner (middle scorer) of the incomplete lines.
func completionWinner(ls []string) int {
	var scores []int
L:
	for _, l := range ls {
		if corruptedScore(l) > 0 {
			continue L
		}
		scores = append(scores, incompleteScore(l))
	}
	sort.Ints(scores)
	return scores[(len(scores)-1)/2]
}

func main() {
	ls, err := extractLines("input.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Input ingestion went wrong: ", err)
		os.Exit(1)
	}
	fmt.Println("Part 1:", corruptedScoreSum(ls))
	fmt.Println("Part 2:", completionWinner(ls))
}
