package parse

import "strconv"

// MustInt either parses the given string as an int, or panics.
func MustInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic("Could not parse int!")
	}
	return i
}
