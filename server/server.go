package server

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net"
	"strings"

	"github.com/willeponken/elvisp/database"
	"github.com/willeponken/elvisp/tasks"
)

// Server holds a database.
type Server struct {
	db database.Database
}

// taskFactory creates a new task based on an string which defines the type.
func (s *Server) taskFactory(input string) tasks.TaskInterface {
	t := tasks.Task{}

	array := strings.Split(input, " ")
	if len(array) >= 2 {
		cmd := strings.ToLower(array[0])
		argv := array[1:]

		t.SetArgs(argv)
		t.SetDB(&s.db)

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
func (s *Server) taskRunner(t tasks.TaskInterface, out chan string) {
	out <- t.Run() + "\n"
}

// requestHandler reads from a TCP connection/session and writes it to a channel.
func (s *Server) requestHandler(conn net.Conn, out chan string) error {
	defer close(out)

	for {
		line, err := bufio.NewReader(conn).ReadBytes('\n')
		if err != nil {
			return err
		}

		t := s.taskFactory(strings.TrimRight(string(line), "\n"))
		go s.taskRunner(t, out)
	}
}

// sendHandler copies all communication from a channel to a TCP connection/session. Empty messages and errors terminates the loop.
func (s *Server) sendHandler(conn net.Conn, in <-chan string) {
	defer conn.Close()

	for {
		message := <-in
		rb, err := io.Copy(conn, bytes.NewBufferString(message))

		if rb == 0 || err != nil {
			return
		}
	}
}

// Listen starts listening on a defined port using TCP and connects to a BoltDB database. It will then initialize two handlers, request and send handler, as goroutines.
func Listen(port, db string) (err error) {
	var s Server

	// First, we need to make sure we are able to communicate with the database
	s.db, err = database.Open(db)
	if err != nil {
		return
	}

	ln, err := net.Listen("tcp", port)
	if err != nil {
		return
	}

	var conn net.Conn
	for {
		conn, err = ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		channel := make(chan string)

		go s.requestHandler(conn, channel)
		go s.sendHandler(conn, channel)
	}
}
