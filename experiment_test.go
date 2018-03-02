package evo

import (
	"errors"
	"testing"

	"github.com/klokare/evo/internal/test"
)

func TestExperimentBatch(t *testing.T) {

	// Test with errors
	t.Run("with errors", func(t *testing.T) {

		// Execute without options and will get an error
		r := 3
		_, errs := Batch(r, 2)

		// There should be errors
		if len(errs) != r {
			t.Errorf("incorrect number of potential erros returned: expected %d, actual %d", r, len(errs))
		} else {
			for i, err := range errs {
				if err == nil {
					t.Errorf("error expected for run %d", i)
				}
			}
		}
	})

	// Test without errors
	t.Run("without errors", func(t *testing.T) {

		// Design the experiment
		var options = []Option{
			WithConfiguration(&MockConfigurer{}),
			WithCompare(ByFitness),
			WithMockCrosser(),
			WithMockEvaluator(),
			WithMockSeeder(),
			WithMockSearcher(),
			WithMockSelector(),
			WithMockSpeciator(),
			WithMockTranscriber(),
			WithMockTranslator(),
			WithMockMutators(),
		}

		// Execute the batch
		r := 3
		pops, errs := Batch(r, 2, options...)

		// There should be the correct number of populations
		if len(pops) != r {
			t.Errorf("incorrect number of populations: expected %d, actual %d", r, len(pops))
		} else {
			for i, p := range pops {
				if len(p.Genomes) == 0 {
					t.Errorf("run %d not executed", i)
				}
			}
		}

		// There should be no errors
		if len(errs) != r {
			t.Errorf("incorrect number of potential errors returned: expected %d, actual %d", r, len(errs))
		} else {
			for i, err := range errs {
				if err != nil {
					t.Errorf("unexpected error for run %d: %v", i, err)
				}
			}
		}
	})
}

// Things to test in Run
// • If solution is found, the experiment should end early returning after the whole population's results are processed
func TestExperimentRun(t *testing.T) {

	var options = []Option{
		WithConfiguration(&MockConfigurer{}),
		WithCompare(ByFitness),
		WithMockCrosser(),
		WithMockEvaluator(),
		WithMockSeeder(),
		WithMockSearcher(),
		WithMockSelector(),
		WithMockSpeciator(),
		WithMockTranscriber(),
		WithMockTranslator(),
		WithMockMutators(),
	}

	// There error from a failed option should be propogated up
	t.Run("failed option", func(t *testing.T) {
		defer func() { options = options[:len(options)-1] }() // pop off the added option
		options = append(options, func(e *Experiment) error {
			return errors.New("option error")
		})
		_, err := Run(1, options...)
		t.Run("error", test.Error(true, err))
	})

	// There should be an error if the required helpers are not present. We can test this by
	// iterating the options, skipping each one incrementally, and testing a Run.
	t.Run("failed verify", func(t *testing.T) {
		_, err := Run(1) // run with no options
		t.Run("error", test.Error(true, err))
	})

	// Failing the population seeding should produce an error
	t.Run("failed seeding", func(t *testing.T) {
		defer func() { options = options[:len(options)-1] }() // pop off the added option
		options = append(options, func(e *Experiment) error {
			e.Seeder = &MockSeeder{HasError: true}
			return nil
		})
		_, err := Run(1, options...)
		t.Run("error", test.Error(true, err))
	})

	// Failed speciation should produce an error
	t.Run("failed speciation", func(t *testing.T) {
		defer func() { options = options[:len(options)-1] }() // pop off the added option
		options = append(options, func(e *Experiment) error {
			e.Speciator = &MockSpeciator{HasError: true}
			return nil
		})
		_, err := Run(1, options...)
		t.Run("error", test.Error(true, err))
	})

	// Failed advancing should produce an error
	t.Run("failed advancing", func(t *testing.T) {
		defer func() { options = options[:len(options)-1] }() // pop off the added option
		options = append(options, func(e *Experiment) error {
			e.Selector = &MockSelector{HasError: true}
			return nil
		})
		_, err := Run(1, options...)
		t.Run("error", test.Error(true, err))
	})

	// Failed evaluation should produce an error
	t.Run("failed evaluation", func(t *testing.T) {
		defer func() { options = options[:len(options)-1] }() // pop off the added option
		options = append(options, func(e *Experiment) error {
			e.Searcher = &MockSearcher{HasError: true}
			return nil
		})
		_, err := Run(1, options...)
		t.Run("error", test.Error(true, err))
	})

	// Failed publishing on advance should produce an error
	t.Run("failed publishing on advance", func(t *testing.T) {
		defer func() { options = options[:len(options)-1] }() // pop off the added option
		options = append(options, func(e *Experiment) error {
			e.Subscribe(Advanced, func(bool, Population) error { return errors.New("callback error") })
			return nil
		})
		_, err := Run(1, options...)
		t.Run("error", test.Error(true, err))
	})

	// Failed publishing on evaluate should produce an error
	t.Run("failed publishing on advance", func(t *testing.T) {
		defer func() { options = options[:len(options)-1] }() // pop off the added option
		options = append(options, func(e *Experiment) error {
			e.Subscribe(Evaluated, func(bool, Population) error { return errors.New("callback error") })
			return nil
		})
		_, err := Run(1, options...)
		t.Run("error", test.Error(true, err))
	})
}

