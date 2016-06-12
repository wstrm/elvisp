package cjdns

import (
	"log"
	"net"

	"github.com/ehmry/go-cjdns/key"
)

// AddUser adds a new user to the database and allows a new iptunnel connection for the user.
func (c *Conn) AddUser(publicKey *key.Public, ip net.IP) error {
	if err := c.Conn.IpTunnel_allowConnection(publicKey, ip); err != nil {
		return err
	}

	return nil
}

// DelUser looks up the user for the defined public key and deauthenticates the user from the iptunnel.
func (c *Conn) DelUser(publicKey *key.Public) error {
	tunIndexes, err := c.Conn.IpTunnel_listConnections()
	if err != nil {
		return err
	}

	// TODO Check the tunConn type
	for tunConn := range tunIndexes {
		log.Println(tunConn)
	}

	return nil
}
