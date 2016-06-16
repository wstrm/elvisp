package tasks

import (
	"log"
	"net"

	"github.com/ehmry/go-cjdns/key"
	"github.com/willeponken/elvisp/cjdns"
	"github.com/willeponken/elvisp/database"
	"github.com/willeponken/elvisp/lease"
)

const (
	errorInvalidLength = "2 Invalid length of arguments"
	errorInvalidTask   = "1 Invalid task specified"
)

// TaskInterface defines the methods needed for a default task
type TaskInterface interface {
	Run() string
}

// Task needs the arguments to use, and a database to save the changes to
type Task struct {
	argv   []string
	db     *database.Database
	admin  *cjdns.Conn
	ip     net.IP
	pubkey *key.Public
	auth   bool
}

// Init returns a new task
func Init(argv []string, db *database.Database, admin *cjdns.Conn, ip net.IP, auth bool) (task Task, err error) {
	task.argv = argv
	task.ip = ip
	task.auth = auth

	var k string
	if auth {
		k = argv[0]
	} else {
		k, err = admin.LookupPubKey(ip.String())

		if err != nil {
			return task, err
		}
	}

	task.pubkey, err = key.DecodePublic(k)
	if err != nil {
		return task, err
	}

	task.db = db
	task.admin = admin

	return
}

// errorString returns a string prefixed with "error"
func (t Task) errorString(e string) string {
	return "error " + e
}

// successString returns a string prefix with "success"
func (t Task) successString(s string) string {
	return "success " + s
}

// Add should implement the add task
type Add struct{ Task }

// Remove should implement the remove task
type Remove struct{ Task }

// Lease should implement the lease task
type Lease struct{ Task }

// Release should implement the release task
type Release struct{ Task }

// Invalid should implement an invalid task which is returned if no task could be found for a command
type Invalid struct{ Task }

// Run adds a user using the public key and a token
func (t Add) Run() string {
	db := t.db
	admin := t.admin
	pubkey := t.pubkey

	id, err := db.AddUser(pubkey)
	if err != nil {
		return t.errorString(err.Error())
	}

	ip, err := lease.Generate("192.168.1.0/24", id)
	if err != nil {
		return t.errorString(err.Error())
	}

	if err := admin.AddUser(pubkey, ip); err != nil {
		// If we failed to add the user, delete it from the database.
		if e := db.DelUser(pubkey); e != nil {
			log.Println(e)
		}

		return t.errorString(err.Error())
	}

	return t.successString(ip.String())
}

// Run removes a user
func (t Remove) Run() string {
	db := t.db
	admin := t.admin
	pubkey := t.pubkey

	if err := db.DelUser(pubkey); err != nil {
		return t.errorString(err.Error())
	}

	if err := admin.DelUser(pubkey); err != nil {
		return t.errorString(err.Error())
	}

	return t.successString("Removed user: " + pubkey.String())
}

// Run leases a new address (if available)
func (t Lease) Run() string {
	if len(t.argv) != 2 {
		return t.errorString(errorInvalidLength)
	}

	return t.successString("<<LEASED ADDRESS HERE>>")
}

// Run releases a lease from a user
func (t Release) Run() string {
	if len(t.argv) != 2 {
		return t.errorString(errorInvalidLength)
	}

	return t.successString("Released lease for user: " + t.argv[0])
}

// Run returns a Invalid Task error
func (t Invalid) Run() string {
	return t.errorString(errorInvalidTask)
}
