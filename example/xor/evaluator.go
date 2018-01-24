package xor

import (
	"context"
	"math"

	"github.com/klokare/evo"
)

type Evaluator struct{}

func (e Evaluator) Evaluate(ctx context.Context, p evo.Phenome) (r evo.Result, err error) {
	in := [][]float64{{0, 0}, {1, 0}, {0, 1}, {1, 1}}
	out := make([]float64, len(in))
	var outputs []float64
	for i, inputs := range in {
		if outputs, err = p.Activate(inputs); err != nil {
			return
		}
		out[i] = outputs[0]
	}
	r = evo.Result{
		ID:      p.ID,
		Fitness: math.Pow(4.0-(out[0]+(1-out[1])+(1-out[2])+out[3]), 2.0),
		Solved:  out[0] < 0.5 && out[1] > 0.5 && out[2] > 0.5 && out[3] < 0.5,
	}
	return
}

func WithEvaluator() evo.Option {
	return func(e *evo.Experiment) error {
		e.Evaluator = &Evaluator{}
		return nil
	}
}
