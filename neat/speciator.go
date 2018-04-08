package neat

import (
	"errors"
	"sort"

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
	lastSID  int
	examples map[int]evo.Genome
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
	a := make(map[int]bool, len(s.examples)+5) // tracks new assignments
	n := make(map[int]int, len(s.examples)+5)

	// TRIAL: use local cache so we can phase out species as a separate object
	// Ensure examples for existing species
	var ok bool
	if s.examples == nil {
		s.examples = make(map[int]evo.Genome, 20)
	}
	for _, g := range pop.Genomes {
		if g.Species > 0 {
			if _, ok = s.examples[g.Species]; !ok {
				s.examples[g.Species] = g
			}
			if s.lastSID < g.Species {
				s.lastSID = g.Species
			}
		}
	}

	// Extract IDs into a sorted list
	sids := make([]int, 0, len(s.examples)+5)
	for sid := range s.examples {
		sids = append(sids, sid)
	}
	sort.Slice(sids, func(i, j int) bool { return sids[i] < sids[j] })

	// Assign genomes to species
	for i, genome := range pop.Genomes {

		// The genomes is already assigned an ID and that species exists
		if _, ok = s.examples[genome.Species]; ok {
			n[genome.Species]++
			continue
		}

		// Look at existing species.
		// Interestingly, the number of XOR failures rises 5x if exmaples are iterated randomly
		// (like using range over Go map). Presenting older (lower IDs) first prevents this.
		genome.Species = 0
		// for sid, example := range s.examples {
		for _, sid := range sids {
			example := s.examples[sid]

			var d float64
			if d, err = s.Distance(example, genome); err != nil {
				return
			}
			if d < s.CompatibilityThreshold {
				genome.Species = sid
				a[genome.Species] = true
				n[genome.Species]++
				break
			}
		}

		// No species found, add a new one
		if genome.Species == 0 {
			s.lastSID++
			genome.Species = s.lastSID
			a[genome.Species] = true
			n[genome.Species]++
			s.examples[genome.Species] = genome
			sids = append(sids, genome.Species)
		}

		// Save the genome back to the population
		pop.Genomes[i] = genome
	}

	// Remove empty species
	for sid, cnt := range n {
		if cnt == 0 {
			delete(s.examples, sid)
		}
	}

	// Adjust the compatible threshold
	if len(a) > s.TargetSpecies {
		s.CompatibilityThreshold += s.CompatibilityModifier
	} else if s.CompatibilityThreshold > s.CompatibilityModifier && len(a) < s.TargetSpecies {
		s.CompatibilityThreshold -= s.CompatibilityModifier
		if s.CompatibilityThreshold < s.CompatibilityModifier {
			s.CompatibilityThreshold = s.CompatibilityModifier
		}
	}
	return
}
