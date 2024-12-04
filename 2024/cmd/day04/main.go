package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type coord struct {
	x, y int
}

var (
	// Part 1: XMAS in a straight line (in any direction: up / down / diagonal)
	xmasVals    []string  = []string{"X", "M", "A", "S"}
	xmasOffsets [][]coord = [][]coord{
		{coord{0, 0}, coord{0, 1}, coord{0, 2}, coord{0, 3}},
		{coord{0, 0}, coord{0, -1}, coord{0, -2}, coord{0, -3}},
		{coord{0, 0}, coord{1, 0}, coord{2, 0}, coord{3, 0}},
		{coord{0, 0}, coord{-1, 0}, coord{-2, 0}, coord{-3, 0}},
		{coord{0, 0}, coord{1, 1}, coord{2, 2}, coord{3, 3}},
		{coord{0, 0}, coord{1, -1}, coord{2, -2}, coord{3, -3}},
		{coord{0, 0}, coord{-1, 1}, coord{-2, 2}, coord{-3, 3}},
		{coord{0, 0}, coord{-1, -1}, coord{-2, -2}, coord{-3, -3}},
	}

	// Part 2: diagonal MAS crossed with diagonal MAS, either in either direction.
	crossMasVals    []string  = []string{"A", "M", "S", "M", "S"}
	crossMasOffsets [][]coord = [][]coord{
		{coord{0, 0}, coord{-1, -1}, coord{1, 1}, coord{1, -1}, coord{-1, 1}},
		{coord{0, 0}, coord{1, 1}, coord{-1, -1}, coord{1, -1}, coord{-1, 1}},
		{coord{0, 0}, coord{-1, -1}, coord{1, 1}, coord{-1, 1}, coord{1, -1}},
		{coord{0, 0}, coord{1, 1}, coord{-1, -1}, coord{-1, 1}, coord{1, -1}},
	}
)

func hasInstance(ws [][]string, start coord, offsets []coord, vals []string) bool {
	maxY := len(ws)
	maxX := len(ws[0])
	for i, o := range offsets {
		x := start.x + o.x
		y := start.y + o.y
		if x >= maxX || y >= maxY || x < 0 || y < 0 {
			return false
		}
		if ws[y][x] != vals[i] {
			return false
		}
	}
	return true
}

func countInstances(ws [][]string, offsets [][]coord, vals []string) int {
	c := 0
	for y, row := range ws {
		for x := range row {
			for _, o := range offsets {
				if hasInstance(ws, coord{x, y}, o, vals) {
					c++
				}
			}
		}
	}
	return c
}

func extractWordSearch(name string) ([][]string, error) {
	fp, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	s := bufio.NewScanner(fp)
	ws := [][]string{}
	for s.Scan() {
		ws = append(ws, strings.Split(s.Text(), ""))
	}
	return ws, nil
}

func main() {
	t := time.Now()
	ws, err := extractWordSearch(os.Args[1])
	if err != nil {
		fmt.Printf("extractWordSearch: %v", err)
	}
	fmt.Printf("Part 1: %d\n", countInstances(ws, xmasOffsets, xmasVals))
	fmt.Printf("Part 2: %d\n", countInstances(ws, crossMasOffsets, crossMasVals))
	fmt.Printf("Time elapsed: %s\n", time.Since(t))
}
