package translator

import (
	"math/rand"
	"testing"

	"github.com/klokare/evo"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNetworkNewNetwork(t *testing.T) {
	Convey("Given a valid substrate", t, func() {
		sub := evo.Substrate{
			Nodes: []evo.Node{
				{Position: evo.Position{Layer: 0.0, X: 0.0}, NeuronType: evo.Bias, ActivationType: evo.Direct}, // Use the Activation function to track below
				{Position: evo.Position{Layer: 0.0, X: 0.5}, NeuronType: evo.Input, ActivationType: evo.Tanh},
				{Position: evo.Position{Layer: 0.0, X: 1.0}, NeuronType: evo.Input, ActivationType: evo.InverseAbs},
				{Position: evo.Position{Layer: 0.5, X: 0.5}, NeuronType: evo.Hidden, ActivationType: evo.SteepenedSigmoid},
				{Position: evo.Position{Layer: 1.0, X: 0.5}, NeuronType: evo.Output, ActivationType: evo.Sigmoid},
			},
			Conns: []evo.Conn{
				{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 1.0, Enabled: true}, // Use the weight to track below
				{Source: evo.Position{Layer: 0.0, X: 0.5}, Target: evo.Position{Layer: 0.0, X: 0.5}, Weight: 5.0, Enabled: true},
				{Source: evo.Position{Layer: 0.0, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 2.0, Enabled: false},
				{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 3.0, Enabled: true},
				{Source: evo.Position{Layer: 0.5, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 0.5}, Weight: 4.0, Enabled: true},
			},
		}

		a := Network{
			Biases:  1,
			Inputs:  2,
			Hiddens: 1,
			Outputs: 1,
			Neurons: []Neuron{
				{NeuronType: evo.Bias, ActivationType: evo.Direct}, // Use the Activation function to track below
				{NeuronType: evo.Input, ActivationType: evo.Tanh},
				{NeuronType: evo.Input, ActivationType: evo.InverseAbs},
				{NeuronType: evo.Hidden, ActivationType: evo.SteepenedSigmoid},
				{NeuronType: evo.Output, ActivationType: evo.Sigmoid},
			},
			Synapses: []Synapse{
				{Source: 0, Target: 4, Weight: 1.0},
				{Source: 1, Target: 3, Weight: 5.0},
				{Source: 2, Target: 4, Weight: 3.0},
				{Source: 3, Target: 4, Weight: 4.0},
			},
		}
		Convey("When creating a new network", func() {
			b := NewNetwork(sub)
			Convey("There should be the correct number of bias nodes", func() { So(b.Biases, ShouldEqual, a.Biases) })
			Convey("There should be the correct number of input nodes", func() { So(b.Inputs, ShouldEqual, a.Inputs) })
			Convey("There should be the correct number of hidden nodes", func() { So(b.Hiddens, ShouldEqual, a.Hiddens) })
			Convey("There should be the correct number of output nodes", func() { So(b.Outputs, ShouldEqual, a.Outputs) })
			Convey("There should be the correct number of neurons", func() { So(len(b.Neurons), ShouldEqual, len(a.Neurons)) })
			Convey("The neurons should be in order", func() {
				for i := 0; i < len(a.Neurons); i++ {
					So(b.Neurons[i].ActivationType, ShouldEqual, a.Neurons[i].ActivationType)
				}
			})
			Convey("There should be the correct number of synapses", func() { So(len(b.Synapses), ShouldEqual, len(a.Synapses)) })
			Convey("The synapses should be in order", func() {
				for i := 0; i < len(a.Synapses); i++ {
					So(b.Synapses[i].Weight, ShouldEqual, a.Synapses[i].Weight)
				}
			})
		})
		Convey("When creating a new network with randomly ordered substrate nodes and conns", func() {
			sub2 := evo.Substrate{
				Nodes: make([]evo.Node, 0, len(sub.Nodes)),
				Conns: make([]evo.Conn, 0, len(sub.Conns)),
			}
			idxs := rand.Perm(len(sub.Nodes))
			for _, i := range idxs {
				sub2.Nodes = append(sub2.Nodes, sub.Nodes[i])
			}
			idxs = rand.Perm(len(sub.Conns))
			for _, i := range idxs {
				sub2.Conns = append(sub2.Conns, sub.Conns[i])
			}
			b := NewNetwork(sub2)
			Convey("There should be the correct number of neurons", func() { So(len(b.Neurons), ShouldEqual, len(a.Neurons)) })
			Convey("The neurons should be in order", func() {
				for i := 0; i < len(a.Neurons); i++ {
					So(b.Neurons[i].ActivationType, ShouldEqual, a.Neurons[i].ActivationType)
				}
			})
			Convey("There should be the correct number of synapses", func() { So(len(b.Synapses), ShouldEqual, len(a.Synapses)) })
			Convey("The synapses should be in order", func() {
				for i := 0; i < len(a.Synapses); i++ {
					So(b.Synapses[i].Weight, ShouldEqual, a.Synapses[i].Weight)
				}
			})
		})
	})
}

func TestNetworkActivate(t *testing.T) {
	Convey("Given a network and inputs", t, func() {
		a := Network{
			Biases:  1,
			Inputs:  2,
			Hiddens: 1,
			Outputs: 1,
			Neurons: []Neuron{
				{NeuronType: evo.Bias, ActivationType: evo.Direct}, // Use the Activation function to track below
				{NeuronType: evo.Input, ActivationType: evo.Direct},
				{NeuronType: evo.Input, ActivationType: evo.Direct},
				{NeuronType: evo.Hidden, ActivationType: evo.Sigmoid},
				{NeuronType: evo.Output, ActivationType: evo.Sigmoid},
			},
			Synapses: []Synapse{
				{Source: 0, Target: 3, Weight: 2.79},
				{Source: 0, Target: 4, Weight: 5.89},
				{Source: 1, Target: 3, Weight: -7.00},
				{Source: 1, Target: 4, Weight: -3.70},
				{Source: 2, Target: 3, Weight: -12.0},
				{Source: 2, Target: 4, Weight: -3.90},
				{Source: 3, Target: 4, Weight: -9.3},
			},
		}
		inputs := [][]float64{
			{0.0, 0.0}, {0.0, 1.0}, {1.0, 0.0}, {1.0, 1.0},
		}
		expected := []float64{0.05, 0.88, 0.89, 0.15}
		Convey("When activating the network", func() {
			actual := make([][]float64, len(expected))
			for i, in := range inputs {
				actual[i] = a.Activate(in)
			}
			Convey("The actual values should equal the expected", func() {
				for i := 0; i < len(expected); i++ {
					So(actual[i][0], ShouldAlmostEqual, expected[i], 0.01)
				}
			})
		})
	})
}
