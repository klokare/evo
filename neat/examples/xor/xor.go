package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/klokare/config/json"
	"github.com/klokare/evo"
	"github.com/klokare/evo/example"
	"github.com/klokare/evo/example/xor"
	"github.com/klokare/evo/neat"
)

var runs = flag.Int("runs", 1, "number of runs to execute")
var compare = flag.Bool("compare", false, "run in comparison mode")

func main() {

	// Parse the command line flags
	flag.Parse()

	// Load the configuration
	var err error
	var cfg evo.Configurer
	if cfg, err = json.NewFromFile("xor.json"); err != nil {
		log.Fatal(err)
	}

	// Describe the experiment
	options := neat.WithOptions(cfg)               // begin with NEAT experiment
	options = append(options, xor.WithEvaluator()) // add the XOR evaluator
	if !*compare {
		options = append(options, example.WithProgress(evo.ByFitness)) // display progress
	}

	// Run the experiment
	pops, errs := evo.Batch(*runs, 100, options...)
	for i, err := range errs {
		if err != nil {
			log.Println("error in run", i, " was", err)
		}
	}

	// Note the number of failures and, for successes, number of generations, mean number of nodes
	// and conns of best
	var failed, nodes, conns, gens int
	for _, pop := range pops {
		evo.SortBy(pop.Genomes, evo.BySolved, evo.ByFitness, evo.ByComplexity, evo.ByAge)
		b := pop.Genomes[len(pop.Genomes)-1]
		if b.Solved {
			nodes += len(b.Encoded.Nodes)
			c2 := 0
			for _, c := range b.Encoded.Conns {
				if c.Enabled {
					c2++
				}
			}
			conns += c2
			gens += pop.Generation
		} else {
			failed++
		}
	}
	if !*compare {
		fmt.Printf("mean generations: %.2f, nodes: %.2f, conns: %.2f\n",
			float64(gens)/float64(len(pops)-failed),
			float64(nodes)/float64(len(pops)-failed),
			float64(conns)/float64(len(pops)-failed),
		)
		fmt.Println(failed, "failures out of", len(pops), "runs")
	} else {
		// fmt.Println(failed, "failures out of", len(pops), "runs")
		fmt.Println(len(pops), failed, float64(nodes)/float64(len(pops)-failed), float64(conns)/float64(len(pops)-failed), 100*float64(gens)/float64(len(pops)-failed))
	}
}
