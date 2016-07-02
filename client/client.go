package main

import (
	"bufio"
	"flag"
	"log"
	"net"
	"os"
	"strings"
)

var flags struct {
	leaseTask, removeTask bool
	serverAddr            string
}

func connect(addr string) (conn net.Conn, err error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp6", addr)
	if err != nil {
		return
	}

	return net.DialTCP("tcp", nil, tcpAddr)
}

func sendCmd(conn net.Conn, cmd string) (resp string, err error) {
	_, err = conn.Write([]byte(cmd + "\n"))
	if err != nil {
		return
	}

	r, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return
	}

	resp = strings.TrimSpace(string(r))

	return
}

func init() {
	flag.BoolVar(&flags.leaseTask, "l", flags.leaseTask, "Request lease.")
	flag.BoolVar(&flags.removeTask, "r", flags.removeTask, "Remove client.")
	flag.StringVar(&flags.serverAddr, "a", flags.serverAddr, "Address for server.")

	flag.Parse()

	if flags.serverAddr == "" {
		log.Fatal("No server address defined")
	}

	if !flags.leaseTask && !flags.removeTask {
		log.Fatal("No task defined")
	}
}

func main() {
	var err error
	var resp string

	conn, err := connect(flags.serverAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	if flags.leaseTask {
		resp, err = sendCmd(conn, "lease")
	}

	if flags.removeTask {
		resp, err = sendCmd(conn, "remove")
	}

	if err != nil {
		log.Fatal(err)
	}

	log.Println(resp)
	os.Exit(0)
}
