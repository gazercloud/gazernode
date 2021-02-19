package system

import (
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/system/cloud"
)

func (c *System) GetCloudChannelValues(channelId string) ([]common_interfaces.Item, error) {
	items, err := c.cloud.GetChannelValues(channelId)
	return items, err
}

func (c *System) AddCloudChannel(channelName string) error {
	err := c.cloud.AddChannel("", "", channelName)
	if err != nil {
		return err
	}
	err = c.SaveConfig()
	return err
}

func (c *System) EditCloudChannel(channelId string, channelName string) error {
	err := c.cloud.EditChannel(channelId, channelName)
	if err != nil {
		return err
	}
	err = c.SaveConfig()
	return err
}

func (c *System) RemoveCloudChannel(channelId string) error {
	err := c.cloud.RemoveChannel(channelId)
	if err != nil {
		return err
	}
	err = c.SaveConfig()
	return err
}

func (c *System) CloudAddItems(channels []string, items []string) error {
	err := c.cloud.AddItems(channels, items)
	if err != nil {
		return err
	}
	err = c.SaveConfig()
	return err
}

func (c *System) CloudRemoveItems(channels []string, items []string) error {
	err := c.cloud.RemoveItems(channels, items)
	if err != nil {
		return err
	}
	err = c.SaveConfig()
	return err
}

func (c *System) GetCloudChannels() ([]cloud.ChannelInfo, error) {
	return c.cloud.GetChannels()
}
