package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type command struct {
	cnt  int
	from int
	to   int
}

type stack struct {
	top *crate
}

type crate struct {
	val  string
	prev *crate
}

func (s *stack) push(v string) {
	c := crate{v, s.top}
	s.top = &c
}

func (s *stack) pop() string {
	v := s.top.val
	s.top = s.top.prev
	return v
}

func (s *stack) peek() string {
	if s.top == nil {
		return "0"
	}
	return s.top.val
}

func (s *stack) String() string {
	out := ""
	cr := s.top
	for cr != nil {
		out += cr.val + " "
		cr = cr.prev
	}
	return out
}

func extractStacks(raw string) ([]*stack, error) {
	rows := strings.Split(raw, "\n")
	rawCounts := strings.TrimSpace(rows[len(rows)-1])
	cnt, err := strconv.Atoi(rawCounts[strings.LastIndex(rawCounts, " ")+1:])
	if err != nil {
		return nil, err
	}
	stacks := []*stack{}
	for i := 0; i < cnt; i++ {
		stacks = append(stacks, &stack{})
	}

	for r := len(rows) - 2; r >= 0; r-- {
		for c := 0; c < cnt; c++ {
			loc := 1 + c*4
			if len(rows[r]) > loc && string(rows[r][loc]) != " " {
				stacks[c].push(string(rows[r][loc]))
			}
		}
	}
	return stacks, nil
}

func extractCommands(raw string) ([]*command, error) {
	cmds := []*command{}
	for _, r := range strings.Split(strings.TrimSpace(raw), "\n") {
		s := strings.Split(strings.TrimSpace(r), " ")
		cnt, err := strconv.Atoi(s[1])
		if err != nil {
			return nil, err
		}
		from, err := strconv.Atoi(s[3])
		if err != nil {
			return nil, err
		}
		to, err := strconv.Atoi(s[5])
		if err != nil {
			return nil, err
		}
		cmds = append(cmds, &command{cnt, from, to})
	}
	return cmds, nil
}

func extractInput(fileName string) ([]*stack, []*command, error) {
	raw, err := os.ReadFile(fileName)
	if err != nil {
		return nil, nil, err
	}

	rc := strings.Split(string(raw), "\n\n")
	stacks, err := extractStacks(rc[0])
	if err != nil {
		return nil, nil, err
	}
	cmds, err := extractCommands(rc[1])
	if err != nil {
		return nil, nil, err
	}
	return stacks, cmds, nil
}

func execute1(stacks []*stack, cmd command) {
	for i := 0; i < cmd.cnt; i++ {
		c := stacks[cmd.from-1].pop()
		stacks[cmd.to-1].push(c)
	}
}

func execute2(stacks []*stack, cmd command) {
	tmp := stack{}
	for i := 0; i < cmd.cnt; i++ {
		c := stacks[cmd.from-1].pop()
		tmp.push(c)
	}
	for i := 0; i < cmd.cnt; i++ {
		c := tmp.pop()
		stacks[cmd.to-1].push(c)
	}
}

func solve(stacks []*stack, cmds []*command, exe func([]*stack, command)) string {
	for _, c := range cmds {
		exe(stacks, *c)
	}
	s := ""
	for i := 0; i < len(stacks); i++ {
		s += stacks[i].peek()
	}
	return s
}

func main() {
	in := "input.txt"
	stacks, cmds, err := extractInput(in)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Input ingestion went wrong: ", err)
		os.Exit(1)
	}
	fmt.Println("Part 1:", solve(stacks, cmds, execute1))

	stacks, cmds, err = extractInput(in)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Input ingestion went wrong: ", err)
		os.Exit(1)
	}
	fmt.Println("Part 2:", solve(stacks, cmds, execute2))
}
