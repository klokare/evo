package configurer

import (
	"io"
	"io/ioutil"
	"os"
)

func LoadFromFile(path string) (b []byte, err error) {

	var f *os.File
	if f, err = os.Open(path); err != nil {
		return
	}

	b, err = ioutil.ReadAll(f)
	if err == io.EOF {
		err = nil
	}
	return
}
