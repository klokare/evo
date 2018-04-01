package neat

import (
	"math"
	"testing"

	"github.com/klokare/evo"
	"github.com/klokare/evo/internal/float"
	"github.com/klokare/evo/internal/mock"
)

func TestPopulatorPopulate(t *testing.T) {
	var cases = []struct {
		Desc     string
		HasError bool
		PopSize  int
		Seeder   evo.Seeder
	}{
		{
			Desc:     "zero pop size should error",
			HasError: true,
			PopSize:  0,
			Seeder:   &mock.Seeder{},
		},
		{
			Desc:     "negative pop size should error",
			HasError: true,
			PopSize:  -1,
			Seeder:   &mock.Seeder{},
		},
		{
			Desc:     "missing seeder should cause failure",
			HasError: true,
			PopSize:  1,
			Seeder:   nil,
		},
		{
			Desc:     "seeder error should cause failure",
			HasError: true,
			PopSize:  1,
			Seeder:   &mock.Seeder{HasError: true},
		},
		{
			Desc:     "no errors",
			HasError: false,
			PopSize:  5,
			Seeder:   &mock.Seeder{},
		},
	}

	seed, _ := new(mock.Seeder).Seed()

	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {

			// Create the populator
			z := &Populator{
				PopulationSize: c.PopSize,
				Seeder:         c.Seeder,
			}

			// Create the population
			pop, err := z.Populate()

			// Test for error
			t.Run("error", mock.Error(c.HasError, err))
			if c.HasError {
				return
			}

			// There should be the correct number of genomes
			if len(pop.Genomes) != c.PopSize {
				t.Errorf("inccorect number of genomes: expected %d, actual %d", len(pop.Genomes), c.PopSize)
			}

			// Each genome should have the same structure as the seed genome and a unquie ID
			// NOTE: bias, weight, and trait values are tested below
			for i, g := range pop.Genomes {

				// Check for unique ID
				for j := 0; j < i-1; j++ {
					if pop.Genomes[i].ID == pop.Genomes[j].ID {
						t.Errorf("id for genome %d is not unique", i)
						break
					}
				}

				// Same traits
				if len(g.Traits) != len(seed.Traits) {
					t.Errorf("incorrect number of traits for genome %d: expected %d, actual %d", i, len(g.Traits), len(seed.Traits))
				}

				// Same nodes
				if len(g.Encoded.Nodes) != len(seed.Encoded.Nodes) {
					t.Errorf("incorrect number of nodes for genome %d: expected %d, actual %d", i, len(g.Encoded.Nodes), len(seed.Encoded.Nodes))
				} else {
					for j := 0; j < len(seed.Encoded.Nodes); j++ {
						en := seed.Encoded.Nodes[j]
						an := g.Encoded.Nodes[j]
						if en.Position.Compare(an.Position) != 0 {
							t.Errorf("incorrect position for node %d of genome %d: expected %v, actual %v", j, i, en.Position, an.Position)
						}
						if en.Neuron != an.Neuron {
							t.Errorf("incorrect neuron for node %d of genome %d: expected %s, actual %s", j, i, en.Neuron, an.Neuron)
						}
						if en.Activation != an.Activation {
							t.Errorf("incorrect activation for node %d of genome %d: expected %s, actual %s", j, i, en.Neuron, an.Neuron)
						}
					}
				}
				// Same connections
				if len(g.Encoded.Conns) != len(seed.Encoded.Conns) {
					t.Errorf("incorrect number of conns for genome %d: expected %d, actual %d", i, len(g.Encoded.Conns), len(seed.Encoded.Conns))
				} else {
					for j := 0; j < len(seed.Encoded.Conns); j++ {
						ec := seed.Encoded.Conns[j]
						ac := g.Encoded.Conns[j]
						if ec.Source.Compare(ac.Source) != 0 {
							t.Errorf("incorrect source for conn %d of genome %d: expected %v, actual %v", j, i, ec.Source, ac.Source)
						}
						if ec.Target.Compare(ac.Target) != 0 {
							t.Errorf("incorrect target for conn %d of genome %d: expected %v, actual %v", j, i, ec.Target, ac.Target)
						}
						if ec.Enabled != ac.Enabled {
							t.Errorf("incorrect enabled for node %d of genome %d: expected %t, actual %t", j, i, ec.Enabled, ac.Enabled)
						}
					}
				}
			}
		})
	}
}

