package forward

import (
	"errors"
	"fmt"
	"sort"

	"github.com/klokare/evo"
	"gonum.org/v1/gonum/mat"
)

// Known errors
var (
	ErrNoSensors = errors.New("network requires at least 1 input neuron")
	ErrNoOutputs = errors.New("network requires at least 1 output neuron")
)

// Translator transforms substrates into networks
type Translator struct{}

// Translate the substrate into a network
func (t Translator) Translate(sub evo.Substrate) (net evo.Network, err error) {

	// Sort the substrate to ensure proper ordering during translation
	nodes := make([]evo.Node, len(sub.Nodes))
	copy(nodes, sub.Nodes)
	sort.Slice(nodes, func(i, j int) bool { return nodes[i].Compare(nodes[j]) < 0 })

	conns := make([]evo.Conn, len(sub.Conns))
	copy(conns, sub.Conns)

	// Map the nodes to layers
	n2l := make([][]*evo.Node, 0, 10)

	var ni, no int
	var lay []*evo.Node
	l := -1
	for i, node := range nodes {
		if lay == nil || lay[0].Layer < nodes[i].Layer {
			lay = make([]*evo.Node, 0, 10) // Better estimate of size to prevent multiple allocations?
			n2l = append(n2l, lay)
			l++
		}
		lay = append(lay, &nodes[i])
		n2l[l] = lay
		if node.Neuron == evo.Input {
			ni++
		} else if node.Neuron == evo.Output {
			no++
		}
	}

	// Check for errors
	if ni == 0 {
		return nil, ErrNoSensors
	} else if no == 0 {
		return nil, ErrNoOutputs
	}

	// Create the network's layers
	layers := make([]Layer, len(n2l))
	for i, ns := range n2l {

		// Create the layer
		lay := Layer{Activations: make([]evo.Activation, len(ns))}

		// Transfer the activations
		for j := 0; j < len(ns); j++ {
			lay.Activations[j] = ns[j].Activation
		}

		// Add the biases
		if i > 0 {
			lay.Biases = make([]float64, len(ns))
			for j := 0; j < len(ns); j++ {
				lay.Biases[j] = ns[j].Bias
			}
		}

		// Transfer the activations
		layers[i] = lay
	}
	// Order the connections by target layer then compare
	sort.Slice(conns, func(i, j int) bool {
		a := conns[i]
		b := conns[j]
		if a.Target.Layer < b.Target.Layer {
			return true
		}
		if a.Target.Layer == b.Target.Layer {
			if a.Source.Layer < b.Source.Layer {
				return true
			}
			if a.Source.Layer == b.Source.Layer {
				return a.Compare(b) < 0
			}
		}
		return false
	})

	// Iterate the connections
	var w *mat.Dense
	var sl, tl, sn, tn int
	slay := n2l[0]
	tlay := n2l[0]
	for _, conn := range conns {

		// Skip disabled connections
		if !conn.Enabled {
			continue
		}

		// Advance the target layer as necessary
		if n2l[tl][0].Layer != conn.Target.Layer {

			// Find the source layer
			for n2l[tl][0].Layer != conn.Target.Layer {
				tl++
			}
			tlay = n2l[tl]

			// Reset the target variables
			sl = 0
			slay = n2l[0]
			sn, tn = 0, 0
			w = nil
		}

		// Advance the source layer if necessary
		if w == nil || n2l[sl][0].Layer != conn.Source.Layer {

			// Find the source layer
			for n2l[sl][0].Layer < conn.Source.Layer {
				sl++
			}
			slay = n2l[sl]

			// Reset target node index
			sn, tn = 0, 0

			// Create the source
			w = mat.NewDense(len(slay), len(tlay), nil)
			layers[tl].Sources = append(layers[tl].Sources, sl)
			layers[tl].Weights = append(layers[tl].Weights, w)
		}

		// Advance to the source node
		for slay[sn].Position.Compare(conn.Source) < 0 {
			sn++
			tn = 0
		}

		// Advance to the target node
		for tlay[tn].Position.Compare(conn.Target) < 0 {
			tn++
		}

		// Set the connection weight
		w.Set(sn, tn, conn.Weight)
	}

	// Return the new network
	// showNetwork("general", Network{Layers: layers})
	return Network{Layers: layers}, nil
}

func showNetwork(name string, n Network) {
	fmt.Println(name)
	for l, layer := range n.Layers {
		fmt.Printf("... Layer %d:\n", l)
		fmt.Printf("... activations: %v\n", layer.Activations)
		fmt.Printf("... biases: %v\n", layer.Biases)
		fmt.Printf("... Sources:\n")
		for i, src := range layer.Sources {
			show(fmt.Sprintf("weights %d -> %d", src, l), layer.Weights[i])
		}
	}
}

// WithTranslator configures the experiment to use this translator
func WithTranslator() evo.Option {
	return func(e *evo.Experiment) error {
		e.Translator = Translator{}
		return nil
	}
}
