package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	raw, err := os.ReadFile("input.txt")
	if err != nil {
		fmt.Println("Input ingestion went wrong: ")
		fmt.Println(err)
		os.Exit(1)
	}
	strDepths := strings.Split(strings.TrimSpace(string(raw)), "\n")
	depths := make([]int, len(strDepths))
	for i, sd := range strDepths {
		d, err := strconv.Atoi(sd)
		if err != nil {
			fmt.Println("Could not convert to int: " + sd)
			os.Exit(1)
		}
		depths[i] = d
	}
	increases := 0
	prev := math.MaxInt64
	for i := range depths {
		if i >= len(depths)-2 {
			break
		}
		sum := depths[i] + depths[i+1] + depths[i+2]
		if sum > prev {
			increases += 1
		}
		prev = sum
	}
	fmt.Println(increases)
}
