package main

import (
	"log"

	"github.com/willeponken/elvisp/server"
)

func main() {
	port := ":4132"
	db := "/tmp/elvisp-db"

	log.Printf("Listening to: %s and using database at: %s", port, db)
	log.Fatal(server.Listen(port, db))
}
