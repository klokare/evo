package hyperneat

import (
	"sort"

	"github.com/klokare/evo"
	"github.com/klokare/evo/neat"
)

// Constants for the output index
const (
	Weight int = iota
	Bias
	LEO
)

// Seeder creates the seed population geared towards Cppns. Each encoded substrate will have 8
// inputs, one for each dimension of the source and target nodes, and 3 ouputs: output weight,
// output enabled check (LEO), and bias value. The bias is used for hidden and output nodes only.
type Seeder struct {
	NumTraits         int
	DisconnectRate    float64
	OutputActivation  evo.Activation
	SeedLocalityLayer bool
	SeedLocalityX     bool
	SeedLocalityY     bool
	SeedLocalityZ     bool
}

// Seed returns the seed genome for a HyperNEAT setup. If SeedLocality<Dim> is set to true then a
// node is added and connected to the appropriate inputs and the LEO output.
func (s Seeder) Seed() (g evo.Genome, err error) {

	// Create the seed genome using the NEAT seeder
	ns := neat.Seeder{
		NumInputs:        8,
		NumOutputs:       3,
		NumTraits:        s.NumTraits,
		DisconnectRate:   s.DisconnectRate,
		OutputActivation: evo.InverseAbs, // need a function that gives us [-1,1].
	}
	if g, err = ns.Seed(); err != nil {
		return
	}

	// Add locality
	enc := g.Encoded
	for i := 0; i < 4; i++ {

		// Skip if locality is not desired for this dimension
		switch i {
		case 0:
			if !s.SeedLocalityLayer {
				continue
			}
		case 1:
			if !s.SeedLocalityX {
				continue
			}
		case 2:
			if !s.SeedLocalityY {
				continue
			}
		case 3:
			if !s.SeedLocalityZ {
				continue
			}
		}

		// Create the new node
		x0 := enc.Nodes[i].Position
		x1 := enc.Nodes[i+4].Position
		n := evo.Node{
			Position:   evo.Midpoint(evo.Midpoint(x0, x1), enc.Nodes[10].Position),
			Neuron:     evo.Hidden,
			Activation: evo.Gauss,
		}
		enc.Nodes = append(enc.Nodes, n)

		// Add the connections
		enc.Conns = append(enc.Conns,
			evo.Conn{
				Source:  evo.Position{Layer: 0.0, X: float64(i) / 7.0},
				Target:  n.Position,
				Enabled: true,
			},
			evo.Conn{
				Source:  evo.Position{Layer: 0.0, X: float64(i+4) / 7.0},
				Target:  n.Position,
				Enabled: true,
			},
			evo.Conn{
				Source:  n.Position,
				Target:  evo.Position{Layer: 1.0, X: 1.0},
				Enabled: true,
			},
		)
	}

	// Ensure sorted substrate and return
	sort.Slice(enc.Nodes, func(i, j int) bool { return enc.Nodes[i].Compare(enc.Nodes[j]) < 0 })
	sort.Slice(enc.Conns, func(i, j int) bool { return enc.Conns[i].Compare(enc.Conns[j]) < 0 })
	g.Encoded = enc
	return
}
