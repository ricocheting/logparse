package storage

import (
	"time"

	"github.com/boltdb/bolt"
	"github.com/ricocheting/logparse/internal"
)

// Store represents the data storage for storing messages received and sent.
type Store struct {
	path string
	db   *bolt.DB
}

// NewStore returns a new instance of Store.
func NewStore(path string) *Store {
	return &Store{
		path: path,
	}
}

// Path returns the data path.
func (st *Store) Path() string { return st.path }

// Open opens and initializes the database.
func (st *Store) Open() error {
	// Open underlying data store.
	db, err := bolt.Open(st.path, 0666, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	st.db = db

	// Initialize all the required buckets
	if err := st.db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists(internal.YearBucket)
		tx.CreateBucketIfNotExists(internal.HitsBucket)
		tx.CreateBucketIfNotExists(internal.ExtensionsBucket)
		tx.CreateBucketIfNotExists(internal.StatusCodesBucket)
		tx.CreateBucketIfNotExists(internal.IPSBucket)
		tx.CreateBucketIfNotExists(internal.NotFoundBucket)
		return nil
	}); err != nil {
		st.Close()
		return err
	}

	return nil
}

// Close closes the store.
func (st *Store) Close() error {
	if st.db != nil {
		st.db.Close()
	}
	return nil
}
