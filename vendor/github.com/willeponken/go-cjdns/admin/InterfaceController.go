package admin

import "github.com/willeponken/go-cjdns/key"

func (c *Conn) InterfaceController_disconnectPeer(pubKey key.Public) error {
	_, err := c.sendCmd(&request{
		AQ: "InterfaceController_disconnectPeer",
		Args: &struct {
			PubKey key.Public `bencode:"pubKey"`
		}{pubKey}})
	return err

}

//type PeerState int

// Peer state values
const (
	Unauthenticated = iota
	Handshake
	Established
	Unresponsive
)

/*
var (
	peerStateStrings = [4]string{
		"UNAUTHENTICATED",
		"HANDSHAKE",
		"ESTABLISHED",
		"UNRESPONSIVE",
	}
)
*/

/*
func (s PeerState) String() string {
	if s.Int() < 0 || s.Int() > len(peerStateStrings)-1 {
		return "INVALID"
	}
	return peerStateStrings[s]
}

func (s PeerState) Int() int {
	return int(s)
}
*/

// Peer statistics
type PeerStats struct {
	PublicKey          *key.Public // Public key of peer
	SwitchLabel        *Path       // Internal switch label for reaching the peer
	IsIncoming         bool        // Is the peer connected to us, or us to them
	BytesOut           int64       // Total number of bytes sent
	BytesIn            int64       // Total number of bytes received
	State              string      // Peer connection state
	Last               int64       // Last time a packet was received from the peer
	ReceivedOutOfRange int
	Duplicates         int
	LostPackets        int
}

// Returns stats on currently connected peers
func (c *Conn) InterfaceController_peerStats() ([]*PeerStats, error) {
	var (
		args = new(struct {
			Page int `bencode:"page"`
		})
		req = request{AQ: "InterfaceController_peerStats", Args: args}

		resp = new(struct {
			More  bool
			Peers []*PeerStats //`bencode:"peers"`
			//Total int
		})

		pack *packet
		err  error
	)

	resp.More = true
	for resp.More {
		resp.More = false
		if pack, err = c.sendCmd(&req); err == nil {
			err = pack.Decode(resp)
		}
		if err != nil {
			break
		}
		args.Page++
	}
	return resp.Peers, err
}
