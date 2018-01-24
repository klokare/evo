package xor

import (
	"context"
	"errors"
	"testing"

	"github.com/klokare/evo"
	"github.com/klokare/evo/internal/test"
)

func TestEvaluator(t *testing.T) {

	var cases = []struct {
		Desc        string
		HasSolution bool
		HasError    bool
		Expected    evo.Result
	}{
		{
			Desc:     "has error",
			HasError: true,
		},
		{
			Desc:        "unsolved",
			HasSolution: false,
			HasError:    false,
			Expected: evo.Result{
				ID:      1,
				Fitness: 0,
				Solved:  false,
			},
		},
		{
			Desc:        "solved",
			HasSolution: true,
			HasError:    false,
			Expected: evo.Result{
				ID:      1,
				Fitness: 16.0,
				Solved:  true,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {

			// Create the phenome
			p := evo.Phenome{
				ID: 1,
				Network: MockNetwork{
					HasSolution: c.HasSolution,
					HasError:    c.HasError,
				},
			}

			// Evaluate
			r, err := new(Evaluator).Evaluate(context.Background(), p)

			// Test for error
			if !t.Run("error", test.Error(c.HasError, err)) || c.HasError {
				return
			}

			// Check result
			if c.Expected.ID != r.ID {
				t.Errorf("incorrect id in result: expected %d, actual %d", c.Expected.ID, r.ID)
			}
			if c.Expected.Fitness != r.Fitness {
				t.Errorf("incorrect fitness in result: expected %f, actual %f", c.Expected.Fitness, r.Fitness)
			}
			if c.Expected.Solved != r.Solved {
				t.Errorf("incorrect solved in result: expected %v, actual %v", c.Expected.Solved, r.Solved)
			}
		})
	}
}

type MockNetwork struct {
	HasSolution bool
	HasError    bool
}

func (n MockNetwork) Activate(inputs []float64) ([]float64, error) {
	if n.HasError {
		return nil, errors.New("mock network error")
	}
	var out float64
	switch {
	case inputs[0] == 0 && inputs[1] == 0:
		out = 0
	case inputs[0] == 0 && inputs[1] == 1:
		out = 1
	case inputs[0] == 1 && inputs[1] == 0:
		out = 1
	case inputs[0] == 1 && inputs[1] == 1:
		out = 0
	}
	if !n.HasSolution {
		out = 1.0 - out
	}
	return []float64{out}, nil
}
