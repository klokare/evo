package example

import (
	"context"
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/klokare/evo"
)

// Progress is a helper that reports on the prorgess of the experiment
type Progress struct {
	OrderBy []evo.Compare
	lastGen int
}

// Show is a callback function which can be subscribed to the experiment to show a progress
// summary (when final = false) or detail (when final = true).
func (d *Progress) Show(ctx context.Context, final bool, pop evo.Population) error {

	// Ignore if there is not a change in generation
	if !final && pop.Generation == d.lastGen {
		return nil
	}
	d.lastGen = pop.Generation

	// Sort the genomes
	evo.SortBy(pop.Genomes, d.OrderBy...)

	// Identify the best
	b := pop.Genomes[len(pop.Genomes)-1]
	sort.Slice(b.Encoded.Nodes, func(i, j int) bool { return b.Encoded.Nodes[i].Compare(b.Encoded.Nodes[j]) < 0 })
	sort.Slice(b.Encoded.Conns, func(i, j int) bool { return b.Encoded.Conns[i].Compare(b.Encoded.Conns[j]) < 0 })

	// This is just a summary
	if !final {
		fmt.Printf("Gen. %d has %d species, best is %d from species %d, with %.4f fitness, %.4f novelty, %d complexity\n",
			pop.Generation, len(pop.Species), b.ID, b.SpeciesID, b.Fitness, b.Novelty, b.Complexity())
		return nil
	}

	// Show detail on final call
	sort.Slice(pop.Species, func(i, j int) bool { return pop.Species[i].ID < pop.Species[j].ID })
	fmt.Println("==================================================================")
	fmt.Printf("Experiment ended on generation %d, solved = %v\n", pop.Generation, b.Solved)
	fmt.Println("Species:")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.TabIndent)
	fmt.Fprintf(w, "ID\tDecay\n")
	for _, s := range pop.Species {
		fmt.Fprintf(w, "%d\t%f\n", s.ID, s.Decay)
	}
	w.Flush()
	fmt.Printf("Best Genome %d in gen. %d and species %d, fitness = %f, novelty = %f, solved = %v\n",
		b.ID, pop.Generation, b.SpeciesID, b.Fitness, b.Novelty, b.Solved)
	fmt.Println("Nodes:")
	w = tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.TabIndent)
	fmt.Fprintf(w, "Position\tNeuron\tActivation\tBias\n")
	for _, n := range b.Encoded.Nodes {
		fmt.Fprintf(w, "%v\t%s\t%s\t%f\n", n.Position, n.Neuron, n.Activation, n.Bias)
	}
	w.Flush()
	fmt.Println("Connections:")
	w = tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.TabIndent)
	fmt.Fprintf(w, "Source\tTarget\tEnabled\tWeight\n")
	for _, c := range b.Encoded.Conns {
		fmt.Fprintf(w, "%v\t%v\t%v\t%f\n", c.Source, c.Target, c.Enabled, c.Weight)
	}
	w.Flush()
	return nil
}

// WithProgress subscribes a display helper to the experiment. The ordering list of compare functions
// is used to identify the best genome in the species.
func WithProgress(compare evo.Compare) evo.Option {
	return func(e *evo.Experiment) error {
		// Create the progress helper
		p := &Progress{OrderBy: []evo.Compare{evo.BySolved, compare, evo.ByComplexity, evo.ByAge}}

		// Add to the experiment
		e.Subscribe(evo.Evaluated, p.Show)
		return nil
	}
}
