package mutator

import (
	"context"
	"math"
	"testing"

	"github.com/klokare/evo"
	"github.com/klokare/evo/internal/float"
	"github.com/klokare/evo/internal/mock"
)

func TestWithBias(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		e := &evo.Experiment{Mutators: make([]evo.Mutator, 0, 1)}
		cfg := &mock.Configurer{HasError: true}
		err := WithBias(cfg)(e)
		if err == nil {
			t.Errorf("error expected")
		}
	})
	t.Run("enabled", func(t *testing.T) {
		e := &evo.Experiment{Mutators: make([]evo.Mutator, 0, 1)}
		cfg := &mock.Configurer{MutateBiasProbability: 1.0}
		err := WithBias(cfg)(e)
		if err != nil {
			t.Errorf("error unxpected: %v", err)
		}
		if len(e.Mutators) == 0 {
			t.Errorf("incorrect number of mutators: expected 1, actual 0")
		} else if _, ok := e.Mutators[0].(*Bias); !ok {
			t.Errorf("mutator is not a bias")
		}
	})
	t.Run("not enabled", func(t *testing.T) {
		e := &evo.Experiment{Mutators: make([]evo.Mutator, 0, 1)}
		cfg := &mock.Configurer{MutateBiasProbability: 0.0}
		err := WithBias(cfg)(e)
		if err != nil {
			t.Errorf("error unxpected: %v", err)
		}
		if len(e.Mutators) > 0 {
			t.Errorf("incorrect number of mutators: expected 0, actual: %d", len(e.Mutators))
		}
	})
}

func TestBias(t *testing.T) {

	var tests = []struct {
		MutateBiasProbability  float64
		ReplaceBiasProbability float64
		BiasPower              float64
		MaxBias                float64
		Mean                   float64
		Stdev                  float64
	}{
		{ // No probabilities so no change
			MutateBiasProbability:  0.0,
			ReplaceBiasProbability: 0.0,
			BiasPower:              2.5,
			MaxBias:                8.0,
			Mean:                   1.0,
			Stdev:                  0.0,
		},
		{ // Always mutate but never replace
			MutateBiasProbability:  1.0,
			ReplaceBiasProbability: 0.0,
			BiasPower:              2.5,
			MaxBias:                8.0,
			Mean:                   1.0,
			Stdev:                  1.0,
		},
		{ // Always mutate and always replace
			MutateBiasProbability:  1.0,
			ReplaceBiasProbability: 1.0,
			BiasPower:              2.5,
			MaxBias:                8.0,
			Mean:                   0.0,
			Stdev:                  2.5,
		},
	}

	for _, test := range tests {
		t.Run("Mutate", func(t *testing.T) {

			// Create the mutator
			mut := &Bias{
				MutateBiasProbability:  test.MutateBiasProbability,
				ReplaceBiasProbability: test.ReplaceBiasProbability,
				BiasPower:              test.BiasPower,
				MaxBias:                test.MaxBias,
			}

			// Run N times to get a sample
			biases := make([]float64, 10000)
			for i := 0; i < len(biases); i++ {

				// Using the same genome
				g := evo.Genome{
					Encoded: evo.Substrate{
						Nodes: []evo.Node{
							{Bias: 1.0},
						},
					},
				}

				// Mutate
				err := mut.Mutate(context.Background(), &g)
				if err != nil {
					t.Errorf("there should be no error. instead, %v", err)
					t.FailNow()
				}

				// Record the new weight
				biases[i] = g.Encoded.Nodes[0].Bias

			}

			// Compare against expected
			m := float.Mean(biases...)
			s := float.Stdev(biases...)
			if math.Abs(m-test.Mean) > 0.1 {
				t.Errorf("incorrect mean bias. expected %f, actual %f", test.Mean, m)
			}
			if math.Abs(s-test.Stdev) > 0.1 {
				t.Errorf("incorrect standard deviation. expected %f, actual %f", test.Stdev, s)
			}
		})
	}
}
