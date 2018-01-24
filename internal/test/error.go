package test

import "testing"

func Error(hasError bool, err error) func(*testing.T) {
	return func(t *testing.T) {

		// An error was expected
		if hasError && err == nil {
			t.Errorf("error was expected but none was returned")
			t.FailNow()
		}

		// No error was expected
		if !hasError && err != nil {
			t.Errorf("no error was expected. actual: %v", err)
			t.FailNow()
		}
	}
}
