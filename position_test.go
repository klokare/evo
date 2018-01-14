package evo

import "testing"

func TestPositionString(t *testing.T) {
	p := Position{} // Given even an empty position
	if p.String() == "" {
		t.Error("String() should return non-empty string")
	}
}

func TestPositionCompare(t *testing.T) {

	var cases = []struct {
		Desc     string
		A, B     Position
		Expected int8
	}{
		{
			Desc:     "equal positions",
			A:        Position{Layer: 1.0, X: 0.5, Y: -0.5, Z: 0.25},
			B:        Position{Layer: 1.0, X: 0.5, Y: -0.5, Z: 0.25},
			Expected: 0,
		},
		{
			Desc:     "lower layer",
			A:        Position{Layer: 0.0, X: 0.5, Y: -0.5, Z: 0.25},
			B:        Position{Layer: 1.0, X: 0.5, Y: -0.5, Z: 0.25},
			Expected: -1,
		},
		{
			Desc:     "higher layer",
			A:        Position{Layer: 1.0, X: 0.5, Y: -0.5, Z: 0.25},
			B:        Position{Layer: 0.0, X: 0.5, Y: -0.5, Z: 0.25},
			Expected: 1,
		},
		{
			Desc:     "lower X",
			A:        Position{Layer: 1.0, X: 0.0, Y: -0.5, Z: 0.25},
			B:        Position{Layer: 1.0, X: 0.5, Y: -0.5, Z: 0.25},
			Expected: -1,
		},
		{
			Desc:     "higher X",
			A:        Position{Layer: 1.0, X: 1.0, Y: -0.5, Z: 0.25},
			B:        Position{Layer: 1.0, X: 0.5, Y: -0.5, Z: 0.25},
			Expected: 1,
		},
		{
			Desc:     "lower Y",
			A:        Position{Layer: 1.0, X: 0.5, Y: -0.8, Z: 0.25},
			B:        Position{Layer: 1.0, X: 0.5, Y: -0.5, Z: 0.25},
			Expected: -1,
		},
		{
			Desc:     "higher Y",
			A:        Position{Layer: 1.0, X: 0.5, Y: -0.2, Z: 0.25},
			B:        Position{Layer: 1.0, X: 0.5, Y: -0.5, Z: 0.25},
			Expected: 1,
		},
		{
			Desc:     "lower Z",
			A:        Position{Layer: 1.0, X: 0.5, Y: -0.5, Z: 0.15},
			B:        Position{Layer: 1.0, X: 0.5, Y: -0.5, Z: 0.25},
			Expected: -1,
		},
		{
			Desc:     "higher Z",
			A:        Position{Layer: 1.0, X: 0.5, Y: -0.5, Z: 0.35},
			B:        Position{Layer: 1.0, X: 0.5, Y: -0.5, Z: 0.25},
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
