package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// mustParseInt parses an int from the given string, or panics.
func mustParseInt(raw string) int {
	n, err := strconv.Atoi(raw)
	if err != nil {
		panic("Could not parse int")
	}
	return n
}

// location represents a location within a heightmap at row i and column j.
type location struct{ i, j int }

// isIn returns whether the given candidate location is in the given slice of locations.
func isIn(ls []location, candidate location) bool {
	for _, l := range ls {
		if l == candidate {
			return true
		}
	}
	return false
}

type heightmap [][]int

// at returns the height of the heightmap at the given location.
func (hm heightmap) at(l location) int {
	return hm[l.i][l.j]
}

// neighbors returns the neighboring (above, below, left, right) locations of the given location.
func (hm heightmap) neighbors(l location) []location {
	var ns []location
	// Above
	if l.i-1 >= 0 {
		ns = append(ns, location{l.i - 1, l.j})
	}
	// Left
	if l.j-1 >= 0 {
		ns = append(ns, location{l.i, l.j - 1})
	}
	// Below
	if l.i+1 < len(hm) {
		ns = append(ns, location{l.i + 1, l.j})
	}
	// Right
	if l.j+1 < len(hm[l.i]) {
		ns = append(ns, location{l.i, l.j + 1})
	}
	return ns
}

// isLocalMin returns true if the given location is lower than all its neighbors.
func (hm heightmap) isLocalMin(l location) bool {
	for _, n := range hm.neighbors(l) {
		if hm.at(l) >= hm.at(n) {
			return false
		}
	}
	return true
}

// lowPointRiskLevelSum returns the sum of the risk levels of the low points of the heightmap.
func (hm heightmap) lowPointRiskLevelSum() int {
	s := 0
	for i, _ := range hm {
		for j, p := range hm[i] {
			if hm.isLocalMin(location{i, j}) {
				s += p + 1
			}
		}
	}
	return s
}

// fullBasin returns all locations that are part of the same basin as the given start location.
func (hm heightmap) fullBasin(start location) []location {
	hms := len(hm) * len(hm[0])
	q := make(chan location, hms)
	bm := make(map[location]bool)
	qd := make(map[location]bool)
	q <- start
L:
	for {
		select {
		case l := <-q:
			bm[l] = true
			ns := hm.neighbors(l)
			for _, n := range ns {
				_, inBasin := bm[n]
				_, queued := qd[n]
				if !inBasin && !queued && hm.at(n) != 9 {
					q <- n
					qd[n] = true

				}
			}
		default:
			break L
		}
	}
	var basin []location
	for b := range bm {
		basin = append(basin, b)
	}
	return basin
}

// bigBasinsProduct returns the product of the sizes of the three biggest basins of the heightmap.
func (hm heightmap) bigBasinsProduct() int {
	var sizes []int
	var seen []location
	for i := range hm {
		for j := range hm[i] {
			l := location{i, j}
			if !isIn(seen, l) && hm.at(l) != 9 {
				b := hm.fullBasin(l)
				seen = append(seen, b...)
				sizes = append(sizes, len(b))
			}
		}
	}
	sort.Sort(sort.Reverse(sort.IntSlice(sizes)))
	return sizes[0] * sizes[1] * sizes[2]
}

// extractHeightmap returns the heighmap found in the input file.
func extractHeightmap(fileName string) (heightmap, error) {
	fp, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	var hm heightmap
	s := bufio.NewScanner(fp)
	s.Split(bufio.ScanLines)
	for s.Scan() {
		rawRow := strings.Split(s.Text(), "")
		var row []int
		for _, v := range rawRow {
			row = append(row, mustParseInt(v))
		}
		hm = append(hm, row)
	}
	return hm, nil
}

func main() {
	hm, err := extractHeightmap("input.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Input ingestion went wrong: ", err)
		os.Exit(1)
	}
	fmt.Println("Part 1:", hm.lowPointRiskLevelSum())
	fmt.Println("Part 2:", hm.bigBasinsProduct())
}
