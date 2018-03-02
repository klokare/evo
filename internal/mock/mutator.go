package mock

import (
	"errors"

	"github.com/klokare/evo"
)

// Mutator ...
type Mutator struct {
	HasError        bool
	MutateStructure bool
	Count           int
}

// Mutate ...
func (m *Mutator) Mutate(g *evo.Genome) error {
	m.Count++
	if m.HasError {
		return errors.New("mock mutator error")
	}
	if m.MutateStructure {
		g.Encoded.Nodes = append(g.Encoded.Nodes, evo.Node{})
	}
	return nil
}

// StructureMutator ...
type StructureMutator struct {
	Called   int
	HasError bool
}

// Mutate ...
func (m *StructureMutator) Mutate(g *evo.Genome) (err error) {
	m.Called++
	if m.HasError {
		err = errors.New("mock structure mutator error")
		return
	}
	// Add a new connection if ID is odd
	if g.ID%2 == 1 {
		g.Encoded.Conns = append(g.Encoded.Conns, evo.Conn{
			Source:  evo.Position{Layer: 0.25, X: 0.25},
			Target:  evo.Position{Layer: 0.75, X: 0.75},
			Weight:  1.23,
			Enabled: true,
		})
	}
	return
}

// WeightMutator ...
type WeightMutator struct {
	Called   int
	HasError bool
}

// Mutate ...
func (m *WeightMutator) Mutate(g *evo.Genome) (err error) {
	m.Called++
	if m.HasError {
		err = errors.New("mock weight mutator error")
		return
	}
	if len(g.Encoded.Conns) > 0 {
		g.Encoded.Conns[0].Weight += 0.001
	}
	return
}

// WithMutators ...
func WithMutators() evo.Option {
	return func(e *evo.Experiment) error {
		e.Mutators = append(e.Mutators, &StructureMutator{}, &WeightMutator{})
		return nil
	}
}
