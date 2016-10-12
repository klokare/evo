package translator

import "github.com/klokare/evo"

// A Simple tanslator that creates a network from a substrate
type Simple struct{}

// Translate a substrate into a network
func (h *Simple) Translate(sub evo.Substrate) (evo.Network, error) {
	return NewNetwork(sub), nil
}
