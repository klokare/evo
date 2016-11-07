package evo

import (
	"bytes"
	"fmt"
	"math"
	"sync"

	"github.com/klokare/errors"
)

// An Experiment groups individuals trials under a single description
type Experiment struct {

	// Properties
	Description string `evo:"description"`
	Trials      int    `evo:"trials"`
	NewTrial    func(i int) (*Trial, error)
}

func (e Experiment) String() string {
	return fmt.Sprintf("evo.Experiment{Description: %s, Trials: %d}", e.Description, e.Trials)
}

// A Trial is a single run of an experiment. It encapsulates all the helpers needed for processing.
type Trial struct {
	Iterations int `evo:"iterations"`
	Stopped    bool

	// Helpers
	Crosser
	Mutator
	Populater
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

func (t Trial) String() string {
	b := bytes.NewBufferString("Experiment:\n")
	b.WriteString(fmt.Sprintf("... Iterations:  %d\n", t.Iterations))
	b.WriteString(fmt.Sprintf("... Stopped:     %v\n", t.Stopped))
	b.WriteString(fmt.Sprintf("... Crosser:     %v\n", t.Crosser))
	b.WriteString(fmt.Sprintf("... Crosser:     %v\n", t.Mutator))
	b.WriteString(fmt.Sprintf("... Populater:   %v\n", t.Populater))
	b.WriteString(fmt.Sprintf("... Searcher:    %v\n", t.Searcher))
	b.WriteString(fmt.Sprintf("... Selector:    %v\n", t.Selector))
	b.WriteString(fmt.Sprintf("... Speciater:   %v\n", t.Speciater))
	b.WriteString(fmt.Sprintf("... Transcriber: %v\n", t.Transcriber))
	b.WriteString(fmt.Sprintf("... Translator:  %v\n", t.Translator))
	b.WriteString(fmt.Sprintf("... Watcher:     %v\n", t.Watcher))
	return b.String()
}

// Run one or more trials of the experiment starting with the population and running for n iterations
func Run(e *Experiment) (err error) {

	// Set the experiment

	// Iterate the trials
	if e.Trials == 0 {
		e.Trials = 1
	}
	for j := 0; j < e.Trials; j++ {

		// Create a new trial
		var t *Trial
		if t, err = e.NewTrial(j); err != nil {
			return
		}

		// Set the trial

		// Infinite iterations requestd. Let's hope the evaluator has a solved condition.
		if t.Iterations <= 0 {
			t.Iterations = math.MaxInt64 // Effectively infinite
		}

		// Iterate
		var p Population
		for i := 0; !t.Stopped && i < t.Iterations; i++ {

			// Initialise the experiment or advance the population
			if i == 0 {
				if p, err = t.Populater.Populate(); err != nil {
					return
				}
				t.initIDs(p)
			} else {
				if err = t.advance(&p); err != nil {
					return err
				}
			}

			// Decode the genomes
			if err = t.transcribe(p.Genomes); err != nil {
				return err
			}

			var ps []Phenome
			if ps, err = t.translate(p.Genomes); err != nil {
				return err
			}

			// Search the phenomes and update the population
			var rs []Result
			if rs, err = t.Searcher.Search(ps); err != nil {
				return err
			}

			var stop bool
			if stop, err = update(p.Genomes, rs); err != nil {
				return err
			}
			stagnate(&p)

			// Inform the watchers of the iteration
			if t.Watcher != nil {
				if err = t.Watcher.Watch(p); err != nil {
					return err
				}
			}

			// A solution was found
			if stop {
				break
			}
		}
	}

	return nil
}

// Initialise the ID sequences used in the trial with the ones already used in the population
func (t *Trial) initIDs(p Population) {
	for _, g := range p.Genomes {
		if t.genomeID < g.ID {
			t.genomeID = g.ID
		}
	}
	for _, s := range p.Species {
		if t.speciesID < s.ID {
			t.speciesID = s.ID
		}
	}
}

// Advance the trial for one iteration. If parents are selected, this will trigger a increment
// in generation.
// Check: fatal error if the population does not have same number of genomes. If parents, speciater called, generation updated
func (t *Trial) advance(p *Population) error {

	// Select which genomes to keep and which to become parents
	curr := *p
	keep, parents, err := t.Selector.Select(*p)
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
		if os, err = t.procreate(parents); err != nil {
			return err
		}
		if len(os)+len(keep) != len(curr.Genomes) {
			return fmt.Errorf("insufficient offspring created: %d (wanted %d)", len(os), len(curr.Genomes))
		}
		p.Genomes = append(p.Genomes, os...)

		// Speciate the new generation
		if err = t.Speciater.Speciate(p); err != nil {
			return err
		}

		// Advance the generation
		p.Generation++
	}

	return nil
}

// Creates offspring from the parent groupings
func (t *Trial) procreate(gss [][]Genome) ([]Genome, error) {
	var err error
	os := make([]Genome, 0, len(gss))
	z := new(errors.Safe)
	for _, gs := range gss {
		var o Genome
		if o, err = t.Crosser.Cross(gs...); err != nil {
			z.Add(err)
		}
		t.genomeID++
		o.ID = t.genomeID
		if err = t.Mutator.Mutate(&o); err != nil {
			z.Add(err)
		}
		os = append(os, o)
	}
	return os, z.Err()
}

// Transcribes encoded genomes into substrates that can be tranlsated into networks
func (t *Trial) transcribe(gs []Genome) error {
	var err error
	z := new(errors.Safe)
	wg := new(sync.WaitGroup)
	for i, g := range gs {
		wg.Add(1)
		go func(i int, g Genome) {
			if gs[i].Decoded, err = t.Transcriber.Transcribe(g.Encoded); err != nil {
				z.Add(fmt.Errorf("error transcribing genome %d: %v", g.ID, err))
			}
			wg.Done()
		}(i, g)
	}
	wg.Wait()
	return z.Err()
}

// Translates decoded genomes into phenomes
func (t *Trial) translate(gs []Genome) ([]Phenome, error) {
	var err error
	ps := make([]Phenome, len(gs))
	z := new(errors.Safe)
	for i, g := range gs {
		ps[i].ID = g.ID
		ps[i].Traits = make([]float64, len(g.Traits))
		copy(ps[i].Traits, g.Traits)
		if g.Decoded.Complexity() > 0 {
			if ps[i].Network, err = t.Translator.Translate(g.Decoded); err != nil {
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
