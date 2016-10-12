package speciator

import (
	"testing"

	"github.com/klokare/evo"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDynamicSpeciate(t *testing.T) {
	Convey("Given a dynamic speciator and a population", t, func() {
		c := &mockComparer{}
		h := &Dynamic{Static: Static{Comparer: c, CompatibilityThreshold: 3}, CompatibilityModifier: 1.0}
		p := &evo.Population{
			Species: []evo.Species{
				{ID: 1, Example: evo.Substrate{Conns: make([]evo.Conn, 3)}},
				{ID: 2, Example: evo.Substrate{Conns: make([]evo.Conn, 9)}},
			},
			Genomes: []evo.Genome{
				{ID: 1, Encoded: evo.Substrate{Conns: make([]evo.Conn, 3)}},  // Species 1
				{ID: 2, Encoded: evo.Substrate{Conns: make([]evo.Conn, 3)}},  // Species 1
				{ID: 3, Encoded: evo.Substrate{Conns: make([]evo.Conn, 10)}}, // Species 3
				{ID: 4, Encoded: evo.Substrate{Conns: make([]evo.Conn, 10)}}, // Species 3
			},
		}
		Convey("When speciating the population", func() {
			Convey("When the number of species equals the target number of species", func() {
				h.TargetSpecies = 2
				err := h.Speciate(p)
				Convey("There should be no error", func() { So(err, ShouldBeNil) })
				Convey("The compatibility threshold should not change", func() {
					So(h.Static.CompatibilityThreshold, ShouldEqual, 3.0)
				})
			})
			Convey("When the number of species is below the target number", func() {
				h.TargetSpecies = 3
				err := h.Speciate(p)
				Convey("There should be no error", func() { So(err, ShouldBeNil) })
				Convey("The compatibility threshold should be incremented by the modifier", func() {
					So(h.Static.CompatibilityThreshold, ShouldEqual, 2.0)
				})
			})
			Convey("When the number of species is above the target number", func() {
				h.TargetSpecies = 1
				err := h.Speciate(p)
				Convey("There should be no error", func() { So(err, ShouldBeNil) })
				Convey("The compatibility threshold should be decremented by the modifier", func() {
					So(h.Static.CompatibilityThreshold, ShouldEqual, 4.0)
				})
			})
			Convey("When changing the compatibility threshold would take the value too low", func() {
				h.TargetSpecies = 3
				h.Static.CompatibilityThreshold = 1.4
				err := h.Speciate(p)
				Convey("There should be no error", func() { So(err, ShouldBeNil) })
				Convey("The compatibility threshold should be bottom out at the modifier value", func() {
					So(h.Static.CompatibilityThreshold, ShouldEqual, 1.0)
				})
			})
			Convey("When there is an error during comparisons", func() {
				c.error = true
				err := h.Speciate(p)
				Convey("There should be an error", func() { So(err, ShouldNotBeNil) })
			})
		})
	})

}
