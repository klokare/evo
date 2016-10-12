package transcriber

import (
	"math/rand"
	"testing"

	"github.com/klokare/evo"
	. "github.com/smartystreets/goconvey/convey"
)

func TestevoTranscribe(t *testing.T) {
	Convey("Given a evo transcriber and a substrate", t, func() {
		h := &NEAT{}
		enc := evo.Substrate{
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

		Convey("When transcribing the substrate", func() {

			// Randomise the substrate to make sure the transcriber puts them back in the correct order
			sub2 := evo.Substrate{
				Nodes: make([]evo.Node, 0, len(enc.Nodes)),
				Conns: make([]evo.Conn, 0, len(enc.Conns)),
			}
			idxs := rand.Perm(len(enc.Nodes))
			for _, i := range idxs {
				sub2.Nodes = append(sub2.Nodes, enc.Nodes[i])
			}
			idxs = rand.Perm(len(enc.Conns))
			for _, i := range idxs {
				sub2.Conns = append(sub2.Conns, enc.Conns[i])
			}

			dec, err := h.Transcribe(sub2)

			Convey("There should be no error", func() { So(err, ShouldBeNil) })

			Convey("All the nodes should be present and in the correct order", func() {
				So(len(dec.Nodes), ShouldEqual, len(enc.Nodes))
				for i, n1 := range enc.Nodes {
					So(dec.Nodes[i].ActivationType, ShouldEqual, n1.ActivationType)
				}
			})

			Convey("Only enabled connections should be present and in the correct order", func() {
				So(len(dec.Conns), ShouldEqual, len(enc.Conns)-1)
				for _, c1 := range enc.Conns {
					if c1.Enabled {
						found := false
						for _, c2 := range dec.Conns {
							if c2.Compare(c1) == 0 {
								found = true
								So(c2.Weight, ShouldEqual, c1.Weight)
							}
						}
						So(found, ShouldBeTrue)
					}
				}
			})
		})
	})
}
