package crosser

import (
	"sort"

	"github.com/klokare/evo"
	"github.com/klokare/random"
)

// A Multiple crosser creates a new genome from one or more parent genomes
type Multiple struct {
	EnableProbability float64
}

// Cross the genomes to create a new child genome
func (h *Multiple) Cross(gs ...evo.Genome) (evo.Genome, error) {

	// Sort the parents by fitness decending
	sort.Sort(sort.Reverse(evo.Genomes(gs)))

	// Clone the first parent
	g := evo.Genome{
		Encoded: evo.Substrate{
			Nodes: make([]evo.Node, len(gs[0].Encoded.Nodes)),
			Conns: make([]evo.Conn, len(gs[0].Encoded.Conns)),
		},
	}
	copy(g.Encoded.Nodes, gs[0].Encoded.Nodes)
	copy(g.Encoded.Conns, gs[0].Encoded.Conns)

	// Merge in the remaining parents
	rng := random.New()
	p := 1.0 / float64(len(gs))
	for i := 1; i < len(gs); i++ {
		same := gs[0].Fitness == gs[i].Fitness
		for _, n := range gs[i].Encoded.Nodes {
			j := findNodeIdx(g.Encoded.Nodes, n.Position)
			if j == -1 {
				if same {
					g.Encoded.Nodes = append(g.Encoded.Nodes, n)
				}
			} else {
				if rng.Float64() < p {
					g.Encoded.Nodes[j].ActivationType = n.ActivationType
				}
			}
		}

		for _, c := range gs[i].Encoded.Conns {
			j := findConnIdx(g.Encoded.Conns, c.Source, c.Target)
			if j == -1 {
				if same {
					g.Encoded.Conns = append(g.Encoded.Conns, c)
				}
			} else {
				if rng.Float64() < p {
					g.Encoded.Conns[j].Weight = c.Weight
				}
				g.Encoded.Conns[j].Enabled = g.Encoded.Conns[j].Enabled && c.Enabled // Disables if disabled in either parent
			}
		}
	}

	// Reenable connections
	for i, c := range g.Encoded.Conns {
		if !c.Enabled && rng.Float64() < h.EnableProbability {
			g.Encoded.Conns[i].Enabled = true
		}
	}

	// Sort the substrate and return
	sort.Sort(g.Encoded.Nodes)
	sort.Sort(g.Encoded.Conns)
	return g, nil
}

func findNodeIdx(ns []evo.Node, p evo.Position) int {
	for i, n := range ns {
		if n.Position.Compare(p) == 0 {
			return i
		}
	}
	return -1
}

func findConnIdx(cs []evo.Conn, s, t evo.Position) int {
	for i, c := range cs {
		if c.Source.Compare(s) == 0 && c.Target.Compare(t) == 0 {
			return i
		}
	}
	return -1
}
