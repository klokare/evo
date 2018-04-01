package neat

import (
	"errors"

	"github.com/klokare/evo"
)

// Known errors
var (
	ErrMissingSeeder = errors.New("neat populator requires a seeder")
)

// Populator provides a population to the experiment
type Populator struct {
	PopulationSize int
	BiasPower      float64
	MaxBias        float64
	WeightPower    float64
	MaxWeight      float64
	evo.Seeder
}

// Populate creates a new population by creating randomised version of a seed genome.
func (p Populator) Populate() (pop evo.Population, err error) {

	// Check for errors
	if p.PopulationSize < 1 {
		err = ErrInvalidPopulationSize
		return
	} else if p.Seeder == nil {
		err = ErrMissingSeeder
		return
	}

	// Create the seed genome
	var seed evo.Genome
	if seed, err = p.Seeder.Seed(); err != nil {
		return
	}

	// Create the genomes
	var g evo.Genome
	pop.Genomes = make([]evo.Genome, p.PopulationSize)
	rng := evo.NewRandom()
	for i := 0; i < p.PopulationSize; i++ {

		// Clone the seed genome
		g = evo.Genome{
			ID: int64(i + 1),
			Encoded: evo.Substrate{
				Nodes: make([]evo.Node, len(seed.Encoded.Nodes)),
				Conns: make([]evo.Conn, len(seed.Encoded.Conns)),
			},
			Traits: make([]float64, len(seed.Traits)),
		}
		copy(g.Encoded.Nodes, seed.Encoded.Nodes)
		copy(g.Encoded.Conns, seed.Encoded.Conns)
		copy(g.Traits, seed.Traits)

		// Set the bias values for hidden and output nodes
		for i, n := range g.Encoded.Nodes {
			if n.Neuron != evo.Input {
				n.Bias = rng.NormFloat64() * p.BiasPower
			}
			if n.Bias > p.MaxBias {
				n.Bias = p.MaxBias
			} else if n.Bias < -p.MaxBias {
				n.Bias = -p.MaxBias
			}
			g.Encoded.Nodes[i] = n
		}

		// Set the weight values for connections
		for i, c := range g.Encoded.Conns {
			c.Weight = rng.NormFloat64() * p.WeightPower
			if c.Weight > p.MaxWeight {
				c.Weight = p.MaxWeight
			} else if c.Weight < -p.MaxWeight {
				c.Weight = -p.MaxWeight
			}
			g.Encoded.Conns[i] = c
		}

		// Save the genome
		pop.Genomes[i] = g
	}
	return
}
