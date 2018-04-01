package serial

import (
	"testing"

	"github.com/klokare/evo"
	"github.com/klokare/evo/internal/mock"
)

func TestSearcherSearch(t *testing.T) {

	// Some phenomes
	ps := []evo.Phenome{{ID: 1}, {ID: 2}, {ID: 3}}

	// Has error if evaluator has error
	t.Run("evaluator error", func(t *testing.T) {
		e := &mock.Evaluator{HasError: true}
		_, err := new(Searcher).Search(e, ps)
		if !t.Run("error", mock.Error(e.HasError, err)) || e.HasError {
			return
		}
	})

	// Execute without errors
	t.Run("evaluator succeeds", func(t *testing.T) {
		e := &mock.Evaluator{HasError: false}
		rs, err := new(Searcher).Search(e, ps)
		if !t.Run("no error", mock.Error(e.HasError, err)) || e.HasError {
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
