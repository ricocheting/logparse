package storage

// NOTE: you can't run more than one of these at a time as the Bolt database will lock each other. so use
// # cd ./logparse/storage
// # go test -run=TestDumpRawBoltHitsData

import (
	"fmt"
	"log"
	"testing"

	"github.com/ricocheting/logparse/internal"

	"github.com/boltdb/bolt"
)

// RAW dump the extensions in the database
func TestDumpRawBoltAll(t *testing.T) {
	fmt.Println("TestDumpRawBoltAll")

	db, err := bolt.Open("../data/db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Batch(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(internal.HitsBucket))

		b.ForEach(func(k []byte, v []byte) error {
			fmt.Printf("HitsBucket key=%s, value=%v\n", k, v)
			return nil
		})

		b = tx.Bucket([]byte(internal.IPSBucket))

		b.ForEach(func(k []byte, v []byte) error {
			fmt.Printf("IPSBucket key=%s, value=%v\n", k, v)
			return nil
		})

		b = tx.Bucket([]byte(internal.ExtensionsBucket))

		b.ForEach(func(k []byte, v []byte) error {
			fmt.Printf("ExtensionsBucket key=%s, value=%s\n", k, v)
			return nil
		})
		b = tx.Bucket([]byte(internal.StatusCodesBucket))

		b.ForEach(func(k []byte, v []byte) error {
			fmt.Printf("StatusCodesBucket key=%s, value=%s\n", k, v)
			return nil
		})
		return nil
	})
}

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
		b := tx.Bucket([]byte(internal.HitsBucket))

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
		b := tx.Bucket([]byte(internal.ExtensionsBucket))

		b.ForEach(func(k []byte, v []byte) error {
			fmt.Printf("key=%s, value=%s\n", k, v)
			return nil
		})
		return nil
	})
}

// RAW dump the extensions in the database
func TestDumpRawBoltDirectoriesData(t *testing.T) {
	fmt.Println("\nTestDumpRawBoltDirectoriesData")

	db, err := bolt.Open("../data/db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(internal.DirectoriesBucket))

		b.ForEach(func(k []byte, v []byte) error {
			fmt.Printf("key=%s, value=%s\n", k, v)
			return nil
		})
		return nil
	})
}

// DANGEROUS clear the contents of a bucket
/*func TestDELETE(t *testing.T) {
	fmt.Println("\nTestDELETE")

	db, err := bolt.Open("../data/db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket(internal.DirectoriesBucket)
		tx.CreateBucketIfNotExists(internal.DirectoriesBucket)
		return nil
	})
}*/
