package client

import (
	"encoding/json"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
)

func (c *Client) ListOfUnits(f func([]nodeinterface.UnitListResponseItem, error)) {
	var call Call
	var req nodeinterface.UnitListRequest
	call.function = nodeinterface.FuncUnitList
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.UnitListResponse
		if err == nil {
			err = json.Unmarshal([]byte(call.response), &resp)
		}
		if f != nil {
			f(resp.Items, err)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *Client) AddUnit(unitType string, unitName string, config string, f func(string, error)) {
	var call Call
	var req nodeinterface.UnitAddRequest
	req.UnitType = unitType
	req.UnitName = unitName
	req.Config = config
	call.function = nodeinterface.FuncUnitAdd
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.UnitAddResponse
		if err == nil {
			err = json.Unmarshal([]byte(call.response), &resp)
		}
		if f != nil {
			f(resp.UnitId, err)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *Client) RemoveUnit(units []string, f func(error)) {
	var call Call
	var req nodeinterface.UnitRemoveRequest
	req.Units = units
	call.function = nodeinterface.FuncUnitRemove
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.UnitRemoveResponse
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

func (c *Client) GetUnitState(unitId string, f func(nodeinterface.UnitStateResponse, error)) {
	var call Call
	var req nodeinterface.UnitStateRequest
	req.UnitId = unitId

	call.function = nodeinterface.FuncUnitState
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.UnitStateResponse
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

func (c *Client) GetUnitStateAll(f func(nodeinterface.UnitStateAllResponse, error)) {
	var call Call
	var req nodeinterface.UnitStateAllRequest

	call.function = nodeinterface.FuncUnitStateAll
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.UnitStateAllResponse
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

func (c *Client) StartUnits(ids []string, f func(error)) {
	var call Call
	var req nodeinterface.UnitStartRequest
	req.Ids = ids
	call.function = nodeinterface.FuncUnitStart
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.UnitStartResponse
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

func (c *Client) StopUnits(ids []string, f func(error)) {
	var call Call
	var req nodeinterface.UnitStopRequest
	req.Ids = ids
	call.function = nodeinterface.FuncUnitStop
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.UnitStopResponse
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

func (c *Client) GetUnitValues(unitName string, f func([]common_interfaces.ItemGetUnitItems, error)) {
	var call Call
	var req nodeinterface.UnitItemsValuesRequest
	req.UnitName = unitName
	call.function = nodeinterface.FuncUnitItemsValues
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.UnitItemsValuesResponse
		if err == nil {
			err = json.Unmarshal([]byte(call.response), &resp)
		}
		if f != nil {
			f(resp.Items, err)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *Client) SetUnitConfig(unitId string, unitName string, unitConfig string, f func(error)) {
	var call Call
	var req nodeinterface.UnitSetConfigRequest
	req.UnitId = unitId
	req.UnitName = unitName
	req.UnitConfig = unitConfig
	call.function = nodeinterface.FuncUnitSetConfig
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.UnitSetConfigResponse
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

func (c *Client) GetUnitConfig(unitId string, f func(string, string, string, string, error)) {
	var call Call
	var req nodeinterface.UnitGetConfigRequest
	req.UnitId = unitId
	call.function = nodeinterface.FuncUnitGetConfig
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.UnitGetConfigResponse
		if err == nil {
			err = json.Unmarshal([]byte(call.response), &resp)
		}
		if f != nil {
			f(resp.UnitName, resp.UnitConfig, resp.UnitConfigMeta, resp.UnitType, err)
		}
	}
	call.client = c
	go c.thCall(&call)
}
