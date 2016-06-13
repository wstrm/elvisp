package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/willeponken/elvisp/server"
)

type flags struct {
	version       bool
	listen        string
	db            string
	password      string
	cjdnsIP       string
	cjdnsPort     int
	cjdnsPassword string
}

var context = flags{
	version:       false,
	listen:        ":4132",
	db:            "/tmp/elvisp-db",
	password:      "test123",
	cjdnsIP:       "127.0.0.1",
	cjdnsPort:     11234,
	cjdnsPassword: "ycdzz73bn17k22c017xtdxgmq7kn7xq",
}

var (
	// Version set with ldflags.
	Version string
	// BuildTime set by build script using ldflags.
	BuildTime string
)

func init() {

	flag.StringVar(&context.listen, "listen", context.listen, "Listen address for TCP.")
	flag.StringVar(&context.db, "db", context.db, "Directory to use for the database.")
	flag.StringVar(&context.password, "password", context.password, "Password for administrating elvisp.")
	flag.StringVar(&context.cjdnsIP, "cjdns-ip", context.cjdnsIP, "IP address for cjdns admin.")
	flag.StringVar(&context.cjdnsPassword, "cjdns-password", context.cjdnsPassword, "Password for cjdns admin.")

	flag.BoolVar(&context.version, "v", context.version, "Print current version and exit.")

	flag.IntVar(&context.cjdnsPort, "cjdns-port", context.cjdnsPort, "Port for cjdns admin.")

	flag.Parse()

	// If version flag is true, print version and exit.
	if context.version {
		fmt.Printf("%s (%s)\n", Version, BuildTime)
		os.Exit(0)
	}

	// Log line file:linenumber.
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// Prefix log output with "[elvisp (<version>)]".
	log.SetPrefix("[\033[32melvisp\033[0m (" + Version + ")] ")

	return
}

func main() {
	log.Printf("Listening to: %s and using database at: %s", context.listen, context.db)
	log.Fatal(server.Listen(context.listen, context.db, context.password, context.cjdnsIP, context.cjdnsPort, context.cjdnsPassword))
}
