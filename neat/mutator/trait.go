package mutator

import (
	"github.com/klokare/evo"
)

// Trait is a helper that mutates the trait value of nodes
type Trait struct {
	MutateTraitProbability  float64
	ReplaceTraitProbability float64
}

// Mutate the the trait values of the genomes based on the settings in the helper.
func (b Trait) Mutate(g *evo.Genome) (err error) {
	rng := evo.NewRandom()
	for i, x := range g.Traits {
		if rng.Float64() < b.MutateTraitProbability {
			if rng.Float64() < b.ReplaceTraitProbability {
				x = rng.Float64()
			} else {
				x += rng.NormFloat64()
				if x > 1.0 {
					x = 1.0
				} else if x < 0.0 {
					x = 0.0
				}
			}
			g.Traits[i] = x
		}
	}
	return
}
