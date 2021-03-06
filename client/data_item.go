package client

import (
	"encoding/json"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/history"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
)

func (c *Client) Write(itemName string, value string, f func(error)) {
	var req nodeinterface.DataItemWriteRequest
	req.ItemName = itemName
	req.Value = value
	var call Call
	call.function = nodeinterface.FuncDataItemWrite
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.DataItemWriteResponse
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

func (c *Client) GetItemsValues(items []string, f func([]common_interfaces.ItemGetUnitItems, error)) {
	var call Call
	var req nodeinterface.DataItemListRequest
	req.Items = items
	call.function = nodeinterface.FuncDataItemList
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.DataItemListResponse
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

func (c *Client) GetAllItems(f func([]common_interfaces.ItemGetUnitItems, error)) {
	var call Call
	var req nodeinterface.DataItemListAllRequest
	call.function = nodeinterface.FuncDataItemListAll
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.DataItemListAllResponse
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

func (c *Client) ReadHistory(name string, dtBegin int64, dtEnd int64, f func(*history.ReadResult, error)) {
	var call Call
	var req nodeinterface.DataItemHistoryRequest
	req.Name = name
	req.DTBegin = dtBegin
	req.DTEnd = dtEnd
	call.function = nodeinterface.FuncDataItemHistory
	call.request, _ = json.MarshalIndent(req, "", " ")
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.DataItemHistoryResponse
		if err == nil {
			err = json.Unmarshal([]byte(call.response), &resp)
		}
		if f != nil {
			f(resp.History, err)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *Client) ReadHistoryChart(name string, dtBegin int64, dtEnd int64, groupTimeRange int64, f func(*nodeinterface.DataItemHistoryChartResponse, error)) {
	var call Call
	var req nodeinterface.DataItemHistoryChartRequest
	req.Name = name
	req.DTBegin = dtBegin
	req.DTEnd = dtEnd
	req.GroupTimeRange = groupTimeRange
	call.function = nodeinterface.FuncDataItemHistoryChart
	call.request, _ = json.MarshalIndent(req, "", " ")
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.DataItemHistoryChartResponse
		if err == nil {
			err = json.Unmarshal([]byte(call.response), &resp)
		}
		if f != nil {
			f(&resp, err)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *Client) DataItemRemove(items []string, f func(error)) {
	var req nodeinterface.DataItemRemoveRequest
	req.Items = items
	var call Call
	call.function = nodeinterface.FuncDataItemRemove
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.DataItemRemoveResponse
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
