package evo

import (
	"context"
	"sync/atomic"
)

// WithIterations creates a cancelable context and return the cancel function and a callback which
// must be subscribed in the experiment. The context will be cancelled when the number of iteations
// has been reached
func WithIterations(ctx context.Context, n int) (context.Context, context.CancelFunc, Callback) {
	var fn context.CancelFunc
	completed := new(int64)
	ctx, fn = context.WithCancel(ctx)
	return ctx, fn, func(Population) error {
		if atomic.AddInt64(completed, 1) >= int64(n) {
			fn() // cancel the context
		}
		return nil
	}
}

// WithSolution creates a cancelable context and return the cancel function and a callback which
// must be subscribed in the experiment. The context will be cancelled when a solution has been
// found.
func WithSolution(ctx context.Context) (context.Context, context.CancelFunc, Callback) {
	var fn context.CancelFunc
	ctx, fn = context.WithCancel(ctx)
	return ctx, fn, func(pop Population) error {
		for _, g := range pop.Genomes {
			if g.Solved {
				fn()
				break
			}
		}
		return nil
	}
}
