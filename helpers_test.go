package evo

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func MutatorsMutate(t *testing.T) {
	Convey("Given a collection of mutators and a genome", t, func() {
		g := &Genome{
			Encoded: Substrate{
				Conns: make([]Conn, 2),
			},
		}
		m1 := &mockMutator{}
		m2 := &mockMutator{}
		m3 := &mockMutator{}
		var h Mutators = []Mutator{m1, m2, m3}

		Convey("The helper should implement the interface", func() {
			_, ok := interface{}(h).(Mutator)
			So(ok, ShouldBeTrue)
		})

		Convey("When mutating the genome and no structural changes occur", func() {
			h.Mutate(g)
			Convey("All mutators should be called", func() {
				So(m1.called, ShouldBeTrue)
				So(m2.called, ShouldBeTrue)
				So(m3.called, ShouldBeTrue)
			})
		})

		Convey("When mutating the genome and a structural change occurs", func() {
			m2.change = true
			h.Mutate(g)
			Convey("Mutators up to the change should be called", func() {
				So(m1.called, ShouldBeTrue)
				So(m2.called, ShouldBeTrue)
			})
			Convey("Mutators after the change should not be called", func() {
				So(m3.called, ShouldBeFalse)
			})
		})
	})
}

type mockMutator struct {
	called     bool
	change     bool
	changedIDs []int
}

func (h *mockMutator) Mutate(g *Genome) error {
	h.called = true
	if h.change {
		g.Encoded.Conns = append(g.Encoded.Conns, Conn{})
		h.changedIDs = append(h.changedIDs, g.ID)
	}
	return nil
}

func TestWatchersWatch(t *testing.T) {
	Convey("Given a set of watchers and a population", t, func() {
		w1 := &mockWatcher{}
		w2 := &mockWatcher{}
		w3 := &mockWatcher{}
		var h Watchers = []Watcher{w1, w2, w3}
		p := Population{}

		Convey("The helper should implement the interface", func() {
			_, ok := interface{}(h).(Watcher)
			So(ok, ShouldBeTrue)
		})

		Convey("When informing each watcher", func() {
			h.Watch(p)
			Convey("Each watcher should be called", func() {
				So(w1.called, ShouldBeTrue)
				So(w2.called, ShouldBeTrue)
				So(w3.called, ShouldBeTrue)
			})
		})
	})
}

type mockWatcher struct {
	called bool
}

func (h *mockWatcher) Watch(Population) error {
	h.called = true
	return nil
}
