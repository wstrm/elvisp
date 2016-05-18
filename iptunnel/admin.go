package iptunnel

import (
	"github.com/fc00/go-cjdns/admin"
)

// Conn wraps around a go-cjdns admin connection
type Conn struct {
	Conn *admin.Conn
}

// Connect returns a connection to cjdns admin
func Connect(addr string, port int, password string) (conn *Conn, err error) {
	conf := admin.CjdnsAdminConfig{
		Addr:     addr,
		Port:     port,
		Password: password,
	}

	conn, err = Connect(conf)

	return
}
