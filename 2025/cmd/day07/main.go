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

func (c coord) String() string {
	return fmt.Sprintf("{%d, %d}", c.x, c.y)
}

type diagram struct {
	start     coord
	splitters map[coord]struct{}
	max       coord
}

func extractDiagram(name string) diagram {
	fp, err := os.Open(name)
	if err != nil {
		panic("Unable to open file")
	}
	defer fp.Close()
	s := bufio.NewScanner(fp)
	d := diagram{
		splitters: make(map[coord]struct{}),
	}
	y := 0
	maxX := 0
	for s.Scan() {
		raw := strings.Split(s.Text(), "")
		maxX = len(raw)
		for i, v := range raw {
			if v == "S" {
				d.start = coord{x: i, y: y}
			}
			if v == "^" {
				d.splitters[coord{x: i, y: y}] = struct{}{}
			}
		}
		y++
	}
	d.max = coord{x: maxX, y: y - 1}
	return d
}

func (d diagram) analyze() (splits, timelines int) {
	beams := make(map[int]int)
	beams[d.start.x] = 1
	for j := 2; j < d.max.y; j += 2 {
		nb := beams
		for i := 0; i <= d.max.x; i++ {
			_, hasSplitter := d.splitters[coord{x: i, y: j}]
			_, hasBeam := beams[i]
			if hasSplitter && hasBeam {
				splits++
				nb[i-1] += beams[i]
				nb[i+1] += beams[i]
				delete(nb, i)
			}
		}
		beams = nb
	}
	for _, t := range beams {
		timelines += t
	}
	return splits, timelines
}

func main() {
	t := time.Now()
	d := extractDiagram(os.Args[1])
	splits, timelines := d.analyze()
	fmt.Printf("Part 1: %d\n", splits)
	fmt.Printf("Part 2: %d\n", timelines)
	fmt.Printf("Time elapsed: %s\n", time.Since(t))
}
