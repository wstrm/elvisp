package admin

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/ehmry/go-bencode"
)

type response struct {
	Txid     string
	StreamId string
	Error    string
}

type packet struct {
	buf *bytes.Reader
	dec *bencode.Decoder
	err error
}

func (p *packet) Decode(v interface{}) error {
	if p.err != nil {
		return p.err
	}
	p.buf.Seek(0, 0)
	return p.dec.Decode(v)
}

func (a *Conn) readFromConn() {
	var (
		b    = make([]byte, 69632)
		r    = new(response)
		n    int
		err  error
		pack *packet
	)
	for {
		// read into local buffer
		n, err = a.Conn.Read(b)
		if err != nil {
			// this probably will mess more stuff up down the line.
			pack = &packet{err: err}
			for _, c := range a.responses {
				c <- pack
			}
		}

		// copy to new buffer
		newB := make([]byte, n)
		copy(newB, b[:n])
		// encapsulate
		br := bytes.NewReader(newB)
		pack = &packet{buf: br, dec: bencode.NewDecoder(br)}

		if err = pack.dec.Decode(r); err != nil {
			pack.err = err
		} else {
			if r.Error != "none" && r.Error != "" {
				pack.err = errors.New(r.Error)
				r.Error = ""
			}
		}
		if c, ok := a.responses[r.Txid]; ok {
			c <- pack
			continue
		}

		var id int
		fmt.Sprintf(r.StreamId, "%x", &id)

		if c, ok := a.logStreams[id]; ok {
			m := new(LogMessage)
			if err = bencode.Unmarshal(b, m); err == nil {
				// this runs in it's own go routine because we can't
				// be sure this user supplied channel is buffered.
				go func() { c <- m }()
			}
		}
	}
}

func (c *Conn) writeToConn() {
	c.enc = bencode.NewEncoder(c.Conn)

	var req *request
	var err error
	for {
		select {
		case req = <-c.queries:
			//b, _ := bencode.Marshal(req)
			//fmt.Printf("-> %q\n", b)
			if err = c.enc.Encode(req); err != nil {
				rc := c.responses[req.Txid]
				rc <- &packet{err: errors.New("Failed to query cjdns, " + err.Error())}
			}

		case <-time.After(8 * time.Second):
			// Ping CJDNS to keep log streams going.
			go c.Ping()

		}
	}
}

var errorSocketClosed = errors.New("Socket closed")

func init() {
	rand.Seed(time.Now().UnixNano())
}

func newTxid() string {
	return fmt.Sprintf("%x", rand.Uint32())
}

type request struct {
	Q      string      `bencode:"q"`
	AQ     string      `bencode:"aq,omitempty"`
	Cookie string      `bencode:"cookie,omitempty"`
	Hash   string      `bencode:"hash,omitempty"`
	Args   interface{} `bencode:"args,omitempty"`
	Txid   string      `bencode:"txid"`
}

func (a *Conn) sendCmd(req *request) (response *packet, err error) {
	req.Txid = newTxid()

	//create the channel to receive data back on
	c := make(chan *packet, 1)
	a.mu.Lock()
	a.responses[req.Txid] = c
	a.mu.Unlock()

	// remove channel from map no matter how this function exits.
	defer func() {
		a.mu.Lock()
		delete(a.responses, req.Txid)
		a.mu.Unlock()
	}()

	if req.AQ != "" {
		// it's an authenticated request
		var cookie string
		cookie, err = a.cookie()
		if err != nil {
			return
		}

		h := sha256.New()
		io.WriteString(h, a.password)
		io.WriteString(h, cookie)

		req.Q = "auth"
		req.Cookie = cookie
		req.Hash = hex.EncodeToString(h.Sum(nil))

		h.Reset()
		if err = bencode.NewEncoder(h).Encode(req); err != nil {
			return
		}

		req.Hash = hex.EncodeToString(h.Sum(nil))
	}

	// Send the query
	a.queries <- req

	// wait for the response
	var ok bool
	response, ok = <-c
	if !ok {
		err = errorSocketClosed
	}
	// check for this error now
	if response.err != nil {
		err = response.err
	}
	return
}

// cookie requests a cookie from CJDNS
func (a *Conn) cookie() (string, error) {
	pack, err := a.sendCmd(&request{Q: "cookie"})
	if err != nil {
		return "", err
	}
	r := new(struct {
		Cookie string `cookie`
	})
	return r.Cookie, pack.Decode(r)
}

func (a *Conn) registerLogChan(streamId string, c chan<- *LogMessage) {
	a.mu.Lock()
	var id int
	fmt.Sprintf(streamId, "%x", &id)
	if a.logStreams == nil {
		a.logStreams = make(map[int]chan<- *LogMessage)
	}
	a.logStreams[id] = c
	a.mu.Unlock()
}
