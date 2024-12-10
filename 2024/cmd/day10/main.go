package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/liviro/aoc/2024/internal/parse"
)

var adjacents = []coord{
	{-1, 0},
	{1, 0},
	{0, -1},
	{0, 1},
}

type pos struct {
	height        int
	reachableTops map[coord]struct{}
	score         int
}

type coord struct {
	x, y int
}

func extractMap(name string) [][]*pos {
	fp, err := os.Open(name)
	if err != nil {
		panic("Unable to open file")
	}
	defer fp.Close()
	s := bufio.NewScanner(fp)
	res := [][]*pos{}
	for s.Scan() {
		r := []*pos{}
		for _, v := range strings.Split(s.Text(), "") {
			r = append(r, &pos{
				height:        parse.MustInt(v),
				reachableTops: map[coord]struct{}{},
			})
		}
		res = append(res, r)
	}
	return res
}

// Returns sum of scores and sums of ratings of trailheads
func scoreTrailheads(m [][]*pos) (int, int) {
	score := 0
	rating := 0
	for h := 9; h >= 0; h-- {
		for i, r := range m {
			for j, v := range r {
				if v.height == h {
					// Special-case tops, they're their own peaks.
					if h == 9 {
						v.reachableTops[coord{i, j}] = struct{}{}
						v.score = 1
					}
					for _, a := range adjacents {
						ni := i + a.x
						nj := j + a.y
						if inMap(m, ni, nj) && m[ni][nj].height == h+1 {
							for t := range m[ni][nj].reachableTops {
								v.reachableTops[t] = struct{}{}
							}
							v.score += m[ni][nj].score
						}
					}
					if h == 0 {
						score += len(v.reachableTops)
						rating += v.score
					}
				}
			}
		}
	}
	return score, rating
}

func inMap(m [][]*pos, i, j int) bool {
	return i >= 0 && j >= 0 && i < len(m) && j < len(m[i])
}

func main() {
	t := time.Now()
	m := extractMap(os.Args[1])
	scores, ratings := scoreTrailheads(m)
	fmt.Printf("Part 1: %d\n", scores)
	fmt.Printf("Part 2: %d\n", ratings)
	fmt.Printf("Time elapsed: %s\n", time.Since(t))
}
