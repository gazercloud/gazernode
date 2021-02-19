package httpserver

import (
	"encoding/json"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
)

func (c *HttpServer) ResourceAdd(request []byte) (response []byte, err error) {
	var req nodeinterface.ResourceAddRequest
	var resp nodeinterface.ResourceAddResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp.Id, err = c.system.ResAdd(req.Name, req.Type, req.Content)

	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) ResourceSet(request []byte) (response []byte, err error) {
	var req nodeinterface.ResourceSetRequest
	var resp nodeinterface.ResourceSetResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	err = c.system.ResSet(req.Id, req.Thumbnail, req.Content)

	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) ResourceGet(request []byte) (response []byte, err error) {
	var req nodeinterface.ResourceGetRequest
	var resp nodeinterface.ResourceGetResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp.Item, err = c.system.ResGet(req.Id)

	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) ResourceRemove(request []byte) (response []byte, err error) {
	var req nodeinterface.ResourceRemoveRequest
	var resp nodeinterface.ResourceRemoveResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	err = c.system.ResRemove(req.Id)

	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) ResourceRename(request []byte) (response []byte, err error) {
	var req nodeinterface.ResourceRenameRequest
	var resp nodeinterface.ResourceRenameResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	err = c.system.ResRename(req.Id, req.Name)

	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) ResourceList(request []byte) (response []byte, err error) {
	var req nodeinterface.ResourceListRequest
	var resp nodeinterface.ResourceListResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp.Items = c.system.ResList(req.Type, req.Filter, req.Offset, req.MaxCount)

	response, err = json.MarshalIndent(resp, "", " ")
	return
}