// Thing to test:
// • no listeners should be ok
// • if listener throws an error then an error should be return
// • no errors in listeners then no error should be returned
// • final and pop should be passed on
// • The right listener is called
func TestSubscribeAndPublish(t *testing.T) {

	// Begin with an empty set of listeners and a population
	pop := Population{Generation: 1}
	e := &Experiment{listeners: make(map[Event][]Listener, 10)}

	// Calling empty set causes no error
	err := e.publish(Evaluated, false, pop)
	if err != nil {
		t.Errorf("unexpected error with empty set: %v", err)
	}

	// Add a new listner for evaluated
	var f1 bool
	var p1 Population
	e.Subscribe(Evaluated, func(final bool, pop Population) error {
		f1, p1 = final, pop
		return nil
	})

	err = e.publish(Evaluated, false, pop)
	if err != nil {
		t.Errorf("unexpected error with first set: %v", err)
	}
	if f1 != false {
		t.Errorf("incorrect final with first set: expected %t, actual %t", false, f1)
	}
	if p1.Generation != pop.Generation {
		t.Errorf("incorrect population with first set: expected %d, actual %d", pop.Generation, p1.Generation)
	}

	// Add a second listener for evaluated and check both results
	var f2 bool
	var p2 Population
	e.Subscribe(Evaluated, func(final bool, pop Population) error {
		f2, p2 = final, pop
		return nil
	})

	err = e.publish(Evaluated, true, Population{Generation: 3})
	if err != nil {
		t.Errorf("unexpected error with second set: %v", err)
	}
	if f1 != true {
		t.Errorf("incorrect final with second set, first listener: expected %t, actual %t", true, f1)
	}
	if p1.Generation != 3 {
		t.Errorf("incorrect population with second set, first listener: expected %d, actual %d", 3, p1.Generation)
	}
	if f2 != true {
		t.Errorf("incorrect final with second set, second listener: expected %t, actual %t", true, f2)
	}
	if p2.Generation != 3 {
		t.Errorf("incorrect population with second set, second listener: expected %d, actual %d", 3, p1.Generation)
	}

	// Add listener for advanced
	f2 = false // reset so we can check that it was not called
	var f3 bool
	var p3 Population
	e.Subscribe(Advanced, func(final bool, pop Population) error {
		f3, p3 = final, pop
		return nil
	})

	err = e.publish(Advanced, true, pop)
	if err != nil {
		t.Errorf("unexpected error with third set: %v", err)
	}
	if f3 != true {
		t.Errorf("incorrect final with third set: expected %t, actual %t", true, f3)
	}
	if p3.Generation != pop.Generation {
		t.Errorf("incorrect population with third set: expected %d, actual %d", pop.Generation, p3.Generation)
	}
	if f2 != false {
		t.Errorf("incorrect final with third set, dormant listner: expected %t, actual %t", false, f2)
	}

	// A listener that throws an error
	e.Subscribe(Advanced, func(final bool, pop Population) error {
		return errors.New("listener error")
	})
	err = e.publish(Advanced, true, pop)
	if err == nil {
		t.Errorf("expected error with fourth set but none received")
	}
}

