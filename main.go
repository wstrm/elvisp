package main

import (
	"log"

	"github.com/willeponken/elvisp/server"
)

func main() {
	log.Fatal(server.Listen(":4132"))
}
