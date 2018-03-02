package evo

import (
	"errors"
	"sort"
	"sync/atomic"

	"github.com/klokare/evo/internal/workers"
)

// MinFitness is the minimum recognised fitness. Anything lower will be set to this. Positive
// fitness is a requirement of the algorithm.
const MinFitness = 1e-10

// An Option modifies the experiment in some way. The most common use is to directly set one or more
// of the experiment's helpers but more sophisticated uses are possible.
type Option func(*Experiment) error

// Known errors
var (
	ErrMissingRequiredHelper        = errors.New("missing required helper")
	ErrMissingNetworkFromTranslator = errors.New("successful translation did not return a network")
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
type Listener func(final bool, pop Population) error

// An Experiment comprises the helpers necessary for creating, evaluating, and advancing a
// population in the search of a solution (or simply a better solver) of a particular problem.
// Experiments are created within the Run method using the options specified.
type Experiment struct {

	// Required helpers
	Crosser
	Evaluator
	Seeder
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
func Batch(r, n int, options ...Option) (pops []Population, errs []error) {

	// Create the tasks
	type task struct {
		pop Population
		err error
	}

	tasks := make([]workers.Task, r)
	for i := 0; i < r; i++ {
		tasks[i] = new(task)
	}

	// Run the tasks
	workers.Do(tasks, func(t workers.Task) {
		x := t.(*task)
		x.pop, x.err = Run(n, options...)
	})

	// Return the populations and errors
	pops = make([]Population, r)
	errs = make([]error, r)
	for i, t := range tasks {
		x := t.(*task)
		pops[i] = x.pop
		errs[i] = x.err
	}
	return
}

// Run performs a single execution of an experiment, defined by the options, for n iterations. The
// final state of the population is returned. Since an experiment may use a different configuration,
// if the it relies on a restored population, only the genomes continue and will be reevaluated.
// The species and previous evaluation results will be discarded.
func Run(n int, options ...Option) (pop Population, err error) {

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

	// Seed the "zero" population. The "zero" population is the one created from the seeds and
	// will be advanced to create the first evaluable population. This is done so that both
	// new and restored populations begin in the same state
	if pop.Genomes, err = e.Seeder.Seed(); err != nil {
		return
	}
	setSequence(e, pop.Genomes) // Determine the next genome ID

	// Speciate the "zero" population
	if err = e.Speciator.Speciate(&pop); err != nil {
		return
	}

	// Iterate the experiment
	solved := false
	for i := 0; !solved && i < n; i++ {

		// This is the last iteration
		final := i == (n - 1)

		// Advance the population
		if err = advance(e.Selector, e.Crosser, e.Mutators, e.Speciator, e.lastGID, &pop); err != nil {
			return
		}
		if err = e.publish(Advanced, solved || final, pop); err != nil {
			return
		}

		// Evaluate the population
		if solved, err = evaluate(e.Searcher, e.Evaluator, e.Transcriber, e.Translator, e.Compare, e.SpeciesDecayRate, &pop); err != nil {
			return
		}
		if err = e.publish(Evaluated, solved || final, pop); err != nil {
			return
		}

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
func (e *Experiment) publish(event Event, final bool, pop Population) (err error) {
	for _, callback := range e.listeners[event] {
		if err = callback(final, pop); err != nil {
			return
		}
	}
	return
}

func verify(e *Experiment) (err error) {
	check := []interface{}{
		e.Crosser, e.Evaluator, e.Seeder, e.Searcher, e.Selector, e.Speciator,
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
	for _, g := range genomes {
		if max < g.ID {
			max = g.ID
		}
	}
	e.lastGID = new(int64)
	*e.lastGID = max
}

// Evaluate the population by transcribing the encoded substrates and then translating the decoded
// versions into phenomes. Those phenomes are sent to the searcher along with the evaluator and
// the results are used to update the genomes and species of the population. Any solution will be
// detected while updating the genomes.
func evaluate(srch Searcher, eval Evaluator, trsc Transcriber, tran Translator, cmp Compare, decay float64, pop *Population) (solved bool, err error) {

	// Create the phenomes
	var phenomes []Phenome
	if phenomes, err = createPhenomes(trsc, tran, pop.Genomes); err != nil {
		return
	}

	// Evaluate the phenomes
	var results []Result
	if results, err = srch.Search(eval, phenomes); err != nil {
		return
	}

	// Update the genomes and species
	solved = updateGenomes(pop.Genomes, results)
	updateSpecies(cmp, decay, pop.Species, pop.Genomes)
	return
}

func createPhenomes(trsc Transcriber, tran Translator, genomes []Genome) (phenomes []Phenome, err error) {

	// Define the tasks
	type task struct {
		genome  Genome
		phenome Phenome
		err     error
	}

	tasks := make([]workers.Task, len(genomes))
	for i, g := range genomes {
		tasks[i] = &task{genome: g}
	}

	// Peform the tasks
	workers.Do(tasks, func(wt workers.Task) {
		t := wt.(*task)
		var dec Substrate
		if dec, t.err = trsc.Transcribe(t.genome.Encoded); t.err != nil {
			return
		}
		var net Network
		if net, t.err = tran.Translate(dec); t.err != nil {
			return
		} else if net == nil {
			t.err = ErrMissingNetworkFromTranslator
			return
		}
		t.phenome = Phenome{
			ID:      t.genome.ID,
			Traits:  make([]float64, len(t.genome.Traits)),
			Network: net,
		}
		copy(t.phenome.Traits, t.genome.Traits)
	})

	// Return the phenomes
	phenomes = make([]Phenome, 0, len(genomes))
	for _, wt := range tasks {
		t := wt.(*task)
		if t.err != nil {
			err = t.err
			return
		}
		phenomes = append(phenomes, t.phenome)
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

	// Species begin to decay but can be reset below with a new champion. Doing it in this loop will
	// update species even if there are no genomes in it.
	sort.Slice(species, func(i, j int) bool { return species[i].ID < species[j].ID })
	for i := 0; i < len(species); i++ {
		species[i].Decay += decay
		if species[i].Decay > 1.0 {
			species[i].Decay = 1.0
		}
	}

	// Sort the genomes first by species then by the comparison function, complexity, and age
	SortBy(genomes,
		func(g1, g2 Genome) int8 {
			// reverse sort the species
			if g1.SpeciesID == g2.SpeciesID {
				return 0
			} else if g1.SpeciesID > g2.SpeciesID {
				return -1
			}
			return 1
		},
		cmp, ByComplexity, ByAge,
	)

	// Iterate the genomes and update the species
	var last int64 = -1
	s := 0
	for i := len(genomes) - 1; i >= 0; i-- {

		// Note the genome
		g := genomes[i]

		// Advance to the correct species
		for species[s].ID < g.SpeciesID {
			s++
		}

		// Update the champion and/or decay
		if species[s].ID != last { // the champion is the first genome we see for a species
			if species[s].Champion != g.ID {
				species[s].Champion = g.ID
				species[s].Decay = 0.0
			}
		}
		last = species[s].ID
	}
}

// Advances the population by one iteration
func advance(sel Selector, crs Crosser, mutators []Mutator, spe Speciator, lastGID *int64, pop *Population) (err error) {

	// Select the genomes to continue to the next iteration and the parents of future offspring
	var continuing []Genome
	var parents [][]Genome
	if continuing, parents, err = sel.Select(*pop); err != nil {
		return
	}

	// Offspring will be created which means a new generation
	if len(parents) > 0 {
		pop.Generation++
	}

	// For each parenting group, create an offspring
	// TODO: make concurrent
	var child Genome
	offspring := make([]Genome, 0, len(parents))
	for _, pgrp := range parents {
		if child, err = makeChild(crs, mutators, atomic.AddInt64(lastGID, 1), pgrp); err != nil {
			return
		}
		offspring = append(offspring, child)
	}

	// Update the population
	pop.Genomes = make([]Genome, 0, len(continuing)+len(offspring))
	pop.Genomes = append(pop.Genomes, continuing...)
	pop.Genomes = append(pop.Genomes, offspring...)

	// Speciate
	if err = spe.Speciate(pop); err != nil {
		return
	}
	return
}

func makeChild(crs Crosser, mutators []Mutator, gid int64, parents []Genome) (child Genome, err error) {
	// Cross the parents to make the offspring
	if child, err = crs.Cross(parents...); err != nil {
		return
	}

	// Assign a new ID
	child.ID = gid

	// Mutate
	n := child.Complexity()
	for _, m := range mutators {
		if err = m.Mutate(&child); err != nil {
			return
		}
		if child.Complexity() != n {
			break // Do not continue to mutate if structure changes. Cannot remember where I read this but there was an admonsiment not to allow other mutations when the structure changes
		}
	}
	return
}

// WithConfiguration configures the experiment using the supplied configurer
func WithConfiguration(cfg Configurer) Option {
	return func(e *Experiment) error {
		return cfg.Configure(e)
	}
}
