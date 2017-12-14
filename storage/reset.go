package storage

import (
	"github.com/boltdb/bolt"
	"github.com/ricocheting/logparse/internal"
)

// clear all entries in internal.ErrorsBucket by deleting the bucket
func (st *Store) ResetErrors() error {

	return st.db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket(internal.ErrorsBucket)
	})
}
