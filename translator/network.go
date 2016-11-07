package translator

import (
	"math"
	"sort"

	"github.com/klokare/evo"
)

// A Neuron collects input signals and provides the source for activation
type Neuron struct {
	evo.NeuronType
	evo.ActivationType
}

// A Synapse connects two neurons with a weighted connection
type Synapse struct {
	Source, Target int
	Weight         float64
}

// A Network is the collection of neurons and synapses
type Network struct {
	Biases, Inputs, Hiddens, Outputs int
	Neurons                          []Neuron
	Synapses                         []Synapse
}

// NewNetwork creates a new network based on the given substrate
func NewNetwork(sub evo.Substrate) *Network {

	// Map the nodes and create the neurons
	sort.Sort(sub.Nodes)
	var b, i, h, o int
	m := make(map[evo.Position]int, len(sub.Nodes))
	ns := make([]Neuron, len(sub.Nodes))
	for z, n := range sub.Nodes {
		switch n.NeuronType {
		case evo.Bias:
			b++
		case evo.Input:
			i++
		case evo.Hidden:
			h++
		case evo.Output:
			o++
		}
		m[n.Position] = z
		ns[z] = Neuron{NeuronType: n.NeuronType, ActivationType: n.ActivationType}
	}

	// Create the synapses
	sort.Sort(sub.Conns)
	ss := make([]Synapse, 0, len(sub.Conns))
	for _, c := range sub.Conns {
		if c.Enabled {
			ss = append(ss, Synapse{
				Source: m[c.Source],
				Target: m[c.Target],
				Weight: c.Weight,
			})
		}
	}

	// Return the new network
	return &Network{
		Biases:   b,
		Inputs:   i,
		Hiddens:  h,
		Outputs:  o,
		Neurons:  ns,
		Synapses: ss,
	}
}

// Activate processes the inputs with the network and returns the outputs
func (net Network) Activate(inputs []float64) []float64 {

	// Set the sensors
	v := make([]float64, len(net.Neurons))
	for i := 0; i < net.Biases; i++ {
		v[i] = 1.0
	}
	for i := 0; i < net.Inputs; i++ {
		v[net.Biases+i] = inputs[i]
	}

	// Iterate the synapses
	for _, s := range net.Synapses {
		v[s.Target] += activate(net.Neurons[s.Source].ActivationType, v[s.Source]) * s.Weight
	}
	// Copy the outputs
	o := make([]float64, net.Outputs)
	offset := net.Biases + net.Inputs + net.Hiddens
	for i := 0; i < len(o); i++ {
		o[i] = activate(net.Neurons[offset+i].ActivationType, v[offset+i])
	}
	return o
}

func activate(a evo.ActivationType, x float64) float64 {
	switch a {
	case evo.Direct:
		return x
	case evo.Sigmoid:
		return 1.0 / (1.0 + math.Exp(-x))
	case evo.SteepenedSigmoid:
		return 1.0 / (1.0 + math.Exp(-4.9*x))
	case evo.Tanh:
		return math.Tanh(0.9 * x)
	case evo.InverseAbs:
		return x / (1.0 + math.Abs(x))
	default:
		panic("Unknown activation type")
	}
}
