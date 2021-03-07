package nodeinterface

import (
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/protocols/lookup"
)

type ServiceLookupRequest struct {
	Entity     string `json:"entity"`
	Parameters string `json:"parameters"`
}

type ServiceLookupResponse struct {
	Result lookup.Result
}

type ServiceStatisticsRequest struct {
}

type ServiceStatisticsResponse struct {
	Stat common_interfaces.Statistics
}

type ServiceApiRequest struct {
}

type ServiceApiResponse struct {
	SupportedFunctions []string `json:"supported_functions"`
}
