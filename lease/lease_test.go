package lease

import (
	"net"
	"testing"
)

var generateTests = []struct {
	cidr string
	id   uint
	ip   net.IP
	err  bool
}{
	{"192.168.1.0/24", 42, net.ParseIP("192.168.1.42"), false},
	{"192.168.1.0/24", 256, net.ParseIP("192.168.2.0"), true},
	{"12.21.0.0/16", 65536, net.ParseIP("12.22.0.0"), true},
	{"1234::1222:0/16", 423411, net.ParseIP("1234::1228:75f3"), false},
	{"1234::1222:0/120", 256, net.ParseIP("1234::1222:100"), true},
	{"3214:1261:afb2::0/96", 4294967296, net.ParseIP("3214:1261:afb2::1:0:0"), true},
}

func TestGenerate(t *testing.T) {
	for row, tests := range generateTests {
		ip, err := Generate(tests.cidr, tests.id)

		if !ip.Equal(tests.ip) {
			t.Errorf("Row: %d returned unexpected IP, got: %v, wanted: %v", row, ip, tests.ip)
		}

		if err != nil && !tests.err {
			t.Errorf("Row: %d returned unexpected error: %s", row, err.Error())
		}

		if err == nil && tests.err {
			t.Errorf("Row: %d expected error but got %v", row, err)
		}
	}
}
