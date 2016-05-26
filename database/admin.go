package database

import "log"

// adminBucket defines the namespace for the admin bucket
const adminBucket = "Users"

// hashKey defines the key for storing hashed administration password
const hashKey = "hash"

// SetAdmin sets the hashed password for the administrator
func (db *Database) SetAdmin(hash string) (err error) {
	err = db.Update(func(tx *Tx) error {

		bucket := tx.Bucket([]byte(adminBucket))

		log.Printf("Updating password hash for administration")
		err = bucket.Put([]byte(hashKey), []byte(hash))

		return nil // End of transaction after data is put
	})

	return
}

// AdminHash retrieves the password hash for administrator
func (db *Database) AdminHash() (hash string, err error) {
	err = db.View(func(tx *Tx) error {

		bucket := tx.Bucket([]byte(adminBucket))

		log.Printf("Retrieving password hash")
		hash = string(bucket.Get([]byte(hashKey)))

		return nil // End of transaction after data is put
	})

	return
}
