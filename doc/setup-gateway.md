# Setup a gateway with Elvisp

## On GNU/Linux
The first step is to decide inside which subnet(s) you want to delegate addresses for users.

In this tutorial we'll use 2 subnets, one for IPv6 and one for IPv4:
 * __IPv6__: fd12:3456::0/64
 * __IPv4__: 172.28.0.0/16

### Address for cjdns' TUN device
`cjdroute` does not automatically assign itself an IP address, therefore we'll set it manually. The server hosting `cjdroute` has probably already been given a IPv6 address, we'll assume it's somewhere between `fd12:3456::2` to `fd12:3456::9` and give the TUN interface the IPv6 address: `fd12:3456::10`.

Set the IPv6 address using:
```ìp -6 addr add dev tun0 fd12:3456::10```

The same probably applies to the IPv4 address. We'll use the address: `172.28.0.10`.

Set the IPv4 address using:
```ip -4 addr add dev tun0 172.28.0.10```

*__Note__: You're required to use these commands everytime cjdns starts/restarts, a good idea would be to automate this by adding it to cjdns' init file.*

### Route to ISP's gateway
The next step is to add a static route for the two subnets to allow them to access Internet. To do this, we route the subnets through a interface that has access to the Internet, in our case this is the `eth0` interface.

For IPv6 routing:
```
ip -6 route add dev eth0 fd12:3456::1
ip -6 route add dev tun0 fd12:3456::0/64
ip -6 route add dev default via fd12:3456::1
```

For IPv4 routing:
```
ip -4 route add dev eth0 172.28.0.1
ip -4 route add dev tun0 172.28.0.0/16
ip -4 route add default via 172.28.0.1
```

### Enable IP forwarding
The Linux kernel does not allow IP forwarding per default, to enable for both IPv6 and IPv4 run:
```
echo 1 > /proc/sys/net/ipv6/conf/all/forwarding
echo 1 > /proc/sys/net/ipv4/conf/all/forwarding
echo "net.ipv6.all.forwarding=1" >> /etc/sysctl.conf
echo "net.ipv4.all.forwarding=1" >> /etc/sysctl.conf
```

### Start Elvisp
Now the only thing left is to run Elvisp with the correct flags:
```./elvisp -cidr fd12:3456::10/64 -cidr 172.28.0.10/16 -password ElvispAdminPasswordHere -cjdns-password cjdnsAdminPasswordHere```

Change the cjdns password and set a good administration password. Also notice how the CIDR's start at `::10` and `.10`, Elvsip will start to lease IP addresses after these. Meaning the first user will get: `::11` and `.11`.

We're done!

### It doesn't work
Please see the documentation at [github.com/cjdelisle/cjdns/doc/tunnel.md#it-doesnt-work](https://github.com/cjdelisle/cjdns/blob/master/doc/tunnel.md#it-doesnt-work)

### References:
 * [github.com/cjdelisle/cjdns/doc/tunnel.md#running-a-gateway](https://github.com/cjdelisle/cjdns/blob/master/doc/tunnel.md#running-a-gateway) 
