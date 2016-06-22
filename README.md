![Elvisp](elvisp.png)

# Elvisp (Work in Progress) [![Build Status](https://travis-ci.org/willeponken/elvisp.svg?branch=master)](https://travis-ci.org/willeponken/elvisp)
Elvisp distributes a public IPv6 address in a cjdns-based IP tunnel using the public key that each user provides. It will add the user's assigned address with cjdns' admin API. Elvisp then returns the public key for the cjdns instance that acts like the tunnel.

### Flags
```
Usage of ./elvisp:
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
  -v	Print current version and exit.
```
__Example:__
```
./elvisp -cidr 192.168.1.0/24 -cidr 1234::0/16 -cidr 172.16.0.0/12 -cjdns-password 6c12zbnNoThisIsntMyRealPasswordn7x1
```

Which will add the users / nodes to subnets:
 * 192.168.1.0/24
 * 1234::0/16
 * 172.16.0.0/12

So the first user will get:
 * 192.168.1.1
 * 1234::1
 * 172.16.0.1

### Supported cjdns versions
__Elvisp requires the follwing cjdns admin methods:__
 * `IpTunnel_allowConnection`
 * `IpTunnel_listConnections`
 * `IpTunnel_showConnection`
 * `IpTunnel_removeConnection`
 * `NodeStore_nodeForAddr`

__Note:__
*`IpTunnel_removeConnection` was first implemented with commit `acbb6a8` into the `crashey` branch. As of 2016-06-16 it has not been merged into the `master` branch.*

__Elvisp works with *(kinda)*:__
```
Cjdns version: cjdns-v17.3-129-g116fa2a
Cjdns protocol version: 17
```

### Protocol
See [protocol-v1](doc/protocol-v1.md).
