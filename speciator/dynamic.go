package speciator

import (
	"fmt"
	"math"

	"github.com/klokare/evo"
)

// A Dynamic speciator adjusts the compatibility threshold of the inner, Static speciator to
// achieve the target number of species
type Dynamic struct {
	Static

	TargetSpecies         int     `evo:"target-species"`
	CompatibilityModifier float64 `evo:"compatibility-modifier"`
}

func (h Dynamic) String() string {
	return fmt.Sprintf("evo.speciator.Dynamic{TargetSpecies: %d, CompatibilityModifier: %f, Static: %v}",
		h.TargetSpecies, h.CompatibilityModifier, h.Comparer)
}

// Speciate the population's genomes using the inner speciator and update its compatibility
// threshold based on the number of species found.
func (h *Dynamic) Speciate(p *evo.Population) error {

	// Speciate using the inner helper
	if err := h.Static.Speciate(p); err != nil {
		return err
	}

	// Adjust the compatibility threshold as necessary
	if len(p.Species) < h.TargetSpecies {
		h.Static.CompatibilityThreshold = math.Max(h.CompatibilityModifier,
			h.Static.CompatibilityThreshold-h.CompatibilityModifier)
	} else if len(p.Species) > h.TargetSpecies {
		h.Static.CompatibilityThreshold += h.CompatibilityModifier
	}

	return nil
}
