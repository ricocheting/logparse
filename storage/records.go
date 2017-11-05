package storage

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/ricocheting/logparse/internal"
)

type Stat = internal.Stat
type Stats = internal.Stats

// SaveHits insert or update the total number of hits to the YYYYMMDD key
func (st *Store) SaveHits(dateKey []byte, hits uint64) error {
	// TODO: needs error checking

	return st.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(hitsBucket)

		oldVal := b.Get(dateKey)
		newVal := itob(hits) // INSERT value

		if oldVal != nil {
			newVal = itob((hits + btoi(oldVal))) // UPDATE value
		}

		//fmt.Printf("Hits: %+v = %+v\n", string(dateKey), newVal)
		b.Put(dateKey, newVal)

		return nil
	})
}

// SaveExtensions insert or update the total number of hits to the YYYYMMDD key
func (st *Store) SaveExtensions(dateKey []byte, data Stats) error {
	// TODO: needs error checking
	sorted := data.ToSlice(0)

	return st.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(extensionsBucket)

		oldVal := b.Get(dateKey)

		// trainwreck code. I need to read the old values,
		var stats []Stat
		json.Unmarshal(oldVal, &stats)

		fmt.Printf("Hits: %+v = %s = %s\n", string(dateKey), oldVal, stats)

		// yes, I'm broken because data is not an array
		for i, e := range sorted {
			fmt.Printf("Extensions: %+v = %s\n", i, e)
			// i is the index, e the element
			// if stats[i] exists, I need to add data[i] value to it.
			// if it doesn't exist, I need to create it
		}

		// save stats back up and insert it into the dtabase
		buf, err := json.Marshal(stats) // INSERT value
		if err != nil {
			return err
		}

		//fmt.Printf("Extensions: %+v = %+v\n", string(dateKey), newVal)
		b.Put(dateKey, buf)

		return nil
	})
}

// itob returns an 8-byte big endian representation of v.
func itob(i uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, i)
	return b
}

// convert a byte array to uint64
func btoi(b []byte) uint64 {
	i := binary.BigEndian.Uint64(b)
	return i
}
