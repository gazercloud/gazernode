package local_user_storage

import (
	"encoding/json"
	"errors"
	"github.com/gazercloud/gazernode/logger"
	"io/ioutil"
	"os"
	"os/user"
	"sync"
)

type NodeConnection struct {
	Address      string `json:"address"`
	UserName     string `json:"user_name"`
	SessionToken string `json:"password"`
}

type PreferencesStruct struct {
	Theme       string           `json:"theme"`
	Connections []NodeConnection `json:"connections"`
}

type Preferences struct {
	pref PreferencesStruct
	mtx  sync.Mutex
}

var pref *Preferences

func init() {
	pref = NewPreferences()
	pref.loadPreferences()
}

func Instance() *Preferences {
	return pref
}

func NewPreferences() *Preferences {
	var c Preferences
	c.pref.Theme = "dark_blue"
	return &c
}

func (c *Preferences) loadPreferences() {
	usr, err := user.Current()
	if err != nil {
		logger.Println(err)
		return
	}
	dir := usr.HomeDir + "/Gazer"
	err = os.MkdirAll(dir, 0700)
	if err != nil {
		logger.Println(err)
		return
	}
	fullPath := dir + "/preferences.json"
	var bs []byte
	bs, err = ioutil.ReadFile(fullPath)
	if err != nil {
		logger.Println(err)
		return
	}
	err = json.Unmarshal(bs, &c.pref)
	if err != nil {
		logger.Println(err)
		return
	}
}

func (c *Preferences) savePreferences() {
	if c == nil {
		logger.Println(errors.New("pref == nil"))
	}
	usr, err := user.Current()
	if err != nil {
		logger.Println(err)
		return
	}
	dir := usr.HomeDir + "/Gazer"
	err = os.MkdirAll(dir, 0700)
	if err != nil {
		logger.Println(err)
		return
	}
	fullPath := dir + "/preferences.json"
	var bs []byte
	bs, err = json.MarshalIndent(pref.pref, "", " ")
	err = ioutil.WriteFile(fullPath, bs, 0600)
}

func (c *Preferences) SetTheme(theme string) {
	c.mtx.Lock()
	c.pref.Theme = theme
	c.savePreferences()
	c.mtx.Unlock()
}

func (c *Preferences) Theme() string {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	return c.pref.Theme
}

func (c *Preferences) AddConnection(connection NodeConnection) {
	c.mtx.Lock()
	c.pref.Connections = append(c.pref.Connections, connection)
	c.savePreferences()
	c.mtx.Unlock()
}

func (c *Preferences) SetConnection(index int, connection NodeConnection) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if index >= 0 && index < len(c.pref.Connections) {
		c.pref.Connections[index] = connection
	}
	c.savePreferences()
}

func (c *Preferences) ConnectionCount() int {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	return len(c.pref.Connections)
}

func (c *Preferences) Connections() []NodeConnection {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	result := make([]NodeConnection, len(c.pref.Connections))
	for i, conn := range c.pref.Connections {
		result[i] = conn
	}
	return result
}

func (c *Preferences) RemoveConnection(index int) {
	c.mtx.Lock()
	if index >= 0 && index < len(c.pref.Connections) {
		c.pref.Connections = append(c.pref.Connections[:index], c.pref.Connections[index+1:]...)
		c.savePreferences()
	}
	c.mtx.Unlock()
}
