package mutator

import (
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
func (b Bias) Mutate(g *evo.Genome) (err error) {
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
