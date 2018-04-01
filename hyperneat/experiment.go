package hyperneat

import (
	"github.com/klokare/evo"
	"github.com/klokare/evo/config"
	"github.com/klokare/evo/neat"
	"github.com/klokare/evo/neat/mutator"
	"github.com/klokare/evo/network/forward"
)

// Ensure the experiment struct implements the experiment interface
var (
	_ evo.Experiment = &Experiment{}
)

// Experiment builds on the NEAT experiment by adding HyperNEAT specific helpers
type Experiment struct {
	neat.Experiment
	Transcriber
	Seeder
}

// NewExperiment creates a new Hyper-NEAT experiment using the configuration. Configurations employ the
// maximum namespace so user can be as specific or lax (depending on depth of namespace used) as
// desired.
func NewExperiment(cfg config.Configurer) (exp *Experiment) {

	// Create the HyperNEAT experiment
	exp = new(Experiment)
	exp.Experiment = *neat.NewExperiment(cfg) // backfill with the NEAT helpers
	exp.Experiment.Populator = neat.Populator{
		Seeder: Seeder{
			NumTraits:         cfg.Int("hyperneat|seeder|num-traits"),
			DisconnectRate:    cfg.Float64("hyperneat|seeder|disconnect-rate"),
			SeedLocalityLayer: cfg.Bool("hyperneat|transcriber|seed-locality-layer"),
			SeedLocalityX:     cfg.Bool("hyperneat|transcriber|seed-locality-x"),
			SeedLocalityY:     cfg.Bool("hyperneat|transcriber|seed-locality-y"),
			SeedLocalityZ:     cfg.Bool("hyperneat|transcriber|seed-locality-z"),
		},
		PopulationSize: cfg.Int("neat|populator|population-size"),
		WeightPower:    cfg.Float64("neat|populator|weight-power"),
		MaxWeight:      cfg.Float64("neat|populator|max-weight"),
		BiasPower:      cfg.Float64("neat|populator|bias-power"),
		MaxBias:        cfg.Float64("neat|populator|max-bias"),
	}

	exp.Transcriber = Transcriber{
		CppnTranscriber:  neat.Transcriber{DisableSortCheck: cfg.Bool("hyperneat|transcriber|cppn-transcriber|disable-sort-check")},
		CppnTranslator:   forward.Translator{DisableSortCheck: cfg.Bool("forward|translator|disable-sort-check")},
		Inspector:        LinkExpressionOutput{},
		WeightPower:      cfg.Float64("hyperneat|transcriber|weight-power"),
		BiasPower:        cfg.Float64("hyperneat|transcriber|bias-power"),
		DisableSortCheck: cfg.Bool("hyperneat|transcriber|disable-sort-check"),
	}

	// Add additional mutator for activations
	am := mutator.Activation{
		ReplaceActivationProbability: cfg.Float64("neat|mutator|activation|replace-activation-probability"),
		Activations:                  cfg.Activations("neat|mutator|activation|mutate-activations"),
	}
	if am.ReplaceActivationProbability > 0.0 {
		exp.Mutators = append(exp.Mutators, am)
	}
	return
}
