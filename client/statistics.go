package client

import (
	"sync"
)

type Statistics struct {
	mtx         sync.Mutex
	secReceived int
	secSent     int
}

var instance Statistics

func init() {
	StatReset()
}

func StatReset() {
	instance.mtx.Lock()
	instance.secReceived = 0
	instance.secSent = 0
	instance.mtx.Unlock()
}

func AddStatReceived(receivedBytes int) {
	instance.mtx.Lock()
	instance.secReceived += receivedBytes
	instance.mtx.Unlock()
}

func AddStatSent(sentBytes int) {
	instance.mtx.Lock()
	instance.secSent += sentBytes
	instance.mtx.Unlock()
}

func StatReceived() int {
	return instance.secReceived
}
func StatSent() int {
	return instance.secSent
}
