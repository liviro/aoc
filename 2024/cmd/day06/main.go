package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"
)

type coord struct {
	x, y int
}

type guard struct {
	position  coord
	direction coord
}

func (g *guard) rotate() {
	rotations := []coord{
		{0, -1}, // ^
		{1, 0},  // >
		{0, 1},  // v
		{-1, 0}, // >
	}
	i := slices.Index(rotations, g.direction)
	g.direction = rotations[(i+1)%4]
}

func (g guard) nextPosition() coord {
	return coord{
		x: g.position.x + g.direction.x,
		y: g.position.y + g.direction.y,
	}
}

func (g *guard) move() {
	g.position = g.nextPosition()
}

func extractMapData(name string) ([][]bool, *guard) {
	fp, err := os.Open(name)
	if err != nil {
		panic("Unable to open file")
	}
	defer fp.Close()
	s := bufio.NewScanner(fp)
	obst := [][]bool{}
	guard := &guard{
		direction: coord{0, -1},
	}
	for s.Scan() {
		r := []bool{}
		raw := strings.Split(s.Text(), "")
		for i, v := range raw {
			switch {
			case v == "#":
				r = append(r, true)
			case v == ".":
				r = append(r, false)
			case v == "^":
				guard.position = coord{x: i, y: len(obst)}
				r = append(r, false)
			}
		}
		obst = append(obst, r)
	}
	return obst, guard
}

func countPositions(obst [][]bool, g guard) int {
	visited := map[coord]struct{}{}
	// Include starting position.
	visited[coord{x: g.position.x, y: g.position.y}] = struct{}{}
	for {
		ahead := g.nextPosition()
		if ahead.x < 0 || ahead.y < 0 || ahead.x >= len(obst[0]) || ahead.y >= len(obst) {
			break
		}
		if obst[ahead.y][ahead.x] {
			g.rotate()
		} else {
			g.move()
			visited[coord{x: g.position.x, y: g.position.y}] = struct{}{}
		}
	}
	return len(visited)
}

func hasLoop(obst [][]bool, g guard) bool {
	visited := map[guard]struct{}{}
	visited[g] = struct{}{}
	for {
		ahead := g.nextPosition()
		if ahead.x < 0 || ahead.y < 0 || ahead.x >= len(obst[0]) || ahead.y >= len(obst) {
			break
		}
		if obst[ahead.y][ahead.x] {
			g.rotate()
		} else {
			g.move()
			if _, ok := visited[g]; ok {
				return true
			}
			visited[g] = struct{}{}
		}
	}
	return false
}

func countLoops(obst [][]bool, g guard) int {
	c := 0
	for i := range obst {
		for j := range obst[i] {
			if !obst[i][j] && !(g.position.x == j && g.position.y == i) {
				// Ick to mutating, but am too lazy to properly copy :s
				obst[i][j] = true
				if hasLoop(obst, g) {
					c++
				}
				obst[i][j] = false
			}
		}
	}
	return c
}

func main() {
	t := time.Now()
	obst, guard := extractMapData(os.Args[1])
	fmt.Printf("Part 1: %d\n", countPositions(obst, *guard))
	fmt.Printf("Part2: %d\n", countLoops(obst, *guard))
	fmt.Printf("Time elapsed: %s\n", time.Since(t))
}
