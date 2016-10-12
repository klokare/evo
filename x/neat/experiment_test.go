package evo

import (
	"testing"

	"github.com/klokare/evo"
	"github.com/klokare/evo/comparer"
	"github.com/klokare/evo/crosser"
	"github.com/klokare/evo/mutator"
	"github.com/klokare/evo/searcher"
	"github.com/klokare/evo/speciator"
	"github.com/klokare/evo/transcriber"
	"github.com/klokare/evo/translator"
	"github.com/klokare/evo/x/configurer"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewExperiment(t *testing.T) {
	Convey("Given a configuration and an evaluator", t, func() {
		c := &configurer.JSON{
			Source: []byte(`{"DifferentConnsCoefficient":1.0,"DifferentNodesCoefficient":0.0,"WeightCoefficient":0.3,"EnableProbability":0.2,"AddNodeProbability":0.03,"AddConnProbability":0.05,"AllowRecurrent":false,"MutateWeightProbability":0.9,"ReplaceWeightProbability":0.1,"SurvivalRate":0.2,"MaxStagnation":15,"MutateOnlyProbability":0.25,"InterspeciesMateProbability":0.001,"CompatibilityThreshold":3.0}`),
		}
		e := &mockEvaluator{}

		f := func(e *evo.Experiment) error {
			e.Crosser.(*crosser.Multiple).EnableProbability *= -1.0
			return nil
		}

		Convey("When creating the new experiment", func() {
			exp, err := NewExperiment(c, e)
			Convey("There should be no error", func() { So(err, ShouldBeNil) })
			Convey("There should be an experiment", func() { So(exp, ShouldNotBeNil) })

			Convey("The crosser should be properly configured", func() {
				var ok bool
				var h *crosser.Multiple
				h, ok = exp.Crosser.(*crosser.Multiple)
				So(ok, ShouldBeTrue)
				So(h.EnableProbability, ShouldEqual, 0.2)
			})

			Convey("The mutator should be properly configured", func() {
				var ok bool
				var m evo.Mutators
				var w *mutator.Weight
				var c *mutator.Complexify
				m, ok = exp.Mutator.(evo.Mutators)
				So(ok, ShouldBeTrue)

				c, ok = m[0].(*mutator.Complexify)
				So(ok, ShouldBeTrue)
				So(c.AddConnProbability, ShouldEqual, 0.05)
				So(c.AddNodeProbability, ShouldEqual, 0.03)
				So(c.AllowRecurrent, ShouldEqual, false)

				w, ok = m[1].(*mutator.Weight)
				So(ok, ShouldBeTrue)
				So(w.MutateWeightProbability, ShouldEqual, 0.9)
				So(w.ReplaceWeightProbability, ShouldEqual, 0.1)
			})

			Convey("The searcher should be properly configured", func() {
				var ok bool
				var s *searcher.Serial
				s, ok = exp.Searcher.(*searcher.Serial)
				So(ok, ShouldBeTrue)

				_, ok = s.Evaluator.(*mockEvaluator)
				So(ok, ShouldEqual, true)
			})

			Convey("The transcriber should be properly configured", func() {
				var ok bool
				_, ok = exp.Transcriber.(*transcriber.NEAT)
				So(ok, ShouldBeTrue)
			})

			Convey("The translater should be properly configured", func() {
				var ok bool
				_, ok = exp.Translator.(*translator.Simple)
				So(ok, ShouldBeTrue)
			})

			Convey("The speciator should be properly configured", func() {
				var ok bool
				var h *speciator.Dynamic
				h, ok = exp.Speciater.(*speciator.Dynamic)
				So(ok, ShouldBeTrue)
				So(h.CompatibilityThreshold, ShouldEqual, 3.0)
			})

			Convey("The comparer should be properly configured", func() {
				var ok bool
				var s *speciator.Dynamic
				s, ok = exp.Speciater.(*speciator.Dynamic)
				So(ok, ShouldBeTrue)

				var h *comparer.Distance
				h, ok = s.Static.Comparer.(*comparer.Distance)
				So(ok, ShouldBeTrue)
				So(h.DifferentConnsCoefficient, ShouldEqual, 1.0)
				So(h.DifferentNodesCoefficient, ShouldEqual, 0.0)
				So(h.WeightCoefficient, ShouldEqual, 0.3)
			})

		})

		Convey("When creating the experiment with options", func() {
			exp, err := NewExperiment(c, e, f)
			Convey("There should be no error", func() { So(err, ShouldBeNil) })
			Convey("The crosser should reflect the optoin", func() {
				var ok bool
				var h *crosser.Multiple
				h, ok = exp.Crosser.(*crosser.Multiple)
				So(ok, ShouldBeTrue)
				So(h.EnableProbability, ShouldEqual, -0.2)
			})
		})
	})
}

type mockEvaluator struct{}

func (h *mockEvaluator) Evaluate(evo.Phenome) (evo.Result, error) {
	return evo.Result{}, nil
}
