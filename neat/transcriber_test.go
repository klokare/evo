package neat

import (
	"context"
	"testing"

	"github.com/klokare/evo"
	"github.com/klokare/evo/internal/test"
)

func TestTranscriber(t *testing.T) {

	var tests = []struct {
		Desc     string
		HasError bool
		Encoded  evo.Substrate
		Decoded  evo.Substrate
	}{
		{
			Desc:     "empty genome",
			HasError: false,
			Encoded:  evo.Substrate{},
			Decoded:  evo.Substrate{},
		},
		{
			Desc:     "just nodes", // won't make a good network but is OK for transcribing
			HasError: false,
			Encoded: evo.Substrate{
				Nodes: []evo.Node{
					{Position: evo.Position{Layer: 0.0, X: 0.5}, Neuron: evo.Input, Activation: evo.Direct, Bias: 0.0},
					{Position: evo.Position{Layer: 1.0, X: 0.5}, Neuron: evo.Output, Activation: evo.Sigmoid, Bias: 1.5},
				},
			},
			Decoded: evo.Substrate{
				Nodes: []evo.Node{
					{Position: evo.Position{Layer: 0.0, X: 0.5}, Neuron: evo.Input, Activation: evo.Direct, Bias: 0.0},
					{Position: evo.Position{Layer: 1.0, X: 0.5}, Neuron: evo.Output, Activation: evo.Sigmoid, Bias: 1.5},
				},
			},
		},
		{
			Desc:     "just conns", // won't make a good network but is OK for transcribing
			HasError: false,
			Encoded: evo.Substrate{
				Conns: []evo.Conn{
					{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 0.5, X: 0.5}, Weight: 1.1, Enabled: true},
					{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 0.5, X: 0.5}, Weight: 2.2, Enabled: true},
					{Source: evo.Position{Layer: 0.5, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 3.3, Enabled: true},
				},
			},
			Decoded: evo.Substrate{
				Conns: []evo.Conn{
					{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 0.5, X: 0.5}, Weight: 1.1, Enabled: true},
					{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 0.5, X: 0.5}, Weight: 2.2, Enabled: true},
					{Source: evo.Position{Layer: 0.5, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 3.3, Enabled: true},
				},
			},
		},

		// All connections are enabled
		{
			Desc:     "all enabled",
			HasError: false,
			Encoded: evo.Substrate{
				Nodes: []evo.Node{
					{Position: evo.Position{Layer: 0.0, X: 0.0}, Neuron: evo.Input, Activation: evo.Direct, Bias: 0.0},
					{Position: evo.Position{Layer: 0.0, X: 1.0}, Neuron: evo.Input, Activation: evo.Direct, Bias: 0.0},
					{Position: evo.Position{Layer: 0.5, X: 0.5}, Neuron: evo.Output, Activation: evo.Sigmoid, Bias: 0.1},
					{Position: evo.Position{Layer: 1.0, X: 0.5}, Neuron: evo.Output, Activation: evo.Sigmoid, Bias: 0.2},
				},
				Conns: []evo.Conn{
					{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 0.5, X: 0.5}, Weight: 1.1, Enabled: true},
					{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 0.5, X: 0.5}, Weight: 2.2, Enabled: true},
					{Source: evo.Position{Layer: 0.5, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 3.3, Enabled: true},
				},
			},
			Decoded: evo.Substrate{
				Nodes: []evo.Node{
					{Position: evo.Position{Layer: 0.0, X: 0.0}, Neuron: evo.Input, Activation: evo.Direct, Bias: 0.0},
					{Position: evo.Position{Layer: 0.0, X: 1.0}, Neuron: evo.Input, Activation: evo.Direct, Bias: 0.0},
					{Position: evo.Position{Layer: 0.5, X: 0.5}, Neuron: evo.Output, Activation: evo.Sigmoid, Bias: 0.1},
					{Position: evo.Position{Layer: 1.0, X: 0.5}, Neuron: evo.Output, Activation: evo.Sigmoid, Bias: 0.2},
				},
				Conns: []evo.Conn{
					{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 0.5, X: 0.5}, Weight: 1.1, Enabled: true},
					{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 0.5, X: 0.5}, Weight: 2.2, Enabled: true},
					{Source: evo.Position{Layer: 0.5, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 3.3, Enabled: true},
				},
			},
		},
		// Some connections are disabled
		{
			Desc:     "some enabled",
			HasError: false,
			Encoded: evo.Substrate{
				Nodes: []evo.Node{
					{Position: evo.Position{Layer: 0.0, X: 0.0}, Neuron: evo.Input, Activation: evo.Direct, Bias: 0.0},
					{Position: evo.Position{Layer: 0.0, X: 1.0}, Neuron: evo.Input, Activation: evo.Direct, Bias: 0.0},
					{Position: evo.Position{Layer: 0.5, X: 0.5}, Neuron: evo.Output, Activation: evo.Sigmoid, Bias: 0.1},
					{Position: evo.Position{Layer: 1.0, X: 0.5}, Neuron: evo.Output, Activation: evo.Sigmoid, Bias: 0.2},
				},
				Conns: []evo.Conn{
					{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 0.5, X: 0.5}, Weight: 1.1, Enabled: true},
					{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 0.5, X: 0.5}, Weight: 2.2, Enabled: true},
					{Source: evo.Position{Layer: 0.5, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 3.3, Enabled: false},
				},
			},
			Decoded: evo.Substrate{
				Nodes: []evo.Node{
					{Position: evo.Position{Layer: 0.0, X: 0.0}, Neuron: evo.Input, Activation: evo.Direct, Bias: 0.0},
					{Position: evo.Position{Layer: 0.0, X: 1.0}, Neuron: evo.Input, Activation: evo.Direct, Bias: 0.0},
					{Position: evo.Position{Layer: 0.5, X: 0.5}, Neuron: evo.Output, Activation: evo.Sigmoid, Bias: 0.1},
					{Position: evo.Position{Layer: 1.0, X: 0.5}, Neuron: evo.Output, Activation: evo.Sigmoid, Bias: 0.2},
				},
				Conns: []evo.Conn{
					{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 0.5, X: 0.5}, Weight: 1.1, Enabled: true},
					{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 0.5, X: 0.5}, Weight: 2.2, Enabled: true},
				},
			},
		},
		// All connections are disabled
		{
			Desc:     "all disabled",
			HasError: false,
			Encoded: evo.Substrate{
				Nodes: []evo.Node{
					{Position: evo.Position{Layer: 0.0, X: 0.0}, Neuron: evo.Input, Activation: evo.Direct, Bias: 0.0},
					{Position: evo.Position{Layer: 0.0, X: 1.0}, Neuron: evo.Input, Activation: evo.Direct, Bias: 0.0},
					{Position: evo.Position{Layer: 0.5, X: 0.5}, Neuron: evo.Output, Activation: evo.Sigmoid, Bias: 0.1},
					{Position: evo.Position{Layer: 1.0, X: 0.5}, Neuron: evo.Output, Activation: evo.Sigmoid, Bias: 0.2},
				},
				Conns: []evo.Conn{
					{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 0.5, X: 0.5}, Weight: 1.1, Enabled: false},
					{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 0.5, X: 0.5}, Weight: 2.2, Enabled: false},
					{Source: evo.Position{Layer: 0.5, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 3.3, Enabled: false},
				},
			},
			Decoded: evo.Substrate{
				Nodes: []evo.Node{
					{Position: evo.Position{Layer: 0.0, X: 0.0}, Neuron: evo.Input, Activation: evo.Direct, Bias: 0.0},
					{Position: evo.Position{Layer: 0.0, X: 1.0}, Neuron: evo.Input, Activation: evo.Direct, Bias: 0.0},
					{Position: evo.Position{Layer: 0.5, X: 0.5}, Neuron: evo.Output, Activation: evo.Sigmoid, Bias: 0.1},
					{Position: evo.Position{Layer: 1.0, X: 0.5}, Neuron: evo.Output, Activation: evo.Sigmoid, Bias: 0.2},
				},
				Conns: []evo.Conn{},
			},
		},
	}

	// Run the tests
	for _, c := range tests {
		t.Run(c.Desc, func(t *testing.T) {

			// Create a new transcriber
			z := new(Transcriber)

			// Transcribe the genome
			dec, err := z.Transcribe(context.Background(), c.Encoded)

			// There should be no error
			t.Run("error", test.Error(c.HasError, err))
			if c.HasError {
				return
			}

			// Only the expected nodes should be present
			t.Run("Nodes", testTranscribeNodes(dec.Nodes, c.Decoded.Nodes))

			// Only the expected conns should be present
			t.Run("Conns", testTranscribeConns(dec.Conns, c.Decoded.Conns))

		})
	}
}

