![Elvisp](elvisp.png)

# Elvisp (Work in Progress)
Elvisp distributes a public IPv6 address in a cjdns-based IP tunnel. Elvisp uses the public key that each user provides to store it in a database. Elvisp adds the user's assigned address with cjdns' admin API. Elvisp then return the public key for the cjdns instance that acts like the tunnel.

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
Data that is sent should be formated as JSON, and contain a password for authorization to the server, and a pubkey for the registration (also used as an user name). The misc field is optional, but can contain any string.
```
{
  "password": "examplePassword",
  "pubkey": "examplePubkey.k",
  "misc": "Fullname: John Doe, Email: john@doe"
}
```
###Receive
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
    password; 'Auth Pasword',
    pubkey: 'Public key for the user',
    misc: 'Misc information, optional'
  }));
});
```
