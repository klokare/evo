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
	DecayRate                   float64
	evo.Comparison

	// Internal state
	prevBest     int64   // The ID of the best genome from previous generations
	lastImproved int     // Number of iterations since the last improvement
	mop          float64 // stored mutate only probability for restoring when toggled
}

// Select the genomes to continue and those to become parents
func (s *Selector) Select(pop evo.Population) (continuing []evo.Genome, parents [][]evo.Genome, err error) {

	// Sort and rank the genomes
	genomes := make([]evo.Genome, len(pop.Genomes))
	copy(genomes, pop.Genomes)
	ranks := sortRank(s.Comparison, genomes)

	// Identify the current champion
	best := genomes[0]

	// TRIAL:
	// Override species decay with the decay rate times the max age of genomes. This is already different than
	// the current implementation in that it does not reset. Or, well, it does if only 1 elite is used.
	bs := evo.GroupBySpecies(genomes)
	decay := make([]float64, len(bs))
	for sid, gs := range bs {
		ma := 0.0
		for _, g := range gs {
			a := float64(g.Age)
			if ma < a {
				ma = a
			}
		}
		d := ma * s.DecayRate
		if d > 1.0 {
			d = 1.0
		}
		decay[sid] = d
	}

	// Calculate the avg ranking by species, adjusted by decay, and the total of avg to be used
	// below for roulette
	var tot float64
	avg := make([]float64, len(bs))
	for sid, gs := range bs {

		// Calculate the average
		for _, g := range gs {
			avg[sid] += ranks[g.ID]
		}
		avg[sid] /= float64(len(gs))

		// Adjust by decay
		avg[sid] *= (1.0 - decay[sid])

		// Append to total
		tot += avg[sid]
	}

	// Determine continuing
	if !s.DisableContinuing {
		var bsid int // index of best's species
		continuing = make([]evo.Genome, 0, len(avg))
		for sid, gs := range bs {
			if avg[sid] > 0.0 { // species is not stagnant
				continuing = append(continuing, gs[0])
			}
			if gs[0].ID == best.ID {
				bsid = sid
			}
		}

		// Ensure best is continuing, if necessary
		if avg[bsid] == 0.0 { // best is stagnant
			found := false
			for _, g := range continuing {
				if s.Compare(best, g) == 0 { // A champion exists in another species
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
	off := make([]int, len(avg))
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
	rng := evo.NewRandom()
	adjCounts(rng, off, cnt, tgt)

	// Generate parents
	parents = make([][]evo.Genome, 0, tgt)
	for sid1, n := range off {

		for i := 0; i < n; i++ {
			// Pick parent 1 from the species
			p1 := roulette(rng, bs[sid1], ranks)

			// This is a single parent
			if rng.Float64() < s.MutateOnlyProbability {
				parents = append(parents, []evo.Genome{p1})
				continue
			}

			// Pick a second parent
			// TODO: add interspecies
			sid := sid1
			if len(avg) > 1 && rng.Float64() < s.InterspeciesMateProbability {
				idxs := rng.Perm(len(bs))
				for _, i := range idxs {
					if i != sid1 {
						sid = i
						break
					}
				}
			}
			p2 := roulette(rng, bs[sid], ranks)
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
func adjCounts(rng evo.Random, off []int, cnt, tgt int) {

	// Calaculate the adjustment
	diff := tgt - cnt
	adj := 1
	if diff < 0 {
		diff = -diff
		adj = -1
	}

	// Make multiple passes
	for cnt != tgt {
		idxs := rng.Perm(len(off))
		for _, i := range idxs {
			if off[i] > 0 && off[i]+adj > 0 {
				off[i] += adj
				cnt += adj
				if cnt == tgt {
					break
				}
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
