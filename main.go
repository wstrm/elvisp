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
	cidrIPv4      string
	cidrIPv6      string
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
	flag.StringVar(&context.cidrIPv4, "cidr-v4", context.cidrIPv4, "IPv4 CIDR to use for IP Leasing.")
	flag.StringVar(&context.cidrIPv6, "cidr-v6", context.cidrIPv6, "IPv6 CIDR to use for IP leasing.")

	flag.BoolVar(&context.version, "v", context.version, "Print current version and exit.")

	flag.IntVar(&context.cjdnsPort, "cjdns-port", context.cjdnsPort, "Port for cjdns admin.")

	flag.Parse()

	// If version flag is true, print version and exit.
	if context.version {
		fmt.Printf("%s (%s)\n", Version, BuildTime)
		os.Exit(0)
	}

	// Log line file:linenumber.
	log.SetFlags(log.LstdFlags | log.Llongfile)
	// Prefix log output with "[elvisp (<version>)]".
	log.SetPrefix("[\033[32melvisp\033[0m (" + Version + ")] ")

	if context.cidrIPv4 == "" && context.cidrIPv6 == "" {
		log.Fatalln("Atleast one CIDR has to be defined")
	}

	return
}

func main() {

	settings := server.Settings{
		Listen:        context.listen,
		DB:            context.db,
		Password:      context.password,
		CjdnsIP:       context.cjdnsIP,
		CjdnsPort:     context.cjdnsPort,
		CjdnsPassword: context.cjdnsPassword,
		IPv4CIDR:      context.cidrIPv4,
		IPv6CIDR:      context.cidrIPv6,
	}

	log.Printf("Listening to: %s and using database at: %s", context.listen, context.db)
	log.Fatal(server.Listen(settings))
}
