package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Move struct {
	dir string
	amt int
}

func extractMoves(fileName string) ([]Move, error) {
	fp, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	var ms []Move
	s := bufio.NewScanner(fp)
	s.Split(bufio.ScanLines)
	for s.Scan() {
		l := strings.Split(s.Text(), " ")
		amt, err := strconv.Atoi(l[1])
		if err != nil {
			return nil, err
		}
		m := Move{l[0], amt}
		ms = append(ms, m)
	}
	return ms, s.Err()
}

func finalLocalePart1(moves []Move) int {
	pos := 0
	dep := 0
	for _, m := range moves {
		switch m.dir {
		case "forward":
			pos += m.amt
		case "down":
			dep += m.amt
		case "up":
			dep -= m.amt
		default:
			fmt.Fprintln(os.Stderr, "unknown command: ", m.dir)
			os.Exit(1)
		}
	}
	return pos * dep
}

func finalLocalePart2(moves []Move) int {
	pos := 0
	dep := 0
	aim := 0
	for _, m := range moves {
		switch m.dir {
		case "forward":
			pos += m.amt
			dep += aim * m.amt
		case "down":
			aim += m.amt
		case "up":
			aim -= m.amt
		default:
			fmt.Fprintln(os.Stderr, "unknown command: ", m.dir)
			os.Exit(1)
		}
	}
	return pos * dep
}

func main() {
	ms, err := extractMoves("2-input.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Input ingestion went wrong: ", err)
		os.Exit(1)
	}
	part1 := finalLocalePart1(ms)
	part2 := finalLocalePart2(ms)
	fmt.Println("Part 1: ", part1)
	fmt.Println("Part 2: ", part2)
}
