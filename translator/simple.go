package translator

import "github.com/klokare/evo"

// A Simple tanslator that creates a network from a substrate
type Simple struct{}

func (h Simple) String() string { return "evo.translator.Simple{}" }

// Translate a substrate into a network
func (h *Simple) Translate(sub evo.Substrate) (evo.Network, error) {
	return NewNetwork(sub), nil
}
