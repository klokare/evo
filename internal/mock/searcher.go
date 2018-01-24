package mock

import (
	"context"
	"errors"

	"github.com/klokare/evo"
)

type Searcher struct {
	Called   int
	HasError bool
}

func (s *Searcher) Search(ctx context.Context, e evo.Evaluator, ps []evo.Phenome) (rs []evo.Result, err error) {
	s.Called++
	if s.HasError {
		err = errors.New("mock searcher error")
		return
	}
	rs = make([]evo.Result, len(ps))
	for i, p := range ps {
		if rs[i], err = e.Evaluate(ctx, p); err != nil {
			return
		}
	}
	return
}

func WithSearcher() evo.Option {
	return func(e *evo.Experiment) error {
		e.Searcher = &Searcher{}
		return nil
	}
}
