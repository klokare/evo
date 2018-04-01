package hyperneat

import (
	"github.com/klokare/evo"
)

// LinkExpressionOutput determines the weight and expression of a connection using the LEO output
type LinkExpressionOutput struct{}

// WeightAndExpression returns the weight and expression values for the connection at row i of the
// matrix. If expression > 0, the HyperNEAT transcriber will create the connection with the given
// weight
func (LinkExpressionOutput) WeightAndExpression(outputs evo.Matrix, i int, wp float64) (w float64, e float64) {
	w = outputs.At(i, Weight) * wp
	e = outputs.At(i, LEO)
	return
}
