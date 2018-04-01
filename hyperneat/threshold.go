package hyperneat

import (
	"math"

	"github.com/klokare/evo"
)

// ConstantThreshold determines the weight and expression of a connection using the weight output compared
// against a constant threshold value
type ConstantThreshold struct {
	Threshold float64
}

// WeightAndExpression returns the weight and expression values for the connection at row i of the
// matrix. If expression > 0, the HyperNEAT transcriber will create the connection with the given
// weight
func (ct ConstantThreshold) WeightAndExpression(outputs evo.Matrix, i int, wp float64) (w float64, e float64) {
	w = outputs.At(i, Weight)
	a := math.Abs(w)
	if a > ct.Threshold {
		w = math.Copysign((a-ct.Threshold)/ct.Threshold, w) * wp
	}
	return
}
