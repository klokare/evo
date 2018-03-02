package mock

import "github.com/klokare/evo"

// Network ...
type Network struct{}

// Activate ...
func (n *Network) Activate(evo.Matrix) (evo.Matrix, error) {
	return nil, nil
}
