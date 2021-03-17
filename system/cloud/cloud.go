package cloud

import (
	"errors"
	"github.com/gazercloud/gazernode/common_interfaces"
	"sync"
)

type Cloud struct {
	channels     []*Channel
	mtx          sync.Mutex
	iDataStorage common_interfaces.IDataStorage
}

func NewCloud(iDataStorage common_interfaces.IDataStorage) *Cloud {
	var c Cloud
	c.iDataStorage = iDataStorage
	c.channels = make([]*Channel, 0)
	return &c
}

func (c *Cloud) Start() {
	c.mtx.Lock()
	for _, ch := range c.channels {
		ch.Start()
	}
	c.mtx.Unlock()
}

func (c *Cloud) Stop() {
	c.mtx.Lock()
	for _, ch := range c.channels {
		ch.Stop()
	}
	c.mtx.Unlock()
}

func (c *Cloud) ReceivedBytes() int {
	var result int
	c.mtx.Lock()
	for _, ch := range c.channels {
		stat := ch.Stat()
		if stat != nil {
			result += stat.Get("rcv")
		}
	}
	c.mtx.Unlock()
	return result
}

func (c *Cloud) SentBytes() int {
	var result int
	c.mtx.Lock()
	for _, ch := range c.channels {
		stat := ch.Stat()
		if stat != nil {
			result += stat.Get("snd")
		}
	}
	c.mtx.Unlock()
	return result
}

func (c *Cloud) ChannelsFullInfo() []ChannelFullInfo {
	channels := make([]ChannelFullInfo, 0)
	c.mtx.Lock()
	for _, ch := range c.channels {
		var item ChannelFullInfo
		item.Id = ch.channelId
		item.Name = ch.name
		item.Password = ch.password
		item.Items = ch.items
		channels = append(channels, item)
	}

	c.mtx.Unlock()
	return channels
}

func (c *Cloud) GetChannels() ([]ChannelInfo, error) {
	channels := make([]ChannelInfo, 0)
	c.mtx.Lock()
	for _, ch := range c.channels {
		var item ChannelInfo
		item.Id = ch.channelId
		item.Name = ch.name
		channels = append(channels, item)
	}

	c.mtx.Unlock()
	return channels, nil
}

func (c *Cloud) AddChannel(channelId string, password string, name string) error {
	ch := NewChannel(c.iDataStorage, channelId, password, name)
	c.mtx.Lock()
	c.channels = append(c.channels, ch)
	c.mtx.Unlock()
	return nil
}

func (c *Cloud) EditChannel(channelId string, name string) error {
	c.mtx.Lock()
	for _, ch := range c.channels {
		if ch.channelId == channelId {
			ch.name = name
			break
		}
	}
	c.mtx.Unlock()
	return nil
}

func (c *Cloud) RemoveChannel(channelId string) error {
	c.mtx.Lock()
	for i, ch := range c.channels {
		if ch.channelId == channelId {
			ch.Stop()
			c.channels = append(c.channels[:i], c.channels[i+1:]...)
			break
		}
	}
	c.mtx.Unlock()
	return nil
}

func (c *Cloud) AddItems(channels []string, items []string) error {
	chs := make([]*Channel, 0)
	channelsAsMap := make(map[string]string)
	for _, ch := range channels {
		channelsAsMap[ch] = ch
	}
	var err error
	c.mtx.Lock()
	for _, ch := range c.channels {
		if _, ok := channelsAsMap[ch.channelId]; ok {
			chs = append(chs, ch)
		}
	}
	c.mtx.Unlock()

	for _, ch := range chs {
		err = ch.AddItems(items)
	}

	return err
}

func (c *Cloud) GetChannelsWithItem(name string) []string {
	result := make([]string, 0)
	for _, ch := range c.channels {
		for _, item := range ch.items {
			if item == name {
				result = append(result, ch.channelId)
				break
			}
		}
	}
	return result
}

func (c *Cloud) GetChannelsNamesWithItem(name string) []string {
	result := make([]string, 0)
	for _, ch := range c.channels {
		for _, item := range ch.items {
			if item == name {
				result = append(result, ch.name)
				break
			}
		}
	}
	return result
}

func (c *Cloud) RemoveItems(channels []string, items []string) error {
	var err error
	removeFromAllChannels := false
	if channels == nil {
		removeFromAllChannels = true
	}

	chs := make([]*Channel, 0)

	if !removeFromAllChannels {
		channelsAsMap := make(map[string]string)
		for _, ch := range channels {
			channelsAsMap[ch] = ch
		}
		c.mtx.Lock()
		for _, ch := range c.channels {
			if _, ok := channelsAsMap[ch.channelId]; ok {
				chs = append(chs, ch)
			}
		}
		c.mtx.Unlock()
	} else {
		c.mtx.Lock()
		for _, ch := range c.channels {
			chs = append(chs, ch)
		}
		c.mtx.Unlock()
	}

	for _, ch := range chs {
		err = ch.RemoveItems(items)
	}

	return err
}

func (c *Cloud) RemoveAllItems(channelsId string) error {
	var channel *Channel
	var err error
	c.mtx.Lock()
	for _, ch := range c.channels {
		if ch.channelId == channelsId {
			channel = ch
		}
	}
	c.mtx.Unlock()

	if channel != nil {
		err = channel.RemoveAllItems()
	} else {
		err = errors.New("no channel found")
	}
	return err
}

func (c *Cloud) GetChannelValues(channelId string) ([]common_interfaces.Item, error) {
	var channel *Channel
	var err error
	c.mtx.Lock()
	for _, ch := range c.channels {
		if ch.channelId == channelId {
			channel = ch
		}
	}
	c.mtx.Unlock()

	var values []common_interfaces.Item
	if channel != nil {
		values, err = channel.GetValues()
	} else {
		err = errors.New("no channel found")
	}
	return values, err
}

func (c *Cloud) RenameItems(oldPrefix string, newPrefix string) {
	c.mtx.Lock()
	for _, ch := range c.channels {
		ch.RenameItems(oldPrefix, newPrefix)
	}
	c.mtx.Unlock()
}
