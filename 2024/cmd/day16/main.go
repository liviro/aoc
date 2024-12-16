package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"slices"
	"strings"
	"time"
)

type coord struct {
	x, y int
}

type cell struct {
	dir  coord
	cost int
}

type maze map[coord]*cell

var adjacents = []coord{
	{x: 0, y: -1}, // ^
	{x: 0, y: 1},  // v
	{x: 1, y: 0},  // >
	{x: -1, y: 0}, // <
}

func extractMaze(name string) (maze, coord, coord) {
	fp, err := os.Open(name)
	if err != nil {
		panic("Unable to open file")
	}
	defer fp.Close()
	s := bufio.NewScanner(fp)
	m := map[coord]*cell{}
	i := 0
	var start, end coord
	for s.Scan() {
		for j, v := range strings.Split(s.Text(), "") {
			c := coord{x: j, y: i}
			switch {
			case v == ".": // Unexplored floor
				m[c] = &cell{cost: math.MaxInt}
			case v == "E": // End
				m[c] = &cell{cost: math.MaxInt}
				end = c
			case v == "S": // Start
				m[c] = &cell{
					cost: 0,
					dir:  coord{x: 1, y: 0}, // Start facing east
				}
				start = c
			}
		}
		i++
	}
	return m, start, end
}

func (m maze) copy() maze {
	nm := map[coord]*cell{}
	for c, v := range m {
		nm[c] = &cell{
			dir:  coord{x: v.dir.x, y: v.dir.y},
			cost: v.cost,
		}
	}
	return nm
}

func flow(m maze, start coord) {
	toProcess := []coord{start}
	for {
		if len(toProcess) == 0 {
			break
		}
		c := toProcess[0]

		for _, a := range adjacents {
			nc := coord{x: c.x + a.x, y: c.y + a.y}
			// Wall: ignore
			if _, ok := m[nc]; !ok {
				continue
			}
			var stepCost int
			// Facing same direction as before:
			// Cost goes up by 1
			if a == m[c].dir {
				stepCost = 1
			} else {
				stepCost = 1001
			}
			// If cost is less than before, go there:
			// Update direction + cost
			// And queue that cell up to be re-processed again
			if m[c].cost+stepCost < m[nc].cost {
				m[nc].cost = m[c].cost + stepCost
				m[nc].dir = a
				toProcess = append(toProcess, nc)
			}
		}
		toProcess = toProcess[1:]
	}
}

func goodSeat(orig, flowed maze, start, end, c coord) bool {
	if c == end || c == start {
		return true
	}
	maze := orig.copy()
	maze[start] = &cell{cost: math.MaxInt}
	maze[c] = &cell{
		cost: 0,
		dir: coord{
			x: flowed[c].dir.x,
			y: flowed[c].dir.y,
		},
	}
	flow(maze, c)
	return flowed[end].cost == maze[end].cost+flowed[c].cost
}

func vis(m maze, gs map[coord]struct{}) string {
	dis := [][]string{}
	for i := 0; i < 141; i++ {
		row := slices.Repeat([]string{"#"}, 141)
		dis = append(dis, row)
	}
	for c := range m {
		dis[c.y][c.x] = "."
	}
	for s := range gs {
		dis[s.y][s.x] = "O"
	}

	var b strings.Builder
	for _, r := range dis {
		for _, v := range r {
			b.WriteString(v)

		}
		b.WriteString("\n")
	}
	return b.String()
}

func part1(m maze, start, end coord) int {
	flow(m, start)
	return m[end].cost
}

func part2(orig, flowed maze, start, end coord) int {
	gs := map[coord]struct{}{}
	cs := []coord{}
	for k := range flowed {
		cs = append(cs, k)
	}
	for _, c := range cs {
		// If more expensive to get to seat then end, discard.
		if flowed[c].cost > flowed[end].cost {
			continue
		}
		if goodSeat(orig, flowed, start, end, c) {
			gs[c] = struct{}{}
		}
	}
	return len(gs)
}

func main() {
	t := time.Now()
	m, start, end := extractMaze(os.Args[1])
	flowed := m.copy()
	fmt.Printf("Part 1: %d\n", part1(flowed, start, end))
	fmt.Printf("Part 2: %d\n", part2(m, flowed, start, end))
	fmt.Printf("Time elapsed: %s\n", time.Since(t))
}
