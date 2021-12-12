package gazer_client

import (
	"encoding/json"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/system/protocols/lookup"
	nodeinterface2 "github.com/gazercloud/gazernode/system/protocols/nodeinterface"
	"net/url"
)

func (c *GazerNodeClient) Lookup(entity string) (lookup.Result, error) {
	var call Call
	var req nodeinterface2.ServiceLookupRequest
	req.Entity = entity
	call.function = nodeinterface2.FuncServiceLookup
	call.request, _ = json.Marshal(req)
	call.client = c
	c.thCall(&call)
	err := call.err
	var resp nodeinterface2.ServiceLookupResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return resp.Result, err
}

func (c *GazerNodeClient) sessionTokenUrl() *url.URL {
	var uu url.URL
	uu.Host = c.address
	uu.Scheme = "http"
	uu.Path = "/api"
	return &uu
}

func (c *GazerNodeClient) GetStatistics() (common_interfaces.Statistics, error) {
	var call Call

	var req nodeinterface2.ServiceStatisticsRequest
	call.function = nodeinterface2.FuncServiceStatistics
	call.request, _ = json.Marshal(req)
	call.client = c
	c.thCall(&call)
	err := call.err
	var resp nodeinterface2.ServiceStatisticsResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return resp.Stat, err
}

func (c *GazerNodeClient) ServiceApi() (nodeinterface2.ServiceApiResponse, error) {
	var call Call

	var req nodeinterface2.ServiceApiRequest
	call.function = nodeinterface2.FuncServiceApi
	call.request, _ = json.Marshal(req)
	call.client = c
	c.thCall(&call)
	err := call.err
	var resp nodeinterface2.ServiceApiResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return resp, err
}

func (c *GazerNodeClient) ServiceSetNodeName(nodeName string) (nodeinterface2.ServiceSetNodeNameResponse, error) {
	var call Call

	var req nodeinterface2.ServiceSetNodeNameRequest
	req.Name = nodeName
	call.function = nodeinterface2.FuncServiceSetNodeName
	call.request, _ = json.Marshal(req)
	call.client = c
	c.thCall(&call)
	err := call.err
	var resp nodeinterface2.ServiceSetNodeNameResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return resp, err
}

func (c *GazerNodeClient) ServiceNodeName() (nodeinterface2.ServiceNodeNameResponse, error) {
	var call Call

	var req nodeinterface2.ServiceNodeNameRequest
	call.function = nodeinterface2.FuncServiceNodeName
	call.request, _ = json.Marshal(req)
	call.client = c
	c.thCall(&call)
	err := call.err
	var resp nodeinterface2.ServiceNodeNameResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return resp, err
}
