package hyperneat

import (
	"errors"
	"sort"

	"github.com/klokare/evo"
	"gonum.org/v1/gonum/mat"
)

// Known errors
var (
	ErrMissingCppnTranscriber = errors.New("hyperneat transcriber requires a cppn transcriber")
	ErrMissingCppnTranslator  = errors.New("hyperneat transcriber requires a cppn translator")
	ErrMissingInspector       = errors.New("hyperneat transcriber requires an inspector")
	ErrTemplateNotSet         = errors.New("hyperneat template not set")
	ErrNoInputNodes           = errors.New("hyperneat template has no input nodes")
	ErrNoOutputNodes          = errors.New("hyperneat template has no output nodes")
	ErrHasExistingConns       = errors.New("hyperneat template should not have existing connections")
	ErrInvalidWeightPower     = errors.New("hyperneat weight power cannot be zero")
)

// Inspector describes a helper that provides the weight and expression values for a connection
type Inspector interface {
	WeightAndExpression(outputs evo.Matrix, idx int, weightPower float64) (w float64, e float64)
}

// Transcriber decodes an encoded substrate by processing it through a Cppn
type Transcriber struct {

	// Helpers
	CppnTranscriber evo.Transcriber
	CppnTranslator  evo.Translator
	Inspector       // used to determine the weight of a connection and if it will be expressed

	// Properties
	WeightPower      float64 // Power by which to multiply the weight's output value
	BiasPower        float64 // Power by which to multiply the weight's output value
	DisableSortCheck bool    // Speed things up by disabling sort check on encoded substrate. Use only if sure the incoming substrate is already sorted.

	// Internal structure
	layers                  map[float64][]evo.Node
	inputs, outputs, hidden int
	keys                    []float64
}

// SetTemplate creates the inner structure of the transcriber using the supplied template.
func (t *Transcriber) SetTemplate(tmpl evo.Substrate) (err error) {

	// Abort if there are existing connections
	if len(tmpl.Conns) > 0 {
		err = ErrHasExistingConns
		return
	}

	// Ensure the template is sorted
	if !t.DisableSortCheck {
		sort.Slice(tmpl.Nodes, func(i, j int) bool { return tmpl.Nodes[i].Compare(tmpl.Nodes[j]) < 0 })
	}

	// Separate the nodes into layers. Template is already sorted
	t.layers = make(map[float64][]evo.Node, 5)
	t.keys = make([]float64, 0, 5)
	for _, n := range tmpl.Nodes {
		var nodes []evo.Node
		var ok bool
		if nodes, ok = t.layers[n.Layer]; !ok {
			t.keys = append(t.keys, n.Layer)
			nodes = make([]evo.Node, 0, len(tmpl.Nodes))
		}
		nodes = append(nodes, n)
		t.layers[n.Layer] = nodes
		switch n.Neuron {
		case evo.Input:
			t.inputs++
		case evo.Hidden:
			t.hidden++
		case evo.Output:
			t.outputs++
		}
	}
	sort.Float64s(t.keys)

	// Check for errors
	if t.inputs == 0 {
		err = ErrNoInputNodes
		return
	} else if t.outputs == 0 {
		err = ErrNoOutputNodes
		return
	}
	return
}

// Transcribe creates a cppn network from the encoded substrate and then decodes it
func (t Transcriber) Transcribe(enc evo.Substrate) (dec evo.Substrate, err error) {

	// Check for known errors
	if len(t.keys) == 0 {
		err = ErrTemplateNotSet
		return
	} else if t.WeightPower == 0.0 {
		err = ErrInvalidWeightPower
		return
	} else if t.CppnTranscriber == nil {
		err = ErrMissingCppnTranscriber
		return
	} else if t.CppnTranslator == nil {
		err = ErrMissingCppnTranslator
		return
	} else if t.Inspector == nil {
		err = ErrMissingInspector
		return
	}

	// Transcribe the encoded substrate
	if dec, err = t.CppnTranscriber.Transcribe(enc); err != nil {
		return
	}

	// Translate into a Cppn
	var net evo.Network
	if net, err = t.CppnTranslator.Translate(dec); err != nil {
		return
	}

	// Initialise the decoded substrate
	dec = evo.Substrate{
		Nodes: make([]evo.Node, 0, t.inputs+t.hidden+t.outputs),
	}
	for _, layer := range t.layers {
		dec.Nodes = append(dec.Nodes, layer...)
	}

	// Connect each layer
	var inputs *mat.Dense
	var outputs evo.Matrix
	row := make([]float64, 8)
	for i := 0; i < len(t.keys)-1; i++ {
		snodes := t.layers[t.keys[i]]
		tnodes := t.layers[t.keys[i+1]]
		inputs = mat.NewDense(len(snodes)*len(tnodes), 8, nil)
		r := 0
		for _, src := range snodes {
			row[0], row[1], row[2], row[3] = src.Layer, src.X, src.Y, src.Z
			for _, tgt := range tnodes {
				row[4], row[5], row[6], row[7] = tgt.Layer, tgt.X, tgt.Y, tgt.Z
				inputs.SetRow(r, row)
				r++
			}
		}

		if outputs, err = net.Activate(inputs); err != nil {
			return
		}
		r = 0
		for _, sn := range snodes {
			for _, tn := range tnodes {
				w, e := t.WeightAndExpression(outputs, r, t.WeightPower)
				if e > 0 {
					dec.Conns = append(dec.Conns, evo.Conn{
						Source:  sn.Position,
						Target:  tn.Position,
						Weight:  w,
						Enabled: true,
					})
				}
				r++
			}
		}
	}

	// Set bias values for hidden and output nodes
	row[4], row[5], row[6], row[7] = 0.0, 0.0, 0.0, 0.0
	inputs = mat.NewDense(t.hidden+t.outputs, 8, nil)
	r := 0
	for _, n := range dec.Nodes {
		if n.Neuron == evo.Input {
			continue
		}
		row[0], row[1], row[2], row[3] = n.Layer, n.X, n.Y, n.Z
		inputs.SetRow(r, row)
		r++
	}
	if outputs, err = net.Activate(inputs); err != nil {
		return
	}
	for i := t.inputs; i < len(dec.Nodes); i++ {
		dec.Nodes[i].Bias = outputs.At(i-t.inputs, Bias) * t.BiasPower
	}

	// Return the decoded substrate
	sort.Slice(dec.Conns, func(i, j int) bool { return dec.Conns[i].Compare(dec.Conns[j]) < 0 })
	return
}
