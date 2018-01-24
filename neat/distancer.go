package neat

import (
	"math"
	"sort"

	"github.com/klokare/evo"
)

// Distancer calculates the compatibility distance between two genomes
type Distancer interface {
	Distance(a, b evo.Genome) (float64, error)
}

// Compatibility distance measurer using the methods described by Stanley
type Compatibility struct {
	NodesCoefficient      float64
	ActivationCoefficient float64
	ConnsCoefficient      float64
	WeightCoefficient     float64
}

// Distance is the compatibility distance based on encoded substrates
func (c Compatibility) Distance(a, b evo.Genome) (d float64, err error) {

	// Calculate the number of disjoint/excess nodes
	n, f := nodeDistance(a.Encoded.Nodes, b.Encoded.Nodes)

	// Calculate the number of disjoint/excess conns as well as the average difference in weight maggnitude
	s, w := connDistance(a.Encoded.Conns, b.Encoded.Conns)

	// Return the compatibility distance
	d = n*c.NodesCoefficient + f*c.ActivationCoefficient + s*c.ConnsCoefficient + w*c.WeightCoefficient
	return
}

func nodeDistance(nodes1, nodes2 []evo.Node) (c, a float64) {

	// Sort the nodes
	sort.Slice(nodes1, func(i, j int) bool { return nodes1[i].Compare(nodes1[j]) < 0 })
	sort.Slice(nodes2, func(i, j int) bool { return nodes2[i].Compare(nodes2[j]) < 0 })

	// Iterate the nodes and look for differences
	var i, j int
	var n float64
	for i < len(nodes1) && j < len(nodes2) {
		switch nodes1[i].Compare(nodes2[j]) {
		case -1:
			c += 1.0
			i++
		case 1.0:
			c += 1.0
			j++
		default:
			n += 1.0
			if nodes1[i].Activation != nodes2[j].Activation {
				a += 1.0
			}
			i++
			j++
		}
	}

	// Add remaining unmatched nodes
	c += float64(len(nodes1) - i + len(nodes2) - j)
	if n > 0 {
		a /= n
	}
	return
}

func connDistance(conns1, conns2 []evo.Conn) (c float64, w float64) {

	// Sort the connections
	sort.Slice(conns1, func(i, j int) bool { return conns1[i].Compare(conns1[j]) < 0 })
	sort.Slice(conns2, func(i, j int) bool { return conns2[i].Compare(conns2[j]) < 0 })

	// Iterate the connections and look for differences
	var i, j int
	var n float64
	for i < len(conns1) && j < len(conns2) {
		switch conns1[i].Compare(conns2[j]) {
		case -1:
			c += 1.0
			i++
		case 1.0:
			c += 1.0
			j++
		default:
			n += 1.0
			w += math.Abs(conns1[i].Weight - conns2[j].Weight)
			i++
			j++
		}
	}

	// Add remaining unmatched connections
	c += float64(len(conns1) - i + len(conns2) - j)

	// Calculate the average difference in weight maggnitude
	if n > 0.0 {
		w = w / n
	}
	return
}
