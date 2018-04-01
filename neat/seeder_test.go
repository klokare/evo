package neat

import (
	"math"
	"testing"

	"github.com/klokare/evo"
	"github.com/klokare/evo/internal/mock"
)

// Tests the population creation. Testing the specific genome creation is done below.
func TestSeederSeed(t *testing.T) {

	var cases = []struct {
		Desc             string
		In, Out, Traits  int
		DisconnectRate   float64
		HasError         bool
		OutputActivation evo.Activation
	}{
		{
			Desc:             "negative number of inputs",
			In:               -1,
			HasError:         true,
			OutputActivation: evo.Sigmoid,
		},
		{
			Desc:             "zero inputs",
			In:               0,
			HasError:         true,
			OutputActivation: evo.Sigmoid,
		},
		{
			Desc:             "negative number of outputs",
			In:               1,
			Out:              -1,
			HasError:         true,
			OutputActivation: evo.Sigmoid,
		},
		{
			Desc:             "zero outputs",
			In:               1,
			Out:              0,
			HasError:         true,
			OutputActivation: evo.Sigmoid,
		},
		{
			Desc:             "normal population",
			In:               2,
			Out:              2,
			Traits:           2,
			DisconnectRate:   0.5,
			HasError:         false,
			OutputActivation: evo.Sigmoid,
		},
	}

	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {

			// Create the seeder
			s := &Seeder{
				NumInputs:        c.In,
				NumOutputs:       c.Out,
				NumTraits:        c.Traits,
				DisconnectRate:   c.DisconnectRate,
				OutputActivation: c.OutputActivation,
			}

			// Create the seed genome
			g, err := s.Seed()

			// Check error
			if !t.Run("error", mock.Error(c.HasError, err)) {
				t.FailNow()
			}
			if c.HasError {
				return // Expected error so stop checking other things
			}

			// The genome should have the correct number of input nodes
			cnt := 0
			for _, n := range g.Encoded.Nodes {
				if n.Neuron == evo.Input {
					cnt++
				}
			}
			if cnt != c.In {
				t.Errorf("incorrect number of input nodes")
			}

			// The genome should have the correct number of output nodes and their activations should also be correct
			cnt = 0
			ac := true
			for _, n := range g.Encoded.Nodes {
				if n.Neuron == evo.Output {
					cnt++
					ac = ac && n.Activation == c.OutputActivation
				}
			}
			if cnt != c.In {
				t.Errorf("incorrect number of output nodes")
			}
			if !ac {
				t.Error("not all output nodes have the correct activation")
			}

		})
	}
}

func TestSeederDisconnectRate(t *testing.T) {

	var cases = []struct {
		Desc           string
		DisconnectRate float64
	}{
		{Desc: "fully connected", DisconnectRate: 0.0},
		{Desc: "partially connected", DisconnectRate: 0.5},
		{Desc: "unconnected", DisconnectRate: 1.0},
	}

	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {

			// Create some genomes
			s := &Seeder{
				NumInputs:      5,
				NumOutputs:     2,
				NumTraits:      2,
				DisconnectRate: c.DisconnectRate,
			}

			sum := 0.0
			for i := 0; i < 10000; i++ {
				g, _ := s.Seed()
				sum += float64(len(g.Encoded.Conns)) / 10.0
			}
			avg := sum / 10000.0 // avg connected rate
			avg = 1.0 - avg      // avg disconnected rate
			if math.Abs(c.DisconnectRate-avg) > 0.1 {
				t.Errorf("incorrect disconnect rate: expected %f, actual: %f", c.DisconnectRate, avg)
			}
		})
	}

}
