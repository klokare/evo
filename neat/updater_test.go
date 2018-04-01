package neat

import (
	"sort"
	"testing"

	"github.com/klokare/evo"
	"github.com/klokare/evo/internal/mock"
)

func TestUpdaterUpdate(t *testing.T) {
	var cases = []struct {
		Desc      string
		HasError  bool
		DecayRate float64
		evo.Comparison
		Actual, Expected evo.Population
		Results          []evo.Result
	}{
		{
			Desc:     "missing compare should cause error",
			HasError: true,
		},
		{
			Desc:       "negative decay rate should cause error",
			HasError:   true,
			Comparison: evo.ByFitness,
			DecayRate:  -1.0,
		},
		{
			Desc:       "no change in champion",
			HasError:   false,
			Comparison: evo.ByFitness,
			DecayRate:  0.5,
			Actual: evo.Population{
				Genomes: []evo.Genome{
					{ID: 1, SpeciesID: 10},
					{ID: 2, SpeciesID: 10},
					{ID: 3, SpeciesID: 20},
					{ID: 4, SpeciesID: 20},
				},
				Species: []evo.Species{
					{ID: 10, Champion: 2, Decay: 0.0},
					{ID: 20, Champion: 4, Decay: 0.8},
				},
			},
			Expected: evo.Population{
				Genomes: []evo.Genome{
					{ID: 1, Fitness: 1.0, Novelty: 1.1, Solved: false},
					{ID: 2, Fitness: 2.0, Novelty: 2.2, Solved: false},
					{ID: 3, Fitness: 3.0, Novelty: 3.3, Solved: false},
					{ID: 4, Fitness: 4.0, Novelty: 4.4, Solved: true},
				},
				Species: []evo.Species{
					{ID: 10, Champion: 2, Decay: 0.5},
					{ID: 20, Champion: 4, Decay: 1.0}, // Max is 1.0
				},
			},
			Results: []evo.Result{
				{ID: 1, Fitness: 1.0, Novelty: 1.1, Solved: false},
				{ID: 2, Fitness: 2.0, Novelty: 2.2, Solved: false},
				{ID: 3, Fitness: 3.0, Novelty: 3.3, Solved: false},
				{ID: 4, Fitness: 4.0, Novelty: 4.4, Solved: true},
			},
		},
		{
			Desc:       "change of champion",
			HasError:   false,
			Comparison: evo.ByFitness,
			DecayRate:  0.5,
			Actual: evo.Population{
				Genomes: []evo.Genome{
					{ID: 3, SpeciesID: 20},
					{ID: 1, SpeciesID: 10},
					{ID: 2, SpeciesID: 10},
					{ID: 4, SpeciesID: 20},
				},
				Species: []evo.Species{
					{ID: 10, Champion: 2, Decay: 0.0},
					{ID: 20, Champion: 4, Decay: 0.8},
				},
			},
			Expected: evo.Population{
				Genomes: []evo.Genome{
					{ID: 3, Fitness: 5.0, Novelty: 3.3, Solved: false},
					{ID: 1, Fitness: 1.0, Novelty: 1.1, Solved: false},
					{ID: 2, Fitness: 2.0, Novelty: 2.2, Solved: false},
					{ID: 4, Fitness: 4.0, Novelty: 4.4, Solved: true},
				},
				Species: []evo.Species{
					{ID: 10, Champion: 2, Decay: 0.5},
					{ID: 20, Champion: 3, Decay: 0.0}, // Resets with new champion
				},
			},
			Results: []evo.Result{
				{ID: 1, Fitness: 1.0, Novelty: 1.1, Solved: false},
				{ID: 2, Fitness: 2.0, Novelty: 2.2, Solved: false},
				{ID: 3, Fitness: 5.0, Novelty: 3.3, Solved: false},
				{ID: 4, Fitness: 4.0, Novelty: 4.4, Solved: true},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {

			// Create the new updater
			u := &Updater{
				SpeciesDecayRate: c.DecayRate,
				Comparison:       c.Comparison,
			}

			// Update the population
			err := u.Update(&c.Actual, c.Results)

			// Test for error
			t.Run("error", mock.Error(c.HasError, err))
			if c.HasError {
				return
			}

			// Test the genomes
			sort.Slice(c.Actual.Genomes, func(i, j int) bool { return c.Actual.Genomes[i].ID < c.Actual.Genomes[j].ID })
			sort.Slice(c.Expected.Genomes, func(i, j int) bool { return c.Expected.Genomes[i].ID < c.Expected.Genomes[j].ID })
			if len(c.Expected.Genomes) != len(c.Actual.Genomes) {
				t.Errorf("incorrect number of genomes: expected %d, actual %d", len(c.Expected.Genomes), len(c.Actual.Genomes))
			} else {
				for i := 0; i < len(c.Expected.Genomes); i++ {
					eg := c.Expected.Genomes[i]
					ag := c.Actual.Genomes[i]
					if eg.Fitness != ag.Fitness {
						t.Errorf("incorrect fitness for genome %d: expected %f, actual %f", i, eg.Fitness, ag.Fitness)
					}
					if eg.Novelty != ag.Novelty {
						t.Errorf("incorrect novelty for genome %d: expected %f, actual %f", i, eg.Novelty, ag.Novelty)
					}
					if eg.Solved != ag.Solved {
						t.Errorf("incorrect solved for genome %d: expected %t, actual %t", i, eg.Solved, ag.Solved)
					}
				}
			}
			// Test the species
			sort.Slice(c.Actual.Species, func(i, j int) bool { return c.Actual.Species[i].ID < c.Actual.Species[j].ID })
			sort.Slice(c.Expected.Species, func(i, j int) bool { return c.Expected.Species[i].ID < c.Expected.Species[j].ID })
			if len(c.Expected.Species) != len(c.Actual.Species) {
				t.Errorf("incorrect number of genomes: expected %d, actual %d", len(c.Expected.Species), len(c.Actual.Species))
			} else {
				for i := 0; i < len(c.Expected.Species); i++ {
					es := c.Expected.Species[i]
					as := c.Actual.Species[i]
					if es.Decay != as.Decay {
						t.Errorf("incorrect decay for species %d: expected %f, actual %f", i, es.Decay, as.Decay)
					}
					if es.Champion != as.Champion {
						t.Errorf("incorrect champion for species %d: expected %d, actual %d", i, es.Champion, as.Champion)
					}
				}
			}
		})
	}
}
