package mutator

import (
	"math"
	"testing"

	"github.com/klokare/evo"
	"github.com/klokare/random"
	. "github.com/smartystreets/goconvey/convey"
)

func TestComplexifyMutate(t *testing.T) {

	Convey("Given a complexify mutator and a genome", t, func() {
		h := &Complexify{}
		g := &evo.Genome{
			Encoded: evo.Substrate{
				Nodes: []evo.Node{
					{Position: evo.Position{Layer: 0.0, X: 0.0}},
					{Position: evo.Position{Layer: 0.5, X: 0.75}},
					{Position: evo.Position{Layer: 1.0, X: 0.5}},
				},
				Conns: []evo.Conn{
					{
						Source:  evo.Position{Layer: 0.0, X: 0.0},
						Target:  evo.Position{Layer: 1.0, X: 0.5},
						Weight:  2.0,
						Enabled: true,
					},
				},
			},
		}
		Convey("When add node probability is 1.0", func() {
			h.AddNodeProbability = 1.0
			a := len(g.Encoded.Nodes)
			err := h.Mutate(g)
			Convey("There should be no error", func() { So(err, ShouldBeNil) })
			Convey("The genome's complexity should increase", func() {
				So(len(g.Encoded.Nodes), ShouldEqual, a+1)
			})
		})
		Convey("When add conn probability is 1.0", func() {
			h.AddConnProbability = 1.0
			a := len(g.Encoded.Conns)
			err := h.Mutate(g)
			Convey("There should be no error", func() { So(err, ShouldBeNil) })
			Convey("The genome's complexity should increase", func() {
				So(len(g.Encoded.Conns), ShouldEqual, a+1)
			})
		})
	})
}

func TestComplexifyAddNode(t *testing.T) {

	Convey("Given a complexify mutator and a genome", t, func() {
		h := &Complexify{}
		g := &evo.Genome{
			Encoded: evo.Substrate{
				Nodes: []evo.Node{
					{Position: evo.Position{Layer: 0.0, X: 0.0}},
					{Position: evo.Position{Layer: 1.0, X: 0.5}},
				},
				Conns: []evo.Conn{
					{
						Source:  evo.Position{Layer: 0.0, X: 0.0},
						Target:  evo.Position{Layer: 1.0, X: 0.5},
						Weight:  2.0,
						Enabled: true,
					},
				},
			},
		}
		rng := random.New()
		Convey("When there is a connection that can be split", func() {
			h.addNode(rng, g)
			Convey("There should be a new node", func() {
				So(len(g.Encoded.Nodes), ShouldEqual, 3)
			})
			Convey("The original connection should be disabled", func() {
				i := findConnIdx(g.Encoded.Conns,
					evo.Position{Layer: 0.0, X: 0.0},
					evo.Position{Layer: 1.0, X: 0.5})
				So(g.Encoded.Conns[i].Enabled, ShouldBeFalse)
			})
			Convey("The incoming connection should be enabled", func() {
				i := findConnIdx(g.Encoded.Conns,
					evo.Position{Layer: 0.0, X: 0.0},
					evo.Position{Layer: 0.5, X: 0.25})
				So(g.Encoded.Conns[i].Enabled, ShouldBeTrue)
			})
			Convey("The incoming connection should have a weight of 1.0", func() {
				i := findConnIdx(g.Encoded.Conns,
					evo.Position{Layer: 0.0, X: 0.0},
					evo.Position{Layer: 0.5, X: 0.25})
				So(g.Encoded.Conns[i].Weight, ShouldEqual, 1.0)

			})
			Convey("The outgoing connection should be enabled", func() {
				i := findConnIdx(g.Encoded.Conns,
					evo.Position{Layer: 0.5, X: 0.25},
					evo.Position{Layer: 1.0, X: 0.5})
				So(g.Encoded.Conns[i].Enabled, ShouldBeTrue)
			})
			Convey("The outgoing connection should have the original weight", func() {
				i := findConnIdx(g.Encoded.Conns,
					evo.Position{Layer: 0.5, X: 0.25},
					evo.Position{Layer: 1.0, X: 0.5})
				So(g.Encoded.Conns[i].Weight, ShouldEqual, 2.0)
			})
		})

		Convey("When there is more than one connection to split", func() {
			cnts := make(map[float64]int, 2)
			for i := 0; i < 1000; i++ {
				g := &evo.Genome{
					Encoded: evo.Substrate{
						Nodes: []evo.Node{
							{Position: evo.Position{Layer: 0.0, X: 0.0}},
							evo.Node{Position: evo.Position{Layer: 1.0, X: 0.0}},
							{Position: evo.Position{Layer: 1.0, X: 0.5}},
						},
						Conns: []evo.Conn{
							{
								Source:  evo.Position{Layer: 0.0, X: 0.0},
								Target:  evo.Position{Layer: 1.0, X: 0.5},
								Weight:  2.0,
								Enabled: true,
							},
							{
								Source:  evo.Position{Layer: 0.0, X: 1.0},
								Target:  evo.Position{Layer: 1.0, X: 0.5},
								Weight:  2.0,
								Enabled: true,
							},
						},
					},
				}
				h.addNode(rng, g)
				for _, c := range g.Encoded.Nodes {
					if c.Position.Layer == 0.5 {
						cnts[c.X]++
					}
				}
			}
			Convey("The connections should be chosen randomly", func() {
				t.Log("cnt", cnts)
				So(float64(cnts[0.25])/1000.0, ShouldAlmostEqual, 0.5, 0.1)
				So(float64(cnts[0.75])/1000.0, ShouldAlmostEqual, 0.5, 0.1)
			})
		})

		Convey("When there is no connection to split", func() {
			g.Encoded.Nodes = append(g.Encoded.Nodes, evo.Node{
				Position: evo.Position{Layer: 0.5, X: 0.25},
			})
			Convey("The genome's complexity should not change", func() {
				So(g.Complexity(), ShouldEqual, 4)
			})
		})
	})
}

