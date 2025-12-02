package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/liviro/aoc/2025/internal/parse"
)

type rotation struct {
	distance int
	dir      string
}

func extractRotations(name string) ([]rotation, error) {
	fp, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	s := bufio.NewScanner(fp)
	var rs []rotation
	for s.Scan() {
		raw := s.Text()
		rs = append(rs, rotation{
			distance: parse.MustInt(raw[1:]),
			dir:      raw[:1],
		})
	}
	return rs, nil
}

func rotate(start int, r rotation) (position, zeros int) {
	position = start
	if r.dir == "L" {
		position -= r.distance
	} else {
		position += r.distance
	}
	if position < 0 {
		zeros = -1 * (position + 1) / 100
		if start != 0 {
			zeros++
		}
	} else {
		zeros = (position - 1) / 100
	}
	position %= 100
	if position < 0 {
		position += 100
	}
	if position == 0 {
		zeros++
	}
	return position, zeros
}

func passwords(rs []rotation) (int, int) {
	p := 50
	zeroPos := 0
	zeroClick := 0
	for _, r := range rs {
		pos, z := rotate(p, r)
		p = pos
		zeroClick += z
		if p == 0 {
			zeroPos++
		}
	}
	return zeroPos, zeroClick
}

func main() {
	t := time.Now()
	rs, err := extractRotations(os.Args[1])
	if err != nil {
		fmt.Printf("extractRotations: %v", err)
		return
	}
	p1, p2 := passwords(rs)
	fmt.Printf("Part 1: %d\n", p1)
	fmt.Printf("Part 2: %d\n", p2)
	fmt.Printf("Time elapsed: %s\n", time.Since(t))
}
