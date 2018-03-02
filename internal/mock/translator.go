package mock

import (
	"errors"

	"github.com/klokare/evo"
)

// Translator ...
type Translator struct {
	Called   int
	HasError bool
}

// Translate ...
func (t *Translator) Translate(evo.Substrate) (net evo.Network, err error) {
	t.Called++
	if t.HasError {
		err = errors.New("mock translator error")
		return
	}
	return &Network{}, nil
}

// WithTranslator ...
func WithTranslator() evo.Option {
	return func(e *evo.Experiment) error {
		e.Translator = &Translator{}
		return nil
	}
}
