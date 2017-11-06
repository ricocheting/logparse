package storage

import "github.com/boltdb/bolt"
import "github.com/ricocheting/logparse/internal"

// ListTasks to return array of all tasks in the database (that match certain criteria)
// priority parameter is option. If a blank string is passed in, all records are returned
// NOTE: does not retrieve acts. they will be blank unless filled in elsewhere
func (st *Store) ListHits() ([]Stat, error) {
	var stats = []Stat{}

	return stats, st.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket(internal.HitsBucket)

		// ForEach instead of b.Cursor() because we know we'll be iterating over all the keys
		b.ForEach(func(k []byte, v []byte) error {
			//fmt.Printf("key=%s, value=%s\n", k, v)

			stat := Stat{Name: string(k[:]), Value: btoi(v)}

			stats = append(stats, stat)

			return nil
		})

		return nil
	})
}
