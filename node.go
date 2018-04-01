package evo

import "fmt"

// A Node describes a neuron in the network
type Node struct {
	Position           // The location of the node on the substrate
	Neuron             // Then neuron type
	Activation         // The activation type
	Bias       float64 // Bias value for the neuron
	Locked     bool    // Locked nodes cannot be removed
}

// String reutnrs the description of the node
func (n Node) String() string {
	return fmt.Sprintf("%s node %s - %s bias %.4f", n.Neuron, n.Position, n.Activation, n.Bias)
}

// Compare two nodes for relative positions on the substrate
func (n Node) Compare(other Node) int8 { return n.Position.Compare(other.Position) }
