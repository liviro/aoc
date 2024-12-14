package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"
)

const (
	width  = 101
	height = 103
)

type coord struct {
	x, y int
}

type robot struct {
	pos coord
	vel coord
}

func (r robot) String() string {
	return fmt.Sprintf("pos: (%d, %d), vel: (%d, %d)", r.pos.x, r.pos.y, r.vel.x, r.vel.y)
}

func (r *robot) move() {
	r.pos = coord{
		x: (r.pos.x + r.vel.x + width) % width,
		y: (r.pos.y + r.vel.y + height) % height,
	}
}

func extractRobots(name string) []*robot {
	fp, err := os.Open(name)
	if err != nil {
		panic("Unable to open file")
	}
	defer fp.Close()
	s := bufio.NewScanner(fp)
	rs := []*robot{}
	for s.Scan() {
		r := robot{}
		fmt.Sscanf(s.Text(), "p=%d,%d v=%d,%d", &r.pos.x, &r.pos.y, &r.vel.x, &r.vel.y)
		rs = append(rs, &r)
	}
	return rs
}

func safetyFactor(rs []*robot) int {
	var q1, q2, q3, q4 = 0, 0, 0, 0
	for _, r := range rs {
		switch {
		case r.pos.x < (width-1)/2 && r.pos.y < (height-1)/2:
			q1++
		case r.pos.x > (width-1)/2 && r.pos.y < (height-1)/2:
			q2++
		case r.pos.x < (width-1)/2 && r.pos.y > (height-1)/2:
			q3++
		case r.pos.x > (width-1)/2 && r.pos.y > (height-1)/2:
			q4++
		}
	}
	return q1 * q2 * q3 * q4
}

func display(rs []*robot) string {
	dis := [][]int{}
	for i := 0; i < height; i++ {
		row := slices.Repeat([]int{0}, width)
		dis = append(dis, row)
	}
	for _, r := range rs {
		dis[r.pos.y][r.pos.x]++
	}
	var b strings.Builder
	for _, r := range dis {
		for _, v := range r {
			if v == 0 {
				b.WriteString(".")
			} else {
				b.WriteString(fmt.Sprintf("%d", v))
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}

func part1(rs []*robot) int {
	for i := 0; i < 100; i++ {
		for _, r := range rs {
			r.move()
		}
	}
	return safetyFactor(rs)
}

func part2(rs []*robot) {
	// After looking through way too many printouts, noticed that a pattern
	// occasionally emerge, roughly every 101 ticks (with one such example at
	// 879), for my particular inputs. Thus, only printing those.
	maybeTreeIdx := 879
	maybeTreePeriod := 101
	for i := 101; i < 10_000; i++ {
		for _, r := range rs {
			r.move()
		}
		if (i-maybeTreeIdx)%maybeTreePeriod == 0 {
			fmt.Printf("After %d ticks:\n%s\n\n", i, display(rs))
		}
	}
}

func main() {
	t := time.Now()
	rs := extractRobots(os.Args[1])
	fmt.Printf("Part 1: %d\n", part1(rs))
	part2(rs)
	fmt.Printf("Time elapsed: %s\n", time.Since(t))
}
