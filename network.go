package evo

import "math"

// A Network provides the ability to process a set of inputs and returns the outputs
type Network interface {
	Activate(Matrix) (Matrix, error)
}

// Neuron is the type of neuron to create within the network
type Neuron byte

// Complete list of neuron types
const (
	Input Neuron = iota + 1
	Hidden
	Output
)

func (n Neuron) String() string {
	switch n {
	case Input:
		return "input"
	case Hidden:
		return "hidden"
	case Output:
		return "output"
	default:
		return "unknown"
	}
}

// Activation is the type of activation function to use with the neuron
type Activation byte

// Known list of activation types
const (
	Direct Activation = iota + 1
	Sigmoid
	SteepenedSigmoid
	Tanh
	InverseAbs
	Sin
	Gauss
	ReLU
)

func (a Activation) String() string {
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
	case Sin:
		return "sin"
	case Gauss:
		return "gauss"
	case ReLU:
		return "relu"
	default:
		return "unknown"
	}
}

// Activate the neuron using the appropriate transformation function.
func (a Activation) Activate(x float64) float64 {
	switch a {
	case Direct:
		return x
	case Sigmoid:
		return 1.0 / (1.0 + math.Exp(-x))
	case SteepenedSigmoid:
		return 1.0 / (1.0 + math.Exp(-4.9*x))
	case Tanh:
		return math.Tanh(x)
	case InverseAbs:
		return x / (1.0 + math.Abs(x))
	case Sin:
		return math.Sin(x)
	case Gauss:
		return math.Exp(-2.0 * x * x)
	case ReLU:
		if x > 0 {
			return x
		}
		return 0
	default:
		panic("unknown activation")
	}
}

// Activations provides map of activation functions by name
var Activations = map[string]Activation{
	"direct":           Direct,
	"sigmoid":          Sigmoid,
	"steepenedsigmoid": SteepenedSigmoid,
	"tanh":             Tanh,
	"inverseabs":       InverseAbs,
	"sin":              Sin,
	"gauss":            Gauss,
	"relu":             ReLU,
}
