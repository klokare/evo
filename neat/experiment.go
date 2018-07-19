package neat

import (
	"github.com/klokare/evo"
	"github.com/klokare/evo/config"
	"github.com/klokare/evo/neat/mutator"
	"github.com/klokare/evo/network/forward"
	"github.com/klokare/evo/searcher/parallel"
)

// Ensure the experiment struct implements the experiment interface
var (
	_ evo.Experiment = &Experiment{}
)

// Experiment implements an EVO experiment with the NEAT helpers.
type Experiment struct {
	Crosser
	Populator
	Selector
	Speciator
	Transcriber
	forward.Translator
	evo.Searcher
	evo.Mutators
	subscriptions []evo.Subscription
}

// NewExperiment creates a new NEAT experiment using the configuration. Configurations employ the
// maximum namespace so user can be as specific or lax (depending on depth of namespace used) as
// desired.
func NewExperiment(cfg config.Configurer) (exp *Experiment) {

	// Create the experiment using the NEAT and other default helpers
	exp = &Experiment{

		// Set the crosser helper
		Crosser: Crosser{
			EnableProbability:       cfg.Float64("neat|crosser|enable-probability"),
			DisableEqualParentCheck: cfg.Bool("neat|crosser|disable-equal-parent-check"),
			Comparison:              cfg.Comparison("neat|crosser|comparison"),
			DisableSortCheck:        cfg.Bool("neat|crosser|disable-sort-check"),
		},

		// Set the populator helper including the seeder helper
		Populator: Populator{
			Seeder: Seeder{
				NumInputs:        cfg.Int("neat|seeder|num-inputs"),
				NumOutputs:       cfg.Int("neat|seeder|num-outputs"),
				NumTraits:        cfg.Int("neat|seeder|num-traits"),
				OutputActivation: cfg.Activation("neat|seeder|output-activation"),
				DisconnectRate:   cfg.Float64("neat|seeder|disconnect-rate"),
			},
			PopulationSize: cfg.Int("neat|populator|population-size"),
			WeightPower:    cfg.Float64("neat|populator|weight-power"),
			MaxWeight:      cfg.Float64("neat|populator|max-weight"),
			BiasPower:      cfg.Float64("neat|populator|bias-power"),
			MaxBias:        cfg.Float64("neat|populator|max-bias"),
		},

		// Set the selector helper
		Selector: Selector{
			PopulationSize:              cfg.Int("neat|selector|population-size"),
			MutateOnlyProbability:       cfg.Float64("neat|selector|mutate-only-probability"),
			InterspeciesMateProbability: cfg.Float64("neat|selector|interspecies-mate-probability"),
			Elitism:                     cfg.Float64("neat|selector|elitism"), // Set to zero for novelty search
			SurvivalRate:                cfg.Float64("neat|selector|survival-rate"),
			Comparison:                  cfg.Comparison("neat|selector|comparison"), // can specify multiple functions separated by comma
			DecayRate:                   cfg.Float64("neat|updater|species-decay-rate"),
		},

		// Set the speciator helper using the compatibility distance helper
		Speciator: Speciator{
			Distancer: Compatibility{
				NodesCoefficient:      cfg.Float64("neat|distancer|nodes-coefficient"),
				ConnsCoefficient:      cfg.Float64("neat|distancer|conns-coefficient"),
				WeightCoefficient:     cfg.Float64("neat|distancer|weight-coefficient"),
				BiasCoefficient:       cfg.Float64("neat|distancer|bias-coefficient"),
				ActivationCoefficient: cfg.Float64("neat|distancer|activation-coefficient"),
				DisableSortCheck:      cfg.Bool("neat|distancer|disable-sort-check"),
			},
			CompatibilityThreshold: cfg.Float64("neat|speciator|compatibility-threshold"),
			CompatibilityModifier:  cfg.Float64("neat|speciator|compatibility-modifier"),
			TargetSpecies:          cfg.Int("neat|speciator|target-species"),
		},

		// Set the transcriber helper
		Transcriber: Transcriber{
			DisableSortCheck: cfg.Bool("neat|transcriber|disable-sort-check"),
		},

		// Create an empty slice for the mutators. Those are added below.
		Mutators: make([]evo.Mutator, 0, 5),

		// Set the translator to the default forward network
		Translator: forward.Translator{DisableSortCheck: cfg.Bool("forward|translator|disable-sort-check")},

		// Initialise the subscriptions slice
		subscriptions: make([]evo.Subscription, 0, 5),

		// Set the default Searcher
		Searcher: parallel.Searcher{},
	}

	// Add the mutators. Only those with a chance of being activated will be added.
	cm := mutator.Complexify{
		AddNodeProbability: cfg.Float64("neat|mutator|complexify|add-node-probability"),
		AddConnProbability: cfg.Float64("neat|mutator|complexify|add-conn-probability"),
		WeightPower:        cfg.Float64("neat|mutator|complexify|weight-power"),
		MaxWeight:          cfg.Float64("neat|mutator|complexify|max-weight"),
		BiasPower:          cfg.Float64("neat|mutator|complexify|bias-power"),
		MaxBias:            cfg.Float64("neat|mutator|complexify|max-bias"),
		HiddenActivation:   cfg.Activation("neat|mutator|complexify|hidden-activation"),
		DisableSortCheck:   cfg.Bool("neat|mutator|complexify|disable-sort-check"),
	}

	if cm.AddNodeProbability > 0.0 || cm.AddConnProbability > 0.0 {
		exp.Mutators = append(exp.Mutators, cm)
	}

	wm := mutator.Weight{
		MutateWeightProbability:  cfg.Float64("neat|mutator|weight|mutate-weight-probability"),
		ReplaceWeightProbability: cfg.Float64("neat|mutator|weight|replace-weight-probability"),
		WeightPower:              cfg.Float64("neat|mutator|weight|weight-power"),
		MaxWeight:                cfg.Float64("neat|mutator|weight|max-weight"),
	}
	if wm.MutateWeightProbability > 0.0 {
		exp.Mutators = append(exp.Mutators, wm)
	}

	bm := mutator.Bias{
		MutateBiasProbability:  cfg.Float64("neat|mutator|bias|mutate-bias-probability"),
		ReplaceBiasProbability: cfg.Float64("neat|mutator|bias|replace-bias-probability"),
		BiasPower:              cfg.Float64("neat|mutator|bias|bias-power"),
		MaxBias:                cfg.Float64("neat|mutator|bias|max-bias"),
	}
	if bm.MutateBiasProbability > 0.0 {
		exp.Mutators = append(exp.Mutators, bm)
	}

	tm := mutator.Trait{
		MutateTraitProbability: cfg.Float64("neat|mutator|trait|mutate-trait-probability"),
	}
	if tm.MutateTraitProbability > 0.0 {
		exp.Mutators = append(exp.Mutators, tm)
	}

	// Add subscriptions
	// TODO: add subscription for phased mutator when that helper is complete.

	// Return the new experiment
	return
}

// Subscriptions returns the subscriptions registered with this experiment.
func (e *Experiment) Subscriptions() []evo.Subscription { return e.subscriptions }

// AddSubscription adds a new subscription to the experiment
func (e *Experiment) AddSubscription(s evo.Subscription) {
	if e.subscriptions == nil {
		e.subscriptions = make([]evo.Subscription, 0, 5)
	}
	e.subscriptions = append(e.subscriptions, s)
}
