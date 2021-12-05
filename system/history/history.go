package history

import (
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/system/settings"
	"github.com/gazercloud/gazernode/utilities/logger"
	"sync"
	"time"
)

type History struct {
	ss             *settings.Settings
	items          map[uint64]*Item
	mtx            sync.Mutex
	flushPeriodSec int
	started        bool
	stopping       bool
}

func NewHistory(ss *settings.Settings) *History {
	var c History
	c.ss = ss
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
		go item.FinishFlush()
	}
	for {
		countOfFinished := 0
		for _, item := range c.items {
			if item.flushFinished {
				countOfFinished++
			}
		}
		if countOfFinished == len(c.items) {
			break
		}
		logger.Println("HISTORY flushing", countOfFinished, " of ", len(c.items), "finished")
		time.Sleep(250 * time.Millisecond)
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
		item = NewItem(id, c.ss)
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

func (c *History) RemoveItem(id uint64) {
	c.mtx.Lock()
	item, ok := c.items[id]
	if ok {
		item.Remove()
		delete(c.items, id)
	}
	c.mtx.Unlock()
}
