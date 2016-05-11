# Protocol / API

### Error codes
  * `1` Invalid task specified
  * `2` Invalid length of arguments
### Standard error
```
error <error code> <error message>
```

## Tasks

### Obtain lease

Send (from user node):
```
lease
```

Send (using admin):
```
lease <master-password-for-admin> <public-key-for-user.k>
```

Get:
```
success ipv6 <ipv6-address-here>
success ipv4 <ipv4-address-here>
```

### Remove user

Send (from user node):
```
remove
```

Send (using admin):
```
remove <master-password-for-admin> <public-key-for-user.k>
```

Get:
```
success Removed user: <public-key-for-user.k>
```

### Release lease

Send (from user node):
```
release
```

Send (using admin):
```
release <master-password-for-admin> <public-key-for-user.k>
```

Get:
```
success Released lease for user: <public-key-for-user.k>
```

### Add user

Send (from user node):
```
add
```

Send (using admin):
```
add <master-password-for-admin> <public-key-for-user.k>
```

Get:
```
sucess Added user: <public-key-for-user.k>
```
