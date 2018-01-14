package json

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"
)

type Simple struct {
	Name string
}

type SimpleWithTag struct {
	Name string `json:"name"`
}

type Nested struct {
	Data Simple
}

type NestedAnonymous struct {
	Simple
}

func TestConfigureFailedReader(t *testing.T) {
	r := &MockReader{HasError: true}
	_, err := New(r)
	if err == nil {
		t.Errorf("error expected for failed reader")
	}
}

func TestConfigureFailedFile(t *testing.T) {
	_, err := NewFromFile("boaty-mcboatface.json")
	if err == nil {
		t.Errorf("error expected for failed file")
	}
}

func TestConfigure(t *testing.T) {

	var tests = []struct {
		Source        []byte
		Item          interface{}
		Expected      interface{}
		IsConfigError bool
	}{
		// Empty source
		{
			Item:          &Simple{},
			Expected:      &Simple{},
			IsConfigError: false,
		},
		// Cannot update item
		{
			Source:        []byte(`{"Name": "Hemmingway"}`),
			Item:          Simple{},
			Expected:      &Simple{},
			IsConfigError: true,
		},
		// Simple
		{
			Source:        []byte(`{"Name": "Hemmingway"}`),
			Item:          &Simple{},
			Expected:      &Simple{Name: "Hemmingway"},
			IsConfigError: false,
		},
		// Simple with tag
		{
			Source:        []byte(`{"name": "Hemmingway"}`),
			Item:          &SimpleWithTag{},
			Expected:      &SimpleWithTag{Name: "Hemmingway"},
			IsConfigError: false,
		},
		// Nested
		{
			Source:        []byte(`{"Name": "Hemmingway"}`),
			Item:          &Nested{},
			Expected:      &Nested{}, // JSON does not support nested, named properties
			IsConfigError: false,
		},
		// Nested with anonymous member
		{
			Source:        []byte(`{"Name": "Hemmingway"}`),
			Item:          &NestedAnonymous{},
			Expected:      &NestedAnonymous{Simple: Simple{Name: "Hemmingway"}},
			IsConfigError: false,
		},
	}

	for _, test := range tests {

		// Test from reader
		t.Run("Reader", fromReader(test.Source, test.Item, test.Expected, test.IsConfigError))

		// Test from reader
		t.Run("File", fromFile(test.Source, test.Item, test.Expected, test.IsConfigError))

	}

}

func fromReader(source []byte, item, expected interface{}, isError bool) func(*testing.T) {
	return func(t *testing.T) {

		// Create the configurer
		cfg, err := New(bytes.NewBuffer(source))
		if err != nil {
			t.Errorf("error not expected creating configurer. actual was %v", err)
			t.FailNow()
		}

		// Configure the item
		err = cfg.Configure(item)
		t.Run("Error", testError(isError, err))
		if err != nil {
			return
		}

		// Compare the results
		if !reflect.DeepEqual(item, expected) {
			t.Errorf("configuration failed. expected %+v, actual %+v", expected, item)
		}
	}
}

func fromFile(source []byte, item, expected interface{}, isError bool) func(*testing.T) {
	return func(t *testing.T) {

		// Name the temporary file
		name := path.Join(os.TempDir(), "config-json-fromFile.json")

		// Remove the file when test completes
		defer func() {
			if _, err := os.Stat(name); err != nil {
				if os.IsNotExist(err) {
					return
				}
				t.Errorf("unexpected error removing file: %v", err)
				return
			}
			if err := os.Remove(name); err != nil {
				t.Errorf("unexpected error removing file: %v", err)
			}
		}()

		// Write the data to the file
		if err := ioutil.WriteFile(name, source, os.ModePerm); err != nil {
			t.Errorf("unexpected error writing file: %v", err)
			t.FailNow()
		}

		// Create the configurer
		cfg, err := NewFromFile(name)
		if err != nil {
			t.Errorf("error not expected creating configurer. actual was %v", err)
			t.FailNow()
		}

		// Configure the item
		err = cfg.Configure(item)
		t.Run("Error", testError(isError, err))
		if err != nil {
			return
		}

		// Compare the results
		if !reflect.DeepEqual(item, expected) {
			t.Errorf("configuration failed. expected %+v, actual %+v", expected, item)
		}

	}
}

func testError(isError bool, err error) func(*testing.T) {
	return func(t *testing.T) {

		// No error was expected but one was found
		if !isError && err != nil {
			t.Errorf("error not expected. actual: %v", err)
		}

		// Error was expected but none found
		if isError && err == nil {
			t.Error("error expected but none returned")
		}
	}
}

type MockReader struct {
	HasError bool
}

func (r MockReader) Read(p []byte) (n int, err error) {
	if r.HasError {
		err = errors.New("mock reader error")
		return
	}
	b := bytes.NewBufferString(`{"name": "Hemmingway"}`)
	return b.Read(p)
}
