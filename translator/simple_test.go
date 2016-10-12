package translator

import (
	"testing"

	"github.com/klokare/evo"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSimpleTranslate(t *testing.T) {
	Convey("Given a simple translator and substrate", t, func() {
		h := &Simple{}
		s := evo.Substrate{
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
		Convey("When translating the substrate", func() {
			n, err := h.Translate(s)
			Convey("There should be no error", func() { So(err, ShouldBeNil) })
			Convey("The network should be created", func() { So(n, ShouldNotBeNil) })
		})
	})
}
