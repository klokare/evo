package evo

import (
	"testing"
)

func TestByFitness(t *testing.T) {
	var cases = []struct {
		Desc     string
		A, B     Genome
		Expected int8
	}{
		{
			Desc:     "equal fitness",
			A:        Genome{Fitness: 1.0},
			B:        Genome{Fitness: 1.0},
			Expected: 0,
		},
		{
			Desc:     "lower fitness",
			A:        Genome{Fitness: 0.0},
			B:        Genome{Fitness: 1.0},
			Expected: -1,
		},
		{
			Desc:     "higher fitness",
			A:        Genome{Fitness: 2.0},
			B:        Genome{Fitness: 1.0},
			Expected: 1,
		},
	}

	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {
			x := ByFitness.Compare(c.A, c.B)
			if x != c.Expected {
				t.Errorf("incorrect comparison result: expected %d, actual %d", c.Expected, x)
			}
		})
	}
}

func TestByNovelty(t *testing.T) {
	var cases = []struct {
		Desc     string
		A, B     Genome
		Expected int8
	}{
		{
			Desc:     "equal novelty",
			A:        Genome{Novelty: 1.0},
			B:        Genome{Novelty: 1.0},
			Expected: 0,
		},
		{
			Desc:     "lower novelty",
			A:        Genome{Novelty: 0.0},
			B:        Genome{Novelty: 1.0},
			Expected: -1,
		},
		{
			Desc:     "higher novelty",
			A:        Genome{Novelty: 2.0},
			B:        Genome{Novelty: 1.0},
			Expected: 1,
		},
	}

	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {
			x := ByNovelty.Compare(c.A, c.B)
			if x != c.Expected {
				t.Errorf("incorrect comparison result: expected %d, actual %d", c.Expected, x)
			}
		})
	}
}

func TestByAge(t *testing.T) {
	var cases = []struct {
		Desc     string
		A, B     Genome
		Expected int8
	}{
		{
			Desc:     "equal age", // unlikely unless genome is a duplicate
			A:        Genome{ID: 2},
			B:        Genome{ID: 2},
			Expected: 0,
		},
		{
			Desc:     "lower age",
			A:        Genome{ID: 1},
			B:        Genome{ID: 2},
			Expected: 1,
		},
		{
			Desc:     "higher age",
			A:        Genome{ID: 3},
			B:        Genome{ID: 2},
			Expected: -1,
		},
	}

	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {
			x := ByAge.Compare(c.A, c.B)
			if x != c.Expected {
				t.Errorf("incorrect comparison result: expected %d, actual %d", c.Expected, x)
			}
		})
	}
}

func TestByComplexity(t *testing.T) {
	var cases = []struct {
		Desc     string
		A, B     Genome
		Expected int8
	}{
		{
			Desc:     "equal complexity",
			A:        Genome{Encoded: Substrate{Nodes: []Node{{}, {}}}},
			B:        Genome{Encoded: Substrate{Nodes: []Node{{}, {}}}},
			Expected: 0,
		},
		{
			Desc:     "lower complexity",
			A:        Genome{Encoded: Substrate{Nodes: []Node{{}, {}}}},
			B:        Genome{Encoded: Substrate{Nodes: []Node{{}, {}, {}}}},
			Expected: 1,
		},
		{
			Desc:     "higher complexity",
			A:        Genome{Encoded: Substrate{Nodes: []Node{{}, {}, {}}}},
			B:        Genome{Encoded: Substrate{Nodes: []Node{{}, {}}}},
			Expected: -1,
		},
	}

	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {
			x := ByComplexity.Compare(c.A, c.B)
			if x != c.Expected {
				t.Errorf("incorrect comparison result: expected %d, actual %d", c.Expected, x)
			}
		})
	}
}

