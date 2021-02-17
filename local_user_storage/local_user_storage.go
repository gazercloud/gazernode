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

type Preferences struct {
	Theme string `json:"theme"`
}

var mtx sync.Mutex

func NewPreferences() *Preferences {
	var c Preferences
	c.Theme = "dark_blue"
	return &c
}

func loadPreferences() *Preferences {
	usr, err := user.Current()
	if err != nil {
		return NewPreferences()
	}
	dir := usr.HomeDir + "/Gazer"
	err = os.MkdirAll(dir, 0600)
	if err != nil {
		return NewPreferences()
	}
	fullPath := dir + "/preferences.json"
	var bs []byte
	bs, err = ioutil.ReadFile(fullPath)
	if err != nil {
		return NewPreferences()
	}
	var result Preferences
	err = json.Unmarshal(bs, &result)
	if err != nil {
		return NewPreferences()
	}
	return &result
}

func savePreferences(pref *Preferences) error {
	if pref == nil {
		return errors.New("pref == nil")
	}
	usr, err := user.Current()
	if err != nil {
		return err
	}
	dir := usr.HomeDir + "/Gazer"
	err = os.MkdirAll(dir, 0600)
	if err != nil {
		return err
	}
	fullPath := dir + "/preferences.json"
	var bs []byte
	bs, err = json.MarshalIndent(pref, "", " ")
	err = ioutil.WriteFile(fullPath, bs, 0600)
	return err
}

func SetTheme(theme string) {
	mtx.Lock()
	p := loadPreferences()
	p.Theme = theme
	if err := savePreferences(p); err != nil {
		logger.Println("Set There error:", err)
	}
	mtx.Unlock()
}

func Theme() string {
	mtx.Lock()
	p := loadPreferences()
	mtx.Unlock()
	return p.Theme
}
