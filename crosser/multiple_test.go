package crosser

import (
	"testing"

	"github.com/klokare/evo"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMultipleCross(t *testing.T) {
	Convey("Given a multiple crosser and some parent genomes", t, func() {
		h := &Multiple{}
		gs := []evo.Genome{
			{
				ID:      1,
				Fitness: 1.0,
				Encoded: evo.Substrate{
					Nodes: []evo.Node{
						{Position: evo.Position{Layer: 0.0}, NeuronType: evo.Input, ActivationType: evo.Direct},
						{Position: evo.Position{Layer: 0.5}, NeuronType: evo.Hidden, ActivationType: evo.Direct},
						{Position: evo.Position{Layer: 0.8}, NeuronType: evo.Hidden, ActivationType: evo.Direct},
						{Position: evo.Position{Layer: 1.0}, NeuronType: evo.Output, ActivationType: evo.Direct},
					},
					Conns: []evo.Conn{
						{Source: evo.Position{Layer: 0.0}, Target: evo.Position{Layer: 0.5}, Weight: 1.1, Enabled: true},
						{Source: evo.Position{Layer: 0.0}, Target: evo.Position{Layer: 0.8}, Weight: 2.1, Enabled: true},
						{Source: evo.Position{Layer: 0.5}, Target: evo.Position{Layer: 1.0}, Weight: 3.1, Enabled: true},
					},
				},
			},
			{
				ID:      2,
				Fitness: 2.0,
				Encoded: evo.Substrate{
					Nodes: []evo.Node{
						{Position: evo.Position{Layer: 0.0}, NeuronType: evo.Input, ActivationType: evo.Sigmoid},
						{Position: evo.Position{Layer: 0.5}, NeuronType: evo.Hidden, ActivationType: evo.Sigmoid},
						{Position: evo.Position{Layer: 0.9}, NeuronType: evo.Hidden, ActivationType: evo.Sigmoid},
						{Position: evo.Position{Layer: 1.0}, NeuronType: evo.Output, ActivationType: evo.Sigmoid},
					},
					Conns: []evo.Conn{
						{Source: evo.Position{Layer: 0.0}, Target: evo.Position{Layer: 0.5}, Weight: 1.2, Enabled: true},
						{Source: evo.Position{Layer: 0.0}, Target: evo.Position{Layer: 0.9}, Weight: 2.2, Enabled: true},
						{Source: evo.Position{Layer: 0.5}, Target: evo.Position{Layer: 1.0}, Weight: 3.2, Enabled: true},
					},
				},
			},
			{
				ID:      3,
				Fitness: 0.5,
				Encoded: evo.Substrate{
					Nodes: []evo.Node{
						{Position: evo.Position{Layer: 0.0}, NeuronType: evo.Input, ActivationType: evo.Tanh},
						{Position: evo.Position{Layer: 0.4}, NeuronType: evo.Hidden, ActivationType: evo.Tanh},
						{Position: evo.Position{Layer: 0.5}, NeuronType: evo.Hidden, ActivationType: evo.Tanh},
						{Position: evo.Position{Layer: 1.0}, NeuronType: evo.Output, ActivationType: evo.Tanh},
					},
					Conns: []evo.Conn{
						{Source: evo.Position{Layer: 0.0}, Target: evo.Position{Layer: 0.5}, Weight: 1.3, Enabled: true},
						{Source: evo.Position{Layer: 0.5}, Target: evo.Position{Layer: 1.0}, Weight: 3.3, Enabled: true},
						{Source: evo.Position{Layer: 0.4}, Target: evo.Position{Layer: 1.0}, Weight: 3.3, Enabled: true},
					},
				},
			},
		}

		Convey("When crossing only one parent", func() {
			p1 := gs[0]
			g, err := h.Cross(p1)
			Convey("There should be no error", func() { So(err, ShouldBeNil) })
			Convey("The child's encoded substrate should be a copy of the parent's", func() {
				So(len(g.Encoded.Nodes), ShouldEqual, len(p1.Encoded.Nodes))
				for _, n0 := range g.Encoded.Nodes {
					i := findNodeIdx(p1.Encoded.Nodes, n0.Position)
					So(i, ShouldNotEqual, -1)
					So(n0.NeuronType, ShouldEqual, p1.Encoded.Nodes[i].NeuronType)
					So(n0.ActivationType, ShouldEqual, p1.Encoded.Nodes[i].ActivationType)
				}

				So(len(g.Encoded.Conns), ShouldEqual, len(p1.Encoded.Conns))
				for _, c0 := range g.Encoded.Conns {
					i := findConnIdx(p1.Encoded.Conns, c0.Source, c0.Target)
					So(i, ShouldNotEqual, -1)
					So(c0.Weight, ShouldEqual, p1.Encoded.Conns[i].Weight)
					So(c0.Enabled, ShouldEqual, p1.Encoded.Conns[i].Enabled)
				}
			})
		})

		Convey("When crossing two parents of different fitness", func() {
			p1 := gs[0]
			p2 := gs[1]
			g, err := h.Cross(p1, p2)
			Convey("There should be no error", func() { So(err, ShouldBeNil) })
			Convey("The child's encoded substrate should contain all the nodes and conns of the more fit parent", func() {
				So(len(g.Encoded.Nodes), ShouldEqual, len(p2.Encoded.Nodes))
				for _, n0 := range g.Encoded.Nodes {
					i := findNodeIdx(p2.Encoded.Nodes, n0.Position)
					So(i, ShouldNotEqual, -1)
					So(n0.NeuronType, ShouldEqual, p2.Encoded.Nodes[i].NeuronType)
				}

				So(len(g.Encoded.Conns), ShouldEqual, len(p2.Encoded.Conns))
				for _, c0 := range g.Encoded.Conns {
					i := findConnIdx(p2.Encoded.Conns, c0.Source, c0.Target)
					So(i, ShouldNotEqual, -1)
				}
			})
			Convey("Matching conns' weights should come from either parent", func() {
				for _, c0 := range g.Encoded.Conns {
					i := findConnIdx(p1.Encoded.Conns, c0.Source, c0.Target)
					j := findConnIdx(p2.Encoded.Conns, c0.Source, c0.Target)
					if i != -1 && j != -1 {
						So([]float64{p1.Encoded.Conns[i].Weight, p2.Encoded.Conns[j].Weight}, ShouldContain, c0.Weight)
					}
				}
			})
			Convey("Matching nodes' activations should come from either parent", func() {
				for _, n0 := range g.Encoded.Nodes {
					i := findNodeIdx(p1.Encoded.Nodes, n0.Position)
					j := findNodeIdx(p2.Encoded.Nodes, n0.Position)
					if i != -1 && j != -1 {
						So([]evo.ActivationType{p1.Encoded.Nodes[i].ActivationType, p2.Encoded.Nodes[j].ActivationType}, ShouldContain, n0.ActivationType)
					}
				}
			})
			Convey("Child conns' enabled should be false if disabled in either parent", func() {
				for _, c0 := range g.Encoded.Conns {
					i := findConnIdx(p1.Encoded.Conns, c0.Source, c0.Target)
					j := findConnIdx(p2.Encoded.Conns, c0.Source, c0.Target)
					if i != -1 && j != -1 {
						So(c0.Enabled, ShouldEqual, p1.Encoded.Conns[i].Enabled && p2.Encoded.Conns[j].Enabled)
					}
				}
			})
		})

		Convey("When crossing two parents of equal fitness", func() {
			p1 := gs[0]
			p2 := gs[1]
			p1.Fitness = p2.Fitness
			g, err := h.Cross(p1, p2)
			Convey("There should be no error", func() { So(err, ShouldBeNil) })
			Convey("The child's encoded substrate should contain all the nodes and conns from both parents", func() {
				for _, n := range p1.Encoded.Nodes {
					i := findNodeIdx(g.Encoded.Nodes, n.Position)
					So(i, ShouldNotEqual, -1)
				}
				for _, n := range p2.Encoded.Nodes {
					i := findNodeIdx(g.Encoded.Nodes, n.Position)
					So(i, ShouldNotEqual, -1)
				}
				for _, c := range p1.Encoded.Conns {
					i := findConnIdx(g.Encoded.Conns, c.Source, c.Target)
					So(i, ShouldNotEqual, -1)
				}
				for _, c := range p2.Encoded.Conns {
					i := findConnIdx(g.Encoded.Conns, c.Source, c.Target)
					So(i, ShouldNotEqual, -1)
				}
			})
		})

		Convey("When crossing three parents", func() {
			p1 := gs[0]
			p2 := gs[1]
			p3 := gs[2]
			g, err := h.Cross(p1, p2, p3)
			Convey("There should be no error", func() { So(err, ShouldBeNil) })
			Convey("Matching conns' weights should come from any parent", func() {
				for _, c0 := range g.Encoded.Conns {
					w := make([]float64, 0, 3)
					if i := findConnIdx(p1.Encoded.Conns, c0.Source, c0.Target); i != -1 {
						w = append(w, p1.Encoded.Conns[i].Weight)
					}
					if i := findConnIdx(p2.Encoded.Conns, c0.Source, c0.Target); i != -1 {
						w = append(w, p2.Encoded.Conns[i].Weight)
					}
					if i := findConnIdx(p3.Encoded.Conns, c0.Source, c0.Target); i != -1 {
						w = append(w, p3.Encoded.Conns[i].Weight)
					}
					if len(w) > 1 {
						So(w, ShouldContain, c0.Weight)
					}
				}
			})
			Convey("Matching nodes' activations should come from any parent", func() {
				for _, n0 := range g.Encoded.Nodes {
					a := make([]evo.ActivationType, 0, 3)
					if i := findNodeIdx(p1.Encoded.Nodes, n0.Position); i != -1 {
						a = append(a, p1.Encoded.Nodes[i].ActivationType)
					}
					if i := findNodeIdx(p2.Encoded.Nodes, n0.Position); i != -1 {
						a = append(a, p2.Encoded.Nodes[i].ActivationType)
					}
					if i := findNodeIdx(p3.Encoded.Nodes, n0.Position); i != -1 {
						a = append(a, p3.Encoded.Nodes[i].ActivationType)
					}
					if len(a) > 1 {
						So(a, ShouldContain, n0.ActivationType)
					}
				}
			})
		})

		Convey("When enable probabiliy is set to 1.0", func() {
			h.EnableProbability = 1.0
			g, err := h.Cross(gs...)
			Convey("There should be no error", func() { So(err, ShouldBeNil) })
			Convey("All the childs's connections should be enabled", func() {
				found := false
				for _, c := range g.Encoded.Conns {
					if !c.Enabled {
						found = true
						break
					}
				}
				So(found, ShouldBeFalse)
			})
		})
	})
}
