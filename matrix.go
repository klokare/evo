package evo

// Matrix descibes data organised as a matrix. It mimics a subset of the signature of [gonum's mat.Matrix](https://godoc.org/gonum.org/v1/gonum/mat) which allows directly passing matrices from that package as inputs as well as any other type that implements it, such as [sparse](https://godoc.org/github.com/james-bowman/sparse). Network implementations, however, may expect a specific type and throw an error if they cannot convert to the desired type.
//
// If you want to use single-row, float64 slices instead, please look at the ToMatrix and
// FromMatrix functions in the [klokare/network](https://github.com/klokare/network) package.
type Matrix interface {

	// Dims returns the dimensions of a Matrix.
	Dims() (r, c int)

	// At returns the value of a matrix element at row i, column j.
	// It will panic if i or j are out of bounds for the matrix.
	At(i, j int) float64
}
