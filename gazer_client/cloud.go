package gazer_client

import (
	"encoding/json"
	nodeinterface2 "github.com/gazercloud/gazernode/system/protocols/nodeinterface"
)

func (c *GazerNodeClient) CloudLogin(userName string, password string) error {
	var call Call
	var req nodeinterface2.CloudLoginRequest
	req.UserName = userName
	req.Password = password
	call.function = nodeinterface2.FuncCloudLogin
	call.request, _ = json.Marshal(req)
	call.client = c
	c.thCall(&call)
	err := call.err
	var resp nodeinterface2.CloudLoginResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return err
}

func (c *GazerNodeClient) CloudLogout() error {
	var call Call
	var req nodeinterface2.CloudLogoutRequest
	call.function = nodeinterface2.FuncCloudLogout
	call.request, _ = json.Marshal(req)
	call.client = c
	c.thCall(&call)
	err := call.err
	var resp nodeinterface2.CloudLogoutResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return err
}

func (c *GazerNodeClient) CloudState() (nodeinterface2.CloudStateResponse, error) {
	var call Call
	var req nodeinterface2.CloudStateRequest
	call.function = nodeinterface2.FuncCloudState
	call.request, _ = json.Marshal(req)
	call.client = c
	c.thCall(&call)
	err := call.err
	var resp nodeinterface2.CloudStateResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return resp, err
}

func (c *GazerNodeClient) CloudNodes() (nodeinterface2.CloudNodesResponse, error) {
	var call Call
	var req nodeinterface2.CloudNodesRequest
	call.function = nodeinterface2.FuncCloudNodes
	call.request, _ = json.Marshal(req)
	call.client = c
	c.thCall(&call)
	err := call.err
	var resp nodeinterface2.CloudNodesResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return resp, err
}

func (c *GazerNodeClient) CloudAddNode(name string) (nodeinterface2.CloudAddNodeResponse, error) {
	var call Call
	var req nodeinterface2.CloudAddNodeRequest
	req.Name = name

	call.function = nodeinterface2.FuncCloudAddNode
	call.request, _ = json.Marshal(req)
	call.client = c
	c.thCall(&call)
	err := call.err
	var resp nodeinterface2.CloudAddNodeResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return resp, err
}

func (c *GazerNodeClient) CloudUpdateNode(nodeId string, name string) (nodeinterface2.CloudUpdateNodeResponse, error) {
	var call Call
	var req nodeinterface2.CloudUpdateNodeRequest
	req.NodeId = nodeId
	req.Name = name
	call.function = nodeinterface2.FuncCloudUpdateNode
	call.request, _ = json.Marshal(req)
	call.client = c
	c.thCall(&call)
	err := call.err
	var resp nodeinterface2.CloudUpdateNodeResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return resp, err
}

func (c *GazerNodeClient) CloudRemoveNode(nodeId string) (nodeinterface2.CloudRemoveNodeResponse, error) {
	var call Call
	var req nodeinterface2.CloudRemoveNodeRequest
	req.NodeId = nodeId
	call.function = nodeinterface2.FuncCloudRemoveNode
	call.request, _ = json.Marshal(req)
	call.client = c
	c.thCall(&call)
	err := call.err
	var resp nodeinterface2.CloudRemoveNodeResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return resp, err
}

func (c *GazerNodeClient) CloudGetSettings() (nodeinterface2.CloudGetSettingsResponse, error) {
	var call Call
	var req nodeinterface2.CloudGetSettingsRequest
	call.function = nodeinterface2.FuncCloudGetSettings
	call.request, _ = json.Marshal(req)
	call.client = c
	c.thCall(&call)
	err := call.err
	var resp nodeinterface2.CloudGetSettingsResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return resp, err
}

func (c *GazerNodeClient) CloudGetSettingsProfiles() (nodeinterface2.CloudGetSettingsProfilesResponse, error) {
	var call Call
	var req nodeinterface2.CloudGetSettingsProfilesRequest
	call.function = nodeinterface2.FuncCloudGetSettingsProfiles
	call.request, _ = json.Marshal(req)
	call.client = c
	c.thCall(&call)
	err := call.err
	var resp nodeinterface2.CloudGetSettingsProfilesResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return resp, err
}

func (c *GazerNodeClient) CloudSetSettings(req nodeinterface2.CloudSetSettingsRequest) (nodeinterface2.CloudSetSettingsResponse, error) {
	var call Call
	//var req nodeinterface.CloudSetSettingsRequest
	call.function = nodeinterface2.FuncCloudSetSettings
	call.request, _ = json.Marshal(req)
	call.client = c
	c.thCall(&call)
	err := call.err
	var resp nodeinterface2.CloudSetSettingsResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return resp, err
}

func (c *GazerNodeClient) CloudAccountInfo() (nodeinterface2.CloudAccountInfoResponse, error) {
	var call Call
	var req nodeinterface2.CloudAccountInfoRequest
	call.function = nodeinterface2.FuncCloudAccountInfo
	call.request, _ = json.Marshal(req)
	call.client = c
	c.thCall(&call)
	err := call.err
	var resp nodeinterface2.CloudAccountInfoResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return resp, err
}

func (c *GazerNodeClient) CloudSetCurrentNodeId(nodeId string) (nodeinterface2.CloudSetCurrentNodeIdResponse, error) {
	var call Call
	var req nodeinterface2.CloudSetCurrentNodeIdRequest
	req.NodeId = nodeId
	call.function = nodeinterface2.FuncCloudSetCurrentNodeId
	call.request, _ = json.Marshal(req)
	call.client = c
	c.thCall(&call)
	err := call.err
	var resp nodeinterface2.CloudSetCurrentNodeIdResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return resp, err
}
