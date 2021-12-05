package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// minMax orders the input int64 pair from smaller to larger.
func minMax(a, b int64) (int64, int64) {
	if a > b {
		return b, a
	} else {
		return a, b
	}
}

type point struct{ x, y int64 }

// parsePoint returns a new point from the raw input of the format "x,y".
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

// parseLine returns a new line from the raw input of the format "a,b -> c,d".
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

// gridPoints returns a slice of all points that are on a horizontal or vertical line.
// If the line is diagonal, an empty slice is returned.
func (l line) gridPoints() []point {
	var ps []point

	if l.from.x == l.to.x { // Vertical
		min, max := minMax(l.from.y, l.to.y)
		for i := min; i <= max; i++ {
			ps = append(ps, point{l.from.x, i})
		}
	} else if l.from.y == l.to.y { // Horizontal
		min, max := minMax(l.from.x, l.to.x)
		for i := min; i <= max; i++ {
			ps = append(ps, point{i, l.from.y})
		}
	}

	return ps
}

// diagonalPoints returns a slice of all points that are on a diagonal line.
// If the line is horizontal or vertical, an empty slice is returned.
func (l line) diagonalPoints() []point {
	if l.from.x == l.to.x || l.from.y == l.to.y {
		return []point{}
	}

	minX, maxX := minMax(l.from.x, l.to.x)

	var top point
	if minX == l.from.x {
		top = l.from
	} else {
		top = l.to
	}
	// Determine if diagonal is going left or right from top to bottom.
	var yDir int64
	if minY, _ := minMax(l.from.y, l.to.y); top.y == minY {
		yDir = 1
	} else {
		yDir = -1
	}

	var ps []point
	for i, j := top.x, top.y; i <= maxX; i, j = i+1, j+yDir {
		ps = append(ps, point{i, j})
	}

	return ps
}

// extractLines parses and returns all lines in the input file.
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

// countOverlaps returns the number of points that have overlaps (more than 1 use).
func countOverlaps(points []point) int {
	m := make(map[point]int64)
	c := 0
	for _, p := range points {
		m[p] += 1
		if m[p] == 2 {
			c += 1
		}
	}
	return c
}

// countGridLineOverlaps returns the number of points that have overlapping grid lines.
// This excludes diagonal lines.
func countGridLineOverlaps(lines []line) int {
	var pts []point
	for _, l := range lines {
		pts = append(pts, l.gridPoints()...)
	}
	return countOverlaps(pts)
}

// countAllOverlaps returns the number of points that have overlapping lines.
func countAllOverlaps(lines []line) int {
	var pts []point
	for _, l := range lines {
		pts = append(pts, l.gridPoints()...)
		pts = append(pts, l.diagonalPoints()...)
	}
	return countOverlaps(pts)
}

func main() {
	ls, err := extractLines("input.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Input ingestion went wrong: ", err)
		os.Exit(1)
	}
	fmt.Println("Part 1:", countGridLineOverlaps(ls))
	fmt.Println("Part 2:", countAllOverlaps(ls))
}
