package tasks

import (
	"fmt"
	"log"
	"net"

	"github.com/willeponken/go-cjdns/key"
	"github.com/willeponken/elvisp/cjdns"
	"github.com/willeponken/elvisp/database"
	"github.com/willeponken/elvisp/lease"
)

// TaskInterface defines the methods needed for a default task
type TaskInterface interface {
	Run() (result string, err error)
}

// Task needs the arguments to use, and a database to save the changes to
type Task struct {
	argv     []string
	db       *database.Database
	admin    *cjdns.Conn
	clientIP net.IP
	pubkey   *key.Public
	cidrs    []lease.CIDR
}

// Init returns a new task
func Init(argv []string, db *database.Database, admin *cjdns.Conn, clientIP net.IP, cidrs []lease.CIDR) (task Task, err error) {
	task.argv = argv
	task.db = db
	task.admin = admin
	task.clientIP = clientIP
	task.cidrs = cidrs

	var k string
	k, err = task.admin.LookupPubKey(clientIP.String())

	if err != nil {
		return task, err
	}

	task.pubkey, err = key.DecodePublic(k)
	if err != nil {
		return task, err
	}

	return
}

// Add should implement the add task
type Add struct{ Task }

// Remove should implement the remove task
type Remove struct{ Task }

// Lease should implement the lease task
type Lease struct{ Task }

// Release should implement the release task
type Release struct{ Task }

// Invalid should implement the invalid task, i.e. take an error
type Invalid struct{ Error error }

// allowIPTunnel adds the defined IP to the cjdns IP tunnel, if it fails, it will delete the user from the database.
func (t Add) allowIPTunnel(c lease.CIDR, id uint64) (ip net.IP, err error) {

	ip, err = lease.Generate(c, id)
	if err != nil {
		return
	}

	if err = t.admin.AddUser(t.pubkey, ip); err != nil {
		if e := t.db.DelUser(t.pubkey); e != nil {
			log.Println(err)
		}

		return
	}

	return
}

// Run Add adds a user using the public key and a token.
func (t Add) Run() (result string, err error) {
	var id uint64
	var ip net.IP

	id, err = t.db.AddUser(t.pubkey)
	if err != nil {
		return
	}

	for _, cidr := range t.cidrs {
		ip, err = t.allowIPTunnel(cidr, id)
		if err != nil {
			return
		}

		result += ip.String() + " "
	}

	return
}

// Run Remove removes a user.
func (t Remove) Run() (result string, err error) {
	db := t.db
	admin := t.admin
	pubkey := t.pubkey

	if err = db.DelUser(pubkey); err != nil {
		return
	}

	if err = admin.DelUser(pubkey); err != nil {
		return
	}

	result = fmt.Sprintf("Removed user: %s", pubkey.String())
	return
}

// Run Lease returns the users leases.
func (t Lease) Run() (result string, err error) {
	var id uint64
	var ip net.IP

	id, err = t.db.GetID(t.pubkey)
	if err != nil {
		return
	}

	for _, cidr := range t.cidrs {
		ip, err = lease.Generate(cidr, id)
		if err != nil {
			return
		}

		result += ip.String() + " "
	}

	return
}

// Run Invalid returns an error and empty result.
func (t Invalid) Run() (result string, err error) {
	err = t.Error
	return
}
