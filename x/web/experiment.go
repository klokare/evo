package web

import (
	"sort"
	"time"

	"github.com/klokare/evo"
)

// An Experiment is a collection of Experiments. Because there is so much randomness in running
// experiments, multiple Experiments must be used and the outcomes studied statistically.
type Experiment struct {
	ID          int
	Description string
	Trials
}

func (e Experiment) LastUpdated() time.Time {
	var u time.Time
	for _, e := range e.Trials {
		ut := e.LastUpdated()
		if u.Before(ut) {
			u = ut
		}
	}
	return u
}

func (e Experiment) Best() evo.Genome {
	var gs evo.Genomes = make([]evo.Genome, len(e.Trials))
	for i, t := range e.Trials {
		gs[i] = t.Best()
	}
	sort.Sort(sort.Reverse(gs))
	return gs[0]
}

func (e Experiment) Solved() bool {
	for _, t := range e.Trials {
		if t.Solved() {
			return true
		}
	}
	return false
}

// Fitness returns the fitnesses of the best genome of each trial
func (e Experiment) Fitness() Floats {
	var x Floats = make([]float64, len(e.Trials))
	for i, t := range e.Trials {
		x[i] = t.Best().Fitness
	}
	return x
}

// Novelty returns the novelty values of the best genome of each trial
func (e Experiment) Novelty() Floats {
	var x Floats = make([]float64, len(e.Trials))
	for i, t := range e.Trials {
		x[i] = t.Best().Novelty
	}
	return x
}

// Complexity returns the complexity of the best genome of each trial
func (e Experiment) Complexity() Floats {
	var x Floats = make([]float64, len(e.Trials))
	for i, t := range e.Trials {
		x[i] = float64(t.Best().Complexity())
	}
	return x
}

// Diversity returns the average number of species in each trial
func (e Experiment) Diversity() Floats {
	var x Floats = make([]float64, len(e.Trials))
	for i, t := range e.Trials {
		x[i] = t.Diversity().Mean()
	}
	return x
}

// Trials returns the number of iterations in each trial
func (e Experiment) Iterations() Floats {
	var x Floats = make([]float64, len(e.Trials))
	for i, t := range e.Trials {
		x[i] = t.Iterations()[0]
	}
	return x
}

type Experiments []Experiment

func (es Experiments) Len() int      { return len(es) }
func (es Experiments) Swap(i, j int) { es[i], es[j] = es[j], es[i] }
func (es Experiments) Less(i, j int) bool {
	ui := es[i].LastUpdated()
	uj := es[j].LastUpdated()
	if ui.Equal(uj) {
		return es[i].ID < es[j].ID
	}
	return ui.Before(uj)
}

// ExperimentClient manages experiments in the peristent store
type ExperimentClient interface {
	AddExperiment(sid int, desc string) (int, error)
	SetExperiment(sid int, eid int, desc string) error
	DelExperiment(eid int) error
	GetExperiment(eid int) (Experiment, error)
	GetExperiments(sid int) ([]Experiment, error)
}
