package evo

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestExperimentRun(t *testing.T) {
	Convey("Given an expriement, a population and a target number of iterations", t, func() {
		h1 := &mockCrosser{}
		h2 := &mockMutator{change: true}
		h3 := &mockSearcher{}
		h4 := &mockSelector{keep: 0, parents: 10}
		h5 := &mockSpeciater{}
		h6 := &mockTranscriber{}
		h7 := &mockTranslator{}
		h8 := &genWatcher{}
		e := &Experiment{
			Crosser:     h1,
			Mutator:     h2,
			Searcher:    h3,
			Selector:    h4,
			Speciater:   h5,
			Transcriber: h6,
			Translator:  h7,
			Watcher:     h8,
		}
		p := Population{
			Generation: 1,
			Species: []Species{
				{ID: 1},
			},
			Genomes: []Genome{
				{ID: 1, SpeciesID: 1, Encoded: Substrate{Nodes: make([]Node, 2)}},
				{ID: 2, SpeciesID: 1, Encoded: Substrate{Nodes: make([]Node, 2)}},
				{ID: 3, SpeciesID: 1, Encoded: Substrate{Nodes: make([]Node, 2)}},
				{ID: 4, SpeciesID: 1, Encoded: Substrate{Nodes: make([]Node, 2)}},
				{ID: 5, SpeciesID: 1, Encoded: Substrate{Nodes: make([]Node, 2)}},
				{ID: 6, SpeciesID: 1, Encoded: Substrate{Nodes: make([]Node, 2)}},
				{ID: 7, SpeciesID: 1, Encoded: Substrate{Nodes: make([]Node, 2)}},
				{ID: 8, SpeciesID: 1, Encoded: Substrate{Nodes: make([]Node, 2)}},
				{ID: 9, SpeciesID: 1, Encoded: Substrate{Nodes: make([]Node, 2)}},
				{ID: 10, SpeciesID: 1, Encoded: Substrate{Nodes: make([]Node, 2)}},
			},
		}
		n := 3
		Convey("When running the experiment", func() {
			Convey("When there are no errors or solutions found", func() {
				err := Run(e, p, n)
				Convey("There should be no error", func() { So(err, ShouldBeNil) })
				Convey("Each helper should have been called", func() {
					So(h1.called, ShouldBeTrue)
					So(h2.called, ShouldBeTrue)
					So(h3.called, ShouldBeTrue)
					So(h4.called, ShouldBeTrue)
					So(h5.called, ShouldBeTrue)
					So(h6.called, ShouldBeTrue)
					So(h7.called, ShouldBeTrue)
					So(h8.called, ShouldBeTrue)
				})
				Convey("The generation should have advanced n-1 times", func() {
					So(h8.generation, ShouldEqual, 3)
				})
			})
			Convey("When there is a solution in the second iteration", func() {
				h3.solvedID = 11
				err := Run(e, p, n)
				Convey("There should be no error", func() { So(err, ShouldBeNil) })
				Convey("The generation should have advanced by 1", func() {
					So(h8.generation, ShouldEqual, 2)
				})
			})
			Convey("When there was an error advancing", func() {
				h1.error = true
				err := Run(e, p, n)
				Convey("There should be an error", func() { So(err, ShouldNotBeNil) })
				Convey("The generation should not have advanced n times", func() {
					So(h8.generation, ShouldBeLessThan, 3)
				})
			})

			Convey("When there was an error transcribing", func() {
				h6.error = true
				err := Run(e, p, n)
				Convey("There should be an error", func() { So(err, ShouldNotBeNil) })
				Convey("The generation should not have advanced n times", func() {
					So(h8.generation, ShouldBeLessThan, 3)
				})
			})
			Convey("When there was an error translating", func() {
				h7.error = true
				err := Run(e, p, n)
				Convey("There should be an error", func() { So(err, ShouldNotBeNil) })
				Convey("The generation should not have advanced n times", func() {
					So(h8.generation, ShouldBeLessThan, 3)
				})
			})
			Convey("When there was an error searching", func() {
				h3.error = true
				err := Run(e, p, n)
				Convey("There should be an error", func() { So(err, ShouldNotBeNil) })
				Convey("The generation should not have advanced n times", func() {
					So(h8.generation, ShouldBeLessThan, 3)
				})
			})
			Convey("When there was an error watching", func() {
				h8.error = true
				err := Run(e, p, n)
				Convey("There should be an error", func() { So(err, ShouldNotBeNil) })
				Convey("The generation should not have advanced n times", func() {
					So(h8.generation, ShouldBeLessThan, 3)
				})
			})

		})
	})
}
func TestExperimentInitIDs(t *testing.T) {
	Convey("Given an experiment and a population", t, func() {
		e := &Experiment{}
		p := Population{
			Genomes: []Genome{
				{ID: 5}, {ID: 4}, {ID: 7},
			},
			Species: []Species{
				{ID: 11}, {ID: 12},
			},
		}
		Convey("When initialising the ID sequences", func() {
			e.initIDs(p)
			Convey("The genome ID sequence should equal the max genome ID", func() {
				So(e.genomeID, ShouldEqual, 7)
			})
			Convey("The species ID sequence should equal the max species ID", func() {
				So(e.speciesID, ShouldEqual, 12)
			})
		})
	})
}

