package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/willeponken/elvisp/server"
)

type cidrList []string

type flags struct {
	version       bool
	listen        string
	db            string
	password      string
	cidrList      cidrList
	cjdnsIP       string
	cjdnsPort     int
	cjdnsPassword string
}

// Default values for flags
var context = flags{
	version:   false,
	listen:    ":4132",
	db:        "/tmp/elvisp-db",
	cjdnsIP:   "127.0.0.1",
	cjdnsPort: 11234,
}

var (
	// Version set with ldflags.
	Version string
	// BuildTime set by build script using ldflags.
	BuildTime string
)

// List cidrList lists all the CIDR's as a slice of strings
func (c cidrList) List() (cidrs []string) {
	for _, cidr := range c {
		cidrs = append(cidrs, cidr)
	}
	return
}

// String cidrList stringifies the list of CIDR's
func (c *cidrList) String() (str string) {
	for _, cidr := range c.List() {
		str += cidr + " "
	}
	return str
}

// Set cidrList appends the list of CIDR's with a new string (hopefully a CIDR)
func (c *cidrList) Set(cidr string) error {
	*c = append(*c, cidr)
	return nil
}

func init() {

	flag.StringVar(&context.listen, "listen", context.listen, "Listen address for TCP.")
	flag.StringVar(&context.db, "db", context.db, "Directory to use for the database.")
	flag.StringVar(&context.password, "password", context.password, "Password for administrating elvisp.")
	flag.StringVar(&context.cjdnsIP, "cjdns-ip", context.cjdnsIP, "IP address for cjdns admin.")
	flag.StringVar(&context.cjdnsPassword, "cjdns-password", context.cjdnsPassword, "Password for cjdns admin.")

	flag.Var(&context.cidrList, "cidr", "CIDR to use for IP leasing, use flag repeatedly for multiple CIDR's.")

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

	if len(context.cidrList) < 1 {
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
		CIDRs:         context.cidrList.List(),
	}

	log.Printf("Listening to: %s and using database at: %s", context.listen, context.db)
	log.Fatal(server.Listen(settings))
}