func TestBySolved(t *testing.T) {
	var cases = []struct {
		Desc     string
		A, B     Genome
		Expected int8
	}{
		{
			Desc:     "both unsolved",
			A:        Genome{Solved: false},
			B:        Genome{Solved: false},
			Expected: 0,
		},
		{
			Desc:     "both solved",
			A:        Genome{Solved: true},
			B:        Genome{Solved: true},
			Expected: 0,
		},
		{
			Desc:     "a solved",
			A:        Genome{Solved: true},
			B:        Genome{Solved: false},
			Expected: 1,
		},
		{
			Desc:     "b solved",
			A:        Genome{Solved: false},
			B:        Genome{Solved: true},
			Expected: -1,
		},
	}

	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {
			x := BySolved.Compare(c.A, c.B)
			if x != c.Expected {
				t.Errorf("incorrect comparison result: expected %d, actual %d", c.Expected, x)
			}
		})
	}
}

func TestSortBy(t *testing.T) {
	var cases = []struct {
		Desc     string
		Genomes  []Genome
		Expected []int64
	}{
		{
			Desc: "equal novelty, equal fitness",
			Genomes: []Genome{
				Genome{ID: 2, Fitness: 1.0, Novelty: 1.0},
				Genome{ID: 1, Fitness: 1.0, Novelty: 1.0},
			},
			Expected: []int64{2, 1},
		},
		{
			Desc: "lower novelty, equal fitness",
			Genomes: []Genome{
				Genome{ID: 2, Fitness: 1.0, Novelty: 0.0},
				Genome{ID: 1, Fitness: 1.0, Novelty: 1.0},
			},
			Expected: []int64{2, 1},
		},
		{
			Desc: "lower novelty, lower fitness",
			Genomes: []Genome{
				Genome{ID: 2, Fitness: 0.0, Novelty: 0.0},
				Genome{ID: 1, Fitness: 1.0, Novelty: 1.0},
			},
			Expected: []int64{2, 1},
		},
		{
			Desc: "lower novelty, higher fitness",
			Genomes: []Genome{
				Genome{ID: 2, Fitness: 2.0, Novelty: 0.0},
				Genome{ID: 1, Fitness: 1.0, Novelty: 1.0},
			},
			Expected: []int64{2, 1},
		},
		{
			Desc: "higher novelty, equal fitness",
			Genomes: []Genome{
				Genome{ID: 2, Fitness: 1.0, Novelty: 3.0},
				Genome{ID: 1, Fitness: 1.0, Novelty: 1.0},
			},
			Expected: []int64{1, 2},
		},
		{
			Desc: "higher novelty, lower fitness",
			Genomes: []Genome{
				Genome{ID: 2, Fitness: 0.0, Novelty: 3.0},
				Genome{ID: 1, Fitness: 1.0, Novelty: 1.0},
			},
			Expected: []int64{1, 2},
		},
		{
			Desc: "higher novelty, higher fitness",
			Genomes: []Genome{
				Genome{ID: 2, Fitness: 2.0, Novelty: 3.0},
				Genome{ID: 1, Fitness: 1.0, Novelty: 1.0},
			},
			Expected: []int64{1, 2},
		},
		{
			Desc: "equal novelty, lower fitness",
			Genomes: []Genome{
				Genome{ID: 2, Fitness: 0.0, Novelty: 1.0},
				Genome{ID: 1, Fitness: 1.0, Novelty: 1.0},
			},
			Expected: []int64{2, 1},
		},
		{
			Desc: "equal novelty, higher fitness",
			Genomes: []Genome{
				Genome{ID: 2, Fitness: 2.0, Novelty: 1.0},
				Genome{ID: 1, Fitness: 1.0, Novelty: 1.0},
			},
			Expected: []int64{1, 2},
		},
	}

	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {
			SortBy(c.Genomes, ByNovelty, ByFitness)
			for i, gid := range c.Expected {
				if gid != c.Genomes[i].ID {
					t.Errorf("incorrect genome in position %d: expected %d, actual %d", i, gid, c.Genomes[i].ID)
				}
			}
		})
	}
}
