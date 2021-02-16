package users

import (
	"errors"
	"github.com/gazercloud/gazernode/utilities/state"
	"sort"
	"sync"
)

type Users struct {
	mtx      sync.Mutex
	users    map[string]*User
	sessions map[string]*Session
}

func NewUsers() *Users {
	var c Users
	c.users = make(map[string]*User, 0)
	c.sessions = make(map[string]*Session)

	return &c
}

func (c *Users) AddUser(userName string, password string) (*User, error) {
	var user *User
	c.mtx.Lock()
	if _, ok := c.users[userName]; !ok {
		u := NewUser(userName, password)
		c.users[userName] = u
		user = u
	}
	c.mtx.Unlock()
	if user == nil {
		return nil, errors.New("user already exists")
	}
	return user, nil
}

func (c *Users) RemoveUser(userName string) error {
	var err error
	c.mtx.Lock()
	if user, ok := c.users[userName]; ok {
		for _, session := range user.sessions {
			err = user.CloseSession(session.sessionId)
			if err != nil {
				break
			}
			delete(c.sessions, session.sessionId)
		}
		if err == nil {
			user.Dispose()
			delete(c.users, userName)
		}
	} else {
		err = errors.New("user not found")
	}
	c.mtx.Unlock()
	return err
}

func (c *Users) OpenSession(userName string, password string) (*Session, error) {
	var err error
	var session *Session
	c.mtx.Lock()
	if user, ok := c.users[userName]; ok {
		session, err = user.OpenSession(password)
		if err == nil {
			c.sessions[session.sessionId] = session
		}
	} else {
		err = errors.New("user not found")
	}
	c.mtx.Unlock()
	return session, err
}

func (c *Users) CloseSession(sessionId string) error {
	var err error
	c.mtx.Lock()
	if session, ok := c.sessions[sessionId]; ok {
		err = session.user.CloseSession(sessionId)
		if err == nil {
			delete(c.sessions, sessionId)
		}
	} else {
		err = errors.New("session not found")
	}
	c.mtx.Unlock()
	return err
}

func (c *Users) CheckUser(userName string, password string) error {
	var err error
	c.mtx.Lock()
	if user, ok := c.users[userName]; ok {
		if user.password != password {
			err = errors.New("wrong password")
		}
	} else {
		err = errors.New("wrong user")
	}
	c.mtx.Unlock()
	return err
}

func (c *Users) Session(sessionId string) (*Session, error) {
	var session *Session
	var err error
	var ok bool
	c.mtx.Lock()
	if session, ok = c.sessions[sessionId]; !ok {
		err = errors.New("unknown session")
	}
	c.mtx.Unlock()
	return session, err
}

func (c *Users) State() *state.Users {
	var st state.Users
	c.mtx.Lock()
	st.Users = make([]state.User, len(c.users))
	index := 0
	for _, user := range c.users {
		st.Users[index].Name = user.name
		st.Users[index].Sessions = make([]state.Session, len(user.sessions))
		sessionIndex := 0
		for _, session := range user.sessions {
			st.Users[index].Sessions[sessionIndex].Id = session.sessionId
			st.Users[index].Sessions[sessionIndex].UserName = session.user.name
			st.Users[index].Sessions[sessionIndex].BeginTime = session.beginTime
			st.Users[index].Sessions[sessionIndex].Disposed = session.disposed
			sessionIndex++
		}
		sort.Slice(st.Users[index].Sessions, func(i, j int) bool {
			if st.Users[index].Sessions[i].Id < st.Users[index].Sessions[j].Id {
				return true
			}
			return false
		})
		index++
	}
	sort.Slice(st.Users, func(i, j int) bool {
		if st.Users[i].Name < st.Users[j].Name {
			return true
		}
		return false
	})

	c.mtx.Unlock()
	return &st
}
