package storage

import (
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/ricocheting/logparse/internal"
)

// where prefix = YYYY or YYYYMM that we want to watch in our search
/*func (st *Store) FilterBaseNumber(bucket []byte, prefix []byte) ([]Stat, error) {
	var stats = []Stat{}

	return stats, st.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		c := tx.Bucket(bucket).Cursor()

		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {

			stat := Stat{Name: string(k[:]), Value: internal.Btoi64(v)}

			stats = append(stats, stat)
		}

		return nil
	})
}*/

// ListBaseNumber
func (st *Store) ListBaseNumber(bucket []byte) (internal.StatTotal, error) {
	var data internal.StatTotal

	return data, st.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket(bucket)

		// ForEach instead of b.Cursor() because we know we'll be iterating over all the keys
		b.ForEach(func(k []byte, v []byte) error {

			//data.Get(k[0:4]).Get(k[4:6]).AddDay(k[6:8], v)
			data.AddTotal(k[0:4], k[4:6], k[6:8], v)

			return nil
		})

		return nil
	})
}

// ListBaseStats
func (st *Store) ListBaseStats(bucket []byte) (internal.StatTotal, error) {
	var data internal.StatTotal

	return data, st.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)

		b.ForEach(func(k []byte, v []byte) error {
			worker := Stats{}
			json.Unmarshal(v, &worker)

			for key, value := range worker {
				//data.Add(string(k), key, value)
				data.AddStat(k[0:4], k[4:6], k[6:8], key, value)
				//data.GrandTotal += value
			}

			return nil
		})

		return nil
	})
}

// ListPages
func (st *Store) ListPages(bucket []byte) (internal.StatTotal, error) {
	var data internal.StatTotal

	return data, st.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)

		b.ForEach(func(k []byte, v []byte) error {
			worker := Stats{}
			json.Unmarshal(v, &worker)
			//stat := Stat{Name: string(k[:]), Value: 0}
			var n uint64 = 0

			for key, value := range worker {
				if key == "" || key == ".shtml" || key == ".php" || key == ".htm" || key == ".html" {
					n += value
				}
			}

			data.AddTotal(k[0:4], k[4:6], k[6:8], internal.Itob(n))

			return nil
		})

		return nil
	})
}

// ListErrors
func (st *Store) ListErrors(bucket []byte) (internal.StatErrors, error) {
	var data internal.StatErrors

	return data, st.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)

		b.ForEach(func(page []byte, v []byte) error {
			worker := Stats{}
			json.Unmarshal(v, &worker)

			for missing, value := range worker {
				//data.Add(string(k), key, value)
				data.SetVal(string(page), missing, value)
				//data.GrandTotal += value
			}

			return nil
		})

		return nil
	})
}
