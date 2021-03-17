package system

import (
	"crypto"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/logger"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
	"github.com/gazercloud/gazernode/settings"
	"io/ioutil"
	"math/rand"
	"sort"
	"strconv"
	"time"
)

var DefaultUserName string
var DefaultUserPassword string

func init() {
	DefaultUserName = "admin"
	DefaultUserPassword = "admin"
}

type UserSession struct {
	SessionToken              string `json:"session_token"`
	UserName                  string `json:"user_name"`
	SessionOpenTime           int64  `json:"session_open_time"`
	SessionOpenTimeForDisplay string `json:"session_open_time_for_display"`
	Host                      string `json:"host"`
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

func (c *System) RemoveSession(sessionToken string) error {
	var err error
	c.mtx.Lock()
	if _, ok := c.sessions[sessionToken]; ok {
		delete(c.sessions, sessionToken)
	} else {
		err = errors.New("wrong session token")
	}
	c.saveSessions()

	c.mtx.Unlock()

	return err
}

func (c *System) SessionList(userName string) (nodeinterface.SessionListResponse, error) {
	var result nodeinterface.SessionListResponse
	var err error
	c.mtx.Lock()
	for _, s := range c.sessions {
		if s.UserName == userName {
			var item nodeinterface.SessionListResponseItem
			item.SessionToken = s.SessionToken
			item.UserName = s.UserName
			item.SessionOpenTime = s.SessionOpenTime
			result.Items = append(result.Items, item)
		}
	}
	c.mtx.Unlock()

	sort.Slice(result.Items, func(i, j int) bool {
		return result.Items[i].SessionOpenTime < result.Items[j].SessionOpenTime
	})

	return result, err
}

func (c *System) OpenSession(name string, password string, host string) (nodeinterface.SessionOpenResponse, error) {
	var result nodeinterface.SessionOpenResponse
	var err error

	c.mtx.Lock()
	if user, ok := c.userByName[name]; ok {
		if c.hashPassword(password) == user.PasswordHash {
			stringForHash := time.Now().Format("2006-01-02-15-04-05") + strconv.FormatInt(rand.Int63(), 10) + "42"
			sessionToken := c.hashSession(stringForHash)
			result.SessionToken = sessionToken

			timeOpenSession := time.Now().UTC()

			var ss UserSession
			ss.UserName = name
			ss.SessionToken = result.SessionToken
			ss.SessionOpenTime = timeOpenSession.UnixNano() / 1000
			ss.SessionOpenTimeForDisplay = timeOpenSession.Format("2006-01-02 15:04:05.999")
			ss.Host = host
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
		err = ioutil.WriteFile(settings.ServerDataPath()+"/sessions.json", bs, 0666)
		if err != nil {
			logger.Println("saveSessions error", err)
		}
	} else {
		logger.Println("saveSessions (marshal) error", err)
	}
}

func (c *System) loadSessions() {
	logger.Println("System loadSessions begin")
	configString, err := ioutil.ReadFile(settings.ServerDataPath() + "/sessions.json")
	if err == nil {
		err = json.Unmarshal(configString, &c.sessions)
		if err != nil {
			logger.Println("loadSessions (unmarshal) error ", err)
		}
	} else {
		logger.Println("loadSessions error ", err)
	}
	logger.Println("System loadSessions")
	logger.Println(c.sessions)
	logger.Println("System loadSessions end")
}

/*func (c *System) saveUsers() {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	bs, err := json.MarshalIndent(c.users, "", " ")
	if err == nil {
		err = ioutil.WriteFile(paths.ProgramDataFolder()+"/gazer/users.json", bs, 0666)
		if err != nil {
			logger.Println("saveUsers error", err)
		}
	} else {
		logger.Println("saveUsers (marshal) error", err)
	}
}

func (c *System) loadUsers() {
	logger.Println("System loadUsers begin")
	configString, err := ioutil.ReadFile(paths.ProgramDataFolder() + "/gazer/users.json")
	if err == nil {
		err = json.Unmarshal(configString, &c.users)
		if err != nil {
			logger.Println("loadUsers (unmarshal) error ", err)
		} else {
			c.userByName = make(map[string]*common_interfaces.User)
			for _, u := range c.users {
				c.userByName[u.Name] = u
			}
		}
	} else {
		logger.Println("loadUsers error ", err)
	}

	logger.Println("System loadUsers")
	for index, u := range c.users {
		logger.Println(index, ":", u.Name)
	}
	logger.Println("System loadUsers end")
}*/

func (c *System) UserList() (nodeinterface.UserListResponse, error) {
	var result nodeinterface.UserListResponse
	c.mtx.Lock()
	result.Items = make([]string, 0)
	for _, u := range c.users {
		result.Items = append(result.Items, u.Name)
	}
	c.mtx.Unlock()
	return result, nil
}

func (c *System) UserAdd(name string, password string) (nodeinterface.UserAddResponse, error) {
	var err error
	var result nodeinterface.UserAddResponse
	c.mtx.Lock()
	if _, ok := c.userByName[name]; !ok {
		var us common_interfaces.User
		us.Name = name
		us.PasswordHash = c.hashPassword(password)
		c.users = append(c.users, &us)
		c.userByName[us.Name] = &us
	} else {
		err = errors.New("user exists already")
	}
	c.mtx.Unlock()

	c.SaveConfig()

	return result, err
}

func (c *System) UserSetPassword(name string, password string) (nodeinterface.UserSetPasswordResponse, error) {
	var err error
	var result nodeinterface.UserSetPasswordResponse
	c.mtx.Lock()
	if u, ok := c.userByName[name]; ok {
		u.PasswordHash = c.hashPassword(password)
	} else {
		err = errors.New("no user found")
	}
	c.mtx.Unlock()

	c.SaveConfig()

	return result, err
}

func (c *System) hashPassword(password string) string {
	s := crypto.SHA256.New()
	return base64.StdEncoding.EncodeToString(s.Sum([]byte(password)))
}

func (c *System) hashSession(sessionData string) string {
	s := crypto.SHA256.New()
	return base64.StdEncoding.EncodeToString(s.Sum([]byte(sessionData)))
}

func (c *System) UserRemove(name string) (nodeinterface.UserRemoveResponse, error) {
	var err error
	var found bool
	var result nodeinterface.UserRemoveResponse
	c.mtx.Lock()
	for index, u := range c.users {
		if u.Name == name {
			c.users = append(c.users[:index], c.users[index+1:]...)
			delete(c.userByName, name)
			found = true
			break
		}
	}
	if !found {
		err = errors.New("no user found")
	}
	c.mtx.Unlock()

	c.SaveConfig()

	return result, err
}
