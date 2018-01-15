package evo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"
)

// Things to test in Run
// 1. If context does not end early and there are no errors or solution, it should run for n iterations
// 2. If solution is found, the experiment should end early returning after the whole population's results are processed
// 3. If the context completes before iterations are done, the experiment should end early
func TestExperimentRun(t *testing.T) {

	var options = []Option{
		WithConfiguration(&MockConfigurer{}),
		WithCompare(ByFitness),
		WithMockCrosser(),
		WithMockEvaluator(),
		WithMockPopulator(),
		WithMockSearcher(),
		WithMockSelector(),
		WithMockSpeciator(),
		WithMockTranscriber(),
		WithMockTranslator(),
		WithMockMutators(),
	}

	// There error from a failed option should be propogated up
	t.Run("failed option", testFailedOption(options))

	// There should be an error if the required helpers are not present. We can test this by
	// iterating the options, skipping each one incrementally, and testing a Run.
	t.Run("missing helper", testMissingHelper(options))

	// When running with a new population, IDs should be setand speciator called
	t.Run("initial population", testInitialPopulation(options))

	// When running with a restored popuation, IDs should be set and speciator called
	t.Run("restored population", testRestoredPopulation(options))

	// If the iterations > 0, evalaute should be called
	t.Run("evaluate", testEvaluate(options))

	// If the iterations = 1, advance should be not be called
	t.Run("no advance", testNoAdvance(options))

	// If the iterations > 1, advance should be called
	t.Run("advance", testAdvance(options))

	// If there is no solution or errors, the experiment should run for N iterations
	t.Run("complete iterations", testIterations(options))

	// There are subscriptions
	t.Run("subscriptions", testSubscriptions(options))

	// There are subscriptions but one fails
	t.Run("failed callbacks", testFailedCallback(options))

	// If the context ends, the experiment should stop stop early.
	t.Run("context ending", testEndContext(options))

	// If a solution is found, the experiment should stop early
	t.Run("solution found", testSolutionFound(options))

	// If any of the major helpers fails, the experiment should stop early with an error.
	t.Run("helper errors", testHelperErrors(options))

}

func testFailedOption(options []Option) func(*testing.T) {
	return func(t *testing.T) {
		failing := make([]Option, len(options))
		copy(failing, options)
		failing = append(failing, WithFailedOption())
		_, err := Run(context.Background(), 10, failing...)
		if err == nil {
			t.Error("error for failed option expected.")
		} else if err.Error() != "option failed" {
			t.Errorf("unexpected error returned: expected \"%s\", actual \"%s\"", "option failed", err)
		}
	}
}

func WithFailedOption() Option {
	return func(e *Experiment) error {
		return errors.New("option failed")
	}
}

func testMissingHelper(original []Option) func(*testing.T) {
	return func(t *testing.T) {

		// Work with a copy of the options
		options := make([]Option, len(original))
		copy(options, original)

		for i := 1; i < len(options); i++ { // Skip the configurer
			missing := make([]Option, len(options))
			copy(missing, options)
			missing = append(missing[:i], missing[i+1:]...)
			_, err := Run(context.Background(), 10, missing...)
			if err == nil {
				t.Errorf("error for missing helpers expected when dropping %d.", i)
			} else if err != ErrMissingRequiredHelper {
				t.Errorf("unexpected error returned for option %d: expected \"%s\", actual \"%s\"",
					i, ErrMissingRequiredHelper, err)
			}
		}
	}
}

