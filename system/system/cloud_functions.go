package system

import (
	"github.com/gazercloud/gazernode/system/protocols/nodeinterface"
)

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

func (c *System) CloudAddNode(name string) (nodeinterface.CloudAddNodeResponse, error) {
	return c.cloudConnection.AddNode(name)
}

func (c *System) CloudUpdateNode(nodeId string, name string) (nodeinterface.CloudUpdateNodeResponse, error) {
	return c.cloudConnection.UpdateNode(nodeId, name)
}

func (c *System) CloudRemoveNode(nodeId string) (nodeinterface.CloudRemoveNodeResponse, error) {
	return c.cloudConnection.RemoveNode(nodeId)
}

func (c *System) CloudGetSettings(request nodeinterface.CloudGetSettingsRequest) (nodeinterface.CloudGetSettingsResponse, error) {
	return c.cloudConnection.GetSettings(request)
}

func (c *System) CloudGetSettingsProfiles(request nodeinterface.CloudGetSettingsProfilesRequest) (nodeinterface.CloudGetSettingsProfilesResponse, error) {
	return c.cloudConnection.GetSettingsProfiles(request)
}

func (c *System) CloudSetSettings(request nodeinterface.CloudSetSettingsRequest) (nodeinterface.CloudSetSettingsResponse, error) {
	return c.cloudConnection.SetSettings(request)
}

func (c *System) CloudAccountInfo(request nodeinterface.CloudAccountInfoRequest) (nodeinterface.CloudAccountInfoResponse, error) {
	return c.cloudConnection.AccountInfo(request)
}

func (c *System) CloudSetCurrentNodeId(request nodeinterface.CloudSetCurrentNodeIdRequest) (nodeinterface.CloudSetCurrentNodeIdResponse, error) {
	return c.cloudConnection.SetCurrentNodeId(request)
}