func TestExperimentAdvance(t *testing.T) {
	Convey("Given an experiment with a selector and a population", t, func() {
		s := &mockSelector{keep: 2, parents: 8}
		h := &mockSpeciater{}
		c := &mockCrosser{}
		m := &mockMutator{}
		e := &Experiment{Crosser: c, Mutator: m, Selector: s, Speciater: h}
		p := &Population{Generation: 2, Genomes: make([]Genome, 10)}
		Convey("When advancing the the experiment", func() {
			Convey("When there is no error", func() {
				err := e.advance(p)
				Convey("There should be no error", func() { So(err, ShouldBeNil) })
				Convey("Total genomes should be the same length as the population's genomes", func() {
					So(len(p.Genomes), ShouldEqual, 10)
				})
				Convey("The population's generation should advance", func() {
					So(p.Generation, ShouldEqual, 3)
				})
				Convey("The speciater should have been called", func() { So(h.called, ShouldBeTrue) })
			})
			Convey("When there are no parents", func() {
				s.parents = 0
				s.keep = 10
				e.advance(p)
				Convey("The population's generation should not advance", func() {
					So(p.Generation, ShouldEqual, 2)
				})
			})
			Convey("When there are not enough offspring", func() {
				s.parents = 5
				err := e.advance(p)
				Convey("There should be an error", func() { So(err, ShouldNotBeNil) })
			})
			Convey("When there is an error", func() {
				s.error = true
				err := e.advance(p)
				Convey("There should be an error", func() { So(err, ShouldNotBeNil) })
			})
		})
	})
}

func TestExperimentProcreate(t *testing.T) {
	Convey("Given an experiment witha a crosser and mutator and some genome couples", t, func() {
		c := &mockCrosser{}
		m := &mockMutator{change: true}
		e := &Experiment{
			Crosser:  c,
			Mutator:  m,
			genomeID: 5,
		}
		gss := make([][]Genome, 4)
		for i := 0; i < len(gss); i++ {
			gss[i] = make([]Genome, 1)
		}
		Convey("When creating offspring", func() {
			Convey("When there are no errors", func() {
				os, err := e.procreate(gss)
				Convey("There should be no error", func() { So(err, ShouldBeNil) })
				Convey("The number of offspring should equal the number of couples", func() {
					So(len(os), ShouldEqual, len(gss))
				})
				Convey("The offspring should have distinct, unique ids", func() {
					ids := []int{6, 7, 8, 9}
					for _, g := range os {
						So(ids, ShouldContain, g.ID)
					}
				})
				Convey("The mutator should have been called on the offspring", func() {
					for _, g := range os {
						So(m.changedIDs, ShouldContain, g.ID)
					}
				})
			})
			Convey("When there is an error", func() {
				c.error = true
				_, err := e.procreate(gss)
				Convey("There should be an error", func() { So(err, ShouldNotBeNil) })
			})
		})
	})
}

func TestExperimentTranscribe(t *testing.T) {
	Convey("Given an experiment with a transcriber and genomes", t, func() {
		h := &mockTranscriber{}
		e := &Experiment{Transcriber: h}
		gs := []Genome{
			{ID: 1, Encoded: Substrate{Nodes: make([]Node, 3)}},
			{ID: 2, Encoded: Substrate{Nodes: make([]Node, 5)}},
		}
		Convey("When transcribing the genomes", func() {
			Convey("When there is no error", func() {
				err := e.transcribe(gs)
				Convey("There should be no error", func() { So(err, ShouldBeNil) })
				Convey("The genome's decoded substrate should not be empty", func() {
					for _, g := range gs {
						So(len(g.Encoded.Nodes), ShouldBeGreaterThan, 0)
					}
				})
			})
			Convey("When there is an error transcribing", func() {
				h.error = true
				err := e.transcribe(gs)
				Convey("There should be an error", func() { So(err, ShouldNotBeNil) })
			})
		})
	})
}

