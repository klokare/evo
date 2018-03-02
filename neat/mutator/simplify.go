package mutator

import (
	"sort"

	"github.com/klokare/evo"
)

// Simplify mutator
type Simplify struct {
	DelNodeProbability float64
	DelConnProbability float64
}

// Mutate a genome by simplifying its structure
func (s Simplify) Mutate(g *evo.Genome) (err error) {
	rng := evo.NewRandom()
	if rng.Float64() < s.DelNodeProbability {
		delNode(rng, &g.Encoded)
	} else if rng.Float64() < s.DelConnProbability {
		delConn(rng, &g.Encoded)
	}
	return
}

// Connection deletion is very simply the deletion of a randomly selected connection, all
// connections are considered to be available for deletion. When a connection is deleted the neurons
// that were at each end of the connection are tested to check if they are no longer connected to by
// other connections, if this is the case then the stranded neuron is also deleted. Note that a more
// thorough cleanup routine could be invoked at this point that cleans up any dead-end structures
// that could not possibly be functional, but this can become complex and so we leave NEAT to
// eliminate such structures naturally.
//
// NOTE: in EVO, the deadend structures are also removed
func delConn(rng evo.Random, sub *evo.Substrate) {

	// Improve search speed
	sort.Slice(sub.Nodes, func(i, j int) bool { return sub.Nodes[i].Compare(sub.Nodes[j]) < 0 })
	sort.Slice(sub.Conns, func(i, j int) bool { return sub.Conns[i].Compare(sub.Conns[j]) < 0 })

	// There are no connections to delete
	if len(sub.Conns) == 0 {
		return
	}

	// Choose a connection at random
	idx := rng.Intn(len(sub.Conns))
	c := sub.Conns[idx]

	// Remove the connection from the substrate
	removeConn(sub, c)
}

// Neuron deletion is slightly more complex. The deletion algorithm attempts to replace neurons with
// connections to maintain any circuits a neuron may have participated in, in further generations
// those connections themselves will be open to deletion. This approach provides NEAT with the
// ability to delete whole structures, not just connections.
//
// Because we replace connected neurons with connections we must be careful which neurons we delete.
// Any neuron with only incoming or only outgoing connections is at a dead-end of a circuit and can
// therefore be safely deleted with all of it's connections. However, a neuron with multiple
// incoming and multiple outgoing connections will require a large number of connections to
// substitute for the loss of the neuron - we must fully connect all of the original neuron's source
// neurons with its target neurons, this could be done but may actually be detrimental since the
// functionality represented by the neuron is now distributed over a number of connections, and this
// cannot easily be reversed. Because of this, such neurons are omitted from the process of
// selecting neurons for deletion.
//
// Neurons with only one incoming or one outgoing connection can be replaced with however many
// connections were on the other side of the neuron, therefore these are candidates for deletion.
func delNode(rng evo.Random, sub *evo.Substrate) {

	// Improve search speed
	sort.Slice(sub.Nodes, func(i, j int) bool { return sub.Nodes[i].Compare(sub.Nodes[j]) < 0 })
	sort.Slice(sub.Conns, func(i, j int) bool { return sub.Conns[i].Compare(sub.Conns[j]) < 0 })

	// Iterate nodes randomly
	idxs := rng.Perm(len(sub.Nodes))
	for _, idx := range idxs {

		// Identify the node
		n := sub.Nodes[idx]

		// Only hidden nodes can be deleted
		if n.Neuron != evo.Hidden {
			continue
		}

		// Count incoming and outgoing connections
		incoming := make([]int, 0, 10)
		outgoing := make([]int, 0, 10)
		for i, c := range sub.Conns {
			if c.Source.Compare(n.Position) == 0 { // outgoing
				outgoing = append(outgoing, i)
			}
			if c.Target.Compare(n.Position) == 0 { // incoming
				incoming = append(incoming, i)
			}
		}

		// There are too many connections on both sides, try a different genome
		if len(incoming) > 1 && len(outgoing) > 1 {
			continue
		}

		// Join the connections
		for _, ic := range incoming {
			for _, oc := range outgoing {
				sub.Conns[ic].Target = sub.Conns[oc].Target
				sub.Conns[oc].Source = sub.Conns[ic].Source
			}
		}
		if len(incoming) == 1 {
			removeConn(sub, sub.Conns[incoming[0]])
		} else if len(outgoing) == 1 {
			removeConn(sub, sub.Conns[outgoing[0]])
		}

		// Remove the node
		removeNode(sub, n)
		break
	}
}

// Removes a connection from the substrate. If this strands nodes, they are removed, too.
func removeConn(sub *evo.Substrate, c evo.Conn) {

	// Remove the connection
	idx := sort.Search(len(sub.Conns), func(i int) bool { return sub.Conns[i].Compare(c) >= 0 })
	if idx < len(sub.Conns) && sub.Conns[idx].Compare(c) == 0 {
		sub.Conns = append(sub.Conns[:idx], sub.Conns[idx+1:]...)
	}

	// Check for stranded neurons
	for _, pos := range []evo.Position{c.Source, c.Target} {

		// Count remaining connections to this node
		found := false
		for _, cx := range sub.Conns {
			if cx.Source.Compare(pos) == 0 || cx.Target.Compare(pos) == 0 {
				found = true
				break
			}
		}
		if found {
			continue // node is not stranded
		}

		// Node is stranded
		for _, n := range sub.Nodes {
			if n.Position.Compare(pos) == 0 {
				removeNode(sub, n)
				break
			}
		}
	}
}

// Removes a node from the substrate. Connections to this node are also removed.
func removeNode(sub *evo.Substrate, n evo.Node) {

	// Only hidden nodes can be removed
	if n.Neuron != evo.Hidden {
		return
	}

	// Find the node
	idx := sort.Search(len(sub.Nodes), func(i int) bool { return sub.Nodes[i].Compare(n) >= 0 })
	if idx < len(sub.Nodes) && sub.Nodes[idx].Compare(n) == 0 {
		sub.Nodes = append(sub.Nodes[:idx], sub.Nodes[idx+1:]...)
	}

	// Remove any connections to this node
	remove := make([]evo.Conn, 0, 10)
	for _, c := range sub.Conns {
		if c.Source.Compare(n.Position) == 0 || c.Target.Compare(n.Position) == 0 {
			remove = append(remove, c)
		}
	}
	for _, c := range remove {
		removeConn(sub, c)
	}
}

// WithComplexify adds a configured complexify mutator to the experiment
func WithSimplify(cfg evo.Configurer) evo.Option {
	return func(e *evo.Experiment) (err error) {
		z := new(Simplify)
		if err = cfg.Configure(z); err != nil {
			return
		}

		// Do not continue if there is no chance for mutation
		if z.DelNodeProbability == 0.0 && z.DelConnProbability == 0.0 {
			return
		}

		e.Mutators = append(e.Mutators, z)
		return
	}
}
