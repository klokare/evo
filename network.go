package evo

// A Network is a black box that processes a set of inputs and returns a set of outputs
type Network interface {
	Activate([]float64) []float64
}

// NeuronType defines the type of neuron to create
type NeuronType byte

// Predfined neuron types
const (
	Bias NeuronType = iota + 1
	Input
	Hidden
	Output
)

// ActivationType defines the type of activation function to use for the neuron
type ActivationType byte

// Predfined activation types
const (
	Direct ActivationType = iota + 1
	Sigmoid
	SteepenedSigmoid
	Tanh
	InverseAbs
)

func (a ActivationType) String() string {
	switch a {
	case Direct:
		return "direct"
	case Sigmoid:
		return "sigmoid"
	case SteepenedSigmoid:
		return "steepened-sigmoid"
	case Tanh:
		return "tanh"
	case InverseAbs:
		return "inverse-abs"
	default:
		return "unknown"
	}
}

// ActivationTypes is a list of available activations.
var (
	ActivationTypes = []ActivationType{Direct, Sigmoid, SteepenedSigmoid}
)
