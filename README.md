# iptd (Work in Progress)
IPTd is developed to distribute a public IPv6 address in a CJDNS-based IP tunnel. IPTd uses the public key that each user provides to store it in a database. IPTd adds the user's assigned address with CJDNS' admin API. IPTd then return the public key for the CJDNS instance that acts like the tunnel.

##TODO
* Registration
* Unit testing
* Both IPv4/6 at the same time
* Documentation

##API
####Status codes
* Error: 0
* Success: 1

###Send
Data that is sent should be formated as JSON, and contain password for authorization to the server, and a pubkey for the registration. The misc field is optional, but can contain any string.
```
{
  "password": "examplePassword",
  "pubkey": "examplePubkey.k",
  "misc": "Fullname: John Doe, Email: john@doe"
}
```
###Recieve
####Success
```
{
  "error": null,
  "data": {
    "address": "IPv6 address",
    "pubkey": "serverPubKey.k"
  },
  "status": 1
}
```

####Error
```
{
  "error": "Error message",
  "status": 0
}
```

###Example
```
var net = require('net');

var client = net.connect({ port: 4123 }, function() {
  client.write(JSON.stringify({
    password; 'Server Pasword',
    pubkey: 'Public key for new user registration',
    misc: 'Random information, optional'
  }));
});
```
