package bolt

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/boltdb/bolt"
	"github.com/klokare/evo"
	"github.com/klokare/evo/x/web"
)

func (c *Client) AddIteration(tid int, p evo.Population) (int, error) {

	// Add the iteration to the data store
	var tx *bolt.Tx
	var err error
	if tx, err = c.db.Begin(true); err != nil {
		return -1, err
	}
	defer tx.Rollback()

	// Create the iteration object
	ib := tx.Bucket(iterationsBucket)
	id, _ := ib.NextSequence()

	sort.Sort(sort.Reverse(p.Genomes))
	i := web.Iteration{
		ID:         int(id),
		Updated:    time.Now(),
		Best:       p.Genomes[0],
		Solved:     p.Genomes[0].Solved,
		Fitness:    make([]float64, len(p.Genomes)),
		Novelty:    make([]float64, len(p.Genomes)),
		Complexity: make([]float64, len(p.Genomes)),
		Diversity:  []float64{float64(len(p.Species))},
	}
	for j, g := range p.Genomes {
		i.Fitness[j] = g.Fitness
		i.Novelty[j] = g.Novelty
		i.Complexity[j] = float64(g.Complexity())
	}

	b := new(bytes.Buffer)
	if err = json.NewEncoder(b).Encode(i); err != nil {
		return i.ID, err
	}

	if err = ib.Put(itob(i.ID), b.Bytes()); err != nil {
		return i.ID, err
	}

	// Add the population
	b = new(bytes.Buffer)
	g := gzip.NewWriter(b)
	if err = json.NewEncoder(g).Encode(p); err != nil {
		return i.ID, err
	}
	g.Flush()
	if err = tx.Bucket(populationsBucket).Put(itob(i.ID), b.Bytes()); err != nil {
		return i.ID, err
	}

	// Add to the trial's index and save the updated time
	if err = idxAddValue(tx, t2iBucket, tid, i.ID); err != nil {
		return i.ID, err
	}

	err = tx.Commit()
	return i.ID, err
}

func (c *Client) GetIteration(iid int) (w web.Iteration, err error) {

	// Begin a read-only transaction
	var tx *bolt.Tx
	if tx, err = c.db.Begin(false); err != nil {
		return
	}
	defer tx.Rollback()

	// Return the iteration
	return getIteration(tx, iid)
}

func getIteration(tx *bolt.Tx, iid int) (i web.Iteration, err error) {

	// Retrieve the iteration
	v := tx.Bucket(iterationsBucket).Get(itob(iid))
	if v == nil {
		err = fmt.Errorf("iteration %d does not exist", iid)
		return
	}

	// Decode the iteration
	if err = json.NewDecoder(bytes.NewBuffer(v)).Decode(&i); err != nil {
		return
	}

	return
}

func (c *Client) GetIterations(tid int) (is []web.Iteration, err error) {

	// Begin a read-only transaction
	var tx *bolt.Tx
	if tx, err = c.db.Begin(false); err != nil {
		return
	}
	defer tx.Rollback()

	// Retieve the iteration ids
	iids := idxValues(tx, t2iBucket, tid)
	is = make([]web.Iteration, 0, len(iids))
	for _, iid := range iids {
		var i web.Iteration
		if i, err = getIteration(tx, iid); err != nil {
			return
		}
		is = append(is, i)
	}

	// Return the iterations
	return
}
