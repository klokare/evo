package neat

import (
	"testing"

	"github.com/klokare/evo"
	"github.com/klokare/evo/internal/mock"
	"github.com/klokare/evo/internal/test"
)

func TestWithOptions(t *testing.T) {

	// Configurer has an error
	t.Run("configurer error", func(t *testing.T) {

		// Retrieve the options
		cfg := &mock.Configurer{HasError: true}
		options := WithOptions(cfg)

		// Configure an experiment
		var err error
		e := new(evo.Experiment)
		for _, option := range options {
			err = option(e)
			if err != nil {
				break
			}
		}
		if !t.Run("error", test.Error(true, err)) {
			return
		}
	})

	// Configurer does not have an error
	t.Run("no configurer error", func(t *testing.T) {

		// Retrieve the options
		cfg := &mock.Configurer{
			HasError:                false,
			MutateBiasProbability:   1.0,
			MutateWeightProbability: 1.0,
			AddConnProbability:      1.0,
			AddNodeProbability:      1.0,
		}
		options := WithOptions(cfg)

		// Configure an experiment
		e := new(evo.Experiment)
		for _, option := range options {
			err := option(e)
			if !t.Run("error", test.Error(false, err)) {
				return
			}
		}

		// Check that all the options have been run
		if e.Compare == nil {
			t.Errorf("compare function should be set")
		}
		if e.Crosser == nil {
			t.Errorf("crosser should be set")
		}
		if e.Seeder == nil {
			t.Errorf("seeder should be set")
		}
		if e.Searcher == nil {
			t.Errorf("searcher should be set")
		}
		if e.Selector == nil {
			t.Errorf("selector should be set")
		}
		if e.Speciator == nil {
			t.Errorf("speciator should be set")
		}
		if e.Translator == nil {
			t.Errorf("translator should be set")
		}
		if e.Transcriber == nil {
			t.Errorf("transcriber should be set")
		}

		// Check that the mutators are all present
		// TODO: Add config for traits, simplify, phased
		if len(e.Mutators) != 3 {
			t.Errorf("incorrect number of mutators: expected %d, actual %d", 3, len(e.Mutators))
		}

	})
}
