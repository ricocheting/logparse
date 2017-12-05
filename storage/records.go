package storage

import (
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/ricocheting/logparse/internal"
)

type Stat = internal.Stat
type Stats = internal.Stats

// SaveBaseNumber insert or update the bucket on the YYYYMMDD dateKey with value (if exists, adds value to existing amount)
func (st *Store) SaveBaseNumber(bucket []byte, dateKey []byte, value uint64) error {
	// TODO: needs error checking

	return st.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)

		oldVal := b.Get(dateKey)
		newVal := internal.Itob(value) // INSERT value

		if oldVal != nil {
			newVal = internal.Itob(value + internal.Btoi64(oldVal)) // UPDATE value
		}

		//fmt.Printf("Hits: %+v = %+v\n", string(dateKey), newVal)
		b.Put(dateKey, newVal)

		return nil
	})
}

// SaveBaseStats insert or update the bucket on the YYYYMMDD dateKey with the Stats collection (if collection item exists, adds new value to existing amount)
func (st *Store) SaveBaseStats(bucket []byte, dateKey []byte, data Stats) error {
	// TODO: needs error checking
	//sorted := data.ToSlice(0)

	return st.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		raw := b.Get(dateKey)

		worker := Stats{}
		json.Unmarshal(raw, &worker)

		// walk through the new data
		for key, value := range data {
			// if stored worker[i] exists, I need to add data[i] value to it.
			if wVal, ok := worker[key]; ok {
				worker[key] = wVal + value
			} else { // if it doesn't exist, I need to create it
				worker[key] = value
			}
		}

		// save stats back up and insert it into the dtabase
		buf, err := json.Marshal(worker)
		if err != nil {
			return err
		}

		b.Put(dateKey, buf)

		return nil
	})
}
