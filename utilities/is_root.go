package utilities

import (
	"os/user"
)

func IsRoot() bool {
	u, _ := user.Current()
	if u.Uid == "0" || u.Username == "root" {
		return true
	}
	return false
}
