package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

type coord struct {
	x, y int
}

var adjacents = []coord{
	{x: 0, y: -1}, // ^
	{x: 0, y: 1},  // v
	{x: 1, y: 0},  // >
	{x: -1, y: 0}, // <
}

var maxCoord = coord{x: 70, y: 70}

func shortestPath(bytes []coord, after int) int {
	corrupted := map[coord]struct{}{}
	for i := 0; i < after; i++ {
		corrupted[bytes[i]] = struct{}{}
	}

	visited := map[coord]int{
		{x: 0, y: 0}: 0,
	}
	toProcess := []coord{{x: 0, y: 0}}
	for {
		if _, ok := visited[maxCoord]; ok {
			break
		}
		if len(toProcess) == 0 {
			return -1
		}
		n := toProcess[0]
		for _, a := range adjacents {
			nc := coord{x: n.x + a.x, y: n.y + a.y}
			// Ignore out of boundaries
			if nc.x < 0 || nc.x > maxCoord.x || nc.y < 0 || nc.y > maxCoord.y {
				continue
			}
			// Ignore corrupted
			if _, ok := corrupted[nc]; ok {
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
	return visited[maxCoord]
}

func blocker(bytes []coord) coord {
	for i := 0; i < len(bytes); i++ {
		if shortestPath(bytes, i) == -1 {
			return bytes[i-1]
		}
	}
	return coord{}
}

func extractBytes(name string) []coord {
	fp, err := os.Open(name)
	if err != nil {
		panic("Unable to open file")
	}
	defer fp.Close()
	s := bufio.NewScanner(fp)
	bs := []coord{}
	for s.Scan() {
		b := coord{}
		fmt.Sscanf(s.Text(), "%d,%d", &b.x, &b.y)
		bs = append(bs, b)
	}
	return bs
}

func main() {
	t := time.Now()
	bytes := extractBytes(os.Args[1])
	fmt.Printf("Part 1: %d\n", shortestPath(bytes, 1024))
	b := blocker(bytes)
	fmt.Printf("Part 2: %d,%d\n", b.x, b.y)
	fmt.Printf("Time elapsed: %s\n", time.Since(t))
}
