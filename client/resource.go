package client

import (
	"encoding/json"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
)

func (c *Client) ResAdd(name string, tp string, content []byte, f func(string, error)) {
	var call Call
	var req nodeinterface.ResourceAddRequest
	req.Name = name
	req.Type = tp
	req.Content = content
	call.function = nodeinterface.FuncResourceAdd
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.ResourceAddResponse
		if err == nil {
			err = json.Unmarshal([]byte(call.response), &resp)
		}
		if f != nil {
			f(resp.Id, err)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *Client) ResSet(id string, thumbnail []byte, content []byte, f func(error)) {
	var call Call
	var req nodeinterface.ResourceSetRequest
	req.Id = id
	req.Thumbnail = thumbnail
	req.Content = content

	call.function = nodeinterface.FuncResourceSet
	call.request, _ = json.MarshalIndent(req, "", " ")
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.ResourceSetResponse
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

func (c *Client) ResGet(id string, f func(*common_interfaces.ResourcesItem, error)) {
	var call Call
	var req nodeinterface.ResourceGetRequest
	req.Id = id
	call.function = nodeinterface.FuncResourceGet
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.ResourceGetResponse
		if err == nil {
			err = json.Unmarshal([]byte(call.response), &resp)
		}
		if f != nil {
			f(resp.Item, err)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *Client) ResList(tp string, filter string, offset int, maxCount int, f func(common_interfaces.ResourcesInfo, error)) {
	var call Call
	var req nodeinterface.ResourceListRequest
	req.Type = tp
	req.Filter = filter
	req.Offset = offset
	req.MaxCount = maxCount
	call.function = nodeinterface.FuncResourceList
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.ResourceListResponse
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

func (c *Client) ResRemove(id string, f func(error)) {
	var call Call
	var req nodeinterface.ResourceRemoveRequest
	req.Id = id

	call.function = nodeinterface.FuncResourceRemove
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.ResourceRemoveResponse
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

func (c *Client) ResRename(id string, name string, f func(error)) {
	var call Call
	var req nodeinterface.ResourceRenameRequest
	req.Id = id
	req.Name = name
	call.function = nodeinterface.FuncResourceRename
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.ResourceRenameResponse
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
