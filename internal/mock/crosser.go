package mock

import (
	"context"
	"errors"

	"github.com/klokare/evo"
)

type Crosser struct {
	Called   int
	HasError bool
}

func (c *Crosser) Cross(ctx context.Context, parents ...evo.Genome) (child evo.Genome, err error) {
	c.Called++
	if c.HasError {
		err = errors.New("mock crosser error")
		return
	}
	return
}

func WithCrosser() evo.Option {
	return func(e *evo.Experiment) error {
		e.Crosser = &Crosser{}
		return nil
	}
}
