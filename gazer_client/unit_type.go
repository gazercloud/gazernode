package gazer_client

import (
	"encoding/json"
	nodeinterface2 "github.com/gazercloud/gazernode/system/protocols/nodeinterface"
)

func (c *GazerNodeClient) UnitTypes(category string, filter string, offset int, maxCount int) (nodeinterface2.UnitTypeListResponse, error) {
	var call Call
	var req nodeinterface2.UnitTypeListRequest
	req.Category = category
	req.Filter = filter
	req.Offset = offset
	req.MaxCount = maxCount
	call.function = nodeinterface2.FuncUnitTypeList
	call.request, _ = json.Marshal(req)
	call.client = c

	c.thCall(&call)
	err := call.err
	var result nodeinterface2.UnitTypeListResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &result)
	}
	return result, err
}

func (c *GazerNodeClient) UnitCategories() (nodeinterface2.UnitTypeCategoriesResponse, error) {
	var call Call
	var req nodeinterface2.UnitTypeCategoriesRequest

	call.function = nodeinterface2.FuncUnitTypeCategories
	call.request, _ = json.Marshal(req)
	call.client = c

	c.thCall(&call)
	err := call.err
	var result nodeinterface2.UnitTypeCategoriesResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &result)
	}
	return result, err
}

func (c *GazerNodeClient) GetUnitConfigByType(unitType string) (string, string, error) {
	var call Call
	var req nodeinterface2.UnitTypeConfigMetaRequest

	req.UnitType = unitType
	call.function = nodeinterface2.FuncUnitTypeConfigMeta
	call.request, _ = json.Marshal(req)
	call.client = c

	c.thCall(&call)
	err := call.err
	var resp nodeinterface2.UnitTypeConfigMetaResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return resp.UnitType, resp.UnitTypeConfigMeta, err
}
