package evo

import (
	"github.com/klokare/evo"
	"github.com/klokare/evo/comparer"
	"github.com/klokare/evo/crosser"
	"github.com/klokare/evo/mutator"
	"github.com/klokare/evo/searcher"
	"github.com/klokare/evo/selector"
	"github.com/klokare/evo/speciator"
	"github.com/klokare/evo/transcriber"
	"github.com/klokare/evo/translator"
	"github.com/klokare/evo/watcher"
)

// NewExperiment creates a new experiment using the helpers associated with evo
func NewExperiment(c evo.Configurer, e evo.Evaluator, options ...func(*evo.Experiment) error) (*evo.Experiment, error) {

	// Create and configure the helpers
	var err error
	cmp := new(comparer.Distance)
	if err = c.Configure(cmp); err != nil {
		return nil, err
	}

	sp1 := new(speciator.Static)
	if err = c.Configure(sp1); err != nil {
		return nil, err
	}
	sp1.Comparer = cmp
	spc := new(speciator.Dynamic)
	if err = c.Configure(spc); err != nil {
		return nil, err
	}
	spc.Static = *sp1

	crs := new(crosser.Multiple)
	if err = c.Configure(crs); err != nil {
		return nil, err
	}

	var mut evo.Mutators = make([]evo.Mutator, 2)
	mu1 := new(mutator.Complexify)
	if err = c.Configure(mu1); err != nil {
		return nil, err
	}
	mu2 := new(mutator.Weight)
	if err = c.Configure(mu2); err != nil {
		return nil, err
	}
	mut[0] = mu1
	mut[1] = mu2

	sch := new(searcher.Serial)
	if err = c.Configure(sch); err != nil {
		return nil, err
	}
	sch.Evaluator = e

	sel := new(selector.Generational)
	if err = c.Configure(sel); err != nil {
		return nil, err
	}

	tcr := new(transcriber.NEAT)
	if err = c.Configure(tcr); err != nil {
		return nil, err
	}
	tsl := new(translator.Simple)
	if err = c.Configure(tsl); err != nil {
		return nil, err
	}

	var wtc evo.Watchers
	con := new(watcher.Console)
	if err = c.Configure(tsl); err != nil {
		return nil, err
	}
	wtc = append(wtc, con)
	wtc = append(wtc, sel)

	// Create the experiment
	exp := &evo.Experiment{
		Crosser:     crs,
		Mutator:     mut,
		Searcher:    sch,
		Selector:    sel,
		Speciater:   spc,
		Transcriber: tcr,
		Translator:  tsl,
		Watcher:     wtc,
	}

	// Apply the options, if any, and return
	for _, option := range options {
		if err = option(exp); err != nil {
			return nil, err
		}
	}
	return exp, nil
}
