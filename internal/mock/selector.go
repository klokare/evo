package mock

import (
	"context"
	"errors"

	"github.com/klokare/evo"
)

type Selector struct {
	Called                 int
	HasError               bool
	FailWithTooManyParents bool // Special case to check for incorrect population size
}

func (s *Selector) Select(ctx context.Context, p evo.Population) (continuing []evo.Genome, parents [][]evo.Genome, err error) {
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

func WithSelector() evo.Option {
	return func(e *evo.Experiment) error {
		e.Selector = &Selector{}
		return nil
	}
}
