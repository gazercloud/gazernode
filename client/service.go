package client

import (
	"encoding/json"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/protocols/lookup"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
	"net/url"
)

func (c *Client) Lookup(entity string, f func(lookup.Result, error)) {
	var call Call
	var req nodeinterface.ServiceLookupRequest
	req.Entity = entity
	call.function = nodeinterface.FuncServiceLookup
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.ServiceLookupResponse
		if err == nil {
			err = json.Unmarshal([]byte(call.response), &resp)
		}
		if f != nil {
			f(resp.Result, err)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *Client) sessionTokenUrl() *url.URL {
	var uu url.URL
	uu.Host = c.address
	uu.Scheme = "http"
	uu.Path = "/api"
	return &uu
}

func (c *Client) GetStatistics(f func(common_interfaces.Statistics, error)) {
	var call Call

	var req nodeinterface.ServiceStatisticsRequest
	call.function = nodeinterface.FuncServiceStatistics
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.ServiceStatisticsResponse
		if err == nil {
			err = json.Unmarshal([]byte(call.response), &resp)
		}
		if f != nil {
			f(resp.Stat, err)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *Client) ServiceApi(f func(nodeinterface.ServiceApiResponse, error)) {
	var call Call

	var req nodeinterface.ServiceApiRequest
	call.function = nodeinterface.FuncServiceApi
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface.ServiceApiResponse
		if err == nil {
			err = json.Unmarshal([]byte(call.response), &resp)
		}
		if f != nil {
			f(resp, err)
		}
	}
	call.client = c
	go c.thCall(&call)
}
