package key

import (
	"fmt"
)

func ExampleDecodePrivate() {
	key, _ := DecodePrivate("751d3db85b848deaf221e0ed2b6cc17f587b29057d74cdd4dc0bd18b7157288e")
	fmt.Println(key.Valid())
	// Output:
	// true
}

func ExamplePrivate_String() {
	key, _ := DecodePrivate("751d3db85b848deaf221e0ed2b6cc17f587b29057d74cdd4dc0bd18b7157288e")
	fmt.Printf("%s\n", key)
	// Output:
	// 751d3db85b848deaf221e0ed2b6cc17f587b29057d74cdd4dc0bd18b7157288e
}

func ExamplePrivate_Pubkey() {
	key, _ := DecodePrivate("751d3db85b848deaf221e0ed2b6cc17f587b29057d74cdd4dc0bd18b7157288e")
	fmt.Printf("%s\n", key.Pubkey())
	// Output:
	// r6jzx210usqbgnm3pdtm1z6btd14pvdtkn5j8qnpgqzknpggkuw0.k
}

func ExampleDecodePublic() {
	key, _ := DecodePublic("r6jzx210usqbgnm3pdtm1z6btd14pvdtkn5j8qnpgqzknpggkuw0.k")
	fmt.Println(key.Valid())
	// Output:
	// true
}

func ExamplePublic_IP() {
	key, _ := DecodePublic("r6jzx210usqbgnm3pdtm1z6btd14pvdtkn5j8qnpgqzknpggkuw0.k")
	fmt.Println(key.IP())
	// Output:
	// fc68:cb2c:60db:cb96:19ac:34a8:fd34:3fc
}
