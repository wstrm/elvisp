package server

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net"
	"strings"

	"github.com/willeponken/elvisp/tasks"
)

// taskFactory creates a new task based on an string which defines the type.
func taskFactory(input string) tasks.TaskInterface {
	t := tasks.Task{}

	array := strings.Split(input, " ")
	if len(array) >= 2 {
		cmd := strings.ToLower(array[0])
		argv := array[1:]

		t.SetArgs(argv)

		switch cmd {
		case "add":
			return tasks.Add{t}
		case "remove":
			return tasks.Remove{t}
		case "lease":
			return tasks.Lease{t}
		case "release":
			return tasks.Release{t}
		}
	}

	return tasks.Invalid{t}
}

// taskRunner runs a task and inputs its output into a channel.
func taskRunner(t tasks.TaskInterface, out chan string) {
	out <- t.Run() + "\n"
}

// requestHandler reads from a TCP connection/session and writes it to a channel.
func requestHandler(conn net.Conn, out chan string) error {
	defer close(out)

	for {
		line, err := bufio.NewReader(conn).ReadBytes('\n')
		if err != nil {
			return err
		}

		t := taskFactory(strings.TrimRight(string(line), "\n"))
		go taskRunner(t, out)
	}
}

// sendHandler copies all communication from a channel to a TCP connection/session. Empty messages and errors terminates the loop.
func sendHandler(conn net.Conn, in <-chan string) {
	defer conn.Close()

	for {
		message := <-in
		rb, err := io.Copy(conn, bytes.NewBufferString(message))

		if rb == 0 || err != nil {
			return
		}
	}
}

// Listen starts listening on a defined port using TCP. It will then initialize two handlers, request and send handler, as goroutines.
func Listen(port string) error {

	ln, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	var conn net.Conn
	for {
		conn, err = ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		channel := make(chan string)

		go requestHandler(conn, channel)
		go sendHandler(conn, channel)
	}
}
