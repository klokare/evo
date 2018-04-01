package mutator

import (
	"errors"

	"github.com/klokare/evo"
)

// Known errors
var (
	ErrMissingActivations = errors.New("activation mutator requires list of possible activation functions")
)

// Activation is a helper that mutates the activation value of nodes
type Activation struct {
	ReplaceActivationProbability float64
	Activations                  []evo.Activation
}

// Mutate the the activation values of the genomes based on the settings in the helper.
func (a Activation) Mutate(g *evo.Genome) (err error) {
	if len(a.Activations) == 0 {
		return ErrMissingActivations
	}

	rng := evo.NewRandom()
	for i, n := range g.Encoded.Nodes {
		if n.Neuron == evo.Hidden {
			if rng.Float64() < a.ReplaceActivationProbability {
				g.Encoded.Nodes[i].Activation = a.Activations[rng.Intn(len(a.Activations))]
			}
		}
	}
	return
}
