package database

import (
	"encoding/binary"
	"log"
)

// UsersBucket defines the namespace for the user bucket
const UsersBucket = "Users"

// uint64ToBin returns an 8-byte big endian representation of v
func uint64ToBin(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

// AddUser inserts a new user into the UserBucket with public key and last lease
func (db *Database) AddUser(pubkey string) (err error) {
	err = db.Update(func(tx *Tx) error {

		usersBucket := tx.Bucket([]byte(UsersBucket))
		lease, _ := usersBucket.NextSequence()

		log.Printf("Adding new user with key: %s and lease: %v.", pubkey, lease)

		return usersBucket.Put(uint64ToBin(lease), []byte(pubkey)) // End of transaction after data is put
	})

	return
}