// Things to test:
// • If required helpers are missing then an error should be returned
func TestVerify(t *testing.T) {
	var cases = []struct {
		Desc     string
		HasError bool
		Crosser
		Evaluator
		Seeder
		Searcher
		Selector
		Speciator
		Translator
		Transcriber
		Mutators []Mutator
		Compare
	}{
		{
			Desc:        "nothing missing",
			HasError:    false,
			Crosser:     &MockCrosser{},
			Evaluator:   &MockEvaluator{},
			Seeder:      &MockSeeder{},
			Searcher:    &MockSearcher{},
			Selector:    &MockSelector{},
			Speciator:   &MockSpeciator{},
			Translator:  &MockTranslator{},
			Transcriber: &MockTranscriber{},
			Mutators:    []Mutator{&MockStructureMutator{}},
			Compare:     ByFitness,
		},
		{
			Desc:        "crosser missing",
			HasError:    true,
			Crosser:     nil,
			Evaluator:   &MockEvaluator{},
			Seeder:      &MockSeeder{},
			Searcher:    &MockSearcher{},
			Selector:    &MockSelector{},
			Speciator:   &MockSpeciator{},
			Translator:  &MockTranslator{},
			Transcriber: &MockTranscriber{},
			Mutators:    []Mutator{&MockStructureMutator{}},
			Compare:     ByFitness,
		},
		{
			Desc:        "evaluator missing",
			HasError:    true,
			Crosser:     &MockCrosser{},
			Evaluator:   nil,
			Seeder:      &MockSeeder{},
			Searcher:    &MockSearcher{},
			Selector:    &MockSelector{},
			Speciator:   &MockSpeciator{},
			Translator:  &MockTranslator{},
			Transcriber: &MockTranscriber{},
			Mutators:    []Mutator{&MockStructureMutator{}},
			Compare:     ByFitness,
		},
		{
			Desc:        "seeder missing",
			HasError:    true,
			Crosser:     &MockCrosser{},
			Evaluator:   &MockEvaluator{},
			Seeder:      nil,
			Searcher:    &MockSearcher{},
			Selector:    &MockSelector{},
			Speciator:   &MockSpeciator{},
			Translator:  &MockTranslator{},
			Transcriber: &MockTranscriber{},
			Mutators:    []Mutator{&MockStructureMutator{}},
			Compare:     ByFitness,
		},
		{
			Desc:        "searcher missing",
			HasError:    true,
			Crosser:     &MockCrosser{},
			Evaluator:   &MockEvaluator{},
			Seeder:      &MockSeeder{},
			Searcher:    nil,
			Selector:    &MockSelector{},
			Speciator:   &MockSpeciator{},
			Translator:  &MockTranslator{},
			Transcriber: &MockTranscriber{},
			Mutators:    []Mutator{&MockStructureMutator{}},
			Compare:     ByFitness,
		},
		{
			Desc:        "selector missing",
			HasError:    true,
			Crosser:     &MockCrosser{},
			Evaluator:   &MockEvaluator{},
			Seeder:      &MockSeeder{},
			Searcher:    &MockSearcher{},
			Selector:    nil,
			Speciator:   &MockSpeciator{},
			Translator:  &MockTranslator{},
			Transcriber: &MockTranscriber{},
			Mutators:    []Mutator{&MockStructureMutator{}},
			Compare:     ByFitness,
		},
		{
			Desc:        "speciator missing",
			HasError:    true,
			Crosser:     &MockCrosser{},
			Evaluator:   &MockEvaluator{},
			Seeder:      &MockSeeder{},
			Searcher:    &MockSearcher{},
			Selector:    &MockSelector{},
			Speciator:   nil,
			Translator:  &MockTranslator{},
			Transcriber: &MockTranscriber{},
			Mutators:    []Mutator{&MockStructureMutator{}},
			Compare:     ByFitness,
		},
		{
			Desc:        "translator missing",
			HasError:    true,
			Crosser:     &MockCrosser{},
			Evaluator:   &MockEvaluator{},
			Seeder:      &MockSeeder{},
			Searcher:    &MockSearcher{},
			Selector:    &MockSelector{},
			Speciator:   &MockSpeciator{},
			Translator:  nil,
			Transcriber: &MockTranscriber{},
			Mutators:    []Mutator{&MockStructureMutator{}},
			Compare:     ByFitness,
		},
		{
			Desc:        "transcriber missing",
			HasError:    true,
			Crosser:     &MockCrosser{},
			Evaluator:   &MockEvaluator{},
			Seeder:      &MockSeeder{},
			Searcher:    &MockSearcher{},
			Selector:    &MockSelector{},
			Speciator:   &MockSpeciator{},
			Translator:  &MockTranslator{},
			Transcriber: nil,
			Mutators:    []Mutator{&MockStructureMutator{}},
			Compare:     ByFitness,
		},
		{
			Desc:        "mutators missing",
			HasError:    true,
			Crosser:     &MockCrosser{},
			Evaluator:   &MockEvaluator{},
			Seeder:      &MockSeeder{},
			Searcher:    &MockSearcher{},
			Selector:    &MockSelector{},
			Speciator:   &MockSpeciator{},
			Translator:  &MockTranslator{},
			Transcriber: &MockTranscriber{},
			Mutators:    nil,
			Compare:     ByFitness,
		},
		{
			Desc:        "mutators empty",
			HasError:    true,
			Crosser:     &MockCrosser{},
			Evaluator:   &MockEvaluator{},
			Seeder:      &MockSeeder{},
			Searcher:    &MockSearcher{},
			Selector:    &MockSelector{},
			Speciator:   &MockSpeciator{},
			Translator:  &MockTranslator{},
			Transcriber: &MockTranscriber{},
			Mutators:    []Mutator{},
			Compare:     ByFitness,
		},
		{
			Desc:        "compare missing",
			HasError:    true,
			Crosser:     &MockCrosser{},
			Evaluator:   &MockEvaluator{},
			Seeder:      &MockSeeder{},
			Searcher:    &MockSearcher{},
			Selector:    &MockSelector{},
			Speciator:   &MockSpeciator{},
			Translator:  &MockTranslator{},
			Transcriber: &MockTranscriber{},
			Mutators:    []Mutator{&MockStructureMutator{}},
			Compare:     nil,
		},
	}

	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {
			e := &Experiment{
				Crosser:     c.Crosser,
				Evaluator:   c.Evaluator,
				Seeder:      c.Seeder,
				Searcher:    c.Searcher,
				Selector:    c.Selector,
				Speciator:   c.Speciator,
				Translator:  c.Translator,
				Transcriber: c.Transcriber,
				Mutators:    c.Mutators,
				Compare:     c.Compare,
			}
			err := verify(e)
			t.Run("error", test.Error(c.HasError, err))
		})
	}
}

// Things to test:
// • The experiment's last genome ID should be the highest ID found in the genome list
func TestSetSequence(t *testing.T) {
	var genomes = []Genome{
		{ID: 1}, {ID: 10}, {ID: 5},
	}
	e := new(Experiment)
	setSequence(e, genomes)
	if *e.lastGID != 10 {
		t.Errorf("incorrect last genome id: expected 10, actual: %d", *e.lastGID)
	}
}

