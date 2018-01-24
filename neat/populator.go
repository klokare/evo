package neat

import (
	"context"
	"errors"

	"github.com/klokare/evo"
)

// Errors related to the Populator
var (
	ErrInvalidPopulationSize = errors.New("invalid population size")
	ErrInvalidNumInputs      = errors.New("num inputs must be greater than or equal to zero")
	ErrInvalidNumOutputs     = errors.New("num outputs must be greater than zero")
)

// Populator creates the initial population
type Populator struct {
	PopulationSize   int
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

// Populate creates an initial population from the specifications. An initial seed genome is
// created and the population is filled with mutated clones of it. Note: the seed genome will be
// created with desired disconnected state but, dependening on the mutators used, there may be
// more connections in some genomes.
func (p Populator) Populate(ctx context.Context) (pop evo.Population, err error) {

	// Check for errors
	if p.PopulationSize <= 0 {
		err = ErrInvalidPopulationSize
		return
	}
	if p.NumInputs <= 0 {
		err = ErrInvalidNumInputs
		return
	}
	if p.NumOutputs <= 0 {
		err = ErrInvalidNumOutputs
		return
	}

	// Create the seed genome
	rng := evo.NewRandom()
	g0 := seed(rng, p.NumInputs, p.NumOutputs, p.NumTraits, p.BiasPower, p.WeightPower, p.OutputActivation)

	// Create the population
	pop.Genomes = make([]evo.Genome, 0, p.PopulationSize)
	for i := 0; i < p.PopulationSize; i++ {

		// Clone the seed genome
		g1 := evo.Genome{
			Traits: make([]float64, p.NumTraits),
			Encoded: evo.Substrate{
				Nodes: make([]evo.Node, len(g0.Encoded.Nodes)),
				Conns: make([]evo.Conn, 0, len(g0.Encoded.Conns)),
			},
		}

		// Create variations on the traits
		copy(g1.Traits, g0.Traits)
		for i := 0; i < p.NumTraits; i++ {
			g1.Traits[i] = rng.Float64()
		}

		// Create variations on the nodes
		copy(g1.Encoded.Nodes, g0.Encoded.Nodes)
		for i, node := range g1.Encoded.Nodes {
			if node.Neuron != evo.Input {
				node.Bias = rng.NormFloat64() * p.BiasPower
				if node.Bias > p.MaxBias {
					node.Bias = p.MaxBias
				} else if node.Bias < -p.MaxBias {
					node.Bias = p.MaxBias
				}
				g1.Encoded.Nodes[i] = node
			}
		}

		// Create variations on the nodes and disconnect at the chosen rate
		connect := 1.0 - p.DisconnectRate
		for _, conn := range g0.Encoded.Conns {
			if rng.Float64() < connect {
				conn.Weight = rng.NormFloat64() * p.WeightPower
				if conn.Weight > p.MaxWeight {
					conn.Weight = p.MaxWeight
				} else if conn.Weight < -p.MaxWeight {
					conn.Weight = p.MaxWeight
				}
				g1.Encoded.Conns = append(g1.Encoded.Conns, conn)
			}
		}

		// Add the genome to the population
		pop.Genomes = append(pop.Genomes, g1)
	}
	return
}

// seed creates the initial genome given the requirements.
func seed(rng evo.Random, inputs, outputs, traits int, bp, wp float64, fn evo.Activation) evo.Genome {
	g := evo.Genome{
		Traits: make([]float64, traits),
		Encoded: evo.Substrate{
			Nodes: make([]evo.Node, 0, inputs+outputs),
			Conns: make([]evo.Conn, 0, inputs*outputs),
		},
	}

	// Create the traits
	for i := 0; i < traits; i++ {
		g.Traits[i] = rng.Float64()
	}

	// Add the input nodes
	for i := 0; i < inputs; i++ {
		x := 0.5
		if inputs > 1 {
			x = float64(i) / float64(inputs-1)
		}
		g.Encoded.Nodes = append(g.Encoded.Nodes, evo.Node{
			Position:   evo.Position{Layer: 0.0, X: x},
			Neuron:     evo.Input,
			Activation: evo.Direct,
		})
	}

	// Add the output nodes
	for i := 0; i < outputs; i++ {
		x := 0.5
		if outputs > 1 {
			x = float64(i) / float64(outputs-1)
		}
		g.Encoded.Nodes = append(g.Encoded.Nodes, evo.Node{
			Position:   evo.Position{Layer: 1.0, X: x},
			Neuron:     evo.Output,
			Activation: fn,
		})
	}

	// Connect the sensors to the outputs
	for _, src := range g.Encoded.Nodes[:inputs] {
		for _, tgt := range g.Encoded.Nodes[inputs:] {
			g.Encoded.Conns = append(g.Encoded.Conns, evo.Conn{
				Source:  src.Position,
				Target:  tgt.Position,
				Enabled: true,
			})
		}
	}

	// Return the new genome
	return g
}

// WithPopulator sets the experiment's populator to a configured NEAT populator
func WithPopulator(cfg evo.Configurer) evo.Option {
	return func(e *evo.Experiment) (err error) {
		p := new(Populator)
		if err = cfg.Configure(p); err != nil {
			return
		}
		e.Populator = p
		return
	}
}
