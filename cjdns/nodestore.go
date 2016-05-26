package cjdns

// LookupPubKey finds the public key for a ipv6 in the cjdns node store.
func (c *Conn) LookupPubKey(ip string) (key string, err error) {
	node, err := c.Conn.NodeStore_nodeForAddr(ip)
	if err != nil {
		return
	}

	key = node.Key
	return
}