func TestExperimentTranslate(t *testing.T) {
	Convey("Given an experiment with a translator and genomes", t, func() {
		h := &mockTranslator{}
		e := &Experiment{Translator: h}
		gs := []Genome{
			{ID: 1, Traits: []float64{0.1, 0.2}, Decoded: Substrate{Nodes: make([]Node, 1)}},
			{ID: 2, Traits: []float64{0.2, 0.3}, Decoded: Substrate{Nodes: make([]Node, 1)}},
		}

		Convey("When translating genomes into phenomes", func() {
			Convey("When there was no error in any translation", func() {
				ps, err := e.translate(gs)
				Convey("There should not be an error", func() { So(err, ShouldBeNil) })
				Convey("Each genome should have a phenome with the matching traits", func() {
					for _, g := range gs {
						found := false
						for _, p := range ps {
							if p.ID == g.ID {
								found = true
								So(p.Traits, ShouldResemble, g.Traits)
							}
						}
						So(found, ShouldBeTrue)
					}
				})
				Convey("Each phenome should have a network", func() {
					for _, p := range ps {
						So(p.Network, ShouldNotBeNil)
					}
				})
			})
			Convey("When there was an error during a translation", func() {
				h.error = true
				_, err := e.translate(gs)
				Convey("There should be an error", func() { So(err, ShouldNotBeNil) })
			})
			Convey("When a genome does not have a decoded substrate", func() {
				gs[0].Decoded = Substrate{}
				ps, err := e.translate(gs)
				Convey("There should be an error", func() { So(err, ShouldNotBeNil) })
				Convey("The phenome for that genome should have no network", func() {
					for _, p := range ps {
						if p.ID == gs[0].ID {
							So(p.Network, ShouldBeNil)
							break
						}
					}
				})
			})
		})
	})
}

func TestExperimentSearch(t *testing.T) {
	Convey("Given an experiment and phenomes", t, func() {
		s := &mockSearcher{}
		e := &Experiment{Searcher: s}
		ps := []Phenome{
			{ID: 1}, {ID: 2}, {ID: 3},
		}
		Convey("When searching the phenomes", func() {
			rs, err := e.Search(ps)
			Convey("When there is no error", func() {
				Convey("There should be no error", func() { So(err, ShouldBeNil) })
				Convey("There should be results for the phenomes", func() {
					for _, p := range ps {
						found := false
						for _, r := range rs {
							if r.ID == p.ID {
								found = true
								break
							}
						}
						So(found, ShouldBeTrue)
					}
				})
			})
			Convey("When there is an error", func() {
				s.errorID = 3
				_, err = e.Search(ps)
				Convey("An error should be returned", func() {
					So(err, ShouldNotBeNil)
				})
			})
		})
	})
}

func TestExperimentUpdate(t *testing.T) {
	Convey("Given genomes and results", t, func() {
		gs := []Genome{
			{ID: 1, Fitness: 1.0, Novelty: 2.0},
			{ID: 2, Fitness: 2.0, Novelty: 3.0},
		}
		rs := []Result{
			{ID: 1, Fitness: 3.0, Novelty: 4.0},
			{ID: 2, Fitness: 5.0, Novelty: 6.0},
		}
		Convey("When updating the genomes and there are no solutions or errors", func() {
			solved, err := update(gs, rs)
			Convey("Solved should be false", func() { So(solved, ShouldBeFalse) })
			Convey("Error should be nil", func() { So(err, ShouldBeNil) })
			Convey("The genomes' fitnesses and novelty should be updated", func() {
				So(gs[0].Fitness, ShouldEqual, 3.0)
				So(gs[0].Novelty, ShouldEqual, 4.0)
				So(gs[1].Fitness, ShouldEqual, 5.0)
				So(gs[1].Novelty, ShouldEqual, 6.0)
			})
		})
		Convey("When updating the genomes and there is a solution", func() {
			rs[0].Solved = true
			solved, _ := update(gs, rs)
			Convey("Solved should be true", func() { So(solved, ShouldBeTrue) })
		})
		Convey("When updating the genomes and there is an error", func() {
			rs[0].Error = fmt.Errorf("Error in evaluation")
			_, err := update(gs, rs)
			Convey("Error should not be nil", func() { So(err, ShouldNotBeNil) })
		})
	})
}

