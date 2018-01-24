package neat

import (
	"context"
	"math"
	"testing"

	"github.com/klokare/evo"
	"github.com/klokare/evo/internal/float"
	"github.com/klokare/evo/internal/mock"
	"github.com/klokare/evo/internal/test"
)

// Tests the population creation. Testing the specific genome creation is done below.
func TestPopulatorPopulate(t *testing.T) {

	var cases = []struct {
		Desc                   string
		PopSize                int
		In, Out, Traits        int
		BiasPower, MaxBias     float64
		WeightPower, MaxWeight float64
		DisconnectRate         float64
		HasError               bool
	}{
		{
			Desc:     "negative population size",
			PopSize:  -1,
			HasError: true,
		},
		{
			Desc:     "zero population size",
			PopSize:  0,
			HasError: true,
		},
		{
			Desc:     "negative number of inputs",
			PopSize:  1,
			In:       -1,
			HasError: true,
		},
		{
			Desc:     "zero inputs",
			PopSize:  1,
			In:       0,
			HasError: true,
		},
		{
			Desc:     "negative number of outputs",
			PopSize:  1,
			In:       1,
			Out:      -1,
			HasError: true,
		},
		{
			Desc:     "zero outputs",
			PopSize:  1,
			In:       1,
			Out:      0,
			HasError: true,
		},
		// No mutators is weird but technically not an error so skipping that case
		{
			Desc:           "normal population",
			PopSize:        10000, // With a sufficiently large population we can also check randomness
			In:             2,
			Out:            2,
			Traits:         2,
			BiasPower:      3.0,
			MaxBias:        99.0, // Will test max bias and weight separately
			WeightPower:    2.0,
			MaxWeight:      99.0,
			DisconnectRate: 0.5,
			HasError:       false,
		},
	}

	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {

			// Create the populator
			p := &Populator{
				PopulationSize: c.PopSize,
				NumInputs:      c.In,
				NumOutputs:     c.Out,
				NumTraits:      c.Traits,
				BiasPower:      c.BiasPower,
				MaxBias:        c.MaxBias,
				WeightPower:    c.WeightPower,
				MaxWeight:      c.MaxWeight,
				DisconnectRate: c.DisconnectRate,
			}

			// Create the population
			pop, err := p.Populate(context.Background())

			// Check error
			if !t.Run("error", test.Error(c.HasError, err)) {
				t.FailNow()
			}
			if c.HasError {
				return // Expected error so stop checking other things
			}

			// There should be the correct number of genomes
			if len(pop.Genomes) != c.PopSize {
				t.Errorf("incorrect population size: expected %d, actual %d", c.PopSize, len(pop.Genomes))
			}

			// The traits should show the right amount of randomness
			t.Run("traits", func(t *testing.T) {
				vals := make([]float64, 0, len(pop.Genomes[0].Traits)*len(pop.Genomes))
				for _, g := range pop.Genomes {
					for _, x := range g.Traits {
						vals = append(vals, x)
					}
				}
				avg := float.Mean(vals...)
				if math.Abs(0.5-avg) > 0.1 {
					t.Errorf("incorrect average random trait value: expected 0.5, actual %f", avg)
				}
			})

			// The bias should show the right amount of randomness
			t.Run("bias", func(t *testing.T) {
				vals := make([]float64, 0, len(pop.Genomes[0].Encoded.Nodes)*len(pop.Genomes))
				for _, g := range pop.Genomes {
					for _, node := range g.Encoded.Nodes {
						if node.Neuron != evo.Input {
							vals = append(vals, node.Bias)
						}
					}
				}
				avg := float.Mean(vals...)
				stdev := float.Stdev(vals...)
				if math.Abs(0.0-avg) > 0.1 {
					t.Errorf("incorrect average random bias value: expected 0.0, actual %f", avg)
				}
				if math.Abs(c.BiasPower-stdev) > 0.1 {
					t.Errorf("incorrect stdev random bias value: expected %f, actual %f", c.BiasPower, stdev)
				}
			})

			// The weights should show the right amount of randomness
			t.Run("weights", func(t *testing.T) {
				vals := make([]float64, 0, len(pop.Genomes[0].Encoded.Conns)*len(pop.Genomes))
				for _, g := range pop.Genomes {
					for _, conn := range g.Encoded.Conns {
						vals = append(vals, conn.Weight)
					}
				}
				min := float.Min(vals...)
				max := float.Min(vals...)
				if min < -c.MaxWeight {
					t.Errorf("incorrect min weight value: expected %f, actual: %f", -c.MaxWeight, min)
				}
				if max > c.MaxWeight {
					t.Errorf("incorrect max weight value: expected %f, actual: %f", c.MaxWeight, max)
				}
			})

			// The disconnected rate should be correct
			t.Run("disconnected", func(t *testing.T) {
				cnt := float64(len(pop.Genomes) * p.NumInputs * p.NumOutputs)
				expected := cnt * (1.0 - c.DisconnectRate)
				sum := 0.0
				for _, g := range pop.Genomes {
					for _ = range g.Encoded.Conns {
						sum += 1.0
					}
				}
				actual := sum / cnt
				if actual-expected > 0 {
					t.Errorf("incorrect connectedness: expected %f, actual %f", expected, actual)
				}
			})
		})
	}
}

