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
success <ipv4-address-here> <ipv6-address-here>
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
success <ipv4-address-here> <ipv6-address-here> 
```