// Things to test:
// • An error creating phenomes should return an error
// • An error searching phenomes should return an error
// • A soltuion found in evaluation should return solved
func TestEvaluate(t *testing.T) {

	var cases = []struct {
		Desc string
		Searcher
		Evaluator
		Transcriber
		Translator
		HasError bool
		Solved   bool
	}{
		{
			Desc:        "create phenomes fails",
			Transcriber: &MockTranscriber{HasError: true},
			Translator:  &MockTranslator{},
			Searcher:    &MockSearcher{},
			Evaluator:   &MockEvaluator{},
			HasError:    true,
		},
		{
			Desc:        "search fails",
			Transcriber: &MockTranscriber{},
			Translator:  &MockTranslator{},
			Searcher:    &MockSearcher{HasError: true},
			Evaluator:   &MockEvaluator{},
			HasError:    true,
		},
		{
			Desc:        "solution expected",
			Transcriber: &MockTranscriber{},
			Translator:  &MockTranslator{},
			Searcher:    &MockSearcher{},
			Evaluator:   &MockEvaluator{HasSolve: true},
			HasError:    false,
			Solved:      true,
		},
		{
			Desc:        "solution not expected",
			Transcriber: &MockTranscriber{},
			Translator:  &MockTranslator{},
			Searcher:    &MockSearcher{},
			Evaluator:   &MockEvaluator{HasSolve: false},
			HasError:    false,
			Solved:      false,
		},
	}

	pop := Population{Species: []Species{{ID: 10}}, Genomes: []Genome{{ID: 1}}}
	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {

			solved, err := evaluate(c.Searcher, c.Evaluator, c.Transcriber, c.Translator, ByFitness, 0.0, &pop)

			if !t.Run("error", test.Error(c.HasError, err)) || c.HasError {
				return
			}

			if c.Solved && !solved {
				t.Errorf("expected solution but none received")
			} else if !c.Solved && solved {
				t.Errorf("received unexpected solution")
			}
		})
	}
}

