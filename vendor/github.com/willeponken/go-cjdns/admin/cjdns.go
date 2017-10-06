// Package admin provides easy methods to the cjdns admin interface
package admin

import (
	"encoding/json"

	"io/ioutil"
	"math/rand"
	"net"
	"os/user"
	"regexp"
	"sync"
	"time"

	"github.com/ehmry/go-bencode"
)

type CjdnsAdminConfig struct {
	Addr     string `json:"addr"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	Config   string `json:"config,omitempty"`
}

// Conn is an object for interacting with the CJDNS administration port
type Conn struct {
	password   string
	addr       *net.UDPAddr
	enc        *bencode.Encoder
	Conn       *net.UDPConn
	mu         sync.Mutex
	queries    chan *request
	responses  map[string]chan *packet
	logStreams map[int]chan<- *LogMessage
}

func Connect(config *CjdnsAdminConfig) (admin *Conn, err error) {
	if config == nil {
		config = new(CjdnsAdminConfig)
		u, err := user.Current()
		if err != nil {
			return nil, err
		}

		rawFile, err := ioutil.ReadFile(u.HomeDir + "/.cjdnsadmin")
		if err != nil {
			return nil, err
		}

		raw, err := stripComments(rawFile)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(raw, &config)
		if err != nil {
			return nil, err
		}
	}

	addr := &net.UDPAddr{
		IP:   net.ParseIP(config.Addr),
		Port: config.Port,
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}

	admin = &Conn{
		password:  config.Password,
		addr:      addr,
		Conn:      conn,
		queries:   make(chan *request),
		responses: make(map[string]chan *packet),
	}

	go admin.readFromConn()
	go admin.writeToConn()
	return admin, err
}

const (
	readerChanSize       = 10
	socketReaderChanSize = 100
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}
func stripComments(b []byte) ([]byte, error) {
	regComment, err := regexp.Compile("(?s)//.*?\n|/\\*.*?\\*/")
	if err != nil {
		return nil, err
	}
	out := regComment.ReplaceAllLiteral(b, nil)
	return out, nil
}
