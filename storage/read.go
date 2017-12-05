package storage

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/ricocheting/logparse/internal"
)

// where prefix = YYYY or YYYYMM that we want to watch in our search
func (st *Store) FilterBaseNumber(bucket []byte, prefix []byte) ([]Stat, error) {
	var stats = []Stat{}

	return stats, st.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		c := tx.Bucket(bucket).Cursor()

		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {

			stat := Stat{Name: string(k[:]), Value: internal.Btoi(v)}

			stats = append(stats, stat)
		}

		return nil
	})
}

// ListBaseNumber
func (st *Store) ListBaseNumber(bucket []byte) (internal.StatYear, error) {
	var data = internal.StatYear{}
	data.Years = make(map[string]internal.StatMonth)

	return data, st.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket(bucket)

		// ForEach instead of b.Cursor() because we know we'll be iterating over all the keys
		b.ForEach(func(k []byte, v []byte) error {
			//fmt.Printf("key=%s, value=%s\n", k, v)

			year, err := data.Years[string(k[0:4])]

			if !err {
				year.Months = make(map[string]internal.StatDay)
			}

			month, err := data.Years[string(k[0:4])].Months[string(k[4:6])]

			if !err {
				month.Days = make(map[string]uint64)
			}

			//data.Years[string(k[0:4])].Months[string(k[4:6])].Days[string(k[6:8])] = internal.Btoi(v)

			//stat := Stat{Name: string(k[:]), Value: internal.Btoi(v)}
			year.GrandTotal += internal.Btoi(v)
			// get year, month, day from k
			fmt.Printf("%s, %s, %s: ", k[0:4], k[4:6], k[6:8]) //2017, 12, 03
			fmt.Printf("%d\n", year.GrandTotal)

			//data.Collect[string(k[0:4])].Collect[k[4:6]].Stats = append(data.Collect[string(k[0:4])].Collect[k[4:6]].Stats, stat)
			//stats = append(stats, stat)

			return nil
		})

		return nil
	})
}

// ListBaseStats
func (st *Store) ListBaseStats(bucket []byte) (internal.StatCollection, error) {
	var data = internal.StatCollection{
		Collect:    map[string]internal.Stats{},
		GrandTotal: 0,
	}

	return data, st.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)

		b.ForEach(func(k []byte, v []byte) error {
			worker := Stats{}
			json.Unmarshal(v, &worker)

			for key, value := range worker {
				data.Add(string(k), key, value)
				data.GrandTotal += value
			}

			return nil
		})

		return nil
	})
}

// ListPages
func (st *Store) ListPages(bucket []byte) ([]Stat, error) {
	var stats = []Stat{}

	return stats, st.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)

		b.ForEach(func(k []byte, v []byte) error {
			worker := Stats{}
			json.Unmarshal(v, &worker)
			stat := Stat{Name: string(k[:]), Value: 0}

			for key, value := range worker {
				if key == "" || key == ".shtml" || key == ".php" || key == ".htm" || key == ".html" {
					stat.Value += value
				}
			}

			stats = append(stats, stat)

			return nil
		})

		return nil
	})
}