// Things to test
// • a phenome should be returned for each genome
// • if a genome has traits, the phenome should have the same traits
// • if the phenome's network is nil, there should be an error
// • if the transcriber fails, there should be an error
// • if the translator fails, there should be an error
func TestCreatePhenomes(t *testing.T) {

	var cases = []struct {
		Desc string
		Transcriber
		Translator
		HasError bool
		Genomes  []Genome
		Phenomes []Phenome
	}{
		{
			Desc:        "transcriber failed",
			HasError:    true,
			Transcriber: &MockTranscriber{HasError: true},
			Translator:  &MockTranslator{},
			Genomes: []Genome{
				{ID: 1},
				{ID: 2, Traits: []float64{1, 2}},
			},
			Phenomes: []Phenome{
				{ID: 1},
				{ID: 2, Traits: []float64{1, 2}},
			},
		},
		{
			Desc:        "translator failed",
			HasError:    true,
			Transcriber: &MockTranscriber{},
			Translator:  &MockTranslator{HasError: true},
			Genomes: []Genome{
				{ID: 1},
				{ID: 2, Traits: []float64{1, 2}},
			},
			Phenomes: []Phenome{
				{ID: 1},
				{ID: 2, Traits: []float64{1, 2}},
			},
		},
		{
			Desc:        "network missing",
			HasError:    true,
			Transcriber: &MockTranscriber{},
			Translator:  &MockTranslator{HasMissing: true},
			Genomes: []Genome{
				{ID: 1},
				{ID: 2, Traits: []float64{1, 2}},
			},
			Phenomes: []Phenome{
				{ID: 1},
				{ID: 2, Traits: []float64{1, 2}},
			},
		},
		{
			Desc:        "successful",
			HasError:    false,
			Transcriber: &MockTranscriber{},
			Translator:  &MockTranslator{},
			Genomes: []Genome{
				{ID: 1},
				{ID: 2, Traits: []float64{1, 2}},
			},
			Phenomes: []Phenome{
				{ID: 1},
				{ID: 2, Traits: []float64{1, 2}},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {

			phenomes, err := createPhenomes(c.Transcriber, c.Translator, c.Genomes)

			if !t.Run("error", test.Error(c.HasError, err)) || c.HasError {
				return
			}

			if len(c.Phenomes) != len(phenomes) {
				t.Errorf("incorrect number of phenomes: expected %d, actual %d", len(c.Phenomes), len(phenomes))
			} else {
				for _, ep := range c.Phenomes {
					found := false
					for _, ap := range phenomes {
						if ep.ID == ap.ID {

							if len(ep.Traits) != len(ap.Traits) {
								t.Errorf("incorrect number of traits for phenome %d: expected %d, actual %d", ep.ID, len(ep.Traits), len(ap.Traits))
							} else {
								for i := 0; i < len(ep.Traits); i++ {
									if ep.Traits[i] != ap.Traits[i] {
										t.Errorf("incorrect trait %d for phenome %d: expected %f, actual %f", i, ep.ID, ep.Traits[i], ap.Traits[i])
									}
								}
							}

							found = true
							break
						}
					}
					if !found {
						t.Errorf("phenome not found for genome %d", ep.ID)
					}
				}
			}
		})
	}
}

// Things to test
// • if any result is solved then solved should be true
// • result for genome not in list is ok
// • genome not in result list is ok, too
// • all parts of the result should be copied to the genome
func TestUpdateGenomes(t *testing.T) {

	var cases = []struct {
		Desc          string
		Solved        bool
		Results       []Result
		Before, After []Genome
	}{
		{
			Desc:   "unsolved",
			Solved: true,
			Results: []Result{
				{ID: 1, Fitness: 1, Novelty: 1, Solved: false},
				{ID: 2, Fitness: 2, Novelty: 2, Solved: true},
				{ID: 3, Fitness: 3, Novelty: 3, Solved: false},
				{ID: 4, Fitness: 4, Novelty: 4, Solved: false},
			},
			Before: []Genome{{ID: 5}, {ID: 2}, {ID: 1}, {ID: 4}},
			After: []Genome{
				{ID: 5},
				{ID: 2, Fitness: 2, Novelty: 2, Solved: true},
				{ID: 1, Fitness: 1, Novelty: 1, Solved: false},
				{ID: 4, Fitness: 4, Novelty: 4, Solved: false},
			},
		},
		{
			Desc:   "solved",
			Solved: true,
			Results: []Result{
				{ID: 1, Fitness: 1, Novelty: 1, Solved: false},
				{ID: 2, Fitness: 2, Novelty: 2, Solved: true},
				{ID: 3, Fitness: 3, Novelty: 3, Solved: false},
				{ID: 4, Fitness: 4, Novelty: 4, Solved: false},
			},
			Before: []Genome{{ID: 5}, {ID: 2}, {ID: 1}, {ID: 4}},
			After: []Genome{
				{ID: 5},
				{ID: 2, Fitness: 2, Novelty: 2, Solved: true},
				{ID: 1, Fitness: 1, Novelty: 1, Solved: false},
				{ID: 4, Fitness: 4, Novelty: 4, Solved: false},
			},
		},
		{
			Desc:   "solution not in existing genome",
			Solved: false,
			Results: []Result{
				{ID: 1, Fitness: 1, Novelty: 1, Solved: false},
				{ID: 2, Fitness: 2, Novelty: 2, Solved: false},
				{ID: 3, Fitness: 3, Novelty: 3, Solved: true},
				{ID: 4, Fitness: 4, Novelty: 4, Solved: false},
			},
			Before: []Genome{{ID: 5}, {ID: 2}, {ID: 1}, {ID: 4}},
			After: []Genome{
				{ID: 5},
				{ID: 2, Fitness: 2, Novelty: 2, Solved: false},
				{ID: 1, Fitness: 1, Novelty: 1, Solved: false},
				{ID: 4, Fitness: 4, Novelty: 4, Solved: false},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {

			exp := c.After
			act := c.Before // would be updated
			solved := updateGenomes(act, c.Results)

			if c.Solved && !solved {
				t.Errorf("expected solution but none received")
			} else if !c.Solved && solved {
				t.Errorf("received unexpected solution")
			}

			if len(exp) != len(act) {
				t.Errorf("incorrect number of genomes: expected %d, actual %d", len(exp), len(act))
			} else {
				for _, eg := range exp {
					found := false
					for _, ag := range act {
						if eg.ID == ag.ID {
							if eg.Fitness != ag.Fitness {
								t.Errorf("incorrect fitness for genome %d: expected %f, actual %f", eg.ID, eg.Fitness, ag.Fitness)
							}
							if eg.Novelty != ag.Novelty {
								t.Errorf("incorrect novelty for genome %d: expected %f, actual %f", eg.ID, eg.Novelty, ag.Novelty)
							}
							if eg.Solved != ag.Solved {
								t.Errorf("incorrect solved state for genome %d: expected %t, actual %t", eg.ID, eg.Solved, ag.Solved)
							}
							found = true
							break
						}
					}
					if !found {
						t.Errorf("genome %d not found", eg.ID)
					}
				}
			}
		})
	}
}

// Things to test
// • If species has new champion, champion's ID should be stored and decay set to 0.0
// • If champion is same, champion ID should not change and decay should increase
// • Decay should not increase above 1.0
func TestUpdateSpecies(t *testing.T) {

	genomes := []Genome{
		{ID: 5, SpeciesID: 30, Fitness: 3.0},
		{ID: 1, SpeciesID: 10, Fitness: 3.0},
		{ID: 3, SpeciesID: 20, Fitness: 3.0},
		{ID: 2, SpeciesID: 10, Fitness: 5.0}, // New champion
		{ID: 4, SpeciesID: 30, Fitness: 3.0},
		{ID: 6, SpeciesID: 30, Fitness: 3.0},
	}
	before := []Species{
		{ID: 10, Champion: 1, Decay: 0.8},
		{ID: 20, Champion: 3, Decay: 0.2},
		{ID: 30, Champion: 4, Decay: 0.9},
	}
	after := []Species{
		{ID: 10, Champion: 2, Decay: 0.0},
		{ID: 20, Champion: 3, Decay: 0.4},
		{ID: 30, Champion: 4, Decay: 1.0},
	}

	exp := after
	act := before
	updateSpecies(ByFitness, 0.2, act, genomes)

	if len(exp) != len(act) {
		t.Errorf("incorrect number of species: expected %d, actual %d", len(exp), len(act))
	} else {
		for _, es := range exp {
			found := false
			for _, as := range act {
				if es.ID == as.ID {
					if es.Champion != as.Champion {
						t.Errorf("incorrect champion for species %d: expected %d, actual %d", es.ID, es.Champion, as.Champion)
					}
					if es.Decay != as.Decay {
						t.Errorf("incorrect decay for species %d: expected %f, actual %f", es.ID, es.Decay, as.Decay)
					}
					found = true
					break
				}
			}
			if !found {
				t.Errorf("species %d not found", es.ID)
			}
		}
	}
}

// Things to test:
// • if select fails then an error should be returned
// • If parents are given then the generation should increase
// • If making a child fails then an error should be given
// • There should be a child for every set of parents
// • The number of genomes should be sum of continuing and offspring
// • The speciator should be called
// • If the speciator fails then an error should be returned
// • The genome ID sequence should increment by at least the number of children
func TestAdvance(t *testing.T) {

	var cases = []struct {
		Desc     string
		HasError bool
		LastID   int64
		Selector
		Crosser
		Speciator
	}{
		{
			Desc:     "select failed",
			HasError: true,
			LastID:   100,
			Selector: &MockSelector{
				NumContinuing: 2,
				NumParents:    4,
				HasError:      true,
			},
			Crosser:   &MockCrosser{},
			Speciator: &MockSpeciator{},
		},
		{
			Desc:     "make child failed",
			HasError: true,
			LastID:   100,
			Selector: &MockSelector{
				NumContinuing: 2,
				NumParents:    4,
			},
			Crosser:   &MockCrosser{HasError: true},
			Speciator: &MockSpeciator{},
		},
		{
			Desc:     "speciate failed",
			HasError: true,
			LastID:   100,
			Selector: &MockSelector{
				NumContinuing: 2,
				NumParents:    4,
			},
			Crosser:   &MockCrosser{},
			Speciator: &MockSpeciator{HasError: true},
		},
		{
			Desc:     "new generation",
			HasError: false,
			LastID:   100,
			Selector: &MockSelector{
				NumContinuing: 2,
				NumParents:    4,
			},
			Crosser:   &MockCrosser{},
			Speciator: &MockSpeciator{},
		},
		{
			Desc:     "no new genration",
			HasError: false,
			LastID:   100,
			Selector: &MockSelector{
				NumContinuing: 2,
				NumParents:    0,
			},
			Crosser:   &MockCrosser{},
			Speciator: &MockSpeciator{},
		},
	}

	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {

			pop := &Population{Generation: 10}
			old := c.LastID
			err := advance(c.Selector, c.Crosser, nil, c.Speciator, &c.LastID, pop)

			if !t.Run("error", test.Error(c.HasError, err)) || c.HasError {
				return
			}

			nc := c.Selector.(*MockSelector).NumContinuing
			np := c.Selector.(*MockSelector).NumParents

			if c.LastID-old < int64(np) {
				t.Errorf("incorrect last id: expected > %d, actual %d", old+int64(np), c.LastID)
			}

			if nc+np != len(pop.Genomes) {
				t.Errorf("incorrect number of genomes: expected %d, actual %d", nc+np, len(pop.Genomes))
			}

			if np == 0 && pop.Generation != 10 {
				t.Errorf("generation increment even though no new offspring")
			} else if np > 0 && pop.Generation == 10 {
				t.Errorf("generation did not increment even though there were new offspring")
			}

			if c.Speciator.(*MockSpeciator).Called == 0 {
				t.Errorf("speciator should be invoked")
			}
		})
	}
}

