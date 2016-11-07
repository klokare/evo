package web

import "github.com/klokare/evo"

// Stats provides data for the object
type Stats interface {
	Best() evo.Genome
	Solved() bool
	Fitness() Floats
	Novelty() Floats
	Complexity() Floats
	Diversity() Floats
	Iterations() Floats
}
