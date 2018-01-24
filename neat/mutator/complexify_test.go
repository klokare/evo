package mutator

import (
	"context"
	"testing"

	"github.com/klokare/evo"
	"github.com/klokare/evo/internal/mock"
)

func TestWithComplexify(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		e := &evo.Experiment{Mutators: make([]evo.Mutator, 0, 1)}
		cfg := &mock.Configurer{HasError: true}
		err := WithComplexify(cfg)(e)
		if err == nil {
			t.Errorf("error expected")
		}
	})
	t.Run("enabled", func(t *testing.T) {
		e := &evo.Experiment{Mutators: make([]evo.Mutator, 0, 1)}
		cfg := &mock.Configurer{AddConnProbability: 1.0, AddNodeProbability: 1.0}
		err := WithComplexify(cfg)(e)
		if err != nil {
			t.Errorf("error unxpected: %v", err)
		}
		if len(e.Mutators) == 0 {
			t.Errorf("incorrect number of mutators: expected 1, actual 0")
		} else if _, ok := e.Mutators[0].(*Complexify); !ok {
			t.Errorf("mutator is not a complexify")
		}
	})
	t.Run("not enabled", func(t *testing.T) {
		e := &evo.Experiment{Mutators: make([]evo.Mutator, 0, 1)}
		cfg := &mock.Configurer{AddNodeProbability: 0.0, AddConnProbability: 0.0}
		err := WithComplexify(cfg)(e)
		if err != nil {
			t.Errorf("error unxpected: %v", err)
		}
		if len(e.Mutators) > 0 {
			t.Errorf("incorrect number of mutators: expected 0, actual: %d", len(e.Mutators))
		}
	})
}

