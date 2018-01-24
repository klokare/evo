package mutator

import (
	"context"

	"github.com/klokare/evo"
)

// Complexify mutates a genome by adding to its structure
type Complexify struct {
	AddNodeProbability float64
	AddConnProbability float64
	WeightPower        float64
	HiddenActivation   evo.Activation
}

// Mutate a genome by adding nodes or connections
func (m Complexify) Mutate(ctx context.Context, g *evo.Genome) (err error) {
	rng := evo.NewRandom()
	if rng.Float64() < m.AddNodeProbability {
		return m.addNode(rng, m.HiddenActivation, &g.Encoded)
	}
	if rng.Float64() < m.AddConnProbability {
		return m.addConn(rng, m.WeightPower, &g.Encoded)
	}
	return
}

// add node by splitting a connection. will not add node if one already exists.
//
// In the add node mutation, an existing connection is split and the new node placed where the old
// connection used to be. The old connection is disabled and two new connections are added to the
// genome. The new connection leading into the new node receives a weight of 1, and the new
// connection leading out receives the same weight as the old connection. (Stanley, 107)
//
// NOTE: Stanley's version does not use a bias property in nodes. Setting that property to zero is
// the equivalent.
func (m Complexify) addNode(rng evo.Random, act evo.Activation, sub *evo.Substrate) (err error) {

	// Iterate connections randomly
	idxs := rng.Perm(len(sub.Conns))
	for _, idx := range idxs {

		// Identify the connection
		c0 := &sub.Conns[idx]

		// Identify the new node's position on the substrate
		p := evo.Position{
			Layer: (c0.Source.Layer + c0.Target.Layer) / 2.0,
			X:     (c0.Source.X + c0.Target.X) / 2.0,
			Y:     (c0.Source.Y + c0.Target.Y) / 2.0,
			Z:     (c0.Source.Z + c0.Target.Z) / 2.0,
		}

		// Look for an existing node
		var n evo.Node
		fn := false
		for _, n := range sub.Nodes {
			if n.Position == p {
				fn = true
				break
			}
		}
		if !fn {
			// Create the new node
			n = evo.Node{
				Position:   p,
				Neuron:     evo.Hidden,
				Activation: act,
				Bias:       0.0,
			}
		}

		// Identify the connections to this node based on the original connection
		c1 := evo.Conn{Source: c0.Source, Target: n.Position, Weight: 1.0, Enabled: true}
		c2 := evo.Conn{Source: n.Position, Target: c0.Target, Weight: c0.Weight, Enabled: true}

		// Skip the original connetion if either of these new connections exist
		fc := false
		for _, c := range sub.Conns {
			if c.Source == c0.Source && c.Target == n.Position {
				fc = true
				break
			} else if c.Source == n.Position && c.Target == c0.Target {
				fc = true
				break
			}
		}
		if fc {
			continue
		}

		// Add the node and connections
		if !fn {
			sub.Nodes = append(sub.Nodes, n)
		}
		sub.Conns = append(sub.Conns, c1, c2)

		// Disable the original connection
		sub.Conns[idx].Enabled = false
		return
	}
	return
}

// add a connection between two unconnected nodes
//
// In the add connection mutation, a single new connection gene with a random weight is added
// connecting two previously unconnected nodes (Stanley, 107).
func (m Complexify) addConn(rng evo.Random, wp float64, sub *evo.Substrate) (err error) {
	// Randomise node order
	sidxs := rng.Perm(len(sub.Nodes))
	tidxs := rng.Perm(len(sub.Nodes))

	// Iterate source and targets
	for _, sidx := range sidxs {
		src := sub.Nodes[sidx]
		for _, tidx := range tidxs {
			tgt := sub.Nodes[tidx]

			// Nodes must be previously unconnected
			found := false
			for _, c := range sub.Conns {
				if c.Source == src.Position && c.Target == tgt.Position {
					found = true
					break
				}
			}
			if found {
				continue
			}

			// Simple tests for recurrence
			// NOTE: recurrence not implemented in this version
			if src.Position == tgt.Position {
				continue // No self-connection
			} else if src.Neuron == evo.Output {
				continue
			} else if tgt.Neuron == evo.Input {
				continue
			} else if tgt.Layer <= src.Layer {
				continue
			}

			// Create the connection
			c := evo.Conn{
				Source:  src.Position,
				Target:  tgt.Position,
				Weight:  rng.NormFloat64() * wp,
				Enabled: true,
			}
			sub.Conns = append(sub.Conns, c)
			return
		}
	}
	return
}

// WithComplexify adds a configured complexify mutator to the experiment
func WithComplexify(cfg evo.Configurer) evo.Option {
	return func(e *evo.Experiment) (err error) {
		z := new(Complexify)
		if err = cfg.Configure(z); err != nil {
			return
		}

		// Do not continue if there is no chance for mutation
		if z.AddNodeProbability == 0.0 && z.AddConnProbability == 0.0 {
			return
		}

		e.Mutators = append(e.Mutators, z)
		return
	}
}
