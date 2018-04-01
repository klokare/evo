package parallel

import (
	"github.com/klokare/evo"
	"github.com/klokare/evo/internal/workers"
)

// Searcher evaluates phenomes one at a time
type Searcher struct{}

// Search the solution space with the phenomes
func (s Searcher) Search(eval evo.Evaluator, phenomes []evo.Phenome) (results []evo.Result, err error) {

	// Receive results
	results = make([]evo.Result, 0, len(phenomes))
	ch := make(chan evo.Result, len(phenomes))
	done := make(chan struct{})
	go func(ch <-chan evo.Result, done chan struct{}) {
		defer close(done)
		for r := range ch {
			results = append(results, r)
		}
	}(ch, done)

	// Create the tasks
	tasks := make([]workers.Task, len(phenomes))
	for i := 0; i < len(phenomes); i++ {
		tasks[i] = phenomes[i]
	}

	// Perform the tasks
	err = workers.Do(tasks, func(wt workers.Task) (err error) {
		var r evo.Result
		p := wt.(evo.Phenome)
		if r, err = eval.Evaluate(p); err != nil {
			return
		}
		ch <- r
		return
	})
	close(ch)
	<-done
	return
}
