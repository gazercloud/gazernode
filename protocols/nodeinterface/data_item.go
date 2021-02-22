package nodeinterface

import (
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/history"
)

type DataItemListRequest struct {
	Items []string `json:"items"`
}

type DataItemListResponse struct {
	UnitValues []common_interfaces.ItemGetUnitItems
}

type DataItemListAllRequest struct {
}

type DataItemListAllResponse struct {
	UnitValues []common_interfaces.ItemGetUnitItems
}

type DataItemWriteRequest struct {
	ItemName string `json:"item_name"`
	Value    string `json:"value"`
}

type DataItemWriteResponse struct {
}

type DataItemHistoryRequest struct {
	Name    string `json:"name"`
	DTBegin int64  `json:"dt_begin"`
	DTEnd   int64  `json:"dt_end"`
}

type DataItemHistoryResponse struct {
	ReadResult *history.ReadResult
}