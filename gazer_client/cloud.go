package gazer_client

import (
	"encoding/json"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
)

func (c *GazerNodeClient) CloudLogin(userName string, password string, f func(error)) {
	var call Call
	var req nodeinterface.CloudLoginRequest
	req.UserName = userName
	req.Password = password
	call.function = nodeinterface.FuncCloudLogin
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.CloudLoginResponse
		if err == nil {
			err = json.Unmarshal([]byte(call.response), &resp)
		}
		if f != nil {
			f(err)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *GazerNodeClient) CloudLogout(f func(error)) {
	var call Call
	var req nodeinterface.CloudLogoutRequest
	call.function = nodeinterface.FuncCloudLogout
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.CloudLogoutResponse
		if err == nil {
			err = json.Unmarshal([]byte(call.response), &resp)
		}
		if f != nil {
			f(err)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *GazerNodeClient) CloudState(f func(nodeinterface.CloudStateResponse, error)) {
	var call Call
	var req nodeinterface.CloudStateRequest
	call.function = nodeinterface.FuncCloudState
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.CloudStateResponse
		if err == nil {
			err = json.Unmarshal([]byte(call.response), &resp)
		}
		if f != nil {
			f(resp, err)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *GazerNodeClient) CloudNodes(f func(nodeinterface.CloudNodesResponse, error)) {
	var call Call
	var req nodeinterface.CloudNodesRequest
	call.function = nodeinterface.FuncCloudNodes
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.CloudNodesResponse
		if err == nil {
			err = json.Unmarshal([]byte(call.response), &resp)
		}
		if f != nil {
			f(resp, err)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *GazerNodeClient) CloudAddNode(name string, f func(nodeinterface.CloudAddNodeResponse, error)) {
	var call Call
	var req nodeinterface.CloudAddNodeRequest
	req.Name = name

	call.function = nodeinterface.FuncCloudAddNode
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.CloudAddNodeResponse
		if err == nil {
			err = json.Unmarshal([]byte(call.response), &resp)
		}
		if f != nil {
			f(resp, err)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *GazerNodeClient) CloudUpdateNode(nodeId string, name string, f func(nodeinterface.CloudUpdateNodeResponse, error)) {
	var call Call
	var req nodeinterface.CloudUpdateNodeRequest
	req.NodeId = nodeId
	req.Name = name
	call.function = nodeinterface.FuncCloudUpdateNode
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.CloudUpdateNodeResponse
		if err == nil {
			err = json.Unmarshal([]byte(call.response), &resp)
		}
		if f != nil {
			f(resp, err)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *GazerNodeClient) CloudRemoveNode(nodeId string, f func(nodeinterface.CloudRemoveNodeResponse, error)) {
	var call Call
	var req nodeinterface.CloudRemoveNodeRequest
	req.NodeId = nodeId
	call.function = nodeinterface.FuncCloudRemoveNode
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.CloudRemoveNodeResponse
		if err == nil {
			err = json.Unmarshal([]byte(call.response), &resp)
		}
		if f != nil {
			f(resp, err)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *GazerNodeClient) CloudGetSettings(f func(nodeinterface.CloudGetSettingsResponse, error)) {
	var call Call
	var req nodeinterface.CloudGetSettingsRequest
	call.function = nodeinterface.FuncCloudGetSettings
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.CloudGetSettingsResponse
		if err == nil {
			err = json.Unmarshal([]byte(call.response), &resp)
		}
		if f != nil {
			f(resp, err)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *GazerNodeClient) CloudGetSettingsProfiles(f func(nodeinterface.CloudGetSettingsProfilesResponse, error)) {
	var call Call
	var req nodeinterface.CloudGetSettingsProfilesRequest
	call.function = nodeinterface.FuncCloudGetSettingsProfiles
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.CloudGetSettingsProfilesResponse
		if err == nil {
			err = json.Unmarshal([]byte(call.response), &resp)
		}
		if f != nil {
			f(resp, err)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *GazerNodeClient) CloudSetSettings(req nodeinterface.CloudSetSettingsRequest, f func(nodeinterface.CloudSetSettingsResponse, error)) {
	var call Call
	//var req nodeinterface.CloudSetSettingsRequest
	call.function = nodeinterface.FuncCloudSetSettings
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.CloudSetSettingsResponse
		if err == nil {
			err = json.Unmarshal([]byte(call.response), &resp)
		}
		if f != nil {
			f(resp, err)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *GazerNodeClient) CloudAccountInfo(f func(nodeinterface.CloudAccountInfoResponse, error)) {
	var call Call
	var req nodeinterface.CloudAccountInfoRequest
	call.function = nodeinterface.FuncCloudAccountInfo
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.CloudAccountInfoResponse
		if err == nil {
			err = json.Unmarshal([]byte(call.response), &resp)
		}
		if f != nil {
			f(resp, err)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *GazerNodeClient) CloudSetCurrentNodeId(nodeId string, f func(nodeinterface.CloudSetCurrentNodeIdResponse, error)) {
	var call Call
	var req nodeinterface.CloudSetCurrentNodeIdRequest
	req.NodeId = nodeId
	call.function = nodeinterface.FuncCloudSetCurrentNodeId
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.CloudSetCurrentNodeIdResponse
		if err == nil {
			err = json.Unmarshal([]byte(call.response), &resp)
		}
		if f != nil {
			f(resp, err)
		}
	}
	call.client = c
	go c.thCall(&call)
}
