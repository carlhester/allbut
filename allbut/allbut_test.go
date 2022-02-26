package allbut

import (
	"io/fs"
	"os"
	"testing"
	"time"
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
	s := sanitizer{}

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
			result := s.sanitizeFilenames(tt.input)
			if !slicesEqual(result, tt.output) {
				t.Errorf("Error.  Expected: [%s]. Got: [%s]", tt.output, result)
			}
		})
	}
}

// func TestSetup(t *testing.T) {
// 	tests := []struct {
// 		descr              string
// 		args               []string
// 		expectedToDelete   []string
// 		expectedToProtect  []string
// 		expectedDeleteFlag bool
// 		expectedErr        bool
// 	}{
// 		{"success case", []string{"-f", "testfile.txt"}, []string{}, []string{}, false, true},
// 		{"no files to protect returns just an error", []string{"-f"}, []string{}, []string{}, false, true},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.descr, func(t *testing.T) {
// 			result, err := Setup(tt.args)

// 			if !slicesEqual(result.toDelete, tt.expectedToDelete) {
// 				t.Errorf("error. Expected %+v, got %+v", tt.expectedToDelete, result.toDelete)
// 			}

// 			if !slicesEqual(result.toProtect, tt.expectedToProtect) {
// 				t.Errorf("error. Expected %+v, got %+v", tt.expectedToProtect, result.toProtect)
// 			}

// 			if result.deleteEnabled != tt.expectedDeleteFlag {
// 				t.Errorf("error. Expected %t, got %t", tt.expectedDeleteFlag, result.deleteEnabled)
// 			}
// 			if !tt.expectedErr {
// 				if err != nil {
// 					t.Errorf("error. Expected an error, got %v", err)
// 				}
// 			}
// 		})
// 	}
// }

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

type fakeFileInfo struct {
	name string
}

func (f *fakeFileInfo) Name() string {
	return f.name
}

func (f *fakeFileInfo) Size() int64 {
	return 0
}

func (f *fakeFileInfo) Mode() fs.FileMode {
	return 0
}

func (f *fakeFileInfo) ModTime() time.Time {
	return time.Now()
}

func (f *fakeFileInfo) IsDir() bool {
	return false
}
func (f *fakeFileInfo) Sys() interface{} {
	return 0
}

func TestIdentifyDeletionCandidates(t *testing.T) {
	fakeCwdFiles := []os.FileInfo{
		&fakeFileInfo{name: "fake1"},
		&fakeFileInfo{name: "fake2"},
		&fakeFileInfo{name: "fake3"},
		&fakeFileInfo{name: "fake4"},
	}

	tests := []struct {
		desc           string
		protectedFiles []string
		filesInCwd     []fs.FileInfo
		expectedResult []string
		expectedErr    error
	}{
		{"one matched file", []string{"fake1"}, fakeCwdFiles, []string{"fake2", "fake3", "fake4"}, nil},
		{"two matched files", []string{"fake1", "fake2"}, fakeCwdFiles, []string{"fake3", "fake4"}, nil},
		{"zero matched files", []string{"no_matched_file"}, fakeCwdFiles, []string{"fake1", "fake2", "fake3", "fake4"}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {

			result, err := identifyDeletionCandidates(tt.protectedFiles, tt.filesInCwd)
			if err != tt.expectedErr {
				t.Errorf("error. Expected: %+v. Got: %+v", tt.expectedErr, err)
			}

			if !slicesEqual(tt.expectedResult, result) {
				t.Errorf("error. Expected: %+v. Got: %+v", tt.expectedResult, result)
			}
		})
	}
}
