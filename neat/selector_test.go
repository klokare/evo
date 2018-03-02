package neat

import (
	"testing"

	"github.com/klokare/evo"
	"github.com/klokare/evo/internal/mock"
	"github.com/klokare/evo/internal/test"
)

func TestSelectorSelect(t *testing.T) {

	var cases = []struct {
		Desc       string
		Population evo.Population
		Continuing []evo.Genome
		Offspring  map[int64]int
		HasError   bool
	}{
		{
			Desc: "all species still active",
			Population: evo.Population{
				Species: []evo.Species{
					{ID: 10},
					{ID: 20},
					{ID: 30},
				},
				Genomes: []evo.Genome{
					{ID: 6, SpeciesID: 10, Fitness: 8.0},
					{ID: 9, SpeciesID: 10, Fitness: 4.0},
					{ID: 2, SpeciesID: 10, Fitness: 5.0, Encoded: evo.Substrate{Nodes: []evo.Node{{}}}},
					{ID: 8, SpeciesID: 20, Fitness: 5.0, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}},
					{ID: 1, SpeciesID: 20, Fitness: 4.0},
					{ID: 3, SpeciesID: 30, Fitness: 3.5},
					{ID: 4, SpeciesID: 30, Fitness: 3.5},
					{ID: 5, SpeciesID: 30, Fitness: 3.5},
				},
			},
			Continuing: []evo.Genome{
				{ID: 6, SpeciesID: 10, Fitness: 8.0},
				{ID: 8, SpeciesID: 20, Fitness: 5.0, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}},
				{ID: 3, SpeciesID: 30, Fitness: 3.5},
			},
			Offspring: map[int64]int{
				10: 2,
				20: 1,
				30: 1,
			},
		},
		{
			Desc: "species 10 decayed",
			Population: evo.Population{
				Species: []evo.Species{
					{ID: 10, Decay: 0.5},
					{ID: 20},
					{ID: 30},
				},
				Genomes: []evo.Genome{
					{ID: 6, SpeciesID: 10, Fitness: 8.0},
					{ID: 9, SpeciesID: 10, Fitness: 4.0},
					{ID: 2, SpeciesID: 10, Fitness: 5.0, Encoded: evo.Substrate{Nodes: []evo.Node{{}}}},
					{ID: 8, SpeciesID: 20, Fitness: 5.0, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}},
					{ID: 1, SpeciesID: 20, Fitness: 4.0},
					{ID: 3, SpeciesID: 30, Fitness: 3.5},
					{ID: 4, SpeciesID: 30, Fitness: 3.5},
					{ID: 5, SpeciesID: 30, Fitness: 3.5},
				},
			},
			Continuing: []evo.Genome{
				{ID: 6, SpeciesID: 10, Fitness: 8.0},                                                    // absolute leader
				{ID: 8, SpeciesID: 20, Fitness: 5.0, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}}, // Same as ID 8 but more complex
				{ID: 3, SpeciesID: 30, Fitness: 3.5},                                                    // same as ID 5 but younger
			},
			Offspring: map[int64]int{
				10: 1, // Fewer because of decay
				20: 2,
				30: 1,
			},
		},
		{
			Desc: "species 10 stagnant",
			Population: evo.Population{
				Species: []evo.Species{
					{ID: 10, Decay: 1.0},
					{ID: 20},
					{ID: 30},
				},
				Genomes: []evo.Genome{
					{ID: 6, SpeciesID: 10, Fitness: 8.0},
					{ID: 9, SpeciesID: 10, Fitness: 4.0},
					{ID: 2, SpeciesID: 10, Fitness: 5.0, Encoded: evo.Substrate{Nodes: []evo.Node{{}}}},
					{ID: 8, SpeciesID: 20, Fitness: 5.0, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}},
					{ID: 1, SpeciesID: 20, Fitness: 4.0},
					{ID: 3, SpeciesID: 30, Fitness: 3.5},
					{ID: 4, SpeciesID: 30, Fitness: 3.5},
					{ID: 5, SpeciesID: 30, Fitness: 3.5},
				},
			},
			Continuing: []evo.Genome{
				{ID: 6, SpeciesID: 10, Fitness: 8.0}, // stays in because overall best
				{ID: 8, SpeciesID: 20, Fitness: 5.0, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}},
				{ID: 3, SpeciesID: 30, Fitness: 3.5},
			},
			Offspring: map[int64]int{
				20: 3,
				30: 1,
			},
		},
		{
			Desc: "species 10 stagnant, best shared with other species",
			Population: evo.Population{
				Species: []evo.Species{
					{ID: 10, Decay: 1.0},
					{ID: 20},
					{ID: 30},
				},
				Genomes: []evo.Genome{
					{ID: 6, SpeciesID: 10, Fitness: 5.0},
					{ID: 9, SpeciesID: 10, Fitness: 4.0},
					{ID: 2, SpeciesID: 10, Fitness: 5.0, Encoded: evo.Substrate{Nodes: []evo.Node{{}}}},
					{ID: 8, SpeciesID: 20, Fitness: 5.0, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}},
					{ID: 1, SpeciesID: 20, Fitness: 4.0},
					{ID: 3, SpeciesID: 30, Fitness: 3.5},
					{ID: 4, SpeciesID: 30, Fitness: 3.5},
					{ID: 5, SpeciesID: 30, Fitness: 3.5},
				},
			},
			Continuing: []evo.Genome{ // None from species 10 because stagnant and another species has same fitness
				{ID: 8, SpeciesID: 20, Fitness: 5.0, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}}, // Same as ID 8 but more complex
				{ID: 3, SpeciesID: 30, Fitness: 3.5},                                                    // same as ID 5 but younger
			},
			Offspring: map[int64]int{
				20: 4,
				30: 1,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {

			// Create the selector
			s := &Selector{
				Compare:        evo.ByFitness,
				PopulationSize: len(c.Population.Genomes),
			}

			// Select
			acs, aps, err := s.Select(c.Population)

			// Check error
			if !t.Run("error", test.Error(c.HasError, err)) || c.HasError {
				return
			}

			// Correct number of continuing + parents
			if len(c.Population.Genomes) != len(acs)+len(aps) {
				t.Errorf("incorrect number of continuing and parents: expected %d, actual %d", len(c.Population.Genomes), len(acs)+len(aps))
			}
			// Correct continuing
			ecs := c.Continuing
			if len(ecs) != len(acs) {
				t.Errorf("incorrect number of continuing: expected %d, actual %d", len(ecs), len(acs))
			} else {
				for _, eg := range ecs {
					found := false
					for _, ag := range acs {
						if eg.ID == ag.ID {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("genome %d not found in continuing", eg.ID)
					}
				}
			}

			// Correct (minimum) count of offspring
			aoff := make(map[int64]int, len(aps))
			for _, ap := range aps {
				aoff[ap[0].SpeciesID]++
			}
			eoff := c.Offspring
			if len(eoff) != len(aoff) {
				t.Errorf("incorrect number of species in parent groups: expected %d, actual %d", len(eoff), len(aoff))
			} else {
				for esid, ecnt := range eoff {
					if acnt, ok := aoff[esid]; ok {
						if acnt < ecnt {
							t.Errorf("incorrect number of offspring for species %d: expected %d, actual %d", esid, ecnt, acnt)
						}
					} else {
						t.Errorf("species %d not in parenting group", esid)
					}
				}
			}
		})
	}
}
func TestSelectorSortRank(t *testing.T) {

	var cases = []struct {
		Desc     string
		Original []evo.Genome
		Expected []evo.Genome
		Ranks    map[int64]float64
	}{
		{
			Desc: "out of order with some times",
			Original: []evo.Genome{
				{ID: 1, Fitness: 4.0},                                                    // middle of pack
				{ID: 3, Fitness: 3.5},                                                    // same as ID 5 but younger
				{ID: 5, Fitness: 3.5},                                                    // same as ID 3 but older
				{ID: 2, Fitness: 5.0, Encoded: evo.Substrate{Nodes: []evo.Node{{}}}},     // Same as ID 8 but less complex
				{ID: 8, Fitness: 5.0, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}}, // Same as ID 8 but more complex
				{ID: 6, Fitness: 8.0},                                                    // absolute leader
			},
			Expected: []evo.Genome{
				{ID: 6, Fitness: 8.0},                                                    // absolute leader
				{ID: 2, Fitness: 5.0, Encoded: evo.Substrate{Nodes: []evo.Node{{}}}},     // Same as ID 8 but less complex
				{ID: 8, Fitness: 5.0, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}}, // Same as ID 8 but more complex
				{ID: 1, Fitness: 4.0},                                                    // middle of pack
				{ID: 3, Fitness: 3.5},                                                    // same as ID 5 but younger
				{ID: 5, Fitness: 3.5},                                                    // same as ID 3 but older
			},
			Ranks: map[int64]float64{
				6: 6.0,
				2: 5.0,
				8: 5.0,
				1: 3.0,
				3: 2.0,
				5: 2.0,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {

			// Sort and rank
			ars := sortRank(evo.ByFitness, c.Original)
			actual := c.Original

			// There should be the same order
			if len(c.Expected) != len(actual) {
				t.Errorf("incorrect number of genomes: expected %d, actual %d", len(c.Expected), len(actual))
			} else {
				for i, expected := range c.Expected {
					if expected.ID != actual[i].ID {
						t.Errorf("incorrect genome at %d: expected %d, actual %d", i, expected.ID, actual[i].ID)
					}
				}
			}

			// There should be the same ranks
			ers := c.Ranks
			if len(ers) != len(ars) {
				t.Errorf("incorrect number of rankings: expected %d, actual %d", len(ers), len(ars))
			} else {
				for gid, er := range ers {
					ar, ok := ars[gid]
					if !ok {
						t.Errorf("ranking not found for genome %d", gid)
					} else if er != ar {
						t.Errorf("incorrect ranking for genome %d: expected %f, actual %f", gid, er, ar)
					}
				}
			}
		})
	}
}