// When run with an initial population and zero iterations should return the initial population,
// after being speciated. The internal ID sequence for new genomes should be greater than the
// last genome of the seed population.
func testInitialPopulation(original []Option) func(*testing.T) {
	return func(t *testing.T) {

		// Work with a copy of the options
		options := make([]Option, len(original))
		copy(options, original)

		// Add an option to capture the experiment
		var z *Experiment
		options = append(options, func(e *Experiment) error {
			z = e
			return nil
		})

		// Expected population is the initial population from the mock populator.
		exp, _ := new(MockPopulator).Populate(context.Background())

		// Run the experiment for 0 iterations
		pop, err := Run(context.Background(), 0, options...)
		if err != nil {
			t.Errorf("unxpected error: %v", err)
			return
		}

		// The popoulation size should be the same
		if len(pop.Genomes) != len(exp.Genomes) {
			t.Errorf("incorrect number of genomes: expected %d, actual %d", len(exp.Genomes), len(pop.Genomes))
		}

		// The species IDs should be assigned. Initial population from the mock populator sends in
		// genomes with species ID unset so we can just look for IDs > 0 as the mock speciator
		// starts species' IDs at 1.
		for _, g := range pop.Genomes {
			if g.SpeciesID == 0 {
				t.Errorf("incorrect species assignment. expected >0, actual 0")
			}
		}

		// The experiment's internal genome ID sequence should be greater than or equal the greatest
		// from the original population
		var max int64
		for _, g := range exp.Genomes {
			if max < g.ID {
				max = g.ID
			}
		}
		if z.lastGID == nil {
			t.Errorf("genome sequence not initialised")
		} else if *z.lastGID < max {
			t.Errorf("incorrect genome ID sequence: expected %d or greater, actual %d", max, *z.lastGID)
		}
	}
}

// Basically the same as the intial population test.
func testRestoredPopulation(original []Option) func(*testing.T) {
	return func(t *testing.T) {

		// Work with a copy of the options
		options := make([]Option, len(original))
		copy(options, original)

		// Add an option to capture the experiment
		var z *Experiment
		options = append(options, func(e *Experiment) error {
			z = e
			return nil
		})

		// Add the restorer to use instead of the populator
		options = append(options, WithMockRestorer())

		// Expected population is the initial population from the mock populator.
		exp, _ := new(MockRestorer).Populate(context.Background())

		// Run the experiment for 0 iterations
		pop, err := Run(context.Background(), 0, options...)
		if err != nil {
			t.Errorf("unxpected error: %v", err)
			return
		}

		// The popoulation size should be the same
		if len(pop.Genomes) != len(exp.Genomes) {
			t.Errorf("incorrect number of genomes: expected %d, actual %d", len(exp.Genomes), len(pop.Genomes))
		}

		// The species IDs should be assigned. Initial population from the mock populator sends in
		// genomes with species ID unset so we can just look for IDs > 0 as the mock speciator
		// starts species' IDs at 1.
		for _, g := range pop.Genomes {
			if g.SpeciesID == 0 {
				t.Errorf("incorrect species assignment. expected >0, actual 0")
			}
		}

		// The experiment's internal genome ID sequence should be greater than or equal the greatest
		// from the original population
		var max int64
		for _, g := range exp.Genomes {
			if max < g.ID {
				max = g.ID
			}
		}
		if z.lastGID == nil {
			t.Errorf("genome sequence not initialised")
		} else if *z.lastGID < max {
			t.Errorf("incorrect genome ID sequence: expected %d or greater, actual %d", max, *z.lastGID)
		}
	}
}

