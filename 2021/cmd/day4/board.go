package main

import (
	"strconv"
	"strings"
)

const size = 5

type cell struct {
	num int64
	hit bool
}

type board struct {
	b [][]cell
}

// Parses and returns a new board from the raw input.
func parseBoard(raw string) (board, error) {
	rows := strings.Split(raw, "\n")
	b := make([][]cell, size)
	for i := 0; i < size; i++ {
		b[i] = make([]cell, size)

		row := strings.Fields(strings.TrimSpace(rows[i]))
		for j := 0; j < size; j++ {
			num, err := strconv.Atoi(row[j])
			if err != nil {
				return board{}, err
			}
			b[i][j] = cell{num: int64(num)}
		}
	}
	return board{b: b}, nil
}

// Marks a cell with the given value on the board as hit (if found).
func (board board) mark(val int64) {
	for i := range board.b {
		for j := range board.b[i] {
			if board.b[i][j].num == val {
				board.b[i][j].hit = true
			}
		}
	}
}

// Determines whether the board has a win (fully hit column or board).
func (board board) hasWin() bool {
	// Check rows
	for i := 0; i < size; i++ {
		w := true
		for j := 0; j < size; j++ {
			w = w && board.b[i][j].hit
		}
		if w {
			return true
		}
	}
	// Check cols
	for j := 0; j < size; j++ {
		w := true
		for i := 0; i < size; i++ {
			w = w && board.b[i][j].hit
		}
		if w {
			return true
		}
	}
	return false
}

// Gets the score of a board: sum of non-hit cells times the just-called value.
func (board board) score(justCalled int64) int64 {
	var s int64
	for i := range board.b {
		for j := range board.b[i] {
			if !board.b[i][j].hit {
				s += board.b[i][j].num
			}
		}
	}
	return s * justCalled
}
