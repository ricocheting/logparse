package storage

// NOTE: you can't run more than one of these at a time as the Bolt database will lock each other. so use
// # go test -run=TestDumpRawBoltHitsData

import (
	"fmt"
	"log"
	"testing"

	"github.com/boltdb/bolt"
)

// RAW dump the users in the database
func TestDumpRawBoltHitsData(t *testing.T) {
	fmt.Println("We're Running")

	db, err := bolt.Open("../data/db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("hits"))

		b.ForEach(func(k []byte, v []byte) error {
			fmt.Printf("key=%s, value=%v\n", k, v)
			return nil
		})
		return nil
	})
}
