package gazer_client

import (
	"encoding/json"
	nodeinterface2 "github.com/gazercloud/gazernode/system/protocols/nodeinterface"
)

func (c *GazerNodeClient) SessionOpen(userName string, password string) (nodeinterface2.SessionOpenResponse, error) {
	var call Call

	var req nodeinterface2.SessionOpenRequest
	req.UserName = userName
	req.Password = password
	c.userName = userName
	call.function = nodeinterface2.FuncSessionOpen
	call.request, _ = json.Marshal(req)
	call.client = c

	c.thCall(&call)

	err := call.err
	var resp nodeinterface2.SessionOpenResponse
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

func (c *GazerNodeClient) SessionActivate(sessionToken string, f func(nodeinterface2.SessionActivateResponse, error)) {
	var call Call
	var req nodeinterface2.SessionActivateRequest
	req.SessionToken = sessionToken
	call.function = nodeinterface2.FuncSessionActivate
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface2.SessionActivateResponse
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

func (c *GazerNodeClient) SessionRemove(sessionToken string, f func(nodeinterface2.SessionRemoveResponse, error)) {
	var call Call
	var req nodeinterface2.SessionRemoveRequest
	req.SessionToken = sessionToken
	call.function = nodeinterface2.FuncSessionRemove
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface2.SessionRemoveResponse
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

func (c *GazerNodeClient) SessionList(userName string, f func(nodeinterface2.SessionListResponse, error)) {
	var call Call
	var req nodeinterface2.SessionListRequest
	req.UserName = userName
	call.function = nodeinterface2.FuncSessionList
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface2.SessionListResponse
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

func (c *GazerNodeClient) UserList(f func(nodeinterface2.UserListResponse, error)) {
	var call Call
	var req nodeinterface2.UserListRequest
	call.function = nodeinterface2.FuncUserList
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface2.UserListResponse
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

func (c *GazerNodeClient) UserAdd(userName string, password string, f func(nodeinterface2.UserAddResponse, error)) {
	var call Call
	var req nodeinterface2.UserAddRequest
	req.UserName = userName
	req.Password = password
	call.function = nodeinterface2.FuncUserAdd
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface2.UserAddResponse
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

func (c *GazerNodeClient) UserSetPassword(userName string, password string, f func(nodeinterface2.UserSetPasswordResponse, error)) {
	var call Call
	var req nodeinterface2.UserSetPasswordRequest
	req.UserName = userName
	req.Password = password
	call.function = nodeinterface2.FuncUserSetPassword
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface2.UserSetPasswordResponse
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

func (c *GazerNodeClient) UserRemove(userName string, f func(nodeinterface2.UserRemoveResponse, error)) {
	var call Call
	var req nodeinterface2.UserRemoveRequest
	req.UserName = userName
	call.function = nodeinterface2.FuncUserRemove
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface2.UserRemoveResponse
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
