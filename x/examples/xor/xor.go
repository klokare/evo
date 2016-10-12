package main

import (
	"log"
	"math"

	"github.com/klokare/evo"
	"github.com/klokare/evo/x/configurer"
	start "github.com/klokare/evo/x/neat"
)

func main() {

	// Load the configuration
	var err error
	var b []byte
	if b, err = configurer.LoadFromFile("xor-config.json"); err != nil {
		log.Fatal(err)
	}
	c := &configurer.JSON{Source: b}

	// Create a new experiment
	var exp *evo.Experiment
	if exp, err = start.NewExperiment(c, &XOR{}); err != nil {
		log.Fatal(err)
	}

	// Create the initial population
	pop := start.Seed(250, 2, 1)

	// Run the experiment over several iterations and sum iterations
	x := make([]int, 10)
	g := &gen{}
	ws := exp.Watcher.(evo.Watchers)
	ws = append(ws, g)
	exp.Watcher = ws
	for i := 0; i < len(x); i++ {
		g.generation = 0
		if err = evo.Run(exp, pop, 150); err != nil {
			log.Fatal(err)
		}
		x[i] = g.generation
	}
	log.Println("trials", x)
}

type gen struct {
	generation int
}

func (g *gen) Watch(p evo.Population) error {
	g.generation = p.Generation
	return nil
}

var (
	inputs = [][]float64{{0, 0}, {0, 1}, {1, 0}, {1, 1}}

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
	f := math.Max(0.000001, 1.0-sse)
	s := stop || f >= 0.95
	return evo.Result{ID: p.ID, Fitness: f, Solved: s}, nil
}
