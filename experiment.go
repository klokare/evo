package evo

import (
	"context"
	"errors"
	"runtime"
	"sort"
	"sync"
	"sync/atomic"
)

// MinFitness is the minimum recognised fitness. Anything lower will be set to this. Positive
// fitness is a requirement of the algorithm.
const MinFitness = 1e-10

// An Option modifies the experiment in some way. The most common use is to directly set one or more
// of the experiment's helpers but more sophisticated uses are possible.
type Option func(*Experiment) error

// Known errors
var (
	ErrMissingRequiredHelper = errors.New("missing required helper")
)

// Event is key used with a Listener
type Event byte

// Events associated with the Experiment
const (
	Evaluated Event = iota + 1
	Advanced
)

// Listener functions are called when the event to which they are subscribed occurs. The final flag
// is true when the experiment is solved or on its final iteration
type Listener func(ctx context.Context, final bool, pop Population) error

// An Experiment comprises the helpers necessary for creating, evaluating, and advancing a
// population in the search of a solution (or simply a better solver) of a particular problem.
// Experiments are created within the Run method using the options specified.
type Experiment struct {

	// Required helpers
	Crosser
	Evaluator
	Populator
	Searcher
	Selector
	Speciator
	Translator
	Transcriber
	Mutators []Mutator

	// External methods
	Compare

	// Properties
	SpeciesDecayRate float64 // Increment amount [0,1], per iteration, to decay a species without improvement

	// Internal state
	lastGID   *int64
	listeners map[Event][]Listener
}

// Batch performs r or more runs of an experiment, each created with n iterations and using the
// options. Batch retuns each of the resulting populations.
func Batch(ctx context.Context, r, n int, options ...Option) (pops []Population, errs []error) {

	ch := make(chan int)
	go func(ch chan int) {
		for i := 0; i < r; i++ {
			ch <- i
		}
		close(ch)
	}(ch)

	// Execute each run
	pops = make([]Population, r)
	errs = make([]error, r)
	wg := new(sync.WaitGroup)
	for w := 0; w < runtime.NumCPU(); w++ {
		wg.Add(1)
		go func(ch <-chan int) {
			defer wg.Done()
			for i := range ch {
				pops[i], errs[i] = Run(ctx, n, options...)
			}
		}(ch)
	}
	wg.Wait()
	return
}

// Run performs a single execution of an experiment, defined by the options, for n iterations. The
// final state of the population is returned. Since an experiment may use a different configuration,
// if the it relies on a restored population, only the genomes continue and will be reevaluated.
// The species and previous evaluation results will be discarded.
func Run(ctx context.Context, n int, options ...Option) (pop Population, err error) {

	// Run the experiment in separate go rountine so we can handle the context closing early
	done := make(chan struct{})
	go func(done chan struct{}) {

		// Close the channel when exiting
		defer close(done)

		// Create a new experiment by applying the options
		e := &Experiment{
			listeners: make(map[Event][]Listener, 2),
		}
		for _, option := range options {
			if err = option(e); err != nil {
				return
			}
		}

		// Validate that the required options are present
		if err = verify(e); err != nil {
			return
		}

		// Create the initial population
		if pop, err = e.Populator.Populate(ctx); err != nil {
			return
		}

		// Set the sequence and reset the genomes
		setSequence(e, pop.Genomes)

		// Speciate the population
		if err = e.Speciator.Speciate(ctx, &pop); err != nil {
			return
		}

		// Iterate the experiment
		solved := false
		for i := 0; !solved && i < n; i++ {

			// This is the last iteration
			final := i == (n - 1)

			// Evaluate the population
			if solved, err = evaluate(ctx, e.Searcher, e.Evaluator, e.Transcriber, e.Translator, e.Compare, e.SpeciesDecayRate, &pop); err != nil {
				return
			}
			if err = e.publish(ctx, Evaluated, solved || final, pop); err != nil {
				return
			}

			// Advance the population
			if !solved && !final {
				if err = advance(ctx, e.Selector, e.Crosser, e.Mutators, e.Speciator, e.lastGID, &pop); err != nil {
					return
				}
				if err = e.publish(ctx, Advanced, solved || final, pop); err != nil {
					return
				}
			}
		}
	}(done)

	// Experiment is complete, return
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case <-done:
		// experiment ended normally
	}
	return
}

// Subscribe a callback function to a particular event. Note: this is not concurrent safe which is
// OK because this is really only used via Options
func (e *Experiment) Subscribe(event Event, callback Listener) {
	var callbacks []Listener
	var ok bool
	if callbacks, ok = e.listeners[event]; !ok {
		callbacks = make([]Listener, 0, 10)
	}
	callbacks = append(callbacks, callback)
	e.listeners[event] = callbacks
}

// Publish an event to the listeners.
// TODO: make concurrent
func (e *Experiment) publish(ctx context.Context, event Event, final bool, pop Population) (err error) {
	for _, callback := range e.listeners[event] {
		if err = callback(ctx, final, pop); err != nil {
			return
		}
	}
	return
}

func verify(e *Experiment) (err error) {
	check := []interface{}{
		e.Crosser, e.Evaluator, e.Populator, e.Searcher, e.Selector, e.Speciator,
		e.Translator, e.Transcriber,
	}
	for _, h := range check {
		if h == nil {
			err = ErrMissingRequiredHelper
			return
		}
	}
	if len(e.Mutators) == 0 {
		err = ErrMissingRequiredHelper
		return
	}
	if e.Compare == nil {
		err = ErrMissingRequiredHelper
		return
	}
	return
}

