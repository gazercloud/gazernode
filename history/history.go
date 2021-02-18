package history

import (
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/logger"
	"sync"
	"time"
)

type History struct {
	items          map[uint64]*Item
	mtx            sync.Mutex
	flushPeriodSec int
	started        bool
	stopping       bool
}

func NewHistory() *History {
	var c History
	c.items = make(map[uint64]*Item)
	c.flushPeriodSec = 10
	return &c
}

func (c *History) Start() {
	logger.Println("HISTORY starting begin")

	c.started = true
	c.stopping = false
	go c.thWorker()

	logger.Println("HISTORY starting end")
}

func (c *History) Stop() {
	logger.Println("HISTORY stopping begin")

	c.stopping = true

	for i := 0; i < 10; i++ {
		time.Sleep(100 * time.Millisecond)
		if !c.started {
			break
		}
	}

	logger.Println("HISTORY flushing ...")
	c.mtx.Lock()
	for _, item := range c.items {
		logger.Println("HISTORY flushing ", item.id)
		res := item.Flush()
		logger.Println("HISTORY flushed ", item.id, "OK items:", res.CountOfItems, "data size:", res.FullDataSize)
	}
	c.mtx.Unlock()
	logger.Println("HISTORY flushing ... OK")

	if c.started {
		logger.Println("HISTORY stopping: timeout")
	}
	logger.Println("HISTORY stopping end")
}

func (c *History) thWorker() {
	logger.Println("HISTORY worker begin")

	lastFlushDT := time.Now()

	for !c.stopping {
		for time.Now().Sub(lastFlushDT) < time.Duration(c.flushPeriodSec)*time.Second {
			time.Sleep(100 * time.Millisecond)
			if c.stopping {
				break
			}
		}

		if c.stopping {
			break
		}

		lastFlushDT = time.Now()

		c.mtx.Lock()
		items := make([]*Item, 0)
		for _, item := range c.items {
			items = append(items, item)
		}
		c.mtx.Unlock()

		for _, item := range items {
			item.Flush()
			item.CheckDepth()
		}
	}

	c.started = false
	logger.Println("HISTORY worker end")
}

func (c *History) Write(id uint64, value common_interfaces.ItemValue) {
	var item *Item
	var ok bool

	if c.stopping || !c.started {
		return
	}

	c.mtx.Lock()
	item, ok = c.items[id]
	if !ok {
		item = NewItem(id)
		c.items[id] = item
	}
	c.mtx.Unlock()

	item.Write(value)
}

func (c *History) Read(id uint64, dtBegin int64, dtEnd int64) *ReadResult {
	var item *Item
	var ok bool

	if c.stopping || !c.started {
		var result ReadResult
		return &result
	}

	c.mtx.Lock()
	item, ok = c.items[id]
	c.mtx.Unlock()

	if !ok {
		var result ReadResult
		result.Id = id
		result.DTBegin = dtBegin
		result.DTEnd = dtEnd
		result.Items = make([]*common_interfaces.ItemValue, 0)
		return &result
	}

	result := item.Read(dtBegin, dtEnd)
	result.DTBegin = dtBegin
	result.DTEnd = dtEnd
	return result
}
