package evo

import (
	"strings"
	"testing"
)

func TestConnString(t *testing.T) {
	c := Conn{} // Given even an empty position
	if c.String() == "" {
		t.Error("String() should return non-empty string")
	}

	c.Enabled = true
	s := c.String()
	if !strings.Contains(strings.ToLower(s), "enabled") {
		t.Errorf("string for enabled connection does not indicate enabled: actual, %s", s)
	}

	c.Enabled = false
	s = c.String()
	if !strings.Contains(strings.ToLower(s), "disabled") {
		t.Errorf("string for disabled connection does not indicate disabled: actual, %s", s)
	}
}

func TestConnCompare(t *testing.T) {

	var cases = []struct {
		Desc     string
		A, B     Conn
		Expected int8
	}{
		{
			Desc:     "equal conns",
			A:        Conn{Source: Position{Layer: 0.0, X: 0.0}, Target: Position{Layer: 1.0, X: 1.0}, Enabled: true, Weight: 1.234},
			B:        Conn{Source: Position{Layer: 0.0, X: 0.0}, Target: Position{Layer: 1.0, X: 1.0}, Enabled: true, Weight: 1.234},
			Expected: 0,
		},
		{
			Desc:     "different enabled states should not matter",
			A:        Conn{Source: Position{Layer: 0.0, X: 0.0}, Target: Position{Layer: 1.0, X: 1.0}, Enabled: true, Weight: 1.234},
			B:        Conn{Source: Position{Layer: 0.0, X: 0.0}, Target: Position{Layer: 1.0, X: 1.0}, Enabled: false, Weight: 1.234},
			Expected: 0,
		},
		{
			Desc:     "different weights should not matter",
			A:        Conn{Source: Position{Layer: 0.0, X: 0.0}, Target: Position{Layer: 1.0, X: 1.0}, Enabled: true, Weight: 1.234},
			B:        Conn{Source: Position{Layer: 0.0, X: 0.0}, Target: Position{Layer: 1.0, X: 1.0}, Enabled: true, Weight: 2.112},
			Expected: 0,
		},
		{
			Desc:     "lower source, equal taget",
			A:        Conn{Source: Position{Layer: 0.0, X: 0.0}, Target: Position{Layer: 1.0, X: 1.0}, Enabled: true, Weight: 1.234},
			B:        Conn{Source: Position{Layer: 1.0, X: 0.0}, Target: Position{Layer: 1.0, X: 1.0}, Enabled: true, Weight: 1.234},
			Expected: -1,
		},
		{
			Desc:     "higher source, equal target",
			A:        Conn{Source: Position{Layer: 1.0, X: 0.0}, Target: Position{Layer: 1.0, X: 1.0}, Enabled: true, Weight: 1.234},
			B:        Conn{Source: Position{Layer: 0.0, X: 0.0}, Target: Position{Layer: 1.0, X: 1.0}, Enabled: true, Weight: 1.234},
			Expected: 1,
		},
		{
			Desc:     "lower source, lower taget",
			A:        Conn{Source: Position{Layer: 0.0, X: 0.0}, Target: Position{Layer: 0.0, X: 1.0}, Enabled: true, Weight: 1.234},
			B:        Conn{Source: Position{Layer: 1.0, X: 0.0}, Target: Position{Layer: 1.0, X: 1.0}, Enabled: true, Weight: 1.234},
			Expected: -1,
		},
		{
			Desc:     "higher source, lower target",
			A:        Conn{Source: Position{Layer: 1.0, X: 0.0}, Target: Position{Layer: 0.0, X: 1.0}, Enabled: true, Weight: 1.234},
			B:        Conn{Source: Position{Layer: 0.0, X: 0.0}, Target: Position{Layer: 1.0, X: 1.0}, Enabled: true, Weight: 1.234},
			Expected: 1,
		},
		{
			Desc:     "lower source, higher taget",
			A:        Conn{Source: Position{Layer: 0.0, X: 0.0}, Target: Position{Layer: 1.0, X: 1.0}, Enabled: true, Weight: 1.234},
			B:        Conn{Source: Position{Layer: 1.0, X: 0.0}, Target: Position{Layer: 0.0, X: 1.0}, Enabled: true, Weight: 1.234},
			Expected: -1,
		},
		{
			Desc:     "higher source, higher target",
			A:        Conn{Source: Position{Layer: 1.0, X: 0.0}, Target: Position{Layer: 1.0, X: 1.0}, Enabled: true, Weight: 1.234},
			B:        Conn{Source: Position{Layer: 0.0, X: 0.0}, Target: Position{Layer: 0.0, X: 1.0}, Enabled: true, Weight: 1.234},
			Expected: 1,
		},
		{
			Desc:     "equal source, lower taget",
			A:        Conn{Source: Position{Layer: 0.0, X: 0.0}, Target: Position{Layer: 0.0, X: 1.0}, Enabled: true, Weight: 1.234},
			B:        Conn{Source: Position{Layer: 0.0, X: 0.0}, Target: Position{Layer: 1.0, X: 1.0}, Enabled: true, Weight: 1.234},
			Expected: -1,
		},
		{
			Desc:     "equal source, higher target",
			A:        Conn{Source: Position{Layer: 0.0, X: 0.0}, Target: Position{Layer: 1.0, X: 1.0}, Enabled: true, Weight: 1.234},
			B:        Conn{Source: Position{Layer: 0.0, X: 0.0}, Target: Position{Layer: 0.0, X: 1.0}, Enabled: true, Weight: 1.234},
			Expected: 1,
		},
	}

	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {
			actual := c.A.Compare(c.B)
			if c.Expected != actual {
				t.Errorf("incorrect comparison value: expected %d, actual %d", c.Expected, actual)
			}
		})
	}
}
