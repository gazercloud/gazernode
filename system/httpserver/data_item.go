package httpserver

import (
	"encoding/json"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
	"time"
)

func (c *HttpServer) DataItemList(request []byte) (response []byte, err error) {
	var req nodeinterface.DataItemListRequest
	var resp nodeinterface.DataItemListResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp.Items = c.system.GetItemsValues(req.Items)

	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) DataItemListAll(request []byte) (response []byte, err error) {
	var req nodeinterface.DataItemListAllRequest
	var resp nodeinterface.DataItemListAllResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp.Items = c.system.GetAllItems()

	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) DataItemWrite(request []byte) (response []byte, err error) {
	var req nodeinterface.DataItemWriteRequest
	var resp nodeinterface.DataItemWriteResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	err = c.system.SetItem(req.ItemName, req.Value, "", time.Now().UTC(), "")

	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) DataItemHistory(request []byte) (response []byte, err error) {
	var req nodeinterface.DataItemHistoryRequest
	var resp nodeinterface.DataItemHistoryResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp.History, err = c.system.ReadHistory(req.Name, req.DTBegin, req.DTEnd)

	response, err = json.MarshalIndent(resp, "", " ")
	return
}