func TestPopulatorSeed(t *testing.T) {

	// Create the seed genome
	rng := evo.NewRandom()
	act := seed(rng, 3, 2, 2, 3.0, 2.0, evo.Sigmoid)
	acts := act.Encoded

	// The seed genome should match.
	exp := evo.Genome{
		Traits: make([]float64, 2),
		Encoded: evo.Substrate{
			Nodes: []evo.Node{
				{Position: evo.Position{Layer: 0.0, X: 0.0}, Neuron: evo.Input, Activation: evo.Direct},
				{Position: evo.Position{Layer: 0.0, X: 0.5}, Neuron: evo.Input, Activation: evo.Direct},
				{Position: evo.Position{Layer: 0.0, X: 1.0}, Neuron: evo.Input, Activation: evo.Direct},
				{Position: evo.Position{Layer: 1.0, X: 0.0}, Neuron: evo.Output, Activation: evo.Sigmoid},
				{Position: evo.Position{Layer: 1.0, X: 1.0}, Neuron: evo.Output, Activation: evo.Sigmoid},
			},
			Conns: []evo.Conn{
				{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.0}, Enabled: true},
				{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 1.0}, Enabled: true},
				{Source: evo.Position{Layer: 0.0, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 0.0}, Enabled: true},
				{Source: evo.Position{Layer: 0.0, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 1.0}, Enabled: true},
				{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.0}, Enabled: true},
				{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 1.0}, Enabled: true},
			},
		},
	}
	exps := exp.Encoded
	if len(exps.Nodes) != len(acts.Nodes) {
		t.Errorf("incorrect number of nodes: expected %d, actual %d", len(exps.Nodes), len(acts.Nodes))
	} else {
		for i := 0; i < len(exps.Nodes); i++ {
			en := exps.Nodes[i]
			an := acts.Nodes[i]
			if en.Position.Compare(an.Position) != 0 {
				t.Errorf("incorrect position for node %d: expected %v, actual %v", i, en.Position, an.Position)
			}
			if en.Neuron != an.Neuron {
				t.Errorf("incorrect neuron type for node %d: expected %s, actual %s", i, en.Neuron, an.Neuron)
			}
			if en.Activation != an.Activation {
				t.Errorf("incorrect activation type for node %d: expected %s, actual %s", i, en.Activation, an.Activation)
			}
		}
	}
	if len(exps.Conns) != len(acts.Conns) {
		t.Errorf("incorrect number of connections: expected %d, actual %d", len(exps.Conns), len(acts.Conns))
	} else {
		for i := 0; i < len(exps.Conns); i++ {
			ec := exps.Conns[i]
			ac := exps.Conns[i]
			if ec.Compare(ac) != 0 {
				t.Errorf("incorrect source, target for connection %d: expected %v->%v, actual %v->%v", i, ec.Source, ec.Target, ac.Source, ac.Target)
			}
			if ec.Enabled != ac.Enabled {
				t.Errorf("incorrect enabled for connection %d: expected %v, actual %v", i, ec.Enabled, ac.Enabled)
			}
		}
	}

	// Ensure traits. Just testing here that there are traits. The randomness tested below.
	if len(exp.Traits) != len(act.Traits) {
		t.Errorf("incorrect number of traits: expected %d, actual: %d", len(exp.Traits), len(act.Traits))
	}

}

func TestPopulatorMinMax(t *testing.T) {

	// Create a population
	p := &Populator{
		PopulationSize: 10000,
		NumInputs:      2,
		NumOutputs:     2,
		NumTraits:      2,
		BiasPower:      9.0,
		MaxBias:        3.0,
		WeightPower:    8.0,
		MaxWeight:      2.0,
	}

	pop, _ := p.Populate(context.Background())

	// Check min and max weights, bias, and traits
	t.Run("traits", func(t *testing.T) {
		vals := make([]float64, 0, len(pop.Genomes[0].Traits)*len(pop.Genomes))
		for _, g := range pop.Genomes {
			for _, x := range g.Traits {
				vals = append(vals, x)
			}
		}
		min := float.Min(vals...)
		max := float.Min(vals...)
		if min < 0.0 {
			t.Errorf("incorrect min random trait value: expected 0.0, actual %f", min)
		}
		if max > 1.0 {
			t.Errorf("incorrect max random trait value: expected 1.0, actual %f", max)
		}
	})

	t.Run("bias", func(t *testing.T) {
		vals := make([]float64, 0, len(pop.Genomes[0].Encoded.Nodes)*len(pop.Genomes))
		for _, g := range pop.Genomes {
			for _, node := range g.Encoded.Nodes {
				if node.Neuron != evo.Input {
					vals = append(vals, node.Bias)
				}
			}
		}
		min := float.Min(vals...)
		max := float.Min(vals...)
		if min < -p.MaxBias {
			t.Errorf("incorrect min bias value: expected %f, actual: %f", -p.MaxBias, min)
		}
		if max > p.MaxBias {
			t.Errorf("incorrect max bias value: expected %f, actual: %f", p.MaxBias, max)
		}
	})

	t.Run("weights", func(t *testing.T) {
		vals := make([]float64, 0, len(pop.Genomes[0].Encoded.Conns)*len(pop.Genomes))
		for _, g := range pop.Genomes {
			for _, conn := range g.Encoded.Conns {
				vals = append(vals, conn.Weight)
			}
		}
		min := float.Min(vals...)
		max := float.Min(vals...)
		if min < -p.MaxWeight {
			t.Errorf("incorrect min weight value: expected %f, actual: %f", -p.MaxWeight, min)
		}
		if max > p.MaxWeight {
			t.Errorf("incorrect max weight value: expected %f, actual: %f", p.MaxWeight, max)
		}
	})
}

func TestWithPopulator(t *testing.T) {
	e := new(evo.Experiment)

	// Configurer has no error
	err := WithPopulator(&mock.Configurer{})(e)
	if err != nil {
		t.Errorf("error not expected, instead %v", err)
	}
	if _, ok := e.Populator.(*Populator); !ok {
		t.Errorf("populator incorrectly set")
	}

	// Configurer has error
	err = WithPopulator(&mock.Configurer{HasError: true})(e)
	if err == nil {
		t.Errorf("error expected but not found")
	}
}
