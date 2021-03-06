package mutator

import (
	"math"
	"testing"

	"github.com/klokare/evo"
	"github.com/klokare/evo/internal/float"
)

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
							{Bias: 0.0, Neuron: evo.Input},
							{Bias: 1.0, Neuron: evo.Output},
						},
					},
				}

				// Mutate
				err := mut.Mutate(&g)
				if err != nil {
					t.Errorf("there should be no error. instead, %v", err)
					t.FailNow()
				}

				// Inputs should not be affected
				if g.Encoded.Nodes[0].Bias != 0.0 {
					t.Errorf("input node should not be affected")
				}

				// Record the new bias
				biases[i] = g.Encoded.Nodes[1].Bias

			}

			// Compare against expected
			m := float.Mean(biases)
			s := float.Stdev(biases)
			if math.Abs(m-test.Mean) > 0.1 {
				t.Errorf("incorrect mean bias. expected %f, actual %f", test.Mean, m)
			}
			if math.Abs(s-test.Stdev) > 0.1 {
				t.Errorf("incorrect standard deviation. expected %f, actual %f", test.Stdev, s)
			}
		})
	}
}
