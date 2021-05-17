package httpserver

import (
	"encoding/json"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
)

func (c *HttpServer) CloudLogin(request []byte) (response []byte, err error) {
	var req nodeinterface.CloudLoginRequest
	var resp nodeinterface.CloudLoginResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	// TODO: logic

	if err != nil {
		return
	}
	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) CloudLogout(request []byte) (response []byte, err error) {
	var req nodeinterface.CloudLogoutRequest
	var resp nodeinterface.CloudLogoutResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	// TODO: logic

	if err != nil {
		return
	}
	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) CloudState(request []byte) (response []byte, err error) {
	var req nodeinterface.CloudStateRequest
	var resp nodeinterface.CloudStateResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	// TODO: logic

	if err != nil {
		return
	}
	response, err = json.MarshalIndent(resp, "", " ")
	return
}
