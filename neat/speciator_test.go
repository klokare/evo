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
				Genomes: []evo.Genome{{ID: 1}, {ID: 2}},
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
					{ID: 1, Species: 10, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}}, // 2 Nodes
					{ID: 2, Species: 10, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}}, // 2 Nodes
				},
			},
			Expected: evo.Population{
				Genomes: []evo.Genome{
					{ID: 1, Species: 10, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}}, // 2 Nodes
					{ID: 2, Species: 10, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}}, // 2 Nodes
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
					{ID: 1, Species: 10, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}},    // 2 Nodes
					{ID: 2, Species: 0, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}, {}}}}, // 3 Nodes
				},
			},
			Expected: evo.Population{
				Genomes: []evo.Genome{
					{ID: 1, Species: 10, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}},      // 2 Nodes
					{ID: 2, Species: 101, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}, {}}}}, // 2 Nodes
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
					{ID: 1, Species: 10, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}}, // 2 Nodes
					{ID: 2, Species: 10, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}}, // 2 Nodes
				},
			},
			Expected: evo.Population{
				Genomes: []evo.Genome{
					{ID: 1, Species: 10, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}}, // 2 Nodes
					{ID: 2, Species: 10, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}}, // 2 Nodes
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

			// Test the genomes
			t.Run("genomes", testSpeciatorGenomes(c.Actual.Genomes, c.Expected.Genomes))

		})
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
						if e.Species != a.Species {
							t.Errorf("incorrect species assignment for genome %d: expected %d, actual %d", e.ID, e.Species, a.Species)
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
					{ID: 1, Species: 0, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}},
					{ID: 2, Species: 0, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}, {}}}},
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
					{ID: 1, Species: 0, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}},
					{ID: 2, Species: 0, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}, {}}}},
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
					{ID: 1, Species: 0, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}},
					{ID: 2, Species: 0, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}, {}}}},
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
					{ID: 1, Species: 0, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}},
					{ID: 2, Species: 0, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}, {}}}},
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
