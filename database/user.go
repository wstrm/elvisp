package database

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"

	"github.com/fc00/go-cjdns/key"
)

// usersBucket defines the namespace for the user bucket.
const usersBucket = "Users"

// uint64ToBin returns an 8-byte big endian representation of v.
func uint64ToBin(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

// binToUint64 takes an 8-byte big endian and converts it into a uint64.
func binToUint64(v []byte) uint64 {
	if len(v) != 8 {
		log.Fatalf("Invalid length of binary: %d", len(v))
	}

	return binary.BigEndian.Uint64(v)
}

// userExists takes a identifier and tries to type cast it into either a public key or a uint64, and then lookups that identifier in the database.
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

// nextUserID iterates over all users and tries to find a available slot between users, if there is no available it will return the next sequence available after all users.
func (db *Database) nextUserID() (id uint64, err error) {
	err = db.View(func(tx *Tx) error {
		bucket := tx.Bucket([]byte(usersBucket))

		cursor := bucket.Cursor()
		lastID := uint64(0)
		for currentID, _ := cursor.First(); currentID != nil; currentID, _ = cursor.Next() {

			log.Printf("current ID: %v", binToUint64(currentID))
			// The current ID is larger than the last found ID, this means that there is a "gap" that can be filled,
			// so we'll stop the iteration and use the last ID + 1 as it's available.
			if binToUint64(currentID) > lastID+1 {
				break
			}

			lastID = binToUint64(currentID)
		}

		id = lastID + 1

		return nil
	})

	return
}

// AddUser inserts a new user into the UserBucket with public key and ID (used as seed for lease).
func (db *Database) AddUser(pubkey *key.Public) (id uint64, err error) {
	k := pubkey.String()

	_, exists := db.userExists(pubkey)
	if exists {
		err = fmt.Errorf("User with public key: %s already exists", k)
		log.Println(err)
		return
	}

	err = db.Update(func(tx *Tx) error {

		bucket := tx.Bucket([]byte(usersBucket))

		id, err = db.nextUserID()
		if err != nil {
			return err
		}

		log.Printf("Adding new user with key: %s and ID: %d", k, id)

		return bucket.Put(uint64ToBin(id), []byte(k)) // End of transaction after data is put
	})

	return
}

// GetID returns the ID for a registered user.
func (db *Database) GetID(pubkey *key.Public) (id uint64, err error) {
	pos, exists := db.userExists(pubkey)
	if !exists {
		err = fmt.Errorf("User with public key: %s does not exist", pubkey.String())
		log.Println(err)
		return
	}

	id = binToUint64(pos)

	return
}

// DelUser removes a registered user using the pubkey as identifier.
func (db *Database) DelUser(identifier interface{}) (err error) {

	pos, exists := db.userExists(identifier)
	if !exists {
		err = fmt.Errorf("User identified as: %v does not exist", identifier)
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
