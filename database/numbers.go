package database

import (
	"encoding/binary"
	"fmt"
)

// uint64ToBin returns an 8-byte big endian representation of v.
func uint64ToBin(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

// binToUint64 takes an 8-byte big endian and converts it into a uint64.
func binToUint64(v []byte) uint64 {
	if len(v) != 8 {
		panic(fmt.Sprintf("invalid length of binary: %d", len(v)))
	}

	return binary.BigEndian.Uint64(v)
}
