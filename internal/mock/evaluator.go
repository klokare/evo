package mock

import (
	"context"
	"errors"

	"github.com/klokare/evo"
)

type Evaluator struct {
	Called   int
	HasError bool
	HasSolve bool
}

func (e *Evaluator) Evaluate(ctx context.Context, p evo.Phenome) (r evo.Result, err error) {
	e.Called++
	if e.HasError {
		err = errors.New("mock evaluator error")
		return
	}
	r = evo.Result{
		ID:      p.ID,
		Fitness: 1.0 / float64(p.ID),
		Novelty: float64(p.ID),
		Solved:  e.HasSolve,
	}
	return
}

func WithEvaluator() evo.Option {
	return func(e *evo.Experiment) error {
		e.Evaluator = &Evaluator{}
		return nil
	}
}
