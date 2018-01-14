package evo

import "testing"

func TestSubstrateString(t *testing.T) {
	s := Substrate{} // Given even an empty position
	if s.String() == "" {
		t.Error("String() should return non-empty string")
	}
}

func TestSubstrateComplexity(t *testing.T) {

	var cases = []struct {
		Desc      string
		Substrate Substrate
		Expected  int
	}{
		{
			Desc:      "empty substrate",
			Substrate: Substrate{},
			Expected:  0,
		},
		{
			Desc: "only nodes",
			Substrate: Substrate{
				Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}},
			},
			Expected: 2,
		},
		{
			Desc: "only connections",
			Substrate: Substrate{
				Conns: []Conn{{Enabled: true}, {Enabled: true}, {Enabled: true}},
			},
			Expected: 3,
		},
		{
			Desc: "enabled state should not matter",
			Substrate: Substrate{
				Conns: []Conn{{Enabled: true}, {Enabled: false}, {Enabled: true}},
			},
			Expected: 3,
		},
		{
			Desc: "both nodes and connections",
			Substrate: Substrate{
				Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}},
				Conns: []Conn{{Enabled: true}, {Enabled: true}, {Enabled: false}},
			},
			Expected: 5,
		},
	}

	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {
			actual := c.Substrate.Complexity()
			if c.Expected != actual {
				t.Errorf("incorrect complexity value: expected %d, actual %d", c.Expected, actual)
			}
		})
	}
}
