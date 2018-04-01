package neat

import (
	"errors"
	"sort"

	"github.com/klokare/evo"
)

// Errors related to the Seeder
var (
	ErrInvalidNumInputs  = errors.New("num inputs must be greater than or equal to zero")
	ErrInvalidNumOutputs = errors.New("num outputs must be greater than zero")
)

// Seeder creates the initial population
type Seeder struct {
	NumInputs        int
	NumOutputs       int
	NumTraits        int
	DisconnectRate   float64
	OutputActivation evo.Activation
}

// Seed creates an unitialised genome from the specifications.
func (s Seeder) Seed() (g evo.Genome, err error) {

	// Check for errors
	if s.NumInputs <= 0 {
		err = ErrInvalidNumInputs
		return
	}
	if s.NumOutputs <= 0 {
		err = ErrInvalidNumOutputs
		return
	}

	// Create the seed genome
	g = evo.Genome{
		Traits: make([]float64, s.NumTraits),
		Encoded: evo.Substrate{
			Nodes: make([]evo.Node, 0, s.NumInputs+s.NumOutputs),
			Conns: make([]evo.Conn, 0, s.NumInputs*s.NumOutputs),
		},
	}

	// Add the input nodes
	for i := 0; i < s.NumInputs; i++ {
		x := 0.5
		if s.NumInputs > 1 {
			x = float64(i) / float64(s.NumInputs-1)
		}
		g.Encoded.Nodes = append(g.Encoded.Nodes, evo.Node{
			Position:   evo.Position{Layer: 0.0, X: x},
			Neuron:     evo.Input,
			Activation: evo.Direct,
		})
	}

	// Add the output nodes
	for i := 0; i < s.NumOutputs; i++ {
		x := 0.5
		if s.NumOutputs > 1 {
			x = float64(i) / float64(s.NumOutputs-1)
		}
		g.Encoded.Nodes = append(g.Encoded.Nodes, evo.Node{
			Position:   evo.Position{Layer: 1.0, X: x},
			Neuron:     evo.Output,
			Activation: s.OutputActivation,
		})
	}

	// Connect the sensors to the outputs
	if s.DisconnectRate < 1.0 {
		rng := evo.NewRandom()
		for _, src := range g.Encoded.Nodes[:s.NumInputs] {
			for _, tgt := range g.Encoded.Nodes[s.NumInputs:] {
				if rng.Float64() > s.DisconnectRate {
					g.Encoded.Conns = append(g.Encoded.Conns, evo.Conn{
						Source:  src.Position,
						Target:  tgt.Position,
						Enabled: true,
					})
				}
			}
		}
	}

	// Ensure substrate is sorted
	sort.Slice(g.Encoded.Nodes, func(i, j int) bool { return g.Encoded.Nodes[i].Compare(g.Encoded.Nodes[j]) < 0 })
	sort.Slice(g.Encoded.Conns, func(i, j int) bool { return g.Encoded.Conns[i].Compare(g.Encoded.Conns[j]) < 0 })
	return
}
