package net

import (
	"bytes"
	"testing"
)

var IP4Tests = []struct {
	str         string // input
	expected    IP     // expected result
	testContent string // test details
}{
	{"37.169.43.146", IP{37, 169, 43, 146}, "Valid IPv4"},
	{"37.169.43146", nil, "Invalid IPv4"},
}

func TestParseIP(t *testing.T) {
	for _, tt := range IP4Tests {
		actual := ParseIP(tt.str)
		if bytes.Compare(actual, tt.expected) != 0 {
			t.Errorf("ParseIP(%s): expect [% x], has [% x] - Test type: \033[31m%s\033[0m",
				tt.str, tt.expected, actual, tt.testContent)
		}
	}
}
