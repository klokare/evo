package mock

import (
	"context"
	"errors"

	"github.com/klokare/evo"
)

type Transcriber struct {
	Called   int
	HasError bool
}

func (t *Transcriber) Transcribe(ctx context.Context, enc evo.Substrate) (dec evo.Substrate, err error) {
	t.Called++
	if t.HasError {
		err = errors.New("mock transcriber error")
		return
	}
	dec = enc // Just return the substrate
	return
}

func WithTranscriber() evo.Option {
	return func(e *evo.Experiment) error {
		e.Transcriber = &Transcriber{}
		return nil
	}
}
