package cjdns

import (
	"log"
	"net"

	"github.com/ehmry/go-cjdns/key"
)

// AddUser adds a new user to the database and allows a new iptunnel connection for the user.
func (c *Conn) AddUser(publicKey *key.Public, ip net.IP) error {
	admin := c.Conn

	if err := admin.IpTunnel_allowConnection(publicKey, ip); err != nil {
		return err
	}

	log.Printf("User: %s added to cjdns IP tunnel", publicKey.String())

	return nil
}

// DelUser looks up the user for the defined public key and deauthenticates the user from the iptunnel.
func (c *Conn) DelUser(publicKey *key.Public) error {
	admin := c.Conn

	tunnels, err := admin.IpTunnel_listConnections()
	if err != nil {
		return err
	}

	for i := range tunnels {
		tunnel, err := admin.IpTunnel_showConnection(tunnels[i])
		if err != nil {
			return err
		}

		if publicKey.Equal(tunnel.Key) {
			if err := admin.IpTunnel_removeConnection(tunnels[i]); err != nil {
				return err
			}
		}
	}

	return nil
}
