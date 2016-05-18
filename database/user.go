package database

import (
	"errors"
	"log"
)

// UsersBucket defines the namespace for the user bucket
const UsersBucket = "Users"

// User contains the public key and last lease (integer is used to increment the IP)
type User struct {
	PublicKey string
	Lease     int
}

// AddUser inserts a new user into the UserBucket with public key and last lease
func (db *Database) AddUser(user User) (err error) {
	err = db.Update(func(tx *Tx) error {
		log.Printf("Adding new user with key: %s.", user.PublicKey)

		parentBucket := tx.Bucket([]byte(UsersBucket))

		userBucket, err := parentBucket.CreateBucketIfNotExists([]byte(user.PublicKey))
		if err != nil {
			return err
		}

		err = userBucket.Put([]byte("lease"), []byte(user.Lease))

		return err // End of transaction
	})

	return
}

// UpdateLease checks if the lease for a user is updated and if so generates a new lease
func (db *Database) UpdateLease(pubkey string) (err error) {
	err = db.Update(func(tx *Tx) error {
		log.Printf("Updating lease for user with key: %s.", pubkey)

		parentBucket := tx.Bucket([]byte(UsersBucket))

		userBucket := parentBucket.Bucket([]byte(pubkey))
		if userBucket != nil {
			err = errors.New("User with public key: " + pubkey + " does not exist in database.")
			log.Println(err)

			return err
		}

		// TODO Generate new lease (integer) by iterating over all the available ones.

		return nil
	})

	return
}
