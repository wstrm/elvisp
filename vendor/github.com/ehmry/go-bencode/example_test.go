package bencode

import "fmt"

func ExampleMarshal() {
	p := new(struct {
		A int    `bencode:"foo"`
		B string `bencode:"bar"`
	})

	p.A = 42
	p.B = "spam"
	b, _ := Marshal(p)
	fmt.Printf("%q", b)
	// Output:
	// "d3:bar4:spam3:fooi42ee"
}

func ExampleUnmarshal() {
	b := []byte("d3:bar4:spam3:fooi42ee")

	p := new(struct {
		A int    `bencode:"foo"`
		B string `bencode:"bar"`
	})
	Unmarshal(b, p)
	fmt.Printf("%+v", *p)
	// Output:
	// {A:42 B:spam}
}
