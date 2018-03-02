package mock

import (
	"errors"

	"github.com/klokare/evo"
)

// Searcher ...
type Searcher struct {
	Called   int
	HasError bool
}

// Search ...
func (s *Searcher) Search(e evo.Evaluator, ps []evo.Phenome) (rs []evo.Result, err error) {
	s.Called++
	if s.HasError {
		err = errors.New("mock searcher error")
		return
	}
	rs = make([]evo.Result, len(ps))
	for i, p := range ps {
		if rs[i], err = e.Evaluate(p); err != nil {
			return
		}
	}
	return
}

// WithSearcher ...
func WithSearcher() evo.Option {
	return func(e *evo.Experiment) error {
		e.Searcher = &Searcher{}
		return nil
	}
}
