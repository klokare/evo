package evo

import "fmt"

// A Conn describes a synapse (a connection between neurons) in the network
type Conn struct {
	Source, Target Position // The positions of the source and target nodes
	Weight         float64  // The connection weight
	Enabled        bool     // True if this connection should be used to create a synapse
}

// String returns a description of the connection
func (c Conn) String() string {
	var s string
	if c.Enabled {
		s = "enabled "
	} else {
		s = "disabled"
	}
	return fmt.Sprintf("conn %v->%v %s %.6f", c.Source, c.Target, s, c.Weight)
}

// Compare two connections for relative position, considering their source and
// target node positions.
func (c Conn) Compare(other Conn) int8 {
	switch c.Source.Compare(other.Source) {
	case -1:
		return -1
	case 1:
		return 1
	default:
		return c.Target.Compare(other.Target)
	}
}
