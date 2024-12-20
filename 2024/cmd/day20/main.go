package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
	"time"
)

type coord struct {
	x, y int
}

func (ca coord) Equal(cb coord) bool {
	return ca.x == cb.x && ca.y == cb.y
}

func (ca coord) dist(cb coord) int {
	return int(math.Abs(float64(ca.x-cb.x))) + int(math.Abs(float64(ca.y-cb.y)))
}

var adjacents = []coord{
	{x: 0, y: -1}, // ^
	{x: 0, y: 1},  // v
	{x: 1, y: 0},  // >
	{x: -1, y: 0}, // <
}

func extractMaze(name string) ([]coord, coord, coord) {
	fp, err := os.Open(name)
	if err != nil {
		panic("Unable to open file")
	}
	defer fp.Close()
	s := bufio.NewScanner(fp)
	walls := []coord{}
	var start, end coord
	i := 0
	for s.Scan() {
		for j, v := range strings.Split(s.Text(), "") {
			c := coord{x: i, y: j}
			switch {
			case v == "#": // Wall
				walls = append(walls, c)
			case v == "E": // End
				end = c
			case v == "S": // Start
				start = c
			}
		}
		i++
	}
	return walls, start, end
}

func shortestPath(walls []coord, start, end coord) map[coord]int {
	wm := map[coord]struct{}{}
	for _, w := range walls {
		wm[w] = struct{}{}
	}

	visited := map[coord]int{
		start: 0,
	}
	toProcess := []coord{start}
	for {
		if _, ok := visited[end]; ok {
			break
		}
		if len(toProcess) == 0 {
			return map[coord]int{}
		}
		n := toProcess[0]
		for _, a := range adjacents {
			nc := coord{x: n.x + a.x, y: n.y + a.y}
			// Ignore walls
			if _, ok := wm[nc]; ok {
				continue
			}
			// Ignore already seen
			if _, ok := visited[nc]; ok {
				continue
			}

			visited[nc] = visited[n] + 1
			toProcess = append(toProcess, nc)
		}
		toProcess = toProcess[1:]
	}

	bestPath := map[coord]int{
		end: visited[end],
	}
	c := end
	for {
		if _, ok := bestPath[start]; ok {
			break
		}
	A:
		for _, a := range adjacents {
			nc := coord{x: c.x + a.x, y: c.y + a.y}
			if v, ok := visited[nc]; ok && v+1 == visited[c] {
				bestPath[nc] = visited[nc]
				c = nc
				break A
			}
		}
	}

	return bestPath
}

func savingCheats(walls []coord, start, end coord, cheatSize int) int {
	cheatsOver := 0
	sp := shortestPath(walls, start, end)

	for c1, v1 := range sp {
		for c2, v2 := range sp {
			if c1.dist(c2) <= cheatSize {
				newDist := v1 + c1.dist(c2) + (sp[end] - v2)
				if sp[end]-newDist >= 100 {
					cheatsOver++
				}
			}
		}
	}

	return cheatsOver
}

func main() {
	t := time.Now()
	walls, start, end := extractMaze(os.Args[1])
	fmt.Printf("Part1 : %d\n", savingCheats(walls, start, end, 2))
	fmt.Printf("Part2 : %d\n", savingCheats(walls, start, end, 20))
	fmt.Printf("Time elapsed: %s\n", time.Since(t))
}