func TestPopulatorRandomness(t *testing.T) {

	// Create a populator, leave the max values really high
	p := Populator{
		PopulationSize: 10000,
		Seeder:         &mock.Seeder{},
		BiasPower:      3.0,
		MaxBias:        math.MaxFloat64,
		WeightPower:    4.0,
		MaxWeight:      math.MaxFloat64,
	}

	// Create the population
	pop, _ := p.Populate()

	// Gather the values
	bias := make([]float64, 0, len(pop.Genomes)*len(pop.Genomes[0].Encoded.Nodes))
	weights := make([]float64, 0, len(pop.Genomes)*len(pop.Genomes[0].Encoded.Conns))
	for _, g := range pop.Genomes {
		for _, n := range g.Encoded.Nodes {
			if n.Neuron != evo.Input {
				bias = append(bias, n.Bias)
			}
		}
		for _, c := range g.Encoded.Conns {
			weights = append(weights, c.Weight)
		}
	}

	// Test for mean and stdev
	mb := float.Mean(bias)
	if math.Abs(mb) > 0.1 {
		t.Errorf("incorrect mean bias: expected 0.0, actual: %f", mb)
	}
	sb := float.Stdev(bias)
	if math.Abs(sb-p.BiasPower) > 0.1 {
		t.Errorf("incorrect mean stdev: expected: %f, actual %f", p.BiasPower, sb)
	}

	mw := float.Mean(weights)
	if math.Abs(mw) > 0.1 {
		t.Errorf("incorrect mean weights: expected 0.0, actual: %f", mw)
	}
	sw := float.Stdev(weights)
	if math.Abs(sw-p.WeightPower) > 0.1 {
		t.Errorf("incorrect mean weights: expected: %f, actual %f", p.WeightPower, sw)
	}
}

func TestPopulatorMinMax(t *testing.T) {

	// Create a populator, this time set the powers high but the max low
	p := Populator{
		PopulationSize: 10000,
		Seeder:         &mock.Seeder{},
		BiasPower:      30.0,
		MaxBias:        3.0,
		WeightPower:    40.0,
		MaxWeight:      4.0,
	}

	// Create the population
	pop, _ := p.Populate()

	// Gather the values
	bias := make([]float64, 0, len(pop.Genomes)*len(pop.Genomes[0].Encoded.Nodes))
	weights := make([]float64, 0, len(pop.Genomes)*len(pop.Genomes[0].Encoded.Conns))
	for _, g := range pop.Genomes {
		for _, n := range g.Encoded.Nodes {
			if n.Neuron != evo.Input {
				bias = append(bias, n.Bias)
			}
		}
		for _, c := range g.Encoded.Conns {
			weights = append(weights, c.Weight)
		}
	}

	// Test for min and max
	min := float.Min(bias)
	if min < -p.MaxBias {
		t.Errorf("incorrect min bias: expected %f, actual: %f", -p.MaxBias, min)
	}
	max := float.Max(bias)
	if max > p.MaxBias {
		t.Errorf("incorrect max bias: expected %f, actual: %f", p.MaxBias, max)
	}

	min = float.Min(weights)
	if min < -p.MaxWeight {
		t.Errorf("incorrect min weights: expected %f, actual: %f", -p.MaxWeight, min)
	}
	max = float.Max(weights)
	if max > p.MaxWeight {
		t.Errorf("incorrect max weights: expected %f, actual: %f", p.MaxWeight, max)
	}

}
