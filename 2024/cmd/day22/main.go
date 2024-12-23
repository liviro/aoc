package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/liviro/aoc/2024/internal/parse"
)

type secret struct {
	val     int
	history []int
	changes []int
}

func extractSecrets(name string) []secret {
	fp, err := os.Open(name)
	if err != nil {
		panic("Unable to open file")
	}
	defer fp.Close()
	s := bufio.NewScanner(fp)
	scs := []secret{}
	for s.Scan() {
		scs = append(scs, secret{
			val:     parse.MustInt(s.Text()),
			history: []int{},
			changes: []int{},
		})
	}
	return scs
}

func (s *secret) evolve() {
	prev := s.val % 10
	m := int(s.val) * 64
	s.mix(m)
	s.prune()

	d := int(float64(int(s.val) / 32))
	s.mix(d)
	s.prune()

	n := int(s.val) * 2048
	s.mix(n)
	s.prune()
	s.changes = append(s.changes, s.val%10-prev)
	s.history = append(s.history, s.val%10)
}

func (s *secret) mix(n int) {
	s.val = int(s.val) ^ n
}

func (s *secret) prune() {
	s.val = s.val % 16777216
}

func (s secret) price(seq []int) int {
	for i := 0; i < len(s.changes)-3; i++ {
		if seq[0] == s.changes[i] &&
			seq[1] == s.changes[i+1] &&
			seq[2] == s.changes[i+2] &&
			seq[3] == s.changes[i+3] {
			return s.history[i+3]
		}
	}
	return 0
}

func part1(scs []secret) int {
	r := 0
	for _, s := range scs {
		for i := 0; i < 2000; i++ {
			s.evolve()
		}
		r += s.val
	}
	return r
}

func key(seq []int) string {
	return fmt.Sprintf("%d,%d,%d,%d", seq[0], seq[1], seq[2], seq[3])
}

func part2(scs []secret) int {
	merchants := []*secret{}
	for _, s := range scs {
		merchants = append(merchants, &s)
	}
	for _, m := range merchants {
		for i := 0; i < 2000; i++ {
			m.evolve()
		}
	}
	priceMemo := map[string]int{}
	maxBananas := 0
	for si, s := range merchants {
		for i := 0; i < len(s.changes)-3; i++ {
			seq := []int{s.changes[i], s.changes[i+1], s.changes[i+2], s.changes[i+3]}
			if _, ok := priceMemo[key(seq)]; ok {
				continue
			}
			seqSum := 0
			for mi := si; mi < len(merchants); mi++ {
				seqSum += merchants[mi].price(seq)
			}
			priceMemo[key(seq)] = seqSum
			if seqSum > maxBananas {
				maxBananas = seqSum
			}
		}
	}
	return maxBananas
}

func main() {
	t := time.Now()
	secrets := extractSecrets(os.Args[1])
	fmt.Printf("Part1 : %d\n", part1(secrets))
	fmt.Printf("Part2 : %d\n", part2(secrets))
	fmt.Printf("Time elapsed: %s\n", time.Since(t))
}