// test evaluation routine
func testEvaluate(original []Option) func(*testing.T) {
	return func(t *testing.T) {

		// Work with a copy of the options
		options := make([]Option, len(original))
		copy(options, original)

		// Add an option to capture the experiment
		var z *Experiment
		options = append(options, func(e *Experiment) error {
			z = e
			return nil
		})

		// Take a copy of the initial population and record number of genomes and species
		pop, _ := new(MockPopulator).Populate(context.Background())
		gcnt := len(pop.Genomes)

		new(MockSpeciator).Speciate(context.Background(), &pop)
		scnt := len(pop.Species)

		// Run the experiment for 1 iteration
		pop, _ = Run(context.Background(), 1, options...)

		// Evaluation should not change the number of genomes or species
		if len(pop.Genomes) != gcnt {
			t.Errorf("genome count should not have changed: expected %d, actual %d", gcnt, len(pop.Genomes))
		}
		if len(pop.Species) != scnt {
			t.Errorf("species count should not have changed: expected %d, actual %d", scnt, len(pop.Species))
		}

		// The transcriber should be called for each genome (check count since IDs are not passed)
		trsc := z.Transcriber.(*MockTranscriber)
		if trsc.Called != gcnt {
			t.Errorf("incorrect number of calls to transcriber: expected %d, actual %d", gcnt, trsc.Called)
		}

		// The translator should be called for each genome (check count since IDs are not passed)
		// TODO: If the transcriber produces a substrate for a valid, but non-functioning network
		// (e.g., no connections), should the translator still be invoked?
		tran := z.Translator.(*MockTranslator)
		if tran.Called != gcnt {
			t.Errorf("incorrect number of calls to translator: expected %d, actual %d", gcnt, tran.Called)
		}
		// The searcher should be called
		srch := z.Searcher.(*MockSearcher)
		if srch.Called == 0 {
			t.Errorf("incorrect number of calls to searcher: expected: 1, actual: 0")
		}

		// The evaluator should be called for each genome
		// TODO: If the translator was not invoked, should the genome be evaluated or just given
		// the default minimum fitness and zero novelty?
		eval := z.Evaluator.(*MockEvaluator)
		if eval.Called == 0 {
			t.Errorf("incorrect number of calls to evaluator: expected: %d, actual: %d", gcnt, eval.Called)
		}

		// Genomes should be updated. MockEvaluator sets the fitness to 1 / ID and novelty to ID.
		// This iteration does not include a solution
		for _, g := range pop.Genomes {
			fit := 1.0 / float64(g.ID)
			if fit != g.Fitness {
				t.Errorf("incorrect fitness for genome %d: expected %f, actual %f", g.ID, fit, g.Fitness)
			}
			nov := float64(g.ID)
			if nov != g.Novelty {
				t.Errorf("incorrect novelty for genome %d: expected %f, actual %f", g.ID, nov, g.Novelty)
			}
		}

		// Species champion should reflect the most fit (that's the comparer we are using)
		for _, s := range pop.Species {

			// Find the species champion in the population
			var c Genome
			for _, g := range pop.Genomes {
				if g.ID == s.Champion {
					c = g
					break
				}
			}
			if c.ID == 0 {
				t.Errorf("incorrect champion id for species %d: expected %d, actual 0", s.ID, s.Champion)
			}

			// There should not be another genome in the species that is more fit
			for _, g := range pop.Genomes {
				if g.SpeciesID == s.ID && g.ID != c.ID {
					x := z.Compare(c, g)
					if x < 0 {
						t.Errorf("incorrect champion selected for species %d: expected %d, actual %d", s.ID, g.ID, c.ID)
					}
				}
			}
		}
	}
}

// tests that advance has not been called
func testNoAdvance(original []Option) func(*testing.T) {
	return func(t *testing.T) {

		// Work with a copy of the options
		options := make([]Option, len(original))
		copy(options, original)

		// Add an option to capture the experiment
		var z *Experiment
		options = append(options, func(e *Experiment) error {
			z = e
			return nil
		})

		// No helpers related to advance should be called
		_, _ = Run(context.Background(), 1, options...)

		if z.Selector.(*MockSelector).Called > 0 {
			t.Errorf("selector should not be called")
		}
		if z.Crosser.(*MockCrosser).Called > 0 {
			t.Errorf("crosser should not be called")
		}
		if z.Mutators[0].(*MockStructureMutator).Called > 0 {
			t.Errorf("mutator 0 should not be called")
		}
		if z.Mutators[1].(*MockWeightMutator).Called > 0 {
			t.Errorf("mutator 1 should not be called")
		}
		if z.Selector.(*MockSelector).Called > 0 {
			t.Errorf("selector should not be called")
		}
		if z.Speciator.(*MockSpeciator).Called > 1 { // Called once when experiment starts
			t.Errorf("speciator should not be called")
		}
	}
}

