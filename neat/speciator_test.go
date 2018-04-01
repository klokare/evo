package neat

import (
	"errors"
	"math"
	"testing"

	"github.com/klokare/evo"
	"github.com/klokare/evo/internal/mock"
)

func TestSpeciatorSpeciate(t *testing.T) {

	// Test cases
	var cases = []struct {
		Desc      string
		Threshold float64
		Actual    evo.Population
		Expected  evo.Population
		Distancer Distancer
		HasError  bool
	}{
		{
			Desc:     "no distancer",
			HasError: true,
		},
		{
			Desc:      "negative threshold",
			Distancer: MockDistancer{},
			Threshold: -1.0,
			HasError:  true,
		},
		{
			Desc:      "distancer has error",
			Distancer: MockDistancer{HasError: true},
			Threshold: 1.0,
			HasError:  true,
			Actual: evo.Population{
				Genomes: []evo.Genome{{ID: 1}},
				Species: []evo.Species{{ID: 10}},
			},
		},
		{
			Desc:      "empty population", // odd but legal
			Distancer: MockDistancer{},
			Threshold: 1.0,
			HasError:  false,
		},
		{
			Desc:      "genomes fit in existing species",
			Distancer: MockDistancer{},
			Threshold: 1.0,
			HasError:  false,
			Actual: evo.Population{
				Genomes: []evo.Genome{
					{ID: 1, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}}, // 2 Nodes
					{ID: 2, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}}, // 2 Nodes
				},
				Species: []evo.Species{
					{ID: 10, Example: evo.Genome{Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}}},
				},
			},
			Expected: evo.Population{
				Genomes: []evo.Genome{
					{ID: 1, SpeciesID: 10, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}}, // 2 Nodes
					{ID: 2, SpeciesID: 10, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}}, // 2 Nodes
				},
				Species: []evo.Species{
					{ID: 10, Example: evo.Genome{Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}}},
				},
			},
		},
		{
			Desc:      "new species required",
			Distancer: MockDistancer{},
			Threshold: 1.0,
			HasError:  false,
			Actual: evo.Population{
				Genomes: []evo.Genome{
					{ID: 1, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}},     // 2 Nodes
					{ID: 2, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}, {}}}}, // 3 Nodes
				},
				Species: []evo.Species{
					{ID: 10, Example: evo.Genome{Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}}},
				},
			},
			Expected: evo.Population{
				Genomes: []evo.Genome{
					{ID: 1, SpeciesID: 10, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}},      // 2 Nodes
					{ID: 2, SpeciesID: 101, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}, {}}}}, // 2 Nodes
				},
				Species: []evo.Species{
					{ID: 10, Example: evo.Genome{Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}}},
					{ID: 101, Example: evo.Genome{Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}, {}}}}},
				},
			},
		},
		{
			Desc:      "species discarded",
			Distancer: MockDistancer{},
			Threshold: 1.0,
			HasError:  false,
			Actual: evo.Population{
				Genomes: []evo.Genome{
					{ID: 1, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}}, // 2 Nodes
					{ID: 2, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}}, // 2 Nodes
				},
				Species: []evo.Species{
					{ID: 10, Example: evo.Genome{Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}}},
					{ID: 11, Example: evo.Genome{Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}, {}}}}},
				},
			},
			Expected: evo.Population{
				Genomes: []evo.Genome{
					{ID: 1, SpeciesID: 10, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}}, // 2 Nodes
					{ID: 2, SpeciesID: 10, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}}, // 2 Nodes
				},
				Species: []evo.Species{
					{ID: 10, Example: evo.Genome{Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}}},
				},
			},
		},
		{
			Desc:      "previous species assignment, species still exists",
			Distancer: MockDistancer{},
			Threshold: 1.0,
			HasError:  false,
			Actual: evo.Population{
				Genomes: []evo.Genome{
					{ID: 1, SpeciesID: 25, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}}, // 2 Nodes
					{ID: 2, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}},                // 2 Nodes
				},
				Species: []evo.Species{
					{ID: 10, Example: evo.Genome{Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}}},
					{ID: 25, Example: evo.Genome{Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}}}, // Same structure so assignment should go to 10 unless previously assigned
				},
			},
			Expected: evo.Population{
				Genomes: []evo.Genome{
					{ID: 1, SpeciesID: 25, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}}, // 2 Nodes
					{ID: 2, SpeciesID: 10, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}}, // 2 Nodes
				},
				Species: []evo.Species{
					{ID: 10, Example: evo.Genome{Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}}},
					{ID: 25, Example: evo.Genome{Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}}}, // Same structure so assignment should go to 10 unless previously assigned
				},
			},
		},
		{
			Desc:      "previous species assignment, species no longer exists",
			Distancer: MockDistancer{},
			Threshold: 1.0,
			HasError:  false,
			Actual: evo.Population{
				Genomes: []evo.Genome{
					{ID: 1, SpeciesID: 25, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}}, // 2 Nodes
					{ID: 2, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}},                // 2 Nodes
				},
				Species: []evo.Species{
					{ID: 10, Example: evo.Genome{Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}}},
				},
			},
			Expected: evo.Population{
				Genomes: []evo.Genome{
					{ID: 1, SpeciesID: 10, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}}, // 2 Nodes
					{ID: 2, SpeciesID: 10, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}}, // 2 Nodes
				},
				Species: []evo.Species{
					{ID: 10, Example: evo.Genome{Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}}},
				},
			},
		},
	}

	// Run the tests
	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {

			// Create the speciator
			s := &Speciator{
				CompatibilityThreshold: c.Threshold, // Any difference in number of nodes is a different species for these tests
				Distancer:              c.Distancer,
				lastSID:                100,
			}

			// Speciate the population
			err := s.Speciate(&c.Actual)

			// Test for errors
			if !t.Run("error", mock.Error(c.HasError, err)) || c.HasError {
				return
			}

			// Test the species
			t.Run("species", testSpeciatorSpecies(c.Actual.Species, c.Expected.Species))

			// Test the genomes
			t.Run("genomes", testSpeciatorGenomes(c.Actual.Genomes, c.Expected.Genomes))

		})
	}
}

