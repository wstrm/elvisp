package main

import (
	"log"

	"github.com/willeponken/elvisp/server"
)

var (
	// Version set with ldflags.
	Version string
	// BuildTime set by build script using ldflags.
	BuildTime string
)

// Default values for flags
var context = flags{
	version:   false,
	listen:    ":4132",
	db:        "/tmp/elvisp-db",
	cjdnsIP:   "127.0.0.1",
	cjdnsPort: 11234,
}

func init() {
	// Log line file:linenumber.
	log.SetFlags(log.LstdFlags | log.Llongfile)
	// Prefix log output with "[elvisp (<version>)]".
	log.SetPrefix("[\033[32melvisp\033[0m (" + Version + ")] ")
}

func main() {

	settings := server.Settings{
		Listen:        context.listen,
		DB:            context.db,
		Password:      context.password,
		CjdnsIP:       context.cjdnsIP,
		CjdnsPort:     context.cjdnsPort,
		CjdnsPassword: context.cjdnsPassword,
		CIDRs:         context.cidrList.List(),
	}

	log.Printf("Listening to: %s and using database at: %s", context.listen, context.db)
	log.Fatal(server.Listen(settings))
}
