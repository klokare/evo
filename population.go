package evo

// A Population is the collection of genomes and species for a given generation.
type Population struct {
	Generation int       // The population's generation number
	Species    []Species // The population's collection of species. The ordering of these is not guranateed.
	Genomes    []Genome  // The population's collection of genomes. The ordering of these is not guranateed.
}
