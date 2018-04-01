package neat

import (
	"errors"
	"math"

	"github.com/klokare/evo"
)

// Known errors
var (
	ErrInvalidPopulationSize = errors.New("invalid population size")
)

// Selector determines which genomes continue and which become parents
type Selector struct {
	PopulationSize              int
	MutateOnlyProbability       float64
	InterspeciesMateProbability float64
	DisableContinuing           bool
	Elitism                     float64
	SurvivalRate                float64
	evo.Comparison

	// Internal state
	prevBest     int64   // The ID of the best genome from previous generations
	lastImproved int     // Number of iterations since the last improvement
	mop          float64 // stored mutate only probability for restoring when toggled
}

// Select the genomes to continue and those to become parents
func (s *Selector) Select(pop evo.Population) (continuing []evo.Genome, parents [][]evo.Genome, err error) {

	// Sort and rank the genomes
	ranks := sortRank(s.Comparison, pop.Genomes)

	// Identify the current champion
	best := pop.Genomes[0]

	// Restart the populaton if necessary
	stagnant := true
	for _, s := range pop.Species {
		if s.Decay < 1.0 {
			stagnant = false
			break
		}
	}
	if stagnant {
		if !s.DisableContinuing {
			continuing = []evo.Genome{best}
		}
		parents = make([][]evo.Genome, s.PopulationSize-len(continuing))
		for i := 0; i < len(parents); i++ {
			parents[i] = []evo.Genome{best}
		}
		for i := 0; i < len(pop.Species); i++ {
			pop.Species[i].Decay = 0.0
			pop.Species[i].Champion = 0
		}
		return
	}

	// Divide the genomes into species. Use this loop to also begin the average ranking calculation
	var genomes []evo.Genome
	var ok bool
	g2s := make(map[int64][]evo.Genome, len(pop.Species))
	avg := make(map[int64]float64, len(pop.Species))
	for _, g := range pop.Genomes {

		// Append the genome
		if genomes, ok = g2s[g.SpeciesID]; !ok {
			genomes = make([]evo.Genome, 0, 10)
		}
		genomes = append(genomes, g)
		g2s[g.SpeciesID] = genomes

		// Add in the ranking
		avg[g.SpeciesID] += ranks[g.ID]
	}

	// Finish the average ranking, apply the decay rates, and calculate the total
	tot := 0.0
	for _, z := range pop.Species {
		n := float64(len(g2s[z.ID]))
		x := avg[z.ID] * (1.0 - z.Decay)
		if n > 0.0 {
			x /= n
		}
		avg[z.ID] = x
		tot += x
	}

	// Determine continuing
	if !s.DisableContinuing {
		continuing = make([]evo.Genome, 0, len(avg))
		for sid, genomes := range g2s {
			if avg[sid] > 0.0 { // species is not stagnant
				continuing = append(continuing, genomes[0])
			}
		}

		// Ensure best is continuing, if necessary
		if avg[best.SpeciesID] == 0.0 { // best is stagnant
			found := false
			for _, g := range continuing {
				if s.Compare(best, g) == 0 { // Another champion exists
					found = true
					break
				}
			}
			if !found {
				continuing = append(continuing, best)
				//	parents = append(parents, []evo.Genome{best}) // Give the species 1 more chance
			}
		}
	}

	// Calculate offspring
	cnt := 0
	tgt := s.PopulationSize - len(continuing) - len(parents)
	if tgt <= 0 {
		return // population size fulfilled with continuing
	}
	off := make(map[int64]int, len(avg))
	for sid, x := range avg {
		if x > 0.0 {
			ns := int(math.Floor(float64(tgt) * x / tot))
			if ns < 1 {
				ns = 1
			}
			off[sid] = ns
			cnt += ns
		}
	}

	// Handle rounding errors by adjusting 1 offspring at a time
	adjCounts(off, cnt, tgt)

	// Generate parents
	rng := evo.NewRandom()
	parents = make([][]evo.Genome, 0, tgt)
	for sid1, n := range off {

		for i := 0; i < n; i++ {
			// Pick parent 1 from the species
			p1 := roulette(rng, g2s[sid1], ranks)

			// This is a single parent
			if rng.Float64() < s.MutateOnlyProbability {
				parents = append(parents, []evo.Genome{p1})
				continue
			}

			// Pick a second parent
			// TODO: add interspecies
			sid := sid1
			if len(avg) > 1 && rng.Float64() < s.InterspeciesMateProbability {
				idxs := rng.Perm(len(pop.Species))
				for _, i := range idxs {
					if pop.Species[i].ID != sid1 {
						sid = pop.Species[i].ID
						break
					}
				}
			}
			p2 := roulette(rng, g2s[sid], ranks)
			parents = append(parents, []evo.Genome{p1, p2})
		}
	}
	return
}

// ToggleMutateOnly puts the selector into a mutate only mode when on is true
func (s *Selector) ToggleMutateOnly(on bool) error {
	if s.mop == 0.0 {
		s.mop = s.MutateOnlyProbability
	}
	if on {
		s.MutateOnlyProbability = 1.0
	} else {
		s.MutateOnlyProbability = s.mop
	}
	return nil
}

// Sort and rank the genomes. The sort is a reverse sort (best in first position) and deterministic
// (always producing the same order). The rank is relative order, allowing for ties, where best is
// float64(len(genomes)) and worst, assuming no tie, is 1.
func sortRank(fn evo.Comparison, genomes []evo.Genome) (ranks map[int64]float64) {

	// Sort the genomes using main compare function and others to get separation in rank
	evo.SortBy(genomes, fn, evo.ByComplexity, evo.ByAge)

	// Reverse the sort
	for i, j := 0, len(genomes)-1; i < j; i, j = i+1, j-1 {
		genomes[i], genomes[j] = genomes[j], genomes[i]
	}

	// Create the rankings
	ranks = make(map[int64]float64, len(genomes))
	n := len(genomes)
	ranks[genomes[0].ID] = float64(n)
	for i := 1; i < n; i++ {
		if fn.Compare(genomes[i], genomes[i-1]) == 0 { // deserves same rank)
			ranks[genomes[i].ID] = ranks[genomes[i-1].ID]
		} else {
			ranks[genomes[i].ID] = float64(n - i)
		}
	}
	return
}

// Adjust counts by assigning the difference to the most fit species which is identified by the
// species with the most offspring assigned and the lowest ID.
func adjCounts(off map[int64]int, cnt, tgt int) {

	// Calaculate the adjustment
	diff := tgt - cnt
	if diff < 0 {
		diff = -diff
	}
	adj := 1
	if cnt > tgt {
		adj = -1
	}

	// Make multiple passes
	for i := 0; i < 2*diff && cnt != tgt; i++ {
		for sid := range off {
			n := off[sid]
			if n+adj <= 0 {
				continue
			}
			off[sid] = n + adj
			cnt += adj
			if cnt == tgt {
				break
			}
		}
	}
	return
}

func roulette(rng evo.Random, genomes []evo.Genome, ranks map[int64]float64) (parent evo.Genome) {

	// Determine total
	rs := make([]float64, len(genomes))
	tot := 0.0
	for i, g := range genomes {
		rs[i] = ranks[g.ID]
		tot += rs[i]
	}

	// Spin the wheel
	tgt := rng.Float64() * tot
	sum := 0.0
	for i, g := range genomes {
		sum += rs[i]
		if sum >= tgt {
			parent = g
			break
		}
	}
	return
}
