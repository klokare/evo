package main

import (
	"encoding/json"
	"flag"
	"log"
	"math"
	"os"

	"github.com/klokare/evo"
	"github.com/klokare/evo/x/config"
	"github.com/klokare/evo/x/neat"
)

var (
	path = flag.String("config", "xor-config.json", "path to the configuration file")
)

func main() {
	flag.Parse()

	// Load the configuration
	var err error
	var f *os.File
	if f, err = os.Open(*path); err != nil {
		log.Fatal(err)
	}

	c := &config.Configurer{Settings: &config.Settings{}}
	if err := json.NewDecoder(f).Decode(&c.Settings); err != nil {
		log.Fatal(err)
	}

	// Create the experiment and run the trials
	var e *evo.Experiment
	if e, err = neat.NewExperiment(c, func() (evo.Evaluator, error) { return &XOR{}, nil }); err != nil {
		log.Fatal(err)
	}
	if err := evo.Run(e); err != nil {
		log.Fatal(err)
	}

}

var (
	inputs   = [][]float64{{0, 0}, {0, 1}, {1, 0}, {1, 1}}
	expected = []float64{0, 1, 1, 0}
)

// XOR is the evaluator for this experiment
type XOR struct{}

// Evaluate computes the error for the XOR problem with the phenome
func (e *XOR) Evaluate(p evo.Phenome) (evo.Result, error) {
	sse := 0.0
	stop := true
	for i, inputs := range inputs {
		outputs := p.Activate(inputs)
		sse += (outputs[0] - expected[i]) * (outputs[0] - expected[i])
		if expected[i] == 0 {
			stop = stop && outputs[0] < 0.5
		} else {
			stop = stop && outputs[0] > 0.5
		}
	}
	f := math.Max(1e-10, 1.0-sse)
	s := stop || f >= 0.95
	return evo.Result{ID: p.ID, Fitness: f, Solved: s}, nil
}
