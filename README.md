# iptd (Work in Progress)
IP Tunnel daemon for CJDNS

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
