package main

import "flag"

type cidrList []string

type flags struct {
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
	listen:    ":4132",
	db:        "/tmp/elvispd-db",
	cjdnsIP:   "127.0.0.1",
	cjdnsPort: 11234,
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
	if len(str) > 1 {
		str = str[:(len(str) - 1)] // Trim superflous whitespace at end
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

	flag.IntVar(&context.cjdnsPort, "cjdns-port", context.cjdnsPort, "Port for cjdns admin.")

	flag.Parse()

	return
}
