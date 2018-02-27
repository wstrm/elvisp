package database

import (
	"reflect"
	"testing"
)

var binUint64Tests = []struct {
	v uint64
	b []byte
}{
	{1, []byte{0, 0, 0, 0, 0, 0, 0, 1}},
	{18446744073709551616 - 1, []byte{255, 255, 255, 255, 255, 255, 255, 255}},
}

func TestUint64ToBin(t *testing.T) {
	for row, test := range binUint64Tests {
		b := uint64ToBin(test.v)

		if !reflect.DeepEqual(b, test.b) {
			t.Errorf("Row: %d returned unexpected binary, got: %v, wanted: %v", row, b, test.b)
		}
	}
}

func TestBinToUint64(t *testing.T) {
	for row, test := range binUint64Tests {
		v := binToUint64(test.b)

		if v != test.v {
			t.Errorf("Row: %d returned unexpected number, got: %d, wanted: %d", row, v, test.v)
		}
	}
}

func TestBinToUint64_invalidLength(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Invalid length not equal to 8 should panic")
		}
	}()

	binToUint64(make([]byte, 1)) // should panic, must be equal to 8
	binToUint64(make([]byte, 9)) // same as above
}
