package config

// Thanks to SashaCrofter for the original layout of these structures
type Config struct {
	//the private key for this node
	PrivateKey string `json:"privateKey"`

	//the public key for this node
	PublicKey string `json:"publicKey"`

	//this node's IPv6 address as (derived from publicKey)
	IPv6 string `json:"ipv6"`

	//authorized passwords
	AuthorizedPasswords []AuthPass `json:"authorizedPasswords"`

	//information for RCP server
	Admin AdminBlock `json:"admin"`

	//interfaces for the switch core
	Interfaces InterfacesList `json:"interfaces"`

	//configuration for the router
	Router RouterBlock `json:"router"`

	//remove cryptoauth sessions after this number of seconds
	ResetAfterInactivitySeconds int `json:"resetAfterInactivitySeconds"`

	//cjdns security options
	RawSecurity interface{} `json:"security"`
	//usable representation of the security info that can not be saved to JSON
	Security SecurityBlock `json:"-"`

	// Where to log to
	Logging Logging

	// Fork to background
	Background int `json:"NoBackground"`

	//the internal config file version (mostly unused)
	Version int `json:"version"`
}
type Logging struct {
	LogTo string
}
type AuthPass struct {
	Password string `json:"password"` //the password for incoming authorization
}

type AdminBlock struct {
	Bind     string `json:"bind"`     //the port to bind the RCP server to
	Password string `json:"password"` //the password for the RCP server
}

type InterfacesList struct {
	// Connections done via UDP over existing networks
	UDPInterface []UDPInterfaceBlock `json:"UDPInterface,omitempty"`

	// Use raw ethernet frames.
	ETHInterface []EthInterfaceBlock `json:"ETHInterface,omitempty"`
}

type UDPInterfaceBlock struct {
	//Address to bind to ("0.0.0.0:port")
	Bind string `json:"bind"`

	//Maps connection information to peer details, where the Key is the peer's
	//IPv4 address and port and the Connection contains all of the information
	//about the peer, such as password and public key
	ConnectTo map[string]Connection `json:"connectTo"`
}

type EthInterfaceBlock struct {
	//Interface to bind to ("eth0")
	Bind string `json:"bind"`

	//Maps connection information to peer details, where the Key is the peer's
	//MAC address and the Connection contains all of the information about the
	//peer, such as password and public key
	ConnectTo map[string]Connection `json:"connectTo"`

	//Sets the beacon state for the ether interface. 0 = disabled, 1 = accept
	//beacons, 2 = send and accept beacons.
	Beacon int `json:"beacon"`
}

type Connection struct {
	//the password to connect to the peer node
	Password string `json:"password"`

	//the peer node's public key
	PublicKey string `json:"publicKey"`
}

type RouterBlock struct {
	//interface used for connecting to the cjdns network
	Interface RouterInterface `json:"interface"`

	//interface used for connecting to the cjdns network
	IPTunnel TunnelInterface `json:"ipTunnel"`
}

type RouterInterface struct {
	//the type of interface
	Type string `json:"type"`

	//the persistent interface to use for cjdns (not usually used)
	TunDevice string `json:"tunDevice,omitempty"`
}
type TunnelInterface struct {
	//A list of details for users connecting to us to form an IP tunnel
	AllowedConnections []TunnelAllowed `json:"allowedConnections"`

	//A list of nodes we will connect to in order to form an IP tunnel
	OutgoingConnections []string `json:"outgoingConnections"`
}
type TunnelAllowed struct {
	//the peer node's public key
	Publickey string `json:"publicKey"`

	//the IPv4 address we will assign to the peer's tunnel (we only need to
	//specify either the IPv4 or IPv6 addresses)
	IP4Address string `json:"ip4Address"`

	//the IPv6 address we will assign to the peer's tunnel (we only need to
	//specify either the IPv4 or IPv6 addresses)
	IP6Address string `json:"ip6Address"`
}

//We can not unmarshall the security section of the config directly to a
//useable structure, so we manually save and restore the values using the
//SecurityBlock We set them by parsing the RawSecurity interface{} This allows
//us to easily edit these values in our program. Note that the RawSecurity
//interface{} is what actaully gets marshalled back in to JSON We must parse
//SecurityBlock and create the proper RawSecurity interface{} before
//marshalling
type SecurityBlock struct {
	NoFiles int
	SetUser string
}
