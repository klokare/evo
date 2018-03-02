package mock

import (
	"errors"

	"github.com/klokare/evo"
)

// Crosser ...
type Crosser struct {
	Called   int
	HasError bool
}

// Cross ...
func (c *Crosser) Cross(parents ...evo.Genome) (child evo.Genome, err error) {
	c.Called++
	if c.HasError {
		err = errors.New("mock crosser error")
		return
	}
	return
}

// WithCrosser ...
func WithCrosser() evo.Option {
	return func(e *evo.Experiment) error {
		e.Crosser = &Crosser{}
		return nil
	}
}
