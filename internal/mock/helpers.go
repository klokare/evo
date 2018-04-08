package mock

import (
	"errors"

	"github.com/klokare/evo"
)

type Experiment struct {
	Crosser
	Mutator
	Populator
	Searcher
	Selector
	Speciator
	Transcriber
	Translator
	callbacks []evo.Subscription
}

func (m *Experiment) Subscribers() []evo.Subscription { return m.callbacks }

type Crosser struct{ Called, HasError bool }

func (m *Crosser) Cross(...evo.Genome) (evo.Genome, error) {
	var g evo.Genome
	m.Called = true
	if m.HasError {
		return g, errors.New("error in  crosser")
	}
	return g, nil
}

type Transcriber struct{ Called, HasError bool }

func (m *Transcriber) Transcribe(evo.Substrate) (evo.Substrate, error) {
	var s evo.Substrate
	m.Called = true
	if m.HasError {
		return s, errors.New("error in mock transcriber")
	}
	return s, nil
}

type Translator struct{ Called, HasError bool }

func (m *Translator) Translate(evo.Substrate) (evo.Network, error) {
	var n evo.Network
	m.Called = true
	if m.HasError {
		return n, errors.New("error in mock translator")
	}
	return n, nil
}

type Evaluator struct{ Called, HasError bool }

func (m *Evaluator) Evaluate(p evo.Phenome) (evo.Result, error) {
	r := evo.Result{ID: p.ID}
	m.Called = true
	if m.HasError {
		return r, errors.New("error in  evaluator")
	}
	return r, nil
}

type Mutator struct{ Called, HasError bool }

func (m *Mutator) Mutate(*evo.Genome) error {
	m.Called = true
	if m.HasError {
		return errors.New("error in  mutator")
	}
	return nil
}

type Populator struct {
	Called, HasError bool
	PopSize          int
	LastGID          int64
}

func (m *Populator) Populate() (evo.Population, error) {
	var p evo.Population
	m.Called = true
	if m.HasError {
		return p, errors.New("error in  populator")
	}
	p.Genomes = make([]evo.Genome, m.PopSize)
	for i := 0; i < len(p.Genomes); i++ {
		p.Genomes[i].ID = m.LastGID
	}
	return p, nil
}

type Searcher struct{ Called, HasError bool }

func (m *Searcher) Search(evo.Evaluator, []evo.Phenome) ([]evo.Result, error) {
	var rs []evo.Result
	m.Called = true
	if m.HasError {
		return rs, errors.New("error in  searcher")
	}
	return rs, nil
}

type Seeder struct{ Called, HasError bool }

func (m *Seeder) Seed() (evo.Genome, error) {
	g := evo.Genome{
		Encoded: evo.Substrate{
			Nodes: []evo.Node{
				{Position: evo.Position{Layer: 0.0, X: 0.0}, Neuron: evo.Input, Activation: evo.Direct},
				{Position: evo.Position{Layer: 0.0, X: 1.0}, Neuron: evo.Input, Activation: evo.Direct},
				{Position: evo.Position{Layer: 1.0, X: 0.5}, Neuron: evo.Output, Activation: evo.Sigmoid},
			},
			Conns: []evo.Conn{
				{Source: evo.Position{Layer: 0.0, X: 0.0}, Target: evo.Position{Layer: 1.0, X: 0.5}, Enabled: true},
				{Source: evo.Position{Layer: 0.0, X: 1.0}, Target: evo.Position{Layer: 1.0, X: 0.5}, Enabled: true},
			},
		},
		Traits: []float64{0.25, 0.75},
	}
	m.Called = true
	if m.HasError {
		return g, errors.New("error in  searcher")
	}
	return g, nil
}

type Selector struct{ Called, HasError bool }

func (m *Selector) Select(evo.Population) ([]evo.Genome, [][]evo.Genome, error) {
	var cs []evo.Genome
	var ps [][]evo.Genome
	m.Called = true
	if m.HasError {
		return cs, ps, errors.New("error in  selector")
	}
	cs = []evo.Genome{{}}
	ps = [][]evo.Genome{[]evo.Genome{{}, {}}}
	return cs, ps, nil
}

type Speciator struct{ Called, ErrorOn int }

func (m *Speciator) Speciate(*evo.Population) error {
	m.Called++
	if m.Called == m.ErrorOn {
		return errors.New("error in  mutator")
	}
	return nil
}
