package mock

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/klokare/evo"
)

type Configurer struct {

	// Common mock properties
	HasError bool
	Count    int

	// Experiment
	SpeciesDecayRate float64 // Increment amount [0,1], per iteration, to decay a species without improvement

	// neat.Crosser
	EnableProbability float64

	// neat.Complexity distancer
	NodesCoefficient      float64
	ActivationCoefficient float64
	ConnsCoefficient      float64

	// neat.Populator
	PopulationSize   int
	NumInputs        int
	NumOutputs       int
	NumTraits        int
	DisconnectRate   float64
	OutputActivation evo.Activation
	BiasPower        float64
	MaxBias          float64
	WeightPower      float64
	MaxWeight        float64

	// neat.Selector
	MaxStagnation               int
	MutateOnlyProbability       float64
	InterspeciesMateProbability float64

	// neat.Speciator
	CompatibilityThreshold float64 // Threshold for determining if a genome is compatible with the species
	CompatibilityModifier  float64 // Adjustment to threshold to help achieve target
	TargetSpecies          int     // The desired number of species

	// neat.mutator.Bias
	MutateBiasProbability  float64
	ReplaceBiasProbability float64
	// BiasPower              float64
	// MaxBias                float64

	// neat.mutator.Complexify
	AddNodeProbability float64
	AddConnProbability float64
	// WeightPower        float64

	// neat.mutator.Weight
	MutateWeightProbability  float64 // The probability that the connection's weight will be mutated
	ReplaceWeightProbability float64 // The probability that, if being mutated, the weight will be replaced
	// WeightPower              float64
	// MaxWeight                float64

	// neat.muator.Simplify
	// TODO: add after simplify mutator implemented

	// neat.muator.Phased
	// TODO: add after phased mutator implemented

	// neat.muator.Traits
	// TODO: add after trait mutator implemented

}

func (c *Configurer) Configure(items ...interface{}) (err error) {

	// Common mock operations
	if c.HasError {
		err = errors.New("mock configurer error")
		return
	}
	c.Count++

	// Use self as the source for configuration
	var data []byte
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(c)
	data = b.Bytes()

	for _, item := range items {
		if err = json.NewDecoder(bytes.NewBuffer(data)).Decode(item); err != nil {
			return
		}
	}
	return
}
