package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/klokare/evo"
)

// Watcher tracks iterations of the experiment and records the results at the server.
type Watcher struct {
	URL          string `evo:"web-url"`
	StudyID      int    `evo:"web-study-id"`
	ExperimentID int    `evo:"web-experiment-id"`
	TrialID      int
}

func (w Watcher) String() string {
	return fmt.Sprintf("evo.x.web.Watcher{URL: %s, StudyID: %d, ExperimentID: %d, TrialID: %d}",
		w.URL, w.StudyID, w.ExperimentID, w.TrialID)
}

// SetExperiment registers the watcher with the server under the specified experiement id. If the
// ID is 0, a new experiement will be created. If the ID is non-zero, trials will be appended to
// the specified experiement.
//
// If the description or configuration differs from the existing experiment a warning will be given
// in the console. The existing values will not be overwritten on the server.
func (w *Watcher) SetExperiment(desc string) error {
	log.Println("watcher.SetExperiment")
	// This is a new experiment
	var err error
	if w.ExperimentID == 0 {

		// Description should not be empty. Use the executable's name
		if desc == "" {
			desc = os.Args[0]
		}

		vs := url.Values{}
		vs.Set("desc", desc)
		vs.Set("sid", strconv.Itoa(w.StudyID))

		// Add the experiment on the server and update the experiment's ID
		var req *http.Request
		url := fmt.Sprintf("%s/api/experiments/add?%s", w.URL, vs.Encode())
		if req, err = http.NewRequest("POST", url, nil); err != nil {
			return err
		}

		var resp *http.Response
		if resp, err = new(http.Client).Do(req); err != nil {
			return err
		}
		defer resp.Body.Close()

		if !(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated) {
			err = fmt.Errorf("did not create experiment: status was %s", resp.Status)
			return err
		}

		var answer struct{ ExperimentID int }
		if err = json.NewDecoder(resp.Body).Decode(&answer); err != nil {
			err = fmt.Errorf("could not decode experiment id: error %s", err.Error())
			return err
		}
		w.ExperimentID = answer.ExperimentID
		return nil
	}

	// This is an existing experiment. Connect to the server and ensure the communications work.
	var resp *http.Response
	if resp, err = http.Get(fmt.Sprintf("%s/api/experiments/%d", w.URL, w.ExperimentID)); err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		msg, _ := ioutil.ReadAll(resp.Body)
		log.Println("Error in body", string(msg))
		return fmt.Errorf("could not confirm experiment %d: status was %s", w.ExperimentID, resp.Status)
	}
	return nil
}

// SetTrial informs the server a new trial is about to begin an requests a new ID.
func (w *Watcher) SetTrial() error {

	// Add the trial to the experiment on the server and update the trial ID
	var err error
	var req *http.Request

	url := fmt.Sprintf("%s/api/experiments/%d/trials/add", w.URL, w.ExperimentID)
	if req, err = http.NewRequest("POST", url, nil); err != nil {
		return err
	}

	var resp *http.Response
	if resp, err = new(http.Client).Do(req); err != nil {
		return err
	}
	defer resp.Body.Close()

	if !(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated) {
		err = fmt.Errorf("did not create trial for experiment %d: status was %s", w.ExperimentID, resp.Status)
		return err
	}

	var answer struct{ TrialID int }
	if err = json.NewDecoder(resp.Body).Decode(&answer); err != nil {
		err = fmt.Errorf("could not decode trial id for experiment %d: error %s", w.ExperimentID, err.Error())
		return err
	}
	w.TrialID = answer.TrialID
	return nil
}

// Watch records the population at the server as an iteration under the current trial.
func (w *Watcher) Watch(p evo.Population) error {

	// Encode the population
	var err error
	b := new(bytes.Buffer)
	if err = json.NewEncoder(b).Encode(p); err != nil {
		return err
	}

	// Add the iteration (aka the population) to the trial
	var req *http.Request
	url := fmt.Sprintf("%s/api/trials/%d/iterations/add", w.URL, w.TrialID)
	if req, err = http.NewRequest("POST", url, b); err != nil {
		return err
	}

	var resp *http.Response
	if resp, err = new(http.Client).Do(req); err != nil {
		return err
	}
	defer resp.Body.Close()

	if !(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated) {
		err = fmt.Errorf("did not create iteration for trial %d: status was %s", w.TrialID, resp.Status)
		return err
	}

	return nil
}
