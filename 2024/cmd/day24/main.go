package main

import (
	"fmt"
	"math"
	"os"
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

func main() {
	t := time.Now()
	wires, rules := extractInput(os.Args[1])
	populateRules(wires, rules)
	fmt.Printf("Part 1: %d\n", computeZ(wires))
	fmt.Printf("Part2: %d\n", 0)
	fmt.Printf("Time elapsed: %s\n", time.Since(t))
}