// Things to test
// • if crosser fails then an error should be returned
// • if any mutator fails then an error should be returned
// • the child's ID should be set
// • the first mutator should be called
// • if the first mutator does not modify complexity then the second mutator should be called
func TestMakeChild(t *testing.T) {
	var cases = []struct {
		Desc     string
		HasError bool
		ID       int64
		Parents  []Genome
		Child    Genome
		Crosser
		Mutators []Mutator
	}{
		{
			Desc:     "crosser failed",
			HasError: true,
			ID:       100,
			Parents:  []Genome{{ID: 1}, {ID: 2}},
			Child:    Genome{ID: 100},
			Crosser:  &MockCrosser{HasError: true},
			Mutators: []Mutator{
				&MockStructureMutator{},
				&MockWeightMutator{},
			},
		},
		{
			Desc:     "mutator 1 failed",
			HasError: true,
			ID:       100,
			Parents:  []Genome{{ID: 1}, {ID: 2}},
			Child:    Genome{ID: 100},
			Crosser:  &MockCrosser{},
			Mutators: []Mutator{
				&MockStructureMutator{HasError: true},
				&MockWeightMutator{},
			},
		},
		{
			Desc:     "mutator 2 failed",
			HasError: true,
			ID:       100,
			Parents:  []Genome{{ID: 1}, {ID: 2}},
			Child:    Genome{ID: 100},
			Crosser:  &MockCrosser{HasError: true},
			Mutators: []Mutator{
				&MockStructureMutator{},
				&MockWeightMutator{HasError: true},
			},
		},
		{
			Desc:     "mutator 1 changes complexity",
			HasError: false,
			ID:       100,
			Parents:  []Genome{{ID: 1}, {ID: 2}},
			Child:    Genome{ID: 100},
			Crosser:  &MockCrosser{},
			Mutators: []Mutator{
				&MockStructureMutator{},
				&MockWeightMutator{},
			},
		},
		{
			Desc:     "mutator 1 does not change complexity",
			HasError: false,
			ID:       100,
			Parents:  []Genome{{ID: 1}, {ID: 2}},
			Child: Genome{
				ID: 100,
				Encoded: Substrate{
					Conns: []Conn{
						{
							Source:  Position{Layer: 0.25, X: 0.25},
							Target:  Position{Layer: 0.75, X: 0.75},
							Weight:  1.23,
							Enabled: true,
						},
					},
				},
			},
			Crosser: &MockCrosser{},
			Mutators: []Mutator{
				&MockStructureMutator{Execute: true},
				&MockWeightMutator{},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {

			exp := c.Child
			act, err := makeChild(c.Crosser, c.Mutators, c.ID, c.Parents)

			if !t.Run("error", test.Error(c.HasError, err)) || c.HasError {
				return
			}

			if c.Mutators[0].(*MockStructureMutator).Called == 0 {
				t.Errorf("mutator 1 should have been called")
			}
			if !c.Mutators[0].(*MockStructureMutator).Execute && c.Mutators[1].(*MockWeightMutator).Called == 0 {
				t.Errorf("mutator 2 should have been called")
			}

			if exp.ID != act.ID {
				t.Errorf("incorrect child ID: expected %d, actual %d", exp.ID, act.ID)
			}

			if exp.Complexity() != act.Complexity() {
				t.Errorf("incorrect child complexity: expected %d, actual %d", exp.Complexity(), act.Complexity())
			}
		})
	}
}

// Things to test:
// • Returns an option
func TestWithConfiguration(t *testing.T) {
	var cases = []struct {
		Desc      string
		HasError  bool
		DecayRate float64
		Configurer
	}{
		{
			Desc:       "configurer failed",
			HasError:   true,
			DecayRate:  0.3,
			Configurer: &MockConfigurer{HasError: true},
		},
		{
			Desc:       "success",
			HasError:   false,
			DecayRate:  0.3,
			Configurer: &MockConfigurer{HasError: false, SpeciesDecayRate: 0.3},
		},
	}

	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {

			opt := WithConfiguration(c.Configurer)
			e := new(Experiment)
			err := opt(e)

			if !t.Run("error", test.Error(c.HasError, err)) || c.HasError {
				return
			}

			if c.DecayRate != e.SpeciesDecayRate {
				t.Errorf("incorrect decay rate: expected %f, actual %f", c.DecayRate, e.SpeciesDecayRate)
			}
		})
	}
}

