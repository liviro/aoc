package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/liviro/aoc/2025/internal/parse"
)

type position struct {
	x, y, z int
}

func (p position) distance(q position) float64 {
	return math.Sqrt(float64(
		(p.x-q.x)*(p.x-q.x) + (p.y-q.y)*(p.y-q.y) + (p.z-q.z)*(p.z-q.z),
	))
}

func (p position) String() string {
	return fmt.Sprintf("{%d, %d, %d}", p.x, p.y, p.z)
}

type conn struct {
	a, b position
	dist float64
}

func (c conn) String() string {
	return fmt.Sprintf("%s - %s (dist = %.2f)", c.a, c.b, c.dist)
}

func extractBoxes(name string) []position {
	fp, err := os.Open(name)
	if err != nil {
		panic("Unable to open file")
	}
	defer fp.Close()
	s := bufio.NewScanner(fp)
	var ps []position
	for s.Scan() {
		raw := strings.Split(s.Text(), ",")
		ps = append(ps, position{
			x: parse.MustInt(raw[0]),
			y: parse.MustInt(raw[1]),
			z: parse.MustInt(raw[2]),
		})
	}
	return ps
}

func allDistances(bs []position) []conn {
	var cs []conn
	for a := 0; a < len(bs); a++ {
		for b := a + 1; b < len(bs); b++ {
			cs = append(cs, conn{
				a: bs[a], b: bs[b],
				dist: bs[a].distance(bs[b]),
			})
		}
	}
	return cs
}

func addToCircuit(circuits []map[position]struct{}, c conn) []map[position]struct{} {
	idxA := -1
	idxB := -1
	for i, m := range circuits {
		if _, ok := m[c.a]; ok {
			idxA = i
		}
		if _, ok := m[c.b]; ok {
			idxB = i
		}
	}
	if idxA != -1 {
		circuits[idxA][c.b] = struct{}{}
	}
	if idxB != -1 {
		circuits[idxB][c.a] = struct{}{}
	}

	if idxA == -1 && idxB == -1 {
		nm := make(map[position]struct{})
		nm[c.a] = struct{}{}
		nm[c.b] = struct{}{}
		circuits = append(circuits, nm)
	}

	if idxA != -1 && idxB != -1 && idxA != idxB {
		for b := range circuits[idxB] {
			circuits[idxA][b] = struct{}{}
		}
		circuits = slices.Delete(circuits, idxB, idxB+1)
	}
	return circuits
}

func part1(bs []position) int {
	all := allDistances(bs)
	sort.Slice(all, func(i, j int) bool { return all[i].dist < all[j].dist })
	top := all[:1000]
	var cs []map[position]struct{}
	for _, c := range top {
		cs = addToCircuit(cs, c)
	}

	var sizes []int
	for _, m := range cs {
		sizes = append(sizes, len(m))
	}
	sort.Ints(sizes)

	res := 1
	for i := len(sizes) - 1; i >= len(sizes)-3; i-- {
		res *= sizes[i]
	}
	return res
}

func part2(bs []position) int {
	all := allDistances(bs)
	sort.Slice(all, func(i, j int) bool { return all[i].dist < all[j].dist })
	var cs []map[position]struct{}
	for _, c := range all {
		cs = addToCircuit(cs, c)
		if len(cs[0]) == len(bs) {
			return c.a.x * c.b.x
		}
	}
	return 0
}

func main() {
	t := time.Now()
	bs := extractBoxes(os.Args[1])
	fmt.Printf("Part 1: %d\n", part1(bs))
	fmt.Printf("Part 2: %d\n", part2(bs))
	fmt.Printf("Time elapsed: %s\n", time.Since(t))
}
