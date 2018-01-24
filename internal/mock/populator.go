package mock

import (
	"context"
	"errors"

	"github.com/klokare/evo"
)

type Populator struct {
	Called   int
	HasError bool
}

func (m *Populator) Populate(context.Context) (p evo.Population, err error) {
	m.Called++
	if m.HasError {
		err = errors.New("mock populator error")
		return
	}
	p = evo.Population{
		Species: []evo.Species{{ID: 5}}, // Species do not persist in load
		Genomes: []evo.Genome{
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
		},
	}
	return
}

func WithPopulator() evo.Option {
	return func(e *evo.Experiment) error {
		e.Populator = &Populator{}
		return nil
	}
}
