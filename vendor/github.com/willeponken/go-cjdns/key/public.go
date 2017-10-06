package key

import (
	"bytes"
	"crypto/sha512"
	"errors"
	"net"
)

var (
	encodeAlphabet = [32]byte{
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'b', 'c', 'd', 'f', 'g', 'h', 'j', 'k', 'l', 'm',
		'n', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	}

	decodeAlphabet = func() [256]byte {
		var ascii [256]byte
		for i := range ascii {
			ascii[i] = 255
		}
		for i, b := range encodeAlphabet {
			ascii[b] = byte(i)
		}
		return ascii
	}()

	ErrInvalidPubKey = errors.New("Invalid public key supplied")

	suffix = []byte{'.', 'k'}
)

// Hashes the supplied slice twice and return the resulting byte slice.
func hashTwice(b []byte) []byte {
	var ip []byte
	h := sha512.New()
	h.Write(b[:])
	ip = h.Sum(ip[:0])
	h.Reset()
	h.Write(ip)
	ip = h.Sum(ip[:0])[0:16]
	return ip
}

// Represents a cjdns public key.
type Public [32]byte

// encodePublic encodes a key array to a byte slice using the cjdns key alphabet.
// dst should have a length of 52 or more.
func encodePublic(dst []byte, k *Public) {
	var wide, bits uint
	in := k[:]
	var i int
	for len(in) > 0 {
		// Add the 8 bits of data from the next `in` byte above the existing bits
		wide, in, bits = wide|uint(in[0])<<bits, in[1:], bits+8
		for bits > 5 {
			dst[i] = encodeAlphabet[int(wide&0x1F)]
			i++
			wide >>= 5
			bits -= 5
		}
	}
	/// If it wasn't a precise multiple of 40 bits, add some padding based on the remaining bits
	if bits > 0 {
		dst[51] = encodeAlphabet[int(wide)]
	}
}

// decodePublic decodes a byte slice to a Public.
func decodePublic(p []byte) (key *Public, err error) {
	key = new(Public)
	var wide, bits uint

	var i int
	for len(p) > 0 && p[0] != '=' {
		// Add the 5 bits of data corresponding to the next `in` character above existing bits
		wide, p, bits = wide|uint(decodeAlphabet[int(p[0])])<<bits, p[1:], bits+5
		if bits >= 8 {
			// Remove the least significant 8 bits of data and add it to out
			key[i] = byte(wide)
			i++
			wide >>= 8
			bits -= 8
		}
	}

	// If there was padding, there will be bits left, but they should be zero
	if wide != 0 {
		err = ErrInvalidPubKey
		return
	}

	// Check the key for validitiy
	if !key.Valid() {
		err = ErrInvalidPubKey
	}
	return
}

// Takes the string representation of a public key and returns a new Public
func DecodePublic(key string) (*Public, error) {
	if len(key) < 52 {
		return nil, ErrInvalidPubKey
	}
	if len(key) > 52 {
		key = key[:52]
	}
	return decodePublic([]byte(key))
}

// Returns true if k is a valid public key.
func (k *Public) Valid() bool {
	// It's a valid key if the IP address begins with FC
	v := hashTwice(k[:])
	return v[0] == 0xFC
}

// Returns the public key in base32 format ending with .k
func (k *Public) String() string {
	if k.isZero() {
		return ""
	}
	b := make([]byte, 54)
	encodePublic(b, k)
	copy(b[52:], suffix)
	return string(b)
}

func (k *Public) isZero() bool {
	for _, c := range k {
		if c != 0 {
			return false
		}
	}
	return true
}

// Implements the encoding.TextMarshaler interface
func (k *Public) MarshalText() ([]byte, error) {
	if k.isZero() {
		return nil, errors.New("MarshalText called on zero key")
	}
	text := make([]byte, 54)
	encodePublic(text, k)
	copy(text[52:], suffix)
	return text, nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (k *Public) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		k = nil
		return nil
	}

	if len(text) < 52 {
		return ErrInvalidPubKey
	}
	if len(text) > 52 {
		text = text[:52]

	}

	key, err := decodePublic(text)
	*k = *key
	return err
}

// Returns the cjdns IPv6 address of the key.
func (k *Public) IP() net.IP {
	return k.makeIPv6()
}

// Returns a string containing the IPv6 address for the public key
func (k *Public) makeIPv6() net.IP {
	out := hashTwice(k[:])
	return net.IP(out)
}

// Equal returns true if key and x and the same public key.
func (key *Public) Equal(x *Public) bool {
	return bytes.Equal(key[:], x[:])
}
