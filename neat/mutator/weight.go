package mutator

import (
	"github.com/klokare/evo"
)

// Weight mutates the genome's connection weights
type Weight struct {
	MutateWeightProbability  float64 // The probability that the connection's weight will be mutated
	ReplaceWeightProbability float64 // The probability that, if being mutated, the weight will be replaced
	WeightPower              float64
	MaxWeight                float64
}

// Mutate a genome by perturbing or replacing its connections' weights
func (z *Weight) Mutate(g *evo.Genome) (err error) {
	rng := evo.NewRandom()
	for i, c := range g.Encoded.Conns {
		if rng.Float64() < z.MutateWeightProbability {
			if rng.Float64() < z.ReplaceWeightProbability {
				c.Weight = rng.NormFloat64() * z.WeightPower
			} else {
				c.Weight += rng.NormFloat64() * z.WeightPower
			}
			if c.Weight >= z.MaxWeight {
				c.Weight = z.MaxWeight
			} else if c.Weight <= -z.MaxWeight {
				c.Weight = -z.MaxWeight
			}
			g.Encoded.Conns[i] = c
		}
	}
	return
}

// WithWeight adds the configured weight mutator to the experiment
func WithWeight(cfg evo.Configurer) evo.Option {
	return func(e *evo.Experiment) (err error) {
		z := new(Weight)
		if err = cfg.Configure(z); err != nil {
			return
		}

		// Do not continue if there is no chance for mutation
		if z.MutateWeightProbability == 0.0 {
			return
		}

		e.Mutators = append(e.Mutators, z)
		return
	}
}
