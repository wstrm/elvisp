package bencode

import (
	"bytes"
	"encoding"
	"fmt"
	"io"
	"net"
	"reflect"
	"testing"
)

type test struct {
	in  string
	ptr interface{}
	out interface{}
	err error
}

type Ambig struct {
	// Given "hello", the first match should win.
	First  int `bencode:"HELLO"`
	Second int `bencode:"Hello"`
}

var tests = []test{
	// basic types
	//{in: `i1e`, ptr: new(bool), out: true},
	{in: `i1e`, ptr: new(int), out: 1},
	{in: `i2e`, ptr: new(int64), out: int64(2)},
	{in: `i-5e`, ptr: new(int16), out: int16(-5)},
	{in: `i2e`, ptr: new(interface{}), out: int64(2)},
	{in: "i0e", ptr: new(interface{}), out: int64(0)},
	{in: "i0e", ptr: new(int), out: 0},

	{in: "0:", ptr: new(string), out: ""},
	{in: "1:a", ptr: new(string), out: "a"},
	{in: "2:a\"", ptr: new(string), out: "a\""},
	{in: "3:abc", ptr: new([]byte), out: []byte("abc")},
	{in: "11:0123456789a", ptr: new(interface{}), out: []byte("0123456789a")},
	{in: "le", ptr: new([]int64), out: []int64{}},
	{in: "li1ei2ee", ptr: new([]int), out: []int{1, 2}},
	{in: "l3:abc3:def0:e", ptr: new([]string), out: []string{"abc", "def", ""}},
	//{in: "li42e3:abce", ptr: new([]interface{}), out: []interface{}{42, []byte("abc")}},
	{in: "de", ptr: new(map[string]interface{}), out: make(map[string]interface{})},
	{in: "d3:cati1e3:dogi2ee", ptr: new(map[string]int), out: map[string]int{"cat": 1, "dog": 2}},
	{in: "9:127.0.0.1", ptr: new(net.IP), out: net.ParseIP("127.0.0.1")},
	{in: "d1:i3:1231:m9:Arith.Adde", ptr: new(request), out: request{"Arith.Add", nil, "123"}},
}

var afs = []byte("d18:availableFunctionsd18:AdminLog_subscribed4:filed8:requiredi0e4:type6:Stringe5:leveld8:requiredi0e4:type6:Stringe4:lined8:requiredi0e4:type3:Intee20:AdminLog_unsubscribed8:streamIdd8:requiredi1e4:type6:Stringee18:Admin_asyncEnabledde24:Admin_availableFunctionsd4:paged8:requiredi0e4:type3:Intee34:InterfaceController_disconnectPeerd6:pubkeyd8:requiredi1e4:type6:Stringee29:InterfaceController_peerStatsd4:paged8:requiredi0e4:type3:Intee17:SwitchPinger_pingd4:datad8:requiredi0e4:type6:Stringe4:pathd8:requiredi1e4:type6:Stringe7:timeoutd8:requiredi0e4:type3:Intee16:UDPInterface_newd11:bindAddressd8:requiredi0e4:type6:Stringeee4:morei1e4:txid8:c37b0faae")

type outer struct {
	More bool   `bencode:"more"`
	Txid string `bencode:"txid"`
}

func TestMarshal(t *testing.T) {
	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	for i, tt := range tests {
		buf.Reset()
		var scan scanner
		in := []byte(tt.in)
		if err := checkValid(in, &scan); err != nil {
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("#%d: %s checkValid: %#v", i, tt.in, err)
				continue
			}
		}
		if err := enc.Encode(tt.out); err != nil {
			t.Errorf("#%d: %q Error: %s", i, tt.in, err)
			continue
		}

		out := buf.String()
		if out != tt.in {
			t.Errorf("#%d: Want %q, got %q", i, tt.in, out)
		}
	}
}

func TestUnmarshal(t *testing.T) {
	for i, tt := range tests {
		var scan scanner
		in := []byte(tt.in)
		if err := checkValid(in, &scan); err != nil {
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("#%d: %s checkValid: %#v", i, tt.in, err)
				continue
			}
		}
		if tt.ptr == nil {
			continue
		}
		v := reflect.New(reflect.TypeOf(tt.ptr).Elem())
		dec := NewDecoder(bytes.NewBuffer(in))
		if err := dec.Decode(v.Interface()); !reflect.DeepEqual(err, tt.err) {
			t.Errorf("#%d: %q %v want %v", i, tt.in, err, tt.err)
			continue
		}
		if !reflect.DeepEqual(v.Elem().Interface(), tt.out) {
			t.Errorf("#%d: %s mismatch\nhave: %#+v\nwant: %#+v", i, tt.in, v.Elem().Interface(), tt.out)
		}
	}
}

func TestSkip(t *testing.T) {
	o := new(outer)
	err := Unmarshal(afs, o)
	if err != nil {
		t.Fatal("error unmarshaling nested struct,", err)
	}
	if o.Txid != "c37b0faa" {
		t.Errorf("got txid %q", o.Txid)
	}
}

type request struct {
	Method string      `bencode:"m"`
	Params interface{} `bencode:"p,omitempty"`
	Id     string      `bencode:"i"`
}

func TestRequest(t *testing.T) {
	var buf bytes.Buffer
	var req request
	dec := NewDecoder(&buf)
	for i := 0; i < 8; i++ {
		fmt.Fprintf(&buf, "d1:ii%de1:m9:Arith.Add1:pd1:Ai%de1:Bi%deee",
			i, i, i+1)
		if err := dec.Decode(&req); err != nil {
			t.Fatal(err)
		}
	}
}

