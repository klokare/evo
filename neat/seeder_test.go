package neat

import (
	"math"
	"testing"

	"github.com/klokare/evo"
	"github.com/klokare/evo/internal/float"
	"github.com/klokare/evo/internal/mock"
	"github.com/klokare/evo/internal/test"
)

// Tests the population creation. Testing the specific genome creation is done below.
func TestSeederSeed(t *testing.T) {

	var cases = []struct {
		Desc                   string
		In, Out, Traits        int
		BiasPower, MaxBias     float64
		WeightPower, MaxWeight float64
		DisconnectRate         float64
		HasError               bool
	}{
		{
			Desc:     "negative number of inputs",
			In:       -1,
			HasError: true,
		},
		{
			Desc:     "zero inputs",
			In:       0,
			HasError: true,
		},
		{
			Desc:     "negative number of outputs",
			In:       1,
			Out:      -1,
			HasError: true,
		},
		{
			Desc:     "zero outputs",
			In:       1,
			Out:      0,
			HasError: true,
		},
		{
			Desc:           "normal population",
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

			// Create the seeder
			s := &Seeder{
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
			genomes, err := s.Seed()

			// Check error
			if !t.Run("error", test.Error(c.HasError, err)) {
				t.FailNow()
			}
			if c.HasError {
				return // Expected error so stop checking other things
			}

			// Create 10000 genomes to test randomness
			genomes = make([]evo.Genome, 0, 10000)
			for i := 0; i < 10000; i++ {
				g2, _ := s.Seed()
				genomes = append(genomes, g2...)
			}

			// The traits should show the right amount of randomness
			t.Run("traits", func(t *testing.T) {
				vals := make([]float64, 0, len(genomes[0].Traits)*len(genomes))
				for _, g := range genomes {
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
				vals := make([]float64, 0, len(genomes[0].Encoded.Nodes)*len(genomes))
				for _, g := range genomes {
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
				vals := make([]float64, 0, len(genomes[0].Encoded.Conns)*len(genomes))
				for _, g := range genomes {
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
				cnt := float64(len(genomes) * s.NumInputs * s.NumOutputs)
				expected := cnt * (1.0 - c.DisconnectRate)
				sum := 0.0
				for _, g := range genomes {
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

func TestSeederMinMax(t *testing.T) {

	// Create a population
	s := &Seeder{
		NumInputs:   2,
		NumOutputs:  2,
		NumTraits:   2,
		BiasPower:   9.0,
		MaxBias:     3.0,
		WeightPower: 8.0,
		MaxWeight:   2.0,
	}

	genomes := make([]evo.Genome, 0, 10000)
	for i := 0; i < 10000; i++ {
		g2, _ := s.Seed()
		genomes = append(genomes, g2...)
	}
	// Check min and max weights, bias, and traits
	t.Run("traits", func(t *testing.T) {
		vals := make([]float64, 0, len(genomes[0].Traits)*len(genomes))
		for _, g := range genomes {
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
		vals := make([]float64, 0, len(genomes[0].Encoded.Nodes)*len(genomes))
		for _, g := range genomes {
			for _, node := range g.Encoded.Nodes {
				if node.Neuron != evo.Input {
					vals = append(vals, node.Bias)
				}
			}
		}
		min := float.Min(vals...)
		max := float.Min(vals...)
		if min < -s.MaxBias {
			t.Errorf("incorrect min bias value: expected %f, actual: %f", -s.MaxBias, min)
		}
		if max > s.MaxBias {
			t.Errorf("incorrect max bias value: expected %f, actual: %f", s.MaxBias, max)
		}
	})

	t.Run("weights", func(t *testing.T) {
		vals := make([]float64, 0, len(genomes[0].Encoded.Conns)*len(genomes))
		for _, g := range genomes {
			for _, conn := range g.Encoded.Conns {
				vals = append(vals, conn.Weight)
			}
		}
		min := float.Min(vals...)
		max := float.Min(vals...)
		if min < -s.MaxWeight {
			t.Errorf("incorrect min weight value: expected %f, actual: %f", -s.MaxWeight, min)
		}
		if max > s.MaxWeight {
			t.Errorf("incorrect max weight value: expected %f, actual: %f", s.MaxWeight, max)
		}
	})
}

func TestWithSeeder(t *testing.T) {
	e := new(evo.Experiment)

	// Configurer has no error
	err := WithSeeder(&mock.Configurer{})(e)
	if err != nil {
		t.Errorf("error not expected, instead %v", err)
	}
	if _, ok := e.Seeder.(*Seeder); !ok {
		t.Errorf("seeder incorrectly set")
	}

	// Configurer has error
	err = WithSeeder(&mock.Configurer{HasError: true})(e)
	if err == nil {
		t.Errorf("error expected but not found")
	}
}
