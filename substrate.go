package evo

import "fmt"

// A Substrate is the structural representation of a neural network
type Substrate struct {
	Nodes Nodes
	Conns Conns
}

// Complexity returns the number of nodes and connections in the substrate
func (s Substrate) Complexity() int { return len(s.Nodes) + len(s.Conns) }

// A Node is the definition of a neuron located at a particular position
// NOTE: No separate node ID is used as this library assumes only 1 node can exist at a single
// position.
type Node struct {
	Position
	NeuronType
	ActivationType
}

// String provides a more readable description
func (a Node) String() string {
	return fmt.Sprintf("Node [%f,%f,%f,%f] %v %v", a.Position.Layer, a.Position.X,
		a.Position.Y, a.Position.Z, a.NeuronType, a.ActivationType)
}

// Compare two nodes for relative position
func (a Node) Compare(b Node) int { return a.Position.Compare(b.Position) }

// Nodes is a sortable collection of nodes
type Nodes []Node

func (n Nodes) Len() int           { return len(n) }
func (n Nodes) Less(i, j int) bool { return n[i].Compare(n[j]) < 0 }
func (n Nodes) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }

// A Conn is the defintion of a synapse connecting neurons a positions source and target
// NOTE: No separate innovation number is used as the library does not reset innovations after each
// generation and the key for checking innovation would be Source and Target.
type Conn struct {
	Source  Position
	Target  Position
	Weight  float64
	Enabled bool
}

// Compare two connections for relative position
func (a Conn) Compare(b Conn) int {
	s := a.Source.Compare(b.Source)
	if s == 0 {
		return a.Target.Compare(b.Target)
	}
	return s
}

// Conns is a sortable collection of connections
type Conns []Conn

func (c Conns) Len() int           { return len(c) }
func (c Conns) Less(i, j int) bool { return c[i].Compare(c[j]) < 0 }
func (c Conns) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

// A Position it the 4-dimensional location of a node on the substrate
type Position struct {
	Layer   float64
	X, Y, Z float64
}

// Compare two positions for the relative order. 0 means the the two positions are equal; -1, means
// position a is less than b; 1, a is greater than b
func (a Position) Compare(b Position) int {
	switch {
	case a.Layer < b.Layer:
		return -1
	case a.Layer > b.Layer:
		return 1
	default:
		switch {
		case a.X < b.X:
			return -1
		case a.X > b.X:
			return 1
		default:
			switch {
			case a.Y < b.Y:
				return -1
			case a.Y > b.Y:
				return 1
			default:
				switch {
				case a.Z < b.Z:
					return -1
				case a.Z > b.Z:
					return 1
				default:
					return 0
				}
			}
		}
	}
}
