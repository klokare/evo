package selector

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sort"

	"github.com/klokare/evo"
	"github.com/klokare/random"
)

// A Generational selector creates a new generation every iteration
type Generational struct {
	SurvivalRate                float64 `evo:"survival-rate"`
	MaxStagnation               int     `evo:"max-stagnation"`
	MutateOnlyProbability       float64 `evo:"mutate-only-probability"`
	InterspeciesMateProbability float64 `evo:"interspecies-mate-probability"`

	maxFitness float64
	stagnation int
}

func (h Generational) String() string {
	return fmt.Sprintf("evo.selector.Generational{SurvivalRate: %f, MaxStagnation: %d, MutateOnlyProbability: %f, InterspeciesMateProbability; %f}",
		h.SurvivalRate, h.MaxStagnation, h.MutateOnlyProbability, h.InterspeciesMateProbability)
}

// Select the genomes to keep and to become parents of the next generation
func (h *Generational) Select(p evo.Population) (keep []evo.Genome, parents [][]evo.Genome, err error) {

	// Identify the best genome in the population
	sort.Sort(sort.Reverse(p.Genomes))
	b := p.Genomes[0]

	// Determine if experiment is super stagnated
	rng := random.New()
	super := h.stagnation >= 2*h.MaxStagnation
	if super {
		log.Println("Super stagnation")

		// Ignore the best, go with the least complex
		var w evo.Genome
		w = b
		for i := len(p.Genomes) - 1; i >= 0; i-- {
			if w.Complexity() > p.Genomes[i].Complexity() {
				w = p.Genomes[i]
			}
		}

		// Seed a whole new population using the best genome
		parents = make([][]evo.Genome, len(p.Genomes))
		for i := 0; i < len(parents); i++ {
			g := evo.Genome{
				Encoded: evo.Substrate{
					Nodes: make([]evo.Node, len(w.Encoded.Nodes)),
					Conns: make([]evo.Conn, len(w.Encoded.Conns)),
				},
			}
			copy(g.Encoded.Nodes, w.Encoded.Nodes)
			copy(g.Encoded.Conns, w.Encoded.Conns)
			for i := 0; i < len(g.Encoded.Conns); i++ {
				g.Encoded.Conns[i].Weight += rng.NormFloat64()
				g.Encoded.Conns[i].Enabled = rng.Float64() < 50
			}
			parents[i] = []evo.Genome{g}
		}

		// Reset the selector's properites
		h.stagnation = 0
		h.maxFitness = 0.0
		return
	}

	// Map the genomes to the species. They will already be sorted in reverse fitness order by
	// above sort statement
	g2s := make(map[int]evo.Genomes, len(p.Species))
	for _, g := range p.Genomes {
		g2s[g.SpeciesID] = append(g2s[g.SpeciesID], g)
	}

	// Map the species
	i2s := make(map[int]evo.Species)
	for _, s := range p.Species {
		i2s[s.ID] = s
	}

	// Keep the elite
	keep = make([]evo.Genome, 0, len(p.Species))
	for id, gs := range g2s {
		if id == b.SpeciesID || i2s[id].Stagnation < h.MaxStagnation {
			keep = append(keep, gs[0])
		}
	}

	// Cull the pool of stagnant species and non-survivors
	for id, gs := range g2s {
		switch {
		case id == b.SpeciesID:
			g2s[id] = []evo.Genome{b}
		case i2s[id].Stagnation >= h.MaxStagnation:
			delete(g2s, id)
		default:
			n := int(math.Max(1.0, float64(len(gs))*h.SurvivalRate))
			g2s[id] = gs[:n]
		}
	}

	// Calculate total of average fitnesses from remaining species
	t := 0.0
	for _, gs := range g2s {
		t += gs.AvgFitness()
	}

	// Create the parent groupings
	parents = make([][]evo.Genome, 0, len(p.Genomes)-len(keep))
	for len(keep)+len(parents) < len(p.Genomes) {
		sid := roulette(rng, g2s, t)
		// Identify parent 1
		couple := make([]evo.Genome, 0, 2)
		gs := g2s[sid]
		couple = append(couple, gs[rand.Intn(len(gs))])

		// Identify parent 2, if any
		if rng.Float64() > h.MutateOnlyProbability {

			if len(g2s) > 1 && rng.Float64() < h.InterspeciesMateProbability {
				old := sid
				for id, _ := range g2s {
					if id != old {
						sid = id
					}
				}
			}
			gs = g2s[sid]
			couple = append(couple, gs[rand.Intn(len(gs))])
		}
		parents = append(parents, couple)
	}
	return
}

// Use a weight-fitness roulette to pick the species and then randomly choose a parent
func roulette(rng *rand.Rand, g2s map[int]evo.Genomes, tot float64) int {

	t := rng.Float64() * tot
	x := 0.0
	for id, gs := range g2s {
		x += gs.AvgFitness()
		if x >= t {
			return id
		}
	}
	return -1 // Should never reach this
}

// Watch the population for generational changes to look out for super stagnation
func (h *Generational) Watch(p evo.Population) error {

	// Find the maximum fitness of the population
	sort.Sort(sort.Reverse(p.Genomes))
	f := p.Genomes[0].Fitness

	// Update the overall stagnation
	if f > h.maxFitness {
		h.maxFitness = f
		h.stagnation = 0
	} else {
		h.stagnation++
	}
	return nil
}
