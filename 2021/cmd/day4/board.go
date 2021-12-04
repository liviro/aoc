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

type board [size][size]cell

// parseBoard parses and returns a new board from the raw input.
func parseBoard(raw string) (board, error) {
	rows := strings.Split(raw, "\n")
	var b board
	for i := 0; i < size; i++ {
		// TrimSpace to handle the right-aligned nums at index 0.
		row := strings.Fields(strings.TrimSpace(rows[i]))
		for j := 0; j < size; j++ {
			num, err := strconv.Atoi(row[j])
			if err != nil {
				return board{}, err
			}
			b[i][j] = cell{num: int64(num)}
		}
	}
	return b, nil
}

// mark marks a cell with the given value on the board as hit (if found).
func (board *board) mark(val int64) {
	for i := range board {
		for j := range board[i] {
			board[i][j].hit = board[i][j].hit || board[i][j].num == val
		}
	}
}

// hasWin determines whether the board has a win (fully hit column or row of board).
func (board *board) hasWin() bool {
Rows:
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if !board[i][j].hit {
				continue Rows
			}
		}
		return true
	}
Cols:
	for j := 0; j < size; j++ {
		for i := 0; i < size; i++ {
			if !board[i][j].hit {
				continue Cols
			}
		}
		return true
	}
	return false
}

// score gets the score of a board: sum of non-hit cells times the just-called value.
func (board *board) score(justCalled int64) int64 {
	var s int64
	for i := range board {
		for j := range board[i] {
			if !board[i][j].hit {
				s += board[i][j].num
			}
		}
	}
	return s * justCalled
}
