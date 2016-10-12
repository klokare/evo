package speciator

import (
	"fmt"
	"math"
	"testing"

	"github.com/klokare/evo"
	. "github.com/smartystreets/goconvey/convey"
)

func TestStaticSpeciate(t *testing.T) {
	Convey("Give a static speciator and a population", t, func() {
		c := &mockComparer{}
		h := &Static{
			Comparer:               c,
			CompatibilityThreshold: 0.5,
		}
		p := &evo.Population{
			Species: []evo.Species{
				{ID: 1, Example: evo.Substrate{Conns: make([]evo.Conn, 3)}},
				{ID: 2, Example: evo.Substrate{Conns: make([]evo.Conn, 5)}},
			},
			Genomes: []evo.Genome{
				{ID: 1, Encoded: evo.Substrate{Conns: make([]evo.Conn, 3)}}, // Species 1
				{ID: 2, Encoded: evo.Substrate{Conns: make([]evo.Conn, 3)}}, // Species 1
				{ID: 3, Encoded: evo.Substrate{Conns: make([]evo.Conn, 6)}}, // Species 3
				{ID: 4, Encoded: evo.Substrate{Conns: make([]evo.Conn, 6)}}, // Species 3
			},
		}
		Convey("When speciating the population", func() {
			Convey("When there is no error during comparisons", func() {
				err := h.Speciate(p)
				Convey("There should be no error", func() { So(err, ShouldBeNil) })
				Convey("A new species should be added for genome 1", func() {
					found := false
					for _, s := range p.Species {
						if s.ID == 3 {
							found = true
							break
						}
					}
					So(found, ShouldBeTrue)
				})
				Convey("All genomes should be assigned to the correct species", func() {
					for _, g := range p.Genomes {
						if g.ID < 3 {
							So(g.SpeciesID, ShouldEqual, 1)
						} else {
							So(g.SpeciesID, ShouldEqual, 3)
						}
					}
				})
				Convey("Species 2 is empty and should be removed", func() {
					found := false
					for _, s := range p.Species {
						if s.ID == 2 {
							found = true
							break
						}
					}
					So(found, ShouldBeFalse)
				})
				Convey("Remaining species should have a example with the correct complexity", func() {
					g2s := map[int]int{
						1: 3,
						3: 6,
					}
					for _, s := range p.Species {
						So(s.Example.Complexity(), ShouldEqual, g2s[s.ID])
					}
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

type mockComparer struct {
	called    bool
	calledIDs []int
	error     bool
}

func (h *mockComparer) Compare(s1, s2 evo.Substrate) (float64, error) {
	h.called = true
	d := math.Abs(float64(len(s1.Conns) - len(s2.Conns)))
	if h.error {
		return d, fmt.Errorf("error while comparing")
	}
	return d, nil
}
