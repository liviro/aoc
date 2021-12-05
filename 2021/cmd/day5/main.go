package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// minMax orders int64 pair passed in from smaller to larger
func minMax(a, b int64) (int64, int64) {
	if a > b {
		return b, a
	} else {
		return a, b
	}
}

type point struct{ x, y int64 }

// parsePoint returns a new point from the raw input of the format "x,y"
func parsePoint(raw string) (point, error) {
	s := strings.Split(raw, ",")

	x, err := strconv.Atoi(s[0])
	if err != nil {
		return point{}, err
	}

	y, err := strconv.Atoi(s[1])
	if err != nil {
		return point{}, err
	}

	return point{int64(x), int64(y)}, nil
}

type line struct{ from, to point }

// parseLine returns a new line from the raw input of the format "a,b -> c,d"
func parseLine(raw string) (line, error) {
	ps := strings.Split(raw, " -> ")

	from, err := parsePoint(ps[0])
	if err != nil {
		return line{}, err
	}

	to, err := parsePoint(ps[1])
	if err != nil {
		return line{}, err
	}

	return line{from, to}, nil
}

// allPoints returns a slice of all points that are on the given line
func (l line) allPoints(includeDiagonals bool) []point {
	var ps []point

	if l.from.x == l.to.x { // Horizontal
		min, max := minMax(l.from.y, l.to.y)
		for i := min; i <= max; i++ {
			ps = append(ps, point{l.from.x, i})
		}
	} else if l.from.y == l.to.y { // Vertical
		min, max := minMax(l.from.x, l.to.x)
		for i := min; i <= max; i++ {
			ps = append(ps, point{i, l.from.y})
		}
	} else if includeDiagonals { // Diagonal
		minX, maxX := minMax(l.from.x, l.to.x)

		var top point
		if minX == l.from.x {
			top = l.from
		} else {
			top = l.to
		}
		// Determine if diagonal is going left or right from top to bottom
		var yMod int64
		if minY, _ := minMax(l.from.y, l.to.y); top.y == minY {
			yMod = 1
		} else {
			yMod = -1
		}

		for i, j := top.x, top.y; i <= maxX; i, j = i+1, j+yMod {
			ps = append(ps, point{i, j})
		}
	}

	return ps
}

// extractLines parses and returns all lines in the input file
func extractLines(fileName string) ([]line, error) {
	fp, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	var ls []line
	s := bufio.NewScanner(fp)
	s.Split(bufio.ScanLines)
	for s.Scan() {
		l, err := parseLine(s.Text())
		if nil != err {
			return nil, err
		}
		ls = append(ls, l)
	}
	return ls, s.Err()
}

// countPoints counts the frequency of each point's presence in the input lines
func countPoints(lines []line, includeDiagonals bool) map[point]int64 {
	m := make(map[point]int64)
	for _, l := range lines {
		ps := l.allPoints(includeDiagonals)
		for _, p := range ps {
			m[p] += 1
		}
	}
	return m
}

// countOverlaps returns the number of points that have overlaps (more than 1 use)
func countOverlaps(uses map[point]int64) int {
	c := 0
	for _, v := range uses {
		if v >= 2 {
			c += 1
		}
	}
	return c
}

func main() {
	ls, err := extractLines("input.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Input ingestion went wrong: ", err)
		os.Exit(1)
	}
	nonDiagonalPoints := countPoints(ls, false)
	allPoints := countPoints(ls, true)
	fmt.Println("Part 1:", countOverlaps(nonDiagonalPoints))
	fmt.Println("Part 2:", countOverlaps(allPoints))
}
