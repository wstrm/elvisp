package bencode

import (
	"bytes"
	"encoding"
	"fmt"
	"io"
	"reflect"
	"runtime"
	"sort"
)

// Encoder writes bencode data to an output stream..
type Encoder struct {
	w   io.Writer
	e   encodeState
	err error
}

// NewEncoder returns a new Encoder that bencodes to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

// Encode writes bencode data to the wrapped stream.
//
// See the documentation for Marshal for details about the
// conversion of Go values to bencode.
func (enc *Encoder) Encode(v interface{}) error {
	if enc.err != nil {
		return enc.err
	}
	enc.e.Reset()
	err := enc.e.marshal(v)
	if err != nil {
		return err
	}

	if _, err = enc.w.Write(enc.e.Bytes()); err != nil {
		enc.err = err
	}
	return err
}

// Marshal returns a bencoded form of x.
//
// Marshal traverses the value v recursively using the the following
// type-dependent encodings:
//
// Integer types encode as bencode integers.
//
// String and []byte values encode as bencode strings.
//
// Struct values encode as bencode dictionaries. Each exported struct
// field becomes a member of the object unless
//   - the field's tag is "-", or
//   - the field is empty and its tag specifies the "omitempty" option.
// The empty values are false, 0, any nil pointer or interface value,
// and any array, slice, string, or map with zero length.
// The values default key string is the struct field name but can be
// specified in the struct field's tag value. The "bencode" key in the
// struct field's tag value is the key name, followed by an optional
// comma and options. Examples:
//
//   // Field is ignored by this package.
//   Field int `bencode:"-"`
//
//   // Field appears in bencode dictionaries as "6:myName".
//   Field int `bencode:"myName"`
//
//   // Field appears in bencode dictionaries as "6:myName" and
//   // will be omitted if it's value is empty,
//   // as defined above.
//   Field int `bencode"myName,omitEmpty"`
//
// // Field appears in bencode dictionaries as "5:Field" (the default),
// // but the field is skipped if empty.
// // Note the leading comma.
// Field int `bencode:",omitempty"`
// The key name will be used if it's a non-empty string consisting of
// only Unicode letters, digits, dollar signs, percent signs, hyphens,
// underscores and slashes.
//
// Anonymous struct fields are usually marshaled as if their inner exported fields
// were fields in the outer struct, subject to the usual Go visibility rules amended
// as described in the next paragraph.
// An anonymous struct field with a name given in its bencode tag is treated as
// having that name, rather than being anonymous.
//
// The Go visibility rules for struct fields are amended for bencode when
// deciding which field to marshal or unmarshal. If there are
// multiple fields at the same level, and that level is the least
// nested (and would therefore be the nesting level selected by the
// usual Go rules), the following extra rules apply:
//
// 1) Of those fields, if any are bencode-tagged, only tagged fields are considered,
// even if there are multiple untagged fields that would otherwise conflict.
// 2) If there is exactly one field (tagged or not according to the first rule), that is selected.
// 3) Otherwise there are multiple fields, and all are ignored; no error occurs.
//
// Handling of anonymous struct fields is new in Go 1.1.
// Prior to Go 1.1, anonymous struct fields were ignored. To force ignoring of
// an anonymous struct field in both current and earlier versions, give the field
// a bencode tag of "-".
//
// Map values encode as bencode objects.
// The map's key type must be string; the object keys are used directly
// as map keys.
//
// Pointer values encode as the value pointed to.
// A nil pointer encodes as the null bencode object.
//
// Interface values encode as the value contained in the interface.
// A nil interface value encodes as the null bencode object.
//
// Channel, complex, and function values cannot be encoded in bencode.
// Attempting to encode such a value causes Marshal to return
// an UnsupportedTypeError.
//
// bencode cannot represent cyclic data structures and Marshal does not
// handle them.  Passing cyclic structures to Marshal will result in
// an infinite recursion.
func Marshal(v interface{}) ([]byte, error) {
	e := &encodeState{}
	err := e.marshal(v)
	if err != nil {
		return nil, err
	}
	return e.Bytes(), err
}

