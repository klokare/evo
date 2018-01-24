package serial

import (
	"context"
	"errors"
	"testing"

	"github.com/klokare/evo"
)

func TestSearcherSearch(t *testing.T) {

	// Some phenomes
	ps := []evo.Phenome{{ID: 1}, {ID: 2}, {ID: 3}}

	// Has error if evaluator has error
	t.Run("evaluator error", func(t *testing.T) {
		e := &MockEvaluator{HasError: true}
		_, err := new(Searcher).Search(context.Background(), e, ps)
		if !t.Run("error", testError(e.HasError, err)) || e.HasError {
			return
		}
	})

	// Execute without errors
	t.Run("evaluator succeeds", func(t *testing.T) {
		e := &MockEvaluator{HasError: false}
		rs, err := new(Searcher).Search(context.Background(), e, ps)
		if !t.Run("no error", testError(e.HasError, err)) || e.HasError {
			return
		}

		// Ensure the results
		if len(ps) != len(rs) {
			t.Errorf("incorrect number of results: expected %d, actual %d", len(ps), len(rs))
		} else {
			for _, p := range ps {
				found := false
				for _, r := range rs {
					if p.ID == r.ID {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("result for phenome %d not found", p.ID)
				}
			}
		}
	})

}

func testError(hasError bool, err error) func(*testing.T) {
	return func(t *testing.T) {

		// An error was expected
		if hasError && err == nil {
			t.Errorf("error was expected but none was returned")
			t.FailNow()
		}

		// No error was expected
		if !hasError && err != nil {
			t.Errorf("no error was expected. actual: %v", err)
			t.FailNow()
		}
	}
}

func TestWithSearcher(t *testing.T) {
	e := new(evo.Experiment)
	err := WithSearcher()(e)
	if err != nil {
		t.Errorf("error was not expected: %v", err)
	}
	if _, ok := e.Searcher.(*Searcher); !ok {
		t.Errorf("searcher incorrectly set")
	}
}

type MockEvaluator struct {
	Called   int
	HasError bool
	HasSolve bool
}

func (e *MockEvaluator) Evaluate(ctx context.Context, p evo.Phenome) (r evo.Result, err error) {
	e.Called++
	if e.HasError {
		err = errors.New("mock transcriber error")
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
