package database

import (
	"encoding/binary"
	"errors"
	"log"

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

func (db *Database) userExists(identifier interface{}) (pos []byte, exists bool) {
	err := db.View(func(tx *Tx) error {
		bucket := tx.Bucket([]byte(usersBucket))

		pubkey, isPubkey := identifier.(*key.Public)
		if isPubkey {
			strPubkey := pubkey.String()
			cursor := bucket.Cursor()

			for k, pk := cursor.First(); pk != nil; k, pk = cursor.Next() {
				if string(pk) == strPubkey {
					pos = k
					exists = true
					return nil
				}
			}

			exists = false
			return nil
		}

		id, isID := identifier.(uint64)
		if isID {
			p := uint64ToBin(id)
			if user := bucket.Get(p); user != nil {
				pos = p
				exists = true
				return nil
			}

			exists = false
			return nil
		}

		return errors.New("Unknown identifier specified")
	})

	if err != nil {
		exists = false
	}

	return
}

// AddUser inserts a new user into the UserBucket with public key and ID (used as seed for lease).
func (db *Database) AddUser(pubkey *key.Public) (id uint, err error) {
	k := pubkey.String()

	_, exists := db.userExists(pubkey)
	if exists {
		err = errors.New("User with public key: " + pubkey.String() + " already exists")
		log.Println(err)
		return
	}

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

	pos, exists := db.userExists(identifier)
	if !exists {
		err = errors.New("User identified as: %v does not exist")
		log.Println(err)
		return
	}

	err = db.Update(func(tx *Tx) error {

		bucket := tx.Bucket([]byte(usersBucket))

		log.Printf("Deleting user identified as: %v", identifier)
		err = bucket.Delete(pos)

		return nil
	})

	return
}
