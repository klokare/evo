package searcher

import "github.com/klokare/evo"

// Serial searches phenomes one at a time
type Serial struct {
	evo.Evaluator
}

// Search all the phenomes and return the results.
func (h *Serial) Search(ps []evo.Phenome) (rs []evo.Result, err error) {
	rs = make([]evo.Result, 0, len(ps))
	for _, p := range ps {
		var r evo.Result
		if r, err = h.Evaluator.Evaluate(p); err != nil {
			return
		}
		rs = append(rs, r)
	}
	return
}
