package neat

import (
	"context"

	"github.com/klokare/evo"
)

// Transcriber provides the nodes and connections to be used in network translation
type Transcriber struct{}

// Transcribe the genome into the nodes and connections used in network translation
func (t Transcriber) Transcribe(ctx context.Context, enc evo.Substrate) (dec evo.Substrate, err error) {

	// Copy the nodes
	dec.Nodes = make([]evo.Node, len(enc.Nodes))
	copy(dec.Nodes, enc.Nodes)

	// Copy the enabled conns
	dec.Conns = make([]evo.Conn, 0, len(enc.Conns))
	for _, c := range enc.Conns {
		if c.Enabled {
			dec.Conns = append(dec.Conns, c)
		}
	}
	return
}

// WithTranscriber sets the experiment's transcriber to a configured NEAT transcriber
func WithTranscriber(evo.Configurer) evo.Option {
	return func(e *evo.Experiment) (err error) {
		e.Transcriber = new(Transcriber)
		return
	}
}
