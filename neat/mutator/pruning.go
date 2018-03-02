package mutator

import (
	"sort"

	"github.com/klokare/evo"
)

// Pruning mutator disables connections and, occasionally, removes them
type Pruning struct {
	DisableProbability float64
	PruneProbability   float64
}

// Mutate a genome by disabling or removing connections. If nodes are left stranded, they, too, are
// removed.
func (p Pruning) Mutate(g *evo.Genome) (err error) {
	rng := evo.NewRandom()
	if rng.Float64() < p.DisableProbability {
		if rng.Float64() < p.PruneProbability {
			pruneConn(rng, &g.Encoded)
		} else {
			disableConn(rng, &g.Encoded)
		}
	}
	return nil
}

func pruneConn(rng evo.Random, sub *evo.Substrate) {

	// Improve search speed
	sort.Slice(sub.Nodes, func(i, j int) bool { return sub.Nodes[i].Compare(sub.Nodes[j]) < 0 })
	sort.Slice(sub.Conns, func(i, j int) bool { return sub.Conns[i].Compare(sub.Conns[j]) < 0 })

	// There are no connections to delete
	if len(sub.Conns) == 0 {
		return
	}

	// Choose a connection at random
	idx := rng.Intn(len(sub.Conns))
	c := sub.Conns[idx]

	// Remove the connection from the substrate
	removeConn(sub, c)
}

func disableConn(rng evo.Random, sub *evo.Substrate) {
	idxs := rng.Perm(len(sub.Conns))
	for _, i := range idxs {
		if sub.Conns[i].Enabled {
			sub.Conns[i].Enabled = false
			break
		}
	}
}

// WithPruning adds the configured pruning mutator to the experiment
func WithPruning(cfg evo.Configurer) evo.Option {
	return func(e *evo.Experiment) (err error) {
		z := new(Pruning)
		if err = cfg.Configure(z); err != nil {
			return
		}

		// Do not continue if there is no chance for mutation
		if z.DisableProbability == 0.0 {
			return
		}

		e.Mutators = append(e.Mutators, z)
		return
	}
}
