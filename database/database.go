package database

import (
	"log"

	"github.com/boltdb/bolt"
)

// Database represents a Bolt-backed data store
type Database struct {
	*bolt.DB
}

// Tx represents a Bolt transaction
type Tx struct {
	*bolt.Tx
}

// View wrapps bolt.DB.View
func (db *Database) View(fn func(*Tx) error) error {
	return db.DB.View(func(tx *bolt.Tx) error {
		return fn(&Tx{tx})
	})
}

// Update wrapps bolt.DB.Update
func (db *Database) Update(fn func(*Tx) error) error {
	return db.DB.Update(func(tx *bolt.Tx) error {
		return fn(&Tx{tx})
	})
}

// initBuckets iterates over every bucket that should always exist
func (db *Database) initBuckets(buckets []string) {
	db.Update(func(tx *Tx) error {
		for _, bucket := range buckets {
			if _, err := tx.CreateBucketIfNotExists([]byte(bucket)); err != nil {
				log.Fatal(err)
			}
		}

		return nil
	})
}

// Open initializes or opens a database from a defined directory
func Open(path string) (db Database, err error) {

	db.DB, err = bolt.Open(path, 0600, nil)
	if err != nil {
		db.Close()

		return
	}

	var buckets = []string{usersBucket, adminBucket}
	db.initBuckets(buckets)

	return
}
