![Elvisp](elvisp.png)

# Elvisp (Work in Progress) [![Build Status](https://travis-ci.org/willeponken/elvisp.svg?branch=master)](https://travis-ci.org/willeponken/elvisp) [![Coverage Status](https://coveralls.io/repos/github/willeponken/elvisp/badge.svg?branch=master)](https://coveralls.io/github/willeponken/elvisp?branch=master)
**El** (Spanish for *the*) **VISP** (*Virtual Internet Service Provider*)

Elvisp assigns IPv6 and/or IPv4 addresses for a cjdns-based IP tunnel using the public key for the connecting node. It will add the user's assigned address with cjdns' admin API. Elvisp then returns the assigned address(es).

## Installation
```
go get github.com/willeponken/elvisp/cmd/...
```

## Usage
### Elvispd flags
```
Usage of elvispd:
  -cidr value
    	CIDR to use for IP leasing, use flag repeatedly for multiple CIDR's.
  -cjdns-ip string
    	IP address for cjdns admin. (default "127.0.0.1")
  -cjdns-password string
    	Password for cjdns admin.
  -cjdns-port int
    	Port for cjdns admin. (default 11234)
  -db string
    	Directory to use for the database. (default "/tmp/elvisp-db")
  -listen string
    	Listen address for TCP. (default ":4132")
  -password string
    	Password for administrating Elvisp.
```
__Example:__
```
elvispd -cidr 192.168.1.0/24 -cidr 1234::0/16 -cidr 172.16.0.0/12 -cjdns-password 6c12zbnNoThisIsntMyRealPasswordn7x1
```

Which will add the users / nodes to subnets:
 * 192.168.1.0/24
 * 1234::0/16
 * 172.16.0.0/12

So the first user will get:
 * 192.168.1.1
 * 1234::1
 * 172.16.0.1

### Elvispc flags
```
Usage of elvispc:
  -a string
    	Address for server.
  -l	Request lease.
  -r	Remove client.
```
__Example:__
```
elvispc -a 127.0.0.1:4132 -l # Request lease
elvispc -a 127.0.0.1:4132 -r # Remove client
```

### Supported cjdns versions
__Elvisp requires the follwing cjdns admin methods:__
 * `IpTunnel_allowConnection`
 * `IpTunnel_listConnections`
 * `IpTunnel_showConnection`
 * `IpTunnel_removeConnection`
 * `NodeStore_nodeForAddr`

__Elvisp works with (atleast):__
```
Cjdns version: cjdns-v17.4
Cjdns protocol version: 17
```

### Documentation
 * Protocol [protocol-v2](docs/protocol-v2.md)
 * Setup a gateway [setup-gateway](docs/setup-gateway.md)
