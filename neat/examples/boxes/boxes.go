package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/klokare/config/json"
	"github.com/klokare/evo"
	"github.com/klokare/evo/example/boxes"
	"github.com/klokare/evo/neat"
	"github.com/klokare/evo/searcher/parallel"
)

var (
	config = flag.String("config", "boxes-default.json", "configuration file for experiment")
	runs   = flag.Int("runs", 1, "number of runs to execute")
	n      = flag.Int("n", 300, "number of iterations")
	name   = flag.String("name", "boxes-snapshot.txt", "output file name for snapshots")
)

func main() {

	// Parse the command line flags
	flag.Parse()
	//defer profile.Start().Stop()

	// Load the configuration
	var err error
	var cfg evo.Configurer
	if cfg, err = json.NewFromFile(*config); err != nil {
		log.Fatal(err)
	}

	// Create a new tracker file
	var t *tracker
	if t, err = newTracker(*name); err != nil {
		log.Fatal(err)
	}
	defer t.Close()

	// Describe the experiment
	options := neat.WithOptions(cfg) // begin with NEAT experiment
	options = append(options,
		parallel.WithSearcher(),  // Parallel searcher
		boxes.WithEvaluator(cfg), // add the OCR evaluator
		withTracker(t),
		//example.WithProgress(evo.ByFitness), // display progress
	)

	// Run the experiment
	pops, errs := evo.Batch(*runs, *n, options...)
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
		nodes += len(b.Encoded.Nodes)
		c2 := 0
		for _, c := range b.Encoded.Conns {
			if c.Enabled {
				c2++
			}
		}
		conns += c2
		gens += pop.Generation
		if !b.Solved {
			failed++
		}
	}
	// fmt.Printf("mean generations: %.2f, nodes: %.2f, conns: %.2f\n",
	// 	float64(gens)/float64(len(pops)-failed),
	// 	float64(nodes)/float64(len(pops)-failed),
	// 	float64(conns)/float64(len(pops)-failed),
	// )
	// fmt.Println(failed, "failures out of", len(pops), "runs")
	fmt.Println(len(pops), failed, float64(nodes)/float64(len(pops)), float64(conns)/float64(len(pops)), 100*float64(gens)/float64(len(pops)))

}

type record struct {
	Elapsed    float64
	Generation int
	Grouping   string
	Complexity float64
	Fitness    float64
}

type tracker struct {
	t0   time.Time
	last float64
	w    *bufio.Writer
	c    io.Closer
	ch   chan record
}

func (t *tracker) Close() error {
	close(t.ch)
	time.Sleep(5 * time.Second)
	t.w.Flush()
	if t.c != nil {
		return t.c.Close()
	}
	return nil
}

func newTracker(name string) (t *tracker, err error) {
	t = &tracker{t0: time.Now()}
	var f *os.File
	if f, err = os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
		return
	}
	t.c = f
	t.w = bufio.NewWriter(f)
	t.w.WriteString("elapsed\tgeneration\tgrouping\tcomplexity\tfitness\n")

	t.ch = make(chan record, *runs*2)
	go func(ch chan record) {
		for r := range ch {
			if _, err = fmt.Fprintf(t.w, "%f\t%d\t%s\t%f\t%f\n", r.Elapsed, r.Generation, r.Grouping, r.Complexity, r.Fitness); err != nil {
				log.Fatal(err)
			}
		}
	}(t.ch)
	return
}

func (t *tracker) Snapshot(final bool, pop evo.Population) (err error) {

	// Write every second
	d := time.Now().Sub(t.t0).Seconds()

	// Calculate the mean and best
	evo.SortBy(pop.Genomes, evo.BySolved, evo.ByFitness)
	c := 0.0
	f := 0.0
	for _, g := range pop.Genomes {
		c += float64(g.Complexity())
		f += g.Fitness
	}
	c /= float64(len(pop.Genomes))
	f /= float64(len(pop.Genomes))
	t.ch <- record{
		Elapsed:    d,
		Generation: pop.Generation,
		Grouping:   "mean",
		Complexity: c,
		Fitness:    f,
	}
	b := pop.Genomes[len(pop.Genomes)-1]
	t.ch <- record{
		Elapsed:    d,
		Generation: pop.Generation,
		Grouping:   "best",
		Complexity: float64(b.Complexity()),
		Fitness:    b.Fitness,
	}
	return
}

func withTracker(t *tracker) evo.Option {
	return func(e *evo.Experiment) error {
		e.Subscribe(evo.Evaluated, t.Snapshot)
		return nil
	}
}