func TestSelectorMutateOnly(t *testing.T) {

	// Create the initial population
	pop := evo.Population{
		Species: []evo.Species{
			{ID: 10},
			{ID: 20},
			{ID: 30},
		},
		Genomes: []evo.Genome{
			{ID: 6, SpeciesID: 10, Fitness: 8.0},
			{ID: 9, SpeciesID: 10, Fitness: 4.0},
			{ID: 2, SpeciesID: 10, Fitness: 5.0, Encoded: evo.Substrate{Nodes: []evo.Node{{}}}},
			{ID: 8, SpeciesID: 20, Fitness: 5.0, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}},
			{ID: 1, SpeciesID: 20, Fitness: 4.0},
			{ID: 3, SpeciesID: 30, Fitness: 3.5},
			{ID: 4, SpeciesID: 30, Fitness: 3.5},
			{ID: 5, SpeciesID: 30, Fitness: 3.5},
		},
	}

	// Create the selector
	// Create the selector
	s := &Selector{
		Compare:               evo.ByFitness,
		PopulationSize:        len(pop.Genomes),
		MutateOnlyProbability: 1.0,
	}

	// Select
	_, aps, err := s.Select(pop)

	// Check error
	if !t.Run("error", test.Error(false, err)) {
		return
	}

	// There should only be single parents
	for i, ps := range aps {
		if len(ps) != 1 {
			t.Errorf("incorrect number of parents in group %d: expected 1, actual %d", i, len(ps))
		}
	}
}

