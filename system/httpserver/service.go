package httpserver

import (
	"encoding/json"
	"github.com/gazercloud/gazernode/system/protocols/nodeinterface"
)

func (c *HttpServer) ServiceLookup(request []byte) (response []byte, err error) {
	var req nodeinterface.ServiceLookupRequest
	var resp nodeinterface.ServiceLookupResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp.Result, err = c.system.Lookup(req.Entity)
	if err != nil {
		return
	}
	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) ServiceStatistics(request []byte) (response []byte, err error) {
	var req nodeinterface.ServiceStatisticsRequest
	var resp nodeinterface.ServiceStatisticsResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp.Stat, err = c.system.GetStatistics()
	if err != nil {
		return
	}
	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) ServiceApi(request []byte) (response []byte, err error) {
	var req nodeinterface.ServiceApiRequest
	var resp nodeinterface.ServiceApiResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp, err = c.system.GetApi()
	if err != nil {
		return
	}
	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) ServiceSetNodeName(request []byte) (response []byte, err error) {
	var req nodeinterface.ServiceSetNodeNameRequest
	var resp nodeinterface.ServiceSetNodeNameResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	err = c.system.SetNodeName(req.Name)
	if err != nil {
		return
	}
	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) ServiceNodeName(request []byte) (response []byte, err error) {
	var req nodeinterface.ServiceNodeNameRequest
	var resp nodeinterface.ServiceNodeNameResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp.Name = c.system.NodeName()
	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) ServiceInfo(request []byte) (response []byte, err error) {
	var req nodeinterface.ServiceInfoRequest
	var resp nodeinterface.ServiceInfoResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp, err = c.system.GetInfo()
	if err != nil {
		return
	}
	response, err = json.MarshalIndent(resp, "", " ")
	return
}
