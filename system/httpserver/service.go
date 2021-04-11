package httpserver

import (
	"encoding/json"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
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