func testSpeciatorSpecies(actual, expected []evo.Species) func(*testing.T) {
	return func(t *testing.T) {
		if len(expected) != len(actual) {
			t.Log("incorrect species:", actual, expected)
			t.Errorf("incorrect number of species: expected %d, actual %d", len(expected), len(actual))
		} else {
			for _, e := range expected {
				found := false
				for _, a := range actual {
					if e.ID == a.ID {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("species %d not found in actual", e.ID)
				}
			}
		}
	}
}

func testSpeciatorGenomes(actual, expected []evo.Genome) func(*testing.T) {
	return func(t *testing.T) {
		if len(expected) != len(actual) {
			t.Errorf("incorrect number of genome: expected %d, actual %d", len(expected), len(actual))
		} else {
			for _, e := range expected {
				found := false
				for _, a := range actual {
					if e.ID == a.ID {
						if e.SpeciesID != a.SpeciesID {
							t.Errorf("incorrect species assignment for genome %d: expected %d, actual %d", e.ID, e.SpeciesID, a.SpeciesID)
						}
						found = true
						break
					}
				}
				if !found {
					t.Errorf("genome %d not found in actual", e.ID)
				}
			}
		}
	}
}

// TODO: add to test a case where there are no new additions to species. this happens when advancing only contains continuing
func TestSpeciatorModify(t *testing.T) {

	// test cases
	var cases = []struct {
		Desc     string
		Modifier float64
		Original float64
		Expected float64
		Target   int
		evo.Population
	}{
		{
			Desc:     "at target",
			Target:   2,
			Modifier: 0.5,
			Original: 1.0,
			Expected: 1.0,
			Population: evo.Population{
				Genomes: []evo.Genome{
					{ID: 1, SpeciesID: 0, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}},
					{ID: 2, SpeciesID: 0, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}, {}}}},
				},
				Species: []evo.Species{
					{ID: 10, Example: evo.Genome{Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}}},
					{ID: 20, Example: evo.Genome{Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}, {}}}}},
				},
			},
		},
		{
			Desc:     "below target",
			Target:   3,
			Modifier: 0.5,
			Original: 1.5,
			Expected: 1.0, // Should decrease
			Population: evo.Population{
				Genomes: []evo.Genome{
					{ID: 1, SpeciesID: 0, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}},
					{ID: 2, SpeciesID: 0, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}, {}}}},
				},
				Species: []evo.Species{
					{ID: 10, Example: evo.Genome{Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}}},
					{ID: 20, Example: evo.Genome{Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}, {}}}}},
				},
			},
		},
		{
			Desc:     "above target",
			Target:   1,
			Modifier: 0.5,
			Original: 1.0,
			Expected: 1.5, // Should increase
			Population: evo.Population{
				Genomes: []evo.Genome{
					{ID: 1, SpeciesID: 0, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}},
					{ID: 2, SpeciesID: 0, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}, {}}}},
				},
				Species: []evo.Species{
					{ID: 10, Example: evo.Genome{Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}}},
					{ID: 20, Example: evo.Genome{Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}, {}}}}},
				},
			},
		},
		{
			Desc:     "below target, near modifier",
			Target:   3,
			Modifier: 0.5,
			Original: 0.8,
			Expected: 0.5, // Should not go below the modifier value as a minimum
			Population: evo.Population{
				Genomes: []evo.Genome{
					{ID: 1, SpeciesID: 0, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}},
					{ID: 2, SpeciesID: 0, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}, {}}}},
				},
				Species: []evo.Species{
					{ID: 10, Example: evo.Genome{Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}}},
					{ID: 20, Example: evo.Genome{Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}, {}}}}},
				},
			},
		},
	}

	// Iterate cases
	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {

			// Create the speciator
			s := &Speciator{
				Distancer:              &MockDistancer{},
				CompatibilityThreshold: c.Original,
				CompatibilityModifier:  c.Modifier,
				TargetSpecies:          c.Target,
			}

			// Speciate
			_ = s.Speciate(&c.Population)

			// Check threshold
			if c.Expected != s.CompatibilityThreshold {
				t.Errorf("incorrect compatibility threshold: expected %f, actual %f", c.Expected, s.CompatibilityThreshold)
			}
		})
	}
}

type MockDistancer struct {
	HasError bool
}

func (d MockDistancer) Distance(a, b evo.Genome) (float64, error) {
	if d.HasError {
		return 0.0, errors.New("mock distancer error")
	}
	return math.Abs(float64(a.Complexity() - b.Complexity())), nil
}
