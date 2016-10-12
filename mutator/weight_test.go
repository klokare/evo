package mutator

import (
	"math"
	"testing"

	"github.com/klokare/evo"
	. "github.com/smartystreets/goconvey/convey"
)

func TestWeightMutate(t *testing.T) {
	trials := 10000
	Convey("Given a weight mutator and a genome", t, func() {

		h := &Weight{}
		g := &evo.Genome{
			Encoded: evo.Substrate{
				Conns: []evo.Conn{
					{Weight: 4.0},
				},
			},
		}

		Convey("The selector should implement the evo interface", func() {
			_, ok := interface{}(h).(evo.Mutator)
			So(ok, ShouldBeTrue)
		})

		Convey("When mutate weight probability is zero", func() {
			h.MutateWeightProbability = 0.0
			h.Mutate(g)
			Convey("There should be no change in the connection weight", func() {
				So(g.Encoded.Conns[0].Weight, ShouldEqual, 4.0)
			})
		})

		Convey("When mutate weight probability is 100% and replace weight probability is 0%", func() {
			h.MutateWeightProbability = 1.0
			h.ReplaceWeightProbability = 0.0
			v := make([]float64, trials)
			for i := 0; i < trials; i++ {
				g.Encoded.Conns[0].Weight = 4.0
				h.Mutate(g)
				v[i] = g.Encoded.Conns[0].Weight
			}

			mean := 0.0
			for _, x := range v {
				mean += x
			}
			mean /= float64(trials)

			stdev := 0.0
			for _, x := range v {
				stdev += (x - mean) * (x - mean)
			}
			stdev = math.Sqrt(stdev / float64(trials))

			Convey("The changes to the connection weight have a mean of the original weight and a standard deviation of 1.0", func() {
				So(mean, ShouldAlmostEqual, 4.0, 0.1)
				So(stdev, ShouldAlmostEqual, 1.0, 0.1)
			})
		})

		Convey("When mutate weight probability is 100% and replace weight probability is 100%", func() {
			h.MutateWeightProbability = 1.0
			h.ReplaceWeightProbability = 1.0
			v := make([]float64, trials)
			for i := 0; i < trials; i++ {
				g.Encoded.Conns[0].Weight = 4.0
				h.Mutate(g)
				v[i] = g.Encoded.Conns[0].Weight
			}

			mean := 0.0
			for _, x := range v {
				mean += x
			}
			mean /= float64(trials)

			stdev := 0.0
			for _, x := range v {
				stdev += (x - mean) * (x - mean)
			}
			stdev = math.Sqrt(stdev / float64(trials))

			Convey("The changes to the connection weight have a mean of 0.0 and a standard deviation of 1.0", func() {
				So(mean, ShouldAlmostEqual, 0.0, 0.1)
				So(stdev, ShouldAlmostEqual, 1.0, 0.1)
			})
		})

	})
}
