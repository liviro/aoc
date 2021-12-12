package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// A cave is a node in the tunnel graph.
type cave string

// start, end are the special caves in the tunnel graph that must be at the respective ends of a path.
const start = cave("start")
const end = cave("end")

// isBig returns whether or not a cave is big.
func (c cave) isBig() bool {
	return string(c) == strings.ToUpper(string(c))
}

// path is a path through caves.
type path []cave

// visited returns whether the given cave has already been visited in the path.
func (p path) visited(c cave) bool {
	for _, s := range p {
		if s == c {
			return true
		}
	}
	return false
}

// tunnelSystem is the cave graph.
// This is an adjacency map graph representation: each key's value is a list of caves it has a tunnel to.
type tunnelSystem map[cave][]cave

// addTunnel adds a connection between the two given caves to the existing tunnel system.
// start and end have special treatment: no tunnel ever ends at start cave or starts at the end cave.
func (t tunnelSystem) addTunnel(a, b cave) {
	if b != start && a != end {
		if _, ok := t[a]; !ok {
			t[a] = []cave{b}
		} else {
			t[a] = append(t[a], b)
		}
	}

	if a != start && b != end {
		if _, ok := t[b]; !ok {
			t[b] = []cave{a}
		} else {
			t[b] = append(t[b], a)
		}
	}
}

// extractTunnelSystem reads in the tunnelSystem from the given file.
func extractTunnelSystem(fileName string) (tunnelSystem, error) {
	fp, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	s := bufio.NewScanner(fp)
	s.Split(bufio.ScanLines)
	ts := make(tunnelSystem)
	for s.Scan() {
		ns := strings.Split(s.Text(), "-")
		ts.addTunnel(cave(ns[0]), cave(ns[1]))
	}
	return ts, nil
}

// onlyOnce returns whether the small cave can be visited after the path if small caves can be visited only once.
func onlyOnce(p path, c cave) bool {
	return !p.visited(c)
}

// revisitOnce returns whether the small cave can be visited after the path if one small cave can be visited twice.
func revisitOne(p path, c cave) bool {
	var hasSmallRevisit bool
	v := make(map[cave]int)
L:
	for _, s := range p {
		if !s.isBig() {
			v[s] += 1
			if v[s] > 1 {
				hasSmallRevisit = true
				break L
			}
		}
	}
	return !hasSmallRevisit || !p.visited(c)
}

// findPaths from returns all allowed paths that begin with the given path and end on the end cave.
// The smallCaveCheck determines whether a small cave may be visited next.
func findPathsFrom(t tunnelSystem, smallCaveCheck func(path, cave) bool, p path) []path {
	soFar := make(path, len(p))
	copy(soFar, p)

	last := soFar[len(soFar)-1]

	nexts := make([]cave, len(t[last]))
	copy(nexts, t[last])

	var ps []path
	for _, n := range nexts {
		if n == end {
			ps = append(ps, append(soFar, n))
		} else if n.isBig() || smallCaveCheck(soFar, n) {
			ps = append(ps, findPathsFrom(t, smallCaveCheck, append(soFar, n))...)
		}
	}
	return ps
}

// quickPathsCount returns the count of paths through the tunnel system which visit each small cave at most once.
func (t tunnelSystem) quickPathsCount() int {
	return len(findPathsFrom(t, onlyOnce, path([]cave{start})))
}

// slowPathsCount returns the count of paths through the tunnel system which allow one small cave to be revisited.
func (t tunnelSystem) slowPathsCount() int {
	return len(findPathsFrom(t, revisitOne, path([]cave{start})))
}

func main() {
	ts, err := extractTunnelSystem("input.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Input ingestion went wrong: ", err)
		os.Exit(1)
	}
	fmt.Println("Part 1:", ts.quickPathsCount())
	fmt.Println("Part 2:", ts.slowPathsCount())
}
