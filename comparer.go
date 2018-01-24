package evo

import "sort"

// SortBy orders the genome by the comparison functions.
func SortBy(genomes []Genome, comparisons ...Compare) {
	sort.Slice(genomes, func(i, j int) bool {
		for _, compare := range comparisons {
			x := compare(genomes[i], genomes[j])
			if x < 0 {
				return true
			} else if x > 0 {
				return false
			}
		}
		return false
	})
}

// Compare two genomes for relative order
type Compare func(a, b Genome) int8

// WithCompare sets the experiment's comparison function to the specified Compare
func WithCompare(fn Compare) Option {
	return func(e *Experiment) error {
		e.Compare = fn
		return nil
	}
}

// ByFitness returns the relative order of two genomes by fitness
func ByFitness(a, b Genome) int8 {
	switch {
	case a.Fitness < b.Fitness:
		return -1
	case b.Fitness < a.Fitness:
		return 1
	default:
		return 0
	}
}

// ByNovelty returns the relative order of two genomes by novelty
func ByNovelty(a, b Genome) int8 {
	switch {
	case a.Novelty < b.Novelty:
		return -1
	case b.Novelty < a.Novelty:
		return 1
	default:
		return 0
	}
}

// BySolved returns the relative order of two genomes by solution state
func BySolved(a, b Genome) int8 {
	switch {
	case a.Solved && !b.Solved:
		return 1
	case b.Solved && !a.Solved:
		return -1
	default:
		return 0
	}
}

// ByComplexity retuns the relative order of two genomes by their complexity. Note: a lower
// complexity is better.
func ByComplexity(a, b Genome) int8 {
	switch {
	case a.Complexity() > b.Complexity():
		return -1
	case a.Complexity() < b.Complexity():
		return 1
	default:
		return 0
	}
}

// ByAge retuns the relative order of two genomes by their age, using ID number as a proxy
func ByAge(a, b Genome) int8 {
	switch {
	case a.ID > b.ID:
		return -1
	case a.ID < b.ID:
		return 1
	default:
		return 0
	}
}
