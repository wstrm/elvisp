package key

import (
	"crypto/rand"
	"encoding/hex"
	"golang.org/x/crypto/curve25519"
)

type (
	// Represents a cjdns private key.
	Private [32]byte
)

// DeocodePrivate returns a private key from a hex encoded string.
func DecodePrivate(s string) (key *Private, err error) {
	key = new(Private)
	_, err = hex.Decode(key[:], []byte(s))
	return
}

// Generate creates a new random private key.
func Generate() (key *Private) {
	key = new(Private)
	for {
		rand.Read(key[:])

		key[0] &= 248
		key[31] &= 127
		key[31] |= 64

		if key.Valid() {
			return
		}
	}
}

// Returns true if the private key is valid.
func (k *Private) Valid() bool { return k.Pubkey().Valid() }

// Returns the public key in base32 format.
func (k *Private) String() string {
	return hex.EncodeToString(k[:])
}

// Implements the encoding.TextMarshaler interface
func (k *Private) MarshalText() ([]byte, error) {
	out := make([]byte, 64)
	hex.Encode(out, k[:])
	return out, nil
}

// Implements the encoding.TextUnmarshaler interface
func (k *Private) UnmarshalText(text []byte) (err error) {
	if len(text) == 0 {
		k = nil
		return
	}

	key := Private{}
	_, err = hex.Decode(key[:], text)
	*k = key
	return
}

// Pubkey returns the associated public key for the supplied private key.
func (k *Private) Pubkey() *Public {
	var pub [32]byte
	priv := [32]byte(*k)

	// Performs ScalarBaseMult on the supplied private key, returning the public key
	curve25519.ScalarBaseMult(&pub, &priv)
	public := Public(pub)
	return &public
}
