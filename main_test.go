package main

import "testing"

func TestParseArgs(t *testing.T) {
	tests := []struct {
		desc  string
		input []string
		res1  string
		res2  bool
	}{
		{"no args no results", []string{}, "", false},
		{"one file arg is just a file ", []string{"testfile"}, "testfile", false},
		{"two valid args works", []string{"testfile", "-f"}, "testfile", true},
		{"two invalid args", []string{"testfile", "-x"}, "testfile", false},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			res1, res2 := parseArgs(tt.input)
			if res1 != tt.res1 {
				t.Errorf("ERROR: got [%s], expected [%s]", res1, tt.res1)
			}
			if res2 != tt.res2 {
				t.Errorf("ERROR: got %t, expected %t", res2, tt.res2)
			}
		})
	}
}
