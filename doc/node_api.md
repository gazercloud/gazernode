# API Groups
- UnitType
- Unit
- DataItem
- Cloud
- Service
- Resource
- User

# API reference
## UnitType
## Unit
## DataItem
## Cloud
## Service
## Resource

## User
___
### user_list
#### input
```
no
```
#### output
```
{ "items" : [ "user1", "user2" ... ] }
```
___
### user_add
#### input
```
user_name
password
```
#### output
```
no
```
### user_set_password
#### input
```
user_name
password
```
#### output
```
no
```

### user_remove
#### input
```
user_name
```
#### output
```
no
```





### session_open
#### input
```
user_name
password
```
#### output
```
session_token
```
### session_remove
#### input
```
session_token
```
#### output
```
no
```
### session_list
#### input
```
user_name
```
#### output
```
items:
[
    {
        session_token : string,
        user_name : string,
        session_open_time: number, // unix time
    },
]
```
