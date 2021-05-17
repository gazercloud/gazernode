package system

import "github.com/gazercloud/gazernode/protocols/nodeinterface"

func (c *System) CloudLogin(userName string, password string) error {
	/*err := c.publicChannels.RemoveChannel(channelId)
	if err != nil {
		return err
	}
	err = c.SaveConfig()
	return err*/
	return nil
}

func (c *System) CloudLogout() error {
	/*err := c.publicChannels.RemoveChannel(channelId)
	if err != nil {
		return err
	}
	err = c.SaveConfig()
	return err*/
	return nil
}

func (c *System) CloudState() (nodeinterface.CloudStateResponse, error) {
	var result nodeinterface.CloudStateResponse
	/*err := c.publicChannels.RemoveChannel(channelId)
	if err != nil {
		return err
	}
	err = c.SaveConfig()
	return result, err*/
	return result, nil
}
