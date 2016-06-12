package lease

import (
	"net"
	"testing"
)

func TestGenerate(t *testing.T) {
	testCIDRv4 := "192.168.1.0/24"
	testCIDRv6 := "1234::1222:0/16"
	testID := 5

	expIPv4 := net.ParseIP("192.168.1.5")
	expIPv6 := net.ParseIP("1234::1222:5")

	ip, err := Generate(testCIDRv4, testID)
	if err != nil {
		t.Errorf("%s", err.Error())
	}

	if !ip.Equal(expIPv4) {
		t.Errorf("unexpected ip, got: %v, wanted: %v", ip, expIPv4)
	}

	ip, err = Generate(testCIDRv6, testID)
	if err != nil {
		t.Errorf("%s", err.Error())
	}

	if !ip.Equal(expIPv6) {
		t.Errorf("unexpected ip, got: %v, wanted: %v", ip, expIPv6)
	}
}
