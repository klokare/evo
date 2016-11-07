package mutator

import (
	"fmt"
	"math/rand"

	"github.com/klokare/evo"
	"github.com/klokare/random"
)

// A Complexify mutator mutates a genome by adding a node or connection
type Complexify struct {
	AddNodeProbability float64 `evo:"add-node-probability"`
	AddConnProbability float64 `evo:"add-conn-probability"`
	AllowRecurrent     bool    `evo:"allow-recurrent"`
}

func (h Complexify) String() string {
	return fmt.Sprintf("evo.mutator.Complexify{AddNodeProbability: %f, AddConnProbability: %f, AllowRecurrent: %v}",
		h.AddNodeProbability, h.AddConnProbability, h.AllowRecurrent)
}

// Mutate a genome by adding a node or a connection
func (h *Complexify) Mutate(g *evo.Genome) error {
	rng := random.New()
	if rng.Float64() < h.AddNodeProbability {
		h.addNode(rng, g)
	} else if rng.Float64() < h.AddConnProbability {
		h.addConn(rng, g)
	}
	return nil
}

// Add a node to the genome by splitting an existing connection. There must not be an node at the
// target position.
func (h *Complexify) addNode(rng *rand.Rand, g *evo.Genome) {

	// Iterate the connections randomly
	idx := rng.Perm(len(g.Encoded.Conns))
	for _, i := range idx {

		// Figure out the location of the new node
		c0 := g.Encoded.Conns[i]
		p := evo.Position{
			Layer: (c0.Source.Layer + c0.Target.Layer) / 2.0,
			X:     (c0.Source.X + c0.Target.X) / 2.0,
			Y:     (c0.Source.Y + c0.Target.Y) / 2.0,
			Z:     (c0.Source.Z + c0.Target.Z) / 2.0,
		}

		// There is no node this position
		if findNodeIdx(g.Encoded.Nodes, p) == -1 {

			// Disable the original connection
			g.Encoded.Conns[i].Enabled = false

			// Add the new node
			g.Encoded.Nodes = append(g.Encoded.Nodes, evo.Node{
				Position:       p,
				NeuronType:     evo.Hidden,
				ActivationType: evo.SteepenedSigmoid,
			})

			// Add the new connections
			g.Encoded.Conns = append(g.Encoded.Conns, evo.Conn{
				Source:  c0.Source,
				Target:  p,
				Weight:  1.0,
				Enabled: true,
			})

			g.Encoded.Conns = append(g.Encoded.Conns, evo.Conn{
				Source:  p,
				Target:  c0.Target,
				Weight:  c0.Weight,
				Enabled: true,
			})
			break
		}
	}
}

func findNodeIdx(ns []evo.Node, p evo.Position) int {
	for i, n := range ns {
		if n.Position.Compare(p) == 0 {
			return i
		}
	}
	return -1
}

// Add a connection between two unconnected nodes
func (h *Complexify) addConn(rng *rand.Rand, g *evo.Genome) {
	src := rand.Perm(len(g.Encoded.Nodes))
	tgt := rand.Perm(len(g.Encoded.Nodes))

	for _, s := range src {
		for _, t := range tgt {

			// Self-recurrence not allowed
			if s == t {
				continue
			}

			// Senors cannot be targets
			if g.Encoded.Nodes[t].NeuronType == evo.Input || g.Encoded.Nodes[t].NeuronType == evo.Bias {
				continue
			}

			// Feed-forward check if recurrance is not allowed
			if !h.AllowRecurrent && g.Encoded.Nodes[s].Compare(g.Encoded.Nodes[t]) >= 0 {
				continue
			}

			// This is an existing connection
			i := findConnIdx(g.Encoded.Conns,
				g.Encoded.Nodes[s].Position,
				g.Encoded.Nodes[t].Position)
			if i != -1 {
				continue
			}

			// Add the new connection
			g.Encoded.Conns = append(g.Encoded.Conns, evo.Conn{
				Source:  g.Encoded.Nodes[s].Position,
				Target:  g.Encoded.Nodes[t].Position,
				Weight:  rng.NormFloat64(),
				Enabled: true,
			})

			return
		}
	}
}

func findConnIdx(cs []evo.Conn, s, t evo.Position) int {
	for i, c := range cs {
		if c.Source.Compare(s) == 0 && c.Target.Compare(t) == 0 {
			return i
		}
	}
	return -1
}
