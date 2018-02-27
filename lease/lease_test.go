package lease_test

import (
	"net"
	"testing"

	"github.com/willeponken/elvisp/lease"
)

func TestGenerate(t *testing.T) {
	var generateTests = []struct {
		cidr string
		id   uint64
		ip   net.IP
		err  bool
	}{
		{"192.168.1.0/24", 42, net.ParseIP("192.168.1.42"), false},
		{"192.168.1.0/24", 256, net.ParseIP("192.168.2.0"), true},
		{"12.21.0.0/16", 65536, net.ParseIP("12.22.0.0"), true},
		{"1234::1222:0/16", 423411, net.ParseIP("1234::1228:75f3"), false},
		{"1234::1222:0/120", 256, net.ParseIP("1234::1222:100"), true},
		{"3214:1261:afb2::0/96", 4294967296, net.ParseIP("3214:1261:afb2::1:0:0"), true},
		{"3214:1261:afb2:ffff:ffff:ffff:ffff:ffff/96", 4294967296, net.ParseIP("3214:1261:afb3::ffff:ffff"), true},
	}

	for row, tests := range generateTests {
		cidr, err := lease.ParseCIDR(tests.cidr)
		if err != nil {
			t.Errorf("Row: %d returned unexpected error: %v", row, err)
		}

		ip, err := lease.Generate(cidr, tests.id)

		if !ip.Equal(tests.ip) {
			t.Errorf("Row: %d returned unexpected IP, got: %v, wanted: %v", row, ip, tests.ip)
		}

		if err != nil && !tests.err {
			t.Errorf("Row: %d returned unexpected error: %v", row, err)
		}

		if err == nil && tests.err {
			t.Errorf("Row: %d expected error but got %v", row, err)
		}
	}
}

func TestGenerate_emptyCIDR(t *testing.T) {
	_, err := lease.Generate(lease.CIDR{}, 0)
	if err == nil {
		t.Errorf("Expected error but got %v", err)
	}
}

func TestCIDR(t *testing.T) {

	var cidrTests = []struct {
		cidr string
		err  bool
	}{
		{"192.168.1.0/24", false},
		{"1234::1222:0/16", false},
		{"192.168.1.0/128", true},
		{"1234::1222:0/512", true},
		{"wow:such:an:invalid:address::yes/lol", true},
	}

	for row, tests := range cidrTests {
		_, err := lease.ParseCIDR(tests.cidr) // No reason to test the net.ParseCIDR() function, therefore ignoring first return

		if err != nil && !tests.err {
			t.Errorf("Row: %d returned unexpected error: %v", row, err)
		}

		if err == nil && tests.err {
			t.Errorf("Row: %d expected error but got %v", row, err)
		}
	}
}
