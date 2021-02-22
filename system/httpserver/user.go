package httpserver

import (
	"encoding/json"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
)

func (c *HttpServer) SessionOpen(request []byte) (response []byte, err error) {
	var req nodeinterface.SessionOpenRequest
	var resp nodeinterface.SessionOpenResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp, err = c.system.OpenSession(req.UserName, req.Password)
	if err != nil {
		return
	}

	response, err = json.MarshalIndent(resp, "", " ")
	return
}
