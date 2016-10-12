package evo

import (
	"math/rand"
	"sort"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPositionCompare(t *testing.T) {

	Convey("Given pairs of positions", t, func() {

		Convey("When positions are equal", func() {
			a := Position{Layer: 1.0, X: 2.0, Y: 3.0, Z: 4.0}
			b := Position{Layer: 1.0, X: 2.0, Y: 3.0, Z: 4.0}
			Convey("The results should be 0", func() {
				So(a.Compare(b), ShouldEqual, 0)
			})
			Convey("The result should be reciprical", func() {
				So(a.Compare(b), ShouldEqual, b.Compare(a))
			})
		})

		Convey("When a.Layer is less than b.Layer", func() {
			a := Position{Layer: 1.0, X: 2.0, Y: 3.0, Z: 4.0}
			b := Position{Layer: 2.0, X: 2.0, Y: 3.0, Z: 4.0}
			Convey("The results should be -1", func() {
				So(a.Compare(b), ShouldEqual, -1)
			})
		})
		Convey("When a.Layer is greater than b.Layer", func() {
			a := Position{Layer: 2.0, X: 2.0, Y: 3.0, Z: 4.0}
			b := Position{Layer: 1.0, X: 2.0, Y: 3.0, Z: 4.0}
			Convey("The results should be 1", func() {
				So(a.Compare(b), ShouldEqual, 1)
			})
		})

		Convey("When layer is the same and a.X is less than b.X", func() {
			a := Position{Layer: 1.0, X: 1.0, Y: 3.0, Z: 4.0}
			b := Position{Layer: 1.0, X: 2.0, Y: 3.0, Z: 4.0}
			Convey("The results should be -1", func() {
				So(a.Compare(b), ShouldEqual, -1)
			})
		})
		Convey("When layer is the same and a.X is greater than b.X", func() {
			a := Position{Layer: 1.0, X: 2.0, Y: 3.0, Z: 4.0}
			b := Position{Layer: 1.0, X: 1.0, Y: 3.0, Z: 4.0}
			Convey("The results should be 1", func() {
				So(a.Compare(b), ShouldEqual, 1)
			})
		})

		Convey("When layer and X are the same and a.Y is less than b.Y", func() {
			a := Position{Layer: 1.0, X: 2.0, Y: 2.0, Z: 4.0}
			b := Position{Layer: 1.0, X: 2.0, Y: 3.0, Z: 4.0}
			Convey("The results should be -1", func() {
				So(a.Compare(b), ShouldEqual, -1)
			})
		})
		Convey("When layer and X are the same and a.Y is greater than b.Y", func() {
			a := Position{Layer: 1.0, X: 2.0, Y: 3.0, Z: 4.0}
			b := Position{Layer: 1.0, X: 2.0, Y: 2.0, Z: 4.0}
			Convey("The results should be 1", func() {
				So(a.Compare(b), ShouldEqual, 1)
			})
		})

		Convey("When layer, X and Y are the same and a.Z is less than b.Z", func() {
			a := Position{Layer: 1.0, X: 2.0, Y: 3.0, Z: 3.0}
			b := Position{Layer: 1.0, X: 2.0, Y: 3.0, Z: 4.0}
			Convey("The results should be -1", func() {
				So(a.Compare(b), ShouldEqual, -1)
			})
		})
		Convey("When layer, X and Y are the same and a.Z is greater than b.Z", func() {
			a := Position{Layer: 1.0, X: 2.0, Y: 3.0, Z: 4.0}
			b := Position{Layer: 1.0, X: 2.0, Y: 3.0, Z: 3.0}
			Convey("The results should be 1", func() {
				So(a.Compare(b), ShouldEqual, 1)
			})
		})
	})
}

func TestNodeCopare(t *testing.T) {
	Convey("Given pairs of nodes", t, func() {
		Convey("When the nodes' positions are equal", func() {
			a := Node{Position: Position{Layer: 1.0}}
			b := Node{Position: Position{Layer: 1.0}}
			Convey("The result should be 0", func() { So(a.Compare(b), ShouldEqual, 0.0) })
			Convey("The result should be reciprical", func() { So(a.Compare(b), ShouldEqual, b.Compare(a)) })
		})
		Convey("When a's position is less than b's", func() {
			a := Node{Position: Position{Layer: 1.0}}
			b := Node{Position: Position{Layer: 2.0}}
			Convey("The result should be -1", func() { So(a.Compare(b), ShouldEqual, -1.0) })
		})
		Convey("When a's position is greater than b's", func() {
			a := Node{Position: Position{Layer: 2.0}}
			b := Node{Position: Position{Layer: 1.0}}
			Convey("The result should be 1", func() { So(a.Compare(b), ShouldEqual, 1.0) })
		})
	})
}

func TestNodesSort(t *testing.T) {
	Convey("Given a list of nodes", t, func() {
		a := []Node{
			{Position: Position{Layer: 1.0}},
			{Position: Position{Layer: 2.0}},
			{Position: Position{Layer: 3.0}},
			{Position: Position{Layer: 4.0}},
			{Position: Position{Layer: 5.0}},
		}
		Convey("When sorting a randomly ordered copy", func() {
			var b Nodes = make([]Node, 0, len(a))
			idxs := rand.Perm(len(a))
			for _, i := range idxs {
				b = append(b, a[i])
			}
			sort.Sort(b)
			Convey("The nodes should be in order", func() {
				for i := 0; i < len(a); i++ {
					So(a[i].Compare(b[i]), ShouldEqual, 0)
				}
			})
		})
	})
}

func TestConnCompare(t *testing.T) {
	Convey("Given pairs of connections", t, func() {
		Convey("When the connections' sources and targets are equal", func() {
			a := Conn{Source: Position{Layer: 1.0}, Target: Position{Layer: 2.0}}
			b := Conn{Source: Position{Layer: 1.0}, Target: Position{Layer: 2.0}}
			Convey("The result should be 0", func() { So(a.Compare(b), ShouldEqual, 0) })
			Convey("The result should be reciprical", func() { So(a.Compare(b), ShouldEqual, (b.Compare(a))) })
		})
		Convey("When connection a's source is less than b's", func() {
			a := Conn{Source: Position{Layer: 1.0}, Target: Position{Layer: 2.0}}
			b := Conn{Source: Position{Layer: 2.0}, Target: Position{Layer: 2.0}}
			Convey("The result should be -1", func() { So(a.Compare(b), ShouldEqual, -1) })
		})
		Convey("When connection a's is greater than b's", func() {
			a := Conn{Source: Position{Layer: 2.0}, Target: Position{Layer: 2.0}}
			b := Conn{Source: Position{Layer: 1.0}, Target: Position{Layer: 2.0}}
			Convey("The result should be 1", func() { So(a.Compare(b), ShouldEqual, 1) })
		})
		Convey("When sources are the same and connection a's target is less than b's", func() {
			a := Conn{Source: Position{Layer: 1.0}, Target: Position{Layer: 1.0}}
			b := Conn{Source: Position{Layer: 1.0}, Target: Position{Layer: 2.0}}
			Convey("The result should be -1", func() { So(a.Compare(b), ShouldEqual, -1) })
		})
		Convey("When sources are the same and connection a's target is greater than b's", func() {
			a := Conn{Source: Position{Layer: 1.0}, Target: Position{Layer: 2.0}}
			b := Conn{Source: Position{Layer: 1.0}, Target: Position{Layer: 1.0}}
			Convey("The result should be 1", func() { So(a.Compare(b), ShouldEqual, 1) })
		})
	})
}

func TestConnsSort(t *testing.T) {
	Convey("Given a list of connections", t, func() {
		a := []Conn{
			{Source: Position{Layer: 1.0}},
			{Source: Position{Layer: 2.0}},
			{Source: Position{Layer: 3.0}},
			{Source: Position{Layer: 4.0}},
			{Source: Position{Layer: 5.0}},
		}
		Convey("When sorting a randomly ordered copy", func() {
			var b Conns = make([]Conn, 0, len(a))
			idxs := rand.Perm(len(a))
			for _, i := range idxs {
				b = append(b, a[i])
			}
			sort.Sort(b)
			Convey("The conns should be in order", func() {
				for i := 0; i < len(a); i++ {
					So(a[i].Compare(b[i]), ShouldEqual, 0)
				}
			})
		})
	})
}

// -------------- benchmarking choice of position --------------
type PositionSlice []float64

func (a PositionSlice) Compare(b PositionSlice) int {
	for i := 0; i < len(a); i++ {
		if a[i] < b[i] {
			return -1
		} else if a[i] > b[i] {
			return 1
		}
	}
	return 0
}

type PositionIF struct {
	Layer   float64
	X, Y, Z float64
}

func (a PositionIF) Compare(b PositionIF) int {
	if a.Layer == b.Layer {
		if a.X == b.X {
			if a.Y == b.Y {
				if a.Z == b.Z {
					return 0
				} else if a.Z < b.Z {
					return -1
				}
				return 1
			} else if a.Y < b.Y {
				return -1
			}
			return 1
		} else if a.X < b.X {
			return -1
		}
		return 1
	} else if a.Layer < b.Layer {
		return -1
	}
	return 1
}

func BenchmarkPositionIF(b *testing.B) {
	ps := [][]PositionIF{
		{{Layer: 1.0, X: 2.0, Y: 3.0, Z: 4.0}, {Layer: 1.0, X: 2.0, Y: 3.0, Z: 4.0}},
		{{Layer: 1.0, X: 1.0, Y: 3.0, Z: 4.0}, {Layer: 1.0, X: 2.0, Y: 3.0, Z: 4.0}},
		{{Layer: 1.0, X: 2.0, Y: 3.0, Z: 4.0}, {Layer: 1.0, X: 1.0, Y: 3.0, Z: 4.0}},
		{{Layer: 1.0, X: 2.0, Y: 2.0, Z: 4.0}, {Layer: 1.0, X: 2.0, Y: 3.0, Z: 4.0}},
		{{Layer: 1.0, X: 2.0, Y: 3.0, Z: 4.0}, {Layer: 1.0, X: 2.0, Y: 2.0, Z: 4.0}},
		{{Layer: 1.0, X: 2.0, Y: 3.0, Z: 3.0}, {Layer: 1.0, X: 2.0, Y: 3.0, Z: 4.0}},
		{{Layer: 1.0, X: 2.0, Y: 3.0, Z: 4.0}, {Layer: 1.0, X: 2.0, Y: 3.0, Z: 3.0}},
	}
	for i := 0; i < b.N; i++ {
		for _, p := range ps {
			p[0].Compare(p[1])
		}
	}
}

func BenchmarkPositionSwitch(b *testing.B) {
	ps := [][]Position{
		{{Layer: 1.0, X: 2.0, Y: 3.0, Z: 4.0}, {Layer: 1.0, X: 2.0, Y: 3.0, Z: 4.0}},
		{{Layer: 1.0, X: 1.0, Y: 3.0, Z: 4.0}, {Layer: 1.0, X: 2.0, Y: 3.0, Z: 4.0}},
		{{Layer: 1.0, X: 2.0, Y: 3.0, Z: 4.0}, {Layer: 1.0, X: 1.0, Y: 3.0, Z: 4.0}},
		{{Layer: 1.0, X: 2.0, Y: 2.0, Z: 4.0}, {Layer: 1.0, X: 2.0, Y: 3.0, Z: 4.0}},
		{{Layer: 1.0, X: 2.0, Y: 3.0, Z: 4.0}, {Layer: 1.0, X: 2.0, Y: 2.0, Z: 4.0}},
		{{Layer: 1.0, X: 2.0, Y: 3.0, Z: 3.0}, {Layer: 1.0, X: 2.0, Y: 3.0, Z: 4.0}},
		{{Layer: 1.0, X: 2.0, Y: 3.0, Z: 4.0}, {Layer: 1.0, X: 2.0, Y: 3.0, Z: 3.0}},
	}
	for i := 0; i < b.N; i++ {
		for _, p := range ps {
			p[0].Compare(p[1])
		}
	}
}

func BenchmarkPositionSlice(b *testing.B) {
	ps := [][]PositionSlice{
		{{1.0, 2.0, 3.0, 4.0}, {1.0, 2.0, 3.0, 4.0}},
		{{1.0, 1.0, 3.0, 4.0}, {1.0, 2.0, 3.0, 4.0}},
		{{1.0, 2.0, 3.0, 4.0}, {1.0, 1.0, 3.0, 4.0}},
		{{1.0, 2.0, 2.0, 4.0}, {1.0, 2.0, 3.0, 4.0}},
		{{1.0, 2.0, 3.0, 4.0}, {1.0, 2.0, 2.0, 4.0}},
		{{1.0, 2.0, 3.0, 3.0}, {1.0, 2.0, 3.0, 4.0}},
		{{1.0, 2.0, 3.0, 4.0}, {1.0, 2.0, 3.0, 3.0}},
	}
	for i := 0; i < b.N; i++ {
		for _, p := range ps {
			p[0].Compare(p[1])
		}
	}
}
