package comparer

import (
	"math"
	"testing"

	"github.com/klokare/evo"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDistanceCompare(t *testing.T) {
	Convey("Given a distance comparer and two genomes", t, func() {
		h := &Distance{}
		s1 := evo.Substrate{
			Nodes: []evo.Node{
				{Position: evo.Position{Layer: 0.0}},
				{Position: evo.Position{Layer: 0.2}},
				{Position: evo.Position{Layer: 0.5}},
				{Position: evo.Position{Layer: 1.0}},
			},
			Conns: []evo.Conn{
				{Source: evo.Position{Layer: 0.0}, Target: evo.Position{Layer: 0.2}, Weight: 1.1},
				{Source: evo.Position{Layer: 0.0}, Target: evo.Position{Layer: 1.0}, Weight: 2.1},
				{Source: evo.Position{Layer: 0.2}, Target: evo.Position{Layer: 0.5}, Weight: 3.1},
				{Source: evo.Position{Layer: 0.5}, Target: evo.Position{Layer: 1.0}, Weight: 4.1},
			},
		}
		s2 := evo.Substrate{
			Nodes: []evo.Node{
				{Position: evo.Position{Layer: 0.0}},
				{Position: evo.Position{Layer: 0.5}},
				{Position: evo.Position{Layer: 1.0}},
			},
			Conns: []evo.Conn{
				{Source: evo.Position{Layer: 0.0}, Target: evo.Position{Layer: 1.0}, Weight: 3.3},
				{Source: evo.Position{Layer: 0.5}, Target: evo.Position{Layer: 1.0}, Weight: 4.4},
			},
		}
		Convey("When comparing different nodes only", func() {
			h.DifferentNodesCoefficient = 2.0
			d, err := h.Compare(s1, s2)
			Convey("There should be no error", func() { So(err, ShouldBeNil) })
			Convey("The distance should be the coefficient * the number of different nodes", func() {
				So(d, ShouldEqual, 1*h.DifferentNodesCoefficient)
			})
			Convey("The result should be reciprocal", func() {
				d2, _ := h.Compare(s2, s1)
				So(d2, ShouldEqual, d)
			})
		})

		Convey("When comparing different conns only", func() {
			h.DifferentConnsCoefficient = 2.0
			d, err := h.Compare(s1, s2)
			Convey("There should be no error", func() { So(err, ShouldBeNil) })
			Convey("The distance should be the coefficient * the number of different conns", func() {
				So(d, ShouldEqual, 2*h.DifferentConnsCoefficient)
			})
			Convey("The result should be reciprocal", func() {
				d2, _ := h.Compare(s2, s1)
				So(d2, ShouldEqual, d)
			})
		})

		Convey("When comparing same conns only", func() {
			h.WeightCoefficient = 2.0
			d, err := h.Compare(s1, s2)
			Convey("There should be no error", func() { So(err, ShouldBeNil) })
			Convey("The distance should be the coefficient * the average absolute difference in weights", func() {
				a := math.Abs(2.1 - 3.3)
				b := math.Abs(4.1 - 4.4)
				So(d, ShouldAlmostEqual, ((a+b)/2.0)*h.WeightCoefficient, 0.01)
			})
			Convey("The result should be reciprocal", func() {
				d2, _ := h.Compare(s2, s1)
				So(d2, ShouldEqual, d)
			})
		})
	})
}
