package lease

import (
	"errors"
	"math"
	"net"
)

// ipToUint32 takes a ip address represented as a slice of bytes and converts it into uint32.
func ipToUint32(ip []byte) uint32 {
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
}

// ipToUint128 takes a ip address represented as a slice of bytes and converts it into a pair of uint64.
func ipToUint128(ip []byte) (a, b uint64) {
	j := uint(56)
	for i := 0; i < 8; i++ {
		a = a | uint64(ip[i])<<j
		b = b | uint64(ip[i+8])<<j

		j = j - uint(8)
	}

	return
}

// uint128ToIP takes a pair of uint64 and converts the first integer to the first 8 blocks of IPv6 and the second integer to the next 8 blocks.
func uint128ToIP(a, b uint64) net.IP {
	ip := make(net.IP, net.IPv6len)

	j := uint(56)
	for i := 0; i < 8; i++ {
		ip[i] = byte(a >> j)
		ip[i+8] = byte(b >> j)

		j = j - uint(8)
	}

	return ip
}

// uint32ToIP converts a uint32 to a IP address.
func uint32ToIP(i uint32) net.IP {
	return net.IPv4(byte(i>>24), byte(i>>16), byte(i>>8), byte(i))
}

// uint128Add takes a pair of uint64 and adds them with another uint64 (as if the pair was a uint128).
func uint128Add(a, b, i uint64) (uint64, uint64) {
	if math.MaxUint64-b < i {
		a++
	}

	b += i

	return a, b
}

// withinNetwork checks if the generated IP address fits within the network specified.
func withinNetwork(network *net.IPNet, ip net.IP) error {
	if !network.Contains(ip) {
		return errors.New("IP address is outside of available network")
	}

	return nil
}

// CIDR holds a start address, the allowed network and a ID to add to the start IP.
type CIDR struct {
	Start   net.IP
	Network *net.IPNet
}

// ParseCIDR acts as a wrapper for net.ParseCIDR and populates a lease.CIDR struct.
func ParseCIDR(cidr string) (c CIDR, err error) {
	c.Start, c.Network, err = net.ParseCIDR(cidr)
	return
}

// Generate takes the CIDR (both IPv4 and IPv6 is supported) and a ID (which is used to increment the IP address from the CIDR). Then the incremented IP address is returned.
func Generate(cidr CIDR, id uint64) (ip net.IP, err error) {
	start := cidr.Start
	network := cidr.Network

	// Is the IP IPv4?
	if s := start.To4(); s != nil {
		ip = uint32ToIP(ipToUint32(s) + uint32(id))
		err = withinNetwork(network, ip)
		return
	}

	// Or is the IP IPv6?
	if s := start.To16(); s != nil {
		a, b := ipToUint128(s)
		ip = uint128ToIP(uint128Add(a, b, uint64(id)))
		err = withinNetwork(network, ip)
		return
	}

	// If ip.To16() returns nil, the IP has an invalid length.
	err = errors.New("Invalid length of IP address")
	return
}
