package main

import (
	"context"
	"flag"
	"log"

	"github.com/klokare/evo"
	"github.com/klokare/evo/config"
	"github.com/klokare/evo/config/source"
	"github.com/klokare/evo/efficacy"
	"github.com/klokare/evo/example"
	"github.com/klokare/evo/example/xor"
	"github.com/klokare/evo/neat"
)

// Define flags to override configuration file settings
var (
	_ = flag.String("neat-hidden-activation", "", "override hidden activation property")
	_ = flag.String("neat-output-activation", "", "override output activation property")
)

func main() {

	// Parse the command-line flags
	var (
		runs  = flag.Int("runs", 1, "number of experiments to run")
		iter  = flag.Int("iterations", 100, "number of iterations for experiment")
		cpath = flag.String("config", "xor.json", "path to the configuration file")
		epath = flag.String("efficacy", "xor-samples.txt", "path for efficacy sample file")
	)
	flag.Parse()

	// Load the configuration
	src, err := source.NewJSONFromFile(*cpath)
	if err != nil {
		log.Fatalf("%+v\n", err)
	}
	cfg := config.Configurer{Source: source.Multi([]config.Source{
		source.Flag{},        // Check flags  first
		source.Environment{}, // Then check environment variables
		src,                  // Lastly, consult the configuration file
	})}

	// Create a sample file if performing multiple runs
	var s *efficacy.Sampler
	if *runs > 1 {
		if s, err = efficacy.NewSampler(*epath); err != nil {
			log.Fatalf("%+v\n", err)
		}
		defer s.Close()
	}

	// Iterate the runs
	for r := 0; r < *runs; r++ {

		// Create the experiment
		exp := neat.NewExperiment(cfg)

		// Add additional subscriptions
		if s == nil {
			exp.AddSubscription(evo.Subscription{Event: evo.Completed, Callback: example.ShowBest}) // Show summary upon completion
		} else {
			c0, c1 := s.Callbacks(r)
			exp.AddSubscription(evo.Subscription{Event: evo.Started, Callback: c0})   // Begin the efficacy sample
			exp.AddSubscription(evo.Subscription{Event: evo.Completed, Callback: c1}) // End the efficacy sample
		}

		// Run the experiment for a set number of iterations
		ctx, fn, cb := evo.WithIterations(context.Background(), *iter)
		defer fn() // ensure the context cancels
		exp.AddSubscription(evo.Subscription{Event: evo.Evaluated, Callback: cb})

		// Stop the experiment if there is a solution
		ctx, fn, cb = evo.WithSolution(ctx)
		defer fn() // ensure the context cancels
		exp.AddSubscription(evo.Subscription{Event: evo.Evaluated, Callback: cb})

		// Execute the experiment
		if _, err = evo.Run(ctx, exp, xor.Evaluator{}); err != nil {
			log.Fatalf("%+v\n", err)
		}
	}
}
