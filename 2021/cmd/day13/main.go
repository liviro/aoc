package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// fold denotes the direction (x or y) and location of a fold.
type fold struct {
	dir string
	loc int
}

// parseFold parses a fold out of a raw string of the format "fold along x=123".
func parseFold(raw string) (fold, error) {
	s := raw[11:] // Ignore "fold along " prefix
	pts := strings.Split(s, "=")
	loc, err := strconv.Atoi(pts[1])
	if err != nil {
		return fold{}, err
	}
	return fold{pts[0], loc}, nil
}

// dot denotes a dot at a given location x, y.
type dot struct{ x, y int }

// parseDot parses a dot out of a raw string of the format "x, y"
func parseDot(raw string) (dot, error) {
	cs := strings.Split(raw, ",")
	x, err := strconv.Atoi(cs[0])
	if err != nil {
		return dot{}, err
	}
	y, err := strconv.Atoi(cs[1])
	if err != nil {
		return dot{}, err
	}
	return dot{x, y}, nil
}

// postFold returns the new dot after a fold was applied.
func (d dot) postFold(f fold) dot {
	if f.dir == "x" {
		if f.loc >= d.x {
			return d
		}
		return dot{2*f.loc - d.x, d.y}
	}
	if f.dir == "y" {
		if f.loc >= d.y {
			return d
		}
		return dot{d.x, 2*f.loc - d.y}
	}
	panic("Unknown fold direction!")
}

// paper represents the collection of dots on the sheet of paper.
type paper map[dot]bool

// parsePaper parses a paper out of the raw input string.
func parsePaper(raw string) (paper, error) {
	rawDots := strings.Split(raw, "\n")
	p := make(paper)
	for _, rd := range rawDots {
		d, err := parseDot(rd)
		if err != nil {
			return nil, err
		}
		p[d] = true
	}
	return p, nil
}

// fold returns a new paper that is the result of the given fold.
func (p paper) fold(f fold) paper {
	np := make(paper)
	for d := range p {
		np[d.postFold(f)] = true
	}
	return np
}

// String() implements the Stringer interface for a paper, allowing it to be printed.
func (p paper) String() string {
	var maxX, maxY int
	for d := range p {
		if maxX < d.x {
			maxX = d.x
		}
		if maxY < d.y {
			maxY = d.y
		}
	}

	var s string
	for y := 0; y <= maxY; y++ {
		var r string
		for x := 0; x <= maxX; x++ {
			if p[dot{x, y}] {
				s += "#"
			} else {
				s += "."
			}
		}
		s += r + "\n"
	}
	return s
}

// dotsAfterFold returns how many dots the paper has after the given fold.
func (p paper) dotsAfterFold(f fold) int {
	np := p.fold(f)
	return len(np)
}

// afterFolds returns a new paper that is the result of applying the given folds to the paper.
func (p paper) afterFolds(fs []fold) paper {
	for _, f := range fs {
		p = p.fold(f)
	}
	return p
}

// extractInput extracts the input paper and folds from the given file.
func extractInput(fileName string) (paper, []fold, error) {
	raw, err := os.ReadFile(fileName)
	if err != nil {
		return nil, nil, err
	}

	blocks := strings.Split(strings.TrimSpace(string(raw)), "\n\n")

	p, err := parsePaper(blocks[0])
	if err != nil {
		return nil, nil, err
	}

	var fs []fold
	for _, rf := range strings.Split(blocks[1], "\n") {
		f, err := parseFold(rf)
		if err != nil {
			return nil, nil, err
		}
		fs = append(fs, f)
	}

	return p, fs, nil
}

func main() {
	p, fs, err := extractInput("input.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Input ingestion went wrong: ", err)
		os.Exit(1)
	}

	fmt.Println("Part 1:", p.dotsAfterFold(fs[0]))
	fmt.Println("Part 2:")
	fmt.Println(p.afterFolds(fs))

}
