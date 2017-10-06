package bencode

import (
	"bytes"
	"encoding"
	"errors"
	"io"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

// A Decoder decodes bencoded data from a stream.
type Decoder struct {
	r    io.Reader
	buf  []byte
	d    decodeState
	scan scanner
	err  error
}

// NewDecoder returns a new decoder that decodes from r.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

// Decode decodes data from the wrapped stream into v.
func (dec *Decoder) Decode(v interface{}) error {
	n, err := dec.readValue()
	if err != nil {
		return err
	}

	// Don't save err from unmarshal into dec.err:
	// the connection is still usable since we read a complete bencode
	// object from it before the error happened.
	dec.d.init(dec.buf[0:n])
	err = dec.d.unmarshal(v)

	// Slide rest of data down.
	rest := copy(dec.buf, dec.buf[n:])
	dec.buf = dec.buf[0:rest]

	return err
}

// Buffered returns a read of the data remaining in the Decoder's buffer.
// The reader is valid until the next call to Decode.
func (dec *Decoder) Buffered() io.Reader {
	return bytes.NewReader(dec.buf)
}

// readValue reads a bencode value into dec.buf.
// It returns the length of the encoding.
func (dec *Decoder) readValue() (int, error) {
	dec.scan.reset()

	var scanp, op, n int
	var err error
Input:
	for {
		for scanp < len(dec.buf) {
			op = dec.scan.step(&dec.scan, int(dec.buf[scanp]))
			scanp++
			if op >= 0 {
				dec.scan.bytes += int64(op)
				scanp += op
				if dec.scan.endTop {
					break Input
				}
			} else {
				dec.scan.bytes++
				switch op {
				case scanEnd:
					break Input

				case scanError:
					dec.err = dec.scan.err
					return 0, dec.scan.err
				}
			}
		}

		// Did the last read have an error?
		// Delayed until now to allow buffer scan.
		if err != nil {
			if err == io.EOF && dec.scan.endTop {
				err = nil
				break Input
			}
			dec.err = err
			return 0, err
		}

		// Make room to read more into the buffer.
		const minRead = 512
		if cap(dec.buf)-len(dec.buf) < minRead {
			newBuf := make([]byte, len(dec.buf), 2*cap(dec.buf)+minRead)
			copy(newBuf, dec.buf)
			dec.buf = newBuf
		}

		// Read. Delay error for the next interation (after scan).
		n, err = dec.r.Read(dec.buf[len(dec.buf):cap(dec.buf)])
		dec.buf = dec.buf[0 : len(dec.buf)+n]
	}
	return scanp, nil
}

// An InvalidUnmarshalError describes an invalid argument passed to Unmarshal.
// (The argument to Unmarshal must be a non-nil pointer.)
type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "bencode: Unmarshal(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "bencode: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "bencode: Unmarshal(nil " + e.Type.String() + ")"
}

func (d *decodeState) unmarshal(v interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
		}
	}()

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}

	d.scan.reset()
	// We decode rv not rv.Elem because the Unmarshaler interface
	// test must be applied at the top level of the value.
	d.value(rv)
	return d.savedError
}

// An UnmarshalTypeError describes a bencode value that was
// not appropriate for a value of a specific Go type.
type UnmarshalTypeError struct {
	Value string       // description of bencode value
	Type  reflect.Type // type of Go value it could not be assigned to
}

func (e *UnmarshalTypeError) Error() string {
	return "bencode: cannot unmarshal " + e.Value + " into Go value of type " + e.Type.String()
}

type decodeState struct {
	data       []byte
	off        int // read offset in data
	scan       scanner
	savedError error
}

// errPhase is used for errors that should not happen unless
// there is a bug in the bencode decoder or something is editing
// the data slice while the decoder executes.
var errPhase = errors.New("bencoder decoder out of sync - data changing underfoot?")

func (d *decodeState) init(data []byte) *decodeState {
	d.data = data
	d.off = 0
	d.savedError = nil
	return d
}

// error aborts the decoding by panicking with err.
func (d *decodeState) error(err error) {
	panic(err)
}

// skip reads d.data with a fresh scanner, skimming over the next value
func (d *decodeState) skip() {
	var skipScan scanner
	skipScan.reset()
	skipScan.step = d.scan.step
	d.scan.step = stateEndValue

	var op int
	for {
		op = skipScan.step(&skipScan, int(d.data[d.off]))
		d.off++
		if op >= 0 {
			d.off += op
			if skipScan.endTop {
				return
			}
		} else {
			switch op {
			case scanEnd:
				return
			case scanError:
				d.error(skipScan.err)
			}
		}
	}
}

