package admin

import "github.com/willeponken/go-cjdns/key"

const (
	BeaconDisable       = 0 //  No beacons are sent and incoming beacon messages are discarded.
	BeaconAccept        = 1 //  No beacons are sent but if an incoming beacon is received, it is acted upon.
	BeaconAcceptAndSend = 2 // Beacons are sent and accepted.
)

// ETHInterface_new creates a new ETHInterface and bind it to a device.
// Use the returned iface number with ETHInterface_beginConnection and
// ETHInterface_beacon.
func (c *Conn) ETHInterface_new(device string) (iface int, err error) {
	args := &struct {
		BindDevice string `bencode:"bindDevice"`
	}{device}

	resp := new(struct{ InterfaceNumber int })

	pack, err := c.sendCmd(&request{AQ: "ETHInterface_new", Args: args})
	if err == nil {
		err = pack.Decode(resp)
	}
	return resp.InterfaceNumber, err
}

// ETHInterface_beginConnection connects an ETHInterface to another computer which has an ETHInterface running.
// Use iface 0 for the first interface.
func (c *Conn) ETHInterface_beginConnection(iface int, mac, pass string, pubKey key.Public) error {
	args := &struct {
		InterfaceNumber int        `bencode:"interfaceNumber"`
		Password        string     `bencode:"password"`
		MacAddress      string     `bencode:"MacAddress"`
		PublicKey       key.Public `bencode:"publicKey"`
	}{iface, pass, mac, pubKey}
	_, err := c.sendCmd(&request{AQ: "ETHInterface_beginConnection", Args: args})
	return err
}

// ETHInterface_beacon enables or disables sending or receiving of  ETHInterface beacon messages.
// ETHInterface uses periodic beacon messages to automatically peer nodes which are on the same LAN.
// Be mindful that if your lan has is open wifi, enabling beaconing will allow anyone to peer with you.
//
// interfaceNumber is the number of the ETHInterface to change the state of. Use 0 for the first interface.
//
// state is the state to switch to, if -1 the current state will be queried only.
// See BeaconDisable, BeaconAccept, and BeaconAcceptAndSend.
func (c *Conn) ETHInterface_beacon(iface int, state int) (currentState int, stateDescription string, err error) {
	args := &struct {
		InterfaceNumber int `bencode:"interfaceNumber"`
		State           int `bencode:"state"`
	}{iface, state}

	resp := new(struct {
		State     int
		StateName string
	})

	pack, err := c.sendCmd(&request{AQ: "ETHInterface_beacon", Args: args})
	if err == nil {
		err = pack.Decode(resp)
	}
	return resp.State, resp.StateName, err
}
