package evo

import (
	"bytes"
	"fmt"

	"github.com/klokare/errors"
)

// Comparer returns the "distance" between two substrates
type Comparer interface {
	Compare(Substrate, Substrate) (float64, error)
}

// Configurer configures a helper
type Configurer interface {
	Configure(interface{}) error
}

// Crosser creates a new genome by crossing parents
type Crosser interface {
	Cross(...Genome) (Genome, error)
}

// Evaluator evaluates a single phenome and returns a result
type Evaluator interface {
	Evaluate(Phenome) (Result, error)
}

// Mutator mutates a genome's encoded structure
type Mutator interface {
	Mutate(*Genome) error
}

// Populater provides the intitial population for the experiment
type Populater interface {
	Populate() (Population, error)
}

// Searcher evalutes all phenomes and collects the results
type Searcher interface {
	Search([]Phenome) ([]Result, error)
}

// Selector selects which genomes are kept and which become parents
type Selector interface {
	Select(Population) ([]Genome, [][]Genome, error)
}

// Speciater divides the population's genome into species
type Speciater interface {
	Speciate(*Population) error
}

// Transcriber creates a decoded substrate from an encoded one
type Transcriber interface {
	Transcribe(Substrate) (Substrate, error)
}

// Translator creates a network from a substrate
type Translator interface {
	Translate(Substrate) (Network, error)
}

// Watcher watches a population and is informed after each iteration
type Watcher interface {
	Watch(Population) error
}

// Mutators is a collection of mutator helpers that can be called as one.
type Mutators []Mutator

func (h Mutators) String() string {
	b := bytes.NewBufferString("evo.Mutators:")
	for i, m := range h {
		b.WriteString(fmt.Sprintf(" [%d] %v", i, m))
	}
	return b.String()
}

// Mutate the genome with the collected mutators. The order of mutators matters as their execution
// stops if one of the mutators changes the genomes structure.
func (h Mutators) Mutate(g *Genome) error {
	c := g.Complexity()
	for _, m := range h {
		if err := m.Mutate(g); err != nil {
			return err
		}
		if c != g.Complexity() {
			break
		}
	}
	return nil
}

// Watchers is a collections of watcher helpers that can be called as one.
type Watchers []Watcher

// Watch informs the helper of an iteration and passes the current population. Order does not
// matter as all watchers are called.
// TODO: make concurrent
// TODO: Handle errors from watchers
func (h Watchers) Watch(p Population) error {
	e := new(errors.Safe)
	for _, w := range h {
		if err := w.Watch(p); err != nil {
			e.Add(err)
		}
	}
	return e.Err()
}

func (h Watchers) String() string {
	b := bytes.NewBufferString("evo.Watchers:")
	for i, w := range h {
		b.WriteString(fmt.Sprintf(" [%d] %v", i, w))
	}
	return b.String()
}
