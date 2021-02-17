package client

import (
	"github.com/gazercloud/gazernode/common_interfaces"
	"sync"
	"time"
)

type ItemsWatcherItem struct {
	Name        string
	value       common_interfaces.ItemValue
	lastGetTime time.Time
}

type ItemsWatcher struct {
	mtx      sync.Mutex
	client   *Client
	stopping bool
	items    map[string]*ItemsWatcherItem
}

func NewItemsWatcher(client *Client) *ItemsWatcher {
	var c ItemsWatcher
	c.client = client
	c.items = make(map[string]*ItemsWatcherItem)
	c.Start()
	return &c
}

func (c *ItemsWatcher) Start() {
	go c.thWorker()
}

func (c *ItemsWatcher) Stop() {
	c.stopping = true
}

func (c *ItemsWatcher) Get(itemName string) common_interfaces.ItemValue {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if foundItem, ok := c.items[itemName]; ok {
		foundItem.lastGetTime = time.Now().UTC()
		return foundItem.value
	} else {
		var i ItemsWatcherItem
		i.Name = itemName
		i.lastGetTime = time.Now().UTC()
		c.items[itemName] = &i
	}
	return common_interfaces.ItemValue{}
}

func (c *ItemsWatcher) thWorker() {
	for !c.stopping {
		time.Sleep(500 * time.Millisecond)
		is := make([]string, 0)
		c.mtx.Lock()
		if len(c.items) > 0 {
			found := true
			for found {
				found = false
				for n, i := range c.items {
					if time.Now().UTC().Sub(i.lastGetTime) > 5*time.Second {
						delete(c.items, n)
						found = true
						break
					}
				}
			}

			for _, i := range c.items {
				is = append(is, i.Name)
			}
			c.client.GetItemsValues(is, func(items []common_interfaces.ItemGetUnitItems, err error) {
				c.mtx.Lock()
				for _, i := range items {
					if wi, ok := c.items[i.Name]; ok {
						wi.value = i.Value
					}
				}
				c.mtx.Unlock()
			})
		}
		c.mtx.Unlock()
	}
}
