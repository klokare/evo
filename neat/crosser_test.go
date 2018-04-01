package neat

import (
	"math"
	"testing"

	"github.com/klokare/evo"
	"github.com/klokare/evo/internal/mock"
)

func TestCrosserParentErrors(t *testing.T) {
	var cases = []struct {
		Desc    string
		Parents []evo.Genome
	}{
		{Desc: "zero parents", Parents: []evo.Genome{}},
		{Desc: "too many parents", Parents: make([]evo.Genome, 3)},
	}

	z := new(Crosser)
	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {
			_, err := z.Cross(c.Parents...)
			t.Run("error", mock.Error(true, err))
		})
	}
}

// Test the crossing of 1 and 2 parents. Errors and enabling disabled connections are tested above.
// Testing probability of inheritance is tested below.
func TestCrosser(t *testing.T) {

	// Though nonsensical in practice, using weight and bias values to track lineage
	var cases = []struct {
		Desc    string
		Parents []evo.Genome
		Child   evo.Genome
	}{
		{
			Desc: "single parent", // Child should be a clone
			Parents: []evo.Genome{
				{
					ID: 2, Fitness: 1.1, Novelty: 2.2, Solved: true,
					Encoded: evo.Substrate{
						Nodes: []evo.Node{
							{Position: evo.Position{Layer: 1.0, X: 0.0}},
							{Position: evo.Position{Layer: 1.0, X: 0.5}},
							{Position: evo.Position{Layer: 1.0, X: 1.0}},
						},
						Conns: []evo.Conn{
							{Source: evo.Position{Layer: 1.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
							{Source: evo.Position{Layer: 1.0, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 0.5}},
						},
					},
				},
			},
			Child: evo.Genome{
				ID: 0, Fitness: 0.0, Novelty: 0.0, Solved: false,
				Encoded: evo.Substrate{
					Nodes: []evo.Node{
						{Position: evo.Position{Layer: 1.0, X: 0.0}},
						{Position: evo.Position{Layer: 1.0, X: 0.5}},
						{Position: evo.Position{Layer: 1.0, X: 1.0}},
					},
					Conns: []evo.Conn{
						{Source: evo.Position{Layer: 1.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
						{Source: evo.Position{Layer: 1.0, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 0.5}},
					},
				},
			},
		},
		{
			Desc: "two parents, equal fitness", // Child should have all nodes and conns
			Parents: []evo.Genome{
				{
					ID: 2, Fitness: 1.1, Novelty: 2.2, Solved: true,
					Encoded: evo.Substrate{
						Nodes: []evo.Node{
							{Position: evo.Position{Layer: 0.0, X: 0.0}},
							{Position: evo.Position{Layer: 0.0, X: 0.5}}, // Not in 2nd parent
							{Position: evo.Position{Layer: 0.0, X: 1.0}},
							{Position: evo.Position{Layer: 1.0, X: 0.5}},
						},
						Conns: []evo.Conn{
							{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
							{Source: evo.Position{Layer: 0.0, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 0.5}}, // not in 2nd parent
							{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
						},
					},
				},
				{
					ID: 3, Fitness: 1.1, Novelty: 3.3, Solved: false,
					Encoded: evo.Substrate{
						Nodes: []evo.Node{
							{Position: evo.Position{Layer: 0.0, X: 0.0}},
							{Position: evo.Position{Layer: 0.0, X: 1.0}},
							{Position: evo.Position{Layer: 0.5, X: 0.5}}, // Not in 1st parent
							{Position: evo.Position{Layer: 1.0, X: 0.5}},
						},
						Conns: []evo.Conn{
							{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
							{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
							{Source: evo.Position{Layer: 0.5, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 0.5}}, // not in 1st parent
						},
					},
				},
			},
			Child: evo.Genome{
				ID: 0, Fitness: 0, Novelty: 0, Solved: false,
				Encoded: evo.Substrate{
					Nodes: []evo.Node{
						{Position: evo.Position{Layer: 0.0, X: 0.0}},
						{Position: evo.Position{Layer: 0.0, X: 0.5}}, // Not in 2nd parent
						{Position: evo.Position{Layer: 0.0, X: 1.0}},
						{Position: evo.Position{Layer: 0.5, X: 0.5}}, // Not in 1st parent
						{Position: evo.Position{Layer: 1.0, X: 0.5}},
					},
					Conns: []evo.Conn{
						{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
						{Source: evo.Position{Layer: 0.0, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 0.5}}, // not in 2nd parent
						{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
						{Source: evo.Position{Layer: 0.5, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 0.5}}, // not in 1st parent
					},
				},
			},
		},
		{
			Desc: "two parents, parent 1 more fit", // Child should have nodes and conns from more fit parent
			Parents: []evo.Genome{
				{
					ID: 2, Fitness: 2.1, Novelty: 2.2, Solved: true,
					Encoded: evo.Substrate{
						Nodes: []evo.Node{
							{Position: evo.Position{Layer: 0.0, X: 0.0}},
							{Position: evo.Position{Layer: 0.0, X: 0.5}}, // Not in 2nd parent
							{Position: evo.Position{Layer: 0.0, X: 1.0}},
							{Position: evo.Position{Layer: 1.0, X: 0.5}},
						},
						Conns: []evo.Conn{
							{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
							{Source: evo.Position{Layer: 0.0, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 0.5}}, // not in 2nd parent
							{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
						},
					},
				},
				{
					ID: 3, Fitness: 1.1, Novelty: 3.3, Solved: false,
					Encoded: evo.Substrate{
						Nodes: []evo.Node{
							{Position: evo.Position{Layer: 0.0, X: 0.0}},
							{Position: evo.Position{Layer: 0.0, X: 1.0}},
							{Position: evo.Position{Layer: 0.5, X: 0.5}}, // Not in 1st parent
							{Position: evo.Position{Layer: 1.0, X: 0.5}},
						},
						Conns: []evo.Conn{
							{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
							{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
							{Source: evo.Position{Layer: 0.5, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 0.5}}, // not in 1st parent
						},
					},
				},
			},
			Child: evo.Genome{
				ID: 0, Fitness: 0, Novelty: 0, Solved: false,
				Encoded: evo.Substrate{
					Nodes: []evo.Node{
						{Position: evo.Position{Layer: 0.0, X: 0.0}},
						{Position: evo.Position{Layer: 0.0, X: 0.5}}, // Not in 2nd parent
						{Position: evo.Position{Layer: 0.0, X: 1.0}},
						{Position: evo.Position{Layer: 1.0, X: 0.5}},
					},
					Conns: []evo.Conn{
						{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
						{Source: evo.Position{Layer: 0.0, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 0.5}}, // not in 2nd parent
						{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
					},
				},
			},
		},
		{
			Desc: "two parents, parent 2 more fit", // Child should have nodes and conns from more fit parent
			Parents: []evo.Genome{
				{
					ID: 2, Fitness: 1.1, Novelty: 2.2, Solved: true,
					Encoded: evo.Substrate{
						Nodes: []evo.Node{
							{Position: evo.Position{Layer: 0.0, X: 0.0}},
							{Position: evo.Position{Layer: 0.0, X: 0.5}}, // Not in 2nd parent
							{Position: evo.Position{Layer: 0.0, X: 1.0}},
							{Position: evo.Position{Layer: 1.0, X: 0.5}},
						},
						Conns: []evo.Conn{
							{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
							{Source: evo.Position{Layer: 0.0, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 0.5}}, // not in 2nd parent
							{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
						},
					},
				},
				{
					ID: 3, Fitness: 2.1, Novelty: 3.3, Solved: false,
					Encoded: evo.Substrate{
						Nodes: []evo.Node{
							{Position: evo.Position{Layer: 0.0, X: 0.0}},
							{Position: evo.Position{Layer: 0.0, X: 1.0}},
							{Position: evo.Position{Layer: 0.5, X: 0.5}}, // Not in 1st parent
							{Position: evo.Position{Layer: 1.0, X: 0.5}},
						},
						Conns: []evo.Conn{
							{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
							{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
							{Source: evo.Position{Layer: 0.5, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 0.5}}, // not in 1st parent
						},
					},
				},
			},
			Child: evo.Genome{
				ID: 0, Fitness: 0, Novelty: 0, Solved: false,
				Encoded: evo.Substrate{
					Nodes: []evo.Node{
						{Position: evo.Position{Layer: 0.0, X: 0.0}},
						{Position: evo.Position{Layer: 0.0, X: 1.0}},
						{Position: evo.Position{Layer: 0.5, X: 0.5}}, // Not in 1st parent
						{Position: evo.Position{Layer: 1.0, X: 0.5}},
					},
					Conns: []evo.Conn{
						{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
						{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
						{Source: evo.Position{Layer: 0.5, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 0.5}}, // not in 1st parent
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {

			// Cross the parents
			cmp := evo.ByFitness
			crs := &Crosser{EnableProbability: 0.0, Comparison: cmp}
			child, err := crs.Cross(c.Parents...)

			// Check if error was expected
			if t.Run("Error", mock.Error(false, err)) == false {
				t.Fail()
			}
			if err != nil {
				return // Error was received so the rest of the tests don't make sense
			}

			// Compare the actual child to the Expected
			t.Run("Child", testCompareChild(child, c.Child))

		})
	}
}

func testCompareChild(actual, expected evo.Genome) func(*testing.T) {
	return func(t *testing.T) {

		// Compare the structure
		t.Run("Nodes", testCompareChildNodes(actual.Encoded.Nodes, expected.Encoded.Nodes))
		t.Run("Conns", testCompareChildConns(actual.Encoded.Conns, expected.Encoded.Conns))
	}
}

func testCompareChildNodes(actual, expected []evo.Node) func(*testing.T) {
	return func(t *testing.T) {

		// There should be the same number of nodes
		if len(actual) != len(expected) {
			t.Errorf("incorrect number of nodes. expected %d, actual %d", len(expected), len(actual))
			t.FailNow()
		}

		// Nodes should match
		for _, en := range expected {
			found := false
			for _, an := range actual {
				if an.Compare(en) == 0 {
					// Properties should match
					if an.Neuron != en.Neuron {
						t.Errorf("incorrect neuron type. expected %v, actual %v", en.Neuron, an.Neuron)
					}
					if an.Activation != en.Activation {
						t.Errorf("incorrect activation type. expected %v, actual %v", en.Activation, an.Activation)
					}

					// Note match and move to next node
					found = true
					break
				}
			}
			if !found {
				t.Errorf("node not found. expected node %v", en.Position)
			}
		}
	}
}

func testCompareChildConns(actual, expected []evo.Conn) func(*testing.T) {
	return func(t *testing.T) {

		// There should be the same number of connections
		if len(actual) != len(expected) {
			t.Errorf("incorrect number of connections. expected %d, actual %d", len(expected), len(actual))
			t.FailNow()
		}

		// Nodes should match
		for _, ec := range expected {
			found := false
			for _, ac := range actual {
				if ac.Compare(ec) == 0 {
					// Properties should match
					if ac.Source != ec.Source {
						t.Errorf("incorrect source. expected %v, actual %v", ec.Source, ac.Source)
					}
					if ac.Target != ec.Target {
						t.Errorf("incorrect target. expected %v, actual %v", ec.Target, ac.Target)
					}
					if ac.Weight != ec.Weight {
						t.Errorf("incorrect weight. expected %f, actual %f", ec.Weight, ac.Weight)
					}
					if ac.Enabled != ec.Enabled {
						t.Log(ac)
						t.Errorf("incorrect enabled. expected %v, actual %v", ec.Enabled, ac.Enabled)
					}

					// Note match and move to next connection
					found = true
					break
				}
			}
			if !found {
				t.Errorf("connection not found. expected conn %v -> %v", ec.Source, ec.Target)
			}
		}
	}
}

func TestCrosserEnable(t *testing.T) {

	var tests = []struct {
		EnableProbability float64
	}{
		{EnableProbability: 0.0},
		{EnableProbability: 0.3},
		{EnableProbability: 0.8},
		{EnableProbability: 1.0},
	}

	for _, test := range tests {

		// Create a new crosser
		crs := &Crosser{EnableProbability: test.EnableProbability}

		// Create the parent
		p1 := evo.Genome{
			Encoded: evo.Substrate{
				Conns: []evo.Conn{
					{Enabled: false},
				},
			},
		}

		// Run test
		tries := 10000
		cnt := 0
		for i := 0; i < tries; i++ {
			child, _ := crs.Cross(p1)
			if child.Encoded.Conns[0].Enabled {
				cnt++
			}
		}

		avg := float64(cnt) / float64(tries)
		if math.Abs(test.EnableProbability-avg) > 0.05 {
			t.Errorf("incorrect enable rate. expected %f, actual, %f", test.EnableProbability, avg)
		}
	}
}

func TestCrosserRandomness(t *testing.T) {

	// Two equivalent parents, differing only in trait, bias, and weight values.
	p1 := evo.Genome{
		Traits: []float64{1.0},
		Encoded: evo.Substrate{
			Nodes: []evo.Node{{Position: evo.Position{Layer: 0.5, X: 0.5}, Bias: 1.0}},
			Conns: []evo.Conn{{Source: evo.Position{Layer: 0.5, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 1.0}},
		},
	}
	p2 := evo.Genome{
		Traits: []float64{2.0},
		Encoded: evo.Substrate{
			Nodes: []evo.Node{{Position: evo.Position{Layer: 0.5, X: 0.5}, Bias: 2.0}},
			Conns: []evo.Conn{{Source: evo.Position{Layer: 0.5, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 2.0}},
		},
	}

	// Create a crosser
	crs := &Crosser{Comparison: evo.ByFitness}

	// Run the tests
	t.Run("traits", func(t *testing.T) {
		var cnt1, cnt2 float64
		n := 10000
		for i := 0; i < n; i++ {
			child, _ := crs.Cross(p1, p2)
			if child.Traits[0] == 1.0 {
				cnt1 += 1.0
			} else if child.Traits[0] == 2.0 {
				cnt2 += 1.0
			}
		}
		avg1 := cnt1 / float64(n)
		avg2 := cnt2 / float64(n)
		if avg1 < 0.45 {
			t.Errorf("incorrect rate of traits from parent 1: expected 0.5, actual %f", avg1)
		}
		if avg2 < 0.45 {
			t.Errorf("incorrect rate of traits from parent 2: expected 0.5, actual %f", avg2)
		}
	})

	t.Run("nodes", func(t *testing.T) {
		var cnt1, cnt2 float64
		n := 10000
		for i := 0; i < n; i++ {
			child, _ := crs.Cross(p1, p2)
			if child.Encoded.Nodes[0].Bias == 1.0 {
				cnt1 += 1.0
			} else if child.Encoded.Nodes[0].Bias == 2.0 {
				cnt2 += 1.0
			}
		}
		avg1 := cnt1 / float64(n)
		avg2 := cnt2 / float64(n)
		if avg1 < 0.45 {
			t.Errorf("incorrect rate of nodes from parent 1: expected 0.5, actual %f", avg1)
		}
		if avg2 < 0.45 {
			t.Errorf("incorrect rate of nodes from parent 2: expected 0.5, actual %f", avg2)
		}
	})
	t.Run("conns", func(t *testing.T) {
		var cnt1, cnt2 float64
		n := 10000
		for i := 0; i < n; i++ {
			child, _ := crs.Cross(p1, p2)
			if child.Encoded.Conns[0].Weight == 1.0 {
				cnt1 += 1.0
			} else if child.Encoded.Conns[0].Weight == 2.0 {
				cnt2 += 1.0
			}
		}
		avg1 := cnt1 / float64(n)
		avg2 := cnt2 / float64(n)
		if avg1 < 0.45 { // even at 10k iterations, rates can be as low as .48. If we get above .45 then we probably have randomness
			t.Errorf("incorrect rate of conns from parent 1: expected 0.5, actual %f", avg1)
		}
		if avg2 < 0.45 {
			t.Errorf("incorrect rate of conns from parent 2: expected 0.5, actual %f", avg2)
		}
	})
}

func TestCrosserDisjointed(t *testing.T) {

	var cases = []struct {
		Desc    string
		Parents []evo.Genome
		Child   evo.Genome
	}{
		{
			Desc: "extra disjoint node p1",
			Parents: []evo.Genome{
				{
					ID: 2, Fitness: 1.1, Novelty: 2.2, Solved: true,
					Encoded: evo.Substrate{
						Nodes: []evo.Node{
							{Position: evo.Position{Layer: 0.0, X: 0.0}},
							{Position: evo.Position{Layer: 0.0, X: 0.5}},
							{Position: evo.Position{Layer: 0.0, X: 1.0}},
							{Position: evo.Position{Layer: 1.0, X: 0.5}},
						},
						Conns: []evo.Conn{
							{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
							{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
						},
					},
				},
				{
					ID: 3, Fitness: 1.1, Novelty: 3.3, Solved: false,
					Encoded: evo.Substrate{
						Nodes: []evo.Node{
							{Position: evo.Position{Layer: 0.0, X: 0.0}},
							{Position: evo.Position{Layer: 0.0, X: 1.0}},
							{Position: evo.Position{Layer: 1.0, X: 0.5}},
						},
						Conns: []evo.Conn{
							{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
							{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
						},
					},
				},
			},
			Child: evo.Genome{
				ID: 0, Fitness: 0, Novelty: 0, Solved: false,
				Encoded: evo.Substrate{
					Nodes: []evo.Node{
						{Position: evo.Position{Layer: 0.0, X: 0.0}},
						{Position: evo.Position{Layer: 0.0, X: 0.5}},
						{Position: evo.Position{Layer: 0.0, X: 1.0}},
						{Position: evo.Position{Layer: 1.0, X: 0.5}},
					},
					Conns: []evo.Conn{
						{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
						{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
					},
				},
			},
		},
		{
			Desc: "extra disjoint node p2",
			Parents: []evo.Genome{
				{
					ID: 2, Fitness: 1.1, Novelty: 2.2, Solved: true,
					Encoded: evo.Substrate{
						Nodes: []evo.Node{
							{Position: evo.Position{Layer: 0.0, X: 0.0}},
							{Position: evo.Position{Layer: 0.0, X: 1.0}},
							{Position: evo.Position{Layer: 1.0, X: 0.5}},
						},
						Conns: []evo.Conn{
							{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
							{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
						},
					},
				},
				{
					ID: 3, Fitness: 1.1, Novelty: 3.3, Solved: false,
					Encoded: evo.Substrate{
						Nodes: []evo.Node{
							{Position: evo.Position{Layer: 0.0, X: 0.0}},
							{Position: evo.Position{Layer: 0.0, X: 0.5}},
							{Position: evo.Position{Layer: 0.0, X: 1.0}},
							{Position: evo.Position{Layer: 1.0, X: 0.5}},
						},
						Conns: []evo.Conn{
							{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
							{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
						},
					},
				},
			},
			Child: evo.Genome{
				ID: 0, Fitness: 0, Novelty: 0, Solved: false,
				Encoded: evo.Substrate{
					Nodes: []evo.Node{
						{Position: evo.Position{Layer: 0.0, X: 0.0}},
						{Position: evo.Position{Layer: 0.0, X: 0.5}},
						{Position: evo.Position{Layer: 0.0, X: 1.0}},
						{Position: evo.Position{Layer: 1.0, X: 0.5}},
					},
					Conns: []evo.Conn{
						{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
						{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
					},
				},
			},
		},
		{
			Desc: "extra excess node p1",
			Parents: []evo.Genome{
				{
					ID: 2, Fitness: 1.1, Novelty: 2.2, Solved: true,
					Encoded: evo.Substrate{
						Nodes: []evo.Node{
							{Position: evo.Position{Layer: 0.0, X: 0.0}},
							{Position: evo.Position{Layer: 0.0, X: 1.0}},
							{Position: evo.Position{Layer: 1.0, X: 0.5}},
							{Position: evo.Position{Layer: 1.0, X: 1.0}},
						},
						Conns: []evo.Conn{
							{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
							{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
						},
					},
				},
				{
					ID: 3, Fitness: 1.1, Novelty: 3.3, Solved: false,
					Encoded: evo.Substrate{
						Nodes: []evo.Node{
							{Position: evo.Position{Layer: 0.0, X: 0.0}},
							{Position: evo.Position{Layer: 0.0, X: 1.0}},
							{Position: evo.Position{Layer: 1.0, X: 0.5}},
						},
						Conns: []evo.Conn{
							{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
							{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
						},
					},
				},
			},
			Child: evo.Genome{
				ID: 0, Fitness: 0, Novelty: 0, Solved: false,
				Encoded: evo.Substrate{
					Nodes: []evo.Node{
						{Position: evo.Position{Layer: 0.0, X: 0.0}},
						{Position: evo.Position{Layer: 0.0, X: 1.0}},
						{Position: evo.Position{Layer: 1.0, X: 0.5}},
						{Position: evo.Position{Layer: 1.0, X: 1.0}},
					},
					Conns: []evo.Conn{
						{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
						{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
					},
				},
			},
		},
		{
			Desc: "extra excess node p2",
			Parents: []evo.Genome{
				{
					ID: 2, Fitness: 1.1, Novelty: 2.2, Solved: true,
					Encoded: evo.Substrate{
						Nodes: []evo.Node{
							{Position: evo.Position{Layer: 0.0, X: 0.0}},
							{Position: evo.Position{Layer: 0.0, X: 1.0}},
							{Position: evo.Position{Layer: 1.0, X: 0.5}},
						},
						Conns: []evo.Conn{
							{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
							{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
						},
					},
				},
				{
					ID: 3, Fitness: 1.1, Novelty: 3.3, Solved: false,
					Encoded: evo.Substrate{
						Nodes: []evo.Node{
							{Position: evo.Position{Layer: 0.0, X: 0.0}},
							{Position: evo.Position{Layer: 0.0, X: 1.0}},
							{Position: evo.Position{Layer: 1.0, X: 0.5}},
							{Position: evo.Position{Layer: 1.0, X: 1.0}},
						},
						Conns: []evo.Conn{
							{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
							{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
						},
					},
				},
			},
			Child: evo.Genome{
				ID: 0, Fitness: 0, Novelty: 0, Solved: false,
				Encoded: evo.Substrate{
					Nodes: []evo.Node{
						{Position: evo.Position{Layer: 0.0, X: 0.0}},
						{Position: evo.Position{Layer: 0.0, X: 1.0}},
						{Position: evo.Position{Layer: 1.0, X: 0.5}},
						{Position: evo.Position{Layer: 1.0, X: 1.0}},
					},
					Conns: []evo.Conn{
						{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
						{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
					},
				},
			},
		},
		{
			Desc: "extra disjoint conn p2",
			Parents: []evo.Genome{
				{
					ID: 2, Fitness: 1.1, Novelty: 2.2, Solved: true,
					Encoded: evo.Substrate{
						Nodes: []evo.Node{
							{Position: evo.Position{Layer: 0.0, X: 0.0}},
							{Position: evo.Position{Layer: 0.0, X: 1.0}},
							{Position: evo.Position{Layer: 1.0, X: 0.5}},
						},
						Conns: []evo.Conn{
							{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
							{Source: evo.Position{Layer: 0.0, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 0.5}},
							{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
						},
					},
				},
				{
					ID: 3, Fitness: 1.1, Novelty: 3.3, Solved: false,
					Encoded: evo.Substrate{
						Nodes: []evo.Node{
							{Position: evo.Position{Layer: 0.0, X: 0.0}},
							{Position: evo.Position{Layer: 0.0, X: 1.0}},
							{Position: evo.Position{Layer: 1.0, X: 0.5}},
						},
						Conns: []evo.Conn{
							{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
							{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
						},
					},
				},
			},
			Child: evo.Genome{
				ID: 0, Fitness: 0, Novelty: 0, Solved: false,
				Encoded: evo.Substrate{
					Nodes: []evo.Node{
						{Position: evo.Position{Layer: 0.0, X: 0.0}},
						{Position: evo.Position{Layer: 0.0, X: 1.0}},
						{Position: evo.Position{Layer: 1.0, X: 0.5}},
					},
					Conns: []evo.Conn{
						{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
						{Source: evo.Position{Layer: 0.0, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 0.5}},
						{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
					},
				},
			},
		},
		{
			Desc: "extra disjoint conn p2",
			Parents: []evo.Genome{
				{
					ID: 2, Fitness: 1.1, Novelty: 2.2, Solved: true,
					Encoded: evo.Substrate{
						Nodes: []evo.Node{
							{Position: evo.Position{Layer: 0.0, X: 0.0}},
							{Position: evo.Position{Layer: 0.0, X: 1.0}},
							{Position: evo.Position{Layer: 1.0, X: 0.5}},
						},
						Conns: []evo.Conn{
							{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
							{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
						},
					},
				},
				{
					ID: 3, Fitness: 1.1, Novelty: 3.3, Solved: false,
					Encoded: evo.Substrate{
						Nodes: []evo.Node{
							{Position: evo.Position{Layer: 0.0, X: 0.0}},
							{Position: evo.Position{Layer: 0.0, X: 1.0}},
							{Position: evo.Position{Layer: 1.0, X: 0.5}},
						},
						Conns: []evo.Conn{
							{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
							{Source: evo.Position{Layer: 0.0, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 0.5}},
							{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
						},
					},
				},
			},
			Child: evo.Genome{
				ID: 0, Fitness: 0, Novelty: 0, Solved: false,
				Encoded: evo.Substrate{
					Nodes: []evo.Node{
						{Position: evo.Position{Layer: 0.0, X: 0.0}},
						{Position: evo.Position{Layer: 0.0, X: 1.0}},
						{Position: evo.Position{Layer: 1.0, X: 0.5}},
					},
					Conns: []evo.Conn{
						{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
						{Source: evo.Position{Layer: 0.0, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 0.5}},
						{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
					},
				},
			},
		},
		{
			Desc: "extra excess conn p2",
			Parents: []evo.Genome{
				{
					ID: 2, Fitness: 1.1, Novelty: 2.2, Solved: true,
					Encoded: evo.Substrate{
						Nodes: []evo.Node{
							{Position: evo.Position{Layer: 0.0, X: 0.0}},
							{Position: evo.Position{Layer: 0.0, X: 1.0}},
							{Position: evo.Position{Layer: 1.0, X: 0.5}},
						},
						Conns: []evo.Conn{
							{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
							{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
							{Source: evo.Position{Layer: 1.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
						},
					},
				},
				{
					ID: 3, Fitness: 1.1, Novelty: 3.3, Solved: false,
					Encoded: evo.Substrate{
						Nodes: []evo.Node{
							{Position: evo.Position{Layer: 0.0, X: 0.0}},
							{Position: evo.Position{Layer: 0.0, X: 1.0}},
							{Position: evo.Position{Layer: 1.0, X: 0.5}},
						},
						Conns: []evo.Conn{
							{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
							{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
						},
					},
				},
			},
			Child: evo.Genome{
				ID: 0, Fitness: 0, Novelty: 0, Solved: false,
				Encoded: evo.Substrate{
					Nodes: []evo.Node{
						{Position: evo.Position{Layer: 0.0, X: 0.0}},
						{Position: evo.Position{Layer: 0.0, X: 1.0}},
						{Position: evo.Position{Layer: 1.0, X: 0.5}},
					},
					Conns: []evo.Conn{
						{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
						{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
						{Source: evo.Position{Layer: 1.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
					},
				},
			},
		},
		{
			Desc: "extra excess conn p2",
			Parents: []evo.Genome{
				{
					ID: 2, Fitness: 1.1, Novelty: 2.2, Solved: true,
					Encoded: evo.Substrate{
						Nodes: []evo.Node{
							{Position: evo.Position{Layer: 0.0, X: 0.0}},
							{Position: evo.Position{Layer: 0.0, X: 1.0}},
							{Position: evo.Position{Layer: 1.0, X: 0.5}},
						},
						Conns: []evo.Conn{
							{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
							{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
						},
					},
				},
				{
					ID: 3, Fitness: 1.1, Novelty: 3.3, Solved: false,
					Encoded: evo.Substrate{
						Nodes: []evo.Node{
							{Position: evo.Position{Layer: 0.0, X: 0.0}},
							{Position: evo.Position{Layer: 0.0, X: 1.0}},
							{Position: evo.Position{Layer: 1.0, X: 0.5}},
						},
						Conns: []evo.Conn{
							{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
							{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
							{Source: evo.Position{Layer: 1.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
						},
					},
				},
			},
			Child: evo.Genome{
				ID: 0, Fitness: 0, Novelty: 0, Solved: false,
				Encoded: evo.Substrate{
					Nodes: []evo.Node{
						{Position: evo.Position{Layer: 0.0, X: 0.0}},
						{Position: evo.Position{Layer: 0.0, X: 1.0}},
						{Position: evo.Position{Layer: 1.0, X: 0.5}},
					},
					Conns: []evo.Conn{
						{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
						{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
						{Source: evo.Position{Layer: 1.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {

			// Cross the parents
			cmp := evo.ByFitness
			crs := &Crosser{EnableProbability: 0.0, Comparison: cmp}
			child, err := crs.Cross(c.Parents...)

			// Check if error was expected
			if t.Run("Error", mock.Error(false, err)) == false {
				t.Fail()
			}
			if err != nil {
				return // Error was received so the rest of the tests don't make sense
			}

			// Compare the actual child to the Expected
			t.Run("Child", testCompareChild(child, c.Child))

		})
	}
}
