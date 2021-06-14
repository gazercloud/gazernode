package client

import (
	"encoding/json"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
)

func (c *Client) CloudLogin(userName string, password string, f func(error)) {
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

func (c *Client) CloudLogout(f func(error)) {
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

func (c *Client) CloudState(f func(nodeinterface.CloudStateResponse, error)) {
	var call Call
	var req nodeinterface.CloudStateResponse
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

func (c *Client) CloudNodes(f func(nodeinterface.CloudNodesResponse, error)) {
	var call Call
	var req nodeinterface.CloudNodesResponse
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

func (c *Client) CloudAddNode(f func(nodeinterface.CloudAddNodeResponse, error)) {
	var call Call
	var req nodeinterface.CloudAddNodeResponse
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

func (c *Client) CloudUpdateNode(f func(nodeinterface.CloudUpdateNodeResponse, error)) {
	var call Call
	var req nodeinterface.CloudUpdateNodeResponse
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

func (c *Client) CloudRemoveNode(f func(nodeinterface.CloudRemoveNodeResponse, error)) {
	var call Call
	var req nodeinterface.CloudRemoveNodeResponse
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

func (c *Client) CloudGetSettings(f func(nodeinterface.CloudGetSettingsResponse, error)) {
	var call Call
	var req nodeinterface.CloudGetSettingsResponse
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

func (c *Client) CloudSetSettings(f func(nodeinterface.CloudSetSettingsResponse, error)) {
	var call Call
	var req nodeinterface.CloudSetSettingsResponse
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