type MockCrosser struct {
	Called   int
	HasError bool
}

func (c *MockCrosser) Cross(parents ...Genome) (child Genome, err error) {
	c.Called++
	if c.HasError {
		err = errors.New("mock transcriber error")
		return
	}
	return
}

func WithMockCrosser() Option {
	return func(e *Experiment) error {
		e.Crosser = &MockCrosser{}
		return nil
	}
}

type MockEvaluator struct {
	Called   int
	HasError bool
	HasSolve bool
}

func (e *MockEvaluator) Evaluate(p Phenome) (r Result, err error) {
	e.Called++
	if e.HasError {
		err = errors.New("mock transcriber error")
		return
	}
	r = Result{
		ID:      p.ID,
		Fitness: 1.0 / float64(p.ID),
		Novelty: float64(p.ID),
		Solved:  e.HasSolve,
	}
	return
}

func WithMockEvaluator() Option {
	return func(e *Experiment) error {
		e.Evaluator = &MockEvaluator{}
		return nil
	}
}

type MockStructureMutator struct {
	Called   int
	Execute  bool
	HasError bool
}

func (m *MockStructureMutator) Mutate(g *Genome) (err error) {
	m.Called++
	if m.HasError {
		err = errors.New("mock structure mutator error")
		return
	}
	// Add a new connection if ID is odd
	if m.Execute {
		g.Encoded.Conns = append(g.Encoded.Conns, Conn{
			Source:  Position{Layer: 0.25, X: 0.25},
			Target:  Position{Layer: 0.75, X: 0.75},
			Weight:  1.23,
			Enabled: true,
		})
	}
	return
}

type MockWeightMutator struct {
	Called   int
	HasError bool
}

func (m *MockWeightMutator) Mutate(g *Genome) (err error) {
	m.Called++
	if m.HasError {
		err = errors.New("mock weight mutator error")
		return
	}
	if len(g.Encoded.Conns) > 0 {
		g.Encoded.Conns[0].Weight += 0.001
	}
	return
}

func WithMockMutators() Option {
	return func(e *Experiment) error {
		e.Mutators = append(e.Mutators, &MockStructureMutator{}, &MockWeightMutator{})
		return nil
	}
}

type MockSeeder struct {
	Called   int
	HasError bool
}

func (m *MockSeeder) Seed() (genomes []Genome, err error) {
	m.Called++
	if m.HasError {
		err = errors.New("mock transcriber error")
		return
	}
	genomes = []Genome{
		{
			ID: 1,
			Encoded: Substrate{
				Nodes: []Node{
					{Position: Position{Layer: 0.0, X: 0.5}, Neuron: Input, Activation: Direct},
					{Position: Position{Layer: 0.5, X: 0.5}, Neuron: Hidden, Activation: InverseAbs},
					{Position: Position{Layer: 1.0, X: 0.5}, Neuron: Output, Activation: Sigmoid},
				},
				Conns: []Conn{
					{Source: Position{Layer: 0.0, X: 0.5}, Target: Position{Layer: 0.0, X: 0.5}, Enabled: false, Weight: 1.5},
					{Source: Position{Layer: 0.5, X: 0.5}, Target: Position{Layer: 1.0, X: 0.5}, Enabled: true, Weight: 2.5},
				},
			},
		},
	}
	return
}

func WithMockSeeder() Option {
	return func(e *Experiment) error {
		e.Seeder = &MockSeeder{}
		return nil
	}
}

type MockRestorer struct {
	Called   int
	HasError bool
}

