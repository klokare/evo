package forward

import (
	"context"
	"errors"
	"sort"

	"github.com/klokare/evo"
)

// Known errors
var (
	ErrNoSensors = errors.New("network requires at least 1 bias or input neuron")
	ErrNoOutputs = errors.New("network requires at least 1 output neuron")
)

// A Translator creates the neural network from the specied nodes and connections
type Translator struct{}

// Translate the nodes and conns into a new neural network. An error will be thrown if
// there are no sensor (bias or input) or output nodes.
func (t Translator) Translate(ctx context.Context, sub evo.Substrate) (net evo.Network, err error) {

	// Construct the empty Network
	z := &Network{
		Neurons:  make([]Neuron, len(sub.Nodes)),
		Synapses: make([]Synapse, len(sub.Conns)),
	}

	// Put the nodes in order
	ns := make([]evo.Node, len(sub.Nodes))
	copy(ns, sub.Nodes)
	sort.Slice(ns, func(i, j int) bool { return ns[i].Compare(ns[j]) < 0 })

	// Create the neurons
	m := make(map[evo.Position]int, len(ns))
	for i, node := range ns {

		// Create the neuron
		z.Neurons[i] = Neuron{Neuron: node.Neuron, Activation: node.Activation, Bias: node.Bias}

		// Note the original index
		m[node.Position] = i

		// Count the neuron type
		switch node.Neuron {
		case evo.Input:
			z.Inputs++
		case evo.Hidden:
			z.Hidden++
		case evo.Output:
			z.Outputs++
		}
	}

	// Check for errors
	if z.Inputs == 0 {
		err = ErrNoSensors
		return
	} else if z.Outputs == 0 {
		err = ErrNoOutputs
		return
	}

	// Create the synapses
	cs := make([]evo.Conn, len(sub.Conns))
	copy(cs, sub.Conns)
	sort.Slice(cs, func(i, j int) bool { return cs[i].Compare(cs[j]) < 0 })
	for i, conn := range sub.Conns {
		z.Synapses[i] = Synapse{
			Source: m[conn.Source],
			Target: m[conn.Target],
			Weight: conn.Weight,
		}
	}

	// Return the new network
	net = z
	return
}
