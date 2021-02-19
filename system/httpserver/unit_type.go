package httpserver

import (
	"encoding/json"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
)

func (c *HttpServer) UnitTypeList(request []byte) (response []byte, err error) {
	var req nodeinterface.UnitTypeListRequest
	var resp nodeinterface.UnitTypeListResponse
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

	resp.UnitName, resp.UnitConfigMeta, err = c.system.GetConfigByType(req.UnitType)

	response, err = json.MarshalIndent(resp, "", " ")
	return
}
