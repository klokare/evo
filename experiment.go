package evo

import (
	"context"
	"errors"
	"sync/atomic"

	"github.com/klokare/evo/internal/workers"
)

// Known errors
var (
	ErrMissingNetworkFromTranslator = errors.New("successful translation did not return a network")
	ErrNoSeedGenomes                = errors.New("seeder produced no genomes")
)

// An Experiment comprises the helpers necessary for creating, evaluating, and advancing a
// population in the search of a solution (or simply a better solver) of a particular problem.
type Experiment interface {
	Crosser
	Mutator
	Populator
	Searcher
	Selector
	Speciator
	Transcriber
	Translator
	Updater
}

// Run the experiment in the given context with the evalutor. The context will decide when the
// experiment ends. See IterationContext and TimoutContext functions. An error is returned if
// any of the composite helpers' methods return an error.
func Run(ctx context.Context, exp Experiment, eval Evaluator) (pop Population, err error) {

	// The experiment provides subscribers so subscribe them
	listeners := make(map[Event][]Callback, 10)
	if sx, ok := exp.(SubscriptionProvider); ok {
		for _, s := range sx.Subscriptions() {
			var ls []Callback
			if ls, ok = listeners[s.Event]; !ok {
				ls = make([]Callback, 0, 10)
			}
			ls = append(ls, s.Callback)
			listeners[s.Event] = ls
		}
	}

	// Create the initial population
	if pop, err = exp.Populate(); err != nil {
		return
	}
	if len(pop.Genomes) == 0 {
		err = ErrNoSeedGenomes
		return
	}
	lastGID := setSequence(pop.Genomes) // Determine the next genome ID

	// Ensure every genome belongs to a species
	if err = exp.Speciate(&pop); err != nil {
		return
	}

	// Inform listeners that the population has started
	if err = publish(listeners, Started, pop); err != nil {
		return
	}

	// Iterate the experiment
	for {

		// Select the continuing genomes and those who will become parents
		var continuing []Genome
		var parents [][]Genome
		if continuing, parents, err = exp.Select(pop); err != nil {
			return
		}
		if len(parents) > 0 {
			pop.Generation++
		}

		// Create the population
		var offspring []Genome
		if offspring, err = createOffspring(exp, lastGID, parents); err != nil {
			return
		}
		pop.Genomes = make([]Genome, 0, len(continuing)+len(offspring))
		pop.Genomes = append(pop.Genomes, continuing...)
		pop.Genomes = append(pop.Genomes, offspring...)

		// Speciate the genomes
		if err = exp.Speciate(&pop); err != nil {
			return
		}

		// Inform listeners that the population has been advanced
		if err = publish(listeners, Advanced, pop); err != nil {
			return
		}

		// Decode the genomes into phenomes
		var phenomes []Phenome
		if phenomes, err = decodeGenomes(exp, pop.Genomes); err != nil {
			return
		}

		// Inform listeners that decoding has completed
		if err = publish(listeners, Decoded, pop); err != nil {
			return
		}

		// Search the problem with the phenomes
		var results []Result
		if results, err = exp.Search(eval, phenomes); err != nil {
			return
		}

		// Update the population with the results
		if err = exp.Update(&pop, results); err != nil {
			return
		}

		// Inform listeners that evaluation has completed
		if err = publish(listeners, Evaluated, pop); err != nil {
			return
		}

		// Check for completion
		select {
		case <-ctx.Done():
			err = publish(listeners, Completed, pop)
			return
		default:
			// continue to next iteration
		}
	}
}

// Determine the starting sequence number for genome IDs
func setSequence(genomes []Genome) (lastGID *int64) {
	var max int64
	for _, g := range genomes {
		if max < g.ID {
			max = g.ID
		}
	}
	lastGID = new(int64)
	*lastGID = max
	return
}

type progenator interface {
	Crosser
	Mutator
}

// Create the offspring from the parents, mutate the children, and set their IDs.
func createOffspring(helper progenator, lastGID *int64, parents [][]Genome) (offspring []Genome, err error) {

	// Receive offspring
	offspring = make([]Genome, 0, len(parents))
	ch := make(chan Genome, len(parents))
	done := make(chan struct{})
	go func(ch <-chan Genome, done chan struct{}) {
		defer close(done)
		for g := range ch {
			offspring = append(offspring, g)
		}
	}(ch, done)

	// Create the tasks
	tasks := make([]workers.Task, len(parents))
	for i, pgrp := range parents {
		tasks[i] = pgrp
	}

	// Do the work
	err = workers.Do(tasks, func(wt workers.Task) (err error) {
		pgrp := wt.([]Genome)
		// Create the child
		var child Genome
		if child, err = helper.Cross(pgrp...); err != nil {
			return
		}
		child.ID = atomic.AddInt64(lastGID, 1) // Assign the next ID

		// Mutate the child and add to the list
		if err = helper.Mutate(&child); err != nil {
			return
		}
		ch <- child
		return
	})
	close(ch)
	<-done
	return
}

type decoder interface {
	Transcriber
	Translator
}

// Decode the genomes into phenomes
func decodeGenomes(dec decoder, genomes []Genome) (phenomes []Phenome, err error) {

	// Receive phenomes
	phenomes = make([]Phenome, 0, len(genomes))
	ch := make(chan Phenome, len(genomes))
	done := make(chan struct{})
	go func(ch <-chan Phenome, done chan struct{}) {
		defer close(done)
		for p := range ch {
			phenomes = append(phenomes, p)
		}
	}(ch, done)

	// Create the tasks
	tasks := make([]workers.Task, len(genomes))
	for i := 0; i < len(genomes); i++ {
		tasks[i] = &genomes[i]
	}

	// Do the work
	err = workers.Do(tasks, func(wt workers.Task) (err error) {

		// Decode the encoded substrate
		g := wt.(*Genome)
		if g.Decoded.Complexity() == 0 {
			if g.Decoded, err = dec.Transcribe(g.Encoded); err != nil {
				return
			}
		}

		// Create the neural network
		var net Network
		if net, err = dec.Translate(g.Decoded); err != nil {
			return
		}

		// Create the phenome and add to the list
		p := Phenome{
			ID:      g.ID,
			Network: net,
			Traits:  make([]float64, len(g.Traits)),
		}
		copy(p.Traits, g.Traits)
		ch <- p
		return
	})
	close(ch)
	<-done
	return
}
