package neat

import (
	"testing"

	"github.com/klokare/evo"
	"github.com/klokare/evo/internal/mock"
)

func TestComplexityDistance(t *testing.T) {

	var tests = []struct {
		Desc       string
		N, C, F, W float64    // Coefficients: Disjoint, Excess, Activation, and Weight
		A, B       evo.Genome // Genomes to compare
		Expected   float64    // Expected distance
		IsError    bool       // Expect an error
	}{
		{ // Empty genomes. Not supposed to happen in NEAT but it is a vaild input for this helper
			N: 0.0, C: 0.0, W: 0.0,
			A: evo.Genome{}, B: evo.Genome{},
			Expected: 0.0, IsError: false,
		},
		{
			Desc: "just nodes, same activations",
			N:    1.0, C: 1.0, W: 1.0,
			A: evo.Genome{
				Encoded: evo.Substrate{
					Nodes: []evo.Node{
						{Position: evo.Position{Layer: 0.0, X: 0.0}, Neuron: evo.Input, Activation: evo.Direct},
						{Position: evo.Position{Layer: 0.0, X: 0.5}, Neuron: evo.Input, Activation: evo.Direct},
						{Position: evo.Position{Layer: 0.0, X: 1.0}, Neuron: evo.Input, Activation: evo.Direct},
						{Position: evo.Position{Layer: 1.0, X: 0.5}, Neuron: evo.Output, Activation: evo.Sigmoid},
					},
				},
			},
			B: evo.Genome{
				Encoded: evo.Substrate{
					Nodes: []evo.Node{
						{Position: evo.Position{Layer: 0.0, X: 0.0}, Neuron: evo.Input, Activation: evo.Direct},
						{Position: evo.Position{Layer: 0.0, X: 0.5}, Neuron: evo.Input, Activation: evo.Direct},
						{Position: evo.Position{Layer: 0.0, X: 1.0}, Neuron: evo.Input, Activation: evo.Direct},
						{Position: evo.Position{Layer: 1.0, X: 0.5}, Neuron: evo.Output, Activation: evo.Sigmoid},
					},
				},
			},
			Expected: 0.0, IsError: false,
		},
		{
			Desc: "disjoint node in genome A",
			N:    2.0, C: 1.0, W: 1.0,
			A: evo.Genome{
				Encoded: evo.Substrate{
					Nodes: []evo.Node{
						{Position: evo.Position{Layer: 0.0, X: 0.0}, Neuron: evo.Input, Activation: evo.Direct},
						{Position: evo.Position{Layer: 0.0, X: 0.5}, Neuron: evo.Input, Activation: evo.Direct},
						{Position: evo.Position{Layer: 0.0, X: 1.0}, Neuron: evo.Input, Activation: evo.Direct},
						{Position: evo.Position{Layer: 1.0, X: 0.5}, Neuron: evo.Output, Activation: evo.Sigmoid},
					},
				},
			},
			B: evo.Genome{
				Encoded: evo.Substrate{
					Nodes: []evo.Node{
						{Position: evo.Position{Layer: 0.0, X: 0.0}, Neuron: evo.Input, Activation: evo.Direct},
						{Position: evo.Position{Layer: 0.0, X: 1.0}, Neuron: evo.Input, Activation: evo.Direct},
						{Position: evo.Position{Layer: 1.0, X: 0.5}, Neuron: evo.Output, Activation: evo.Sigmoid},
					},
				},
			},
			Expected: 2.0, IsError: false,
		},
		{
			Desc: "disjoint node in genome B",
			N:    2.0, C: 1.0, W: 1.0,
			A: evo.Genome{
				Encoded: evo.Substrate{
					Nodes: []evo.Node{
						{Position: evo.Position{Layer: 0.0, X: 0.0}, Neuron: evo.Input, Activation: evo.Direct},
						{Position: evo.Position{Layer: 0.0, X: 1.0}, Neuron: evo.Input, Activation: evo.Direct},
						{Position: evo.Position{Layer: 1.0, X: 0.5}, Neuron: evo.Output, Activation: evo.Sigmoid},
					},
				},
			},
			B: evo.Genome{
				Encoded: evo.Substrate{
					Nodes: []evo.Node{
						{Position: evo.Position{Layer: 0.0, X: 0.0}, Neuron: evo.Input, Activation: evo.Direct},
						{Position: evo.Position{Layer: 0.0, X: 0.5}, Neuron: evo.Input, Activation: evo.Direct},
						{Position: evo.Position{Layer: 0.0, X: 1.0}, Neuron: evo.Input, Activation: evo.Direct},
						{Position: evo.Position{Layer: 1.0, X: 0.5}, Neuron: evo.Output, Activation: evo.Sigmoid},
					},
				},
			},
			Expected: 2.0, IsError: false,
		},
		{
			Desc: "excess node in genome A",
			N:    2.0, C: 1.0, W: 1.0,
			A: evo.Genome{
				Encoded: evo.Substrate{
					Nodes: []evo.Node{
						{Position: evo.Position{Layer: 0.0, X: 0.0}, Neuron: evo.Input, Activation: evo.Direct},
						{Position: evo.Position{Layer: 0.0, X: 1.0}, Neuron: evo.Input, Activation: evo.Direct},
						{Position: evo.Position{Layer: 1.0, X: 0.5}, Neuron: evo.Output, Activation: evo.Sigmoid},
						{Position: evo.Position{Layer: 1.0, X: 1.0}, Neuron: evo.Output, Activation: evo.Sigmoid},
					},
				},
			},
			B: evo.Genome{
				Encoded: evo.Substrate{
					Nodes: []evo.Node{
						{Position: evo.Position{Layer: 0.0, X: 0.0}, Neuron: evo.Input, Activation: evo.Direct},
						{Position: evo.Position{Layer: 0.0, X: 1.0}, Neuron: evo.Input, Activation: evo.Direct},
						{Position: evo.Position{Layer: 1.0, X: 0.5}, Neuron: evo.Output, Activation: evo.Sigmoid},
					},
				},
			},
			Expected: 2.0, IsError: false,
		},
		{
			Desc: "Excess node in genome b",
			N:    2.0, C: 1.0, W: 1.0,
			A: evo.Genome{
				Encoded: evo.Substrate{
					Nodes: []evo.Node{
						{Position: evo.Position{Layer: 0.0, X: 0.0}, Neuron: evo.Input, Activation: evo.Direct},
						{Position: evo.Position{Layer: 0.0, X: 1.0}, Neuron: evo.Input, Activation: evo.Direct},
						{Position: evo.Position{Layer: 1.0, X: 0.5}, Neuron: evo.Output, Activation: evo.Sigmoid},
					},
				},
			},
			B: evo.Genome{
				Encoded: evo.Substrate{
					Nodes: []evo.Node{
						{Position: evo.Position{Layer: 0.0, X: 0.0}, Neuron: evo.Input, Activation: evo.Direct},
						{Position: evo.Position{Layer: 0.0, X: 1.0}, Neuron: evo.Input, Activation: evo.Direct},
						{Position: evo.Position{Layer: 1.0, X: 0.5}, Neuron: evo.Output, Activation: evo.Sigmoid},
						{Position: evo.Position{Layer: 1.0, X: 1.0}, Neuron: evo.Output, Activation: evo.Sigmoid},
					},
				},
			},
			Expected: 2.0, IsError: false,
		},
		{
			Desc: "same nodes, different activations",
			N:    1.0, C: 1.0, W: 1.0, F: 2.0,
			A: evo.Genome{
				Encoded: evo.Substrate{
					Nodes: []evo.Node{
						{Position: evo.Position{Layer: 0.0, X: 0.0}, Neuron: evo.Input, Activation: evo.Direct},
						{Position: evo.Position{Layer: 0.0, X: 0.5}, Neuron: evo.Input, Activation: evo.Direct},
						{Position: evo.Position{Layer: 0.0, X: 1.0}, Neuron: evo.Input, Activation: evo.Direct},
						{Position: evo.Position{Layer: 0.5, X: 0.5}, Neuron: evo.Hidden, Activation: evo.Tanh},
						{Position: evo.Position{Layer: 1.0, X: 0.5}, Neuron: evo.Output, Activation: evo.Sigmoid},
					},
				},
			},
			B: evo.Genome{
				Encoded: evo.Substrate{
					Nodes: []evo.Node{
						{Position: evo.Position{Layer: 0.0, X: 0.0}, Neuron: evo.Input, Activation: evo.Direct},
						{Position: evo.Position{Layer: 0.0, X: 0.5}, Neuron: evo.Input, Activation: evo.Direct},
						{Position: evo.Position{Layer: 0.0, X: 1.0}, Neuron: evo.Input, Activation: evo.Direct},
						{Position: evo.Position{Layer: 0.5, X: 0.5}, Neuron: evo.Hidden, Activation: evo.InverseAbs},
						{Position: evo.Position{Layer: 1.0, X: 0.5}, Neuron: evo.Output, Activation: evo.Sigmoid},
					},
				},
			},
			Expected: 0.4, IsError: false,
		},
		{
			Desc: "equal genomes",
			N:    1.0, C: 1.0, W: 1.0,
			A: evo.Genome{
				Encoded: evo.Substrate{
					Nodes: []evo.Node{
						{Position: evo.Position{Layer: 0.0, X: 0.0}, Neuron: evo.Input, Activation: evo.Direct},
						{Position: evo.Position{Layer: 0.0, X: 0.5}, Neuron: evo.Input, Activation: evo.Direct},
						{Position: evo.Position{Layer: 0.0, X: 1.0}, Neuron: evo.Input, Activation: evo.Direct},
						{Position: evo.Position{Layer: 1.0, X: 0.5}, Neuron: evo.Output, Activation: evo.Sigmoid},
					},
					Conns: []evo.Conn{
						{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 1.1, Enabled: true},
						{Source: evo.Position{Layer: 0.0, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 2.2, Enabled: true},
						{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 3.3, Enabled: true},
					},
				},
			},
			B: evo.Genome{
				Encoded: evo.Substrate{
					Nodes: []evo.Node{
						{Position: evo.Position{Layer: 0.0, X: 0.0}, Neuron: evo.Input, Activation: evo.Direct},
						{Position: evo.Position{Layer: 0.0, X: 0.5}, Neuron: evo.Input, Activation: evo.Direct},
						{Position: evo.Position{Layer: 0.0, X: 1.0}, Neuron: evo.Input, Activation: evo.Direct},
						{Position: evo.Position{Layer: 1.0, X: 0.5}, Neuron: evo.Output, Activation: evo.Sigmoid},
					},
					Conns: []evo.Conn{
						{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 1.1, Enabled: true},
						{Source: evo.Position{Layer: 0.0, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 2.2, Enabled: true},
						{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 3.3, Enabled: true},
					},
				},
			},
			Expected: 0.0, IsError: false,
		},
		{
			Desc: "disabled connections do not affect calculation",
			N:    1.0, C: 1.0, W: 1.0,
			A: evo.Genome{
				Encoded: evo.Substrate{
					Conns: []evo.Conn{
						{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 1.1, Enabled: false},
						{Source: evo.Position{Layer: 0.0, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 2.2, Enabled: false},
						{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 3.3, Enabled: true},
					},
				},
			},
			B: evo.Genome{
				Encoded: evo.Substrate{
					Conns: []evo.Conn{
						{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 1.1, Enabled: true},
						{Source: evo.Position{Layer: 0.0, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 2.2, Enabled: false},
						{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 3.3, Enabled: false},
					},
				},
			},
			Expected: 0.0, IsError: false,
		},
		{
			Desc: "excess genome connection in genome A",
			N:    1.0, C: 2.0, W: 1.0,
			A: evo.Genome{
				Encoded: evo.Substrate{
					Conns: []evo.Conn{
						{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 1.1, Enabled: true},
						{Source: evo.Position{Layer: 0.0, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 2.2, Enabled: true},
						{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 3.3, Enabled: true},
						{Source: evo.Position{Layer: 1.0, X: 0.5}, Target: evo.Position{Layer: 0.0, X: 1.0}, Weight: 4.4, Enabled: true},
					},
				},
			},
			B: evo.Genome{
				Encoded: evo.Substrate{
					Conns: []evo.Conn{
						{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 1.1, Enabled: true},
						{Source: evo.Position{Layer: 0.0, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 2.2, Enabled: true},
						{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 3.3, Enabled: true},
					},
				},
			},
			Expected: 2.0, IsError: false,
		},
		{
			Desc: "excess genome connection in genome B",
			N:    1.0, C: 2.0, W: 1.0,
			A: evo.Genome{
				Encoded: evo.Substrate{
					Conns: []evo.Conn{
						{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 1.1, Enabled: true},
						{Source: evo.Position{Layer: 0.0, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 2.2, Enabled: true},
						{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 3.3, Enabled: true},
					},
				},
			},
			B: evo.Genome{
				Encoded: evo.Substrate{
					Conns: []evo.Conn{
						{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 1.1, Enabled: true},
						{Source: evo.Position{Layer: 0.0, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 2.2, Enabled: true},
						{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 3.3, Enabled: true},
						{Source: evo.Position{Layer: 1.0, X: 0.5}, Target: evo.Position{Layer: 0.0, X: 1.0}, Weight: 4.4, Enabled: true},
					},
				},
			},
			Expected: 2.0, IsError: false,
		},
		{
			Desc: "disjoint connection in genome A",
			N:    2.0, C: 2.0, W: 1.0,
			A: evo.Genome{
				Encoded: evo.Substrate{
					Conns: []evo.Conn{
						{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 1.1, Enabled: true},
						{Source: evo.Position{Layer: 0.0, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 2.2, Enabled: true},
						{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 3.3, Enabled: true},
					},
				},
			},
			B: evo.Genome{
				Encoded: evo.Substrate{
					Conns: []evo.Conn{
						{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 1.1, Enabled: true},
						{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 3.3, Enabled: true},
					},
				},
			},
			Expected: 2.0, IsError: false,
		},
		{
			Desc: "disjoint connection in genome B",
			N:    2.0, C: 2.0, W: 1.0,
			A: evo.Genome{
				Encoded: evo.Substrate{
					Conns: []evo.Conn{
						{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 1.1, Enabled: true},
						{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 3.3, Enabled: true},
					},
				},
			},
			B: evo.Genome{
				Encoded: evo.Substrate{
					Conns: []evo.Conn{
						{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 1.1, Enabled: true},
						{Source: evo.Position{Layer: 0.0, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 2.2, Enabled: true},
						{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 3.3, Enabled: true},
					},
				},
			},
			Expected: 2.0, IsError: false,
		},
		{
			Desc: "weight differences",
			N:    1.0, C: 1.0, W: 3.0,
			A: evo.Genome{
				Encoded: evo.Substrate{
					Conns: []evo.Conn{
						{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 1.1, Enabled: true}, // Different
						{Source: evo.Position{Layer: 0.0, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 4.0, Enabled: true}, // Different
						{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 3.3, Enabled: true}, // Same
					},
				},
			},
			B: evo.Genome{
				Encoded: evo.Substrate{
					Conns: []evo.Conn{
						{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 0.5, Enabled: true},  // Different
						{Source: evo.Position{Layer: 0.0, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: -2.0, Enabled: true}, // Different
						{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 3.3, Enabled: true},  // Same
					},
				},
			},
			Expected: 6.6, IsError: false,
		},
	}

	for _, c := range tests {
		t.Run(c.Desc, func(t *testing.T) {

			// Compare the Genomes
			cmp := &Compatibility{
				NodesCoefficient:      c.N,
				ActivationCoefficient: c.F,
				ConnsCoefficient:      c.C,
				WeightCoefficient:     c.W,
			}

			dist, err := cmp.Distance(c.A, c.B)

			// Check if error was expected
			if t.Run("Error", mock.Error(c.IsError, err)) == false {
				t.Fail()
			}
			if err != nil {
				return // Error was received so the rest of the tests don't make sense
			}

			// Check the actual distance
			if dist != c.Expected {
				t.Errorf("incorrect compatibility distance. expected %f, actual %f", c.Expected, dist)
			}
		})
	}
}