func testAdvance(original []Option) func(*testing.T) {
	return func(t *testing.T) {

		// Work with a copy of the options
		options := make([]Option, len(original))
		copy(options, original)

		// Add an option to capture the experiment
		var z *Experiment
		options = append(options, func(e *Experiment) error {
			z = e
			return nil
		})

		// Remember what the initial population looks like
		var maxID int64
		pop, _ := new(MockPopulator).Populate(context.Background())
		size := len(pop.Genomes)
		for _, g := range pop.Genomes {
			if maxID < g.ID {
				maxID = g.ID
			}
		}

		// Run the experiment for 2 iterations
		pop, _ = Run(context.Background(), 2, options...)

		// Population should have same size
		if len(pop.Genomes) != size {
			t.Errorf("incorrect population size: expected %d, actual %d", size, len(pop.Genomes))
		}

		// Selector should have been called
		if z.Selector.(*MockSelector).Called == 0 {
			t.Errorf("selector should have be called")
		}

		// NOTE: The mock version continues the first 2 genomes and makes the remaining genomes single parents

		//Crosser should have been called for each parenting group
		if z.Crosser.(*MockCrosser).Called != size-2 {
			t.Errorf("incorrect number of calls to crosser: expected %d, actual %d", size-2, z.Crosser.(*MockCrosser).Called)
		}

		// Structure mutator should have been called evey time (first in our list)
		if z.Mutators[0].(*MockStructureMutator).Called != size-2 {
			t.Errorf("incorrect number of calls to structure mutator: expected %d, actual %d", size-2, z.Mutators[0].(*MockStructureMutator).Called)
		}

		// Weight mutator should only be called for the even ID genomes as a structural mutation should stop others
		if z.Mutators[1].(*MockWeightMutator).Called != (size-2)/2 {
			t.Errorf("incorrect number of calls to weight mutator: expected %d, actual %d", (size-2)/2, z.Mutators[1].(*MockWeightMutator).Called)
		}

		// Speciator should have been called
		if z.Speciator.(*MockSpeciator).Called < 2 { // Called once after initial population
			t.Errorf("selector should have be called")
		}
	}
}

func testIterations(original []Option) func(*testing.T) {
	return func(t *testing.T) {
		// Work with a copy of the options
		options := make([]Option, len(original))
		copy(options, original)

		// Add an option to capture the experiment
		var z *Experiment
		options = append(options, func(e *Experiment) error {
			z = e
			return nil
		})

		// Run the experiment
		n := 5
		_, _ = Run(context.Background(), n, options...)

		// The searcher should have been called the same number of times as iterations
		if z.Searcher.(*MockSearcher).Called != n {
			t.Errorf("incorrect number of calls for searcher: expected %d, actual %d", n, z.Searcher.(*MockSearcher).Called)
		}

	}
}

func testHelperErrors(original []Option) func(*testing.T) {
	return func(t *testing.T) {
		t.Run("crosser", testHelperError(original, func(e *Experiment) error {
			e.Crosser = &MockCrosser{HasError: true}
			return nil
		}))
		t.Run("evaluator", testHelperError(original, func(e *Experiment) error {
			e.Evaluator = &MockEvaluator{HasError: true}
			return nil
		}))
		t.Run("populator", testHelperError(original, func(e *Experiment) error {
			e.Populator = &MockPopulator{HasError: true}
			return nil
		}))
		t.Run("searcher", testHelperError(original, func(e *Experiment) error {
			e.Searcher = &MockSearcher{HasError: true}
			return nil
		}))
		t.Run("selector", testHelperError(original, func(e *Experiment) error {
			e.Selector = &MockSelector{HasError: true}
			return nil
		}))
		t.Run("speciator initial", testHelperError(original, func(e *Experiment) error {
			e.Speciator = &MockSpeciator{HasError: true}
			return nil
		}))
		t.Run("speciator advance", testHelperError(original, func(e *Experiment) error {
			e.Speciator = &MockSpeciator{HasError2: true}
			return nil
		}))
		t.Run("transcriber", testHelperError(original, func(e *Experiment) error {
			e.Transcriber = &MockTranscriber{HasError: true}
			return nil
		}))
		t.Run("translator", testHelperError(original, func(e *Experiment) error {
			e.Translator = &MockTranslator{HasError: true}
			return nil
		}))
		t.Run("mutator 0", testHelperError(original, func(e *Experiment) error {
			e.Mutators[0] = &MockStructureMutator{HasError: true}
			return nil
		}))
		t.Run("mutator 1", testHelperError(original, func(e *Experiment) error {
			e.Mutators[1] = &MockWeightMutator{HasError: true}
			return nil
		}))
	}
}