var (
	unmarshalerType     = reflect.TypeOf(new(Unmarshaler)).Elem()
	textUnmarshalerType = reflect.TypeOf(new(encoding.TextUnmarshaler)).Elem()
)

// value decodes the next item from d.data[d.off:] into v, updating d.off.
func (d *decodeState) value(v reflect.Value) {
	if !v.IsValid() {
		d.skip()
		return
	}

	if v.Type().Implements(unmarshalerType) {
		d.unmarshaler(v)
		return
	}

	op := d.scan.step(&d.scan, int(d.data[d.off]))
	d.off++

	switch op {
	case scanBeginStringLen:
		d.string(v)

	case scanBeginInteger:
		d.integer(v)

	case scanBeginList:
		d.list(v)

	case scanBeginDict:
		d.dict(v)

	case scanEnd:
		return

	case scanError:
		d.error(d.scan.err)

	default:
		d.error(errPhase)
	}
}

// valueInterface is like value but returns interface{}
func (d *decodeState) valueInterface() (x interface{}) {
	c := int(d.data[d.off])
	d.off++

	switch d.scan.step(&d.scan, c) {

	case scanBeginStringLen:
		x = d.readString()

	case scanBeginInteger:
		x = d.integerInterface()

	case scanBeginList:
		x = d.listInterface()

	case scanBeginDict:
		x = d.dictInterface()

	case scanError:
		d.error(d.scan.err)
	default:
		d.error(errPhase)
	}
	return
}

type integerUnmarshalError struct {
	err error
}

// readInteger reads an integer value from d.data[d.off:],
//  and returns the binary form
func (d *decodeState) readInteger() []byte {
	i := d.off
	var c int
Read:
	for {
		c = int(d.data[d.off])
		d.off++
		switch d.scan.step(&d.scan, c) {
		case scanParseInteger:
			continue
		case scanEndInteger, scanEnd:
			break Read
		case scanError:
			d.error(d.scan.err)
		default:
			d.error(errPhase)
		}
	}
	return d.data[i : d.off-1]
}

// integer consumes an integer from d.data[d.off:], decoding into the value v.
func (d *decodeState) integer(v reflect.Value) {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	k := v.Kind()
	if k == reflect.Interface {
		v.Set(reflect.ValueOf(d.integerInterface()))
		return
	}

	s := string(d.readInteger())

	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			d.error(err)
		}
		v.SetInt(n)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			d.error(err)
		}
		v.SetUint(n)

	case reflect.Float32, reflect.Float64:
		n, err := strconv.ParseFloat(s, v.Type().Bits())
		if err != nil || v.OverflowFloat(n) {
			d.error(&UnmarshalTypeError{"integer " + s, v.Type()})
		}
		v.SetFloat(n)

	case reflect.String:
		v.SetString(s)

	case reflect.Bool:
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			d.error(err)
		}
		v.SetBool(n != 0)

	default:
		d.error(&UnmarshalTypeError{"integer " + s, v.Type()})
	}
}

// intergerInterface consumes an integer from d.data[d.off:], and returns an interface{}.
func (d *decodeState) integerInterface() (x interface{}) {
	s := string(d.readInteger())

	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		d.error(err)
	}
	return n
}

func (d *decodeState) readString() []byte {
	var c, i, op int
Read:
	for {
		c = int(d.data[d.off])
		d.off++
		op = d.scan.step(&d.scan, c)
		if op >= 0 {
			i = d.off
			d.off += op
			break Read
		} else {
			switch op {
			case scanParseStringLen, scanParseString:
				continue
			case scanError:
				d.error(d.scan.err)
			default:
				d.error(errPhase)
			}
		}
	}
	return d.data[i:d.off]
}

// string consumes a string from d.data[d.off:], decoding into the value v.
func (d *decodeState) string(v reflect.Value) {
	for {
		if v.Type().Implements(textUnmarshalerType) {
			if v.Kind() == reflect.Ptr && v.IsNil() {
				v.Set(reflect.New(v.Type().Elem()))
			}

			//if u, ok := v.Interface().(encoding.TextUnmarshaler); ok {
			u := v.Interface().(encoding.TextUnmarshaler)
			if err := u.UnmarshalText(d.readString()); err != nil {
				d.error(err)
			}
			return
		}

		if v.Kind() != reflect.Ptr {
			break
		}
		v = v.Elem()
	}

	switch v.Kind() {
	default:
		d.error(&UnmarshalTypeError{"string", v.Type()})

	case reflect.Slice:
		if v.Type() != byteSliceType {
			d.error(&UnmarshalTypeError{"string", v.Type()})
		}

		v.SetBytes(d.readString())

	case reflect.String:
		v.SetString(string(d.readString()))

	case reflect.Interface:
		if v.NumMethod() != 0 {
			d.error(&UnmarshalTypeError{"string", v.Type()})
		}

		x := d.readString()
		v.Set(reflect.ValueOf(x))
	}
}

