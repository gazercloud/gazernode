package crunner

import (
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"sync"
)

type CRunner struct {
	window   uiinterfaces.Window
	received []*Call
	mtx      sync.Mutex
	tm       *uievents.FormTimer
}

type Call struct {
	request interface{}

	response interface{}
	err      error

	thFunc     func(interface{}) (interface{}, error)
	onResponse func(interface{}, error)
}

func New(window uiinterfaces.Window) *CRunner {
	var c CRunner
	c.initRunner(window)
	return &c
}

func (c *CRunner) Call(thFunc func(thParameters interface{}) (interface{}, error), resultFunc func(result interface{}, err error), parameters ...interface{}) {
	var call Call
	call.thFunc = thFunc
	call.onResponse = resultFunc
	call.request = parameters
	go c.thCall(&call)
}

func (c *CRunner) initRunner(window uiinterfaces.Window) {
	c.tm = window.NewTimer(100, c.timer)
	c.tm.StartTimer()
}

func (c *CRunner) timer() {
	c.mtx.Lock()
	for _, call := range c.received {
		call.onResponse(call.response, call.err)
	}
	c.received = make([]*Call, 0)
	c.mtx.Unlock()
}

func (c *CRunner) thCall(call *Call) {
	call.response, call.err = call.thFunc(call.request)
	c.mtx.Lock()
	c.received = append(c.received, call)
	c.mtx.Unlock()
}
