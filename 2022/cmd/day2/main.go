package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	rock int = iota
	paper
	scissors
)

type round []int

func (r round) score() int {
	s := r[1] + 1 // Shape we selected: 1 for rock, 2 for paper, 3 for scissors

	// Draw
	if r[0] == r[1] {
		return s + 3
	}

	// Victory
	if r[0] == rock && r[1] == paper ||
		r[0] == paper && r[1] == scissors ||
		r[0] == scissors && r[1] == rock {
		return s + 6
	}

	// Loss
	return s
}

func sumScores(rs []round) int {
	s := 0
	for _, r := range rs {
		s += r.score()
	}
	return s
}

func strToRoundPt1(raw string) round {
	them := int(raw[0] - 'A')
	us := int(raw[2] - 'X')
	return []int{them, us}
}

func strToRoundPt2(raw string) round {
	them := int(raw[0] - 'A')
	switch raw[2] {
	// Lose
	case 'X':
		return []int{them, (them + 2) % 3}
	// Draw
	case 'Y':
		return []int{them, them}
	// Victory
	case 'Z':
		return []int{them, (them + 1) % 3}
	}

	panic("Unexpected input!")
}

func extractRounds(fileName string, mode string) ([]round, error) {
	fp, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	s := bufio.NewScanner(fp)
	var rs []round
	for s.Scan() {
		if mode == "part1" {
			rs = append(rs, strToRoundPt1(s.Text()))
		} else {
			rs = append(rs, strToRoundPt2(s.Text()))
		}
	}
	return rs, nil
}

func main() {
	in := "input.txt"
	rs1, err := extractRounds(in, "part1")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Input ingestion went wrong: ", err)
		os.Exit(1)
	}
	fmt.Println("Part 1:", sumScores(rs1))

	rs2, err := extractRounds(in, "part2")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Input ingestion went wrong: ", err)
		os.Exit(1)
	}
	fmt.Println("Part 2:", sumScores(rs2))
}
