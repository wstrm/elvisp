package main

import (
	"log"

	"github.com/willeponken/elvisp/server"
)

func main() {
	port := ":4132"
	db := "/tmp/elvisp-db"
	password := "test123"
	cjdnsPassword := "test321"
	cjdnsIP := "127.0.0.1"
	cjdnsPort := 11234

	log.Printf("Listening to: %s and using database at: %s", port, db)
	log.Fatal(server.Listen(port, db, password, cjdnsIP, cjdnsPort, cjdnsPassword))
}
