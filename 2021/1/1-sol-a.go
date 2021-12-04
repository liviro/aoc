package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	raw, err := os.ReadFile("1-input-test.txt")
	if err != nil {
		fmt.Println("Input ingestion went wrong: ")
		fmt.Println(err)
	}
	strDepths := strings.Split(strings.TrimSpace(string(raw)), "\n")
	depths := make([]int, len(strDepths))
	for i, sd := range strDepths {
		d, err := strconv.Atoi(sd)
		if err != nil {
			fmt.Println("Could not convert to int: " + sd)
		}
		depths[i] = d

	}
	increases := 0
	prev := math.MaxInt64
	for _, d := range depths {
		if d > prev {
			increases += 1
		}
		prev = d
	}
	fmt.Println(increases)
}