func testHelperError(original []Option, fail Option) func(*testing.T) {
	return func(t *testing.T) {
		options := make([]Option, len(original), len(original)+1)
		copy(options, original)
		options = append(options, fail)
		_, err := Run(context.Background(), 2, options...)
		if err == nil {
			t.Errorf("expected error")
		}
	}
}

func testSubscriptions(original []Option) func(*testing.T) {
	return func(t *testing.T) {

		// Work with a copy of the options
		options := make([]Option, len(original))
		copy(options, original)

		// Add options to subscribe
		var evaluated, advanced int
		options = append(options,
			func(e *Experiment) error {
				e.Subscribe(Evaluated, func(context.Context, bool, Population) error {
					evaluated++
					return nil
				})
				return nil
			},
			func(e *Experiment) error {
				e.Subscribe(Advanced, func(context.Context, bool, Population) error {
					advanced++
					return nil
				})
				return nil
			},
		)

		// Run the experiment
		_, _ = Run(context.Background(), 2, options...)

		// The listeners should be called
		if evaluated != 2 {
			t.Errorf("incorrect count for evaluated: expected 2, actual %d", evaluated)
		}
		if advanced != 1 {
			t.Errorf("incorrect count for advanced: expected 1, actual %d", advanced)
		}
	}
}

func testFailedCallback(original []Option) func(*testing.T) {
	return func(t *testing.T) {
		for _, event := range []Event{Evaluated, Advanced} {
			// Work with a copy of the options
			options := make([]Option, len(original))
			copy(options, original)

			// Add options to subscribe
			options = append(options,
				func(e *Experiment) error {
					e.Subscribe(event, func(context.Context, bool, Population) error {
						return errors.New("failed callback")
					})
					return nil
				},
			)

			// Run the experiment
			_, err := Run(context.Background(), 2, options...)
			if err == nil {
				t.Errorf("expected error from callback")
			}
		}
	}
}

func testSolutionFound(original []Option) func(*testing.T) {
	return func(t *testing.T) {

		// Work with a copy of the options
		options := make([]Option, len(original))
		copy(options, original)

		// Add an option to capture the experiment
		var z *Experiment
		options = append(options, func(e *Experiment) error {
			z = e
			return nil
		})

		// Add option that solves the experiment on second iteration
		options = append(options, func(e *Experiment) error {
			e.Evaluator = &MockEvaluator{HasSolve: true}
			return nil
		})

		// Run the experiment
		_, _ = Run(context.Background(), 5, options...)

		// The searcher should have been called the same number of times as iterations
		if z.Searcher.(*MockSearcher).Called == 5 {
			t.Errorf("incorrect number of calls for searcher: expected <5, actual 5")
		}

	}
}

func testEndContext(original []Option) func(*testing.T) {
	return func(t *testing.T) {

		// Work with a copy of the options
		options := make([]Option, len(original))
		copy(options, original)

		// Add an option to capture the experiment
		var z *Experiment
		options = append(options, func(e *Experiment) error {
			z = e
			return nil
		})

		// Create a cancelable context and add an option that cancels it after the first evaluation
		ctx, cancel := context.WithCancel(context.Background())
		options = append(options, func(e *Experiment) error {
			e.Subscribe(Evaluated, func(context.Context, bool, Population) error {
				cancel()
				time.Sleep(1)
				return nil
			})
			return nil
		})

		// Run the experiment
		_, err := Run(ctx, 5, options...)

		// There should be an error from the context
		if err == nil {
			t.Errorf("error from canceled context expected")
		}

		// The searcher should have been called the same number of times as iterations
		if z.Searcher.(*MockSearcher).Called != 1 {
			t.Errorf("incorrect number of calls for searcher: expected 1, actual %d", z.Searcher.(*MockSearcher).Called)
		}
	}
}

type MockCrosser struct {
	Called   int
	HasError bool
}

