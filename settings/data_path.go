package settings

import (
	"fmt"
	"os/user"
	"path/filepath"
	"strings"
)

var serverDataPath = "~/gazer_node"

func SetServerDataPath(path string) {

	usr, _ := user.Current()
	dir := usr.HomeDir

	if path == "~" {
		path = dir
	} else if strings.HasPrefix(path, "~/") {
		path = filepath.Join(dir, path[2:])
	}

	serverDataPath = path
	fmt.Println("Server Path:", serverDataPath)
}

func ServerDataPath() string {
	return serverDataPath
}
