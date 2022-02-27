package allbut

import (
	"os"
	"testing"
)

type fakeStatter struct {
}

func (fs *fakeStatter) Stat(f string) (os.FileInfo, error) {
	return &fakeFileInfo{}, nil
}

func TestValidator(t *testing.T) {
	fakestatter := &fakeStatter{}

	v := &validator{
		s: fakestatter,
	}

	err := v.validate([]string{"test1"})
	if err != nil {
		t.Errorf("error. Expected: %+v. Got: %+v", nil, err)
	}

}
