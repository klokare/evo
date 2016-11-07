package searcher

import (
	"sync"

	"github.com/klokare/errors"
	"github.com/klokare/evo"
)

// Concurrent searches phenomes, well, concurrently
type Concurrent struct {
	evo.Evaluator
}

func (h Concurrent) String() string { return "evo.searcher.Concurrent{}" }

// Search all the phenomes and return the results.
func (h *Concurrent) Search(ps []evo.Phenome) (rs []evo.Result, err error) {
	rc := make(chan evo.Result, len(ps))
	wg := new(sync.WaitGroup)
	z := new(errors.Safe)
	for _, p := range ps {
		wg.Add(1)
		go func(p evo.Phenome) {
			var r evo.Result
			if r, err = h.Evaluator.Evaluate(p); err != nil {
				z.Add(err)
			}
			rc <- r
			wg.Done()
		}(p)
	}
	wg.Wait()
	close(rc)

	rs = make([]evo.Result, 0, len(ps))
	for r := range rc {
		rs = append(rs, r)
	}
	return rs, z.Err()
}
