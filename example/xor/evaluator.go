package xor

import (
	"math"

	"gonum.org/v1/gonum/mat"

	"github.com/klokare/evo"
)

// Evaluator runs the XOR experiment
type Evaluator struct{}

// Evaluate the XOR experiment with this phenome
func (e Evaluator) Evaluate(p evo.Phenome) (r evo.Result, err error) {

	var outputs evo.Matrix
	in := mat.NewDense(4, 2, []float64{
		0, 0,
		1, 0,
		0, 1,
		1, 1,
	})
	if outputs, err = p.Activate(in); err != nil {
		return
	}

	var out []float64
	if rcv, ok := outputs.(mat.RawColViewer); ok {
		out = rcv.RawColView(0)
	} else {
		out = make([]float64, 4)
		for i := 0; i < 4; i++ {
			out[i] = outputs.At(i, 0)
		}
	}
	r = evo.Result{
		ID:      p.ID,
		Fitness: math.Pow(4.0-(out[0]+(1-out[1])+(1-out[2])+out[3]), 2.0),
		Solved:  out[0] < 0.5 && out[1] > 0.5 && out[2] > 0.5 && out[3] < 0.5,
	}
	return
}

// WithEvaluator configures the experiment with the XOR evaluator
func WithEvaluator() evo.Option {
	return func(e *evo.Experiment) error {
		e.Evaluator = Evaluator{}
		return nil
	}
}
