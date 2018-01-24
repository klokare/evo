package evo

// A Species represents the current state of a group of like genomes
type Species struct {
	ID       int64   // The species's unique identifer
	Decay    float64 // The current decay amount applied to the species when calculating offspring or checking for stagnation
	Champion int64   // ID of best genome, according to experiment's Comparer, within the species
	Example  Genome  // The genome against whom other genomes are compared to see if they belong to this species
}
