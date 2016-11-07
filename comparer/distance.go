package comparer

import (
	"fmt"
	"math"

	"github.com/klokare/evo"
)

// A Distance comparer measures the structual differences of two substrates
//
// NOTE: ConnsCoefficient is applied to count of connections not found in the other
// genome. Replaces disjoint and excess coefficients of Stanley's NEAT paper.
type Distance struct {
	ConnsCoefficient  float64 `evo:"conns-coefficient"`
	NodesCoefficient  float64 `evo:"nodes-coefficient"`
	WeightCoefficient float64 `evo:"weight-coefficient"`
}

func (h Distance) String() string {
	return fmt.Sprintf("evo.comparer.Distance{ConnsCoefficient: %f, NodesCoefficient: %f, WeightCoefficient: %f}",
		h.ConnsCoefficient, h.NodesCoefficient, h.WeightCoefficient)
}

// Compare two substrates and return the structual distance between the two.
func (h *Distance) Compare(s1, s2 evo.Substrate) (float64, error) {

	n := compareNodes(s1.Nodes, s2.Nodes)
	c, w := compareConns(s1.Conns, s2.Conns)

	return n*h.NodesCoefficient +
		c*h.ConnsCoefficient +
		w*h.WeightCoefficient, nil
}

func compareNodes(ns1, ns2 evo.Nodes) float64 {
	cnt1 := len(ns1)
	cnt2 := len(ns2)
	for _, n := range ns1 {
		if findNodeIdx(ns2, n.Position) != -1 {
			cnt1--
			cnt2--
		}
	}
	return float64(cnt1 + cnt2)
}

func compareConns(cs1, cs2 evo.Conns) (float64, float64) {
	var w, n float64
	cnt1 := len(cs1)
	cnt2 := len(cs2)
	for _, c := range cs1 {
		if i := findConnIdx(cs2, c.Source, c.Target); i != -1 {
			n += 1.0
			w += math.Abs(c.Weight - cs2[i].Weight)
			cnt1--
			cnt2--
		}
	}
	if n > 0.0 {
		w /= n
	}
	return float64(cnt1 + cnt2), w
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
