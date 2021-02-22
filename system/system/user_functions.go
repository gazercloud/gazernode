package system

import (
	"crypto"
	"encoding/json"
	"errors"
	"github.com/gazercloud/gazernode/logger"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
	"github.com/gazercloud/gazernode/utilities/paths"
	"io/ioutil"
	"time"
)

type User struct {
	Name         string
	PasswordHash string
}

var DefaultUserName string
var DefaultUserPassword string

func init() {
	DefaultUserName = "admin"
	DefaultUserPassword = "admin"
}

type UserSession struct {
	SessionToken string
	UserName     string
}

func (c *System) UserAdd(name string, password string) error {
	c.mtx.Lock()
	c.mtx.Unlock()
	return nil
}

func (c *System) UserRename(oldName string, newName string) error {
	c.mtx.Lock()
	c.mtx.Unlock()
	return nil
}

func (c *System) UserSetPassword(name string, password string) error {
	c.mtx.Lock()
	c.mtx.Unlock()
	return nil
}

func (c *System) hashPassword(password string) string {
	password = password + "_salt"
	s := crypto.SHA256.New()
	return string(s.Sum([]byte(password)))
}

func (c *System) UserRemove(name string) error {
	c.mtx.Lock()
	c.mtx.Unlock()
	return nil
}

func (c *System) CheckSession(sessionToken string) (string, error) {
	var userName string
	var err error
	c.mtx.Lock()
	if session, ok := c.sessions[sessionToken]; ok {
		userName = session.UserName
	} else {
		err = errors.New("wrong session token")
	}
	c.mtx.Unlock()
	return userName, err
}

func (c *System) OpenSession(name string, password string) (nodeinterface.SessionOpenResponse, error) {
	var result nodeinterface.SessionOpenResponse
	var err error

	c.mtx.Lock()
	if user, ok := c.userByName[name]; ok {
		if c.hashPassword(password) == c.hashPassword("123") {
			result.SessionToken = "session_" + user.Name + "_" + time.Now().Format("2006-01-02-15-04-05")
			var ss UserSession
			ss.UserName = name
			ss.SessionToken = result.SessionToken
			c.sessions[result.SessionToken] = &ss
		} else {
			err = errors.New("wrong password")
		}
	} else {
		err = errors.New("user not found")
	}

	c.saveSessions()

	c.mtx.Unlock()
	return result, err
}

func (c *System) saveSessions() {
	bs, err := json.MarshalIndent(c.sessions, "", " ")
	if err == nil {
		err = ioutil.WriteFile(paths.ProgramDataFolder()+"/gazer/sessions.json", bs, 0666)
		if err != nil {
			logger.Println("saveSessions error", err)
		}
	}
}

func (c *System) loadSessions() {
	configString, err := ioutil.ReadFile(paths.ProgramDataFolder() + "/gazer/sessions.json")
	if err == nil {
		err = json.Unmarshal(configString, &c.sessions)
		if err != nil {
			logger.Println("loadSessions error ", err)
		}
	}
}
