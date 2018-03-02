package neat

import (
	"errors"

	"github.com/klokare/evo"
	"github.com/klokare/evo/internal/float"
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
	BiasPower        float64
	MaxBias          float64
	WeightPower      float64
	MaxWeight        float64
}

// Seed creates an initial genome from the specifications.
func (s Seeder) Seed() (genomes []evo.Genome, err error) {

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
	rng := evo.NewRandom()
	g := evo.Genome{
		Fitness: evo.MinFitness,
		Traits:  make([]float64, s.NumTraits),
		Encoded: evo.Substrate{
			Nodes: make([]evo.Node, 0, s.NumInputs+s.NumOutputs),
			Conns: make([]evo.Conn, 0, s.NumInputs*s.NumOutputs),
		},
	}

	// Create the traits
	for i := 0; i < s.NumTraits; i++ {
		g.Traits[i] = rng.Float64()
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
			Bias:       float.Min(s.MaxBias, float.Max(-s.MaxBias, rng.NormFloat64()*s.BiasPower)),
		})
	}

	// Connect the sensors to the outputs
	for _, src := range g.Encoded.Nodes[:s.NumInputs] {
		for _, tgt := range g.Encoded.Nodes[s.NumInputs:] {
			g.Encoded.Conns = append(g.Encoded.Conns, evo.Conn{
				Source:  src.Position,
				Target:  tgt.Position,
				Weight:  float.Min(s.MaxWeight, float.Max(-s.MaxWeight, rng.NormFloat64()*s.WeightPower)),
				Enabled: true,
			})
		}
	}

	// Return the new genome
	genomes = []evo.Genome{g}
	return
}

// WithSeeder sets the experiment's seeder to a configured NEAT seeder
func WithSeeder(cfg evo.Configurer) evo.Option {
	return func(e *evo.Experiment) (err error) {
		s := new(Seeder)
		if err = cfg.Configure(s); err != nil {
			return
		}
		e.Seeder = s
		return
	}
}
