package server

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"
	"net"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/willeponken/elvisp/cjdns"
	"github.com/willeponken/elvisp/database"
	"github.com/willeponken/elvisp/tasks"
)

// Server holds a database and a connection to cjdns admin.
type Server struct {
	db    *database.Database
	admin *cjdns.Conn
}

// authAdmin checks the password with the saved hash in the database.
func (s *Server) authAdmin(password string) error {
	hash, err := s.db.AdminHash()
	if err != nil {
		return err
	}

	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// initAdmin sets the hashed admin password in the database.
func (s *Server) initAdmin(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.db.SetAdmin(string(hash))
}

// resolveIPv6 takes a string and converts it into IP address.
func resolveIPv6(addr string) (ip net.IP, err error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp6", addr)
	if err != nil {
		return
	}

	ip = tcpAddr.IP

	str := ip.String()

	if ip.To4() != nil { // If able to parse to IPv4 the address is invalid.
		err = errors.New(ip.String() + "is not IPv6")
		return
	}
	if ip.To16()[0] != 0xFC { // If the first bit is not 0xFC then it's not in the cjdns address space.
		err = errors.New(str + " is not in the cjdns address space (0xFC)")
	}

	return
}

// taskFactory creates a new task based on an string which defines the type.
func (s *Server) taskFactory(conn net.Conn, input string) tasks.TaskInterface {
	var t tasks.Task

	array := strings.Split(input, " ")
	if len(array) < 1 {
		return tasks.Invalid{t}
	}

	var password, address string

	cmd := strings.ToLower(array[0])
	argv := array[1:]
	auth := false

	// If longer than 3, the second element should be a password for the administrator, and the third the address.
	if len(array) == 3 {
		password = array[1]
		address = array[2]

		if err := s.authAdmin(password); err != nil {
			log.Printf("Failed to authenticate administrator: %s", err)

			return tasks.Invalid{t}
		}

		auth = true
	} else {
		// If not admin, use the remote address that is currently connecting.
		address = conn.RemoteAddr().String()
	}

	ipv6, err := resolveIPv6(address)
	if err != nil {
		log.Printf("Unable to resolve IPv6: %s", err)
		return tasks.Invalid{t}
	}

	t, err = tasks.Init(argv, s.db, s.admin, ipv6, auth)
	if err != nil {
		log.Printf("Unable to initialize task: %s", err)
		return tasks.Invalid{t}
	}

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

		t := s.taskFactory(conn, strings.TrimRight(string(line), "\n"))
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

// Listen starts listening on a defined port using TCP6, connects to a BoltDB database and sets a admin password if defined. It will then initialize two handlers, request and send handler, as goroutines.
func Listen(port, db, password, cjdnsIP string, cjdnsPort int, cjdnsPassword string) (err error) {
	var s Server

	// First, we need to make sure we are able to communicate with the database.
	d, err := database.Open(db)
	if err != nil {
		log.Printf("Unable to open database: %s", err)

		return
	}
	s.db = &d

	if password != "" {
		s.initAdmin(password)
	}

	// Listen only to IPv6 network. Administrators can connect locally using [::1].
	ln, err := net.Listen("tcp6", port)
	if err != nil {
		log.Printf("Unable to listen to port: %s, due to error: %s", port, err)

		return
	}

	// Connect to the cjdns admin interface.
	s.admin, err = cjdns.Connect(cjdnsIP, cjdnsPort, cjdnsPassword)
	if err != nil {
		log.Printf("Unable to connect to cjdns admin on: %s:%d, due to error: %s", cjdnsIP, cjdnsPort, err)

		return
	}

	var conn net.Conn
	for {
		conn, err = ln.Accept()
		if err != nil {
			log.Printf("TCP connection returned error: %s", err)

			continue
		}

		channel := make(chan string)

		go s.requestHandler(conn, channel)
		go s.sendHandler(conn, channel)
	}
}
