package client

import (
	"encoding/json"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
	"github.com/gazercloud/gazernode/system/cloud"
)

func (c *Client) GetCloudChannelValues(channelId string, f func([]common_interfaces.Item, error)) {
	var call Call
	var req nodeinterface.PublicChannelItemsStateRequest
	req.ChannelId = channelId
	call.function = nodeinterface.FuncPublicChannelItemsState
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.PublicChannelItemsStateResponse
		if err == nil {
			err = json.Unmarshal([]byte(call.response), &resp)
		}
		if f != nil {
			f(resp.UnitValues, err)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *Client) AddCloudChannel(channelName string, f func(error)) {
	var call Call
	var req nodeinterface.PublicChannelAddRequest
	req.ChannelName = channelName
	call.function = nodeinterface.FuncPublicChannelAdd
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.PublicChannelAddResponse
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

func (c *Client) EditCloudChannel(channelId string, channelName string, f func(error)) {
	var call Call
	var req nodeinterface.PublicChannelSetNameRequest
	req.ChannelId = channelId
	req.ChannelName = channelName
	call.function = nodeinterface.FuncPublicChannelSetName
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.PublicChannelSetNameResponse
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

func (c *Client) RemoveCloudChannel(channelId string, f func(error)) {
	var call Call
	var req nodeinterface.PublicChannelRemoveRequest
	req.ChannelId = channelId
	call.function = nodeinterface.FuncPublicChannelRemove
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.PublicChannelRemoveResponse
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

func (c *Client) GetCloudChannels(f func([]cloud.ChannelInfo, error)) {
	var call Call
	var req nodeinterface.PublicChannelListRequest
	call.function = nodeinterface.FuncPublicChannelList
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.PublicChannelListResponse
		if err == nil {
			err = json.Unmarshal([]byte(call.response), &resp)
		}
		if f != nil {
			f(resp.Channels, err)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *Client) CloudAddItems(channels []string, items []string, f func(error)) {
	var call Call
	var req nodeinterface.PublicChannelItemAddRequest
	req.Channels = channels
	req.Items = items
	call.function = nodeinterface.FuncPublicChannelItemAdd
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.PublicChannelItemAddResponse
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

func (c *Client) CloudRemoveItems(channels []string, items []string, f func(error)) {
	var call Call
	var req nodeinterface.PublicChannelItemRemoveRequest
	req.Channels = channels
	req.Items = items
	call.function = nodeinterface.FuncPublicChannelItemRemove
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.PublicChannelItemRemoveResponse
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
