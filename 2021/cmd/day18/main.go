package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"os"
)

// elem denotes a single entity in a pair: either a regular integer, or a pair.
type elem struct {
	reg  int
	pair *pair
}

// String returns the prinable representation of an element.
func (e *elem) String() string {
	if e.pair == nil {
		return fmt.Sprintf("%v", e.reg)
	}
	return fmt.Sprintf("%v", e.pair)
}

// pair denotes a snailfish number, with left & right children and a reference to its parent, if applicable.
type pair struct {
	left   *elem
	right  *elem
	parent *pair
}

// strToPair parses a pair from its string representation.
func strToPair(raw string) *pair {
	var rawPair []interface{}
	if err := json.Unmarshal([]byte(raw), &rawPair); err != nil {
		panic(err)
	}
	return mustParsePair(rawPair, nil)
}

// mustParsePair parses a pair out of a slice of raw interfaces and a pointer to the parent.
func mustParsePair(raw []interface{}, parent *pair) *pair {
	root := &pair{parent: parent}
	root.left = mustParseElem(raw[0], root)
	root.right = mustParseElem(raw[1], root)
	return root
}

// mustParseElem parses out an element (regular integer or pair) out of a raw interface and a pointer to the parent.
func mustParseElem(raw interface{}, root *pair) *elem {
	switch e := raw.(type) {
	case float64:
		return &elem{reg: int(e)}
	case []interface{}:
		return &elem{pair: mustParsePair(e, root)}
	default:
		panic("Unexpected raw type!")
	}
}

// String returns the printable string representation of the pair.
func (p *pair) String() string {
	return fmt.Sprintf("[%v,%v]", p.left, p.right)
}

// deepCopy returns a deep copy of the given pair.
// Dirty? Sure. But it works.
func (p *pair) deepCopy() *pair {
	s := p.String()
	return strToPair(s)
}

// addToLeftRegNeighbor adds the given value to the leftmost regular neighbor of the given pair.
// If the pair is leftmost in its number, this is no-op.
func (p *pair) addToLeftRegNeighbor(v int) {
	curr := p
	// Climb up until there is a left sibling.
	for curr.parent != nil {
		parent := curr.parent
		if parent.right.pair == curr {
			curr = parent
			// If the left sibling is a regular number, add to it and exit.
			if curr.left.pair == nil {
				curr.left.reg += v
				return
			}
			// If the left sibling is a pair, climb down to its rightmost regular integer.
			curr = curr.left.pair
			for curr.right.pair != nil {
				curr = curr.right.pair
			}
			curr.right.reg += v
			return
		}
		curr = parent
	}
}

// addToRightRegNeighbor adds the given value to the rightmost regular neighbor of the given pair.
// If the pair is rightmost in its number, this is no-op.
func (p *pair) addToRightRegNeighbor(v int) {
	curr := p
	// Climb up until there is a right sibling.
	for curr.parent != nil {
		parent := curr.parent
		if parent.left.pair == curr {
			curr = parent
			// If the right sibling is a regular number, add to it and exit.
			if curr.right.pair == nil {
				curr.right.reg += v
				return
			}
			// If the right sibling is a pair, climb down to its leftmost regular integer.
			curr = curr.right.pair
			for curr.left.pair != nil {
				curr = curr.left.pair
			}
			curr.left.reg += v
			return
		}
		curr = parent
	}
}

// pairToExplode finds a pair that should be exploded.
// Note that it is possible for no child pairs to need exploding: if this is the case, return is nil.
func (p *pair) pairToExplode(depth int) *pair {
	if depth == 4 {
		return p
	}
	if p.left.pair != nil {
		e := p.left.pair.pairToExplode(depth + 1)
		if e != nil {
			return e
		}
	}
	if p.right.pair != nil {
		e := p.right.pair.pairToExplode(depth + 1)
		if e != nil {
			return e
		}
	}
	return nil
}

// explode explodes the leftmost explodable pair in the given number and returns true.
// If no child pair is explodable, no-ops and returns false.
func (p *pair) explode() bool {
	target := p.pairToExplode(0)
	if target == nil {
		return false
	}
	target.addToLeftRegNeighbor(target.left.reg)
	target.addToRightRegNeighbor(target.right.reg)

	parent := target.parent
	if parent.left.pair == target {
		parent.left.pair = nil
		parent.left.reg = 0
	}
	if parent.right.pair == target {
		parent.right.pair = nil
		parent.right.reg = 0
	}

	return true
}

// split splits the leftmost splittable pair in the given number and returns true.
// If no child pair is splittable, no-ops and returns false.
func (p *pair) split() bool {
	if p.left.pair == nil && p.left.reg >= 10 {
		p.left.pair = &pair{
			left:   &elem{reg: int(math.Floor(float64(p.left.reg) / 2))},
			right:  &elem{reg: int(math.Ceil(float64(p.left.reg) / 2))},
			parent: p,
		}
		p.left.reg = 0
		return true
	}

	if p.left.pair != nil && p.left.pair.split() {
		return true
	}

	if p.right.pair == nil && p.right.reg >= 10 {
		p.right.pair = &pair{
			left:   &elem{reg: int(math.Floor(float64(p.right.reg) / 2))},
			right:  &elem{reg: int(math.Ceil(float64(p.right.reg) / 2))},
			parent: p,
		}
		p.right.reg = 0
		return true
	}

	if p.right.pair != nil && p.right.pair.split() {
		return true
	}

	return false
}

// reduce reduces a given number by exploding and splitting until neither operations do anything.
func (p *pair) reduce() {
	for p.explode() || p.split() {
	}
}

// magnitude returns the magnitude of the given pair.
func (p *pair) magnitude() int {
	var left int
	if p.left.pair == nil {
		left = p.left.reg
	} else {
		left = p.left.pair.magnitude()
	}

	var right int
	if p.right.pair == nil {
		right = p.right.reg
	} else {
		right = p.right.pair.magnitude()
	}

	return 3*left + 2*right
}

// add adds the two pairs and reduces the result.
// Note that this operation does not modify the parameters.
func add(p1, p2 *pair) *pair {
	root := &pair{}
	l := p1.deepCopy()
	r := p2.deepCopy()
	l.parent = root
	r.parent = root
	root.left = &elem{pair: l}
	root.right = &elem{pair: r}
	root.reduce()
	return root
}

// extractNumbers extracts the input numbers from the given file.
func extractNumbers(fileName string) ([]*pair, error) {
	fp, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	s := bufio.NewScanner(fp)
	var ps []*pair
	for s.Scan() {
		ps = append(ps, strToPair(s.Text()))
	}
	return ps, nil
}

// sumPairs sums all given pairs.
func sumPairs(ps []*pair) *pair {
	var s *pair = ps[0]
	for i, p := range ps {
		if i == 0 {
			continue
		}
		s = add(s, p)
	}
	return s
}

// maxSumMagnitude returns the highest magnitude from summing any two of the input pairs.
func maxSumMagnitude(ps []*pair) int {
	max := 0
	for i, p1 := range ps {
		for j, p2 := range ps {
			if i == j {
				continue
			}
			m := add(p1, p2).magnitude()
			if m > max {
				max = m
			}
		}
	}
	return max
}

func main() {
	ps, err := extractNumbers("input.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Input ingestion went wrong: ", err)
		os.Exit(1)
	}
	fmt.Println("Part 1:", sumPairs(ps).magnitude())
	fmt.Println("Part 2:", maxSumMagnitude(ps))

}
