package bolt

import (
	"fmt"
	"sort"

	"github.com/boltdb/bolt"
	"github.com/klokare/evo/x/web"
)

func (c *Client) AddStudy(desc string) (sid int, err error) {

	// Begin a writable transaction
	var tx *bolt.Tx
	if tx, err = c.db.Begin(true); err != nil {
		return
	}
	defer tx.Rollback()

	// Save the study
	id, _ := tx.Bucket(studiesBucket).NextSequence()
	sid = int(id)
	if setStudy(tx, sid, desc); err != nil {
		return
	}

	// Commit the changes
	err = tx.Commit()
	return
}

func (c *Client) SetStudy(sid int, desc string) (err error) {

	// Begin a writable transaction
	var tx *bolt.Tx
	if tx, err = c.db.Begin(true); err != nil {
		return
	}
	defer tx.Rollback()

	// Save the study
	if setStudy(tx, sid, desc); err != nil {
		return
	}

	// Commit the changes
	err = tx.Commit()
	return
}

func setStudy(tx *bolt.Tx, sid int, desc string) (err error) {
	return tx.Bucket(studiesBucket).Put(itob(sid), []byte(desc))
}

func (c *Client) DelStudy(sid int) (err error) {

	// Begin a writable transaction
	var tx *bolt.Tx
	if tx, err = c.db.Begin(true); err != nil {
		return
	}
	defer tx.Rollback()

	// Delete the study
	if err = tx.Bucket(studiesBucket).Delete(itob(sid)); err != nil {
		return
	}

	// Delete experiments
	eids := idxValues(tx, s2eBucket, sid)
	for _, eid := range eids {
		if err = delExperiment(tx, eid); err != nil {
			return
		}
	}
	// Commit the changes
	err = tx.Commit()
	return
}

func (c *Client) GetStudy(sid int) (s web.Study, err error) {

	// Begin a read-only transaction
	var tx *bolt.Tx
	if tx, err = c.db.Begin(false); err != nil {
		return
	}
	defer tx.Rollback()

	return getStudy(tx, sid)
}

func getStudy(tx *bolt.Tx, sid int) (s web.Study, err error) {

	// Retrieve the study's description
	x := tx.Bucket(studiesBucket).Get(itob(sid))
	if x == nil {
		err = fmt.Errorf("could not get study %d: not found", sid)
		return
	}

	s.ID = sid
	s.Description = string(x)

	// Add the experiments in reverse order
	eids := idxValues(tx, s2eBucket, sid)
	if len(eids) > 0 {
		s.Experiments = make([]web.Experiment, len(eids))
		for i, eid := range eids {
			var e web.Experiment
			if e, err = getExperiment(tx, eid); err != nil {
				return
			}
			s.Experiments[i] = e
		}
		sort.Sort(sort.Reverse(s.Experiments))
	}
	return
}

func (c *Client) GetStudies() (ss []web.Study, err error) {

	// Begin a read-only transaction
	var tx *bolt.Tx
	if tx, err = c.db.Begin(false); err != nil {
		return
	}
	defer tx.Rollback()

	// Retrieve the studies
	ss = make([]web.Study, 0, 100)
	sc := tx.Bucket(studiesBucket).Cursor()
	for k, _ := sc.First(); k != nil; k, _ = sc.Next() {
		var s web.Study
		if s, err = getStudy(tx, btoi(k)); err != nil {
			return
		}
		ss = append(ss, s)
	}
	sort.Sort(sort.Reverse(web.Studies(ss)))
	return
}
