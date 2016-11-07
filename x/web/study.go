package web

import (
	"sort"
	"time"

	"github.com/klokare/evo"
)

// A Study is a collection of experiments. Under this heading, variations can be compared.
type Study struct {
	ID          int
	Description string
	Experiments
}

// Best returns the best genome of the experiments
func (s Study) Best() evo.Genome {
	if len(s.Experiments) == 0 {
		return evo.Genome{}
	}
	var gs evo.Genomes = make([]evo.Genome, len(s.Experiments))
	for i, e := range s.Experiments {
		gs[i] = e.Best()
	}
	sort.Sort(sort.Reverse(gs))
	return gs[0]
}

// Last updated returns the time of the most recent update
func (s Study) LastUpdated() time.Time {
	var u time.Time
	for _, e := range s.Experiments {
		ut := e.LastUpdated()
		if u.Before(ut) {
			u = ut
		}
	}
	return u
}

func (s Study) Solved() bool {
	for _, e := range s.Experiments {
		if e.Solved() {
			return true
		}
	}
	return false
}

// Fitness returns the fitnesses of the best genome of each experiment
func (s Study) Fitness() Floats {
	var x Floats = make([]float64, len(s.Experiments))
	for i, e := range s.Experiments {
		x[i] = e.Best().Fitness
	}
	return x
}

// Novelty returns the novelty values of the best genome of each experiment
func (s Study) Novelty() Floats {
	var x Floats = make([]float64, len(s.Experiments))
	for i, e := range s.Experiments {
		x[i] = e.Best().Novelty
	}
	return x
}

// Complexity returns the complexity of the best genome of each experiment
func (s Study) Complexity() Floats {
	var x Floats = make([]float64, len(s.Experiments))
	for i, e := range s.Experiments {
		x[i] = float64(e.Best().Complexity())
	}
	return x
}

// Diversity returns the average number of species in each experiment
func (s Study) Diversity() Floats {
	var x Floats = make([]float64, len(s.Experiments))
	for i, e := range s.Experiments {
		x[i] = e.Diversity().Mean()
	}
	return x
}

// Iterations returns the average number of iterations in each experiment
func (s Study) Iterations() Floats {
	var x Floats = make([]float64, len(s.Experiments))
	for i, e := range s.Experiments {
		x[i] = e.Iterations().Mean()
	}
	return x
}

type Studies []Study

func (ss Studies) Len() int      { return len(ss) }
func (ss Studies) Swap(i, j int) { ss[i], ss[j] = ss[j], ss[i] }
func (ss Studies) Less(i, j int) bool {
	ui := ss[i].LastUpdated()
	uj := ss[j].LastUpdated()
	if ui.Equal(uj) {
		return ss[i].ID < ss[j].ID
	}
	return ui.Before(uj)
}

// StudyClient manages studies in a persistent store
type StudyClient interface {
	AddStudy(desc string) (int, error)
	SetStudy(sid int, desc string) error
	DelStudy(sid int) error
	GetStudy(sid int) (Study, error)
	GetStudies() ([]Study, error)
}
