package forward

import (
	"errors"
	"fmt"

	"github.com/klokare/evo"
)

// Known errors
var (
	ErrInvalidInputs = errors.New("too many inputs for activation")
)

// Neuron based on type and activation
type Neuron struct {
	evo.Neuron
	evo.Activation
	Bias float64
}

// Synapse connecting two neurons, adjusting by a weight
type Synapse struct {
	Source, Target int
	Weight         float64
}

// Network definition
type Network struct {
	Inputs, Hidden, Outputs int
	Neurons                 []Neuron
	Synapses                []Synapse
}

// Activate the network using the inputs and return the output. An error is thrown if there are more
// inputs in the incoming inputs arrary than there are in the network. An error will also be thrown
// if a panic is caused when calling any of the neurons' activation functions.
func (z Network) Activate(inputs []float64) (outputs []float64, err error) {
	// Check for error
	if len(inputs) > z.Inputs {
		err = ErrInvalidInputs
		return
	}

	// Recover from panic in neuron activation and return an error
	defer func() {
		if msg := recover(); msg != nil {
			err = fmt.Errorf("panic in activation: %v", msg)
		}
	}()

	// Set the Inputs
	active := make([]bool, len(z.Neurons))
	values := make([]float64, len(z.Neurons))
	copy(values, inputs)
	for i, x := range z.Neurons {
		values[i] += x.Bias
		if x.Neuron == evo.Input {
			active[i] = true
		}
	}

	// Iterate the synapses
	tmp := make([]float64, len(z.Neurons))
	// TODO: Add ability to be recurrent. This might be a set number of iterations or using booleans to see if neurons need to "fire" again. Not straightforward as all downstream neurons need to be rechecked.
	for _, s := range z.Synapses {
		if !active[s.Source] {
			values[s.Source] = z.Neurons[s.Source].Activate(tmp[s.Source] + z.Neurons[s.Source].Bias)
			active[s.Source] = true
		}
		tmp[s.Target] += values[s.Source] * s.Weight
		active[s.Target] = false
	}

	// Return the Outputs
	outputs = make([]float64, z.Outputs)
	offset := len(z.Neurons) - z.Outputs
	for i := 0; i < z.Outputs; i++ {
		if !active[i+offset] {
			values[i+offset] = z.Neurons[i+offset].Activate(tmp[i+offset] + z.Neurons[i+offset].Bias)
			active[i+offset] = true
		}
		outputs[i] = values[i+offset]
	}
	return
}
