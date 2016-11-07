package web

import (
	"sort"
	"time"

	"github.com/klokare/evo"
)

// A Trial is a single run of an experiment. The summary of each generation is stored here.
type Trial struct {
	ID          int
	Description string
	Iters       Iterations
}

func (t Trial) LastUpdated() time.Time {
	var u time.Time
	for _, i := range t.Iters {
		if u.Before(i.Updated) {
			u = i.Updated
		}
	}
	return u
}

func (t Trial) Best() evo.Genome {
	var gs evo.Genomes = make([]evo.Genome, len(t.Iters))
	for i, it := range t.Iters {
		gs[i] = it.Best
	}
	sort.Sort(sort.Reverse(gs))
	return gs[0]
}

func (t Trial) Solved() bool {
	for _, it := range t.Iters {
		if it.Solved {
			return true
		}
	}
	return false
}

// Fitness returns the fitnesses of the best genome of each iteration
func (t Trial) Fitness() Floats {
	var x Floats = make([]float64, len(t.Iters))
	for i, it := range t.Iters {
		x[i] = it.Best.Fitness
	}
	return x
}

// Novelty returns the novelty values of the best genome of each iteration
func (t Trial) Novelty() Floats {
	var x Floats = make([]float64, len(t.Iters))
	for i, it := range t.Iters {
		x[i] = it.Best.Novelty
	}
	return x
}

// Complexity returns the complexity of the population
func (t Trial) Complexity() Floats {
	var x Floats = make([]float64, len(t.Iters))
	for i, it := range t.Iters {
		x[i] = float64(it.Best.Complexity())
	}
	return x
}

// Diversity returns the number of species in each iteration
func (t Trial) Diversity() Floats {
	var x Floats = make([]float64, len(t.Iters))
	for i, it := range t.Iters {
		x[i] = it.Diversity[0]
	}
	return x
}

// Iterations returns the number of iterations in this trial
func (t Trial) Iterations() Floats {
	return []float64{float64(len(t.Iters))}
}

type Trials []Trial

func (ts Trials) Len() int      { return len(ts) }
func (ts Trials) Swap(i, j int) { ts[i], ts[j] = ts[j], ts[i] }
func (ts Trials) Less(i, j int) bool {
	ui := ts[i].LastUpdated()
	uj := ts[j].LastUpdated()
	if ui.Equal(uj) {
		return ts[i].ID < ts[j].ID
	}
	return ui.Before(uj)
}

// TrialClient manages trials in the peristent store
type TrialClient interface {
	AddTrial(eid int, desc string) (int, error)
	DelTrial(tid int) error
	GetTrial(tid int) (Trial, error)
	GetTrials(eid int) ([]Trial, error)
}
