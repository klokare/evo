package trials

import (
	"fmt"

	"github.com/klokare/errors"
	"github.com/klokare/evo"
)

// A Runner executes multiple trials of an experiment
type Runner struct {
	Trials  int `evo:"trials"`
	Factory func() (*evo.Experiment, error)
}

// Run the trials
func (h *Runner) Run() error {

	// Iterate the trials
	if h.Trials == 0 {
		h.Trials = 1
	}
	me := new(errors.Multiple)
	for t := 0; t < h.Trials; t++ {

		// Create the new experiment
		var e *evo.Experiment
		var err error
		if e, err = h.Factory(); err != nil {
			me.Add(fmt.Errorf("trial %d: %v", t, err))
			continue
		}

		// Inform helpers of experiment and trial
		if err = h.inform(e, t); err != nil {
			me.Add(fmt.Errorf("trial %d: %v", t, err))
			continue
		}

		// Run the trial
		if err = evo.Run(e); err != nil {
			me.Add(fmt.Errorf("trial %d: %v", t, err))
		}
	}

	return me.Err()
}

// TODO: Expand to all helpers not just watchers
func (h *Runner) inform(e *evo.Experiment, t int) (err error) {

	var ws evo.Watchers
	var ok bool
	if ws, ok = e.Watcher.(evo.Watchers); !ok {
		ws = evo.Watchers{e.Watcher}
	}

	for _, w := range ws {
		if t == 0 {
			if h, ok := w.(interface {
				SetExperiment(*evo.Experiment) error
			}); ok {
				if err = h.SetExperiment(e); err != nil {
					return
				}
			}
		}
		if h, ok := w.(interface {
			SetTrial(int) error
		}); ok {
			if err = h.SetTrial(t); err != nil {
				return
			}
		}
	}
	return
}
