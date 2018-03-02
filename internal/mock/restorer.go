package mock

import (
	"errors"

	"github.com/klokare/evo"
)

// Restorer ...
type Restorer struct {
	Called   int
	HasError bool
}

// Seed ...
func (m *Restorer) Seed() (genomes []evo.Genome, err error) {
	m.Called++
	if m.HasError {
		err = errors.New("mock restorer error")
		return
	}
	genomes = []evo.Genome{
		{
			ID:        8,
			SpeciesID: 5,
			Encoded: evo.Substrate{
				Nodes: []evo.Node{
					{Position: evo.Position{Layer: 0.5, X: 0.5}, Neuron: evo.Hidden, Activation: evo.InverseAbs},
					{Position: evo.Position{Layer: 1.0, X: 0.5}, Neuron: evo.Output, Activation: evo.Sigmoid},
				},
				Conns: []evo.Conn{
					{Source: evo.Position{Layer: 0.5, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 0.5}, Enabled: true, Weight: 2.5},
				},
			},
		},
		{
			ID:        9,
			SpeciesID: 6,
			Encoded: evo.Substrate{
				Nodes: []evo.Node{
					{Position: evo.Position{Layer: 0.0, X: 0.5}, Neuron: evo.Input, Activation: evo.Direct},
					{Position: evo.Position{Layer: 0.5, X: 0.5}, Neuron: evo.Hidden, Activation: evo.InverseAbs},
					{Position: evo.Position{Layer: 1.0, X: 0.5}, Neuron: evo.Output, Activation: evo.Sigmoid},
				},
				Conns: []evo.Conn{
					{Source: evo.Position{Layer: 0.0, X: 0.5}, Target: evo.Position{Layer: 0.0, X: 0.5}, Enabled: false, Weight: 1.5},
					{Source: evo.Position{Layer: 0.5, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 0.5}, Enabled: true, Weight: 2.5},
				},
			},
		},
	}
	return
}

// WithRestorer ...
func WithRestorer() evo.Option {
	return func(e *evo.Experiment) error {
		e.Seeder = &Restorer{}
		return nil
	}
}
