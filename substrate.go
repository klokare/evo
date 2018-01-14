package evo

import "fmt"

// A Substrate lays out a neural network on a multidimensional s
type Substrate struct {
	Nodes []Node
	Conns []Conn
}

// String returns a description of the substrate
func (s Substrate) String() string {
	return fmt.Sprintf("substrate with %d nodes, %d conns", len(s.Nodes), len(s.Conns))
}

// Complexity returns the sum of the sizes of the substrate's nodes and connections
func (s Substrate) Complexity() int { return len(s.Nodes) + len(s.Conns) }
