package source

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
)

// NewJSON returns a new map source based on the JSON data stored in the reader.
func NewJSON(r io.Reader) (m Map, err error) {

	// Decode the data in the reader
	var b []byte
	if b, err = ioutil.ReadAll(r); err != nil {
		return
	}
	err = json.Unmarshal(b, &m)
	return
}

// NewJSONFromFile is a convenience function that creates a new JSON source from the file specified
// by filename.
func NewJSONFromFile(filename string) (m Map, err error) {

	// Open the file
	var f *os.File
	if f, err = os.Open(filename); err != nil {
		return
	}

	// Decode the data
	if m, err = NewJSON(f); err != nil {
		f.Close() // ignore error as it would overwrite the decoding one
		return
	}

	// Close the file and return
	err = f.Close()
	return
}