// list consumes a list from d.data[d.off-1:], decoding into the value v.
func (d *decodeState) list(v reflect.Value) {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Check type of target.
	switch v.Kind() {
	case reflect.Interface:
		if v.NumMethod() == 0 {
			// Decoding into nil interface? Switch to non-reflect code.
			x := d.listInterface()
			v.Set(reflect.ValueOf(x))
			return
		}
		// Otherwilse it's invalid
		fallthrough
	default:
		d.error(&UnmarshalTypeError{"list", v.Type()})

	case reflect.Array:
	case reflect.Slice:
	}

	var c, op int
	i := v.Len()
Read:
	for {
		c = int(d.data[d.off])
		d.off++
		op = d.scan.step(&d.scan, c)

		switch op {
		case scanEndList, scanEnd:
			break Read
		}

		// Get element of array, growing if necessary.
		if v.Kind() == reflect.Slice {
			// Grow slice if necessary
			if i >= v.Cap() {
				newcap := v.Cap() + v.Cap()/2
				if newcap < 4 {
					newcap = 4
				}
				newv := reflect.MakeSlice(v.Type(), v.Len(), newcap)
				reflect.Copy(newv, v)
				v.Set(newv)
			}
			if i >= v.Len() {
				v.SetLen(i + 1)
			}
		}

		var subv reflect.Value
		if i < v.Len() {
			// Decode into element.
			subv = v.Index(i)
		} else {
			// Ran out of fixed array: skip.
			subv = reflect.Value{}
		}

		for subv.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
			if subv.IsNil() {
				subv.Set(reflect.New(subv.Type().Elem()))
			}
			subv = subv.Elem()
		}
		switch op {
		case scanBeginStringLen:
			d.string(subv)

		case scanBeginInteger:
			d.integer(subv)

		case scanBeginList:
			d.list(subv)

		case scanBeginDict:
			d.dict(subv)

		case scanError:
			d.error(d.scan.err)
		default:
			d.error(errPhase)

		}
		i++
	}

	if i < v.Len() {
		if v.Kind() == reflect.Array {
			// Array. Zero the rest.
			z := reflect.Zero(v.Type().Elem())
			for ; i < v.Len(); i++ {
				v.Index(i).Set(z)
			}
		} else {
			v.SetLen(i)
		}
	}
	if i == 0 && v.Kind() == reflect.Slice {
		v.Set(reflect.MakeSlice(v.Type(), 0, 0))
	}
}

// listInterface is like list but returns []interface{}
func (d *decodeState) listInterface() interface{} {
	var (
		v = make([]interface{}, 0)
		x interface{}
		c int
	)
Read:
	for {
		c = int(d.data[d.off])
		d.off++

		switch op := d.scan.step(&d.scan, c); op {
		case scanEndList:
			break Read

		case scanBeginStringLen:
			x = d.readString()
		case scanBeginInteger:
			x = d.integerInterface()
		case scanBeginList:
			x = d.listInterface()
		case scanBeginDict:
			x = d.dictInterface()

		case scanError:
			d.error(d.scan.err)
		default:
			d.error(errPhase)
		}
		v = append(v, x)
	}
	return v
}

// dict consumes a dict from d.data[d.off:], decoding into the value v.
func (d *decodeState) dict(v reflect.Value) {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return
		}
		v = v.Elem()
	}

	if v.Kind() == reflect.Interface && v.NumMethod() == 0 {
		// Decoding into nil interface? Switch to non-reflect code.
		v.Set(reflect.ValueOf(d.dictInterface()))
		return
	}

	// Check type of target: struct or map[string]interface{}
	switch v.Kind() {
	case reflect.Map:
		// map must have string kind
		t := v.Type()
		if t.Key().Kind() != reflect.String {
			d.error(&UnmarshalTypeError{"dictionary", v.Type()})
		}
		if v.IsNil() {
			v.Set(reflect.MakeMap(t))
		}

	case reflect.Struct:

	default:
		d.error(&UnmarshalTypeError{"dictionary", v.Type()})
	}

	var mapElem reflect.Value
	var key string
	var c, op, p int
