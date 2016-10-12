package evo

import (
	"math"
	"testing"

	"github.com/klokare/evo"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSeed(t *testing.T) {
	Convey("Given the starting parameters", t, func() {
		n := 10
		i := 2
		o := 2

		ns := []evo.Node{
			{Position: evo.Position{Layer: 0.0, X: 0.0}, NeuronType: evo.Bias, ActivationType: evo.Direct},
			{Position: evo.Position{Layer: 0.0, X: 0.5}, NeuronType: evo.Input, ActivationType: evo.Direct},
			{Position: evo.Position{Layer: 0.0, X: 1.0}, NeuronType: evo.Input, ActivationType: evo.Direct},
			{Position: evo.Position{Layer: 1.0, X: 0.0}, NeuronType: evo.Output, ActivationType: evo.SteepenedSigmoid},
			{Position: evo.Position{Layer: 1.0, X: 1.0}, NeuronType: evo.Output, ActivationType: evo.SteepenedSigmoid},
		}

		cs := []evo.Conn{
			{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.0}, Enabled: true},
			{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 1.0}, Enabled: true},
			{Source: evo.Position{Layer: 0.0, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 0.0}, Enabled: true},
			{Source: evo.Position{Layer: 0.0, X: 0.5}, Target: evo.Position{Layer: 1.0, X: 1.0}, Enabled: true},
			{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.0}, Enabled: true},
			{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 1.0}, Enabled: true},
		}

		Convey("When creating the population", func() {
			p := Seed(n, i, o)

			Convey("The population generation should be 1", func() {
				// Though technially not necessary, it does make for a more human-readable event log
				So(p.Generation, ShouldEqual, 1)
			})

			Convey("The population should have 1 species", func() {
				So(len(p.Species), ShouldEqual, 1)
			})
			Convey("The species should have ID 1", func() {
				So(p.Species[0].ID, ShouldEqual, 1)
			})

			Convey("The population should have the correct number of genomes", func() {
				So(len(p.Genomes), ShouldEqual, n)
			})
			Convey("Each genome should have a unique ID starting with 1", func() {
				cnts := make(map[int]int, n)
				for _, g := range p.Genomes {
					cnts[g.ID]++
				}
				So(cnts[0], ShouldEqual, 0)
				for _, cnt := range cnts {
					So(cnt, ShouldEqual, 1)
				}
			})
			Convey("The genomes have the correct nodes", func() {
				for _, g := range p.Genomes {
					s := g.Encoded
					So(len(s.Nodes), ShouldEqual, len(ns))
					for i, n := range ns {
						So(s.Nodes[i].Compare(n), ShouldEqual, 0)
						So(s.Nodes[i].NeuronType, ShouldEqual, n.NeuronType)
						So(s.Nodes[i].ActivationType, ShouldEqual, n.ActivationType)
					}

				}
			})
			Convey("The genomes should have the correct connections", func() {
				for _, g := range p.Genomes {
					s := g.Encoded
					So(len(s.Conns), ShouldEqual, len(cs))
					for i, c := range cs {
						So(s.Conns[i].Compare(c), ShouldEqual, 0)
						So(s.Conns[i].Enabled, ShouldEqual, c.Enabled)
					}
				}
			})

			Convey("The weights of all connections should have a mean of 0 and a standard deviation of 1", func() {
				v := make([]float64, 0, n*i*o*100)
				for j := 0; j < 10; j++ { // need more samples
					p = Seed(n, i, o)
					for _, g := range p.Genomes {
						for _, c := range g.Encoded.Conns {
							v = append(v, c.Weight)
						}
					}
				}

				m := 0.0
				for _, x := range v {
					m += x
				}
				m /= float64(len(v))

				s := 0.0
				for _, x := range v {
					s += (x - m) * (x - m)
				}
				s = math.Sqrt(s / float64(len(v)))
				So(m, ShouldAlmostEqual, 0.0, 0.1)
				So(s, ShouldAlmostEqual, 1.0, 0.1)
			})
		})
	})
}
