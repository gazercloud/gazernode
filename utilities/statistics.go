package utilities

import "sync"

type Statistics struct {
	mtx    sync.Mutex
	values map[string]int
}

func NewStatistics() *Statistics {
	var c Statistics
	c.Reset()
	return &c
}

func (c *Statistics) Reset() {
	c.mtx.Lock()
	c.values = make(map[string]int)
	c.mtx.Unlock()
}

func (c *Statistics) Add(code string, value int) {
	c.mtx.Lock()
	if _, ok := c.values[code]; ok {
		c.values[code] += value
	} else {
		c.values[code] = value
	}
	c.mtx.Unlock()
}

func (c *Statistics) Get(code string) int {
	result := 0
	c.mtx.Lock()
	if _, ok := c.values[code]; ok {
		result = c.values[code]
	}
	c.mtx.Unlock()
	return result
}