Read:
	for {
		// Read string key.
	ReadKey:
		for {
			c = int(d.data[d.off])
			d.off++
			op = d.scan.step(&d.scan, c)
			if op > 0 {
				p = d.off
				d.off += op
				break ReadKey
			} else {
				switch op {
				case scanEndDict, scanEnd:
					break Read
				case scanBeginKeyLen, scanParseKeyLen, scanParseKey:
				case scanEndKeyLen:
					p = d.off
				default:
					d.error(errPhase)
				}
			}
		}
		key = string(d.data[p:d.off])

		// Figure out field corresponding to key.
		var subv reflect.Value

		if v.Kind() == reflect.Map {
			elemType := v.Type().Elem()
			if !mapElem.IsValid() {
				mapElem = reflect.New(elemType).Elem()
			} else {
				mapElem.Set(reflect.Zero(elemType))
			}
			subv = mapElem
		} else {
			var f *field
			fields := cachedTypeFields(v.Type())
			for i := range fields {
				ff := &fields[i]
				if ff.name == key {
					f = ff
					break
				}
				if f == nil && strings.EqualFold(ff.name, key) {
					f = ff
				}
			}
			if f != nil {
				subv = v
				for _, i := range f.index {
					if subv.Kind() == reflect.Ptr {
						if subv.IsNil() {
							subv.Set(reflect.New(subv.Type().Elem()))
						}
						subv = subv.Elem()
					}
					subv = subv.Field(i)
				}
			}
		}

		// Read value.
		d.value(subv)

		// Write value back to map;
		// if using struct, subv points into struct already.
		if v.Kind() == reflect.Map {
			kv := reflect.ValueOf(key).Convert(v.Type().Key())
			v.SetMapIndex(kv, subv)
		}
	}
}

// dictInterface is like dict but returns a map[string]interface{}.
func (d *decodeState) dictInterface() map[string]interface{} {
	m := make(map[string]interface{})
	var key string
	var c, op, p int
Read:
	for {
	ReadKey:
		for {
			c = int(d.data[d.off])
			d.off++
			op = d.scan.step(&d.scan, c)
			if op > 0 {
				p = d.off
				d.off += op
				break ReadKey
			} else {
				switch op {
				case scanEndDict:
					break Read
				case scanBeginKeyLen, scanParseKeyLen, scanParseKey:
				case scanEndKeyLen:
					p = d.off
				default:
					d.error(errPhase)
				}
			}
		}
		key = string(d.data[p:d.off])
		m[key] = d.valueInterface()
	}
	return m
}

// unmarshaler reads raw bencode into an Unmarshaler.
func (d *decodeState) unmarshaler(v reflect.Value) {
	if v.Kind() == reflect.Ptr && v.IsNil() {
		v.Set(reflect.New(v.Type().Elem()))
	}
	u := v.Interface().(Unmarshaler)

	var tmpScan scanner
	tmpScan.reset()
	tmpScan.step = d.scan.step
	d.scan.step = stateEndValue

	start := d.off

	var op int
ReadRaw:
	for {
		if d.off > len(d.data) {
			d.error(errors.New("readed end of data"))
		}

		op = tmpScan.step(&tmpScan, int(d.data[d.off]))
		if op > 0 {
			d.off += op + 1
		} else {
			d.off++
			switch op {
			case scanEnd:
				break ReadRaw
			case scanError:
				d.error(tmpScan.err)
			}
		}
	}

	if err := u.UnmarshalBencode(d.data[start:d.off]); err != nil {
		d.error(err)
	}
}

// Unmarshal parses the bencode-encoded data and stores the
// result in the value pointed to by v.
//
// Unmarshal uses the inverse of the encodings that Marshal uses,
// allocating maps, silces, and pointers as necessary, with the
// following additional rules:
//
// To unmarshal bencode into a struct, Unmarshal matches incoming
// dictionaries to the keys used by Marshal (either the struct field
// name or its tag), preferring an exact match but also accepting a
// case-insensitive match.
//
// To unmarshal bencode into an interface value, Unmarshal unmarshals
// data into the concrete value contained in the interface value. If
// the interface value is nil, that is, has no concrete value stored in it,
// Unmarshal stores one of these in the interface value:
//
// int64, for integers
// []byte, for a byte string
// []interface{} for a list
// map[string]interface{} for a dictionary
//
func Unmarshal(data []byte, v interface{}) error {
	return NewDecoder(bytes.NewReader(data)).Decode(v)
}
