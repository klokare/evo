package serial

import (
	"github.com/klokare/evo"
)

// Searcher evaluates phenomes one at a time
type Searcher struct{}

// Search the solution space with the phenomes
func (s Searcher) Search(eval evo.Evaluator, phenomes []evo.Phenome) (results []evo.Result, err error) {
	results = make([]evo.Result, len(phenomes))
	for i, p := range phenomes {
		if results[i], err = eval.Evaluate(p); err != nil {
			return
		}
	}
	return
}
