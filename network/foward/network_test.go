package forward

import (
	"testing"

	"github.com/klokare/evo"
	"github.com/klokare/network"
)

func TestActivateForward(t *testing.T) {
	t.Run(name, network.XorTest(expectedForward))
}

func TestActivateWrongInputSize(t *testing.T) {
	net := expectedForward
	_, err := net.Activate([]float64{1.0, 1.0, 1.0})
	if err == nil {
		t.Errorf("error expected")
	}
}

func TestActivateInvalidActivation(t *testing.T) {
	net := &Network{
		Inputs:   expectedForward.Inputs,
		Hidden:   expectedForward.Hidden,
		Outputs:  expectedForward.Outputs,
		Neurons:  make([]Neuron, len(expectedForward.Neurons)),
		Synapses: make([]Synapse, len(expectedForward.Synapses)),
	}
	copy(net.Neurons, expectedForward.Neurons)
	copy(net.Synapses, expectedForward.Synapses)
	net.Neurons[3].Activation = 0
	_, err := net.Activate([]float64{1.0, 1.0})
	if err == nil {
		t.Errorf("error expected")
	}
}

func XorEvaluate(net evo.Network) (in [][]float64, out []float64, solved bool) {
	in = [][]float64{{0, 0}, {1, 0}, {0, 1}, {1, 1}}
	out = make([]float64, len(in))
	solved = true // be hopeful :)
	for i, inputs := range in {
		outputs, _ := net.Activate(inputs)
		out[i] = outputs[0]
		if i == 0 || i == 3 {
			solved = solved && outputs[0] >= 0.5
		} else {
			solved = solved && outputs[0] <= 0.5
		}
	}
	return
}
