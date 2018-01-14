package evo

import (
	"reflect"
	"testing"
)

func TestWithByFitness(t *testing.T) {

	// Use the option
	e := new(Experiment)
	option := WithByFitness()
	err := option(e)

	// There should be no error
	if err != nil {
		t.Errorf("there should be no error. actual %s", err)
	}

	// The comparer should be set correctly
	if _, ok := e.Comparer.(ByFitness); !ok {
		t.Errorf("comparer not set properly: expected: ByFitness, actual %v", reflect.TypeOf(e.Comparer))
	}
}

func TestWithByNovelty(t *testing.T) {

	// Use the option
	e := new(Experiment)
	option := WithByNovelty()
	err := option(e)

	// There should be no error
	if err != nil {
		t.Errorf("there should be no error. actual %s", err)
	}

	// The comparer should be set correctly
	if _, ok := e.Comparer.(ByNovelty); !ok {
		t.Errorf("comparer not set properly: expected: ByNovelty, actual %v", reflect.TypeOf(e.Comparer))
	}
}

func TestByFitnessCompare(t *testing.T) {

	var cases = []struct {
		Desc     string
		A, B     Genome
		Expected int8
	}{
		{
			Desc:     "two empty genomes",
			A:        Genome{},
			B:        Genome{},
			Expected: 0,
		},
		{
			Desc: "first genome is empty",
			A:    Genome{},
			B: Genome{
				ID: 2, Fitness: 2.2, Novelty: 1.2, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			Expected: -1,
		},
		{
			Desc: "second genome is empty",
			A: Genome{
				ID: 1, Fitness: 2.1, Novelty: 1.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			B:        Genome{},
			Expected: 1,
		},
		{
			Desc: "same genome",
			A: Genome{
				ID: 1, Fitness: 2.1, Novelty: 1.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			B: Genome{
				ID: 1, Fitness: 2.1, Novelty: 1.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			Expected: 0,
		},
		{
			Desc: "genome a has solution",
			A: Genome{
				ID: 1, Fitness: 2.1, Novelty: 1.1, Solved: true,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			B: Genome{
				ID: 1, Fitness: 2.1, Novelty: 1.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			Expected: 1,
		},
		{
			Desc: "genome b has solution",
			A: Genome{
				ID: 1, Fitness: 2.1, Novelty: 1.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			B: Genome{
				ID: 1, Fitness: 2.1, Novelty: 1.1, Solved: true,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			Expected: -1,
		},
		{
			Desc: "genome a has higher fitness",
			A: Genome{
				ID: 1, Fitness: 2.1, Novelty: 1.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			B: Genome{
				ID: 1, Fitness: 1.1, Novelty: 1.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			Expected: 1,
		},
		{
			Desc: "genome b has higher fitness",
			A: Genome{
				ID: 1, Fitness: 1.1, Novelty: 1.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			B: Genome{
				ID: 1, Fitness: 2.1, Novelty: 1.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			Expected: -1,
		},
		{
			Desc: "genome a has higher novelty",
			A: Genome{
				ID: 1, Fitness: 2.1, Novelty: 2.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			B: Genome{
				ID: 1, Fitness: 2.1, Novelty: 1.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			Expected: 0,
		},
		{
			Desc: "genome b has higher novelty",
			A: Genome{
				ID: 1, Fitness: 2.1, Novelty: 1.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			B: Genome{
				ID: 1, Fitness: 2.1, Novelty: 2.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			Expected: 0,
		},
		{
			Desc: "genome a has higher encoded complexity",
			A: Genome{
				ID: 1, Fitness: 2.1, Novelty: 1.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}, {Neuron: Output}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			B: Genome{
				ID: 1, Fitness: 2.1, Novelty: 1.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			Expected: -1,
		},
		{
			Desc: "genome b has higher encoded complexity",
			A: Genome{
				ID: 1, Fitness: 2.1, Novelty: 1.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			B: Genome{
				ID: 1, Fitness: 2.1, Novelty: 1.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}, {Neuron: Output}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			Expected: 1,
		},
		{
			Desc: "genome a has higher decoded complexity",
			A: Genome{
				ID: 1, Fitness: 2.1, Novelty: 1.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}, {Neuron: Output}}},
			},
			B: Genome{
				ID: 1, Fitness: 2.1, Novelty: 1.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			Expected: -1,
		},
		{
			Desc: "genome b has higher decoded complexity",
			A: Genome{
				ID: 1, Fitness: 2.1, Novelty: 1.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			B: Genome{
				ID: 1, Fitness: 2.1, Novelty: 1.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}, {Neuron: Output}}},
			},
			Expected: 1,
		},
		{
			Desc: "genome a has lower ID",
			A: Genome{
				ID: 1, Fitness: 2.1, Novelty: 2.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			B: Genome{
				ID: 2, Fitness: 2.1, Novelty: 1.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			Expected: 1,
		},
		{
			Desc: "genome b has lower ID",
			A: Genome{
				ID: 2, Fitness: 2.1, Novelty: 1.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			B: Genome{
				ID: 1, Fitness: 2.1, Novelty: 2.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			Expected: -1,
		},
	}

	cmp := ByFitness{}

	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {
			actual := cmp.Compare(c.A, c.B)
			if c.Expected != actual {
				t.Errorf("incorrect comparison value: expected %d, actual %d", c.Expected, actual)
			}
		})
	}
}

func TestByNoveltyCompare(t *testing.T) {

	var cases = []struct {
		Desc     string
		A, B     Genome
		Expected int8
	}{
		{
			Desc:     "two empty genomes",
			A:        Genome{},
			B:        Genome{},
			Expected: 0,
		},
		{
			Desc: "first genome is empty",
			A:    Genome{},
			B: Genome{
				ID: 2, Fitness: 2.2, Novelty: 1.2, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			Expected: -1,
		},
		{
			Desc: "second genome is empty",
			A: Genome{
				ID: 1, Fitness: 2.1, Novelty: 1.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			B:        Genome{},
			Expected: 1,
		},
		{
			Desc: "same genome",
			A: Genome{
				ID: 1, Fitness: 2.1, Novelty: 1.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			B: Genome{
				ID: 1, Fitness: 2.1, Novelty: 1.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			Expected: 0,
		},
		{
			Desc: "genome a has solution",
			A: Genome{
				ID: 1, Fitness: 2.1, Novelty: 1.1, Solved: true,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			B: Genome{
				ID: 1, Fitness: 2.1, Novelty: 1.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			Expected: 1,
		},
		{
			Desc: "genome b has solution",
			A: Genome{
				ID: 1, Fitness: 2.1, Novelty: 1.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			B: Genome{
				ID: 1, Fitness: 2.1, Novelty: 1.1, Solved: true,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			Expected: -1,
		},
		{
			Desc: "genome a has higher fitness",
			A: Genome{
				ID: 1, Fitness: 2.1, Novelty: 1.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			B: Genome{
				ID: 1, Fitness: 1.1, Novelty: 1.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			Expected: 0,
		},
		{
			Desc: "genome b has higher fitness",
			A: Genome{
				ID: 1, Fitness: 1.1, Novelty: 1.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			B: Genome{
				ID: 1, Fitness: 2.1, Novelty: 1.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			Expected: 0,
		},
		{
			Desc: "genome a has higher novelty",
			A: Genome{
				ID: 1, Fitness: 2.1, Novelty: 2.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			B: Genome{
				ID: 1, Fitness: 2.1, Novelty: 1.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			Expected: 1,
		},
		{
			Desc: "genome b has higher novelty",
			A: Genome{
				ID: 1, Fitness: 2.1, Novelty: 1.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			B: Genome{
				ID: 1, Fitness: 2.1, Novelty: 2.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			Expected: -1,
		},
		{
			Desc: "genome a has higher encoded complexity",
			A: Genome{
				ID: 1, Fitness: 2.1, Novelty: 1.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}, {Neuron: Output}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			B: Genome{
				ID: 1, Fitness: 2.1, Novelty: 1.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			Expected: 1,
		},
		{
			Desc: "genome b has higher encoded complexity",
			A: Genome{
				ID: 1, Fitness: 2.1, Novelty: 1.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			B: Genome{
				ID: 1, Fitness: 2.1, Novelty: 1.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}, {Neuron: Output}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			Expected: -1,
		},
		{
			Desc: "genome a has higher decoded complexity",
			A: Genome{
				ID: 1, Fitness: 2.1, Novelty: 1.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}, {Neuron: Output}}},
			},
			B: Genome{
				ID: 1, Fitness: 2.1, Novelty: 1.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			Expected: 1,
		},
		{
			Desc: "genome b has higher decoded complexity",
			A: Genome{
				ID: 1, Fitness: 2.1, Novelty: 1.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			B: Genome{
				ID: 1, Fitness: 2.1, Novelty: 1.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}, {Neuron: Output}}},
			},
			Expected: -1,
		},
		{
			Desc: "genome a has lower ID",
			A: Genome{
				ID: 1, Fitness: 2.1, Novelty: 2.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			B: Genome{
				ID: 2, Fitness: 2.1, Novelty: 2.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			Expected: 1,
		},
		{
			Desc: "genome b has lower ID",
			A: Genome{
				ID: 2, Fitness: 2.1, Novelty: 2.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			B: Genome{
				ID: 1, Fitness: 2.1, Novelty: 2.1, Solved: false,
				Encoded: Substrate{Nodes: []Node{{Neuron: Input}, {Neuron: Hidden}}},
				Decoded: Substrate{Nodes: []Node{{Neuron: Output}}},
			},
			Expected: -1,
		},
	}

	cmp := ByNovelty{}

	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {
			actual := cmp.Compare(c.A, c.B)
			if c.Expected != actual {
				t.Errorf("incorrect comparison value: expected %d, actual %d", c.Expected, actual)
			}
		})
	}
}
