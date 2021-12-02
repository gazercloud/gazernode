package gazer_client

import (
	"encoding/json"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
)

func (c *GazerNodeClient) ListOfUnits() ([]nodeinterface.UnitListResponseItem, error) {
	var call Call
	var req nodeinterface.UnitListRequest
	call.function = nodeinterface.FuncUnitList
	call.request, _ = json.Marshal(req)
	call.client = c
	c.thCall(&call)
	err := call.err
	var resp nodeinterface.UnitListResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return resp.Items, err
}

func (c *GazerNodeClient) AddUnit(unitType string, unitName string, config string) (nodeinterface.UnitAddResponse, error) {
	var call Call
	var req nodeinterface.UnitAddRequest
	req.UnitType = unitType
	req.UnitName = unitName
	req.Config = config
	call.function = nodeinterface.FuncUnitAdd
	call.request, _ = json.Marshal(req)
	call.client = c
	c.thCall(&call)
	err := call.err
	var resp nodeinterface.UnitAddResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return resp, err
}

func (c *GazerNodeClient) RemoveUnit(units []string) error {
	var call Call
	var req nodeinterface.UnitRemoveRequest
	req.Units = units
	call.function = nodeinterface.FuncUnitRemove
	call.request, _ = json.Marshal(req)
	call.client = c
	c.thCall(&call)
	err := call.err
	var resp nodeinterface.UnitRemoveResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return err
}

func (c *GazerNodeClient) GetUnitState(unitId string) (nodeinterface.UnitStateResponse, error) {
	var call Call
	var req nodeinterface.UnitStateRequest
	req.UnitId = unitId

	call.function = nodeinterface.FuncUnitState
	call.request, _ = json.Marshal(req)
	call.client = c
	c.thCall(&call)
	err := call.err
	var resp nodeinterface.UnitStateResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return resp, err
}

func (c *GazerNodeClient) GetUnitStateAll() (nodeinterface.UnitStateAllResponse, error) {
	var call Call
	var req nodeinterface.UnitStateAllRequest

	call.function = nodeinterface.FuncUnitStateAll
	call.request, _ = json.Marshal(req)
	call.client = c
	c.thCall(&call)
	err := call.err
	var resp nodeinterface.UnitStateAllResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return resp, err
}

func (c *GazerNodeClient) StartUnits(ids []string) error {
	var call Call
	var req nodeinterface.UnitStartRequest
	req.Ids = ids
	call.function = nodeinterface.FuncUnitStart
	call.request, _ = json.Marshal(req)
	call.client = c

	c.thCall(&call)

	err := call.err
	var resp nodeinterface.UnitStartResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}

	return err
}

func (c *GazerNodeClient) StopUnits(ids []string) error {
	var call Call
	var req nodeinterface.UnitStopRequest
	req.Ids = ids
	call.function = nodeinterface.FuncUnitStop
	call.request, _ = json.Marshal(req)
	call.client = c
	c.thCall(&call)
	err := call.err
	var resp nodeinterface.UnitStopResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return err
}

func (c *GazerNodeClient) GetUnitValues(unitName string) ([]common_interfaces.ItemGetUnitItems, error) {
	var call Call
	var req nodeinterface.UnitItemsValuesRequest
	req.UnitName = unitName
	call.function = nodeinterface.FuncUnitItemsValues
	call.request, _ = json.Marshal(req)
	call.client = c
	c.thCall(&call)
	err := call.err
	var resp nodeinterface.UnitItemsValuesResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return resp.Items, err
}

func (c *GazerNodeClient) SetUnitConfig(unitId string, unitName string, unitConfig string) error {
	var call Call
	var req nodeinterface.UnitSetConfigRequest
	req.UnitId = unitId
	req.UnitName = unitName
	req.UnitConfig = unitConfig
	call.function = nodeinterface.FuncUnitSetConfig
	call.request, _ = json.Marshal(req)
	call.client = c
	c.thCall(&call)
	err := call.err
	var resp nodeinterface.UnitSetConfigResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return err
}

func (c *GazerNodeClient) GetUnitConfig(unitId string) (nodeinterface.UnitGetConfigResponse, error) {
	var call Call
	var req nodeinterface.UnitGetConfigRequest
	req.UnitId = unitId
	call.function = nodeinterface.FuncUnitGetConfig
	call.request, _ = json.Marshal(req)
	call.client = c
	c.thCall(&call)
	err := call.err
	var resp nodeinterface.UnitGetConfigResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return resp, err
}
