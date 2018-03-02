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

func TestMidpoint(t *testing.T) {
	var cases = []struct {
		Desc      string
		Positions []Position
		Expected  Position
	}{
		{
			Desc: "one position",
			Positions: []Position{
				{Layer: 1.0, X: 1.0, Y: 1.0, Z: 1.0},
			},
			Expected: Position{Layer: 1.0, X: 1.0, Y: 1.0, Z: 1.0},
		},
		{
			Desc: "two positions",
			Positions: []Position{
				{Layer: 1.0, X: 3.0, Y: 0.0, Z: 2.0},
				{Layer: 3.0, X: 1.0, Y: 4.0, Z: 2.0},
			},
			Expected: Position{Layer: 2.0, X: 2.0, Y: 2.0, Z: 2.0},
		},
		{
			Desc: "three positions",
			Positions: []Position{
				{Layer: 1.0, X: 1.0, Y: 8.0, Z: 3.0},
				{Layer: 3.0, X: 8.0, Y: 3.0, Z: 1.0},
				{Layer: 8.0, X: 3.0, Y: 1.0, Z: 8.0},
			},
			Expected: Position{Layer: 4.0, X: 4.0, Y: 4.0, Z: 4.0},
		},
	}

	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {
			actual := Midpoint(c.Positions...)
			if c.Expected.Layer != actual.Layer {
				t.Errorf("incorrect layer: expected %f, actual %f", c.Expected.Layer, actual.Layer)
			}
			if c.Expected.X != actual.X {
				t.Errorf("incorrect x: expected %f, actual %f", c.Expected.X, actual.X)
			}
			if c.Expected.Y != actual.Y {
				t.Errorf("incorrect y: expected %f, actual %f", c.Expected.Y, actual.Y)
			}
			if c.Expected.Z != actual.Z {
				t.Errorf("incorrect z: expected %f, actual %f", c.Expected.Z, actual.Z)
			}
		})
	}
}
