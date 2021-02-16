package state

import "time"

type Session struct {
	Disposed  bool      `json:"disposed"`
	Id        string    `json:"id"`
	BeginTime time.Time `json:"begin_time"`
	UserName  string    `json:"user_name"`
}

type User struct {
	Name     string    `json:"name"`
	Sessions []Session `json:"sessions"`
}

type Users struct {
	Users []User `json:"users"`
}