func (c *MockCrosser) Cross(ctx context.Context, parents ...Genome) (child Genome, err error) {
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

func (e *MockEvaluator) Evaluate(ctx context.Context, p Phenome) (r Result, err error) {
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
	HasError bool
}

func (m *MockStructureMutator) Mutate(ctx context.Context, g *Genome) (err error) {
	m.Called++
	if m.HasError {
		err = errors.New("mock structure mutator error")
		return
	}
	// Add a new connection if ID is odd
	if g.ID%2 == 1 {
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

func (m *MockWeightMutator) Mutate(ctx context.Context, g *Genome) (err error) {
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

type MockPopulator struct {
	Called   int
	HasError bool
}

func (m *MockPopulator) Populate(context.Context) (p Population, err error) {
	m.Called++
	if m.HasError {
		err = errors.New("mock transcriber error")
		return
	}
	p = Population{
		Species: []Species{{ID: 5}}, // Species do not persist in load
		Genomes: []Genome{
			{
				ID: 3,
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
				ID: 2,
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
			{
				ID: 1,
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
				ID: 4,
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
		},
	}
	return
}

func WithMockPopulator() Option {
	return func(e *Experiment) error {
		e.Populator = &MockPopulator{}
		return nil
	}
}

type MockRestorer struct {
	Called   int
	HasError bool
}

func (m *MockRestorer) Populate(context.Context) (p Population, err error) {
	m.Called++
	if m.HasError {
		err = errors.New("mock transcriber error")
		return
	}
	p = Population{
		Species: []Species{
			{ID: 5},
			{ID: 6},
		}, // Species do not persist in load
		Genomes: []Genome{
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
		},
	}
	return
}

func WithMockRestorer() Option {
	return func(e *Experiment) error {
		e.Populator = &MockRestorer{}
		return nil
	}
}

type MockSearcher struct {
	Called   int
	HasError bool
}

func (s *MockSearcher) Search(ctx context.Context, e Evaluator, ps []Phenome) (rs []Result, err error) {
	s.Called++
	if s.HasError {
		err = errors.New("mock transcriber error")
		return
	}
	rs = make([]Result, len(ps))
	for i, p := range ps {
		if rs[i], err = e.Evaluate(ctx, p); err != nil {
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
	Called                 int
	HasError               bool
	FailWithTooManyParents bool // Special case to check for incorrect population size
}

func (s *MockSelector) Select(ctx context.Context, p Population) (continuing []Genome, parents [][]Genome, err error) {
	s.Called++
	if s.HasError {
		err = errors.New("mock transcriber error")
		return
	}
	continuing = p.Genomes[:2]
	parents = make([][]Genome, len(p.Genomes)-2)
	for i := 0; i < len(parents); i++ {
		parents[i] = []Genome{p.Genomes[i+2]}
	}
	if s.FailWithTooManyParents {
		parents = append(parents, continuing)
	}
	return
}

func WithMockSelector() Option {
	return func(e *Experiment) error {
		e.Selector = &MockSelector{}
		return nil
	}
}

type MockSpeciator struct {
	Called    int
	HasError  bool
	HasError2 bool // fail on second iteration
	SetChamp  bool // Sets the champion after intial speciation so we can test the decay call
}

func (s *MockSpeciator) Speciate(ctx context.Context, pop *Population) (err error) {
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
	Called   int
	HasError bool
}

func (t *MockTranslator) Translate(context.Context, Substrate) (net Network, err error) {
	t.Called++
	if t.HasError {
		err = errors.New("mock transcriber error")
		return
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

func (n *MockNetwork) Activate(inputs []float64) (outputs []float64, err error) {
	return
}

type MockTranscriber struct {
	Called   int
	HasError bool
}

func (t *MockTranscriber) Transcribe(ctx context.Context, enc Substrate) (dec Substrate, err error) {
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

type MockConfigurer struct{}

func (c *MockConfigurer) Configure(items ...interface{}) (err error) {
	data := `{"SpeciesDecayRate": 0.3}`
	for _, item := range items {
		if err = json.NewDecoder(bytes.NewBufferString(data)).Decode(item); err != nil {
			return
		}
	}
	return
}
