package allbut

import (
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

func TestSanitizeFileNames(t *testing.T) {
	tests := []struct {
		descr  string
		input  []string
		output []string
	}{
		{"one name", []string{"test1"}, []string{"./test1"}},
		{"two names", []string{"test1", "test2"}, []string{"./test1", "./test2"}},
		{"one with prefix", []string{"./test1"}, []string{"./test1"}},
	}
	for _, tt := range tests {
		t.Run(tt.descr, func(t *testing.T) {
			result := sanitizeFileNames(tt.input)
			if !slicesEqual(result, tt.output) {
				t.Errorf("Error.  Expected: [%s]. Got: [%s]", tt.output, result)
			}
		})
	}
}

func TestAddDotSlashPrefix(t *testing.T) {
	tests := []struct {
		descr  string
		input  string
		output string
	}{
		{"no prefix", "test1", "./test1"},
		{"existing prefix", "./test1", "./test1"},
		{"wrong prefix", ".test1", "./.test1"},
	}
	for _, tt := range tests {
		t.Run(tt.descr, func(t *testing.T) {
			result := addDotSlashPrefix(tt.input)
			if result != tt.output {
				t.Errorf("Error.  Expected: [%s]. Got: [%s]", tt.output, result)
			}
		})
	}
}

func TestRemoveInvalidChars(t *testing.T) {
	tests := []struct {
		descr  string
		input  string
		output string
	}{
		{"no removal", "test1", "test1"},
		{"one bad char at end of string", "test1[", "test1"},
		{"one bad char at start of string", "[test1", "test1"},
		{"one bad char at middle of string", "te[st1", "test1"},
		{"two bad char at start of string", "[[test1", "test1"},
		{"two bad char at mid of string", "te[s[t1", "test1"},
		{"only bad chars", "[~,\\", ""},
	}
	for _, tt := range tests {
		t.Run(tt.descr, func(t *testing.T) {
			result := removeInvalidChars(tt.input)
			if result != tt.output {
				t.Errorf("Error. Expected: [%s]. Got: [%s]", tt.output, result)
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