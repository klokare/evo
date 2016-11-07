package bolt

import (
	"fmt"
	"sort"

	"github.com/boltdb/bolt"
	"github.com/klokare/evo/x/web"
)

// AddTrial adds the trial to the data store
func (c *Client) AddTrial(eid int, desc string) (int, error) {

	// Begin a new writable transaction
	var tx *bolt.Tx
	var err error
	if tx, err = c.db.Begin(true); err != nil {
		return -1, err
	}
	defer tx.Rollback()

	// Add the trial to the data store
	tb := tx.Bucket(trialsBucket)
	id, _ := tb.NextSequence()
	tid := int(id)

	if err = tb.Put(itob(tid), []byte(desc)); err != nil {
		return tid, err
	}

	// Attach the trial to the experiment
	if err = idxAddValue(tx, e2tBucket, eid, tid); err != nil {
		return tid, err
	}

	// Commit and return
	err = tx.Commit()
	return tid, err
}

// DelTrial removes the trial from the data store
func (c *Client) DelTrial(tid int) error {

	// Begin a new writable transaction
	var tx *bolt.Tx
	var err error
	if tx, err = c.db.Begin(true); err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete the trial
	if err = delTrial(tx, tid); err != nil {
		return err
	}

	// Save the changes
	return tx.Commit()
}

func delTrial(tx *bolt.Tx, tid int) (err error) {

	// Remove the trial
	key := itob(tid)
	if err = tx.Bucket(trialsBucket).Delete(key); err != nil {
		return err
	}

	// Remove any iterations
	if err = idxDelKey(tx, t2iBucket, tid); err != nil {
		return err
	}

	// Remove from the parent experiment
	err = idxDelValues(tx, e2tBucket, tid)
	return
}

// GetTrial returns the trial from the data store
func (c *Client) GetTrial(tid int) (t web.Trial, err error) {
	var tx *bolt.Tx
	if tx, err = c.db.Begin(false); err != nil {
		return
	}
	defer tx.Rollback()

	return getTrial(tx, tid)
}

func getTrial(tx *bolt.Tx, tid int) (t web.Trial, err error) {

	// Retrieve the trial
	v := tx.Bucket(trialsBucket).Get(itob(tid))
	if v == nil {
		err = fmt.Errorf("trial %d does not exist", tid)
		return
	}

	t.ID = tid
	t.Description = string(v)

	// Add the Iterations in reverse order
	iids := idxValues(tx, t2iBucket, tid)
	if len(iids) > 0 {
		t.Iters = make([]web.Iteration, len(iids))
		for i, iid := range iids {
			var iter web.Iteration
			if iter, err = getIteration(tx, iid); err != nil {
				return
			}
			t.Iters[i] = iter
		}
		sort.Sort(sort.Reverse(t.Iters))
	}
	return
}

func (c *Client) GetTrials(eid int) (ts []web.Trial, err error) {

	// Begin a read-only transaction
	var tx *bolt.Tx
	if tx, err = c.db.Begin(false); err != nil {
		return
	}
	defer tx.Rollback()

	// Retrieve the trials associated with this experiment
	tids := idxValues(tx, e2tBucket, eid)
	ts = make([]web.Trial, 0, len(tids))
	for _, tid := range tids {
		var t web.Trial
		if t, err = getTrial(tx, tid); err != nil {
			return
		}
		ts = append(ts, t)
	}
	return

}
