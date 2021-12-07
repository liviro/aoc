package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func abs(a int) int {
  if a < 0 {
    return -1*a
  } else {
    return a
  }
}

// crabArmy represents the crab swarm.
// It contains their most extreme positions and a map of their positions to their numbers there.
type crabArmy struct {
	min, max int
	pos      map[int]int
}

// parseCrabArmy extracts a crab army from the list of string representations of their positions.
func parseCrabArmy(raw []string) (crabArmy, error) {
	min := math.MaxInt64
	max := 0
	pos := make(map[int]int)
	for _, s := range raw {
		c, err := strconv.Atoi(s)
		if err != nil {
			return crabArmy{}, err
		}
		pos[c] += 1
		if min > c {
			min = c
		}
		if max < c {
			max = c
		}
	}
	return crabArmy{min, max, pos}, nil
}

// linearFuelCost computes the cost of moving the crab army to the given place if each step costs 1 fuel.
func (ca crabArmy) linearFuelCost(place int) int {
  c := 0
  for p, n := range ca.pos {
    c += abs(place - p) * n
  }
  return c
}

// triangularFuelCost computes the cost of moving the crab army to the given place if each additional step costs 1 more fuel than the previous.
func (ca crabArmy) triangularFuelCost(place int) int {
  c := 0
  for p, n := range ca.pos {
    steps := abs(place - p)
    c += steps * (steps+1) * n / 2
  }
  return c
}

// optimalFuelCost gets the optimal fuel cost of lining up the crab army, given the fuel cost function provided.
func (ca crabArmy) optimalFuelCost(costCalculator func(int) int) int {
  minCost := math.MaxInt64
  for p := ca.min; p <= ca.max; p++ {
    cost := costCalculator(p)
    if cost < minCost {
      minCost = cost
    }
  }
  return minCost
}

// extractCrabArmy extracts the crab army from the given input file.
func extractCrabArmy(fileName string) (crabArmy, error) {
	raw, err := os.ReadFile(fileName)
	if err != nil {
		return crabArmy{}, err
	}
	rs := strings.Split(strings.TrimSpace(string(raw)), ",")
	ca, err := parseCrabArmy(rs)
	if err != nil {
		return crabArmy{}, err
	}
	return ca, nil
}


func main() {
	ca, err := extractCrabArmy("input.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Input ingestion went wrong: ", err)
		os.Exit(1)
	}

    fmt.Println("Part 1:", ca.optimalFuelCost(ca.linearFuelCost))
    fmt.Println("Part 2:", ca.optimalFuelCost(ca.triangularFuelCost))
}
