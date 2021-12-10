package gazer_client

import (
	"encoding/json"
	nodeinterface2 "github.com/gazercloud/gazernode/system/protocols/nodeinterface"
)

func (c *GazerNodeClient) UnitTypes(category string, filter string, offset int, maxCount int, f func(nodeinterface2.UnitTypeListResponse, error)) {
	var call Call
	var req nodeinterface2.UnitTypeListRequest
	req.Category = category
	req.Filter = filter
	req.Offset = offset
	req.MaxCount = maxCount
	call.function = nodeinterface2.FuncUnitTypeList
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var result nodeinterface2.UnitTypeListResponse
		if err == nil {
			err = json.Unmarshal([]byte(call.response), &result)
		}
		if f != nil {
			f(result, err)
		}
	}
	call.client = c

	go c.thCall(&call)
}

func (c *GazerNodeClient) UnitCategories(f func(nodeinterface2.UnitTypeCategoriesResponse, error)) {
	var call Call
	var req nodeinterface2.UnitTypeCategoriesRequest

	call.function = nodeinterface2.FuncUnitTypeCategories
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var result nodeinterface2.UnitTypeCategoriesResponse
		if err == nil {
			err = json.Unmarshal([]byte(call.response), &result)
		}
		if f != nil {
			f(result, err)
		}
	}
	call.client = c

	go c.thCall(&call)
}

func (c *GazerNodeClient) GetUnitConfigByType(unitType string, f func(string, string, error)) {
	var call Call
	var req nodeinterface2.UnitTypeConfigMetaRequest

	req.UnitType = unitType
	call.function = nodeinterface2.FuncUnitTypeConfigMeta
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface2.UnitTypeConfigMetaResponse
		if err == nil {
			err = json.Unmarshal([]byte(call.response), &resp)
		}
		if f != nil {
			f(resp.UnitType, resp.UnitTypeConfigMeta, err)
		}
	}
	call.client = c

	go c.thCall(&call)
}
