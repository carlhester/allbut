package main

import (
	"fmt"
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		desc  string
		input []string
		res1  []string
		res2  bool
	}{
		{"no args returns no results", []string{}, []string{}, false},
		{"one arg with flag returns results", []string{"testfile", "-f"}, []string{"testfile"}, true},
		{"two files with no flag returns files", []string{"testfile", "testfile2"}, []string{"testfile", "testfile2"}, false},
		{"two files with flag last returns files with flag", []string{"testfile", "testfile2", "-f"}, []string{"testfile", "testfile2"}, true},
		{"two files with flag second returns files with flag", []string{"testfile", "-f", "testfile2"}, []string{"testfile", "testfile2"}, true},
		{"two files with flag first returns files with flag", []string{"-f", "testfile2", "testfile1"}, []string{"testfile2", "testfile1"}, true},
		{"two files with multiple flags returns files with flag", []string{"-f", "-f", "testfile2", "testfile1", "-f"}, []string{"testfile2", "testfile1"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			fmt.Println(tt)
			res1, res2 := parseArgs(tt.input)
			if !slicesEqual(res1, tt.res1) {
				t.Errorf("ERROR: got [%s], expected [%s]", res1, tt.res1)
			}
			if res2 != tt.res2 {
				t.Errorf("ERROR: got %t, expected %t", res2, tt.res2)
			}
		})
	}
}

// slicesEqual compares two slices
func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
