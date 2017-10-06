package bencode

import (
	"errors"
)

// Marshaler is the interface implemented by objects that
// can marshal themselves into valid bencode.
type Marshaler interface {
	MarshalBencode() ([]byte, error)
}

// Unmarshaler is the interface implemented by objects that
// can unmarshal themselves from bencode.
type Unmarshaler interface {
	UnmarshalBencode([]byte) error
}

// RawMessage is a raw encoded bencode object.
// It is intedended to delay decoding or precomute an encoding.
type RawMessage []byte

// MarshalText returns *m as the bencode encoding of m.
func (m *RawMessage) MarshalBencode() ([]byte, error) {
	if m == nil {
		return []byte{'0', ':'}, nil
	}
	return *m, nil
}

// UnmarshalText sets *m to a copy of data.
func (m *RawMessage) UnmarshalBencode(text []byte) error {
	if m == nil {
		return errors.New("bencode.RawMessage: UnmarshalText on nil pointer")
	}
	*m = append((*m)[0:0], text...)
	return nil
}

var _ Marshaler = (*RawMessage)(nil)
var _ Unmarshaler = (*RawMessage)(nil)
