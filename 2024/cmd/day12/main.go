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

type plot struct {
	plant   string
	visited bool
}

var adjacents = []coord{
	// Above
	{0, -1},
	// To right
	{1, 0},
	// Below
	{0, 1},
	// To left
	{-1, 0},
}

func adjacentCoords(c coord, max coord) []coord {
	cs := []coord{}
	for _, a := range adjacents {
		n := coord{x: c.x + a.x, y: c.y + a.y}
		if n.x >= 0 && n.y >= 0 && n.x <= max.x && n.y <= max.y {
			cs = append(cs, n)
		}
	}
	return cs
}

func extractMap(name string) [][]*plot {
	fp, err := os.Open(name)
	if err != nil {
		panic("Unable to open file")
	}
	defer fp.Close()
	s := bufio.NewScanner(fp)
	res := [][]*plot{}
	for s.Scan() {
		r := []*plot{}
		for _, v := range strings.Split(s.Text(), "") {
			r = append(r, &plot{
				plant:   v,
				visited: false,
			})
		}
		res = append(res, r)
	}
	return res
}

// Mutates plot to note what has been visited, returns the
// (perimeter, area, sides) of the region that contains the starting
// index.
func processRegion(m [][]*plot, start coord) (int, int, int) {
	max := coord{x: len(m[0]) - 1, y: len(m) - 1}
	tbc := map[coord]struct{}{}
	done := map[coord]struct{}{}
	tbc[start] = struct{}{}
	area := 0
	perimeter := 0
	for {
		if len(tbc) == 0 {
			break
		}
		for p := range tbc {
			area++
			regionAdjs := 0
			adjs := adjacentCoords(p, max)
			for _, a := range adjs {
				if m[a.y][a.x].plant == m[p.y][p.x].plant {
					regionAdjs++
					if _, ok := done[a]; !ok {
						tbc[a] = struct{}{}
					}
				}
			}
			perimeter += (4 - regionAdjs)
			m[p.y][p.x].visited = true
			done[p] = struct{}{}
			delete(tbc, p)
		}

	}
	sides := countSides(done, max)
	return perimeter, area, sides
}

func countSides(region map[coord]struct{}, max coord) int {
	s := 0
	// Count horizontal sides
	for i := 0; i <= max.y+1; i++ {
		isOn := false
		isAbove := false
		for j := 0; j <= max.x; j++ {
			nextAbove := i != 0 && isIn(region, coord{x: j, y: i - 1})
			nextBelow := i != max.y+1 && isIn(region, coord{x: j, y: i})
			nextOn := nextAbove != nextBelow
			// New side is starting, either:
			// - Was previously off, is now on
			// - Region side flipped from above to below (or vice versa)
			if (nextOn && !isOn) || (nextOn && isOn && isAbove != nextAbove) {
				s++
			}
			isOn = nextOn
			isAbove = nextAbove
		}
	}
	// Count vertical sides
	for j := 0; j <= max.x+1; j++ {
		isOn := false
		isLeft := false
		for i := 0; i <= max.y; i++ {
			nextLeft := j != 0 && isIn(region, coord{x: j - 1, y: i})
			nextRight := j != max.x+1 && isIn(region, coord{x: j, y: i})
			nextOn := nextLeft != nextRight
			// New side is starting, either:
			// - Was previously off, is now on
			// - Region side flipped from left to right (or vice versa)
			if (nextOn && !isOn) || (nextOn && isOn && isLeft != nextLeft) {
				s++
			}
			isOn = nextOn
			isLeft = nextLeft
		}
	}
	return s
}

func isIn(region map[coord]struct{}, c coord) bool {
	_, ok := region[c]
	return ok
}

func totalPrices(m [][]*plot) (int, int) {
	part1 := 0
	part2 := 0
	for i := range m {
		for j := range m[i] {
			if !m[i][j].visited {
				p, a, s := processRegion(m, coord{x: j, y: i})
				part1 += p * a
				part2 += s * a
			}
		}
	}
	return part1, part2
}

func main() {
	t := time.Now()
	m := extractMap(os.Args[1])
	p1, p2 := totalPrices(m)
	fmt.Printf("Part 1: %d\n", p1)
	fmt.Printf("Part 2: %d\n", p2)
	fmt.Printf("Time elapsed: %s\n", time.Since(t))
}
