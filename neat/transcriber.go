package neat

import (
	"sort"

	"github.com/klokare/evo"
)

// Transcriber produces the decoded substrate
type Transcriber struct {
	DisableSortCheck bool
}

// Transcribe the encoded substrate into a decoded one
func (t Transcriber) Transcribe(enc evo.Substrate) (dec evo.Substrate, err error) {

	// Ensure the substrate is ordered properly
	if !t.DisableSortCheck {
		sort.Slice(enc.Nodes, func(i, j int) bool { return enc.Nodes[i].Compare(enc.Nodes[j]) < 0 })
		sort.Slice(enc.Conns, func(i, j int) bool { return enc.Conns[i].Compare(enc.Conns[j]) < 0 })
	}

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
