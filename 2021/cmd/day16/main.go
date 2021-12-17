package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// rawPacket is the raw, binary string representation of a packet.
type rawPacket string

// version returns the version of the given raw packet.
func (p rawPacket) version() int64 {
	return mustParseBinInt(string(p[:3]))
}

// typeID returns the type ID of the given raw packet.
func (p rawPacket) typeID() int64 {
	return mustParseBinInt(string(p[3:6]))
}

// literal returns the literal of the given raw packet.
func (p rawPacket) literal() (int64, int64) {
	if p.typeID() != 4 {
		panic("Trying to get a literal from a non-literal packet!")
	}
	rs := string(p[6:])
	s := ""
	var l int64
	for {
		s += rs[1:5]
		l += 5
		if rs[0] == '0' {
			break
		}
		rs = rs[5:]
	}
	return mustParseBinInt(s), l
}

// packet is the parsed out, structured representation of a packet.
// Note that literal and subPackets are mutually exclusive: a packet only has one or the other set.
type packet struct {
	version    int64
	typeID     int64
	literal    int64
	subPackets []packet
}

// mustParseBinInt either parses a string as the binary representation of an integer, or panics.
func mustParseBinInt(s string) int64 {
	n, err := strconv.ParseInt(s, 2, 64)
	if err != nil {
		panic("Could not parse binary int!")
	}
	return n
}

// parsePacket yields the parsed out version of a raw packet and its length (in raw binary form).
func parsePacket(p rawPacket) (packet, int64) {
	v := p.version()
	tID := p.typeID()
	if tID == 4 {
		l, ll := p.literal()
		return packet{v, tID, l, nil}, 6 + ll
	} else {
		pps, ll := p.subPackets()
		return packet{v, tID, 0, pps}, 7 + ll
	}
}

// subPackets returns the parsed subpackets of the given raw packet, and their combined length.
func (p rawPacket) subPackets() ([]packet, int64) {
	lengthTypeID := p[6]
	if lengthTypeID == '0' {
		subPacketBits := mustParseBinInt(string(p[7:22]))
		rawSubPackets := p[22 : 22+subPacketBits]
		pps, ppl := parsePackets(rawSubPackets, nil)
		return pps, ppl + 15
	} else {
		count := mustParseBinInt(string(p[7:18]))
		pps, ppl := parsePackets(p[18:], &count)
		return pps, ppl + 11
	}
}

// parsePackets returns the list of parsed packets contained in the raw packet and their combined length.
// If a max. count of packets are specified, the parser stops there. Otherwise, it goes to the end of the packet.
func parsePackets(p rawPacket, max *int64) ([]packet, int64) {
	var ps []packet
	var l int64
	var n int64
	if max != nil {
		n = *max
	}
	for len(p) != 0 && (max == nil || n > 0) {
		next, nl := parsePacket(p)
		ps = append(ps, next)
		l += nl
		p = p[nl:]
		n -= 1
	}
	return ps, l
}

// extractTransmission extracts the raw packet in the given file.
func extractTransmission(fileName string) (rawPacket, error) {
	raw, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	s := ""
	for _, c := range strings.TrimSpace(string(raw)) {
		h, err := strconv.ParseUint(string(c), 16, 8)
		if err != nil {
			return "", err
		}
		s += fmt.Sprintf("%04b", h)
	}
	return rawPacket(s), nil
}

// versionSum returns the sum of all the versions found within a packet and its subpackets.
func (p packet) versionSum() int64 {
	s := p.version
	for _, sp := range p.subPackets {
		s += sp.versionSum()
	}
	return s
}

// eval returns the result of evaluating the expression in a given packet.
func (p packet) eval() int64 {
	switch p.typeID {
	case 0:
		s := int64(0)
		for _, sp := range p.subPackets {
			s += sp.eval()
		}
		return s

	case 1:
		s := int64(1)
		for _, sp := range p.subPackets {
			s *= sp.eval()
		}
		return s

	case 2:
		m := int64(math.MaxInt64)
		for _, sp := range p.subPackets {
			e := sp.eval()
			if e < m {
				m = e
			}
		}
		return m

	case 3:
		m := int64(0)
		for _, sp := range p.subPackets {
			e := sp.eval()
			if e > m {
				m = e
			}
		}
		return m

	case 4:
		return p.literal

	case 5:
		if p.subPackets[0].eval() > p.subPackets[1].eval() {
			return 1
		} else {
			return 0
		}

	case 6:
		if p.subPackets[0].eval() < p.subPackets[1].eval() {
			return 1
		} else {
			return 0
		}

	case 7:
		if p.subPackets[0].eval() == p.subPackets[1].eval() {
			return 1
		} else {
			return 0
		}

	default:
		panic("Unknown type ID!")
	}
}

func main() {
	ft, err := extractTransmission("input.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Input ingestion went wrong: ", err)
		os.Exit(1)
	}

	p, _ := parsePacket(ft)
	fmt.Println("Part 1:", p.versionSum())
	fmt.Println("Part 2:", p.eval())
}