func TestSelectorInterspecies(t *testing.T) {

	// Create the initial population
	pop := evo.Population{
		Species: []evo.Species{
			{ID: 10},
			{ID: 20},
			{ID: 30},
		},
		Genomes: []evo.Genome{
			{ID: 6, SpeciesID: 10, Fitness: 8.0},
			{ID: 9, SpeciesID: 10, Fitness: 4.0},
			{ID: 2, SpeciesID: 10, Fitness: 5.0, Encoded: evo.Substrate{Nodes: []evo.Node{{}}}},
			{ID: 8, SpeciesID: 20, Fitness: 5.0, Encoded: evo.Substrate{Nodes: []evo.Node{{}, {}}}},
			{ID: 1, SpeciesID: 20, Fitness: 4.0},
			{ID: 3, SpeciesID: 30, Fitness: 3.5},
			{ID: 4, SpeciesID: 30, Fitness: 3.5},
			{ID: 5, SpeciesID: 30, Fitness: 3.5},
		},
	}

	// Create the selector
	// Create the selector
	s := &Selector{
		Compare:                     evo.ByFitness,
		PopulationSize:              len(pop.Genomes),
		InterspeciesMateProbability: 1.0,
	}

	// Select
	_, aps, err := s.Select(pop)

	// Check error
	if !t.Run("error", test.Error(false, err)) {
		return
	}

	// There couplees should be interspecies
	for i, ps := range aps {
		if ps[0].SpeciesID == ps[1].SpeciesID {
			t.Errorf("expected parents to be of different species in group %d", i)
		}
	}
}

