package evo

import "testing"

func TestNewRandom(t *testing.T) {

	// The function should return a generator of type *rand.Rand
	var rng0 = NewRandom()
	if rng0 == nil {
		t.Error("NewRandom should return a *rand.Rand")
		t.Fail()
	}

	// New random generators should provide different initial values indicating that they have different seeds
	x0 := NewRandom().Float64()
	x1 := NewRandom().Float64()
	if x0 == x1 {
		t.Errorf("different random number generators should not produce the same initial value, x0 %f and x1 %f", x0, x1)
	}

	// When overriding the seed value, the sequence of random numbers should be consistent
	SetSeed(100)
	x0 = NewRandom().Float64()
	SetSeed(100)
	x1 = NewRandom().Float64()
	if x0 != x1 {
		t.Errorf("random number generators using the same seed should produce the same sequence of values, x0 %f and x1 %f", x0, x1)
	}

}
