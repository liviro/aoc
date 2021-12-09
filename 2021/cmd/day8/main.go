package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

// digit is an ordered string of the segments indicators, concatenated.
// The ordering allows to compare digits even if the segments lit up in a different order.
type digit string

// parseDigit parses out a digit from the raw concatenated string of segments.
// The parsing is just ordering.
func parseDigit(raw string) digit {
	s := strings.Split(raw, "")
	sort.Strings(s)
	return digit(strings.Join(s, ""))
}

// segmentCount returns the number of segments used by the digit.
func (d digit) segmentCount() int {
	return len(d)
}

// contains returns whether the given digit fully contains the other digits' segments.
// For example, "0" contains "1".
func (d digit) contains(other digit) bool {
	for _, o := range other {
		if !strings.Contains(string(d), string(o)) {
			return false
		}
	}
	return true
}

// display contains the digits 0-9 of this display, and the int values of the 4-digit output.
type display struct {
	digits  [10]digit
	outputs [4]int
}

// extractOutputValue represents the display's output as a 4-digit int.
func (d display) extractOutputValue() int {
	return d.outputs[0]*1000 + d.outputs[1]*100 + d.outputs[2]*10 + d.outputs[3]
}

// parseDisplay extracts a display from the given raw string.
func parseDisplay(raw string) display {
	sections := strings.Split(raw, " | ")
	var ds []digit
	for _, p := range strings.Split(sections[0], " ") {
		ds = append(ds, parseDigit(p))
	}
	ordered := orderDigits(ds)
	var outputs [4]int
	for i, o := range strings.Split(sections[1], " ") {
		outputs[i] = determineDigit(&ordered, parseDigit(o))
	}
	return display{ordered, outputs}
}

// orderDigits orders the digits by the number they represent.
func orderDigits(all []digit) [10]digit {
	var ordered [10]digit

	unknown := make(map[digit]bool)
	for _, d := range all {
		unknown[d] = true
	}

	// Pluck out the digits defined by the number of segments they have.
	nu := make(map[digit]bool)
	for d, _ := range unknown {
		switch d.segmentCount() {
		case 2:
			ordered[1] = d
		case 4:
			ordered[4] = d
		case 3:
			ordered[7] = d
		case 7:
			ordered[8] = d
		default:
			nu[d] = true
		}
	}
	unknown = nu

	// Out of the remaining unknowns, only "9" has all segments of "4" + 2 more.
	for d, _ := range unknown {
		if d.segmentCount() == 6 && d.contains(ordered[4]) {
			ordered[9] = d
			break
		}
	}
	delete(unknown, ordered[9])

	// Out of the remaining unknowns, only "0" has all segments of "1" + 4 more.
	for d, _ := range unknown {
		if d.segmentCount() == 6 && d.contains(ordered[1]) {
			ordered[0] = d
			break
		}
	}
	delete(unknown, ordered[0])

	// Out of the remaining unknowns, only "6" has 6 segments.
	for d, _ := range unknown {
		if d.segmentCount() == 6 {
			ordered[6] = d
			break
		}
	}
	delete(unknown, ordered[6])

	// Out of the remaining unknowns, only "3" has all segments of "1" + 3 more.
	for d, _ := range unknown {
		if d.contains(ordered[1]) {
			ordered[3] = d
			break
		}
	}
	delete(unknown, ordered[3])

	// Out of the remaining unknowns, only "5" has all but one segments of "6".
	for d, _ := range unknown {
		if ordered[6].contains(d) {
			ordered[5] = d
			break
		}
	}
	delete(unknown, ordered[5])

	// Only unaccounted for left is "2".
	for d, _ := range unknown {
		ordered[2] = d
	}

	return ordered
}

// determineDigit returns the int value of the passed in unknown digit.
func determineDigit(ordered *[10]digit, unknown digit) int {
	for i, d := range ordered {
		if unknown == d {
			return i
		}
	}
	panic("got a truly unknown digit!")
}

// count1478 returns the number of instances when the display outputs a 1, 4, 7, or 8.
func (d display) count1478() int {
	c := 0
	for _, o := range d.outputs {
		switch o {
		case 1, 4, 7, 8:
			c += 1
		}
	}
	return c
}

// extractDisplays extracts displays from the given file.
func extractDisplays(fileName string) ([]display, error) {
	fp, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	var ds []display
	s := bufio.NewScanner(fp)
	s.Split(bufio.ScanLines)
	for s.Scan() {
		ds = append(ds, parseDisplay(s.Text()))
	}
	return ds, nil
}

// countAll1478 counts number of instances of 1, 4, 7, and 8 across all given displays.
func countAll1478(ds []display) int {
	c := 0
	for _, d := range ds {
		c += d.count1478()
	}
	return c
}

// sumAllOutputs returns the sum of the 4-digit int representation of the outputs of all displays.
func sumAllOutputs(ds []display) int {
	s := 0
	for _, d := range ds {
		s += d.extractOutputValue()
	}
	return s
}

func main() {
	ds, err := extractDisplays("input.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Input ingestion went wrong: ", err)
		os.Exit(1)
	}
	fmt.Println("Part 1:", countAll1478(ds))
	fmt.Println("Part 2:", sumAllOutputs(ds))
}
