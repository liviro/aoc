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

// Extract the coordinates per frequency and size of map
// (latter expressed as extreme-most coordinate).
func extractMapData(name string) (map[string][]coord, coord) {
	fp, err := os.Open(name)
	if err != nil {
		panic("Unable to open file")
	}
	defer fp.Close()
	s := bufio.NewScanner(fp)
	freqs := map[string][]coord{}
	i := 0
	max := coord{}
	for s.Scan() {
		for j, v := range strings.Split(s.Text(), "") {
			if v != "." {
				if _, ok := freqs[v]; !ok {
					freqs[v] = []coord{}
				}
				freqs[v] = append(freqs[v], coord{j, i})
			}
			max = coord{j, i}
		}
		i++
	}
	return freqs, max
}

func countAntinodes(freqs map[string][]coord, max coord) int {
	ns := map[coord]struct{}{}
	for _, cs := range freqs {
		for i, a := range cs {
			for j := i + 1; j < len(cs); j++ {
				b := cs[j]
				n1 := coord{
					x: 2*a.x - b.x,
					y: 2*a.y - b.y,
				}
				if inRange(n1, max) {
					ns[n1] = struct{}{}
				}
				n2 := coord{
					x: 2*b.x - a.x,
					y: 2*b.y - a.y,
				}
				if inRange(n2, max) {
					ns[n2] = struct{}{}
				}
			}
		}
	}
	return len(ns)
}
func countAntinodesWithHarmonics(freqs map[string][]coord, max coord) int {
	ns := map[coord]struct{}{}
	for _, cs := range freqs {
		for i, a := range cs {
			for j := i + 1; j < len(cs); j++ {
				b := cs[j]
				delta := coord{
					x: b.x - a.x,
					y: b.y - a.y,
				}
				// Go back (includes node a)
				k := 0
				for {
					nn := coord{
						x: a.x - k*delta.x,
						y: a.y - k*delta.y,
					}
					if !inRange(nn, max) {
						break
					}
					ns[nn] = struct{}{}
					k++
				}
				// Go forward (includes node b)
				k = 1
				for {
					nn := coord{
						x: a.x + k*delta.x,
						y: a.y + k*delta.y,
					}
					if !inRange(nn, max) {
						break
					}
					ns[nn] = struct{}{}
					k++
				}
			}
		}
	}
	return len(ns)
}

func inRange(c coord, max coord) bool {
	return c.x >= 0 && c.y >= 0 && c.x <= max.x && c.y <= max.y
}

func main() {
	t := time.Now()
	freqs, max := extractMapData(os.Args[1])
	fmt.Printf("Part 1: %d\n", countAntinodes(freqs, max))
	fmt.Printf("Part 2: %d\n", countAntinodesWithHarmonics(freqs, max))
	fmt.Printf("Time elapsed: %s\n", time.Since(t))
}
