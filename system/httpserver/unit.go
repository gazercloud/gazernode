package httpserver

import (
	"encoding/json"
	"github.com/gazercloud/gazernode/system/protocols/nodeinterface"
)

func (c *HttpServer) UnitAdd(request []byte, fromCloud bool) (response []byte, err error) {
	var req nodeinterface.UnitAddRequest
	var resp nodeinterface.UnitAddResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp.UnitId, err = c.system.AddUnit(req.UnitName, req.UnitType, req.Config, fromCloud)
	if err == nil {
		response, err = json.MarshalIndent(resp, "", " ")
	}
	return
}

func (c *HttpServer) UnitRemove(request []byte) (response []byte, err error) {
	var req nodeinterface.UnitRemoveRequest
	var resp nodeinterface.UnitRemoveResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	err = c.system.RemoveUnits(req.Units)
	if err != nil {
		return
	}
	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) UnitState(request []byte) (response []byte, err error) {
	var req nodeinterface.UnitStateRequest
	var resp nodeinterface.UnitStateResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp, err = c.system.GetUnitState(req.UnitId)
	if err != nil {
		return
	}

	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) UnitStateAll(request []byte) (response []byte, err error) {
	var req nodeinterface.UnitStateAllRequest
	var resp nodeinterface.UnitStateAllResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp, err = c.system.GetUnitStateAll()
	if err != nil {
		return
	}

	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) UnitItemsValues(request []byte) (response []byte, err error) {
	var req nodeinterface.UnitItemsValuesRequest
	var resp nodeinterface.UnitItemsValuesResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp.Items = c.system.GetUnitValues(req.UnitName)

	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) UnitList(request []byte) (response []byte, err error) {
	var req nodeinterface.UnitListRequest
	var resp nodeinterface.UnitListResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp = c.system.ListOfUnits()

	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) UnitStart(request []byte) (response []byte, err error) {
	var req nodeinterface.UnitStartRequest
	var resp nodeinterface.UnitStartResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	err = c.system.StartUnits(req.Ids)
	if err != nil {
		return
	}

	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) UnitStop(request []byte) (response []byte, err error) {
	var req nodeinterface.UnitStopRequest
	var resp nodeinterface.UnitStopResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	err = c.system.StopUnits(req.Ids)
	if err != nil {
		return
	}

	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) UnitSetConfig(request []byte, fromCloud bool) (response []byte, err error) {
	var req nodeinterface.UnitSetConfigRequest
	var resp nodeinterface.UnitSetConfigResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	err = c.system.SetConfig(req.UnitId, req.UnitName, req.UnitConfig, fromCloud)
	if err != nil {
		return
	}

	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) UnitGetConfig(request []byte) (response []byte, err error) {
	var req nodeinterface.UnitGetConfigRequest
	var resp nodeinterface.UnitGetConfigResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp.UnitName, resp.UnitConfig, resp.UnitConfigMeta, resp.UnitType, err = c.system.GetConfig(req.UnitId)
	resp.UnitId = req.UnitId
	if err != nil {
		return
	}

	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) UnitPropSet(request []byte) (response []byte, err error) {
	var req nodeinterface.UnitPropSetRequest
	var resp nodeinterface.UnitPropSetResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	err = c.system.UnitPropSet(req.UnitId, req.Props)
	if err != nil {
		return
	}

	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) UnitPropGet(request []byte) (response []byte, err error) {
	var req nodeinterface.UnitPropGetRequest
	var resp nodeinterface.UnitPropGetResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp.Props, err = c.system.UnitPropGet(req.UnitId)
	if err != nil {
		return
	}

	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) ResourcePropSet(request []byte) (response []byte, err error) {
	var req nodeinterface.ResourcePropSetRequest
	var resp nodeinterface.ResourcePropSetResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	err = c.system.ResourcePropSet(req.Id, req.Props)
	if err != nil {
		return
	}

	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) ResourcePropGet(request []byte) (response []byte, err error) {
	var req nodeinterface.ResourcePropGetRequest
	var resp nodeinterface.ResourcePropGetResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp.Props, err = c.system.ResourcePropGet(req.Id)
	if err != nil {
		return
	}

	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) UserPropSet(request []byte) (response []byte, err error) {
	var req nodeinterface.UserPropSetRequest
	var resp nodeinterface.UserPropSetResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	err = c.system.UserPropSet(req.UserName, req.Props)
	if err != nil {
		return
	}

	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) UserPropGet(request []byte) (response []byte, err error) {
	var req nodeinterface.UserPropGetRequest
	var resp nodeinterface.UserPropGetResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp.Props, err = c.system.UserPropGet(req.UserName)
	if err != nil {
		return
	}

	response, err = json.MarshalIndent(resp, "", " ")
	return
}
