package storage

import (
	"encoding/binary"

	"github.com/boltdb/bolt"
)

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
/*func (st *Store) SaveExtensions(dateKey []byte, data Stat) error {
	// TODO: needs error checking

	return st.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(extensionsBucket)

		oldVal := b.Get(dateKey)
		var task []Stat
		json.Unmarshal(oldVal, &task)

		fmt.Printf("Hits: %+v = %s = %s\n", string(dateKey), oldVal, task)

		for i, e := range task {
			fmt.Printf("Hits: %+v = %s\n", i, e)
			// i is the index, e the element
		}

		buf, err := json.Marshal(data) // INSERT value
		if err != nil {
			return err
		}


			if oldVal != nil {
				newVal = []byte(data) // UPDATE value
			}

		//fmt.Printf("Hits: %+v = %+v\n", string(dateKey), newVal)
		b.Put(dateKey, buf)

		return nil
	})
}*/

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
