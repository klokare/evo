package mock

import (
	"context"
	"errors"

	"github.com/klokare/evo"
)

type Translator struct {
	Called   int
	HasError bool
}

func (t *Translator) Translate(context.Context, evo.Substrate) (net evo.Network, err error) {
	t.Called++
	if t.HasError {
		err = errors.New("mock translator error")
		return
	}
	return &Network{}, nil
}

func WithTranslator() evo.Option {
	return func(e *evo.Experiment) error {
		e.Translator = &Translator{}
		return nil
	}
}
