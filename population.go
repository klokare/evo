package evo

import "sort"

// A Population is the collection of genomes and species for a given generation.
type Population struct {
	Generation int      // The population's generation number
	Genomes    []Genome // The population's collection of genomes. The ordering of these is not guranateed.
}

// GroupBySpecies returns the genoems orgainised by their species. These are copies and do not
// point back to the original slice.
// TODO: is there a more efficient way to build this than going through a map?
func GroupBySpecies(genomes []Genome) [][]Genome {

	// Map genomes to species
	m := make(map[int][]Genome, 20)
	for _, g := range genomes {
		gs := m[g.Species]
		gs = append(gs, g)
		m[g.Species] = gs
	}

	// Transform into slice
	bs := make([][]Genome, 0, len(m))
	for _, gs := range m {
		bs = append(bs, gs)
	}
	sort.Slice(bs, func(i, j int) bool { return bs[i][0].Species < bs[j][0].Species })
	return bs
}
