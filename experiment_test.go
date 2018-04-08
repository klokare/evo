package evo

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestExperimentRunErrors(t *testing.T) {
	var cases = []struct {
		Desc     string
		HasError bool
		Options  []func(e *mockExperiment)
	}{
		{
			Desc:     "populator returns no genomes",
			HasError: true,
		},
		{
			Desc:     "populator returns genomes",
			HasError: false,
			Options: []func(*mockExperiment){
				func(e *mockExperiment) { e.mockPopulator.PopSize = 1 },
			},
		},
		{
			Desc:     "populator fails",
			HasError: true,
			Options: []func(*mockExperiment){
				func(e *mockExperiment) {
					e.mockPopulator.PopSize = 1
					e.mockPopulator.HasError = true
				},
			},
		},
		{
			Desc:     "speciator fails - initial",
			HasError: true,
			Options: []func(*mockExperiment){
				func(e *mockExperiment) {
					e.mockPopulator.PopSize = 1
					e.mockSpeciator.ErrorOn = 1
				},
			},
		},
		{
			Desc:     "selector fails",
			HasError: true,
			Options: []func(*mockExperiment){
				func(e *mockExperiment) {
					e.mockPopulator.PopSize = 1
					e.mockSelector.HasError = true
				},
			},
		},
		{
			Desc:     "crosser fails",
			HasError: true,
			Options: []func(*mockExperiment){
				func(e *mockExperiment) {
					e.mockPopulator.PopSize = 1
					e.mockCrosser.HasError = true
				},
			},
		},
		{
			Desc:     "mutator fails",
			HasError: true,
			Options: []func(*mockExperiment){
				func(e *mockExperiment) {
					e.mockPopulator.PopSize = 1
					e.mockMutator.HasError = true
				},
			},
		},
		{
			Desc:     "speciator fails - during loop",
			HasError: true,
			Options: []func(*mockExperiment){
				func(e *mockExperiment) {
					e.mockPopulator.PopSize = 1
					e.mockSpeciator.ErrorOn = 2
				},
			},
		},
		{
			Desc:     "transcriber fails",
			HasError: true,
			Options: []func(*mockExperiment){
				func(e *mockExperiment) {
					e.mockPopulator.PopSize = 1
					e.mockTranscriber.HasError = true
				},
			},
		},
		{
			Desc:     "translator fails",
			HasError: true,
			Options: []func(*mockExperiment){
				func(e *mockExperiment) {
					e.mockPopulator.PopSize = 1
					e.mockTranslator.HasError = true
				},
			},
		},
		{
			Desc:     "searcher fails",
			HasError: true,
			Options: []func(*mockExperiment){
				func(e *mockExperiment) {
					e.mockPopulator.PopSize = 1
					e.mockSearcher.HasError = true
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {
			ctx, fn := context.WithTimeout(context.Background(), time.Nanosecond*100)
			defer fn()

			// setup the experiment and run
			exp := new(mockExperiment)
			for _, f := range c.Options {
				f(exp)
			}
			_, err := Run(ctx, exp, &mockEvaluator{})

			if c.HasError {
				if err == nil {
					t.Error("expected error not found")
				}
			} else {
				if err != nil {
					t.Errorf("error not expected: %v", err)
				}
			}
		})
	}
}

func TestExperimentSubscriptions(t *testing.T) {

	const eval = 1
	const advd = 2
	const comp = 3
	const decd = 4

	var called int
	var cases = []struct {
		Desc     string
		HasError bool
		Expected int
		Subscription
	}{
		{
			Desc:     "successful decoded callback",
			HasError: false,
			Expected: decd,
			Subscription: Subscription{Event: Decoded,
				Callback: func(Population) error {
					called = decd
					return nil
				},
			},
		},
		{
			Desc:     "successful decoded callback",
			HasError: true,
			Expected: decd,
			Subscription: Subscription{Event: Decoded,
				Callback: func(Population) error {
					return errors.New("decoded callback error")
				},
			},
		},
		{
			Desc:     "successful evaluated callback",
			HasError: false,
			Expected: eval,
			Subscription: Subscription{Event: Evaluated,
				Callback: func(Population) error {
					called = eval
					return nil
				},
			},
		},
		{
			Desc:     "successful evaluated callback",
			HasError: true,
			Expected: eval,
			Subscription: Subscription{Event: Evaluated,
				Callback: func(Population) error {
					return errors.New("evaluated callback error")
				},
			},
		},
		{
			Desc:     "successful advanced callback",
			HasError: false,
			Expected: advd,
			Subscription: Subscription{Event: Advanced,
				Callback: func(Population) error {
					called = advd
					return nil
				},
			},
		},
		{
			Desc:     "successful advanced callback",
			HasError: true,
			Expected: advd,
			Subscription: Subscription{Event: Advanced,
				Callback: func(Population) error {
					return errors.New("advanced callback error")
				},
			},
		},
		{
			Desc:     "successful completed callback",
			HasError: false,
			Expected: comp,
			Subscription: Subscription{Event: Completed,
				Callback: func(Population) error {
					called = comp
					return nil
				},
			},
		},
		{
			Desc:     "successful completed callback",
			HasError: true,
			Expected: comp,
			Subscription: Subscription{Event: Completed,
				Callback: func(Population) error {
					return errors.New("completed callback error")
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Desc, func(t *testing.T) {

			// Reset the check
			called = 0

			// Create a new experiment
			exp := &mockExperiment{}
			exp.mockPopulator.PopSize = 1 // Otherwise it will fail for having 0
			exp.callbacks = []Subscription{
				c.Subscription,
			}

			// Run the experiment
			ctx, fn := context.WithTimeout(context.Background(), time.Nanosecond*100)
			defer fn()

			_, err := Run(ctx, exp, &mockEvaluator{})

			if c.HasError {
				if err == nil {
					t.Error("expected error not found")
				}
				return
			}
			if err != nil {
				t.Errorf("error not expected: %v", err)
			}

			if called != c.Expected {
				t.Errorf("incorrect called value: expected %d, actual %d", c.Expected, called)
			}
		})

	}

}

func TestExperimentRestorer(t *testing.T) {

	// Create a new experiment
	exp := &mockExperiment{}
	exp.mockPopulator.PopSize = 1   // Otherwise it will fail for having 0
	exp.mockPopulator.LastGID = 100 // as if restoring

	// Run the experiment
	ctx, fn := context.WithTimeout(context.Background(), time.Nanosecond*100)
	defer fn()

	pop, err := Run(ctx, exp, &mockEvaluator{})
	if err != nil {
		t.Errorf("error not expected: %v", err)
	}

	// Need to check in the loop because the mockSelector returns a blank genome as
	// continuing
	found := false
	for _, g := range pop.Genomes {
		if g.ID > 100 {
			found = true
			break
		}
	}
	if !found {
		t.Error("offspring do not start after last genome ID")
	}
}

type mockExperiment struct {
	mockCrosser
	mockMutator
	mockPopulator
	mockSearcher
	mockSelector
	mockSpeciator
	mockTranscriber
	mockTranslator
	callbacks []Subscription
}

func (m *mockExperiment) Subscriptions() []Subscription { return m.callbacks }

type mockCrosser struct{ Called, HasError bool }

func (m *mockCrosser) Cross(...Genome) (Genome, error) {
	var g Genome
	m.Called = true
	if m.HasError {
		return g, errors.New("error in mock crosser")
	}
	return g, nil
}

type mockTranscriber struct{ Called, HasError bool }

func (m *mockTranscriber) Transcribe(Substrate) (Substrate, error) {
	var s Substrate
	m.Called = true
	if m.HasError {
		return s, errors.New("error in mock transcriber")
	}
	return s, nil
}

type mockTranslator struct{ Called, HasError bool }

func (m *mockTranslator) Translate(Substrate) (Network, error) {
	var n Network
	m.Called = true
	if m.HasError {
		return n, errors.New("error in mock translator")
	}
	return n, nil
}

type mockEvaluator struct{ Called, HasError bool }

func (m *mockEvaluator) Evaluate(Phenome) (Result, error) {
	var r Result
	m.Called = true
	if m.HasError {
		return r, errors.New("error in mock evaluator")
	}
	return r, nil
}

type mockMutator struct{ Called, HasError bool }

func (m *mockMutator) Mutate(*Genome) error {
	m.Called = true
	if m.HasError {
		return errors.New("error in mock mutator")
	}
	return nil
}

type mockPopulator struct {
	Called, HasError bool
	PopSize          int
	LastGID          int64
}

func (m *mockPopulator) Populate() (Population, error) {
	var p Population
	m.Called = true
	if m.HasError {
		return p, errors.New("error in mock populator")
	}
	p.Genomes = make([]Genome, m.PopSize)
	for i := 0; i < len(p.Genomes); i++ {
		p.Genomes[i].ID = m.LastGID
	}
	return p, nil
}

type mockSearcher struct{ Called, HasError bool }

func (m *mockSearcher) Search(Evaluator, []Phenome) ([]Result, error) {
	var rs []Result
	m.Called = true
	if m.HasError {
		return rs, errors.New("error in mock searcher")
	}
	return rs, nil
}

type mockSelector struct{ Called, HasError bool }

func (m *mockSelector) Select(Population) ([]Genome, [][]Genome, error) {
	var cs []Genome
	var ps [][]Genome
	m.Called = true
	if m.HasError {
		return cs, ps, errors.New("error in mock selector")
	}
	cs = []Genome{{}}
	ps = [][]Genome{[]Genome{{}, {}}}
	return cs, ps, nil
}

type mockSpeciator struct{ Called, ErrorOn int }

func (m *mockSpeciator) Speciate(*Population) error {
	m.Called++
	if m.Called == m.ErrorOn {
		return errors.New("error in mock mutator")
	}
	return nil
}
