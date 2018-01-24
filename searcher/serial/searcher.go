package serial

import (
	"context"

	"github.com/klokare/evo"
)

// Searcher evaluates phenomes one at a time
type Searcher struct{}

// Search the solution space with the phenomes
func (s Searcher) Search(ctx context.Context, eval evo.Evaluator, phenomes []evo.Phenome) (results []evo.Result, err error) {
	results = make([]evo.Result, len(phenomes))
	for i, p := range phenomes {
		if results[i], err = eval.Evaluate(ctx, p); err != nil {
			return
		}
	}
	return
}

// WithSearcher sets the experiment's transcriber to a configured serial searcher
func WithSearcher() evo.Option {
	return func(e *evo.Experiment) (err error) {
		e.Searcher = new(Searcher)
		return
	}
}
