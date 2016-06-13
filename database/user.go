package database

import (
	"encoding/binary"
	"errors"
	"log"
	"reflect"

	"github.com/ehmry/go-cjdns/key"
)

// usersBucket defines the namespace for the user bucket
const usersBucket = "Users"

// uint64ToBin returns an 8-byte big endian representation of v
func uint64ToBin(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

// AddUser inserts a new user into the UserBucket with public key and ID (used as seed for lease).
func (db *Database) AddUser(pubkey *key.Public) (id uint, err error) {
	k := pubkey.String()

	err = db.Update(func(tx *Tx) error {

		bucket := tx.Bucket([]byte(usersBucket))
		seq, _ := bucket.NextSequence()
		id = uint(seq)

		log.Printf("Adding new user with key: %s and ID: %d", k, id)

		return bucket.Put(uint64ToBin(seq), []byte(k)) // End of transaction after data is put
	})

	return
}

// DelUser removes a registered user using the pubkey as identifier.
func (db *Database) DelUser(identifier interface{}) (err error) {
	id := reflect.ValueOf(identifier)
	idType := id.Type()

	err = db.Update(func(tx *Tx) error {

		bucket := tx.Bucket([]byte(usersBucket))

		if idType.Kind() == reflect.String {
			log.Printf("Identifier for user interpreted to string")

			pubkey, ok := identifier.(string)
			if !ok {
				err = errors.New("Failed to convert identifier to string")
				return nil
			}
			cursor := bucket.Cursor()

			for _, pk := cursor.First(); string(pk) != pubkey; _, pk = cursor.Next() {
				err = cursor.Delete()
				return nil
			}

			err = errors.New("Unable to delete user with public key:" + pubkey + ", because it does not exist")
		}

		if idType.Kind() == reflect.Int {
			log.Println("Identifier for user interpreted as integer")

			id, ok := identifier.(uint64)
			if !ok {
				err = errors.New("Failed to convert identifier to integer")
				return nil
			}

			return bucket.Delete(uint64ToBin(id))
		}

		return nil
	})

	return
}
