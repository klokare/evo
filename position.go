package evo

import "fmt"

// Position of a node on the substrate
type Position struct {
	Layer   float64
	X, Y, Z float64
}

// String returns a description of the position
func (p Position) String() string { return fmt.Sprintf("[%.4f:%.4f,%.4f,%.4f]", p.Layer, p.X, p.Y, p.Z) }

// Compare two positions for relative order on the substrate
func (p Position) Compare(other Position) int8 {
	switch {
	case p.Layer < other.Layer:
		return -1
	case p.Layer > other.Layer:
		return 1
	default:
		switch {
		case p.X < other.X:
			return -1
		case p.X > other.X:
			return 1
		default:
			switch {
			case p.Y < other.Y:
				return -1
			case p.Y > other.Y:
				return 1
			default:
				switch {
				case p.Z < other.Z:
					return -1
				case p.Z > other.Z:
					return 1
				default:
					return 0
				}
			}
		}
	}
}