func TestComplexify(t *testing.T) {

	var tests = []struct {
		AddNodeProbability float64
		AddConnProbability float64
		Original           evo.Genome
		Expected           evo.Genome
	}{
		{ // No probabilities, no change
			AddNodeProbability: 0.0,
			AddConnProbability: 0.0,
			Original: evo.Genome{
				Encoded: evo.Substrate{
					Nodes: []evo.Node{
						{Position: evo.Position{Layer: 0.0, X: 1.0}, Neuron: evo.Input, Activation: evo.Direct},
						{Position: evo.Position{Layer: 1.0, X: 0.5}, Neuron: evo.Output, Activation: evo.Sigmoid},
					},
					Conns: []evo.Conn{
						{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 1.1, Enabled: true},
					},
				},
			},
			Expected: evo.Genome{
				Encoded: evo.Substrate{
					Nodes: []evo.Node{
						{Position: evo.Position{Layer: 0.0, X: 1.0}, Neuron: evo.Input, Activation: evo.Direct},
						{Position: evo.Position{Layer: 1.0, X: 0.5}, Neuron: evo.Output, Activation: evo.Sigmoid},
					},
					Conns: []evo.Conn{
						{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 1.1, Enabled: true},
					},
				},
			},
		},
		//{ // Add node but no possibiliity, would try to create an existing node. Is this a real possibility?
		//},
		{ // Add node with possibility
			AddNodeProbability: 1.0,
			AddConnProbability: 0.0,
			Original: evo.Genome{
				Encoded: evo.Substrate{
					Nodes: []evo.Node{
						{Position: evo.Position{Layer: 0.0, X: 0.0}, Neuron: evo.Input, Activation: evo.Direct},
						{Position: evo.Position{Layer: 1.0, X: 1.0}, Neuron: evo.Output, Activation: evo.Sigmoid},
					},
					Conns: []evo.Conn{
						{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 1.0}, Weight: 2.0, Enabled: true},
					},
				},
			},
			Expected: evo.Genome{
				Encoded: evo.Substrate{
					Nodes: []evo.Node{
						{Position: evo.Position{Layer: 0.0, X: 0.0}, Neuron: evo.Input, Activation: evo.Direct},
						{Position: evo.Position{Layer: 0.5, X: 0.5}, Neuron: evo.Hidden, Activation: evo.SteepenedSigmoid},
						{Position: evo.Position{Layer: 1.0, X: 1.0}, Neuron: evo.Output, Activation: evo.Sigmoid},
					},
					Conns: []evo.Conn{
						{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 1.0}, Weight: 2.0, Enabled: false},
						{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 0.5, X: 0.5}, Weight: 1.0, Enabled: true},
						{Source: evo.Position{Layer: 0.5, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 1.0}, Weight: 2.0, Enabled: true},
					},
				},
			},
		},
		{ // Add conn but all nodes connected (feed-forward only) so no new connection
			AddNodeProbability: 0.0,
			AddConnProbability: 1.0,
			Original: evo.Genome{
				Encoded: evo.Substrate{
					Nodes: []evo.Node{
						{Position: evo.Position{Layer: 0.0, X: 0.0}, Neuron: evo.Input, Activation: evo.Direct},
						{Position: evo.Position{Layer: 0.5, X: 0.5}, Neuron: evo.Hidden, Activation: evo.SteepenedSigmoid},
						{Position: evo.Position{Layer: 1.0, X: 1.0}, Neuron: evo.Output, Activation: evo.Sigmoid},
					},
					Conns: []evo.Conn{
						{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 1.0}, Weight: 2.0, Enabled: true},
						{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 0.5, X: 0.5}, Weight: 3.0, Enabled: true},
						{Source: evo.Position{Layer: 0.5, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 1.0}, Weight: 4.0, Enabled: true},
					},
				},
			},
			Expected: evo.Genome{
				Encoded: evo.Substrate{
					Nodes: []evo.Node{
						{Position: evo.Position{Layer: 0.0, X: 0.0}, Neuron: evo.Input, Activation: evo.Direct},
						{Position: evo.Position{Layer: 0.5, X: 0.5}, Neuron: evo.Hidden, Activation: evo.SteepenedSigmoid},
						{Position: evo.Position{Layer: 1.0, X: 1.0}, Neuron: evo.Output, Activation: evo.Sigmoid},
					},
					Conns: []evo.Conn{
						{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 1.0}, Weight: 2.0, Enabled: true},
						{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 0.5, X: 0.5}, Weight: 3.0, Enabled: true},
						{Source: evo.Position{Layer: 0.5, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 1.0}, Weight: 4.0, Enabled: true},
					},
				},
			},
		},
		{ // Missing connection possible
			AddNodeProbability: 0.0,
			AddConnProbability: 1.0,
			Original: evo.Genome{
				Encoded: evo.Substrate{
					Nodes: []evo.Node{
						{Position: evo.Position{Layer: 0.0, X: 0.0}, Neuron: evo.Input, Activation: evo.Direct},
						{Position: evo.Position{Layer: 0.5, X: 0.5}, Neuron: evo.Hidden, Activation: evo.SteepenedSigmoid},
						{Position: evo.Position{Layer: 1.0, X: 1.0}, Neuron: evo.Output, Activation: evo.Sigmoid},
					},
					Conns: []evo.Conn{
						{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 1.0}, Weight: 2.0, Enabled: true},
						{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 0.5, X: 0.5}, Weight: 3.0, Enabled: true},
						// Missing hidden to output
					},
				},
			},
			Expected: evo.Genome{
				Encoded: evo.Substrate{
					Nodes: []evo.Node{
						{Position: evo.Position{Layer: 0.0, X: 0.0}, Neuron: evo.Input, Activation: evo.Direct},
						{Position: evo.Position{Layer: 0.5, X: 0.5}, Neuron: evo.Hidden, Activation: evo.SteepenedSigmoid},
						{Position: evo.Position{Layer: 1.0, X: 1.0}, Neuron: evo.Output, Activation: evo.Sigmoid},
					},
					Conns: []evo.Conn{
						{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 1.0}, Weight: 2.0, Enabled: true},
						{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 0.5, X: 0.5}, Weight: 3.0, Enabled: true},
						{Source: evo.Position{Layer: 0.5, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 1.0}, Weight: 4.0, Enabled: true},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run("Mutate", func(t *testing.T) {

			// Create a mutator
			mut := &Complexify{
				AddNodeProbability: test.AddNodeProbability,
				AddConnProbability: test.AddConnProbability,
				HiddenActivation:   evo.SteepenedSigmoid,
			}

			// Mutate the genome
			err := mut.Mutate(context.Background(), &test.Original)
			t.Logf("After: %+v\n", test.Original)

			// There should be no error
			if err != nil {
				t.Errorf("there should be no error. instead, %v", err)
				t.FailNow()
			}

			// The expected structure should be found
			t.Run("Genome", testComplexifyGenome(test.Original, test.Expected, test.AddNodeProbability == 1.0 && test.AddConnProbability == 0.0))

		})
	}
}

