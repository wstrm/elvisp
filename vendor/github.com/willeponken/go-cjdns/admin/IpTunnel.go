package admin

import (
	"net"

	"github.com/willeponken/go-cjdns/key"
)

func (c *Conn) IpTunnel_allowConnection(publicKey *key.Public, addr net.IP) (err error) {
	if b := addr.To4(); b != nil {
		_, err = c.sendCmd(&request{AQ: "IpTunnel_allowConnection",
			Args: &struct {
				Ip     net.IP      `bencode:"ip4Address"`
				PubKey *key.Public `bencode:"publicKeyOfAuthorizedNode"`
			}{addr, publicKey}})
	} else {
		_, err = c.sendCmd(&request{AQ: "IpTunnel_allowConnection",
			Args: &struct {
				Ip     net.IP      `bencode:"ip6Address"`
				PubKey *key.Public `bencode:"publicKeyOfAuthorizedNode"`
			}{addr, publicKey}})
	}
	return
}

func (c *Conn) IpTunnel_connectTo(publicKey *key.Public) error {
	_, err := c.sendCmd(&request{AQ: "IpTunnel_connectTo",
		Args: &struct {
			PubKey *key.Public `bencode:"publicKeyOfNodeToConnectTo"`
		}{publicKey}})

	return err
}

// IpTunnel_listConnections returns a list of all current IP tunnels
func (c *Conn) IpTunnel_listConnections() (tunnelIndexes []int, err error) {
	resp := new(struct {
		Connections []int
	})

	var pack *packet
	pack, err = c.sendCmd(&request{AQ: "IpTunnel_listConnections"})
	if err == nil {
		err = pack.Decode(resp)
	}
	return resp.Connections, err
}

func (c *Conn) IpTunnel_removeConnection(connection int) error {
	_, err := c.sendCmd(&request{AQ: "IpTunnel_removeConnection",
		Args: &struct {
			Connection int `bencode:"connection"`
		}{connection}})
	return err
}

type IpTunnelConnection struct {
	Ip4Address *net.IP
	Ip6Address *net.IP
	Key        *key.Public
	Outgoing   bool
}

func (c *Conn) IpTunnel_showConnection(connection int) (*IpTunnelConnection, error) {
	resp := new(IpTunnelConnection)

	pack, err := c.sendCmd(&request{AQ: "IpTunnel_showConnection",
		Args: &struct {
			Connection int `bencode:"connection"`
		}{connection}})
	if err == nil {
		err = pack.Decode(resp)
	}
	return resp, err
}
