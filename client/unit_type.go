package client

import (
	"encoding/json"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
)

func (c *Client) UnitTypes(category string, filter string, offset int, maxCount int, f func(nodeinterface.UnitTypeListResponse, error)) {
	var call Call
	var req nodeinterface.UnitTypeListRequest
	req.Category = category
	req.Filter = filter
	req.Offset = offset
	req.MaxCount = maxCount
	call.function = nodeinterface.FuncUnitTypeList
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var result nodeinterface.UnitTypeListResponse
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

func (c *Client) UnitCategories(f func(nodeinterface.UnitTypeCategoriesResponse, error)) {
	var call Call
	var req nodeinterface.UnitTypeCategoriesRequest

	call.function = nodeinterface.FuncUnitTypeCategories
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var result nodeinterface.UnitTypeCategoriesResponse
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

func (c *Client) GetUnitConfigByType(unitType string, f func(string, string, error)) {
	var call Call
	var req nodeinterface.UnitTypeConfigMetaRequest

	req.UnitType = unitType
	call.function = nodeinterface.FuncUnitTypeConfigMeta
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.UnitTypeConfigMetaResponse
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
