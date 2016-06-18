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
	"github.com/willeponken/elvisp/lease"
	"github.com/willeponken/elvisp/tasks"
)

// Server holds a database and a connection to cjdns admin.
type Server struct {
	db                       *database.Database
	admin                    *cjdns.Conn
	ipv4Enabled, ipv6Enabled bool
	ipv4Lease, ipv6Lease     lease.Lease
}

// Settings holds settings needed to setup the server.
type Settings struct {
	Listen             string
	DB                 string
	Password           string
	CjdnsIP            string
	CjdnsPort          int
	CjdnsPassword      string
	IPv4CIDR, IPv6CIDR string
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

// parseCjdnsIPv6 takes a string and converts it into a IPv6 and checks if it's part of the cjdns address space.
func parseCjdnsIPv6(addr string) (ip net.IP, err error) {
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

	ipv6, err := parseCjdnsIPv6(address)
	if err != nil {
		log.Printf("Unable to resolve IPv6: %s", err)
		return tasks.Invalid{t}
	}

	context := tasks.Context{
		Argv:      argv,
		DB:        s.db,
		Admin:     s.admin,
		ClientIP:  ipv6,
		Auth:      auth,
		IPv4:      s.ipv4Enabled,
		IPv6:      s.ipv6Enabled,
		IPv4Lease: s.ipv4Lease,
		IPv6Lease: s.ipv6Lease,
	}

	t, err = tasks.Init(context)
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
func Listen(settings Settings) (err error) {
	var s Server

	if settings.IPv4CIDR != "" {
		var ipv4Lease lease.Lease

		ipv4Lease.Start, ipv4Lease.Network, err = net.ParseCIDR(settings.IPv4CIDR)
		if err != nil {
			return
		}

		s.ipv4Enabled = true
		s.ipv4Lease = ipv4Lease
	}

	if settings.IPv6CIDR != "" {
		var ipv6Lease lease.Lease

		ipv6Lease.Start, ipv6Lease.Network, err = net.ParseCIDR(settings.IPv6CIDR)
		if err != nil {
			return
		}

		s.ipv6Enabled = true
		s.ipv6Lease = ipv6Lease
	}

	// First, we need to make sure we are able to communicate with the database.
	db, err := database.Open(settings.DB)
	if err != nil {
		log.Printf("Unable to open database: %s", err)

		return
	}
	s.db = &db

	if settings.Password != "" {
		s.initAdmin(settings.Password)
	}

	// Listen only to IPv6 network. Administrators can connect locally using [::1].
	ln, err := net.Listen("tcp6", settings.Listen)
	if err != nil {
		log.Printf("Unable to listen to port: %s, due to error: %s", settings.Listen, err)

		return
	}

	// Connect to the cjdns admin interface.
	s.admin, err = cjdns.Connect(settings.CjdnsIP, settings.CjdnsPort, settings.CjdnsPassword)
	if err != nil {
		log.Printf("Unable to connect to cjdns admin on: %s:%d, due to error: %s", settings.CjdnsIP, settings.CjdnsPort, err)

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
