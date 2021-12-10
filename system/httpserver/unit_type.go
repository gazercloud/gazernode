package httpserver

import (
	"encoding/json"
	"github.com/gazercloud/gazernode/system/protocols/nodeinterface"
)

func (c *HttpServer) UnitTypeList(request []byte) (response []byte, err error) {
	var req nodeinterface.UnitTypeListRequest
	var resp nodeinterface.UnitTypeListResponse
	req.Offset = 0
	req.Category = ""
	req.Filter = ""
	req.MaxCount = 1000
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp = c.system.UnitTypes(req.Category, req.Filter, req.Offset, req.MaxCount)

	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) UnitTypeCategories(request []byte) (response []byte, err error) {
	var req nodeinterface.UnitTypeCategoriesRequest
	var resp nodeinterface.UnitTypeCategoriesResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp = c.system.UnitCategories()

	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) UnitTypeConfigMeta(request []byte) (response []byte, err error) {
	var req nodeinterface.UnitTypeConfigMetaRequest
	var resp nodeinterface.UnitTypeConfigMetaResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp.UnitType, resp.UnitTypeConfigMeta, err = c.system.GetConfigByType(req.UnitType)
	if err != nil {
		return
	}

	response, err = json.MarshalIndent(resp, "", " ")
	return
}
