package searcher

import (
	"fmt"
	"testing"

	"github.com/klokare/evo"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSerialSearch(t *testing.T) {
	Convey("Given a serial searcher and some phenomes", t, func() {
		e := &mockEvaluator{}
		h := &Serial{Evaluator: e}
		ps := []evo.Phenome{
			{ID: 1}, {ID: 2}, {ID: 3},
		}
		Convey("When searching the phenomes", func() {
			Convey("When there is no error", func() {
				_, err := h.Search(ps)
				Convey("There should not be an error", func() { So(err, ShouldBeNil) })
				Convey("The evaluator should be called for each phenome", func() {
					for _, p := range ps {
						So(e.calledIDs, ShouldContain, p.ID)
					}
				})
			})
			Convey("When there is an error", func() {
				e.error = true
				_, err := h.Search(ps)
				Convey("There should be an error", func() { So(err, ShouldNotBeNil) })
			})
		})
	})
}

type mockEvaluator struct {
	calledIDs []int
	error     bool
}

func (h *mockEvaluator) Evaluate(p evo.Phenome) (r evo.Result, err error) {
	h.calledIDs = append(h.calledIDs, p.ID)
	r = evo.Result{ID: p.ID}
	if h.error {
		err = fmt.Errorf("error evaluating")
	}
	return
}
