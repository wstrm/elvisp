package main

import (
	"log"

	"github.com/willeponken/elvisp/server"
)

func init() {
	// Log line file:linenumber.
	log.SetFlags(log.LstdFlags | log.Llongfile)
	// Prefix log output with "[elvisp]".
	log.SetPrefix("[\033[32melvisp\033[0m] ")
}

func main() {
	if len(context.cidrList) < 1 {
		log.Fatalln("Atleast one CIDR has to be defined")
	}

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
