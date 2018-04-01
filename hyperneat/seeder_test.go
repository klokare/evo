package hyperneat

import (
	"testing"

	"github.com/klokare/evo"
	"github.com/klokare/evo/internal/mock"
)

func TestSeederSeed(t *testing.T) {

	// Create the test cases
	const inx = 1.0 / 7.0
	var cases = []struct {
		Desc              string
		SeedLocalityLayer bool
		SeedLocalityX     bool
		SeedLocalityY     bool
		SeedLocalityZ     bool
		Expected          evo.Substrate
	}{
		{
			Desc: "no locality set",
			Expected: evo.Substrate{
				Nodes: []evo.Node{
					// Inputs
					{Position: evo.Position{Layer: 0.0, X: 0 * inx}, Neuron: evo.Input, Activation: evo.Direct},
					{Position: evo.Position{Layer: 0.0, X: 1 * inx}, Neuron: evo.Input, Activation: evo.Direct},
					{Position: evo.Position{Layer: 0.0, X: 2 * inx}, Neuron: evo.Input, Activation: evo.Direct},
					{Position: evo.Position{Layer: 0.0, X: 3 * inx}, Neuron: evo.Input, Activation: evo.Direct},
					{Position: evo.Position{Layer: 0.0, X: 4 * inx}, Neuron: evo.Input, Activation: evo.Direct},
					{Position: evo.Position{Layer: 0.0, X: 5 * inx}, Neuron: evo.Input, Activation: evo.Direct},
					{Position: evo.Position{Layer: 0.0, X: 6 * inx}, Neuron: evo.Input, Activation: evo.Direct},
					{Position: evo.Position{Layer: 0.0, X: 7 * inx}, Neuron: evo.Input, Activation: evo.Direct},

					// Outputs
					{Position: evo.Position{Layer: 1.0, X: 0.0}, Neuron: evo.Output, Activation: evo.InverseAbs},
					{Position: evo.Position{Layer: 1.0, X: 0.5}, Neuron: evo.Output, Activation: evo.InverseAbs},
					{Position: evo.Position{Layer: 1.0, X: 1.0}, Neuron: evo.Output, Activation: evo.InverseAbs},
				},
			},
		},
		{
			Desc:              "layer locality set",
			SeedLocalityLayer: true,
			Expected: evo.Substrate{
				Nodes: []evo.Node{
					// Inputs
					{Position: evo.Position{Layer: 0.0, X: 0 * inx}, Neuron: evo.Input, Activation: evo.Direct},
					{Position: evo.Position{Layer: 0.0, X: 1 * inx}, Neuron: evo.Input, Activation: evo.Direct},
					{Position: evo.Position{Layer: 0.0, X: 2 * inx}, Neuron: evo.Input, Activation: evo.Direct},
					{Position: evo.Position{Layer: 0.0, X: 3 * inx}, Neuron: evo.Input, Activation: evo.Direct},
					{Position: evo.Position{Layer: 0.0, X: 4 * inx}, Neuron: evo.Input, Activation: evo.Direct},
					{Position: evo.Position{Layer: 0.0, X: 5 * inx}, Neuron: evo.Input, Activation: evo.Direct},
					{Position: evo.Position{Layer: 0.0, X: 6 * inx}, Neuron: evo.Input, Activation: evo.Direct},
					{Position: evo.Position{Layer: 0.0, X: 7 * inx}, Neuron: evo.Input, Activation: evo.Direct},

					// Hiddens
					{Position: evo.Midpoint(evo.Midpoint(evo.Position{Layer: 0.0, X: 0 * inx}, evo.Position{Layer: 0.0, X: 4 * inx}), evo.Position{Layer: 1.0, X: 1.0}), Neuron: evo.Hidden, Activation: evo.Gauss},

					// Outputs
					{Position: evo.Position{Layer: 1.0, X: 0.0}, Neuron: evo.Output, Activation: evo.InverseAbs},
					{Position: evo.Position{Layer: 1.0, X: 0.5}, Neuron: evo.Output, Activation: evo.InverseAbs},
					{Position: evo.Position{Layer: 1.0, X: 1.0}, Neuron: evo.Output, Activation: evo.InverseAbs},
				},
			},
		},
		{
			Desc:          "x locality set",
			SeedLocalityX: true,
			Expected: evo.Substrate{
				Nodes: []evo.Node{
					// Inputs
					{Position: evo.Position{Layer: 0.0, X: 0 * inx}, Neuron: evo.Input, Activation: evo.Direct},
					{Position: evo.Position{Layer: 0.0, X: 1 * inx}, Neuron: evo.Input, Activation: evo.Direct},
					{Position: evo.Position{Layer: 0.0, X: 2 * inx}, Neuron: evo.Input, Activation: evo.Direct},
					{Position: evo.Position{Layer: 0.0, X: 3 * inx}, Neuron: evo.Input, Activation: evo.Direct},
					{Position: evo.Position{Layer: 0.0, X: 4 * inx}, Neuron: evo.Input, Activation: evo.Direct},
					{Position: evo.Position{Layer: 0.0, X: 5 * inx}, Neuron: evo.Input, Activation: evo.Direct},
					{Position: evo.Position{Layer: 0.0, X: 6 * inx}, Neuron: evo.Input, Activation: evo.Direct},
					{Position: evo.Position{Layer: 0.0, X: 7 * inx}, Neuron: evo.Input, Activation: evo.Direct},

					// Hiddens
					{Position: evo.Midpoint(evo.Midpoint(evo.Position{Layer: 0.0, X: 1 * inx}, evo.Position{Layer: 0.0, X: 5 * inx}), evo.Position{Layer: 1.0, X: 1.0}), Neuron: evo.Hidden, Activation: evo.Gauss},

					// Outputs
					{Position: evo.Position{Layer: 1.0, X: 0.0}, Neuron: evo.Output, Activation: evo.InverseAbs},
					{Position: evo.Position{Layer: 1.0, X: 0.5}, Neuron: evo.Output, Activation: evo.InverseAbs},
					{Position: evo.Position{Layer: 1.0, X: 1.0}, Neuron: evo.Output, Activation: evo.InverseAbs},
				},
			},
		},
		{
			Desc:          "y locality set",
			SeedLocalityY: true,
			Expected: evo.Substrate{
				Nodes: []evo.Node{
					// Inputs
					{Position: evo.Position{Layer: 0.0, X: 0 * inx}, Neuron: evo.Input, Activation: evo.Direct},
					{Position: evo.Position{Layer: 0.0, X: 1 * inx}, Neuron: evo.Input, Activation: evo.Direct},
					{Position: evo.Position{Layer: 0.0, X: 2 * inx}, Neuron: evo.Input, Activation: evo.Direct},
					{Position: evo.Position{Layer: 0.0, X: 3 * inx}, Neuron: evo.Input, Activation: evo.Direct},
					{Position: evo.Position{Layer: 0.0, X: 4 * inx}, Neuron: evo.Input, Activation: evo.Direct},
					{Position: evo.Position{Layer: 0.0, X: 5 * inx}, Neuron: evo.Input, Activation: evo.Direct},
					{Position: evo.Position{Layer: 0.0, X: 6 * inx}, Neuron: evo.Input, Activation: evo.Direct},
					{Position: evo.Position{Layer: 0.0, X: 7 * inx}, Neuron: evo.Input, Activation: evo.Direct},

					// Hiddens
					{Position: evo.Midpoint(evo.Midpoint(evo.Position{Layer: 0.0, X: 2 * inx}, evo.Position{Layer: 0.0, X: 6 * inx}), evo.Position{Layer: 1.0, X: 1.0}), Neuron: evo.Hidden, Activation: evo.Gauss},

					// Outputs
					{Position: evo.Position{Layer: 1.0, X: 0.0}, Neuron: evo.Output, Activation: evo.InverseAbs},
					{Position: evo.Position{Layer: 1.0, X: 0.5}, Neuron: evo.Output, Activation: evo.InverseAbs},
					{Position: evo.Position{Layer: 1.0, X: 1.0}, Neuron: evo.Output, Activation: evo.InverseAbs},
				},
			},
		},
		{
			Desc:          "z locality set",
			SeedLocalityZ: true,
			Expected: evo.Substrate{
				Nodes: []evo.Node{
					// Inputs
					{Position: evo.Position{Layer: 0.0, X: 0 * inx}, Neuron: evo.Input, Activation: evo.Direct},
					{Position: evo.Position{Layer: 0.0, X: 1 * inx}, Neuron: evo.Input, Activation: evo.Direct},
					{Position: evo.Position{Layer: 0.0, X: 2 * inx}, Neuron: evo.Input, Activation: evo.Direct},
					{Position: evo.Position{Layer: 0.0, X: 3 * inx}, Neuron: evo.Input, Activation: evo.Direct},
					{Position: evo.Position{Layer: 0.0, X: 4 * inx}, Neuron: evo.Input, Activation: evo.Direct},
					{Position: evo.Position{Layer: 0.0, X: 5 * inx}, Neuron: evo.Input, Activation: evo.Direct},
					{Position: evo.Position{Layer: 0.0, X: 6 * inx}, Neuron: evo.Input, Activation: evo.Direct},
					{Position: evo.Position{Layer: 0.0, X: 7 * inx}, Neuron: evo.Input, Activation: evo.Direct},

					// Hiddens
					{Position: evo.Midpoint(evo.Midpoint(evo.Position{Layer: 0.0, X: 3 * inx}, evo.Position{Layer: 0.0, X: 7 * inx}), evo.Position{Layer: 1.0, X: 1.0}), Neuron: evo.Hidden, Activation: evo.Gauss},

					// Outputs
					{Position: evo.Position{Layer: 1.0, X: 0.0}, Neuron: evo.Output, Activation: evo.InverseAbs},
					{Position: evo.Position{Layer: 1.0, X: 0.5}, Neuron: evo.Output, Activation: evo.InverseAbs},
					{Position: evo.Position{Layer: 1.0, X: 1.0}, Neuron: evo.Output, Activation: evo.InverseAbs},
				},
			},
		},
		{
			Desc:              "all localities set",
			SeedLocalityLayer: true,
			SeedLocalityX:     true,
			SeedLocalityY:     true,
			SeedLocalityZ:     true,
			Expected: evo.Substrate{
				Nodes: []evo.Node{
					// Inputs
					{Position: evo.Position{Layer: 0.0, X: 0 * inx}, Neuron: evo.Input, Activation: evo.Direct},
					{Position: evo.Position{Layer: 0.0, X: 1 * inx}, Neuron: evo.Input, Activation: evo.Direct},
					{Position: evo.Position{Layer: 0.0, X: 2 * inx}, Neuron: evo.Input, Activation: evo.Direct},
					{Position: evo.Position{Layer: 0.0, X: 3 * inx}, Neuron: evo.Input, Activation: evo.Direct},
					{Position: evo.Position{Layer: 0.0, X: 4 * inx}, Neuron: evo.Input, Activation: evo.Direct},
					{Position: evo.Position{Layer: 0.0, X: 5 * inx}, Neuron: evo.Input, Activation: evo.Direct},
					{Position: evo.Position{Layer: 0.0, X: 6 * inx}, Neuron: evo.Input, Activation: evo.Direct},
					{Position: evo.Position{Layer: 0.0, X: 7 * inx}, Neuron: evo.Input, Activation: evo.Direct},

					// Hiddens
					{Position: evo.Midpoint(evo.Midpoint(evo.Position{Layer: 0.0, X: 0 * inx}, evo.Position{Layer: 0.0, X: 4 * inx}), evo.Position{Layer: 1.0, X: 1.0}), Neuron: evo.Hidden, Activation: evo.Gauss},
					{Position: evo.Midpoint(evo.Midpoint(evo.Position{Layer: 0.0, X: 1 * inx}, evo.Position{Layer: 0.0, X: 5 * inx}), evo.Position{Layer: 1.0, X: 1.0}), Neuron: evo.Hidden, Activation: evo.Gauss},
					{Position: evo.Midpoint(evo.Midpoint(evo.Position{Layer: 0.0, X: 2 * inx}, evo.Position{Layer: 0.0, X: 6 * inx}), evo.Position{Layer: 1.0, X: 1.0}), Neuron: evo.Hidden, Activation: evo.Gauss},
					{Position: evo.Midpoint(evo.Midpoint(evo.Position{Layer: 0.0, X: 3 * inx}, evo.Position{Layer: 0.0, X: 7 * inx}), evo.Position{Layer: 1.0, X: 1.0}), Neuron: evo.Hidden, Activation: evo.Gauss},

					// Outputs
					{Position: evo.Position{Layer: 1.0, X: 0.0}, Neuron: evo.Output, Activation: evo.InverseAbs},
					{Position: evo.Position{Layer: 1.0, X: 0.5}, Neuron: evo.Output, Activation: evo.InverseAbs},
					{Position: evo.Position{Layer: 1.0, X: 1.0}, Neuron: evo.Output, Activation: evo.InverseAbs},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {

			// Create a new seeder
			s := &Seeder{
				NumTraits:         0,
				DisconnectRate:    0.0,
				OutputActivation:  evo.InverseAbs,
				SeedLocalityLayer: c.SeedLocalityLayer,
				SeedLocalityX:     c.SeedLocalityX,
				SeedLocalityY:     c.SeedLocalityY,
				SeedLocalityZ:     c.SeedLocalityZ,
			}

			// Create the seed genome
			g, err := s.Seed()
			if !t.Run("seeding", mock.Error(false, err)) {
				t.FailNow()
			}

			// Compare substrate
			exp := c.Expected
			act := g.Encoded

			// Compare the nodes
			ens := exp.Nodes
			ans := act.Nodes
			if len(ens) != len(ans) {
				t.Errorf("incorrect number of nodes: expected %d, actual %d", len(ens), len(ans))
			} else {
				for _, en := range ens {
					found := false
					for _, an := range ans {
						if en.Compare(an) == 0 {
							if en.Neuron != an.Neuron {
								t.Errorf("incorrect neuron type: expected %v, actual %v", en.Neuron, an.Neuron)
							}
							if en.Activation != an.Activation {
								t.Errorf("incorrect activation type: expected %v, actual %v", en.Activation, an.Activation)
							}
							found = true
							break
						}
					}
					if !found {
						t.Errorf("expected node not found: %v", en.Position)
					}
				}
			}
		})
	}
}
