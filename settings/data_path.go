package settings

import (
	"fmt"
	"os/user"
	"path/filepath"
	"strings"
)

type Settings struct {
	serverDataPath string
}

func NewSettings() *Settings {
	var c Settings
	c.serverDataPath = "~/gazernode"
	return &c
}

func (c *Settings) SetServerDataPath(path string) {

	usr, _ := user.Current()
	dir := usr.HomeDir

	if path == "~" {
		path = dir
	} else if strings.HasPrefix(path, "~/") {
		path = filepath.Join(dir, path[2:])
	}

	c.serverDataPath = path
	fmt.Println("Server Path:", c.serverDataPath)
}

func (c *Settings) ServerDataPath() string {
	return c.serverDataPath
}
