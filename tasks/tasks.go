package tasks

import (
	"net"

	"github.com/fc00/go-cjdns/key"
	"github.com/willeponken/elvisp/cjdns"
	"github.com/willeponken/elvisp/database"
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

	k, err := admin.LookupPubKey(ip.String())
	if err != nil {
		return
	}

	task.pubkey, err = key.DecodePublic(k)

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
	if (len(t.argv) != 1) || (len(t.argv) != 2) {
		return t.errorString(errorInvalidLength)
	}

	return t.successString("<<LEASED ADDRESS HERE>>")
}

// Run removes a user
func (t Remove) Run() string {
	if len(t.argv) != 2 {
		return t.errorString(errorInvalidLength)
	}

	return t.successString("Removed user: " + t.argv[0])
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
