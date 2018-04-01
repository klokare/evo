package boxes

import (
	"sort"

	"github.com/klokare/evo"
)

// Template returns a substrate for the given resolution suitable for HyperNEAT.
// NOTE: the experiment in the paper does not utilise a hidden layer.
func Template(res int, hidden bool) (enc evo.Substrate) {

	// Create the new substrate
	n := res * res * 2
	if hidden {
		n += res * res
	}
	enc.Nodes = make([]evo.Node, 0, n)

	// Add the nodes
	for i := 0; i < res; i++ {
		x := (float64(i)/float64(res-1))*2.0 - 1.0
		for j := 0; j < res; j++ {
			y := (float64(j)/float64(res-1))*2.0 - 1.0

			// Add the input node
			enc.Nodes = append(enc.Nodes, evo.Node{Position: evo.Position{Layer: 0.0, X: x, Y: y}, Neuron: evo.Input, Activation: evo.Direct})

			// Add the hidden node, if using
			if hidden {
				enc.Nodes = append(enc.Nodes, evo.Node{Position: evo.Position{Layer: 0.5, X: x, Y: y}, Neuron: evo.Hidden, Activation: evo.Sigmoid})
			}

			// Add the input node
			enc.Nodes = append(enc.Nodes, evo.Node{Position: evo.Position{Layer: 1.0, X: x, Y: y}, Neuron: evo.Output, Activation: evo.Sigmoid})
		}
	}

	// Ensure the substrate order and return
	sort.Slice(enc.Nodes, func(i, j int) bool { return enc.Nodes[i].Compare(enc.Nodes[j]) < 0 })
	return
}
