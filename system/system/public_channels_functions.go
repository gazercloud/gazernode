package system

import (
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/system/public_channel"
)

func (c *System) GetCloudChannelValues(channelId string) ([]common_interfaces.Item, error) {
	items, err := c.publicChannels.GetChannelValues(channelId)
	return items, err
}

func (c *System) AddCloudChannel(channelName string) error {
	err := c.publicChannels.AddChannel("", "", channelName)
	if err != nil {
		return err
	}
	err = c.SaveConfig()
	return err
}

func (c *System) EditCloudChannel(channelId string, channelName string) error {
	err := c.publicChannels.EditChannel(channelId, channelName)
	if err != nil {
		return err
	}
	err = c.SaveConfig()
	return err
}

func (c *System) RemoveCloudChannel(channelId string) error {
	err := c.publicChannels.RemoveChannel(channelId)
	if err != nil {
		return err
	}
	err = c.SaveConfig()
	return err
}

func (c *System) CloudAddItems(channels []string, items []string) error {
	err := c.publicChannels.AddItems(channels, items)
	if err != nil {
		return err
	}
	err = c.SaveConfig()
	return err
}

func (c *System) CloudRemoveItems(channels []string, items []string) error {
	err := c.publicChannels.RemoveItems(channels, items)
	if err != nil {
		return err
	}
	err = c.SaveConfig()
	return err
}

func (c *System) GetCloudChannels() ([]public_channel.ChannelInfo, error) {
	return c.publicChannels.GetChannels()
}