type nestA struct {
	A, B int
	C    *nestB
}

type nestB struct {
	D, E int
	F    *nestA `bencode:",omitempty"`
}

func TestNesting(t *testing.T) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)
	dec := NewDecoder(&buf)

	var j int
	top := new(nestA)
	nest := top
	for i := 0; i < 32; i++ {
		nest.A = j
		j++
		nest.B = j
		j++

		nest.C = new(nestB)
		nest.C.D = j
		j++
		nest.C.E = j
		j++

		if err := enc.Encode(top); err != nil {
			t.Error("nesting:", err)
			return
		}

		out := new(nestA)
		if err := dec.Decode(out); err != nil {
			t.Error("nesting:", err)
			return
		}

		nest.C.F = new(nestA)
		nest = nest.C.F
	}
}

type testInterface struct {
	s string
}

func (i *testInterface) MarshalText() (text []byte, err error) {
	return []byte(i.s), nil
}

func (i *testInterface) UnmarshalText(b []byte) error {
	i.s = string(b)
	return nil
}

func TestTextInterface(t *testing.T) {
	var controlA encoding.TextMarshaler
	controlA = &testInterface{"The Cultural Myth of Female Hair in the Victorian Imagination"}

	bA, err := Marshal(controlA)
	if err != nil {
		t.Error("failed to marshal", err)
		return
	}

	controlB := &testInterface{"The Cultural Myth of Female Hair in the Victorian Imagination"}

	bB, err := Marshal(controlB)
	if err != nil {
		t.Error("failed to marshal", err)
		return
	}

	var testA encoding.TextUnmarshaler
	testA = new(testInterface)

	err = Unmarshal(bB, &testA)
	if err != nil {
		t.Error("failed to unmarshal what was marshaled,", err)
		return
	}

	testB := new(testInterface)

	err = Unmarshal(bA, &testB)
	if err != nil {
		t.Error("failed to unmarshal what was marshaled,", err)
		return
	}

	if testB.s != controlB.s {
		t.Errorf("wanted %q, got %q", controlB.s, testB.s)
	}
}

type cooked struct {
	A, B, C int
}

type raw struct {
	A, B int
	C    *RawMessage
}

func TestRawMessage(t *testing.T) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)
	dec := NewDecoder(&buf)

	c := &cooked{1, 2, 3}
	err := enc.Encode(c)
	if err != nil {
		t.Error("encoding cooked:", err)
	}
	r := new(raw)
	err = dec.Decode(r)
	if err != nil {
		t.Error("decoding to raw:", err)
	}

	if len(*r.C) == 0 {
		t.Fatal("nested RawMessage had zero length")
	}

	var x int
	err = Unmarshal(*r.C, &x)
	if err != nil {
		t.Error("decoding RawMessage:", err)
	}
	if x != 3 {
		t.Errorf("RawMessage: want %d, got %d", 3, x)
	}

	err = enc.Encode(r)
	if err != nil {
		t.Error("encoding RawMessage:", err)
	}
	err = dec.Decode(c)
	if err != nil {
		t.Error("decoding encoding RawMessage:", err)
	}
	if c.C != 3 {
		t.Error("mismatch")
	}
}

func TestPipe(t *testing.T) {
	r, w := io.Pipe()
	dec := NewDecoder(r)
	enc := NewEncoder(w)

	for i, tt := range tests {
		if tt.ptr == nil {
			continue
		}
		v := reflect.New(reflect.TypeOf(tt.ptr).Elem())

		go enc.Encode(tt.out)
		if err := dec.Decode(v.Interface()); err != nil {
			t.Errorf("#%d: %q %v want %v", i, tt.in, err, tt.err)
		}

		if !reflect.DeepEqual(v.Elem().Interface(), tt.out) {
			t.Errorf("#%d: %s mismatch\nhave: %#+v\nwant: %#+v", i, tt.in, v.Elem().Interface(), tt.out)
		}
	}
}

type serverResponse struct {
	Error  interface{} `bencode:"e"` //,omitempty"`
	Id     *RawMessage `bencode:"i"`
	Result interface{} `bencode:"r"`
}

func TestEncodeNil(t *testing.T) {
	var resp serverResponse
	resp.Error = nil
	resp.Result = nil
	_, err := Marshal(resp)
	if err != nil {
		t.Error("EncodeNil:", err)
	}
}

type benchmarkStruct struct {
	Q      string      `bencode:"q"`
	AQ     string      `bencode:"aq,omitempty"`
	Cookie string      `bencode:"cookie,omitempty"`
	Hash   string      `bencode:"hash,omitempty"`
	Args   interface{} `bencode:"args,omitempty"`
	Txid   string      `bencode:"txid"`
}

var benchmarkTest = []byte("d1:q4:auth2:aq4:ping6:cookie10:13536270564:hash64:d1e4881e30895e2ee3e14c9bbce4537288a72a242dbd1786e8f1cc512e4cb4674:txid8:37199054e")

func BenchmarkUnmarshal(b *testing.B) {
	x := new(benchmarkStruct)
	var err error
	for i := 0; i < b.N; i++ {
		if err = Unmarshal(benchmarkTest, x); err != nil {
			b.Fatal(err.Error())
		}
	}
}
