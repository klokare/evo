package mutator

import (
	"errors"
	"log"

	"github.com/klokare/evo"
)

// Phased mutator
type Phased struct {

	// Properties
	PhaseThreshold float64
	HoldPhase      int // Number of generations to hold phase if not improving
	evo.Compare

	// Composite mutators
	Complexify
	Simplify

	// Interanal state
	simplifying     bool
	threshold       float64
	champion        int64 // last known champion. used to delay simplification if improving
	stagnant        int   // number of generations that no new champion has been found
	stuck           int   // number of generations since the minMpc has fallen
	lastGen         int
	lastMpc, minMpc float64
	toggler
}

// Mutate the genome by either adding to or removing from its structure, depending on the phase
func (p Phased) Mutate(g *evo.Genome) (err error) {
	if p.simplifying {
		return p.Simplify.Mutate(g)
	}
	return p.Complexify.Mutate(g)
}

// Update the phase as needed
func (p *Phased) update(final bool, pop evo.Population) (err error) {

	// Wait for a generation change
	if pop.Generation == p.lastGen {
		return
	}
	p.lastGen = pop.Generation

	// Calculate the mean complexity
	sum := 0
	for _, g := range pop.Genomes {
		sum += g.Complexity()
	}
	mpc := float64(sum) / float64(len(pop.Genomes))
	log.Println("mpc", mpc, "threshold", p.threshold, "simplifying", p.simplifying, "min", p.minMpc, "stuck", p.stuck, "champ", p.champion, "stag", p.stagnant)

	// This is the initial call, set the threshold and continue
	if p.threshold == 0.0 {
		p.threshold = mpc + p.PhaseThreshold
		return
	}

	// If fitness is improving, then continue in phase
	evo.SortBy(pop.Genomes, evo.BySolved, p.Compare, evo.ByComplexity, evo.ByAge)
	b := pop.Genomes[len(pop.Genomes)-1]
	if p.champion == b.ID {
		p.stagnant++
	} else {
		p.champion = b.ID
		p.stagnant = 0
	}

	// Attempt to switch phase
	if p.simplifying {
		if mpc < p.minMpc {
			p.minMpc = mpc
			p.stuck = 0
		} else {
			p.stuck++
		}
		if p.stuck > p.HoldPhase {
			log.Println("complexifying")
			p.simplifying = false
			p.threshold = p.minMpc + p.PhaseThreshold
			if p.toggler != nil {
				p.ToggleMutateOnly(false)
			}
		}
	} else {
		if mpc >= p.threshold && p.stagnant > p.HoldPhase {
			log.Println("simplifying")
			p.simplifying = true
			p.stuck = 0
			p.minMpc = mpc
			if p.toggler != nil {
				p.ToggleMutateOnly(true)
			}
		}
	}
	return
}

// WithComplexify adds a configured complexify mutator to the experiment
func WithPhased(cfg evo.Configurer) evo.Option {
	return func(e *evo.Experiment) (err error) {
		zc := new(Complexify)
		zs := new(Simplify)
		zp := new(Phased)
		if err = cfg.Configure(zc, zs, zp); err != nil {
			return
		}

		// Do not continue if there is no chance for mutation
		if zc.AddNodeProbability == 0.0 && zc.AddConnProbability == 0.0 && zs.DelNodeProbability == 0.0 && zs.DelConnProbability == 0.0 {
			return
		}

		// There is only complexify
		if zs.DelNodeProbability == 0.0 && zs.DelConnProbability == 0.0 {
			e.Mutators = append(e.Mutators, zc)
			return
		}

		// There is only simplify
		if zc.AddNodeProbability == 0.0 && zc.AddConnProbability == 0.0 {
			e.Mutators = append(e.Mutators, zs)
			return
		}

		// Phase threshold not set, just use composite mutators
		if zp.PhaseThreshold == 0 {
			e.Mutators = append(e.Mutators, zc, zs)
			return
		}

		// Assemble the phased mutator
		if e.Compare == nil {
			err = errors.New("experiment should have compare function set before adding phased mutator")
			return
		}
		zp.Compare = e.Compare
		zp.Complexify = *zc
		zp.Simplify = *zs
		if x, ok := e.Selector.(toggler); ok {
			zp.toggler = x
		}
		e.Mutators = append(e.Mutators, zp)
		e.Subscribe(evo.Evaluated, zp.update)
		return
	}
}

type toggler interface {
	ToggleMutateOnly(bool) error
}
