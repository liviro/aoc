package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type elem struct {
	reg  int
	pair *pair
}

func (e *elem) String() string {
	if e.pair == nil {
		return fmt.Sprintf("%v", e.reg)
	}
	return fmt.Sprintf("%v", e.pair)
}

type pair struct {
	left  *elem
	right *elem
}

func strToPair(raw string) *pair {
	var rawPair []interface{}
	if err := json.Unmarshal([]byte(raw), &rawPair); err != nil {
		panic(err)
	}
	return mustParsePair(rawPair)
}

func mustParsePair(raw []interface{}) *pair {
	root := &pair{}
	root.left = mustParseElem(raw[0])
	root.right = mustParseElem(raw[1])
	return root
}

func mustParseElem(raw interface{}) *elem {
	switch e := raw.(type) {
	case float64:
		return &elem{reg: int(e)}
	case []interface{}:
		return &elem{pair: mustParsePair(e)}
	default:
		fmt.Printf("Don't know how to handle %v of type %T!\n", e, e)
		panic("AAAA")
	}
}

func (p *pair) String() string {
	return fmt.Sprintf("[%v, %v]", p.left, p.right)
}

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

func main() {
	ps, err := extractNumbers("input-test.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Input ingestion went wrong: ", err)
		os.Exit(1)
	}

	for _, p := range ps {
		fmt.Println(p)
	}
}