func testComplexifyGenome(actual, expected evo.Genome, checkWeight bool) func(*testing.T) {
	return func(t *testing.T) {

		// Ensure the nodes
		t.Run("Nodes", testComplexifyNodes(actual.Encoded.Nodes, expected.Encoded.Nodes))

		// Ensure the conns
		t.Run("Conns", testComplexifyConns(actual.Encoded.Conns, expected.Encoded.Conns, checkWeight))
	}
}

func testComplexifyNodes(actual, expected []evo.Node) func(*testing.T) {
	return func(t *testing.T) {

		// There should be the same number of nodes
		if len(actual) != len(expected) {
			t.Errorf("incorrect number of nodes. expected %d, actual %d", len(expected), len(actual))
			t.FailNow()
		}

		// Nodes should match
		for _, en := range expected {
			found := false
			for _, an := range actual {
				if an.Compare(en) == 0 {
					// Properties should match
					if an.Neuron != en.Neuron {
						t.Errorf("incorrect neuron type. expected %v, actual %v", en.Neuron, an.Neuron)
					}
					if an.Activation != en.Activation {
						t.Errorf("incorrect activation type. expected %v, actual %v", en.Activation, an.Activation)
					}

					// Note match and move to next node
					found = true
					break
				}
			}
			if !found {
				t.Errorf("node not found. expected node %v", en.Position)
			}
		}
	}
}

func testComplexifyConns(actual, expected []evo.Conn, checkWeight bool) func(*testing.T) {
	return func(t *testing.T) {

		// There should be the same number of connections
		if len(actual) != len(expected) {
			t.Errorf("incorrect number of connections. expected %d, actual %d", len(expected), len(actual))
			t.FailNow()
		}

		// Nodes should match
		for _, ec := range expected {
			found := false
			for _, ac := range actual {
				if ac.Compare(ec) == 0 {
					// Properties should match
					if ac.Source != ec.Source {
						t.Errorf("incorrect source. expected %v, actual %v", ec.Source, ac.Source)
					}
					if ac.Target != ec.Target {
						t.Errorf("incorrect target. expected %v, actual %v", ec.Target, ac.Target)
					}
					if checkWeight && ac.Weight != ec.Weight {
						t.Errorf("incorrect weight. expected %f, actual %f", ec.Weight, ac.Weight)
					}
					if ac.Enabled != ec.Enabled {
						t.Errorf("incorrect enabled. expected %v, actual %v", ec.Enabled, ac.Enabled)
					}

					// Note match and move to next connection
					found = true
					break
				}
			}
			if !found {
				t.Errorf("connection not found. expected conn ID %v -> %v", ec.Source, ec.Target)
			}
		}
	}
}

// Tests to ensure that nodes (add conn) and conns (add node) are chosen randomly as well as
// weights are set randomly during add connection.
func TestComplexifyRandom(t *testing.T) {

}
