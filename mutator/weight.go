package mutator

import (
	"github.com/klokare/evo"
	"github.com/klokare/random"
)

// A Weight mutator mutates a genome's encoded connection weights
type Weight struct {
	MutateWeightProbability  float64
	ReplaceWeightProbability float64
}

// Mutate a genome's encoded connection weights
func (h *Weight) Mutate(g *evo.Genome) error {
	rng := random.New()
	for i := 0; i < len(g.Encoded.Conns); i++ {
		if rng.Float64() < h.MutateWeightProbability {
			if rng.Float64() < h.ReplaceWeightProbability {
				g.Encoded.Conns[i].Weight = rng.NormFloat64()
			} else {
				g.Encoded.Conns[i].Weight += rng.NormFloat64()
			}
		}
	}
	return nil
}
