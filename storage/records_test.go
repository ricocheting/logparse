package storage

// NOTE: you can't run more than one of these at a time as the Bolt database will lock each other. so use
// # go test -run=TestDumpRawBoltHitsData

import (
	"fmt"
	"log"
	"testing"

	"github.com/boltdb/bolt"
)

// RAW dump the hits in the database
func TestDumpRawBoltHitsData(t *testing.T) {
	fmt.Println("TestDumpRawBoltHitsData")

	db, err := bolt.Open("../data/db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(hitsBucket))

		b.ForEach(func(k []byte, v []byte) error {
			fmt.Printf("key=%s, value=%v\n", k, v)
			return nil
		})
		return nil
	})
}

// RAW dump the extensions in the database
func TestDumpRawBoltExtensionsData(t *testing.T) {
	fmt.Println("\nTestDumpRawBoltExtensionsData")

	db, err := bolt.Open("../data/db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(extensionsBucket))

		b.ForEach(func(k []byte, v []byte) error {
			fmt.Printf("key=%s, value=%s\n", k, v)
			return nil
		})
		return nil
	})
}
