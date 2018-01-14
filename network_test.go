package evo

import (
	"math"
	"testing"
)

func TestNeuronString(t *testing.T) {

	var cases = []struct {
		Desc   string
		Neuron Neuron
	}{
		{Desc: "unknown", Neuron: 0},
		{Desc: "input", Neuron: Input},
		{Desc: "hidden", Neuron: Hidden},
		{Desc: "output", Neuron: Output},
		{Desc: "unknown", Neuron: 4},
	}

	for _, c := range cases {
		if c.Neuron.String() != c.Desc {
			t.Errorf("incorrect String() value: expected %s, actual %s", c.Desc, c.Neuron.String())
		}
	}
}

func TestActivationString(t *testing.T) {

	var cases = []struct {
		Desc       string
		Activation Activation
	}{
		{Desc: "unknown", Activation: 0},
		{Desc: "direct", Activation: Direct},
		{Desc: "sigmoid", Activation: Sigmoid},
		{Desc: "steepened-sigmoid", Activation: SteepenedSigmoid},
		{Desc: "tanh", Activation: Tanh},
		{Desc: "inverse-abs", Activation: InverseAbs},
		{Desc: "sin", Activation: Sin},
		{Desc: "gauss", Activation: Gauss},
		{Desc: "relu", Activation: ReLU},
		{Desc: "unknown", Activation: 9},
	}

	for _, c := range cases {
		if c.Activation.String() != c.Desc {
			t.Errorf("incorrect String() value: expected %s, actual %s", c.Desc, c.Activation.String())
		}
	}
}

func TestActivation(t *testing.T) {

	sig := 1e-6

	var xvals = []float64{-1e10, -10, -1, -0.1, -0.01, 0, 0.01, 0.1, 1, 10, 1e10}
	var cases = []struct {
		Desc       string
		Activation Activation
		Expected   []float64
	}{
		{
			Desc:       "direct",
			Activation: Direct,
			Expected:   []float64{-1e10, -10, -1, -0.1, -0.01, 0, 0.01, 0.1, 1, 10, 1e10},
		},
		{
			Desc:       "sigmoid",
			Activation: Sigmoid,
			Expected:   []float64{0.000000, 0.000045, 0.268941, 0.475021, 0.497500, 0.500000, 0.502500, 0.524979, 0.731059, 0.999955, 1.000000},
		},
		{
			Desc:       "steepened-sigmoid",
			Activation: SteepenedSigmoid,
			Expected:   []float64{0.000000, 0.000000, 0.007392, 0.379894, 0.487752, 0.500000, 0.512248, 0.620106, 0.992608, 1.000000, 1.000000},
		},
		{
			Desc:       "tanh",
			Activation: Tanh,
			Expected:   []float64{-1.000000, -1.000000, -0.761594, -0.099668, -0.010000, 0.000000, 0.010000, 0.099668, 0.761594, 1.000000, 1.000000},
		},
		{
			Desc:       "inverse-abs",
			Activation: InverseAbs,
			Expected:   []float64{-1.000000, -0.909091, -0.500000, -0.090909, -0.009901, 0.000000, 0.009901, 0.090909, 0.500000, 0.909091, 1.000000},
		},
		{
			Desc:       "sin",
			Activation: Sin,
			Expected:   []float64{0.487506, 0.544021, -0.841471, -0.099833, -0.010000, 0.000000, 0.010000, 0.099833, 0.841471, -0.544021, -0.487506},
		},
		{
			Desc:       "gauss",
			Activation: Gauss,
			Expected:   []float64{0.000000, 0.000000, 0.135335, 0.980199, 0.999800, 1.000000, 0.999800, 0.980199, 0.135335, 0.000000, 0.000000},
		},
		{
			Desc:       "relu",
			Activation: ReLU,
			Expected:   []float64{0, 0, 0, 0, 0, 0.000000, 0.010000, 0.100000, 1.000000, 10.000000, 10000000000.000000},
		},
	}

	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {
			for i, x := range xvals {
				y := c.Activation.Activate(x)
				if math.Abs(c.Expected[i]-y) > sig {
					t.Errorf("invalid activation value for x = %f: expected %f, actual %f", x, c.Expected[i], y)
				}
			}
		})
	}

	// Special case: should panic for unknown activation type
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected panic")
			}
		}()
		act := Activation(0)
		_ = act.Activate(-1)
	}()
}
