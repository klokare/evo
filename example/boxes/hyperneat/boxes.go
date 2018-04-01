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
	"github.com/klokare/evo/example/boxes"
	"github.com/klokare/evo/hyperneat"
)

// Define flags to override configuration file settings
var (
	_ = flag.String("neat-hidden-activation", "", "override hidden activation property")
	_ = flag.String("neat-output-activation", "", "override output activation property")
	_ = flag.Int("neat-population-size", 100, "override population size property")
)

func main() {
	// defer profile.Start(profile.MemProfile).Stop()

	// Parse the command-line flags
	var (
		runs   = flag.Int("runs", 1, "number of experiments to run")
		iter   = flag.Int("iterations", 300, "number of iterations for experiment")
		cpath  = flag.String("config", "boxes.json", "path to the configuration file")
		epath  = flag.String("efficacy", "boxes-samples.txt", "path for efficacy sample file")
		ipath  = flag.String("image", "boxes-champ.png", "path for output image on single run")
		hidden = flag.Bool("hidden", false, "use hidden nodes in boxes template")
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
		exp := hyperneat.NewExperiment(cfg)

		// Create the evaluator
		eval := boxes.NewEvaluator(cfg.Int("boxes|resolution"))

		// Initialise the template
		exp.Transcriber.SetTemplate(boxes.Template(cfg.Int("boxes|resolution"), *hidden))

		// Add additional subscriptions
		if s == nil {
			exp.AddSubscription(evo.Subscription{Event: evo.Completed, Callback: example.ShowBest})                         // Show summary upon completion
			exp.AddSubscription(evo.Subscription{Event: evo.Completed, Callback: writeImage(*ipath, exp.Translator, eval)}) // Write output image
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
		if _, err = evo.Run(ctx, exp, eval); err != nil {
			log.Fatalf("%+v\n", err)
		}
	}
}

func writeImage(filename string, trans evo.Translator, eval *boxes.Evaluator) evo.Callback {
	return func(pop evo.Population) (err error) {

		// Determine the best genome in the population
		genomes := make([]evo.Genome, len(pop.Genomes))
		copy(genomes, pop.Genomes) // copy so as not to affect other callbacks
		evo.SortBy(genomes, evo.BySolved, evo.ByFitness, evo.ByComplexity, evo.ByAge)
		best := genomes[len(genomes)-1]

		// Create a phenome
		p := evo.Phenome{ID: best.ID}
		if p.Network, err = trans.Translate(best.Decoded); err != nil {
			return
		}

		// Tell the evaluator to write the solution and then evaluate
		eval.SetOutput(filename)
		_, err = eval.Evaluate(p)
		return
	}
}
