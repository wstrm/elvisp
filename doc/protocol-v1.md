# Protocol / API

### Status codes
* Error: 0
* Success: 1
### Standard error
```
{
  "error": "error-message-here",
  "status": 0
}
```

## Tasks
### Register

Send:
```
{
  "task": "register",
  "data": {
    "password": "password-for-user",
    "pubkey": "public-key-for-user.k",
    "token": "token-obtained-for-registration",
    "misc": "misc-information-for-user"
  }
}
```

Get:
```
{
  "error": null,
  "data": {
    "address": "2001::an::ipv6::address"
  },
  "status": 1
}
```

### Obtain lease

Send:
```
{
  "task": "lease",
  "data": {
    "password": "password-for-user",
    "pubkey": "public-key-for-user.k"
  }
}
```

Get:
```
{
  "error": null,
  "data": {
    "address": "2001::an::ipv6::address"
  },
  "status": 1
}
```

### Remove user

Send:
```
{
  "task": "remove",
  "data": {
    "password": "password-for-user",
    "pubkey": "public-key-for-user.k"
  }
}
```

Get:
```
{
  "error": null,
  "status": 1
}
```

### Release lease

Send:
```
{
  "task": "release",
  "data": {
    "password": "password-for-user",
    "pubkey": "public-key-for-user.k"
  }
}
```

Get:
```
{
  "error": null,
  "status": 1
}
```

### Add user

Send:
```
{
  "task": "add",
  "data": {
    "password": "master-password-for-admin",
    "pubkey": "public-key-for-new-user.k"
  }
}
```

Get:
```
{
  "error": null,
  "data": {
    "token": "token-for-user-registration"
  },
  "status": 1
}
```