func setSequence(e *Experiment, genomes []Genome) {
	var max int64
	for i, g := range genomes {

		// This is a higher ID
		if max < g.ID {
			max = g.ID
		}

		// Reset the genome
		g.Fitness = 0.0
		g.Novelty = 0.0
		g.Solved = false
		g.SpeciesID = 0
		genomes[i] = g
	}

	e.lastGID = new(int64)
	*e.lastGID = max
}

// Evaluate the population by transcribing the encoded substrates and then translating the decoded
// versions into phenomes. Those phenomes are sent to the searcher along with the evaluator and
// the results are used to update the genomes and species of the population. Any solution will be
// detected while updating the genomes.
func evaluate(ctx context.Context, srch Searcher, eval Evaluator, trsc Transcriber, tran Translator, cmp Compare, decay float64, pop *Population) (solved bool, err error) {

	// Create the phenomes
	var phenomes []Phenome
	if phenomes, err = createPhenomes(ctx, trsc, tran, pop.Genomes); err != nil {
		return
	}

	// Evaluate the phenomes
	var results []Result
	if results, err = srch.Search(ctx, eval, phenomes); err != nil {
		return
	}

	// Update the genomes and species
	solved = updateGenomes(pop.Genomes, results)
	updateSpecies(cmp, decay, pop.Species, pop.Genomes)
	return
}

func createPhenomes(ctx context.Context, trsc Transcriber, tran Translator, genomes []Genome) (phenomes []Phenome, err error) {
	phenomes = make([]Phenome, len(genomes))
	for i, g := range genomes {
		var dec Substrate
		if dec, err = trsc.Transcribe(ctx, g.Encoded); err != nil {
			return
		}
		var net Network
		if net, err = tran.Translate(ctx, dec); err != nil {
			return
		}
		phenomes[i] = Phenome{
			ID:      g.ID,
			Traits:  make([]float64, len(g.Traits)),
			Network: net,
		}
		copy(phenomes[i].Traits, g.Traits)
	}
	return
}

func updateGenomes(genomes []Genome, results []Result) (solved bool) {
	sort.Slice(results, func(i, j int) bool { return results[i].ID < results[j].ID })
	for i, g := range genomes {
		idx := sort.Search(len(results), func(i int) bool { return results[i].ID >= g.ID })
		if idx < len(results) && results[idx].ID == g.ID {
			g.Fitness = results[idx].Fitness
			g.Novelty = results[idx].Novelty
			g.Solved = results[idx].Solved
			solved = solved || g.Solved
		}
		genomes[i] = g
	}
	return
}

// Update the champion of the species. This will be the "best" genome, according to the Comparer,
// which belongs to the species.
func updateSpecies(cmp Compare, decay float64, species []Species, genomes []Genome) {

	// Separate the genomes by species
	var gs []Genome
	var ok bool
	m := make(map[int64][]Genome, len(species))
	for _, g := range genomes {
		if gs, ok = m[g.SpeciesID]; !ok {
			gs = make([]Genome, 0, len(genomes))
		}
		gs = append(gs, g)
		m[g.SpeciesID] = gs
	}

	// Iterate the species
	for i, s := range species {

		// Determine the current champion
		gs = m[s.ID]
		SortBy(gs, BySolved, cmp, ByComplexity, ByAge) // This produces and order with no ties
		champ := gs[len(gs)-1]

		// Update the decay and champion properties
		if s.Champion != champ.ID {
			// New champion, reset decay
			s.Decay = 0.0
			s.Champion = champ.ID
		} else {
			s.Decay += decay
			if s.Decay > 1.0 {
				s.Decay = 1.0
			}
		}
		species[i] = s
	}
}

// Advances the population by one iteration
func advance(ctx context.Context, sel Selector, crs Crosser, mutators []Mutator, spe Speciator, lastGID *int64, pop *Population) (err error) {

	// Select the genomes to continue to the next iteration and the parents of future offspring
	var continuing []Genome
	var parents [][]Genome
	if continuing, parents, err = sel.Select(ctx, *pop); err != nil {
		return
	}

	// Offspring will be created which means a new generation
	if len(parents) > 0 {
		pop.Generation++
	}

	// For each parenting group, create an offspring
	offspring := make([]Genome, 0, len(parents))
	for _, pgrp := range parents {

		// Cross the parents to make the offspring
		var child Genome
		if child, err = crs.Cross(ctx, pgrp...); err != nil {
			return
		}

		// Assign a new ID
		child.ID = atomic.AddInt64(lastGID, 1)

		// Mutate
		n := child.Complexity()
		for _, m := range mutators {
			if err = m.Mutate(ctx, &child); err != nil {
				return
			}
			if child.Complexity() != n {
				break // Do not continue to mutate if structure changes. Cannot remember where I read this but there was an admonsiment not to do allow other mutations when the structure changes
			}
		}

		// Save the child
		offspring = append(offspring, child)
	}

	// Update the population
	pop.Genomes = make([]Genome, 0, len(continuing)+len(offspring))
	pop.Genomes = append(pop.Genomes, continuing...)
	pop.Genomes = append(pop.Genomes, offspring...)

	// Speciate
	if err = spe.Speciate(ctx, pop); err != nil {
		return
	}
	return
}

// WithConfiguration configures the experiment using the supplied configurer
func WithConfiguration(cfg Configurer) Option {
	return func(e *Experiment) error {
		return cfg.Configure(e)
	}
}
