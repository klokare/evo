package mutator

import (
	"testing"

	"github.com/klokare/evo"
	"github.com/klokare/evo/internal/float"
)

func TestTrait(t *testing.T) {

	var tests = []struct {
		Desc                    string
		MutateTraitProbability  float64
		ReplaceTraitProbability float64
	}{
		{
			Desc: "No probabilities so no change",
			MutateTraitProbability:  0.0,
			ReplaceTraitProbability: 0.0,
		},
		{
			Desc: "Always mutate but never replace",
			MutateTraitProbability:  1.0,
			ReplaceTraitProbability: 0.0,
		},
		{
			Desc: "Always mutate and always replace",
			MutateTraitProbability:  1.0,
			ReplaceTraitProbability: 1.0,
		},
	}

	for _, test := range tests {
		t.Run(test.Desc, func(t *testing.T) {

			// Create the mutator
			mut := &Trait{
				MutateTraitProbability:  test.MutateTraitProbability,
				ReplaceTraitProbability: test.ReplaceTraitProbability,
			}

			// Run N times to get a sample
			traits := make([]float64, 10000)
			g := evo.Genome{Traits: make([]float64, 1)}
			for i := 0; i < len(traits); i++ {

				// Randomly set the trait to an extreme
				g.Traits[0] = 0.5

				// Mutate
				err := mut.Mutate(&g)
				if err != nil {
					t.Errorf("there should be no error. instead, %v", err)
					t.FailNow()
				}

				// Record the new trait
				traits[i] = g.Traits[0]

			}

			// Compare against expected
			if test.MutateTraitProbability > 0.0 || test.ReplaceTraitProbability > 0.0 {
				if float.Variance(traits) == 0 {
					t.Error("variance should not be zero")
				}
				if float.Max(traits) > 1.0 {
					t.Errorf("invalid max trait: expected 1.0, actual %f", float.Max(traits))
				}
				if float.Min(traits) < 0.0 {
					t.Errorf("invalid min trait: expected 0.0, actual %f", float.Min(traits))
				}
			} else {
				if float.Variance(traits) != 0.0 {
					t.Error("variance should be zero")
				}
			}

		})
	}
}
