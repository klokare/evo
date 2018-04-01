package neat

import (
	"errors"
	"sort"

	"github.com/klokare/evo"
)

// Known errors
var (
	ErrNoParents      = errors.New("NEAT crosser requires at least 1 parent")
	ErrTooManyParents = errors.New("NEAT crosser does not support more than 2 parents")
)

// Crosser combines 1 or more parents to create an offspring genome.
type Crosser struct {
	EnableProbability       float64
	DisableEqualParentCheck bool
	DisableSortCheck        bool
	evo.Comparison
}

// Cross the parents and create a new offspring, using the sequence to assign a new ID. There is a
// chance that connections disabled in one of the parents will also be disabled in the child.
func (z *Crosser) Cross(parents ...evo.Genome) (child evo.Genome, err error) {

	// Check for errors
	if len(parents) == 0 {
		err = ErrNoParents
		return
	} else if len(parents) > 2 {
		err = ErrTooManyParents
		return
	}

	// Special case: single parent
	rng := evo.NewRandom()
	p1 := parents[0]
	if len(parents) == 1 {

		// Clone the parent and return
		child = evo.Genome{
			Encoded: evo.Substrate{
				Nodes: make([]evo.Node, len(p1.Encoded.Nodes)),
				Conns: make([]evo.Conn, len(p1.Encoded.Conns)),
			},
			Traits: make([]float64, len(p1.Traits)),
		}
		copy(child.Encoded.Nodes, p1.Encoded.Nodes)
		copy(child.Encoded.Conns, p1.Encoded.Conns)
		copy(child.Traits, p1.Traits)

	} else {

		// Ensure the more fit parent is in position 1
		p2 := parents[1]
		if z.Compare(p1, p2) < 0 {
			p1, p2 = p2, p1
		}
		same := p1.Fitness == p2.Fitness && !z.DisableEqualParentCheck

		// Create the child
		child = evo.Genome{
			Encoded: evo.Substrate{
				Nodes: crossNodes(rng, p1.Encoded.Nodes, p2.Encoded.Nodes, same, !z.DisableSortCheck),
				Conns: crossConns(rng, p1.Encoded.Conns, p2.Encoded.Conns, same, !z.DisableSortCheck),
			},
			Traits: crossTraits(rng, p1.Traits, p2.Traits),
		}
	}

	// Possibly re-enable disabled connections
	for i, c := range child.Encoded.Conns {
		if !c.Enabled {
			if rng.Float64() < z.EnableProbability {
				child.Encoded.Conns[i].Enabled = true
			}
		}
	}

	// Ensure substrate is sorted
	sort.Slice(child.Encoded.Nodes, func(i, j int) bool { return child.Encoded.Nodes[i].Compare(child.Encoded.Nodes[j]) < 0 })
	sort.Slice(child.Encoded.Conns, func(i, j int) bool { return child.Encoded.Conns[i].Compare(child.Encoded.Conns[j]) < 0 })
	return
}

func crossNodes(rng evo.Random, nodes1, nodes2 []evo.Node, same, check bool) (nodes []evo.Node) {

	// Sort the nodes
	if check {
		tmp := make([]evo.Node, len(nodes1)) // work with a copy to be safe for concurrency
		copy(tmp, nodes1)
		sort.Slice(tmp, func(i, j int) bool { return tmp[i].Compare(tmp[j]) < 0 })
		nodes1 = tmp

		tmp = make([]evo.Node, len(nodes2)) // work with a copy to be safe for concurrency
		copy(tmp, nodes2)
		sort.Slice(tmp, func(i, j int) bool { return tmp[i].Compare(tmp[j]) < 0 })
		nodes2 = tmp
	}

	// Iterate the nodes and look for differences
	var i, j int
	nodes = make([]evo.Node, 0, len(nodes1)+5)
	for i < len(nodes1) && j < len(nodes2) {
		switch nodes1[i].Compare(nodes2[j]) {
		case -1:
			nodes = append(nodes, nodes1[i])
			i++
		case 1.0:
			if same {
				nodes = append(nodes, nodes2[j])
			}
			j++
		default:
			if rng.Float64() < 0.5 {
				nodes = append(nodes, nodes1[i])
			} else {
				nodes = append(nodes, nodes2[j])
			}
			i++
			j++
		}
	}

	// Add remaining unmatched nodes
	for i < len(nodes1) {
		nodes = append(nodes, nodes1[i])
		i++
	}

	for same && j < len(nodes2) {
		nodes = append(nodes, nodes2[j])
		j++
	}
	return
}

func crossConns(rng evo.Random, conns1, conns2 []evo.Conn, same, check bool) (conns []evo.Conn) {

	// Sort the connections
	if check {
		tmp := make([]evo.Conn, len(conns1)) // work with a copy to be safe for concurrency
		copy(tmp, conns1)
		sort.Slice(tmp, func(i, j int) bool { return tmp[i].Compare(tmp[j]) < 0 })
		conns1 = tmp

		tmp = make([]evo.Conn, len(conns2)) // work with a copy to be safe for concurrency
		copy(tmp, conns2)
		sort.Slice(tmp, func(i, j int) bool { return tmp[i].Compare(tmp[j]) < 0 })
		conns2 = tmp
	}

	// Iterate the connections and look for differences
	var i, j int
	for i < len(conns1) && j < len(conns2) {
		switch conns1[i].Compare(conns2[j]) {
		case -1:
			conns = append(conns, conns1[i])
			i++
		case 1.0:
			if same {
				conns = append(conns, conns2[j])
			}
			j++
		default:
			var c evo.Conn
			if rng.Float64() < 0.5 {
				c = conns1[i]
			} else {
				c = conns2[j]
			}
			c.Enabled = conns1[i].Enabled && conns2[j].Enabled // Will disable if disabled in either parent
			conns = append(conns, c)
			i++
			j++
		}
	}

	// Add remaining unmatched connections
	for i < len(conns1) {
		conns = append(conns, conns1[i])
		i++
	}

	for same && j < len(conns2) {
		conns = append(conns, conns2[j])
		j++
	}
	return
}

func crossTraits(rng evo.Random, traits1, traits2 []float64) (traits []float64) {
	traits = make([]float64, len(traits1))
	for i := 0; i < len(traits); i++ {
		if rng.Float64() < 0.5 {
			traits[i] = traits1[i]
		} else {
			traits[i] = traits2[i]
		}
	}
	return
}
