package evo

import (
	"math/rand"
	"sort"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGenomeComplexity(t *testing.T) {
	Convey("Given a genome", t, func() {
		g := Genome{
			Encoded: Substrate{
				Nodes: make([]Node, 3),
				Conns: make([]Conn, 5),
			},
		}
		Convey("When calculating complexity", func() {
			c := g.Complexity()
			Convey("The value should be correct", func() { So(c, ShouldEqual, 8) })
		})
	})
}

func TestGenomesSort(t *testing.T) {
	Convey("Given a list of genomes", t, func() {
		a := []Genome{
			{ID: 1, Fitness: 1.0, Encoded: Substrate{Nodes: make([]Node, 1), Conns: make([]Conn, 1)}},
			{ID: 2, Fitness: 2.0, Encoded: Substrate{Nodes: make([]Node, 2), Conns: make([]Conn, 2)}},
			{ID: 4, Fitness: 2.0, Encoded: Substrate{Nodes: make([]Node, 2), Conns: make([]Conn, 1)}},
			{ID: 3, Fitness: 2.0, Encoded: Substrate{Nodes: make([]Node, 1), Conns: make([]Conn, 2)}},
			{ID: 5, Fitness: 5.0, Encoded: Substrate{Nodes: make([]Node, 2), Conns: make([]Conn, 1)}},
		}
		Convey("When sorting a radnomized copy", func() {
			var b Genomes = make([]Genome, 0, len(a))
			idxs := rand.Perm(len(a))
			for _, i := range idxs {
				b = append(b, a[i])
			}
			sort.Sort(b)
			Convey("The sorted list should be in order", func() {
				for i := 0; i < len(a); i++ {
					So(a[i].ID, ShouldEqual, b[i].ID)
				}
			})
		})
	})
}

func TestGenomesAvgFitness(t *testing.T) {
	Convey("Given a list of genomes", t, func() {
		var gs Genomes = []Genome{
			{Fitness: 2.0},
			{Fitness: 4.0},
			{Fitness: 5.0},
		}
		Convey("When calculating the average fitness", func() {
			x := gs.AvgFitness()
			Convey("The value should be correct", func() {
				So(x, ShouldAlmostEqual, 3.67, 0.01)
			})
		})
	})
	Convey("Given an empty list of genomes", t, func() {
		var gs Genomes
		Convey("When calculating the average fitness", func() {
			x := gs.AvgFitness()
			Convey("The value should be 0.0", func() { So(x, ShouldEqual, 0.0) })
		})
	})
}

func TestGenomesMaxFitness(t *testing.T) {
	Convey("Given a list of genomes", t, func() {
		var gs Genomes = []Genome{
			{Fitness: 2.0},
			{Fitness: 4.0},
			{Fitness: 5.0},
		}
		Convey("When calculating the average fitness", func() {
			x := gs.MaxFitness()
			Convey("The value should be correct", func() { So(x, ShouldEqual, 5.0) })
		})
	})
	Convey("Given an empty list of genomes", t, func() {
		var gs Genomes
		Convey("When calculating the average fitness", func() {
			x := gs.MaxFitness()
			Convey("The value should be 0.0", func() { So(x, ShouldEqual, 0.0) })
		})
	})
}
