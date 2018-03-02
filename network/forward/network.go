package forward

import (
	"fmt"

	"github.com/klokare/evo"
	"gonum.org/v1/gonum/mat"
)

// Layer in the neural network
type Layer struct {
	Activations []evo.Activation // activation functions for neurons in this layer
	Biases      []float64        // bias values for this layer
	Sources     []int            // indexes of source layers
	Weights     []*mat.Dense     // weights to each source layer
}

// Network of layers
type Network struct {
	Layers []Layer
}

// Activate activates the network using the incoming matrix of values
func (net Network) Activate(inputs evo.Matrix) (outputs evo.Matrix, err error) {

	// Prebuild the values with bias values
	values := make([]*mat.Dense, len(net.Layers))
	n, _ := inputs.Dims()

	for i := 1; i < len(values); i++ {
		biases := net.Layers[i].Biases
		tmp := mat.NewDense(n, len(biases), nil)
		for j := 0; j < n; j++ {
			tmp.SetRow(j, biases)
		}
		values[i] = tmp
	}

	// Preallocate the temporary matrices for each layer
	tmp := make([]*mat.Dense, len(values))
	for i := 1; i < len(values); i++ {
		r, c := values[i].Dims()
		tmp[i] = mat.NewDense(r, c, nil)
	}

	// Inputs may not be a dense matrix
	var ok bool
	if values[0], ok = inputs.(*mat.Dense); !ok {
		if tmp, ok := inputs.(mat.Matrix); ok {
			values[0] = mat.DenseCopyOf(tmp)
		} else {
			r, c := inputs.Dims()
			tmp := mat.NewDense(r, c, nil)
			for i := 0; i < r; i++ {
				for j := 0; j < c; j++ {
					tmp.Set(i, j, inputs.At(i, j))
				}
			}
			values[0] = tmp
		}
	}

	// Iterate the layers
	var tgt *mat.Dense
	for i := 1; i < len(net.Layers); i++ {

		// Identify the target
		tgt = values[i]

		// Iterate the sources
		lay := net.Layers[i]
		for j, idx := range lay.Sources {

			// Multiply the source values by the weights
			tmp[i].Mul(values[idx], lay.Weights[j])

			// Add to the values. The target matrix already exists because of bias step above.
			tgt.Add(tgt, tmp[i])
		}

		// Activate the neurons
		r, c := tgt.Dims()
		for j := 0; j < r; j++ {
			for k := 0; k < c; k++ {
				x := tgt.At(j, k)
				tgt.Set(j, k, lay.Activations[k].Activate(x))
			}
		}
		values[i] = tgt
	}

	// Return the output matrix
	outputs = values[len(values)-1]
	return
}

func show(name string, m mat.Matrix) {
	fa := mat.Formatted(m, mat.Prefix("    "), mat.Squeeze())
	fmt.Printf("%s:\na = %v\n\n", name, fa)
}