func TestExperimentStagnate(t *testing.T) {
	Convey("Given a population", t, func() {
		p := Population{
			Species: []Species{
				{ID: 1, Fitness: 4.0, Stagnation: 2},
				{ID: 2, Fitness: 3.0, Stagnation: 4},
			},
			Genomes: []Genome{
				{ID: 5, SpeciesID: 1, Fitness: 3.0},
				{ID: 6, SpeciesID: 2, Fitness: 5.0},
			},
		}
		Convey("When the species does not improve", func() {
			stagnate(&p)
			Convey("The species' stagnation should increment by 1", func() {
				So(p.Species[0].Stagnation, ShouldEqual, 3)
			})
			Convey("The species' fitness should not change", func() {
				So(p.Species[0].Fitness, ShouldEqual, 4.0)
			})
		})
		Convey("When calculating stagnation", func() {
			stagnate(&p)
			Convey("When the species improves", func() {
				Convey("The stagnation should be set to 0", func() {
					So(p.Species[1].Stagnation, ShouldEqual, 0)
				})
				Convey("The species should show the new fitness", func() {
					So(p.Species[1].Fitness, ShouldEqual, 5.0)
				})
			})
		})
	})
}

type mockCrosser struct {
	called bool
	error  bool
}

func (h *mockCrosser) Cross(ps ...Genome) (Genome, error) {
	var err error
	if h.error {
		err = fmt.Errorf("error crossing genomes")
	}
	h.called = true
	o := Genome{}
	if len(ps[0].Encoded.Nodes) > 0 {
		o.Encoded.Nodes = make([]Node, len(ps[0].Encoded.Nodes))
		copy(o.Encoded.Nodes, ps[0].Encoded.Nodes)
	}
	if len(ps[0].Encoded.Conns) > 0 {
		o.Encoded.Conns = make([]Conn, len(ps[0].Encoded.Conns))
		copy(o.Encoded.Conns, ps[0].Encoded.Conns)
	}

	return o, err
}

type mockSearcher struct {
	errorID  int
	solvedID int
	called   bool
	error    bool
}

func (m *mockSearcher) Search(ps []Phenome) (rs []Result, err error) {
	m.called = true
	rs = make([]Result, len(ps))
	for i, p := range ps {
		rs[i] = Result{ID: p.ID}
		if p.ID == m.errorID {
			err = fmt.Errorf("error during search")
		}
		if p.ID == m.solvedID {
			rs[i].Solved = true
		}
	}
	if m.error {
		err = fmt.Errorf("error during search")
	}
	return
}

type mockSpeciater struct {
	called bool
}

func (h *mockSpeciater) Speciate(*Population) error {
	h.called = true
	return nil
}

type mockSelector struct {
	called        bool
	error         bool
	keep, parents int
}

func (h *mockSelector) Select(Population) ([]Genome, [][]Genome, error) {
	var err error
	if h.error {
		err = fmt.Errorf("error selecting genomes")
	}
	h.called = true
	keep := make([]Genome, h.keep)
	parents := make([][]Genome, h.parents)
	for i := 0; i < len(parents); i++ {
		parents[i] = make([]Genome, 2)
	}
	return keep, parents, err
}

type mockTranscriber struct {
	called bool
	error  bool
}

func (h *mockTranscriber) Transcribe(s Substrate) (Substrate, error) {
	var err error
	if h.error {
		err = fmt.Errorf("error transcribing substrate")
	}
	h.called = true
	return s, err
}

type mockTranslator struct {
	called bool
	error  bool
}

func (h *mockTranslator) Translate(s Substrate) (Network, error) {
	var err error
	if h.error {
		err = fmt.Errorf("Error during translation")
	}
	h.called = true
	return &mockNetwork{}, err
}

type mockNetwork struct{}

func (h *mockNetwork) Activate(inputs []float64) []float64 { return inputs }

type genWatcher struct {
	generation int
	called     bool
	error      bool
}

func (h *genWatcher) Watch(p Population) error {
	h.generation = p.Generation
	h.called = true
	if h.error {
		return fmt.Errorf("error watching")
	}
	return nil
}
