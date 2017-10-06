package admin

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"sort"
	"strconv"
)

var (
	ErrNotInTable = errors.New("Node not in local routing table")
	ErrParseIP    = errors.New("CJDNS node failed to parse IP")
)

const magicalLinkConstant = 5366870 //Determined by cjd way back in the dark ages.

type Route struct {
	IP      *net.IP
	Link    Link
	Path    *Path
	Version int
}

type (
	Link uint32
	Path uint64
)

func (l Link) String() string {
	return strconv.FormatUint(uint64(l)/magicalLinkConstant, 10)
}

func (p Path) String() string {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(p))
	text := make([]byte, 19)
	hex.Encode(text, b)
	copy(text[15:19], text[12:16])
	text[14] = '.'
	copy(text[10:14], text[8:12])
	text[9] = '.'
	copy(text[5:9], text[4:8])
	text[4] = '.'
	return string(text)
}

func (p Path) MarshalText() (text []byte, err error) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(p))
	text = make([]byte, 19)
	hex.Encode(text, b)
	copy(text[15:19], text[12:16])
	text[14] = '.'
	copy(text[10:14], text[8:12])
	text[9] = '.'
	copy(text[5:9], text[4:8])
	text[4] = '.'
	return
}

func ParsePath(path string) Path {
	if len(path) != 19 {
		return Path(0)
	}
	text := []byte(path)
	copy(text[4:8], text[5:9])
	copy(text[8:12], text[10:14])
	copy(text[12:16], text[15:19])
	text = text[:16]

	b := make([]byte, 16)

	_, err := hex.Decode(b, text)
	if err != nil {
		return Path(0)
	}
	return Path(binary.BigEndian.Uint64(b))
}

func (p *Path) UnmarshalText(text []byte) error {
	if len(text) != 19 {
		return fmt.Errorf("bad path %q", text)
	}
	copy(text[4:8], text[5:9])
	copy(text[8:12], text[10:14])
	copy(text[12:16], text[15:19])
	text = text[:16]

	b := make([]byte, 16)

	if _, err := hex.Decode(b, text); err != nil {
		return err
	}
	*p = Path(binary.BigEndian.Uint64(b))
	return nil
}

type Routes []*Route

func (rs Routes) Len() int      { return len(rs) }
func (rs Routes) Swap(i, j int) { rs[i], rs[j] = rs[j], rs[i] }

// SortByPath sorts Routes by link quality.
func (r Routes) SortByPath() {
	if len(r) < 2 {
		return
	}
	sort.Sort(byPath{r})
}

type byPath struct{ Routes }

func (s byPath) Less(i, j int) bool { return *s.Routes[i].Path < *s.Routes[j].Path }

// SortByQuality sorts Routes by link quality.
func (r Routes) SortByQuality() {
	if len(r) < 2 {
		return
	}
	sort.Sort(byQuality{r})
}

type byQuality struct{ Routes }

func (s byQuality) Less(i, j int) bool { return s.Routes[i].Link > s.Routes[j].Link }

// Log base 2 of a Path
func log2x64(p Path) (out Path) {
	// return Path(math.Log2(float64(p)))
	// Second method is faster
	for p > 1 {
		p >>= 1
		out++
	}
	return
}

// IsBehind returns true if midpath is routed through p.
func (p Path) IsBehind(midPath Path) bool {
	if midPath > p {
		return false
	}
	mask := ^Path(0) >> (64 - log2x64(midPath))
	return (p & mask) == (midPath & mask)
}

// IsOneHop Returns true if midPath is one hop from p.
// WARNING: this depends on implementation quirks of the router and will be broken in the future.
func (p Path) IsOneHop(node Path) bool {
	// NOTE: This may have false positives which isBehind() will remove.
	if !p.IsBehind(node) {
		return false
	}

	var c Path
	switch {
	case p > node:
		c = p >> log2x64(node)
	case p < node:
		c = node >> log2x64(p)
	default:
		return false
	}

	// The "why" is here:
	// http://gitboria.com/cjd/cjdns/tree/master/switch/NumberCompress.h#L143
	//c := destination >> log2x64(p)
	if c&1 != 0 {
		return log2x64(c) == 4
	}
	if c&3 != 0 {
		return log2x64(c) == 7
	}
	return log2x64(c) == 10
}

// Hops returns a Routes object representing a set of hops to a path
func (rs Routes) Hops(destination Path) (hops Routes) {
	for _, r := range rs {
		if destination.IsBehind(*r.Path) {
			hops = append(hops, r)
		}
	}
	return
}

// NodeStore_dumpTable will return cjdns's routing table.
func (c *Conn) NodeStore_dumpTable() (routingTable Routes, err error) {
	var (
		args = new(struct {
			Page int `bencode:"page"`
		})
		req = request{Q: "NodeStore_dumpTable", Args: args}

		resp = new(struct {
			More bool
			// skip this for now just to get the function to work
			RoutingTable Routes
		})

		pack *packet
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

	return resp.RoutingTable, err
}

type Parent struct {
	IP               string
	ParentChildLabel string
}

type EncodingScheme struct {
	BitCount  int
	Prefix    string
	PrefixLen int
}

type StoreNode struct {
	RouteLabel      string
	BestParent      Parent
	EncodingScheme  []*EncodingScheme
	Key             string
	LinkCount       int
	ProtocolVersion int
	Reach           int
}

func (n *StoreNode) String() string {
	return n.RouteLabel
}

type StoreLink struct {
	LinkState                     int
	Parent, Child                 string
	CannonicalLabel               string
	InverseLinkEncodingFormNumber int
	IsOneHop                      bool
}

func (c *Conn) NodeStore_getLink(parent string, link int) (l *StoreLink, err error) {
	req := request{
		AQ: "NodeStore_getLink",
		Args: &struct {
			Parent string `bencode:"parent"`
			Link   int    `bencode:"linkNum"`
		}{parent, link},
	}

	var pack *packet
	l = new(StoreLink)
	if pack, err = c.sendCmd(&req); err == nil {
		err = pack.Decode(&struct{ Result *StoreLink }{l})
	}
	return
}

func (c *Conn) NodeStore_nodeForAddr(ip string) (n *StoreNode, err error) {
	var (
		req  = request{AQ: "NodeStore_nodeForAddr"}
		pack *packet
	)

	if ip != "" {
		req.Args = &struct {
			Ip string `bencode:"ip"`
		}{ip}
	}

	n = new(StoreNode)
	if pack, err = c.sendCmd(&req); err == nil {
		err = pack.Decode(&struct{ Result *StoreNode }{n})
	}
	if err != nil && err.Error() == "parse_ip" {
		err = ErrParseIP
	}
	if err == nil && n.RouteLabel == "" {
		n = nil
		err = ErrNotInTable
	}
	return
}

// Peers returns a Routes object representing routes
// directly connected to a given IP.
func (rs Routes) Peers(ip net.IP) (peerRoutes Routes) {
	pm := make(map[string]*Route)

	for _, target := range rs {
		if !target.IP.Equal(ip) {
			continue
		}

		for _, node := range rs {
			if node.Path.IsOneHop(*target.Path) || target.Path.IsOneHop(*node.Path) {
				nodeIp := node.IP.String()
				if prev, ok := pm[nodeIp]; !ok || *node.Path < *prev.Path {
					// route has not be stored or it is shorter than the previous
					pm[nodeIp] = node
				}
			}
		}
	}

	peerRoutes = make(Routes, len(pm))
	var i int
	for _, route := range pm {
		peerRoutes[i] = route
		i++
	}
	return
}
