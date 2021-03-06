package tasks

import (
	"fmt"
	"log"
	"net"

	"github.com/willeponken/elvisp/cjdns"
	"github.com/willeponken/elvisp/database"
	"github.com/willeponken/elvisp/lease"
	"github.com/willeponken/go-cjdns/key"
)

// TaskInterface defines the methods needed for a default task
type TaskInterface interface {
	Run() (result string, err error)
}

// Task needs the arguments to use, and a database to save the changes to
type Task struct {
	argv                 []string
	db                   *database.Database
	admin                *cjdns.Conn
	clientIP, serverIP   net.IP
	clientKey, serverKey *key.Public
	cidrs                []lease.CIDR
}

// Init returns a new task
func Init(argv []string, db *database.Database, admin *cjdns.Conn, clientIP, serverIP net.IP, cidrs []lease.CIDR) (task Task, err error) {
	task.argv = argv
	task.db = db
	task.admin = admin
	task.clientIP = clientIP
	task.cidrs = cidrs

	var clientKey, serverKey string
	clientKey, err = task.admin.LookupPubKey(clientIP.String())
	serverKey, err = task.admin.LookupPubKey(serverIP.String())
	if err != nil {
		return task, err
	}

	task.clientKey, err = key.DecodePublic(clientKey)
	task.serverKey, err = key.DecodePublic(serverKey)
	if err != nil {
		return task, err
	}

	return
}

// Remove should implement the remove task
type Remove struct{ Task }

// Lease should implement the lease task
type Lease struct{ Task }

// Release should implement the release task
type Release struct{ Task }

// Info should implement the info task
type Info struct{ Task }

// Invalid should implement the invalid task, i.e. take an error
type Invalid struct{ Error error }

// allowIPTunnel adds the defined IP to the cjdns IP tunnel, if it fails, it will delete the user from the database.
func (t Lease) allowIPTunnel(ips []net.IP) (err error) {

	for _, ip := range ips {
		if err = t.admin.AddUser(t.clientKey, ip); err != nil {
			if e := t.db.DelUser(t.clientKey); e != nil {
				log.Println(err)
			}

			return
		}
	}

	return
}

func (t Lease) generateIPs(cidrs []lease.CIDR, id uint64) (ips []net.IP, str string, err error) {
	var ip net.IP
	for _, cidr := range cidrs {
		ip, err = lease.Generate(cidr, id)
		if err != nil {
			return
		}

		ips = append(ips, ip)
		str += ip.String() + " "
	}

	return
}

// Run Lease adds a user using the public key and a token.
func (t Lease) Run() (result string, err error) {
	var id uint64
	var ips []net.IP
	db := t.db

	// Check if the user already exists
	id, exists := db.GetID(t.clientKey)
	if exists != nil { // User does not exist, add to database
		id, err = db.AddUser(t.clientKey)
		if err != nil {
			return
		}
	}

	ips, result, err = t.generateIPs(t.cidrs, id)
	if err != nil {
		return
	}

	err = t.allowIPTunnel(ips)
	if err != nil {
		return
	}

	return
}

// Run Remove removes a user.
func (t Remove) Run() (result string, err error) {
	db := t.db
	admin := t.admin
	pubkey := t.clientKey

	if err = db.DelUser(pubkey); err != nil {
		return
	}

	if err = admin.DelUser(pubkey); err != nil {
		return
	}

	result = fmt.Sprintf("Removed user: %s", pubkey.String())
	return
}

// Run Info returns information about the Elvisp server
func (t Info) Run() (result string, err error) {
	log.Println(t)
	result = t.serverKey.String()
	return
}

// Run Invalid returns an error and empty result.
func (t Invalid) Run() (result string, err error) {
	err = t.Error
	return
}
