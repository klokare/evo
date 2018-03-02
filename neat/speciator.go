package neat

import (
	"errors"

	"github.com/klokare/evo"
)

// Known errors for the speciator
var (
	ErrMissingDistancer  = errors.New("speciator requires a distancer helper")
	ErrNegativeThreshold = errors.New("compatibility threshold should be nonzero")
)

// Speciator assigns genomes to the correct species, adding and removing species as necessary
type Speciator struct {

	// Properties
	CompatibilityThreshold float64 // Threshold for determining if a genome is compatible with the species
	CompatibilityModifier  float64 // Adjustment to threshold to help achieve target
	TargetSpecies          int     // The desired number of species

	// Helper
	Distancer // Calculates the distance between a genome and the species's example genome

	// State
	lastSID int64
}

// Speciate the population
func (s *Speciator) Speciate(pop *evo.Population) (err error) {

	// Check for known errors
	if s.Distancer == nil {
		err = ErrMissingDistancer
		return
	} else if s.CompatibilityThreshold < 0.0 {
		err = ErrNegativeThreshold
		return
	}

	// Map existing species
	m := make(map[int64][]evo.Genome, len(pop.Species)+5)
	a := make(map[int64]bool, len(pop.Species)+5) // tracks new assignments
	for _, species := range pop.Species {
		m[species.ID] = make([]evo.Genome, 0, 10)
	}

	// Assign genomes to species
	var genomes []evo.Genome
	var ok bool
	for i, genome := range pop.Genomes {

		// The genomes is already assigned an ID and that species exists
		if genomes, ok = m[genome.SpeciesID]; ok {
			genomes = append(genomes, genome)
			m[genome.SpeciesID] = genomes
			continue
		}

		// Look at existing species
		genome.SpeciesID = 0
		for _, species := range pop.Species {
			var d float64
			if d, err = s.Distance(species.Example, genome); err != nil {
				return
			}
			if d < s.CompatibilityThreshold {
				genome.SpeciesID = species.ID
				genomes = m[genome.SpeciesID]
				genomes = append(genomes, genome)
				m[genome.SpeciesID] = genomes
				a[genome.SpeciesID] = true
				break
			}
		}

		// No species found, add a new one
		if genome.SpeciesID == 0 {
			s.lastSID++
			species := evo.Species{
				ID:      s.lastSID,
				Example: genome,
			}
			pop.Species = append(pop.Species, species)
			genome.SpeciesID = species.ID
			genomes := make([]evo.Genome, 0, 10)
			genomes = append(genomes, genome)
			m[genome.SpeciesID] = genomes
			a[genome.SpeciesID] = true
		}

		// Save the genome back to the population
		pop.Genomes[i] = genome
	}

	// Remove empty species
	tmp := pop.Species
	pop.Species = make([]evo.Species, 0, len(m))
	for _, species := range tmp {
		if genomes, ok := m[species.ID]; ok {
			if len(genomes) > 0 {
				pop.Species = append(pop.Species, species)
			}
		}
	}

	// Adjust the compatibile threshold
	if len(a) > s.TargetSpecies {
		s.CompatibilityThreshold += s.CompatibilityModifier
	} else if len(a) < s.TargetSpecies {
		s.CompatibilityThreshold -= s.CompatibilityModifier
		if s.CompatibilityThreshold < s.CompatibilityModifier {
			s.CompatibilityThreshold = s.CompatibilityModifier
		}
	}
	return
}

// WithSpeciator sets the experiment's speciator to a configured NEAT speciator using default distancer
func WithSpeciator(cfg evo.Configurer) evo.Option {
	return func(e *evo.Experiment) (err error) {
		s := new(Speciator)
		d := new(Compatibility)
		if err = cfg.Configure(s, d); err != nil {
			return
		}
		s.Distancer = d
		e.Speciator = s
		return
	}
}
