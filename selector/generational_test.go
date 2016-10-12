package selector

import (
	"testing"

	"github.com/klokare/evo"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGenerationalSelect(t *testing.T) {
	Convey("Given a generational selector and a population", t, func() {
		h := &Generational{
			MaxStagnation:               4,
			SurvivalRate:                0.7,
			MutateOnlyProbability:       0.2,
			InterspeciesMateProbability: 0.1,
		}
		p := evo.Population{
			Species: []evo.Species{
				{ID: 1, Fitness: 10.0, Stagnation: 5}, // Stagnant but best
				{ID: 2, Fitness: 5.0, Stagnation: 5},  // Stagnant
				{ID: 3, Fitness: 8.0, Stagnation: 2},
				{ID: 4, Fitness: 4.0, Stagnation: 2},
			},
			Genomes: []evo.Genome{
				{ID: 1, SpeciesID: 1, Fitness: 10.0},
				{ID: 2, SpeciesID: 1, Fitness: 8.0},
				{ID: 3, SpeciesID: 1, Fitness: 12.0}, // Overall Best
				{ID: 4, SpeciesID: 2, Fitness: 2.0},
				{ID: 5, SpeciesID: 2, Fitness: 7.0},
				{ID: 6, SpeciesID: 3, Fitness: 7.0},
				{ID: 7, SpeciesID: 3, Fitness: 9.0}, // Best for species
				{ID: 8, SpeciesID: 3, Fitness: 8.0},
				{ID: 9, SpeciesID: 4, Fitness: 4.0}, // Best for species
			},
		}

		Convey("When selecting genomes", func() {
			keep, gss, err := h.Select(p)
			Convey("There should be no error", func() { So(err, ShouldBeNil) })

			Convey("Genomes (except best) from stagnant species should not appear in keep", func() {
				ids := []int{3, 7, 9}
				So(len(keep), ShouldEqual, len(ids))
				for _, g := range keep {
					So(ids, ShouldContain, g.ID)
				}
			})
			Convey("Genomes (except best) from stagnant species should not appear in parents", func() {
				ids := []int{3, 6, 7, 8, 9}
				for _, gs := range gss {
					for _, g := range gs {
						So(ids, ShouldContain, g.ID)
					}
				}
			})
			Convey("Parents should only coming from the surving members", func() {
				ids := []int{3, 7, 8, 9} // 6 should be culled
				for _, gs := range gss {
					for _, g := range gs {
						So(ids, ShouldContain, g.ID)
					}
				}
			})
			Convey("The percent of single parents should equal the mutate only probability", func() {
				// Need to run many tests as this involves random numbers
				h.InterspeciesMateProbability = 0.0
				n := 0
				c := 0
				for i := 0; i < 1000; i++ {
					_, gss, _ := h.Select(p)
					n += len(gss)
					for _, gs := range gss {
						if len(gs) == 1 {
							c++
						}
					}
				}
				t.Log("c", c, "n", n)
				So(float64(c)/float64(n), ShouldAlmostEqual, h.MutateOnlyProbability, 0.01)
			})
			Convey("The percent of cross species parents should equal the interspeciaes probability", func() {
				// Need to run many tests as this involves random numbers
				h.MutateOnlyProbability = 0.0
				n := 0
				c := 0
				for i := 0; i < 1000; i++ {
					_, gss, _ := h.Select(p)
					n += len(gss)
					for _, gs := range gss {
						t.Log("len gs", len(gs))
						if gs[0].SpeciesID != gs[1].SpeciesID {
							c++
						}
					}
				}
				So(float64(c)/float64(n), ShouldAlmostEqual, h.InterspeciesMateProbability, 0.01)
			})
			Convey("The percent of offspring by species should equal the relative fitness of the remaining species", func() {

				// Need more genomes test cases as 9 is too few to get proper proportions
				for _, g := range p.Genomes {
					p.Genomes = append(p.Genomes, g)
				}

				// Need to run many tests as this involves random numbers
				n := 0
				c := make(map[int]int, 2)
				exp := map[int]float64{
					1: 11.0 / 23.5,
					3: 8.5 / 23.5,
					4: 4.0 / 23.5,
				}
				keep, gss, _ := h.Select(p)
				t.Log("keep", len(keep), "gss", len(gss))
				for i := 0; i < n; i++ {
					_, gss, _ := h.Select(p)
					n += len(gss)
					for _, gs := range gss {
						c[gs[0].SpeciesID]++
					}
				}
				t.Log("counts by species", c)
				for id, cnt := range c {
					So(float64(cnt)/float64(n), ShouldAlmostEqual, exp[id], 0.01)
				}

			})
			Convey("The total number of genomes kept and to be created should equal the original number", func() {
				So(len(keep)+len(gss), ShouldEqual, len(p.Genomes))
			})
		})

		Convey("When population is superstagnant (2*MaxStagnation)", func() {
			for i := 0; i < 2*h.MaxStagnation+1; i++ {
				h.Watch(p)
			}
			t.Log("overall stagnation", h.stagnation)
			keep, gss, err := h.Select(p)
			Convey("There should be no error", func() { So(err, ShouldBeNil) })
			Convey("There should be no kept genomes", func() { So(len(keep), ShouldEqual, 0) })
			Convey("The number of offspring should equal the original number of genomes", func() {
				So(len(gss), ShouldEqual, len(p.Genomes))
			})
			Convey("The parents should all be single parents", func() {
				found := false
				for _, gs := range gss {
					if len(gs) > 1 {
						found = true
						break
					}
				}
				So(found, ShouldBeFalse)
			})
		})

	})
}
