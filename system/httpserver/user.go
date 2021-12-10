package httpserver

import (
	"encoding/json"
	"github.com/gazercloud/gazernode/system/protocols/nodeinterface"
)

func (c *HttpServer) SessionOpen(request []byte, host string) (response []byte, err error) {
	var req nodeinterface.SessionOpenRequest
	var resp nodeinterface.SessionOpenResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp, err = c.system.OpenSession(req.UserName, req.Password, host)
	if err != nil {
		return
	}

	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) SessionActivate(request []byte) (response []byte, err error) {
	var req nodeinterface.SessionActivateRequest
	var resp nodeinterface.SessionActivateResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	_, err = c.system.CheckSession(req.SessionToken)
	if err != nil {
		return
	}

	resp.SessionToken = req.SessionToken

	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) SessionRemove(request []byte) (response []byte, err error) {
	var req nodeinterface.SessionRemoveRequest
	var resp nodeinterface.SessionRemoveResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	err = c.system.RemoveSession(req.SessionToken)
	if err != nil {
		return
	}

	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) SessionList(request []byte) (response []byte, err error) {
	var req nodeinterface.SessionListRequest
	var resp nodeinterface.SessionListResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp, err = c.system.SessionList(req.UserName)
	if err != nil {
		return
	}

	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) UserList(request []byte) (response []byte, err error) {
	var req nodeinterface.UserListRequest
	var resp nodeinterface.UserListResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp, err = c.system.UserList()
	if err != nil {
		return
	}

	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) UserAdd(request []byte) (response []byte, err error) {
	var req nodeinterface.UserAddRequest
	var resp nodeinterface.UserAddResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp, err = c.system.UserAdd(req.UserName, req.Password)
	if err != nil {
		return
	}

	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) UserSetPassword(request []byte) (response []byte, err error) {
	var req nodeinterface.UserSetPasswordRequest
	var resp nodeinterface.UserSetPasswordResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp, err = c.system.UserSetPassword(req.UserName, req.Password)
	if err != nil {
		return
	}

	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) UserRemove(request []byte) (response []byte, err error) {
	var req nodeinterface.UserRemoveRequest
	var resp nodeinterface.UserRemoveResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp, err = c.system.UserRemove(req.UserName)
	if err != nil {
		return
	}

	response, err = json.MarshalIndent(resp, "", " ")
	return
}
