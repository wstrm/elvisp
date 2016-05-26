package cjdns

import (
	"errors"

	"github.com/fc00/go-cjdns/key"
)

// AddUser adds a new user to the database and allows a new iptunnel connection for the user
func (c *Conn) AddUser(publicKey string) error {
	key, err := key.DecodePublic(publicKey)
	if err != nil {
		return err
	}

	if key.Valid() == false {
		return errors.New("Invalid public key")
	}

	return nil
}
