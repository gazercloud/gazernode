package history

import (
	"github.com/gazercloud/gazernode/common_interfaces"
)

type ReadResultItem struct {
	DatetimeFirst int64
	DatetimeLast  int64
	FirstValue    float64
	LastValue     float64
	MinValue      float64
	MaxValue      float64
	AvgValue      float64
	Qualities     []int64
	CountOfValues int64
}

type ReadResult struct {
	Id      uint64                         `json:"id"`
	DTBegin int64                          `json:"dt_begin"`
	DTEnd   int64                          `json:"dt_end"`
	Items   []*common_interfaces.ItemValue `json:"items"`
}
