package bolt

import (
	"fmt"
	"sort"

	"github.com/boltdb/bolt"
	"github.com/klokare/evo/x/web"
)

func (c *Client) AddExperiment(sid int, desc string) (eid int, err error) {

	// Begin a writable transaction
	var tx *bolt.Tx
	if tx, err = c.db.Begin(true); err != nil {
		return
	}
	defer tx.Rollback()

	// Save the experiment
	id, _ := tx.Bucket(experimentsBucket).NextSequence()
	eid = int(id)
	if setExperiment(tx, sid, eid, desc); err != nil {
		return
	}

	// Commit the changes
	err = tx.Commit()
	return
}

func (c *Client) SetExperiment(sid int, eid int, desc string) (err error) {

	// Begin a writable transaction
	var tx *bolt.Tx
	if tx, err = c.db.Begin(true); err != nil {
		return
	}
	defer tx.Rollback()

	// Save the experiment
	if setExperiment(tx, sid, eid, desc); err != nil {
		return
	}

	// Commit the changes
	err = tx.Commit()
	return
}

func setExperiment(tx *bolt.Tx, sid, eid int, desc string) (err error) {

	// Store the experiment
	if err = tx.Bucket(experimentsBucket).Put(itob(eid), []byte(desc)); err != nil {
		return
	}

	// Unassociate the experiment from previous study, if any
	if err = idxDelValues(tx, s2eBucket, eid); err != nil {
		return
	}

	// Associate the experiment with the study
	if err = idxAddValue(tx, s2eBucket, sid, eid); err != nil {
		return
	}

	return
}

func (c *Client) DelExperiment(eid int) (err error) {

	// Begin a writable transaction
	var tx *bolt.Tx
	if tx, err = c.db.Begin(true); err != nil {
		return
	}
	defer tx.Rollback()

	// Delete the experiment
	if err = delExperiment(tx, eid); err != nil {
		return
	}

	// Save the changed
	return tx.Commit()
}

func delExperiment(tx *bolt.Tx, eid int) (err error) {

	// Remove the experiment
	if err = tx.Bucket(experimentsBucket).Delete(itob(eid)); err != nil {
		return
	}

	// Remove all trials
	tids := idxValues(tx, e2tBucket, eid)
	for _, tid := range tids {
		if err = delTrial(tx, tid); err != nil {
			return
		}
	}

	// Remove the experiment from the parent study
	return idxDelValues(tx, s2eBucket, eid)
}

func (c *Client) GetExperiment(eid int) (e web.Experiment, err error) {

	// Begin a read-only transaction
	var tx *bolt.Tx
	if tx, err = c.db.Begin(false); err != nil {
		return
	}
	defer tx.Rollback()

	// Return the experiment
	return getExperiment(tx, eid)
}

func getExperiment(tx *bolt.Tx, eid int) (e web.Experiment, err error) {

	// Retrieve the experiment's description
	x := tx.Bucket(experimentsBucket).Get(itob(eid))
	if x == nil {
		err = fmt.Errorf("could not get experiment %d: not found", eid)
		return
	}

	e.ID = eid
	e.Description = string(x)

	// Add the trials in reverse order
	tids := idxValues(tx, e2tBucket, eid)
	if len(tids) > 0 {
		e.Trials = make([]web.Trial, len(tids))
		for i, tid := range tids {
			var t web.Trial
			if t, err = getTrial(tx, tid); err != nil {
				return
			}
			e.Trials[i] = t
		}
		sort.Sort(sort.Reverse(e.Trials))
	}
	return
}

func (c *Client) GetExperiments(sid int) (es []web.Experiment, err error) {

	// Begin a read-only transaction
	var tx *bolt.Tx
	if tx, err = c.db.Begin(false); err != nil {
		return
	}
	defer tx.Rollback()

	// Retrieve the experiments
	eids := idxValues(tx, s2eBucket, sid)
	if len(eids) > 0 {
		es = make([]web.Experiment, 0, len(eids))
		for _, eid := range eids {
			var e web.Experiment
			if e, err = getExperiment(tx, eid); err != nil {
				return
			}
			es = append(es, e)
		}
	}
	return
}
