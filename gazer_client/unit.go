package gazer_client

import (
	"encoding/json"
	"github.com/gazercloud/gazernode/common_interfaces"
	nodeinterface2 "github.com/gazercloud/gazernode/system/protocols/nodeinterface"
)

func (c *GazerNodeClient) ListOfUnits() ([]nodeinterface2.UnitListResponseItem, error) {
	var call Call
	var req nodeinterface2.UnitListRequest
	call.function = nodeinterface2.FuncUnitList
	call.request, _ = json.Marshal(req)
	call.client = c
	c.thCall(&call)
	err := call.err
	var resp nodeinterface2.UnitListResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return resp.Items, err
}

func (c *GazerNodeClient) AddUnit(unitType string, unitName string, config string) (nodeinterface2.UnitAddResponse, error) {
	var call Call
	var req nodeinterface2.UnitAddRequest
	req.UnitType = unitType
	req.UnitName = unitName
	req.Config = config
	call.function = nodeinterface2.FuncUnitAdd
	call.request, _ = json.Marshal(req)
	call.client = c
	c.thCall(&call)
	err := call.err
	var resp nodeinterface2.UnitAddResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return resp, err
}

func (c *GazerNodeClient) RemoveUnit(units []string) error {
	var call Call
	var req nodeinterface2.UnitRemoveRequest
	req.Units = units
	call.function = nodeinterface2.FuncUnitRemove
	call.request, _ = json.Marshal(req)
	call.client = c
	c.thCall(&call)
	err := call.err
	var resp nodeinterface2.UnitRemoveResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return err
}

func (c *GazerNodeClient) GetUnitState(unitId string) (nodeinterface2.UnitStateResponse, error) {
	var call Call
	var req nodeinterface2.UnitStateRequest
	req.UnitId = unitId

	call.function = nodeinterface2.FuncUnitState
	call.request, _ = json.Marshal(req)
	call.client = c
	c.thCall(&call)
	err := call.err
	var resp nodeinterface2.UnitStateResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return resp, err
}

func (c *GazerNodeClient) GetUnitStateAll() (nodeinterface2.UnitStateAllResponse, error) {
	var call Call
	var req nodeinterface2.UnitStateAllRequest

	call.function = nodeinterface2.FuncUnitStateAll
	call.request, _ = json.Marshal(req)
	call.client = c
	c.thCall(&call)
	err := call.err
	var resp nodeinterface2.UnitStateAllResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return resp, err
}

func (c *GazerNodeClient) StartUnits(ids []string) error {
	var call Call
	var req nodeinterface2.UnitStartRequest
	req.Ids = ids
	call.function = nodeinterface2.FuncUnitStart
	call.request, _ = json.Marshal(req)
	call.client = c

	c.thCall(&call)

	err := call.err
	var resp nodeinterface2.UnitStartResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}

	return err
}

func (c *GazerNodeClient) StopUnits(ids []string) error {
	var call Call
	var req nodeinterface2.UnitStopRequest
	req.Ids = ids
	call.function = nodeinterface2.FuncUnitStop
	call.request, _ = json.Marshal(req)
	call.client = c
	c.thCall(&call)
	err := call.err
	var resp nodeinterface2.UnitStopResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return err
}

func (c *GazerNodeClient) GetUnitValues(unitName string) ([]common_interfaces.ItemGetUnitItems, error) {
	var call Call
	var req nodeinterface2.UnitItemsValuesRequest
	req.UnitName = unitName
	call.function = nodeinterface2.FuncUnitItemsValues
	call.request, _ = json.Marshal(req)
	call.client = c
	c.thCall(&call)
	err := call.err
	var resp nodeinterface2.UnitItemsValuesResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return resp.Items, err
}

func (c *GazerNodeClient) SetUnitConfig(unitId string, unitName string, unitConfig string) error {
	var call Call
	var req nodeinterface2.UnitSetConfigRequest
	req.UnitId = unitId
	req.UnitName = unitName
	req.UnitConfig = unitConfig
	call.function = nodeinterface2.FuncUnitSetConfig
	call.request, _ = json.Marshal(req)
	call.client = c
	c.thCall(&call)
	err := call.err
	var resp nodeinterface2.UnitSetConfigResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return err
}

func (c *GazerNodeClient) GetUnitConfig(unitId string) (nodeinterface2.UnitGetConfigResponse, error) {
	var call Call
	var req nodeinterface2.UnitGetConfigRequest
	req.UnitId = unitId
	call.function = nodeinterface2.FuncUnitGetConfig
	call.request, _ = json.Marshal(req)
	call.client = c
	c.thCall(&call)
	err := call.err
	var resp nodeinterface2.UnitGetConfigResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return resp, err
}
