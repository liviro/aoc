package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/liviro/aoc/2024/internal/parse"
)

type block struct {
	idx, len int
}

func extractDiskMap(name string) []int {
	fp, err := os.Open(name)
	if err != nil {
		panic("Unable to open file")
	}
	defer fp.Close()
	s := bufio.NewScanner(fp)
	res := []int{}
	for s.Scan() {
		for _, v := range s.Text() {
			res = append(res, parse.MustInt(string(v)))
		}
	}
	return res
}

func individualCheckSum(diskMap []int) int {
	cs := 0
	fwd := 0
	rew := len(diskMap) - 1
	pos := 0
	compSpace := 0
	compProg := 0
F:
	for {
		// Exit conditions: fwd and rew have met
		if fwd == rew {
			for i := 0; i < diskMap[fwd]-compProg; i++ {
				cs += (fwd / 2) * pos
				pos++
			}
			break F
		}
		// Add uncompacted section
		id := fwd / 2
		for i := 0; i < diskMap[fwd]; i++ {
			cs += id * pos
			pos++
		}
		// Compact from behind. Make sure to only fill up
		// the space allowed.
	B:
		for {
			id = rew / 2
			// If we finished putting in the last block, move to previous.
			if compProg == diskMap[rew] {
				rew -= 2
				compProg = 0
				continue B
			}
			// If we run out of compacting space, move to next uncompacted
			// section.
			if compSpace == diskMap[fwd+1] {
				fwd += 2
				compSpace = 0
				continue F
			}
			cs += id * pos
			pos++
			compSpace++
			compProg++
		}
	}
	return cs
}

func sortedSpaceKeys(s map[int]int) []int {
	ks := []int{}
	for k := range s {
		ks = append(ks, k)
	}
	sort.Ints(ks)
	return ks
}

func blockCheckSum(diskMap []int) int {
	m := map[int]block{}
	s := 0
	for i := 0; i < len(diskMap); i++ {
		if i%2 == 0 {
			m[i/2] = block{
				idx: s,
				len: diskMap[i],
			}
		}
		s += diskMap[i]
	}
	// index->length available
	spaces := map[int]int{}
	for i := 0; i+1 < len(m); i++ {
		// Space after a file block starts at the block's index + its length.
		// Length of space is delta between the index of the next block and
		// the index of the space.
		spaceIdx := m[i].idx + m[i].len
		spaces[spaceIdx] = m[i+1].idx - spaceIdx
	}
	// Remove any empty-length spaces
	for k, v := range spaces {
		if v == 0 {
			delete(spaces, k)
		}
	}
	compacted := map[int]block{}

	// Compact.
	rew := len(m) - 1
	for {
		if rew < 0 {
			break
		}
		if _, ok := m[rew]; !ok {
			continue
		}
	S:
		for _, si := range sortedSpaceKeys(spaces) {
			sl := spaces[si]
			// If the space is ahead of the block index, give up
			// for the block: can't move it back.
			if si > m[rew].idx {
				break S
			}
			if sl >= m[rew].len {
				// Add the compacted block
				compacted[rew] = block{
					idx: si,
					len: m[rew].len,
				}
				// If there's space left over, update its index and lengths
				// to account for the newly moved block
				if sl > m[rew].len {
					spaces[si+m[rew].len] = sl - m[rew].len
				}
				delete(spaces, si)
				delete(m, rew)
				break S
			}
		}
		rew--
	}

	cs := 0
	for id, v := range m {
		for i := 0; i < v.len; i++ {
			cs += id * (v.idx + i)
		}
	}
	for id, v := range compacted {
		for i := 0; i < v.len; i++ {
			cs += id * (v.idx + i)
		}
	}
	return cs
}

func main() {
	t := time.Now()
	dm := extractDiskMap(os.Args[1])
	fmt.Printf("Part 1: %d\n", individualCheckSum(dm))
	fmt.Printf("Part 2: %d\n", blockCheckSum(dm))
	fmt.Printf("Time elapsed: %s\n", time.Since(t))
}
