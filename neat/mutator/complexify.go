package mutator

import (
	"sort"

	"github.com/klokare/evo"
)

// Complexify mutates a genome by adding to its structure
type Complexify struct {
	AddNodeProbability float64
	AddConnProbability float64
	WeightPower        float64
	MaxWeight          float64
	BiasPower          float64
	MaxBias            float64
	HiddenActivation   evo.Activation
	DisableSortCheck   bool
}

// Mutate a genome by adding nodes or connections
func (m Complexify) Mutate(g *evo.Genome) (err error) {
	rng := evo.NewRandom()
	if rng.Float64() < m.AddNodeProbability {
		return m.addNode(rng, &g.Encoded, !m.DisableSortCheck)
	}
	if rng.Float64() < m.AddConnProbability {
		return m.addConn(rng, &g.Encoded, !m.DisableSortCheck)
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
func (m Complexify) addNode(rng evo.Random, sub *evo.Substrate, check bool) (err error) {

	// Improve search speed
	if check {
		sort.Slice(sub.Nodes, func(i, j int) bool { return sub.Nodes[i].Compare(sub.Nodes[j]) < 0 })
		sort.Slice(sub.Conns, func(i, j int) bool { return sub.Conns[i].Compare(sub.Conns[j]) < 0 })
	}

	// Iterate connections randomly
	idxs := rng.Perm(len(sub.Conns))
	for _, idx := range idxs {

		// Identify the connection
		c0 := &sub.Conns[idx]
		if c0.Locked {
			continue // do not split a locked connection
		}

		// Create the new node
		n := evo.Node{
			Position:   evo.Midpoint(c0.Source, c0.Target),
			Neuron:     evo.Hidden,
			Activation: m.HiddenActivation,
			Bias:       rng.NormFloat64() * m.BiasPower,
		}
		if n.Bias > m.MaxBias {
			n.Bias = m.MaxBias
		} else if n.Bias < -m.MaxBias {
			n.Bias = -m.MaxBias
		}

		// Look for an existing node
		i := sort.Search(len(sub.Nodes), func(i int) bool { return sub.Nodes[i].Compare(n) >= 0 })
		if i < len(sub.Nodes) && sub.Nodes[i].Compare(n) == 0 {
			continue
		}
		sub.Nodes = append(sub.Nodes, n)

		// Identify the connections to this node based on the original connection
		c1 := evo.Conn{Source: c0.Source, Target: n.Position, Weight: 1.0, Enabled: true}
		c2 := evo.Conn{Source: n.Position, Target: c0.Target, Weight: c0.Weight, Enabled: true}
		sub.Conns = append(sub.Conns, c1, c2)

		// Disable the original connection
		sub.Conns[idx].Enabled = false

		// Ensure the order of the substrate
		sort.Slice(sub.Nodes, func(i, j int) bool { return sub.Nodes[i].Compare(sub.Nodes[j]) < 0 })
		sort.Slice(sub.Conns, func(i, j int) bool { return sub.Conns[i].Compare(sub.Conns[j]) < 0 })
		return
	}
	return
}

// add a connection between two unconnected nodes
//
// In the add connection mutation, a single new connection gene with a random weight is added
// connecting two previously unconnected nodes (Stanley, 107).
func (m Complexify) addConn(rng evo.Random, sub *evo.Substrate, check bool) (err error) {

	// Improve search speed
	if check {
		sort.Slice(sub.Conns, func(i, j int) bool { return sub.Conns[i].Compare(sub.Conns[j]) < 0 })
	}

	// Randomise node order
	sidxs := rng.Perm(len(sub.Nodes))
	tidxs := rng.Perm(len(sub.Nodes))

	// Iterate source and targets
	for _, sidx := range sidxs {
		src := sub.Nodes[sidx]
		for _, tidx := range tidxs {
			tgt := sub.Nodes[tidx]

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
				Weight:  rng.NormFloat64() * m.WeightPower,
				Enabled: true,
			}
			if c.Weight > m.MaxWeight {
				c.Weight = m.MaxWeight
			} else if c.Weight < -m.MaxWeight {
				c.Weight = -m.MaxWeight
			}

			// Nodes must be previously unconnected
			idx := sort.Search(len(sub.Conns), func(i int) bool { return sub.Conns[i].Compare(c) >= 0 })
			if idx < len(sub.Conns) && sub.Conns[idx].Compare(c) == 0 {
				continue
			}

			// Append the connection
			sub.Conns = append(sub.Conns, c)
			sort.Slice(sub.Conns, func(i, j int) bool { return sub.Conns[i].Compare(sub.Conns[j]) < 0 })
			return
		}
	}
	return
}
