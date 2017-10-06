package admin

import "errors"

// SwitchPinger_ping send a switch level ping. There is no routing table lookup and the router is not involved.
// Pinging IP addresses this way is not possible.
//
// data is an optional string that the destination switch will echo back.
func (c *Conn) SwitchPinger_ping(path, dataIn string, timeout int) (dataOut string, ms int, err error) {
	var (
		args = &struct {
			Data    string `bencode:"data,omitempty"`
			Path    string `bencode:"path"`
			Timeout int    `bencode:"timeout,omitempty"`
		}{dataIn, path, timeout}

		resp = new(struct {
			Data, Result string
			Ms           int
		})

		pack *packet
	)
	if pack, err = c.sendCmd(&request{AQ: "SwitchPinger_ping", Args: args}); err == nil {
		err = pack.Decode(resp)
	}
	if err == nil && resp.Result != "pong" {
		err = errors.New(resp.Result)
	}
	return resp.Data, resp.Ms, err
}
