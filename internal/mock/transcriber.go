package mock

import (
	"errors"

	"github.com/klokare/evo"
)

// Transcriber ...
type Transcriber struct {
	Called   int
	HasError bool
}

// Transcribe ...
func (t *Transcriber) Transcribe(enc evo.Substrate) (dec evo.Substrate, err error) {
	t.Called++
	if t.HasError {
		err = errors.New("mock transcriber error")
		return
	}
	dec = enc // Just return the substrate
	return
}

// WithTranscriber ...
func WithTranscriber() evo.Option {
	return func(e *evo.Experiment) error {
		e.Transcriber = &Transcriber{}
		return nil
	}
}
