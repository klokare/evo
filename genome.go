package evo

// A Genome is a single solution
type Genome struct {
	ID        int       // Unique identifier for the genome
	SpeciesID int       // The ID of the species to which this genome belongs
	Fitness   float64   // The most recent fitness value
	Novelty   float64   // The most recent novelty value, if any
	Traits    []float64 // Traits encoded with this genome
	Solved    bool      // True if the genome represents a completed solution
	Encoded   Substrate // The encoded neural network
	Decoded   Substrate // The decoded neural network
}

// Complexity describes the genome by the number of encoded nodes and connections
func (g Genome) Complexity() int { return g.Encoded.Complexity() }

// Genomes is a sortable list of genomes.
type Genomes []Genome

func (gs Genomes) Len() int      { return len(gs) }
func (gs Genomes) Swap(i, j int) { gs[i], gs[j] = gs[j], gs[i] }
func (gs Genomes) Less(i, j int) bool {
	if gs[i].Fitness < gs[j].Fitness {
		return true // Lower fitness is "less"
	} else if gs[i].Fitness == gs[j].Fitness {
		ci := gs[i].Complexity()
		cj := gs[j].Complexity()
		if ci > cj {
			return true // Higher complexity is "less"
		} else if ci == cj {
			return gs[i].ID > gs[j].ID // Greater ID is "less"
		}
	}
	return false
}

// AvgFitness returns the average fitness of the genomes in the set
func (gs Genomes) AvgFitness() float64 {
	if len(gs) == 0 {
		return 0.0
	}
	sum := 0.0
	for _, g := range gs {
		sum += g.Fitness
	}
	return sum / float64(len(gs))
}

// MaxFitness returns the maximum fitness of the genomes in the set
func (gs Genomes) MaxFitness() float64 {
	max := 0.0
	for _, g := range gs {
		if max < g.Fitness {
			max = g.Fitness
		}
	}
	return max
}

// A Phenome is a working instance of the genome
type Phenome struct {
	ID     int
	Traits []float64
	Network
}

// Result of an evaluation
type Result struct {
	ID       int       // ID of the phenome evaluated
	Fitness  float64   // Fitness of phenome after evaluation
	Novelty  float64   // Novelty of phenome, if any, after evaluation
	Behavior []float64 // The behaviour of the phenome recorded during the evaluation
	Solved   bool      // Phenome is a completed solution
	Error    error     // Error, if any, that occurred during evaluation of the phenome
}
