package storage

import (
	"encoding/binary"
	"fmt"

	"github.com/boltdb/bolt"
)

// SaveHits to find next ID and save task to database
func (st *Store) SaveHits(dateKey []byte, hits uint64) error {
	// TODO: needs error checking

	return st.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(hitsBucket)

		b.Put(dateKey, itob(hits))
		fmt.Printf("SAVE Hits: %+v = %+v\n", string(dateKey), hits)

		return nil
	})
}

// itob returns an 8-byte big endian representation of v.
func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

func inIntArray(a uint64, list []uint64) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
