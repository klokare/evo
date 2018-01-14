package evo

import "testing"

func TestNodeString(t *testing.T) {
	n := Node{} // Given even an empty position
	if n.String() == "" {
		t.Error("String() should return non-empty string")
	}
}

func TestNodeCompare(t *testing.T) {

	var cases = []struct {
		Desc     string
		A, B     Node
		Expected int8
	}{
		{
			Desc:     "equal nodes",
			A:        Node{Position: Position{Layer: 0.5, X: 0.5, Y: 0.5, Z: 0.5}, Neuron: Hidden, Activation: Sigmoid, Bias: 1.234},
			B:        Node{Position: Position{Layer: 0.5, X: 0.5, Y: 0.5, Z: 0.5}, Neuron: Hidden, Activation: Sigmoid, Bias: 1.234},
			Expected: 0,
		},
		{
			Desc:     "different neuron types should not matter",
			A:        Node{Position: Position{Layer: 0.5, X: 0.5, Y: 0.5, Z: 0.5}, Neuron: Input, Activation: Sigmoid, Bias: 1.234},
			B:        Node{Position: Position{Layer: 0.5, X: 0.5, Y: 0.5, Z: 0.5}, Neuron: Hidden, Activation: Sigmoid, Bias: 1.234},
			Expected: 0,
		},
		{
			Desc:     "different activation types should not matter",
			A:        Node{Position: Position{Layer: 0.5, X: 0.5, Y: 0.5, Z: 0.5}, Neuron: Hidden, Activation: Sigmoid, Bias: 1.234},
			B:        Node{Position: Position{Layer: 0.5, X: 0.5, Y: 0.5, Z: 0.5}, Neuron: Hidden, Activation: Direct, Bias: 1.234},
			Expected: 0,
		},
		{
			Desc:     "different bias values should not matter",
			A:        Node{Position: Position{Layer: 0.5, X: 0.5, Y: 0.5, Z: 0.5}, Neuron: Hidden, Activation: Sigmoid, Bias: 3.212},
			B:        Node{Position: Position{Layer: 0.5, X: 0.5, Y: 0.5, Z: 0.5}, Neuron: Hidden, Activation: Sigmoid, Bias: 1.234},
			Expected: 0,
		},
		{
			Desc:     "lower position",
			A:        Node{Position: Position{Layer: 0.0, X: 0.5, Y: 0.5, Z: 0.5}, Neuron: Hidden, Activation: Sigmoid, Bias: 1.234},
			B:        Node{Position: Position{Layer: 0.5, X: 0.5, Y: 0.5, Z: 0.5}, Neuron: Hidden, Activation: Sigmoid, Bias: 1.234},
			Expected: -1,
		},
		{
			Desc:     "higher position",
			A:        Node{Position: Position{Layer: 1.0, X: 0.5, Y: 0.5, Z: 0.5}, Neuron: Hidden, Activation: Sigmoid, Bias: 1.234},
			B:        Node{Position: Position{Layer: 0.5, X: 0.5, Y: 0.5, Z: 0.5}, Neuron: Hidden, Activation: Sigmoid, Bias: 1.234},
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
