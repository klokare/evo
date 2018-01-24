package mutator

import (
	"context"

	"github.com/klokare/evo"
)

// Bias is a helper that mutates the bias value of nodes
type Bias struct {
	MutateBiasProbability  float64
	ReplaceBiasProbability float64
	BiasPower              float64
	MaxBias                float64
}

// Mutate the the bias values of the genomes based on the settings in the helper. In Stanley's
// version of NEAT, bias is a separate node type with connections to other nodes. In evo, bias is a
// property of the node. Effectively, though, they are the same.
func (b Bias) Mutate(ctx context.Context, g *evo.Genome) (err error) {
	rng := evo.NewRandom()
	for i, n := range g.Encoded.Nodes {
		if n.Neuron == evo.Input {
			continue
		}
		if rng.Float64() < b.MutateBiasProbability {
			if rng.Float64() < b.ReplaceBiasProbability {
				n.Bias = rng.NormFloat64() * b.BiasPower
			} else {
				n.Bias += rng.NormFloat64()
			}
			if n.Bias < -b.MaxBias {
				n.Bias = -b.MaxBias
			} else if n.Bias > b.MaxBias {
				n.Bias = b.MaxBias
			}
			g.Encoded.Nodes[i] = n
		}
	}
	return
}

// WithBias adds the configured bias mutator to the experiment
func WithBias(cfg evo.Configurer) evo.Option {
	return func(e *evo.Experiment) (err error) {
		z := new(Bias)
		if err = cfg.Configure(z); err != nil {
			return
		}

		// Do not continue if there is no chance for mutation
		if z.MutateBiasProbability == 0.0 {
			return
		}

		e.Mutators = append(e.Mutators, z)
		return
	}
}
