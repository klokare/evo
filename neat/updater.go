package neat

import (
	"errors"
	"sort"

	"github.com/klokare/evo"
)

// Known errors
var (
	ErrInvalidDecayRate = errors.New("species decay rate must not be negative")
	ErrMissingCompare   = errors.New("compare function is required")
)

// Updater updates a population's genomes and species from a given result set
type Updater struct {
	SpeciesDecayRate float64
	evo.Comparison
}

// Update the population using the reults
func (u *Updater) Update(pop *evo.Population, results []evo.Result) (err error) {

	// Check for errors
	if u.SpeciesDecayRate < 0.0 {
		err = ErrInvalidDecayRate
		return
	} else if u.Comparison == 0 {
		err = ErrMissingCompare
		return
	}

	// Update the genomes
	sort.Slice(results, func(i, j int) bool { return results[i].ID < results[j].ID })
	for i, g := range pop.Genomes {
		idx := sort.Search(len(results), func(i int) bool { return results[i].ID >= g.ID })
		if idx < len(results) && results[idx].ID == g.ID {
			g.Fitness = results[idx].Fitness
			g.Novelty = results[idx].Novelty
			g.Solved = results[idx].Solved
		}
		pop.Genomes[i] = g
	}

	// Species begin to decay but can be reset below with a new champion. Doing it in this loop will
	// update species even if there are no genomes in it.
	sort.Slice(pop.Species, func(i, j int) bool { return pop.Species[i].ID < pop.Species[j].ID })
	for i := 0; i < len(pop.Species); i++ {
		pop.Species[i].Decay += u.SpeciesDecayRate
		if pop.Species[i].Decay > 1.0 {
			pop.Species[i].Decay = 1.0
		}
	}

	// Sort the genomes first by species then by the comparison function, complexity, and age
	evo.SortBy(pop.Genomes, evo.BySpecies, u.Comparison, evo.ByComplexity, evo.ByAge)

	// Iterate the genomes and update the species
	var last int64 = -1
	s := 0
	for i := len(pop.Genomes) - 1; i >= 0; i-- {

		// Note the genome
		g := pop.Genomes[i]

		// Advance to the correct species
		for pop.Species[s].ID < g.SpeciesID {
			s++
		}

		// Update the champion and/or decay
		if pop.Species[s].ID != last { // the champion is the first genome we see for a species
			if pop.Species[s].Champion != g.ID {
				pop.Species[s].Champion = g.ID
				pop.Species[s].Decay = 0.0
			}
		}
		last = pop.Species[s].ID
	}
	return
}
