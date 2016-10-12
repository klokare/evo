package evo

import (
	"sort"

	"github.com/klokare/evo"
	"github.com/klokare/random"
)

func Seed(popsize int, inputs, outputs int) evo.Population {

	// Create the prototype
	g := evo.Genome{
		ID:        1,
		SpeciesID: 1,
		Encoded: evo.Substrate{
			Nodes: make([]evo.Node, 0, 1+inputs+outputs),
			Conns: make([]evo.Conn, 0, (1+inputs)*outputs),
		},
	}

	g.Encoded.Nodes = append(g.Encoded.Nodes, evo.Node{
		Position: evo.Position{Layer: 0.0, X: 0.0}, NeuronType: evo.Bias, ActivationType: evo.Direct,
	})

	for i := 0; i < inputs; i++ {
		x := float64(i+1) / float64(inputs)
		g.Encoded.Nodes = append(g.Encoded.Nodes, evo.Node{
			Position: evo.Position{Layer: 0.0, X: x}, NeuronType: evo.Input, ActivationType: evo.Direct,
		})
	}

	for i := 0; i < outputs; i++ {
		x := 0.5
		if outputs > 1 {
			x = float64(i) / float64(outputs-1)
		}
		g.Encoded.Nodes = append(g.Encoded.Nodes, evo.Node{
			Position: evo.Position{Layer: 1.0, X: x}, NeuronType: evo.Output, ActivationType: evo.SteepenedSigmoid,
		})
	}

	rng := random.New()
	for i := 0; i < 1+inputs; i++ {
		for j := 0; j < outputs; j++ {
			g.Encoded.Conns = append(g.Encoded.Conns, evo.Conn{
				Source:  g.Encoded.Nodes[i].Position,
				Target:  g.Encoded.Nodes[j+1+inputs].Position,
				Weight:  rng.NormFloat64(),
				Enabled: true,
			})
		}
	}

	sort.Sort(g.Encoded.Nodes)
	sort.Sort(g.Encoded.Conns)

	// Build the population
	p := evo.Population{
		Generation: 1,
		Species: []evo.Species{
			{ID: 1, Example: g.Encoded},
		},
		Genomes: make([]evo.Genome, 0, popsize),
	}
	p.Genomes = append(p.Genomes, g)

	for i := 1; i < popsize; i++ {
		g2 := evo.Genome{
			ID:        i + 1,
			SpeciesID: 1,
			Encoded: evo.Substrate{
				Nodes: make([]evo.Node, len(g.Encoded.Nodes)),
				Conns: make([]evo.Conn, len(g.Encoded.Conns)),
			},
		}
		copy(g2.Encoded.Nodes, g.Encoded.Nodes)
		copy(g2.Encoded.Conns, g.Encoded.Conns)
		for j := 0; j < len(g2.Encoded.Conns); j++ {
			g2.Encoded.Conns[j].Weight = rng.NormFloat64()
		}
		p.Genomes = append(p.Genomes, g2)
	}

	return p
}
