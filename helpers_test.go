package evo

import (
	"errors"
	"testing"
)

func TestMutatorsMutate(t *testing.T) {
	var cases = []struct {
		Desc     string
		HasError bool
		Expected []bool
		Mutators []Mutator
	}{
		{
			Desc:     "no errors, no structural",
			HasError: false,
			Expected: []bool{true, true},
			Mutators: []Mutator{
				&mockNonstructural{},
				&mockNonstructural{},
			},
		},
		{
			Desc:     "has errors",
			HasError: true,
			Expected: []bool{true, true},
			Mutators: []Mutator{
				&mockNonstructural{HasError: true},
				&mockNonstructural{},
			},
		},
		{
			Desc:     "has structural, second mutator not called",
			HasError: false,
			Expected: []bool{true, false},
			Mutators: []Mutator{
				&mockStructural{},
				&mockNonstructural{},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {
			// Build a new collection
			mut := Mutators(c.Mutators)

			// Mutate a genome
			err := mut.Mutate(&Genome{})
			if c.HasError {
				if err == nil {
					t.Error("expected error not found")
				}
				return
			}
			if err != nil {
				t.Errorf("error not expected: %v", err)
			}

			// Check that right mutators were called
			for i := 0; i < len(c.Mutators); i++ {
				act := c.Mutators[i].(callable).Called()
				if act != c.Expected[i] {
					t.Errorf("incorrect called value for mutator %d: expected %t, actual %t", i, c.Expected, act)
				}
			}
		})
	}
}

type mockStructural struct {
	called   bool
	HasError bool
}

func (m *mockStructural) Mutate(g *Genome) error {
	m.called = true
	if m.HasError {
		return errors.New("mock structural error")
	}
	g.Encoded.Nodes = append(g.Encoded.Nodes, Node{})
	return nil
}

func (m *mockStructural) Called() bool { return m.called }

type mockNonstructural struct {
	called   bool
	HasError bool
}

func (m *mockNonstructural) Mutate(g *Genome) error {
	m.called = true
	if m.HasError {
		return errors.New("mock nonstructural error")
	}
	return nil
}

func (m *mockNonstructural) Called() bool { return m.called }

type callable interface{ Called() bool }
