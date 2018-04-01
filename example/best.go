package example

import (
	"log"

	"github.com/klokare/evo"
)

// ShowBest is an EVO listener which will output a summary of the best genome in the population to the log
func ShowBest(pop evo.Population) error {

	// Copy the genomes so we can sort them without affecting other listeners
	genomes := make([]evo.Genome, len(pop.Genomes))
	copy(genomes, pop.Genomes)

	// Sort so the best genome is at the end
	evo.SortBy(genomes, evo.BySolved, evo.ByFitness, evo.ByComplexity, evo.ByAge)

	// Output the best
	best := genomes[len(genomes)-1]
	log.Printf("generation %d, id %d, species %d, fitness %f, solved %t, complexity %d\n", pop.Generation, best.ID, best.SpeciesID, best.Fitness, best.Solved, best.Complexity())
	return nil
}
