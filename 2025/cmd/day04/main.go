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

type grid struct {
	rolls map[coord]struct{}
}

func (c coord) neighbors() []coord {
	var ns []coord
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			}
			ns = append(ns, coord{
				x: c.x + i,
				y: c.y + j,
			})
		}
	}
	return ns
}

func countNeighborRolls(g grid, c coord) int {
	ns := c.neighbors()
	r := 0
	for _, n := range ns {
		if _, ok := g.rolls[n]; ok {
			r++
		}
	}
	return r
}

func clear(g grid) (grid, int) {
	ng := grid{
		rolls: make(map[coord]struct{}),
	}
	c := 0
	for r := range g.rolls {
		if countNeighborRolls(g, r) >= 4 {
			ng.rolls[r] = struct{}{}
		} else {
			c++
		}
	}
	return ng, c
}

func part1(g grid) int {
	s := 0
	for r := range g.rolls {
		if countNeighborRolls(g, r) < 4 {
			s++
		}
	}
	return s
}

func part2(start grid) int {
	g := start
	s := 0
	for {
		ng, c := clear(g)
		s += c
		g = ng
		if c == 0 {
			break
		}
	}
	return s
}

func extractGrid(name string) grid {
	fp, err := os.Open(name)
	if err != nil {
		panic("Unable to open file")
	}
	defer fp.Close()
	s := bufio.NewScanner(fp)
	grid := grid{
		rolls: make(map[coord]struct{}),
	}
	y := 0
	for s.Scan() {
		raw := strings.Split(s.Text(), "")
		for i, v := range raw {
			if v == "@" {
				grid.rolls[coord{x: i, y: y}] = struct{}{}
			}
		}
		y++
	}
	return grid
}

func main() {
	t := time.Now()
	grid := extractGrid(os.Args[1])
	fmt.Printf("Part 1: %d\n", part1(grid))
	fmt.Printf("Part2: %d\n", part2(grid))
	fmt.Printf("Time elapsed: %s\n", time.Since(t))
}
