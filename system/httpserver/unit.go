package httpserver

import (
	"encoding/json"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
)

func (c *HttpServer) UnitAdd(request []byte) (response []byte, err error) {
	var req nodeinterface.UnitAddRequest
	var resp nodeinterface.UnitAddResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp.UnitId, err = c.system.AddUnit(req.UnitName, req.UnitType)

	response, err = json.MarshalIndent(resp, "", " ")
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

	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) UnitSetConfig(request []byte) (response []byte, err error) {
	var req nodeinterface.UnitSetConfigRequest
	var resp nodeinterface.UnitSetConfigResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	err = c.system.SetConfig(req.UnitId, req.UnitName, req.UnitConfig)

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

	response, err = json.MarshalIndent(resp, "", " ")
	return
}
