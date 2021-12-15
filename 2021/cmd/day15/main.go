package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// max returns the largest of the inputs.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// min returns the smallest of the inputs.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// riskMap represents the 2-dimensional map of risks.
type riskMap [][]int

// extracRiskMap returns the risk map represented in the given input file.
func extractRiskMap(fileName string) (riskMap, error) {
	raw, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var rm riskMap
	for _, rr := range strings.Split(strings.TrimSpace(string(raw)), "\n") {
		var row []int
		for _, e := range rr {
			d, err := strconv.Atoi(string(e))
			if err != nil {
				return nil, err
			}
			row = append(row, d)
		}
		rm = append(rm, row)
	}
	return rm, nil
}

// blowUp returns the 5x blown up version of the risk map.
func (rm riskMap) blowUp() riskMap {
	s := len(rm)
	big := make([][]int, s*5)
	for i := range big {
		big[i] = make([]int, s*5)
	}

	for i := 0; i < s; i++ {
		for j := 0; j < s; j++ {

			for x := 0; x < 5; x++ {
				for y := 0; y < 5; y++ {
					v := rm[i][j] + x + y
					if v > 9 {
						v -= 9
					}
					big[x*s+i][y*s+j] = v
				}
			}
		}
	}
	return big
}

// lowestRiskPathCost returns the cost of the lowest-risk path from top left to bottom right in the risk map.
// This is done by building a memoization matrix for each position in the risk map by diagonal.
// If the need for a backtrack is noticed, the memoization matrix is re-built from the first diagonal.
func (rm riskMap) lowestRiskPathCost() int {
	size := len(rm)
	memo := make([][]int, size)
	for i := range memo {
		memo[i] = make([]int, size)
	}
	memo[0][0] = 0
	for d := 1; d < 2*size; {
		memoOnDiagonal(rm, memo, d)
		if diagonalOkWithoutBacktrack(rm, memo, d-1) {
			d += 1
		} else {
			// Conservatively start anew if backtracking is useful.
			// Probably possible to optimize...
			d = 1
		}
	}
	return memo[size-1][size-1]
}

// memoOnDiagonal updates the memoization 2D array for the risk map for all positions on the diagonal.
// A diagonal of a position is defined by the sum of its vertical and horizontal indices.
// Note that this considers all valid previous steps (incl. backtracking moves), if the memoization matrix
// has been filled out for them.
func memoOnDiagonal(rm riskMap, memo [][]int, d int) {
	size := len(rm)
	for i, j := max(0, d-size), min(d, size-1); i <= min(d, size-1); i, j = i+1, j-1 {
		m := math.MaxInt64
		if j != 0 {
			m = min(m, memo[i][j-1])
		}
		if i != 0 {
			m = min(m, memo[i-1][j])
		}
		if j != size-1 && memo[i][j+1] != 0 {
			m = min(m, memo[i][j+1])
		}
		if i != size-1 && memo[i+1][j] != 0 {
			m = min(m, memo[i+1][j])
		}

		memo[i][j] = m + rm[i][j]
	}
}

// diagonalOkWithoutBacktrack returns whether the positions on the diagonal could benefit from backtracking.
func diagonalOkWithoutBacktrack(rm riskMap, memo [][]int, d int) bool {
	size := len(memo)
	for i, j := max(0, d-size), min(d, size-1); i <= min(d, size-1); i, j = i+1, j-1 {
		if i < len(memo)-1 && memo[i+1][j] != 0 && memo[i][j] > memo[i+1][j]+rm[i][j] {
			return false
		}
		if j < len(memo)-1 && memo[i][j+1] != 0 && memo[i][j] > memo[i][j+1]+rm[i][j] {
			return false
		}
	}
	return true
}

func main() {
	rm, err := extractRiskMap("input.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Input ingestion went wrong: ", err)
		os.Exit(1)
	}

	fmt.Println("Part 1:", rm.lowestRiskPathCost())
	fmt.Println("Part 2:", rm.blowUp().lowestRiskPathCost())
}
