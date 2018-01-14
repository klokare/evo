package evo

import "testing"

func TestGenomeComplexity(t *testing.T) {

	var cases = []struct {
		Desc     string
		Genome   Genome
		Expected int
	}{
		{
			Desc:     "empty genome",
			Genome:   Genome{},
			Expected: 0,
		},
		{
			Desc: "should come from encoded substrate",
			Genome: Genome{
				Encoded: Substrate{
					Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}},
					Conns: []Conn{{Enabled: true}, {Enabled: true}, {Enabled: false}},
				},
			},
			Expected: 5,
		},
		{
			Desc: "should not come from decoded substrate",
			Genome: Genome{
				Decoded: Substrate{
					Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}},
					Conns: []Conn{{Enabled: true}, {Enabled: true}, {Enabled: false}},
				},
			},
			Expected: 0,
		},
	}

	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {
			actual := c.Genome.Complexity()
			if c.Expected != actual {
				t.Errorf("incorrect complexity value: expected %d, actual %d", c.Expected, actual)
			}
		})
	}
}
