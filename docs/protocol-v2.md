# Protocol / API v2

### Standard error
```
error <error message>
```

### Standard success
```
success <success message>
```

## Tasks

### Obtain lease

Send (from user node):
```
lease
```

Send (using admin):
```
lease <master-password-for-admin> <cjdns-ipv6-address>
```

Get:
```
success <ipv4-address-here> <ipv6-address-here>
```

### Remove user

Send (from user node):
```
remove
```

Send (using admin):
```
remove <master-password-for-admin> <cjdns-ipv6-address>
```

Get:
```
success Removed user: <public-key-for-user.k>
```

### Retrieve server info

Send (from user node or admin):
```
info
```

Get:
```
success <public-key-for-server.k>
```
