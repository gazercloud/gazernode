package httpserver

import (
	"encoding/json"
	"github.com/gazercloud/gazernode/system/protocols/nodeinterface"
)

func (c *HttpServer) CloudLogin(request []byte) (response []byte, err error) {
	var req nodeinterface.CloudLoginRequest
	var resp nodeinterface.CloudLoginResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	err = c.system.CloudLogin(req.UserName, req.Password)

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

	resp, err = c.system.CloudLogout()

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

	resp, err = c.system.CloudState()

	if err != nil {
		return
	}
	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) CloudNodes(request []byte) (response []byte, err error) {
	var req nodeinterface.CloudNodesRequest
	var resp nodeinterface.CloudNodesResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp, err = c.system.CloudNodes()

	if err != nil {
		return
	}
	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) CloudAddNode(request []byte) (response []byte, err error) {
	var req nodeinterface.CloudAddNodeRequest
	var resp nodeinterface.CloudAddNodeResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp, err = c.system.CloudAddNode(req.Name)

	if err != nil {
		return
	}
	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) CloudUpdateNode(request []byte) (response []byte, err error) {
	var req nodeinterface.CloudUpdateNodeRequest
	var resp nodeinterface.CloudUpdateNodeResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp, err = c.system.CloudUpdateNode(req.NodeId, req.Name)

	if err != nil {
		return
	}
	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) CloudRemoveNode(request []byte) (response []byte, err error) {
	var req nodeinterface.CloudRemoveNodeRequest
	var resp nodeinterface.CloudRemoveNodeResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp, err = c.system.CloudRemoveNode(req.NodeId)

	if err != nil {
		return
	}
	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) CloudGetSettings(request []byte) (response []byte, err error) {
	var req nodeinterface.CloudGetSettingsRequest
	var resp nodeinterface.CloudGetSettingsResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp, err = c.system.CloudGetSettings(req)

	if err != nil {
		return
	}
	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) CloudGetSettingsProfiles(request []byte) (response []byte, err error) {
	var req nodeinterface.CloudGetSettingsProfilesRequest
	var resp nodeinterface.CloudGetSettingsProfilesResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp, err = c.system.CloudGetSettingsProfiles(req)

	if err != nil {
		return
	}
	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) CloudSetSettings(request []byte) (response []byte, err error) {
	var req nodeinterface.CloudSetSettingsRequest
	var resp nodeinterface.CloudSetSettingsResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp, err = c.system.CloudSetSettings(req)

	if err != nil {
		return
	}
	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) CloudAccountInfo(request []byte) (response []byte, err error) {
	var req nodeinterface.CloudAccountInfoRequest
	var resp nodeinterface.CloudAccountInfoResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp, err = c.system.CloudAccountInfo(req)

	if err != nil {
		return
	}
	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) CloudSetCurrentNodeId(request []byte) (response []byte, err error) {
	var req nodeinterface.CloudSetCurrentNodeIdRequest
	var resp nodeinterface.CloudSetCurrentNodeIdResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp, err = c.system.CloudSetCurrentNodeId(req)

	if err != nil {
		return
	}
	response, err = json.MarshalIndent(resp, "", " ")
	return
}