func testTranscribeNodes(actual, expected []evo.Node) func(*testing.T) {
	return func(t *testing.T) {

		// The nodes should match
		if len(actual) != len(expected) {
			t.Errorf("incorrect number of nodes. expected %d, actual %d", len(expected), len(actual))
		}
		for _, en := range expected {
			found := false
			for _, an := range actual {
				if an.Compare(en) == 0 {

					// Neuron types should match
					if an.Neuron != en.Neuron {
						t.Errorf("incorrect neuron type. expected %v, actual %v", en.Neuron, an.Neuron)
					}

					// Activation types should match
					if an.Neuron != en.Neuron {
						t.Errorf("incorrect activation type. expected %v, actual %v", en.Activation, an.Activation)
					}

					// Note match and continue to next node
					found = true
					break
				}
			}
			if !found {
				t.Errorf("node not found. expected node at %v", en.Position)
			}
		}
	}
}

func testTranscribeConns(actual, expected []evo.Conn) func(*testing.T) {
	return func(t *testing.T) {

		// Only the enabled connections should exist
		cnt := 0
		for _, c := range expected {
			if c.Enabled {
				cnt++
			}
		}
		if len(actual) != cnt {
			t.Errorf("incorrect number of connections. expected %d, actual %d", cnt, len(actual))
		}
		for _, ec := range expected {
			if ec.Enabled {
				found := false
				for _, ac := range actual {
					if ac.Compare(ec) == 0 {

						// Source should match
						if ac.Source != ec.Source {
							t.Errorf("incorrect source. expected %v, actual %v", ec.Source, ac.Source)
						}

						// Target should match
						if ac.Target != ec.Target {
							t.Errorf("incorrect target. expected %v, actual %v", ec.Target, ac.Target)
						}

						// Weight should match
						if ac.Weight != ec.Weight {
							t.Errorf("incorrect weight. expected %f, actual %f", ec.Weight, ac.Weight)
						}

						// Note match and continue to next connection
						found = true
						break
					}
				}
				if !found {
					t.Errorf("connection not found. expected connection from %v to %v", ec.Source, ec.Target)
				}
			}
		}
	}
}

func TestWithTranscriber(t *testing.T) {
	e := new(evo.Experiment)
	err := WithTranscriber(nil)(e)
	if err != nil {
		t.Errorf("error was not expected: %v", err)
	}
	if _, ok := e.Transcriber.(*Transcriber); !ok {
		t.Errorf("transcriber incorrectly set")
	}
}