// An UnsupportedTypeError is returned by Marshal when attempting
// to encode an unsupported value type.
type UnsupportedTypeError struct {
	Type reflect.Type
}

func (e *UnsupportedTypeError) Error() string {
	return "bencode: unsupported type: " + e.Type.String()
}

type UnsupportedValueError struct {
	Value reflect.Value
	Str   string
}

func (e *UnsupportedValueError) Error() string {
	return "bencode: unsupported value: " + e.Str
}

// An encodeState encodes bencode into a bytes.Buffer.
type encodeState struct {
	bytes.Buffer // accumulated output
	scratch      [64]byte
}

func (e *encodeState) marshal(v interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
		}
	}()
	e.reflectValue(reflect.ValueOf(v))
	return nil
}

func (e *encodeState) error(err error) {
	panic(err)
}

var byteSliceType = reflect.TypeOf([]byte(nil))

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

var (
	marshalerType     = reflect.TypeOf(new(Marshaler)).Elem()
	textMarshalerType = reflect.TypeOf(new(encoding.TextMarshaler)).Elem()
)

// reflectValue writes the value in v to the output.
func (e *encodeState) reflectValue(v reflect.Value) {
	if !v.IsValid() {
		e.Write([]byte{'0', ':'})
		return
	}

	if v.Type().Implements(marshalerType) {

		m := v.Interface().(Marshaler)
		b, err := m.MarshalBencode()
		if err == nil {
			_, err = e.Write(b)
		}
		if err != nil {
			e.error(err)
		}
		return
	}

	if v.Type().Implements(textMarshalerType) {
		m := v.Interface().(encoding.TextMarshaler)
		b, err := m.MarshalText()
		if err != nil {
			e.error(err)
		}
		fmt.Fprintf(e, "%d:", len(b))
		_, err = e.Write(b)
		if err != nil {
			e.error(err)
		}
		return
	}

	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fmt.Fprintf(e, "i%de", v.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		fmt.Fprintf(e, "i%de", v.Uint())

	case reflect.String:
		s := v.String()
		fmt.Fprintf(e, "%d:%s", len(s), s)

	case reflect.Struct:
		e.WriteByte('d')
		for _, f := range cachedTypeFields(v.Type()) {
			fv := fieldByIndex(v, f.index)
			if !fv.IsValid() || f.omitEmpty && isEmptyValue(fv) {
				continue
			}
			fmt.Fprintf(e, "%d:%s", len(f.name), f.name)
			e.reflectValue(fv)
		}
		e.WriteByte('e')

	case reflect.Map:
		if v.Type().Key().Kind() != reflect.String {
			e.error(&UnsupportedTypeError{v.Type()})
		}
		e.WriteByte('d')
		if !v.IsNil() {
			var keys stringValues = v.MapKeys()
			sort.Sort(keys)
			for _, k := range keys {
				fmt.Fprintf(e, "%d:%s", k.Len(), k)
				e.reflectValue(v.MapIndex(k))
			}
		}
		e.WriteByte('e')

	case reflect.Slice, reflect.Array:
		if t := v.Type(); t == byteSliceType || t.Elem().Kind() == reflect.Uint8 {
			fmt.Fprintf(e, "%d:", v.Len())
			_, err := e.Write(v.Bytes())
			if err != nil {
				e.error(err)
			}
			return
		}
		e.WriteByte('l')
		n := v.Len()
		for i := 0; i < n; i++ {
			e.reflectValue(v.Index(i))
		}
		e.WriteByte('e')

	case reflect.Interface, reflect.Ptr:
		e.reflectValue(v.Elem())

	default:
		e.error(&UnsupportedTypeError{v.Type()})
	}
}
