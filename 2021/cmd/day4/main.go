package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func extractInput(fileName string) ([]int64, []board, error) {
	raw, err := os.ReadFile(fileName)
	if err != nil {
		return nil, nil, err
	}

	blocks := strings.Split(strings.TrimSpace(string(raw)), "\n\n")

	ds, err := parseDraws(blocks[0])
	if err != nil {
		return nil, nil, err
	}

	bs, err := parseBoards(blocks[1:])
	if err != nil {
		return nil, nil, err
	}

	return ds, bs, nil
}

func parseDraws(raw string) ([]int64, error) {
	split := strings.Split(raw, ",")
	draws := make([]int64, len(split))
	for i, s := range split {
		d, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		draws[i] = int64(d)
	}
	return draws, nil
}

func parseBoards(raw []string) ([]board, error) {
	bs := make([]board, len(raw))
	for i, s := range raw {
		board, err := parseBoard(s)
		if err != nil {
			return nil, err
		}
		bs[i] = board
	}
	return bs, nil
}

func firstWinningScore(draws []int64, boards []board) int64 {
	bs := make([]board, len(boards))
	copy(bs, boards)

	for _, d := range draws {
		for i := range bs {
			bs[i].mark(d)
			if bs[i].hasWin() {
				return bs[i].score(d)
			}
		}
	}
	return 0
}

func lastWinningScore(draws []int64, boards []board) int64 {
	bs := make([]board, len(boards))
	copy(bs, boards)

	for _, d := range draws {
		nextBs := []board{}
		for _, b := range bs {
			b.mark(d)
			if len(bs) == 1 && b.hasWin() {
				return b.score(d)
			} else {
				if !b.hasWin() {
					nextBs = append(nextBs, b)
				}
			}
		}
		bs = nextBs
	}
	return 0
}

func main() {
	ds, bs, err := extractInput("input.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Input ingestion went wrong: ", err)
		os.Exit(1)
	}

	fmt.Println("Part 1:", firstWinningScore(ds, bs))
	fmt.Println("Part 2:", lastWinningScore(ds, bs))
}
