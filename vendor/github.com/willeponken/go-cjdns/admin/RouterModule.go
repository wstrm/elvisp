package admin

import "errors"

// RouterModule_getPeers sends a request to a node for a list of its peers.
// The format of the peers list is 'v${version}.${label}.${publicKey}'.
//
// Timeout and nearbyPath default when there are 0 and "" respectively.
func (c *Conn) RouterModule_getPeers(path string, timeout int, nearbyPath string) (peers []string, ms int, err error) {
	req := request{
		AQ: "RouterModule_getPeers",
		Args: &struct {
			Path       string `bencode:"path"`
			Timeout    int    `bencode:"timeout,omitempty"`
			NearbyPath string `bencode:"nearbyPath,omitempty"`
		}{path, timeout, nearbyPath},
	}

	var pack *packet
	if pack, err = c.sendCmd(&req); err == nil {
		err = pack.Decode(&struct {
			Peers *[]string
			Ms    *int
		}{&peers, &ms})
	}
	return
}

//RouterModule_lookup returns a single path for an address. Not sure what this is used for
func (c *Conn) RouterModule_lookup(address string) (response map[string]interface{}, err error) {
	var (
		args = &struct {
			Address string `bencode:"address"`
		}{address}

		pack *packet
	)

	pack, err = c.sendCmd(&request{AQ: "RouterModule_lookup", Args: args})
	if err == nil {
		err = pack.Decode(response)
	}
	return
}

// Pings the specified IPv6 address or switch label and will timeout if it takes longer than the specified timeout period.
// CJDNS will fallback to its own timeout if the a zero timeout is given.
func (c *Conn) RouterModule_pingNode(addr string, timeout int) (ms int, version string, err error) {
	args := &struct {
		Path    string `bencode:"path"`
		Timeout int    `bencode:"timeout,omitempty"`
	}{addr, timeout}

	resp := new(struct {
		Ms      int    // number of milliseconds since the original ping
		Result  string // set when ping times out
		Version string // git hash of the source code which the node was built on
	})

	var pack *packet
	pack, err = c.sendCmd(&request{AQ: "RouterModule_pingNode", Args: args})
	if err == nil {
		err = pack.Decode(resp)
	}
	if err == nil && resp.Ms == 0 {
		err = errors.New(resp.Result)
	}
	return resp.Ms, resp.Version, err
}
