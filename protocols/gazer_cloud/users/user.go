package users

import (
	"errors"
	"sync"
)

type User struct {
	mtx      sync.Mutex
	name     string
	password string
	sessions map[string]*Session
}

func NewUser(userName string, password string) *User {
	var c User
	c.name = userName
	c.password = password
	c.sessions = make(map[string]*Session)
	return &c
}

func (c *User) Dispose() {
	c.sessions = make(map[string]*Session)
}

func (c *User) OpenSession(password string) (*Session, error) {
	var err error
	var session *Session
	c.mtx.Lock()
	if c.password == password {
		session = NewSession(c)
		c.sessions[session.sessionId] = session
	} else {
		err = errors.New("wrong password")
	}
	c.mtx.Unlock()
	return session, err
}

func (c *User) CloseSession(sessionId string) error {
	var err error
	c.mtx.Lock()
	if session, ok := c.sessions[sessionId]; ok {
		session.Dispose()
		delete(c.sessions, sessionId)
	} else {
		err = errors.New("session not found in user")
	}
	c.mtx.Unlock()
	return err
}
