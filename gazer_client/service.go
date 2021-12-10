package gazer_client

import (
	"encoding/json"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/system/protocols/lookup"
	nodeinterface2 "github.com/gazercloud/gazernode/system/protocols/nodeinterface"
	"net/url"
)

func (c *GazerNodeClient) Lookup(entity string, f func(lookup.Result, error)) {
	var call Call
	var req nodeinterface2.ServiceLookupRequest
	req.Entity = entity
	call.function = nodeinterface2.FuncServiceLookup
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface2.ServiceLookupResponse
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

func (c *GazerNodeClient) sessionTokenUrl() *url.URL {
	var uu url.URL
	uu.Host = c.address
	uu.Scheme = "http"
	uu.Path = "/api"
	return &uu
}

func (c *GazerNodeClient) GetStatistics(f func(common_interfaces.Statistics, error)) {
	var call Call

	var req nodeinterface2.ServiceStatisticsRequest
	call.function = nodeinterface2.FuncServiceStatistics
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface2.ServiceStatisticsResponse
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

func (c *GazerNodeClient) ServiceApi(f func(nodeinterface2.ServiceApiResponse, error)) {
	var call Call

	var req nodeinterface2.ServiceApiRequest
	call.function = nodeinterface2.FuncServiceApi
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface2.ServiceApiResponse
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

func (c *GazerNodeClient) ServiceSetNodeName(nodeName string, f func(nodeinterface2.ServiceSetNodeNameResponse, error)) {
	var call Call

	var req nodeinterface2.ServiceSetNodeNameRequest
	req.Name = nodeName
	call.function = nodeinterface2.FuncServiceSetNodeName
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface2.ServiceSetNodeNameResponse
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

func (c *GazerNodeClient) ServiceNodeName(f func(nodeinterface2.ServiceNodeNameResponse, error)) {
	var call Call

	var req nodeinterface2.ServiceNodeNameRequest
	call.function = nodeinterface2.FuncServiceNodeName
	call.request, _ = json.Marshal(req)
	call.onResponse = func(call *Call) {
		err := call.err
		var resp nodeinterface2.ServiceNodeNameResponse
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
