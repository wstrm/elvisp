package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"strings"
)

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

func main() {
	conn, err := connect("[fc38:f1bc:28ad:21be:2c9d:a543:a091:3087]:4132")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	resp, err := sendCmd(conn, "lease")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(resp)
	os.Exit(0)
}
