package bolt

import (
	"bytes"
	"encoding/binary"
	"time"

	"github.com/boltdb/bolt"
)

// Client stores and retrieves data
type Client struct {
	Path string
	db   *bolt.DB
}

// Top-level buckets
var (
	studiesBucket     = []byte("Studies")
	experimentsBucket = []byte("Experiments")
	trialsBucket      = []byte("Trials")
	iterationsBucket  = []byte("Iterations")
	populationsBucket = []byte("Population")
	s2eBucket         = []byte("Studies-Experiment Index") // Study to experiments index
	e2tBucket         = []byte("Experiment-Trials Index")  // Experiment to trials index
	t2iBucket         = []byte("Trials-Iterations Index")  // Trial to iterations index

	null = []byte{0}
)

// Open and initialise the data file
func (c *Client) Open() error {

	// Open database file.
	var err error
	c.db, err = bolt.Open(c.Path, 0666, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}

	// Start writable transaction.
	var tx *bolt.Tx
	tx, err = c.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Initialize top-level buckets.
	var sb *bolt.Bucket
	if sb, err = tx.CreateBucketIfNotExists(studiesBucket); err != nil {
		return err
	}
	v := sb.Get(itob(0))
	if v == nil {
		if err = setStudy(tx, 0, "Unfiled experiments"); err != nil {
			return err
		}
	}

	for _, k := range [][]byte{experimentsBucket, trialsBucket, iterationsBucket, populationsBucket, s2eBucket, e2tBucket, t2iBucket} {
		if _, err = tx.CreateBucketIfNotExists(k); err != nil {
			return err
		}
	}

	// Save transaction to disk.
	return tx.Commit()
}

// Close the connection to the client
func (c *Client) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func btoi(b []byte) int {
	v := binary.BigEndian.Uint64(b)
	return int(v)
}

func idxAddValue(tx *bolt.Tx, idx []byte, k, v int) error {
	key := append(itob(k), itob(v)...)
	return tx.Bucket(idx).Put(key, null)
}

func idxDelValue(tx *bolt.Tx, idx []byte, k, v int) error {
	key := append(itob(k), itob(v)...)
	return tx.Bucket(idx).Delete(key)
}

func idxValues(tx *bolt.Tx, idx []byte, k int) []int {
	c := tx.Bucket(idx).Cursor()
	prefix := itob(k)
	values := make([]int, 0, 100)
	for key, _ := c.Seek(prefix); bytes.HasPrefix(key, prefix); key, _ = c.Next() {
		values = append(values, btoi(key[8:]))
	}
	return values
}

func idxDelKey(tx *bolt.Tx, idx []byte, k int) error {
	b := tx.Bucket(idx)
	c := b.Cursor()
	prefix := itob(k)
	for key, _ := c.Seek(prefix); bytes.HasPrefix(key, prefix); key, _ = c.Next() {
		if err := b.Delete(key); err != nil {
			return err
		}
	}
	return nil
}

func idxDelValues(tx *bolt.Tx, idx []byte, v int) error {
	b := tx.Bucket(idx)
	c := b.Cursor()
	suffix := itob(v)
	for key, _ := c.First(); key != nil; key, _ = c.Next() {
		if bytes.HasSuffix(key, suffix) {
			if err := b.Delete(key); err != nil {
				return err
			}
		}
	}
	return nil
}
