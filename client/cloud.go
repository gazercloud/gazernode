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
