package parallel

import (
	"github.com/klokare/evo"
	"github.com/klokare/evo/internal/workers"
)

// Searcher evaluates phenomes one at a time
type Searcher struct{}

// Search the solution space with the phenomes
func (s Searcher) Search(eval evo.Evaluator, phenomes []evo.Phenome) (results []evo.Result, err error) {

	// Create the tasks
	type task struct {
		phenome evo.Phenome
		result  evo.Result
		err     error
	}

	tasks := make([]workers.Task, len(phenomes))
	for i, p := range phenomes {
		tasks[i] = &task{phenome: p}
	}

	// Perform the tasks
	workers.Do(tasks, func(wt workers.Task) {
		t := wt.(*task)
		t.result, err = eval.Evaluate(t.phenome)
	})

	// Collect the outputs
	results = make([]evo.Result, 0, len(phenomes))
	for _, wt := range tasks {
		t := wt.(*task)
		if t.err != nil {
			err = t.err
			return
		}
		results = append(results, t.result)
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
