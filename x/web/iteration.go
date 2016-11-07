package web

import (
	"time"

	"github.com/klokare/evo"
)

// Iteration represents one iteration in a trial
type Iteration struct {
	ID          int
	Description string
	Updated     time.Time
	Best        evo.Genome
	Solved      bool
	Fitness     Floats
	Novelty     Floats
	Complexity  Floats
	Diversity   Floats
}

// Iterations is a sortable collection of iterations
type Iterations []Iteration

func (is Iterations) Len() int      { return len(is) }
func (is Iterations) Swap(i, j int) { is[i], is[j] = is[j], is[i] }
func (is Iterations) Less(i, j int) bool {
	if is[i].Updated.Equal(is[j].Updated) {
		return is[i].ID < is[j].ID
	}
	return is[i].Updated.Before(is[j].Updated)
}

// IterationClient manages iterations in the persistent store
type IterationClient interface {
	AddIteration(tid int, p evo.Population) (int, error)
	DelIteration(iid int) error // Do we need? Individual iterations should not be deleted. Instead delete the whole trial
	GetIteration(iid int) (Iteration, error)
	GetIterations(tid int) ([]Iteration, error)
	GetPopulation(iid int) (evo.Population, error)
	SetIteration(iid int, desc string) error
}
