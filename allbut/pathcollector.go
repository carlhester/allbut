package allbut

import (
	"io/ioutil"
	"os"
)

type cwdCollector struct {
}

func (c *cwdCollector) collect() ([]os.FileInfo, error) {
	return ioutil.ReadDir("./")
}
