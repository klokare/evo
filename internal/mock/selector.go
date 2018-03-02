package mock

import (
	"errors"

	"github.com/klokare/evo"
)

// Selector ...
type Selector struct {
	Called                 int
	HasError               bool
	FailWithTooManyParents bool // Special case to check for incorrect population size
	MutateOnlyProbability  float64
	mop                    float64
}

// Select ...
func (s *Selector) Select(p evo.Population) (continuing []evo.Genome, parents [][]evo.Genome, err error) {
	s.Called++
	if s.HasError {
		err = errors.New("mock selector error")
		return
	}
	continuing = p.Genomes[:2]
	parents = make([][]evo.Genome, len(p.Genomes)-2)
	for i := 0; i < len(parents); i++ {
		parents[i] = []evo.Genome{p.Genomes[i+2]}
	}
	if s.FailWithTooManyParents {
		parents = append(parents, continuing)
	}
	return
}

// ToggleMutateOnly ...
func (s *Selector) ToggleMutateOnly(on bool) error {
	if s.mop == 0.0 {
		s.mop = s.MutateOnlyProbability
	}
	if on {
		s.MutateOnlyProbability = 1.0
	} else {
		s.MutateOnlyProbability = s.mop
	}
	return nil
}

// WithSelector ...
func WithSelector() evo.Option {
	return func(e *evo.Experiment) error {
		e.Selector = &Selector{}
		return nil
	}
}
