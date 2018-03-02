package mock

import (
	"errors"

	"github.com/klokare/evo"
)

// Speciator ...
type Speciator struct {
	Called    int
	HasError  bool
	HasError2 bool // fail on second iteration
	SetChamp  bool // Sets the champion after intial speciation so we can test the decay call
}

// Speciate ...
func (s *Speciator) Speciate(pop *evo.Population) (err error) {
	s.Called++
	if s.HasError {
		err = errors.New("mock speciator error")
		return
	}

	if s.Called > 1 && s.HasError2 {
		err = errors.New("mock speciator error 2")
		return
	}
	if s.Called > 1 && s.SetChamp {
		// Set to the known best so we can see decay happening
		pop.Species[1].Champion = 1
		pop.Species[1].Decay = 0.9
	}
	// Divide on complexity
	for i, g := range pop.Genomes {
		c := g.Complexity()
		found := false
		for _, s := range pop.Species {
			if s.ID == int64(c) {
				g.SpeciesID = s.ID
				found = true
				break
			}
		}
		if !found {
			s := evo.Species{ID: int64(c)}
			g.SpeciesID = s.ID
			pop.Species = append(pop.Species, s)
		}
		pop.Genomes[i] = g
	}
	return
}

// WithSpeciator ...
func WithSpeciator() evo.Option {
	return func(e *evo.Experiment) error {
		e.Speciator = &Speciator{}
		return nil
	}
}
