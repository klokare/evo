package forward

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"sync/atomic"
	"time"

	"github.com/klokare/evo"
	"gonum.org/v1/gonum/mat"
)

// Known errors
var (
	ErrNoSensors = errors.New("network requires at least 1 input neuron")
	ErrNoOutputs = errors.New("network requires at least 1 output neuron")
)

// Translator transforms substrates into networks
type Translator struct {
	DisableSortCheck bool
}

var tmpid *int64 = new(int64)

// Translate the substrate into a network
func (t Translator) Translate(sub evo.Substrate) (net evo.Network, err error) {

	defer func() {
		msg := recover()
		if msg != nil {
			log.Println("PANIC", msg)
			var f *os.File
			if f, err = os.Create(fmt.Sprintf("/tmp/panic-translate-%d.txt", atomic.AddInt64(tmpid, 1))); err != nil {
				panic(err)
			}
			defer f.Close()
			fmt.Fprintln(f, msg)
			fmt.Fprintln(f)

			for _, n := range sub.Nodes {
				fmt.Fprintln(f, n)
			}
			for _, c := range sub.Conns {
				fmt.Fprintln(f, c)
			}
			f.Close()
			time.Sleep(time.Second * 2)
		}
	}()

	// Sort the substrate to ensure proper ordering during translation
	nodes := make([]evo.Node, len(sub.Nodes))
	copy(nodes, sub.Nodes)
	if !t.DisableSortCheck {
		sort.Slice(nodes, func(i, j int) bool { return nodes[i].Compare(nodes[j]) < 0 })
	}

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
			if tl == len(n2l) {
				log.Println("bad a", "tl", tl, "n2l", n2l)
			}
			for n2l[tl][0].Layer != conn.Target.Layer {
				if tl == len(n2l) {
					log.Println("bad b", "tl", tl, "n2l", n2l)
				}
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
			if sn == len(slay) {
				log.Println("bad", "sn", sn, "slay", slay)
			}
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
	// fmt.Println("decoded")
	// for _, n := range sub.Nodes {
	// 	fmt.Println(n)
	// }
	// for _, c := range sub.Conns {
	// 	fmt.Println(c)
	// }
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
			showMatrix(fmt.Sprintf("weights %d -> %d", src, l), layer.Weights[i])
		}
	}
}
