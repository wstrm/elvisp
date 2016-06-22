package main

import (
	"flag"
	"fmt"
	"log"
	"os"
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
	return
}

// Set cidrList appends the list of CIDR's with a new string (hopefully a CIDR)
func (c *cidrList) Set(cidr string) error {
	*c = append(*c, cidr)
	return nil
}

func init() {

	flag.StringVar(&context.listen, "listen", context.listen, "Listen address for TCP.")
	flag.StringVar(&context.db, "db", context.db, "Directory to use for the database.")
	flag.StringVar(&context.password, "password", context.password, "Password for administrating Elvisp.")
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

	if len(context.cidrList) < 1 {
		log.Fatalln("Atleast one CIDR has to be defined")
	}

	return
}
