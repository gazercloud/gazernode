package nodeinterface

import (
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/system/history"
)

type DataItemListRequest struct {
	Items []string `json:"items"`
}

type DataItemListResponse struct {
	Items []common_interfaces.ItemStateInfo `json:"items"`
}

type DataItemListAllRequest struct {
}

type DataItemListAllResponse struct {
	Items []common_interfaces.ItemGetUnitItems `json:"items"`
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
	History *history.ReadResult `json:"history"`
}

type DataItemHistoryChartRequestItem struct {
	Name           string `json:"name"`
	DTBegin        int64  `json:"dt_begin"`
	DTEnd          int64  `json:"dt_end"`
	GroupTimeRange int64  `json:"group_time_range"`
}

type DataItemHistoryChartRequest struct {
	Items []DataItemHistoryChartRequestItem `json:"items"`
}

type DataItemHistoryChartResponseDataItemValue struct {
	DatetimeFirst int64   `json:"tf"`
	DatetimeLast  int64   `json:"tl"`
	FirstValue    float64 `json:"vf"`
	LastValue     float64 `json:"vl"`
	MinValue      float64 `json:"vd"`
	MaxValue      float64 `json:"vu"`
	AvgValue      float64 `json:"va"`
	SumValue      float64 `json:"vs"`
	CountOfValues int     `json:"c"`
	Qualities     []int64 `json:"qualities"`
	HasGood       bool    `json:"has_good"`
	HasBad        bool    `json:"has_bad"`
	UOM           string  `json:"uom"`
}

type DataItemHistoryChartResponseDataItem struct {
	Name           string                                       `json:"name"`
	DTBegin        int64                                        `json:"dt_begin"`
	DTEnd          int64                                        `json:"dt_end"`
	GroupTimeRange int64                                        `json:"group_time_range"`
	Items          []*DataItemHistoryChartResponseDataItemValue `json:"items"`
	HistoryError   string                                       `json:"history_error"`
	Value          *common_interfaces.ItemStateInfo             `json:"value"`
}

type DataItemHistoryChartResponse struct {
	Items []*DataItemHistoryChartResponseDataItem `json:"items"`
}

type DataItemRemoveRequest struct {
	Items []string `json:"items"`
}

type DataItemRemoveResponse struct {
}

type PropItem struct {
	PropName  string `json:"prop_name"`
	PropValue string `json:"prop_value"`
}

type DataItemPropSetRequest struct {
	ItemName string     `json:"item_name"`
	Props    []PropItem `json:"props"`
}

type DataItemPropSetResponse struct {
}

type DataItemPropGetRequest struct {
	ItemName string `json:"item_name"`
}

type DataItemPropGetResponse struct {
	Props []PropItem `json:"props"`
}
