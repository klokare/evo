package neat

import (
	"github.com/klokare/evo"
	"github.com/klokare/evo/comparer"
	"github.com/klokare/evo/crosser"
	"github.com/klokare/evo/mutator"
	"github.com/klokare/evo/searcher"
	"github.com/klokare/evo/selector"
	"github.com/klokare/evo/speciator"
	"github.com/klokare/evo/translator"
	"github.com/klokare/evo/watcher"
	"github.com/klokare/evo/x/web"
)

// NewExperiment creates a new trial using the helpers associated with NEAT
func NewExperiment(c evo.Configurer, eval func() (evo.Evaluator, error), options ...func(*evo.Trial) error) (*evo.Experiment, error) {

	// Create and configure the experiment
	exp := new(evo.Experiment)
	if err := c.Configure(exp); err != nil {
		return nil, err
	}

	// Attempt a connection to the web application
	ww := new(web.Watcher)
	if err := c.Configure(ww); err != nil {
		return nil, err
	}
	if ww.URL != "" {

		// Register the experiment with the web application
		if err := ww.SetExperiment(exp.Description); err != nil {
			return nil, err
		}

		// Add code to attach the web watcher
		options = append(options, func(t *evo.Trial) error {

			// Append the web watcher to the list of watchers
			var ws evo.Watchers
			var ok bool
			if ws, ok = t.Watcher.(evo.Watchers); !ok {
				ws = []evo.Watcher{t.Watcher}
			}
			t.Watcher = append(ws, ww)

			// Register the trial with the web application, if applicable
			if ww != nil {
				if err := ww.SetTrial(); err != nil {
					return err
				}
			}
			return nil
		})
	} else {
		ww = nil // get rid of the web watcher
	}

	// Create the function that generates new trials
	exp.NewTrial = func(i int) (*evo.Trial, error) {

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

		pop := new(Seed)
		if err = c.Configure(pop); err != nil {
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
		if sch.Evaluator, err = eval(); err != nil {
			return nil, err
		}

		sel := new(selector.Generational)
		if err = c.Configure(sel); err != nil {
			return nil, err
		}

		tcr := new(Transcriber)
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
		t := new(evo.Trial)
		if err = c.Configure(t); err != nil {
			return nil, err
		}
		t.Crosser = crs
		t.Mutator = mut
		t.Populater = pop
		t.Searcher = sch
		t.Selector = sel
		t.Speciater = spc
		t.Transcriber = tcr
		t.Translator = tsl
		t.Watcher = wtc

		// Apply the options, if any, and return
		for _, option := range options {
			if err = option(t); err != nil {
				return nil, err
			}
		}
		return t, nil
	}

	return exp, nil
}
