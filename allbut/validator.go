package allbut

import (
	"fmt"
	"os"
)

type validator struct {
	s statter
}

type statter interface {
	Stat(string) (os.FileInfo, error)
}

type fileStatter struct {
}

func (fs *fileStatter) Stat(f string) (os.FileInfo, error) {
	return os.Stat(f)
}

func (v *validator) validate(files []string) error {
	for _, file := range files {
		if v.s == nil {
			v.s = &fileStatter{}
		}

		f, err := v.s.Stat(file)
		if err != nil {
			return fmt.Errorf("unable to read file: %s", f)
		}

		if f.IsDir() {
			return fmt.Errorf("%s is a directory, not a plain file", f.Name())
		}
	}
	return nil
}
