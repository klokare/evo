package mutator

import (
	"math"
	"testing"

	"github.com/klokare/evo"
	"github.com/klokare/evo/internal/float"
)

func TestWeight(t *testing.T) {

	var tests = []struct {
		MutateWeightProbability  float64
		ReplaceWeightProbability float64
		WeightPower              float64
		MaxWeight                float64
		Mean                     float64
		Stdev                    float64
	}{
		{ // No probabilities so no change
			MutateWeightProbability:  0.0,
			ReplaceWeightProbability: 0.0,
			WeightPower:              2.5,
			MaxWeight:                8.0,
			Mean:                     1.0,
			Stdev:                    0.0,
		},
		{ // Always mutate but never replace
			MutateWeightProbability:  1.0,
			ReplaceWeightProbability: 0.0,
			WeightPower:              2.5,
			MaxWeight:                8.0,
			Mean:                     1.0,
			Stdev:                    2.5,
		},
		{ // Always mutate and always replace
			MutateWeightProbability:  1.0,
			ReplaceWeightProbability: 1.0,
			WeightPower:              2.5,
			MaxWeight:                8.0,
			Mean:                     0.0,
			Stdev:                    2.5,
		},
	}

	for _, test := range tests {
		t.Run("Mutate", func(t *testing.T) {

			// Create the mutator
			mut := &Weight{
				MutateWeightProbability:  test.MutateWeightProbability,
				ReplaceWeightProbability: test.ReplaceWeightProbability,
				WeightPower:              test.WeightPower,
				MaxWeight:                test.MaxWeight,
			}

			// Run N times to get a sample
			weights := make([]float64, 10000)
			for i := 0; i < len(weights); i++ {

				// Using the same genome
				g := evo.Genome{
					Encoded: evo.Substrate{
						Conns: []evo.Conn{
							{Weight: 1.0},
						},
					},
				}

				// Mutate
				err := mut.Mutate(&g)
				if err != nil {
					t.Errorf("there should be no error. instead, %v", err)
					t.FailNow()
				}

				// Record the new weight
				weights[i] = g.Encoded.Conns[0].Weight

			}

			// Compare against expected
			m := float.Mean(weights)
			s := float.Stdev(weights)
			if math.Abs(m-test.Mean) > 0.1 {
				t.Errorf("incorrect mean weight. expected %f, actual %f", test.Mean, m)
			}
			if math.Abs(s-test.Stdev) > 0.1 {
				t.Errorf("incorrect standard deviation. expected %f, actual %f", test.Stdev, s)
			}
		})
	}
}