func TestWithSelector(t *testing.T) {
	e := new(evo.Experiment)

	// Experiment has no Comparer
	err := WithSelector(&mock.Configurer{})(e)
	if err == nil {
		t.Errorf("error expected but not found")
	}

	// Experiment has Comparer
	e.Compare = evo.ByFitness
	err = WithSelector(&mock.Configurer{})(e)
	if err != nil {
		t.Errorf("error not expected but had: %v", err)
	}
	if _, ok := e.Selector.(*Selector); !ok {
		t.Errorf("selector incorrectly set")
	}

	// Configurer has error
	e.Compare = evo.ByFitness
	err = WithSelector(&mock.Configurer{HasError: true})(e)
	if err == nil {
		t.Errorf("error expected but not found")
	}
}

func TestSelectorAdjCountHigh(t *testing.T) {
	// the tests above never produce a count above target, this checks that scenario
	cnt := 4
	tgt := 3
	off := map[int64]int{
		1: 0,
		2: 1,
		3: 3,
	}

	adjCounts(off, cnt, tgt)
	cnt = 0
	for _, x := range off {
		cnt += x
	}
	if cnt != tgt {
		t.Errorf("incorrect count after adjustment: expected %d, actual %d", tgt, cnt)
	}
}

func TestSelectorStagnant(t *testing.T) {

	// Create a stagnant population
	p := evo.Population{
		Species: []evo.Species{
			{ID: 1, Decay: 1.0},
			{ID: 2, Decay: 1.0}, // has best
		},
		Genomes: []evo.Genome{
			{ID: 1, SpeciesID: 1, Fitness: 1.0},
			{ID: 2, SpeciesID: 1, Fitness: 2.0},
			{ID: 3, SpeciesID: 2, Fitness: 2.0},
			{ID: 4, SpeciesID: 2, Fitness: 3.0}, // best
		},
	}

	// Select
	s := &Selector{
		Compare:        evo.ByFitness,
		PopulationSize: len(p.Genomes),
	}
	cs, ps, err := s.Select(p)

	// Test for error
	t.Run("error", test.Error(false, err))

	// Test for continuing
	if len(cs) != 1 {
		t.Errorf("incorrect number of continuing: expected 1, actual: %d", len(cs))
	} else {
		if cs[0].ID != 4 {
			t.Errorf("incorrect continuing genome: expected 4, actual %d", cs[0].ID)
		}
	}

	// Test for parents
	if len(ps) != 3 {
		t.Errorf("incorrect number of parent groups: expected 3, actual %d", len(ps))
	} else {
		for i, pg := range ps {
			if len(pg) != 1 {
				t.Errorf("incorrect number of parents in group %d: expected 1, actaul %d", i, len(pg))
			} else if pg[0].ID != 4 {
				t.Errorf("incorrect genome in parent group %d: expected 5, actual %d", i, pg[i].ID)
			}
		}
	}

	// Species decay and champion should be reset
	for _, z := range p.Species {
		if z.Decay != 0.0 {
			t.Errorf("incorrect decay value for species %d: expected 0.0, actual %f", z.ID, z.Decay)
		}
		if z.Champion != 0 {
			t.Errorf("incorrect champion for species %d: expected 0.0, actual %d", z.ID, z.Champion)
		}
	}
}

func TestToggleMutateOnly(t *testing.T) {

	// Create a new selector with a non-zero (so we can check) mutate only probability
	s := &Selector{MutateOnlyProbability: 0.5}

	// Selector begins toggled "off" so toggle "on"
	s.ToggleMutateOnly(true)
	if s.mop == 0 {
		t.Errorf("incorrect mop value: expected %f, actual %f", 0.5, s.mop)
	}
	if s.MutateOnlyProbability != 1.0 {
		t.Errorf("incorrect mutate only value when toggled on: expected 1.0, actual %f", s.MutateOnlyProbability)
	}

	// Toggle back off
	s.ToggleMutateOnly(false)
	if s.MutateOnlyProbability != 0.5 {
		t.Errorf("incorrect mutate only value when toggled off: expected %f, actual %f", 0.5, s.MutateOnlyProbability)
	}
}
