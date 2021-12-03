package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func extractReport(fileName string) ([]string, error) {
	fp, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	var r []string
	s := bufio.NewScanner(fp)
	s.Split(bufio.ScanLines)
	for s.Scan() {
		r = append(r, s.Text())
	}
	return r, s.Err()
}

func flip(b string) string {
	if b == "0" {
		return "1"
	} else {
		return "0"
	}
}

func safeStringToDecimal(b string) int64 {
	i, err := strconv.ParseInt(b, 2, 0)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not parse: ", b)
		os.Exit(1)
	}
	return i
}

func mostCommonCharAt(report []string, index int) string {
	zeros, ones := 0, 0
	for _, r := range report {
		switch r[index] {
		case '0':
			zeros += 1
		case '1':
			ones += 1
		default:
			fmt.Fprintln(os.Stderr, "Unknown binary digit found: ", r[index])
			os.Exit(1)
		}
	}
	if zeros > ones {
		return "0"
	} else {
		return "1"
	}
}

func getPower(report []string) int64 {
	var gs, es string
	for i := range report[0] {
		gc := mostCommonCharAt(report, i)
		gs += gc
		es += flip(gc)
	}
	g := safeStringToDecimal(gs)
	e := safeStringToDecimal(es)
	return g * e
}

func getRating(report []string, kind string) int64 {
	l := make([]string, len(report))
	copy(l, report)
	var i int
	for len(l) > 1 {
		var c string
		switch kind {
		case "oxygen":
			c = mostCommonCharAt(l, i)
		case "co2":
			c = flip(mostCommonCharAt(l, i))
		default:
			fmt.Fprintln(os.Stderr, "Unrecognized rating kind: ", kind)
			os.Exit(1)
		}
		var n []string
		for _, v := range l {
			if string(v[i]) == c {
				n = append(n, v)
			}
		}
		l = n
		i += 1
	}
	return safeStringToDecimal(l[0])
}

func getLifeSupportRating(report []string) int64 {
	return getRating(report, "oxygen") * getRating(report, "co2")
}

func main() {
	r, err := extractReport("3-input.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Input ingestion went wrong: ", err)
		os.Exit(1)
	}
	part1 := getPower(r)
	part2 := getLifeSupportRating(r)
	fmt.Println("Part 1: ", part1)
	fmt.Println("Part 2: ", part2)
}
