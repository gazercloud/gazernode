package gazer_client

import (
	"encoding/json"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
)

func (c *GazerNodeClient) SessionOpen(userName string, password string) (nodeinterface.SessionOpenResponse, error) {
	var call Call

	var req nodeinterface.SessionOpenRequest
	req.UserName = userName
	req.Password = password
	c.userName = userName
	call.function = nodeinterface.FuncSessionOpen
	call.request, _ = json.Marshal(req)
	call.client = c

	c.thCall(&call)

	err := call.err
	var resp nodeinterface.SessionOpenResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)

		// Save session token
		if err == nil {
			c.sessionToken = resp.SessionToken
			if c.OnSessionOpen != nil {
				c.OnSessionOpen()
			}
		}
	}
	return resp, err
}

func (c *GazerNodeClient) SessionActivate(sessionToken string, f func(nodeinterface.SessionActivateResponse, error)) {
	var call Call
	var req nodeinterface.SessionActivateRequest
	req.SessionToken = sessionToken
	call.function = nodeinterface.FuncSessionActivate
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.SessionActivateResponse
		if err == nil {
			err = json.Unmarshal([]byte(call.response), &resp)

			// Save session token
			if err == nil {
				c.sessionToken = resp.SessionToken
				if c.OnSessionOpen != nil {
					c.OnSessionOpen()
				}
			}
		}
		if f != nil {
			f(resp, err)
		}
	}
	call.client = c

	go c.thCall(&call)
}

func (c *GazerNodeClient) SessionRemove(sessionToken string, f func(nodeinterface.SessionRemoveResponse, error)) {
	var call Call
	var req nodeinterface.SessionRemoveRequest
	req.SessionToken = sessionToken
	call.function = nodeinterface.FuncSessionRemove
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.SessionRemoveResponse
		if c.sessionToken == sessionToken {
			c.sessionToken = ""
		}
		if err == nil {
			err = json.Unmarshal([]byte(call.response), &resp)
			if c.OnSessionClose != nil {
				c.OnSessionClose()
			}
		}
		if f != nil {
			f(resp, err)
		}
	}
	call.client = c

	go c.thCall(&call)
}

func (c *GazerNodeClient) SessionList(userName string, f func(nodeinterface.SessionListResponse, error)) {
	var call Call
	var req nodeinterface.SessionListRequest
	req.UserName = userName
	call.function = nodeinterface.FuncSessionList
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.SessionListResponse
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

func (c *GazerNodeClient) UserList(f func(nodeinterface.UserListResponse, error)) {
	var call Call
	var req nodeinterface.UserListRequest
	call.function = nodeinterface.FuncUserList
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.UserListResponse
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

func (c *GazerNodeClient) UserAdd(userName string, password string, f func(nodeinterface.UserAddResponse, error)) {
	var call Call
	var req nodeinterface.UserAddRequest
	req.UserName = userName
	req.Password = password
	call.function = nodeinterface.FuncUserAdd
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.UserAddResponse
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

func (c *GazerNodeClient) UserSetPassword(userName string, password string, f func(nodeinterface.UserSetPasswordResponse, error)) {
	var call Call
	var req nodeinterface.UserSetPasswordRequest
	req.UserName = userName
	req.Password = password
	call.function = nodeinterface.FuncUserSetPassword
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.UserSetPasswordResponse
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

func (c *GazerNodeClient) UserRemove(userName string, f func(nodeinterface.UserRemoveResponse, error)) {
	var call Call
	var req nodeinterface.UserRemoveRequest
	req.UserName = userName
	call.function = nodeinterface.FuncUserRemove
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.UserRemoveResponse
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
