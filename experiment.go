package evo

import (
	"fmt"

	"github.com/klokare/errors"
)

// An Experiment encapsulates the helpers and settings for running
type Experiment struct {

	// Helpers
	Crosser
	Mutator
	Searcher
	Selector
	Speciater
	Transcriber
	Translator
	Watcher

	// Internal sequences
	genomeID  int
	speciesID int
}

// Run one trial of the experiment starting with the population and running for n iterations
// Check via watcher. Generation of last should be n unless solved is trigger then should be that generation. Use population of 10 to make this easy. Fitness = ID
// Checks: each helper is called
func Run(e *Experiment, p Population, n int) error {
	var err error
	for i := 0; i < n; i++ {

		// Initialise the experiment or advance the population
		if i == 0 {
			e.initIDs(p)
		} else {
			if err = e.advance(&p); err != nil {
				return err
			}
		}

		// Decode the genomes
		if err = e.transcribe(p.Genomes); err != nil {
			return err
		}

		var ps []Phenome
		if ps, err = e.translate(p.Genomes); err != nil {
			return err
		}

		// Search the phenomes and update the population
		var rs []Result
		if rs, err = e.Searcher.Search(ps); err != nil {
			return err
		}

		var stop bool
		if stop, err = update(p.Genomes, rs); err != nil {
			return err
		}
		stagnate(&p)

		// Inform the watchers of the iteration
		if e.Watcher != nil {
			if err = e.Watcher.Watch(p); err != nil {
				return err
			}
		}

		// A solution was found
		if stop {
			break
		}
	}

	return nil
}

// Initialise the ID sequences used in the experiment with the ones already used in the population
func (e *Experiment) initIDs(p Population) {
	for _, g := range p.Genomes {
		if e.genomeID < g.ID {
			e.genomeID = g.ID
		}
	}
	for _, s := range p.Species {
		if e.speciesID < s.ID {
			e.speciesID = s.ID
		}
	}
}

// Advance the experiment for one iteration. If parents are selected, this will trigger a increment
// in generation.
// Check: fatal error if the population does not have same number of genomes. If parents, speciater called, generation updated
func (e *Experiment) advance(p *Population) error {

	// Select which genomes to keep and which to become parents
	curr := *p
	keep, parents, err := e.Selector.Select(*p)
	if err != nil {
		return err
	}

	// Begin the new population
	p.Genomes = make([]Genome, len(keep), len(curr.Genomes))
	copy(p.Genomes, keep)

	// There is a new generation
	if len(parents) > 0 {

		// Procreate
		var os []Genome
		if os, err = e.procreate(parents); err != nil {
			return err
		}
		if len(os)+len(keep) != len(curr.Genomes) {
			return fmt.Errorf("insufficient offspring created: %d (wanted %d)", len(os), len(curr.Genomes))
		}
		p.Genomes = append(p.Genomes, os...)

		// Speciate the new generation
		if err = e.Speciater.Speciate(p); err != nil {
			return err
		}

		// Advance the generation
		p.Generation++
	}

	return nil
}

// Creates offspring from the parent groupings
func (e *Experiment) procreate(gss [][]Genome) ([]Genome, error) {
	var err error
	os := make([]Genome, 0, len(gss))
	z := new(errors.Safe)
	for _, gs := range gss {
		var o Genome
		if o, err = e.Crosser.Cross(gs...); err != nil {
			z.Add(err)
		}
		e.genomeID++
		o.ID = e.genomeID
		if err = e.Mutator.Mutate(&o); err != nil {
			z.Add(err)
		}
		os = append(os, o)
	}
	return os, z.Err()
}

// Transcribes encoded genomes into substrates that can be tranlsated into networks
func (e *Experiment) transcribe(gs []Genome) error {
	var err error
	z := new(errors.Safe)
	for i, g := range gs {
		if gs[i].Decoded, err = e.Transcriber.Transcribe(g.Encoded); err != nil {
			z.Add(fmt.Errorf("error transcribing genome %d: %v", g.ID, err))
		}
	}
	return z.Err()
}

// Translates decoded genomes into phenomes
func (e *Experiment) translate(gs []Genome) ([]Phenome, error) {
	var err error
	ps := make([]Phenome, len(gs))
	z := new(errors.Safe)
	for i, g := range gs {
		ps[i].ID = g.ID
		ps[i].Traits = make([]float64, len(g.Traits))
		copy(ps[i].Traits, g.Traits)
		if g.Decoded.Complexity() > 0 {
			if ps[i].Network, err = e.Translator.Translate(g.Decoded); err != nil {
				z.Add(fmt.Errorf("error translating genome %d: %v", g.ID, err))
			}
		} else {
			z.Add(fmt.Errorf("genome %d has empty decoded substrate", g.ID))
		}
	}
	return ps, z.Err()
}

// Updates the genomes with the results
// Check, each genome in result receives update. Solved is returned if declare (ID=50)
// What to do with result error? If simply bad phenome, genome's fitness = MinFitness. Otherwise propogate if trouble running experiment
func update(gs []Genome, rs []Result) (bool, error) {

	// Map genomes by id
	i2g := make(map[int]*Genome, len(gs))
	for i := 0; i < len(gs); i++ {
		i2g[gs[i].ID] = &gs[i]
	}

	// Iterate the results
	var solved bool
	z := new(errors.Safe)
	for _, r := range rs {

		// Set the fitness, novelty and solved
		g := i2g[r.ID]
		g.Fitness = r.Fitness
		g.Novelty = r.Novelty
		g.Solved = r.Solved

		// Update the overall
		solved = solved || g.Solved
		if r.Error != nil {
			z.Add(r.Error)
		}
	}
	return solved, z.Err()
}

// Updates the species' fitness and stagnation
func stagnate(p *Population) {

	// Map species by ID
	i2s := make(map[int]*Species)
	for i := 0; i < len(p.Species); i++ {
		i2s[p.Species[i].ID] = &p.Species[i]
	}

	// Map genomes to species
	g2s := make(map[int]Genomes, len(i2s))
	for _, g := range p.Genomes {
		g2s[g.SpeciesID] = append(g2s[g.SpeciesID], g)
	}

	// Update the species
	for i, s := range i2s {
		f := g2s[i].MaxFitness()
		if f > s.Fitness {
			s.Fitness = f
			s.Stagnation = 0
		} else {
			s.Stagnation++
		}
	}
}
