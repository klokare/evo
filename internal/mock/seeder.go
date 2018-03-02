package mock

import (
	"errors"

	"github.com/klokare/evo"
)

// Seeder provides the initial genome
type Seeder struct {
	Called   int
	HasError bool
}

// Seed creates a new genome to be used for seeding a population
func (m *Seeder) Seed() (genomes []evo.Genome, err error) {
	m.Called++
	if m.HasError {
		err = errors.New("mock seeder error")
		return
	}
	genomes = []evo.Genome{
		{
			ID: 3,
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
			ID: 2,
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
		{
			ID: 1,
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
			ID: 4,
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

// WithSeeder sets the experiment's seeder helper
func WithSeeder() evo.Option {
	return func(e *evo.Experiment) error {
		e.Seeder = &Seeder{}
		return nil
	}
}
