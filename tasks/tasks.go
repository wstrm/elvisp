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
	argv                 []string
	db                   *database.Database
	admin                *cjdns.Conn
	clientIP             net.IP
	ipv4, ipv6           bool
	ipv4Lease, ipv6Lease lease.Lease
	pubkey               *key.Public
	auth                 bool
}

// Context holds current settings for the task that should be initialized with Init().
type Context struct {
	Argv                 []string
	DB                   *database.Database
	Admin                *cjdns.Conn
	ClientIP             net.IP
	Auth                 bool
	IPv4, IPv6           bool
	IPv4Lease, IPv6Lease lease.Lease
}

// Init returns a new task
func Init(context Context) (task Task, err error) {
	task.argv = context.Argv
	task.clientIP = context.ClientIP
	task.auth = context.Auth
	task.admin = context.Admin
	task.ipv4 = context.IPv4
	task.ipv6 = context.IPv6
	task.ipv4Lease = context.IPv4Lease
	task.ipv6Lease = context.IPv6Lease

	var k string
	if context.Auth {
		k = context.Argv[0]
	} else {
		k, err = task.admin.LookupPubKey(context.ClientIP.String())

		if err != nil {
			return task, err
		}
	}

	task.pubkey, err = key.DecodePublic(k)
	if err != nil {
		return task, err
	}

	task.db = context.DB
	task.admin = context.Admin

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

// allowIPTunnel adds the defined IP to the cjdns IP tunnel, if it fails, it will delete the user from the database.
func (t Add) allowIPTunnel(l lease.Lease, id uint64) (ip net.IP, err error) {

	ip, err = lease.Generate(l, id)
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

// Run adds a user using the public key and a token
func (t Add) Run() string {
	var id uint64
	var ipv4, ipv6 net.IP
	var err error

	id, err = t.db.AddUser(t.pubkey)
	if err != nil {
		return t.errorString(err.Error())
	}

	if t.ipv4 {
		if ipv4, err = t.allowIPTunnel(t.ipv4Lease, id); err != nil {
			return t.errorString(err.Error())
		}
	}

	if t.ipv6 {
		if ipv6, err = t.allowIPTunnel(t.ipv6Lease, id); err != nil {
			return t.errorString(err.Error())
		}
	}

	return t.successString(ipv4.String() + " " + ipv6.String())
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

// Run returns a Invalid Task error
func (t Invalid) Run() string {
	return t.errorString(errorInvalidTask)
}
