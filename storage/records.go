package storage

import (
	"encoding/binary"
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/ricocheting/logparse/internal"
)

type Stat = internal.Stat
type Stats = internal.Stats

// SaveHits insert or update the total number of hits to the YYYYMMDD key
/*func (st *Store) SaveHits(dateKey []byte, hits uint64) error {
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
}*/

// SaveExtensions insert or update the total number of hits to the YYYYMMDD key
/*func (st *Store) SaveExtensions(dateKey []byte, data Stats) error {
	// TODO: needs error checking
	//sorted := data.ToSlice(0)

	return st.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(extensionsBucket)

		oldRaw := b.Get(dateKey)

		worker := Stats{}
		json.Unmarshal(oldRaw, &worker)

		// walk through the new data
		for key, value := range data {
			// if stats[i] exists, I need to add data[i] value to it.
			if wVal, ok := worker[key]; ok {
				worker[key] = wVal + value
			} else { // if it doesn't exist, I need to create it
				worker[key] = value
			}
		}

		// save stats back up and insert it into the dtabase
		buf, err := json.Marshal(worker) // INSERT value
		if err != nil {
			return err
		}

		b.Put(dateKey, buf)

		return nil
	})
}*/

// SaveBaseNumber insert or update the total number of hits to the YYYYMMDD key
func (st *Store) SaveBaseNumber(bucket []byte, dateKey []byte, value uint64) error {
	// TODO: needs error checking

	return st.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)

		oldVal := b.Get(dateKey)
		newVal := itob(value) // INSERT value

		if oldVal != nil {
			newVal = itob((value + btoi(oldVal))) // UPDATE value
		}

		//fmt.Printf("Hits: %+v = %+v\n", string(dateKey), newVal)
		b.Put(dateKey, newVal)

		return nil
	})
}

func (st *Store) SaveBaseStats(bucket []byte, dateKey []byte, data Stats) error {
	// TODO: needs error checking
	//sorted := data.ToSlice(0)

	return st.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)

		oldRaw := b.Get(dateKey)

		worker := Stats{}
		json.Unmarshal(oldRaw, &worker)

		// walk through the new data
		for key, value := range data {
			// if stats[i] exists, I need to add data[i] value to it.
			if wVal, ok := worker[key]; ok {
				worker[key] = wVal + value
			} else { // if it doesn't exist, I need to create it
				worker[key] = value
			}
		}

		// save stats back up and insert it into the dtabase
		buf, err := json.Marshal(worker) // INSERT value
		if err != nil {
			return err
		}

		b.Put(dateKey, buf)

		return nil
	})
}

// itob converts an uint64 to byte array
func itob(i uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, i)
	return b
}

// btoi converts a byte array to uint64
func btoi(b []byte) uint64 {
	i := binary.BigEndian.Uint64(b)
	return i
}
