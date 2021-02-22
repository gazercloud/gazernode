package client

import (
	"encoding/json"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
)

func (c *Client) SessionOpen(userName string, password string, f func(nodeinterface.SessionOpenResponse, error)) {
	var call Call
	var req nodeinterface.SessionOpenRequest
	req.UserName = userName
	req.Password = password
	call.function = nodeinterface.FuncSessionOpen
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
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
		if f != nil {
			f(resp, err)
		}
	}
	call.client = c

	go c.thCall(&call)
}
