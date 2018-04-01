package evo

import (
	"sort"
)

// SortBy orders the genome by the comparison functions.
func SortBy(genomes []Genome, comparisons ...Comparison) {
	sort.Slice(genomes, func(i, j int) bool {
		for _, c := range comparisons {
			x := c.Compare(genomes[i], genomes[j])
			if x < 0 {
				return true
			} else if x > 0 {
				return false
			}
		}
		return false
	})
}

// Comparison two genomes for relative order
type Comparison byte

// Known comparison types
const (
	ByFitness Comparison = iota + 1
	ByNovelty
	ByComplexity
	ByAge
	BySolved
	BySpecies
)

func (c Comparison) String() string {
	switch c {
	case ByFitness:
		return "fitness"
	case ByNovelty:
		return "novelty"
	case ByComplexity:
		return "complexity"
	case ByAge:
		return "age"
	case BySolved:
		return "solved"
	case BySpecies:
		return "species"
	default:
		return "unknown comparison"
	}
}

// Compare two genomes using the appropriate method
func (c Comparison) Compare(a, b Genome) int8 {
	switch c {

	case ByNovelty:
		switch {
		case a.Novelty < b.Novelty:
			return -1
		case b.Novelty < a.Novelty:
			return 1
		default:
			return 0
		}

	case ByComplexity:
		switch {
		case a.Complexity() > b.Complexity():
			return -1
		case a.Complexity() < b.Complexity():
			return 1
		default:
			return 0
		}

	case ByAge:
		switch {
		case a.ID > b.ID:
			return -1
		case a.ID < b.ID:
			return 1
		default:
			return 0
		}

	case BySolved:
		switch {
		case a.Solved && !b.Solved:
			return 1
		case b.Solved && !a.Solved:
			return -1
		default:
			return 0
		}

	case BySpecies:
		switch {
		case a.SpeciesID < b.SpeciesID:
			return -1
		case b.SpeciesID < a.SpeciesID:
			return 1
		default:
			return 0
		}

	default: // by fitness
		switch {
		case a.Fitness < b.Fitness:
			return -1
		case b.Fitness < a.Fitness:
			return 1
		default:
			return 0
		}
	}
}

// Compares provides map of compare functions by name
var (
	Comparisons = map[string]Comparison{
		"fitness":    ByFitness,
		"age":        ByAge,
		"solved":     BySolved,
		"novelty":    ByNovelty,
		"complexity": ByComplexity,
		"species":    BySpecies,
	}
)
