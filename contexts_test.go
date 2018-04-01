package evo

import (
	"context"
	"testing"
	"time"
)

func TestWithIterations(t *testing.T) {

	// Create our context and listener
	ctx, fn, cb := WithIterations(context.Background(), 2)
	defer fn()

	// Create a tick just in case our method fails
	tick := time.Tick(time.Second)

	// Run our "experiment" by making the callback
	go func(cb Callback) {
		for i := 0; i < 2; i++ {
			cb(Population{})
		}
	}(cb)

	select {
	case <-ctx.Done():
		// Success
	case <-tick:
		t.Error("context did not complete")
	}
}

func TestWithSolution(t *testing.T) {

	// Create our context and listener
	ctx, fn, cb := WithSolution(context.Background())
	defer fn()

	// Create a tick just in case our method fails
	tick := time.Tick(time.Second)

	// Run our "experiment" by making the callback
	go func(cb Callback) {
		cb(Population{
			Genomes: []Genome{{Solved: true}},
		})
	}(cb)

	select {
	case <-ctx.Done():
		// Success
	case <-tick:
		t.Error("context did not complete")
	}
}
