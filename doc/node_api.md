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
#### description
```
returns list of user's names.
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
#### description
```
creates a user
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
#### description
```
changes user's password
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
#### description
```
removes a user
```
