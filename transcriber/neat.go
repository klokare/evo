package transcriber

import (
	"sort"

	"github.com/klokare/evo"
)

type NEAT struct{}

func (h *NEAT) Transcribe(enc evo.Substrate) (evo.Substrate, error) {

	// Create the decoded substrate
	dec := evo.Substrate{
		Nodes: make([]evo.Node, len(enc.Nodes)),
		Conns: make([]evo.Conn, 0, len(enc.Conns)),
	}
	copy(dec.Nodes, enc.Nodes)
	for _, c := range enc.Conns {
		if c.Enabled {
			dec.Conns = append(dec.Conns, c)
		}
	}

	// Sort the decoded substrate
	sort.Sort(dec.Nodes)
	sort.Sort(dec.Conns)

	return dec, nil
}
