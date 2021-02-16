package users

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"strconv"
	"time"
)

type Session struct {
	disposed  bool
	sessionId string
	beginTime time.Time
	user      *User
}

func NewSession(user *User) *Session {
	var c Session
	c.user = user
	c.beginTime = time.Now().UTC()
	hashId := user.name + strconv.Itoa(rand.Int()) + time.Now().String()
	hash := sha256.New()
	hash.Write([]byte(hashId))
	sha := hex.EncodeToString(hash.Sum(nil))
	c.sessionId = user.name + "_" + sha
	return &c
}

func (c *Session) Dispose() {
	c.user = nil
	c.disposed = true
}

func (c *Session) Id() string {
	return c.sessionId
}
