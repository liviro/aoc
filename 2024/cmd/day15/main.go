package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type coord struct {
	x, y int
}

func (c coord) Print() string {
	return fmt.Sprintf("(%d, %d)", c.x, c.y)
}

// 0: nothing
// 1: wall
// 2: box (or left side of wide box)
// 3: right side of wide box
type warehouse [][]int

func printMap(wh warehouse, robot coord, isWide bool) string {
	var b strings.Builder
	for i, r := range wh {
		for j, v := range r {
			switch {
			case i == robot.y && j == robot.x:
				b.WriteString("@")
			case v == 0:
				b.WriteString(".")
			case v == 1:
				b.WriteString("#")
			case v == 2 && isWide:
				b.WriteString("[")
			case v == 2:
				b.WriteString("O")
			case v == 3:
				b.WriteString("]")
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}

func extractWarehouse(raw string) (warehouse, *coord) {
	wh := [][]int{}
	robot := &coord{}
	for i, r := range strings.Split(raw, "\n") {
		row := []int{}
		for j, v := range strings.Split(r, "") {
			switch {
			case v == ".":
				row = append(row, 0)
			case v == "#":
				row = append(row, 1)
			case v == "O":
				row = append(row, 2)
			case v == "@":
				row = append(row, 0)
				robot.x = j
				robot.y = i
			}
		}
		wh = append(wh, row)
	}
	return wh, robot
}

func extractWideWarehouse(raw string) (warehouse, *coord) {
	wh := [][]int{}
	robot := &coord{}
	for i, r := range strings.Split(raw, "\n") {
		row := []int{}
		for j, v := range strings.Split(r, "") {
			switch {
			case v == ".":
				row = append(row, 0)
				row = append(row, 0)
			case v == "#":
				row = append(row, 1)
				row = append(row, 1)
			case v == "O":
				row = append(row, 2)
				row = append(row, 3)
			case v == "@":
				row = append(row, 0)
				row = append(row, 0)
				robot.x = 2 * j
				robot.y = i
			}
		}
		wh = append(wh, row)
	}
	return wh, robot
}

func extractMoves(raw string) []coord {
	ms := []coord{}
	for _, m := range strings.Split(raw, "") {
		switch {
		case m == "^":
			ms = append(ms, coord{x: 0, y: -1})
		case m == ">":
			ms = append(ms, coord{x: 1, y: 0})
		case m == "v":
			ms = append(ms, coord{x: 0, y: 1})
		case m == "<":
			ms = append(ms, coord{x: -1, y: 0})
		}
	}
	return ms
}

func extractInput(name string) (warehouse, *coord, warehouse, *coord, []coord) {
	raw, err := os.ReadFile(name)
	if err != nil {
		panic("Cannot open file!")
	}
	sections := strings.Split(strings.TrimSpace(string(raw)), "\n\n")
	wh, robot := extractWarehouse(sections[0])
	wwh, wRobot := extractWideWarehouse(sections[0])
	moves := extractMoves(sections[1])
	return wh, robot, wwh, wRobot, moves
}

func attemptMove(wh warehouse, robot *coord, move coord) {
	stack := []int{}
	next := coord{
		x: robot.x + move.x,
		y: robot.y + move.y,
	}
S:
	for {
		// Bumped wall: abort and do nothing
		if wh[next.y][next.x] == 1 {
			return
		}
		// Empty space: stop stacking
		if wh[next.y][next.x] == 0 {
			break S
		}
		// Box: add to stack
		if wh[next.y][next.x] == 2 || wh[next.y][next.x] == 3 {
			stack = append(stack, wh[next.y][next.x])
		}
		next = coord{
			x: next.x + move.x,
			y: next.y + move.y,
		}
	}
	for s := len(stack) - 1; s >= 0; s-- {
		wh[next.y][next.x] = stack[s]
		next = coord{
			x: next.x - move.x,
			y: next.y - move.y,
		}
	}
	wh[robot.y][robot.x] = 0
	robot.x = robot.x + move.x
	robot.y = robot.y + move.y
	wh[robot.y][robot.x] = 0
}

func attemptWideMove(wh warehouse, robot *coord, move coord) {
	// Can use old move attempter
	if move.y == 0 {
		attemptMove(wh, robot, move)
		return
	}
	init := map[coord]struct{}{}
	init[*robot] = struct{}{}
	stack := []map[coord]struct{}{init}
S:
	for {
		toCheck := stack[len(stack)-1]
		nextStack := map[coord]struct{}{}
		for s := range toCheck {
			// Bumped wall: abort and do nothing
			if wh[s.y+move.y][s.x] == 1 {
				return
			}
			if wh[s.y+move.y][s.x] == 2 {
				nextStack[coord{x: s.x, y: s.y + move.y}] = struct{}{}
				nextStack[coord{x: s.x + 1, y: s.y + move.y}] = struct{}{}
			}
			if wh[s.y+move.y][s.x] == 3 {
				nextStack[coord{x: s.x - 1, y: s.y + move.y}] = struct{}{}
				nextStack[coord{x: s.x, y: s.y + move.y}] = struct{}{}
			}
		}
		// All empty space detected: stop stacking
		if len(nextStack) == 0 {
			break S
		}
		stack = append(stack, nextStack)
	}
	// For each stack, move that row up.
	for sri := len(stack) - 1; sri > 0; sri-- {
		for se := range stack[sri] {
			wh[se.y+move.y][se.x] = wh[se.y][se.x]
			wh[se.y][se.x] = 0
		}
	}
	wh[robot.y][robot.x] = 0
	robot.y = robot.y + move.y
	wh[robot.y][robot.x] = 0
}

func gps(wh warehouse) int {
	s := 0
	for i, r := range wh {
		for j, v := range r {
			if v == 2 {
				s += 100*i + j
			}
		}
	}
	return s
}

func main() {
	t := time.Now()
	wh, robot, wwh, wrobot, moves := extractInput(os.Args[1])
	// Moves for part 1
	for _, m := range moves {
		attemptMove(wh, robot, m)
	}
	fmt.Printf("Part 1: %d\n", gps(wh))
	// Moves for part 2
	for _, m := range moves {
		attemptWideMove(wwh, wrobot, m)
	}
	fmt.Printf("Part 2: %d\n", gps(wwh))
	fmt.Printf("Time elapsed: %s\n", time.Since(t))
}