func TestComplexifyAddConn(t *testing.T) {
	rng := random.New()
	Convey("Given a complexify mutator and a genome", t, func() {
		h := &Complexify{}
		g := &evo.Genome{
			Encoded: evo.Substrate{
				Nodes: []evo.Node{
					{Position: evo.Position{Layer: 0.0, X: 0.0}, NeuronType: evo.Input},
					{Position: evo.Position{Layer: 0.5, X: 0.25}, NeuronType: evo.Hidden},
					{Position: evo.Position{Layer: 1.0, X: 0.5}, NeuronType: evo.Output},
				},
				Conns: []evo.Conn{
					{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 0.5, X: 0.25}},
					{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
					{Source: evo.Position{Layer: 0.5, X: 0.25}, Target: evo.Position{Layer: 1.0, X: 0.5}},
				},
			},
		}
		Convey("When the mutator allows recurrence", func() {
			h.AllowRecurrent = true
			Convey("When there is a pair of nodes to join", func() {
				// Add in the in the feed forward one so there's only the backward one left
				a := len(g.Encoded.Conns)
				h.addConn(rng, g)
				Convey("There should be a new connection", func() {
					So(len(g.Encoded.Conns), ShouldEqual, a+1)
				})
				Convey("The new connection should be enabled", func() {
					i := findConnIdx(g.Encoded.Conns,
						evo.Position{Layer: 1.0, X: 0.5},
						evo.Position{Layer: 0.5, X: 0.25})
					c := g.Encoded.Conns[i]
					So(c.Enabled, ShouldBeTrue)
				})
				Convey("Over multiple trials, the new connection should have a mean wight of 0.0 and a standard deviation of 1.0", func() {
					v := make([]float64, 1000)
					for i := 0; i < len(v); i++ {
						g.Encoded.Conns = []evo.Conn{
							{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 0.5, X: 0.25}},
							{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
							{Source: evo.Position{Layer: 0.5, X: 0.25}, Target: evo.Position{Layer: 1.0, X: 0.5}},
						}
						h.addConn(rng, g)
						j := findConnIdx(g.Encoded.Conns, evo.Position{Layer: 1.0, X: 0.5}, evo.Position{Layer: 0.5, X: 0.25})
						v[i] = g.Encoded.Conns[j].Weight
					}

					m := 0.0
					for _, x := range v {
						m += x
					}
					m /= float64(len(v))

					s := 0.0
					for _, x := range v {
						s += (x - m) * (x - m)
					}
					s = math.Sqrt(s / float64(len(v)))
					So(m, ShouldAlmostEqual, 0.0, 0.1)
					So(s, ShouldAlmostEqual, 1.0, 0.1)
				})
			})
			Convey("When there is no pair of nodes to join", func() {
				g.Encoded.Conns = append(g.Encoded.Conns, evo.Conn{
					Source: evo.Position{Layer: 0.5, X: 0.25},
					Target: evo.Position{Layer: 1.0, X: 0.5},
				})
				g.Encoded.Conns = append(g.Encoded.Conns, evo.Conn{
					Source: evo.Position{Layer: 1.0, X: 0.5},
					Target: evo.Position{Layer: 0.5, X: 0.25},
				})
				a := len(g.Encoded.Conns)
				h.addConn(rng, g)
				Convey("There should be no new connections", func() {
					So(len(g.Encoded.Conns), ShouldEqual, a)
				})
			})
		})
		Convey("When the mutator is feed-forward only", func() {
			h.AllowRecurrent = false
			i := findConnIdx(g.Encoded.Conns,
				evo.Position{Layer: 0.5, X: 0.25},
				evo.Position{Layer: 1.0, X: 0.5})
			g.Encoded.Conns = append(g.Encoded.Conns[:i], g.Encoded.Conns[i+1:]...)
			Convey("When there is a pair of nodes to join", func() {
				h.addConn(rng, g)
				Convey("There should be a new connection", func() {
					So(len(g.Encoded.Conns), ShouldEqual, 3)
				})
				Convey("The new connection should be enabled", func() {
					i := findConnIdx(g.Encoded.Conns,
						evo.Position{Layer: 0.5, X: 0.25},
						evo.Position{Layer: 1.0, X: 0.5})
					c := g.Encoded.Conns[i]
					So(c.Enabled, ShouldBeTrue)
				})
				Convey("Over multiple trials, the new connection should have a mean wight of 0.0 and a standard deviation of 1.0", func() {
					v := make([]float64, 1000)
					for i := 0; i < len(v); i++ {
						g.Encoded.Conns = []evo.Conn{
							{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 0.5, X: 0.25}},
							{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}},
						}
						h.addConn(rng, g)
						j := findConnIdx(g.Encoded.Conns,
							evo.Position{Layer: 0.5, X: 0.25},
							evo.Position{Layer: 1.0, X: 0.5})
						v[i] = g.Encoded.Conns[j].Weight
					}

					m := 0.0
					for _, x := range v {
						m += x
					}
					m /= float64(len(v))

					s := 0.0
					for _, x := range v {
						s += (x - m) * (x - m)
					}
					s = math.Sqrt(s / float64(len(v)))
					So(m, ShouldAlmostEqual, 0.0, 0.1)
					So(s, ShouldAlmostEqual, 1.0, 0.1)
				})
			})
			Convey("When there is no pair of nodes to join", func() {
				g.Encoded.Conns = append(g.Encoded.Conns, evo.Conn{
					Source: evo.Position{Layer: 0.5, X: 0.25},
					Target: evo.Position{Layer: 1.0, X: 0.5},
				})
				a := len(g.Encoded.Conns)
				h.addConn(rng, g)
				Convey("There should be no new connections", func() {
					So(len(g.Encoded.Conns), ShouldEqual, a)
				})
			})
		})
	})
}
