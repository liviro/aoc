package main

import "fmt"

// target represents the box that makes up the target range, with its max-min x & y values.
type target struct{ minX, maxX, minY, maxY int }

// test
// var t = target {20, 30, -10, -5}

// actual input
var t = target{70, 96, -179, -124}

// max returns the highest among the inputs.
// This should really go into some common package at this point...
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// highestY returns the highest Y position a probe that lands in the target will achieve.
func (t target) highestY() int {
	top := 0
	for x := 1; x <= t.maxX; x++ {
		// The stop condition is stupidly capped here - can probably be explained better...
		for y := 1; y < -1*t.minY*5; y++ {
			vx := x
			vy := y

			px := 0
			py := 0

			my := 0
		N:
			for n := 1; true; n++ {
				px += vx
				py += vy

				if py > my {
					my = py
				}

				vx = max(0, vx-1)
				vy = vy - 1

				if px >= t.minX && px <= t.maxX && py >= t.minY && py <= t.maxY {
					top = max(top, my)
				}

				if vx == 0 && py < t.minY {
					break N
				}
			}
		}
	}
	return top
}

// validVelocitiesCount returns number of initial velocities that will land a probe in the target area.
func (t target) validVelocitiesCount() int {
	c := 0
	for x := 1; x <= t.maxX; x++ {
		// Again, no idea what the "proper" limit is.
		for y := t.minY; y < -1*t.minY*5; y++ {
			vx := x
			vy := y

			px := 0
			py := 0

		N:
			for n := 1; true; n++ {
				px += vx
				py += vy

				vx = max(0, vx-1)
				vy = vy - 1

				if px >= t.minX && px <= t.maxX && py >= t.minY && py <= t.maxY {
					c += 1
					break N
				}

				if vx == 0 && py < t.minY {
					break N
				}
			}
		}
	}
	return c
}

func main() {
	fmt.Println("Part 1:", t.highestY())
	fmt.Println("Part 2:", t.validVelocitiesCount())
}
