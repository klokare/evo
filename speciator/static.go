package speciator

import (
	"fmt"

	"github.com/klokare/evo"
)

// A Static speciator speciates the population without adjusting its compatiblity threshold
type Static struct {
	CompatibilityThreshold float64 `evo:"compatibility-threshold"`
	evo.Comparer
}

func (h Static) String() string {
	return fmt.Sprintf("evo.speciator.Static{CompatibilityThreshold: %f, Comparer: %v}",
		h.CompatibilityThreshold, h.Comparer)
}

// Speciate the population's genomes. If a genome matches a species' example, then it is assigned
// to that species. If no matches occur then a new species is created. Any existing species with
// no members after speciation will be removed.
func (h *Static) Speciate(p *evo.Population) error {

	// Note the last id used and initialise counts
	id := p.Species[0].ID
	cnts := make(map[int]int, len(p.Species))
	for _, s := range p.Species {
		if id < s.ID {
			id = s.ID
		}
		cnts[s.ID] = 0
	}

	// Iterate the genomes
	for i, g := range p.Genomes {

		// Look for an existing species match
		found := false
		for _, s := range p.Species {
			d, err := h.Comparer.Compare(g.Encoded, s.Example)
			if err != nil {
				return err
			}
			if d < h.CompatibilityThreshold {
				cnts[s.ID]++
				p.Genomes[i].SpeciesID = s.ID
				found = true
				break
			}
		}

		// No species found, create a new one
		if !found {
			id++
			s := evo.Species{ID: id, Example: g.Encoded}
			p.Species = append(p.Species, s)
			cnts[s.ID]++
			p.Genomes[i].SpeciesID = s.ID
		}
	}

	// Remove empty species
	for id, cnt := range cnts {
		if cnt == 0 {
			for i := 0; i < len(p.Species); i++ {
				if p.Species[i].ID == id {
					p.Species = append(p.Species[:i], p.Species[i+1:]...)
					break
				}
			}
		}
	}

	return nil
}
