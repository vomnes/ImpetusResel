package net

import (
	"strconv"
	"strings"
)

// IP is the IP format
type IP []byte

func parseIPv4(s string) IP {
	var ip [4]byte
	bytesStr := strings.Split(s, ".")
	if bytesStr == nil || len(bytesStr) != 4 {
		return nil
	}
	for i := 0; i < len(bytesStr); i++ {
		n, err := strconv.Atoi(bytesStr[i])
		if err != nil || n > 0xFF || n < 0 {
			return nil
		}
		ip[i] = byte(n)
	}
	return IP{ip[0], ip[1], ip[2], ip[3]}
}

// ParseIP parses a IP address and return the []byte result
// If not a valid IP return nil
// Handle:
// - IPv4
func ParseIP(s string) IP {
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '.':
			return parseIPv4(s)
		}
	}
	return nil
}
