package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"
)

func extractConns(name string) map[string][]string {
	fp, err := os.Open(name)
	if err != nil {
		panic("Unable to open file")
	}
	defer fp.Close()
	s := bufio.NewScanner(fp)
	conns := map[string][]string{}
	for s.Scan() {
		computers := strings.Split(s.Text(), "-")
		if cs, ok := conns[computers[0]]; ok {
			conns[computers[0]] = append(cs, computers[1])
		} else {
			conns[computers[0]] = []string{computers[1]}
		}

		if cs, ok := conns[computers[1]]; ok {
			conns[computers[1]] = append(cs, computers[0])
		} else {
			conns[computers[1]] = []string{computers[0]}
		}
	}
	return conns
}

func part1(clusters [][]string) int {
	count := 0
C:
	for _, c := range clusters {
		for _, l := range c {
			if strings.HasPrefix(l, "t") {
				count++
				continue C
			}
		}
	}
	return count
}

func cluster3(conns map[string][]string) [][]string {
	clusters := [][]string{}
	for comp, links := range conns {
		for i, l1 := range links {
			for j := i + 1; j < len(links); j++ {
				l2 := links[j]
				if slices.Contains(conns[l1], l2) {
					cluster := []string{comp, l1, l2}
					slices.Sort(cluster)
					if !slices.ContainsFunc(clusters, func(c []string) bool {
						return c[0] == cluster[0] && c[1] == cluster[1] && c[2] == cluster[2]
					}) {
						clusters = append(clusters, cluster)
					}
				}
			}
		}
	}
	return clusters
}

func clusterBig(conns map[string][]string) [][]string {
	clusters := [][]string{}
	for comp, links := range conns {
		for i, l1 := range links {
			cluster := []string{comp, l1}
			j := i + 1
		CB:
			for {
				if j == len(links) {
					break CB
				}
				ln := links[j]
				isGood := true
				for _, c := range cluster {
					if !slices.Contains(conns[c], ln) {
						isGood = false
					}
				}
				if isGood {
					cluster = append(cluster, ln)
					slices.Sort(cluster)
				}
				j++
			}
			if !slices.ContainsFunc(clusters, func(c []string) bool {
				if len(c) != len(cluster) {
					return false
				}
				for i := range c {
					if c[i] != cluster[i] {
						return false
					}
				}
				return true
			}) {
				clusters = append(clusters, cluster)
			}
		}
	}
	return clusters
}

func part2(clusters [][]string) string {
	biggest := clusters[0]
	for _, c := range clusters {
		if len(c) > len(biggest) {
			biggest = c
		}
	}
	return strings.Join(biggest, ",")
}

func main() {
	t := time.Now()
	connections := extractConns(os.Args[1])
	clusters := cluster3(connections)
	fmt.Printf("Part1 : %d\n", part1(clusters))
	fmt.Printf("Part2 : %s\n", part2(clusterBig(connections)))
	fmt.Printf("Time elapsed: %s\n", time.Since(t))
}
