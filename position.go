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

// Midpoint returns the midpoint (between to positions) or centre (between three or more positions).
// This function is necessary because of how Go treats floating point [constants](https://www.ardanlabs.com/blog/2014/04/introduction-to-numeric-constants-in-go.html) and how that can lead to [differeing values](https://github.com/golang/go/issues/23876) for similar calculations.
func Midpoint(positions ...Position) (m Position) {
	if len(positions) == 1 {
		return positions[0]
	}
	var l, x, y, z float64
	for _, p := range positions {
		l += p.Layer
		x += p.X
		y += p.Y
		z += p.Z
	}
	n := float64(len(positions))
	m = Position{
		Layer: l / n,
		X:     x / n,
		Y:     y / n,
		Z:     y / n,
	}
	return
}
