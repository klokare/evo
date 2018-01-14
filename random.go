package evo

import (
	"math/rand"
	"time"
)

// Seed the default generator on startup. This ensures we get a random number the first time we call
// for a new seed below in NewRandom.
func init() {
	rand.Seed(time.Now().UnixNano())
}

// Random provides the necessary functions used by this package without restricting use to the standard library's Rand
type Random interface {
	Float64() float64
	NormFloat64() float64
	Intn(n int) int
	Perm(n int) []int
}

// SetSeed reinitialises the internal random number generator's seed value. This function is not safe for concurrent calls and really only should be used to control seed values for debugging.
func SetSeed(seed int64) {
	rand.Seed(seed)
}

// NewRandom returns a new random number generator
func NewRandom() Random {
	return rand.New(rand.NewSource(rand.Int63()))
}
