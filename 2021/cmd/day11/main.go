package main

import (
	"fmt"
	"os"
	"strings"
)

// min returns the smaller of the two given ints.
func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

// max returns the larger of the two given ints.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// size is the side length of the octopus grid.
const size = 10

// octopus is the cute, flashy critter with an energy level.
type octopus struct {
	energy   int
	flashing bool
}

// parseOctopus returns an octopus of the energy level given in the raw input byte
func parseOctopus(r byte) octopus {
	return octopus{energy: int(r - '0')}
}

// bump increases the octopus' energy level, and, if applicable, makes it flash.
func (octo *octopus) bump() {
	if octo.flashing {
		return
	}
	octo.energy += 1
	if octo.energy > 9 {
		octo.flashing = true
	}
}

// location identifies a location within a grid.
type location struct{ i, j int }

// adjacents returns all adjacent locations to the given location.
func (l location) adjacents() []location {
	var ls []location
	for i := max(0, l.i-1); i <= min(size-1, l.i+1); i++ {
		for j := max(0, l.j-1); j <= min(size-1, l.j+1); j++ {
			if !(i == l.i && j == l.j) {
				ls = append(ls, location{i, j})
			}
		}
	}
	return ls
}

// grid is the 10x10 grid of octopuses.
type grid [size][size]octopus

// at returns a pointer to the octopus at the given location in the grid.
func (g *grid) at(l location) *octopus {
	return &g[l.i][l.j]
}

// reset puts a fresh octopus (not flashing, energy level at 0) at the given location in the grid.
func (g *grid) reset(l location) {
	g[l.i][l.j] = octopus{}
}

// step advances the grid by one step and returns the number of flashes that happpened during it.
func (g *grid) step() int {
	// Each location may have at most 8 flashing neighbors.
	flashNeighbors := make(chan location, size*size*8)
	flashes := make(map[location]bool)
	// Bump everyone by 1.
	for i := range g {
		for j := range g[i] {
			l := location{i, j}
			bumpInStep(g, l, flashNeighbors, flashes)
		}
	}
L:
	for {
		// Bump all adjacents of flashing octopuses.
		select {
		case l := <-flashNeighbors:
			bumpInStep(g, l, flashNeighbors, flashes)
		default:
			break L
		}
	}

	// Reset the energy and flash status of the flashing octopuses, and count them.
	c := 0
	for f := range flashes {
		c += 1
		g.reset(f)
	}
	return c
}

// bumpInStep bumps the octopus at the given location, and if it then flashes, queues its adjacents to be bumped.
// This is a helper to the step method.
func bumpInStep(g *grid, l location, flashNeighbors chan location, flashes map[location]bool) {
	octo := g.at(l)
	octo.bump()
	if _, ok := flashes[l]; octo.flashing && !ok {
		flashes[l] = true
		for _, a := range l.adjacents() {
			// Already flashing adjacents won't be affected by this extra bump.
			if _, ok := flashes[a]; !ok {
				flashNeighbors <- a
			}
		}
	}
}

// flashesAfter returns the number of total flashes the grid will see after the given number of steps.
// Note that this method words on a copy of a grid and does not mutate it.
func (g grid) flashesAfter(steps int) int {
	f := 0
	for i := 0; i < steps; i++ {
		f += g.step()
	}
	return f
}

// firstSyncStep returns the number of steps it takes for all octopuses to flash in sync for the first time.
// Note that this method words on a copy of a grid and does not mutate it.
func (g grid) firstSyncStep() int {
	for i := 1; ; i++ {
		f := g.step()
		if f == size*size {
			return i
		}
	}
}

// extractGrid extracts the grid in the given file.
func extractGrid(fileName string) (grid, error) {
	raw, err := os.ReadFile(fileName)
	if err != nil {
		return grid{}, err
	}
	rows := strings.Split(string(raw), "\n")
	var g grid
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			g[i][j] = parseOctopus(rows[i][j])
		}
	}
	return g, nil
}

func main() {
	g, err := extractGrid("input.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Input ingestion went wrong: ", err)
		os.Exit(1)
	}

	fmt.Println("Part 1:", g.flashesAfter(100))
	fmt.Println("Part 2:", g.firstSyncStep())
}
