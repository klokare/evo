package evo

// A Population is the collection of genomes and species at a particular point in time
type Population struct {
	Generation int
	Genomes    Genomes
	Species    []Species
}

// Species is a collection of similar genomes
type Species struct {
	ID         int
	Stagnation int
	Fitness    float64
	Example    Substrate
}