func (m *MockRestorer) Seed() (genomes []Genome, err error) {
	m.Called++
	if m.HasError {
		err = errors.New("mock transcriber error")
		return
	}
	genomes = []Genome{
		{
			ID:        8,
			SpeciesID: 5,
			Encoded: Substrate{
				Nodes: []Node{
					{Position: Position{Layer: 0.5, X: 0.5}, Neuron: Hidden, Activation: InverseAbs},
					{Position: Position{Layer: 1.0, X: 0.5}, Neuron: Output, Activation: Sigmoid},
				},
				Conns: []Conn{
					{Source: Position{Layer: 0.5, X: 0.5}, Target: Position{Layer: 1.0, X: 0.5}, Enabled: true, Weight: 2.5},
				},
			},
		},
		{
			ID:        9,
			SpeciesID: 6,
			Encoded: Substrate{
				Nodes: []Node{
					{Position: Position{Layer: 0.0, X: 0.5}, Neuron: Input, Activation: Direct},
					{Position: Position{Layer: 0.5, X: 0.5}, Neuron: Hidden, Activation: InverseAbs},
					{Position: Position{Layer: 1.0, X: 0.5}, Neuron: Output, Activation: Sigmoid},
				},
				Conns: []Conn{
					{Source: Position{Layer: 0.0, X: 0.5}, Target: Position{Layer: 0.0, X: 0.5}, Enabled: false, Weight: 1.5},
					{Source: Position{Layer: 0.5, X: 0.5}, Target: Position{Layer: 1.0, X: 0.5}, Enabled: true, Weight: 2.5},
				},
			},
		},
	}

	return
}

func WithMockRestorer() Option {
	return func(e *Experiment) error {
		e.Seeder = &MockRestorer{}
		return nil
	}
}

type MockSearcher struct {
	Called   int
	HasError bool
}

func (s *MockSearcher) Search(e Evaluator, ps []Phenome) (rs []Result, err error) {
	s.Called++
	if s.HasError {
		err = errors.New("mock transcriber error")
		return
	}
	rs = make([]Result, len(ps))
	for i, p := range ps {
		if rs[i], err = e.Evaluate(p); err != nil {
			return
		}
	}
	return
}

func WithMockSearcher() Option {
	return func(e *Experiment) error {
		e.Searcher = &MockSearcher{}
		return nil
	}
}

type MockSelector struct {
	NumContinuing          int
	NumParents             int
	Called                 int
	HasError               bool
	FailWithTooManyParents bool // Special case to check for incorrect population size
}

func (s *MockSelector) Select(p Population) (continuing []Genome, parents [][]Genome, err error) {
	s.Called++
	if s.HasError {
		err = errors.New("mock transcriber error")
		return
	}
	continuing = make([]Genome, s.NumContinuing)
	parents = make([][]Genome, s.NumParents)
	for i := 0; i < len(parents); i++ {
		parents[i] = []Genome{{ID: 1}, {ID: 2}}
	}
	if s.FailWithTooManyParents {
		parents = append(parents, continuing)
	}
	return
}

func WithMockSelector() Option {
	return func(e *Experiment) error {
		e.Selector = &MockSelector{NumContinuing: 2, NumParents: 3}
		return nil
	}
}

type MockSpeciator struct {
	Called    int
	HasError  bool
	HasError2 bool // fail on second iteration
	SetChamp  bool // Sets the champion after intial speciation so we can test the decay call
}

func (s *MockSpeciator) Speciate(pop *Population) (err error) {
	s.Called++
	if s.HasError {
		err = errors.New("mock transcriber error")
		return
	}

	if s.Called > 1 && s.HasError2 {
		err = errors.New("mock transcriber error 2")
		return
	}
	if s.Called > 1 && s.SetChamp {
		// Set to the known best so we can see decay happening
		pop.Species[1].Champion = 1
		pop.Species[1].Decay = 0.9
	}
	// Divide on complexity
	for i, g := range pop.Genomes {
		c := g.Complexity()
		found := false
		for _, s := range pop.Species {
			if s.ID == int64(c) {
				g.SpeciesID = s.ID
				found = true
				break
			}
		}
		if !found {
			s := Species{ID: int64(c)}
			g.SpeciesID = s.ID
			pop.Species = append(pop.Species, s)
		}
		pop.Genomes[i] = g
	}
	return
}

func WithMockSpeciator() Option {
	return func(e *Experiment) error {
		e.Speciator = &MockSpeciator{}
		return nil
	}
}

type MockTranslator struct {
	Called     int
	HasError   bool
	HasMissing bool
}

func (t *MockTranslator) Translate(Substrate) (net Network, err error) {
	t.Called++
	if t.HasError {
		err = errors.New("mock transcriber error")
		return
	}
	if t.HasMissing {
		return nil, nil
	}
	return &MockNetwork{}, nil
}

func WithMockTranslator() Option {
	return func(e *Experiment) error {
		e.Translator = &MockTranslator{}
		return nil
	}
}

type MockNetwork struct{}

func (n *MockNetwork) Activate(Matrix) (Matrix, error) {
	return nil, nil
}

type MockTranscriber struct {
	Called   int
	HasError bool
}

func (t *MockTranscriber) Transcribe(enc Substrate) (dec Substrate, err error) {
	t.Called++
	if t.HasError {
		err = errors.New("mock transcriber error")
		return
	}
	dec = enc // Just return the substrate
	return
}

func WithMockTranscriber() Option {
	return func(e *Experiment) error {
		e.Transcriber = &MockTranscriber{}
		return nil
	}
}

type MockConfigurer struct {
	SpeciesDecayRate float64
	HasError         bool
}

func (c *MockConfigurer) Configure(items ...interface{}) (err error) {
	if c.HasError {
		return errors.New("mock configurer error")
	}
	for _, item := range items {
		if e, ok := item.(*Experiment); ok {
			e.SpeciesDecayRate = c.SpeciesDecayRate
		}
	}
	return
}
