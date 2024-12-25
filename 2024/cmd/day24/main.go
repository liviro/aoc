package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/liviro/aoc/2024/internal/parse"
)

type rule struct {
	in1, in2 string
	out      string
	op       string
}

func extractInput(name string) (map[string]int, []rule) {
	raw, err := os.ReadFile(name)
	if err != nil {
		panic("Cannot open file!")
	}
	sections := strings.Split(strings.TrimSpace(string(raw)), "\n\n")

	wires := map[string]int{}
	for _, raw := range strings.Split(sections[0], "\n") {
		var reg string
		var val int
		fmt.Sscanf(raw, "%s %d", &reg, &val)
		wires[strings.TrimSuffix(reg, ":")] = val
	}

	rules := []rule{}
	for _, raw := range strings.Split(sections[1], "\n") {
		var r rule
		fmt.Sscanf(raw, "%s %s %s -> %s", &r.in1, &r.op, &r.in2, &r.out)
		rules = append(rules, r)
	}
	return wires, rules
}

func populateRules(wires map[string]int, rules []rule) {
	for {
		updated := false
	R:
		for _, r := range rules {
			// Already processed, skip.
			if _, ok := wires[r.out]; ok {
				continue R
			}
			// Wires in haven't been seen yet, skip.
			if _, ok := wires[r.in1]; !ok {
				continue R
			}
			if _, ok := wires[r.in2]; !ok {
				continue R
			}
			switch {
			case r.op == "AND":
				wires[r.out] = wires[r.in1] & wires[r.in2]
			case r.op == "OR":
				wires[r.out] = wires[r.in1] | wires[r.in2]
			case r.op == "XOR":
				wires[r.out] = wires[r.in1] ^ wires[r.in2]
			}
			updated = true
		}
		if !updated {
			return
		}
	}
}

func computeZ(wires map[string]int) int {
	res := 0
	for w, v := range wires {
		if strings.HasPrefix(w, "z") && v == 1 {
			pos := parse.MustInt(strings.TrimPrefix(w, "z"))
			res += int(math.Pow(2, float64(pos)))
		}
	}
	return res
}

func setWires(x, y int) map[string]int {
	m := map[string]int{}
	for i := 0; i < 45; i++ {
		k := fmt.Sprintf("x%02d", i)
		b := x & 1
		m[k] = b
		x = x >> 1
	}
	for i := 0; i < 45; i++ {
		k := fmt.Sprintf("y%02d", i)
		b := y & 1
		m[k] = b
		y = y >> 1
	}
	return m
}

func countMistakes(rules []rule) int {
	m := 0
	for i := 0; i < 45; i++ {
		w0 := setWires(0, 1<<i)
		populateRules(w0, rules)
		if computeZ(w0) != 1<<i {
			m++
		}

		w1 := setWires(1<<i, 0)
		populateRules(w1, rules)
		if computeZ(w1) != 1<<i {
			m++
		}

		w2 := setWires(1<<i, 1<<i)
		populateRules(w2, rules)
		if computeZ(w2) != (1<<i + 1<<i) {
			m++
		}
	}
	return m
}

func swapRules(orig []rule, a, b int) []rule {
	rs := []rule{}
	for i := 0; i < len(orig); i++ {
		if i != a && i != b {
			rs = append(rs, orig[i])
		}
		if i == a {
			rs = append(rs, rule{
				in1: orig[a].in1, in2: orig[a].in2,
				op:  orig[a].op,
				out: orig[b].out,
			})
		}
		if i == b {
			rs = append(rs, rule{
				in1: orig[b].in1, in2: orig[b].in2,
				op:  orig[b].op,
				out: orig[a].out,
			})
		}
	}
	return rs
}

type swap struct {
	a, b     int
	mistakes int
}

// Return the swaps that return 0 mistakes for 0+2^i, 2^i+0, 2^i+2^i
func potentialSwaps(rules []rule) [][4]swap {
	baseline := countMistakes(rules)
	fmt.Printf("Baseline mistakes: %d\n", baseline)
	better := []swap{}
	for i := 0; i < len(rules); i++ {
		for j := i + 1; j < len(rules); j++ {
			m := countMistakes(swapRules(rules, i, j))
			if m+2 < baseline {
				better = append(better, swap{
					a:        i,
					b:        j,
					mistakes: m,
				})
			}
		}
	}
	slices.SortFunc(better, func(a, b swap) int {
		return a.mistakes - b.mistakes
	})
	for _, b := range better {
		fmt.Printf("Swapping %s, %s: %d mistakes\n", rules[b.a].out, rules[b.b].out, b.mistakes)
	}

	sws := [][4]swap{}
	for i := 0; i < len(better); i++ {
		for j := i + 1; j < len(better); j++ {
			for k := j + 1; k < len(better); k++ {
				for l := k + 1; l < len(better); l++ {
					s1 := better[i]
					s2 := better[j]
					s3 := better[k]
					s4 := better[l]
					gs := map[int]struct{}{}
					gs[s1.a] = struct{}{}
					gs[s1.b] = struct{}{}
					gs[s2.a] = struct{}{}
					gs[s2.b] = struct{}{}
					gs[s3.a] = struct{}{}
					gs[s3.b] = struct{}{}
					gs[s4.a] = struct{}{}
					gs[s4.b] = struct{}{}
					if len(gs) < 8 {
						continue
					}
					m := countMistakes(
						swapRules(
							swapRules(
								swapRules(
									swapRules(rules, s1.a, s1.b),
									s2.a, s2.b),
								s3.a, s3.b),
							s4.a, s4.b))
					if m == 0 {
						sws = append(sws, [4]swap{s1, s2, s3, s4})
						fmt.Printf("Swapping {%d, %d}, {%d, %d}, {%d, %d}, {%d, %d}: %d mistakes\n",
							s1.a, s1.b, s2.a, s2.b, s3.a, s3.b, s4.a, s4.b, m)
					}
				}
			}
		}
	}
	fmt.Println("Finished finding all worthy swaps.")
	return sws
}

func goodSwap(rules []rule, swaps [][4]swap) [4]swap {
	bad := []int{}
S:
	for {
		x := int(rand.Float64() * math.Pow(2, float64(40)))
		y := int(rand.Float64() * math.Pow(2, float64(40)))
		want := x + y
		for i, s := range swaps {
			if slices.Contains(bad, i) {
				continue
			}
			nr := swapRules(
				swapRules(
					swapRules(
						swapRules(rules, s[0].a, s[0].b),
						s[1].a, s[1].b),
					s[2].a, s[2].b),
				s[3].a, s[3].b)

			w := setWires(x, y)
			populateRules(w, nr)
			if computeZ(w) != want {
				bad = append(bad, i)
			}
		}
		if len(bad) == len(swaps)-1 {
			break S
		}
	}
	var g [4]swap
	for i := 0; i < len(swaps); i++ {
		if !slices.Contains(bad, i) {
			g = swaps[i]
		}
	}
	return g
}

func part2(rules []rule) string {
	ps := potentialSwaps(rules)
	ss := goodSwap(rules, ps)
	gates := []string{}
	for _, s := range ss {
		gates = append(gates, rules[s.a].out)
		gates = append(gates, rules[s.b].out)
	}
	sort.Strings(gates)
	return strings.Join(gates, ",")
}

func main() {
	t := time.Now()
	wires, rules := extractInput(os.Args[1])
	populateRules(wires, rules)
	fmt.Printf("Part 1: %d\n", computeZ(wires))
	fmt.Printf("Part 2: %s\n", part2(rules))
	fmt.Printf("Time elapsed: %s\n", time.Since(t))
}
