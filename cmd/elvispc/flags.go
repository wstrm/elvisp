package main

import (
	"flag"
	"log"
)

type flags struct {
	leaseTask, removeTask bool
	serverAddr            string
}

var context = flags{
	leaseTask:  false,
	removeTask: false,
	serverAddr: "",
}

func init() {
	flag.BoolVar(&context.leaseTask, "l", context.leaseTask, "Request lease.")
	flag.BoolVar(&context.removeTask, "r", context.removeTask, "Remove client.")
	flag.StringVar(&context.serverAddr, "a", context.serverAddr, "Address for server.")

	flag.Parse()

	if context.serverAddr == "" {
		log.Fatal("No server address defined")
	}

	if !context.leaseTask && !context.removeTask {
		log.Fatal("No task defined")
	}
}
