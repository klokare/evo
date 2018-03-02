package neat

import (
	"github.com/klokare/evo"
	"github.com/klokare/evo/neat/mutator"
	"github.com/klokare/evo/network/forward"
	"github.com/klokare/evo/searcher/serial"
)

// WithOptions returns a set of options that fully configure an experiment using NEAT and default
// helpers.
func WithOptions(cfg evo.Configurer) []evo.Option {
	return []evo.Option{
		evo.WithConfiguration(cfg),
		evo.WithCompare(evo.ByFitness),
		serial.WithSearcher(),
		forward.WithTranslator(),
		WithCrosser(cfg),
		WithSeeder(cfg),
		WithSelector(cfg),
		WithSpeciator(cfg),
		WithTranscriber(cfg),
		mutator.WithComplexify(cfg),
		// mutator.WithSimplify(cfg),
		// mutator.WithPhased(cfg),
		mutator.WithPruning(cfg),
		mutator.WithBias(cfg),
		mutator.WithWeight(cfg),
	}
}
