package nodeinterface

import (
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/system/protocols/lookup"
)

type ServiceLookupRequest struct {
	Entity     string `json:"entity"`
	Parameters string `json:"parameters"`
}

type ServiceLookupResponse struct {
	Result lookup.Result `json:"result"`
}

type ServiceStatisticsRequest struct {
}

type ServiceStatisticsResponse struct {
	Stat common_interfaces.Statistics `json:"stat"`
}

type ServiceApiRequest struct {
}

type ServiceApiResponse struct {
	Product            string   `json:"product"`
	Version            string   `json:"version"`
	BuildTime          string   `json:"build_time"`
	SupportedFunctions []string `json:"supported_functions"`
}

type ServiceSetNodeNameRequest struct {
	Name string `json:"name"`
}

type ServiceSetNodeNameResponse struct {
}

type ServiceNodeNameRequest struct {
}

type ServiceNodeNameResponse struct {
	Name string `json:"name"`
}

type ServiceInfoRequest struct {
}

type ServiceInfoResponse struct {
	NodeName  string `json:"node_name"`
	Version   string `json:"version"`
	BuildTime string `json:"build_time"`
}
