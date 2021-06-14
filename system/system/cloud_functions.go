package system

import "github.com/gazercloud/gazernode/protocols/nodeinterface"

func (c *System) CloudLogin(userName string, password string) error {
	c.cloudConnection.Login(userName, password)
	return nil
}

func (c *System) CloudLogout() (nodeinterface.CloudLogoutResponse, error) {
	var result nodeinterface.CloudLogoutResponse
	c.cloudConnection.Logout()
	return result, nil
}

func (c *System) CloudState() (nodeinterface.CloudStateResponse, error) {
	return c.cloudConnection.State()
}

func (c *System) CloudNodes() (nodeinterface.CloudNodesResponse, error) {
	return c.cloudConnection.Nodes()
}

func (c *System) CloudAddNode() (nodeinterface.CloudAddNodeResponse, error) {
	return c.cloudConnection.AddNode()
}

func (c *System) CloudUpdateNode() (nodeinterface.CloudUpdateNodeResponse, error) {
	return c.cloudConnection.UpdateNode()
}

func (c *System) CloudRemoveNode() (nodeinterface.CloudRemoveNodeResponse, error) {
	return c.cloudConnection.RemoveNode()
}

func (c *System) CloudGetSettings() (nodeinterface.CloudGetSettingsResponse, error) {
	return c.cloudConnection.GetSettings()
}

func (c *System) CloudSetSettings() (nodeinterface.CloudSetSettingsResponse, error) {
	return c.cloudConnection.SetSettings()
}
