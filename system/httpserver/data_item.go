package httpserver

import (
	"encoding/json"
	"errors"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/system/history"
	"github.com/gazercloud/gazernode/system/protocols/nodeinterface"
	"github.com/gazercloud/gazernode/utilities/logger"
	"github.com/gazercloud/gazernode/utilities/uom"
	"math"
	"strconv"
	"strings"
	"time"
)

func (c *HttpServer) DataItemList(request []byte) (response []byte, err error) {
	var req nodeinterface.DataItemListRequest
	var resp nodeinterface.DataItemListResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp.Items = c.system.GetItemsValues(req.Items)

	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) DataItemListAll(request []byte) (response []byte, err error) {
	var req nodeinterface.DataItemListAllRequest
	var resp nodeinterface.DataItemListAllResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp.Items = c.system.GetAllItems()

	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) DataItemWrite(request []byte) (response []byte, err error) {
	var req nodeinterface.DataItemWriteRequest
	var resp nodeinterface.DataItemWriteResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	err = c.system.SetItemByName(req.ItemName, req.Value, "", time.Now().UTC(), true)
	if err != nil {
		return
	}

	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) DataItemRemove(request []byte) (response []byte, err error) {
	var req nodeinterface.DataItemRemoveRequest
	var resp nodeinterface.DataItemRemoveResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	err = c.system.RemoveItems(req.Items)
	if err != nil {
		return
	}

	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) DataItemPropSet(request []byte) (response []byte, err error) {
	var req nodeinterface.DataItemPropSetRequest
	var resp nodeinterface.DataItemPropSetResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	err = c.system.DataItemPropSet(req.ItemName, req.Props)
	if err != nil {
		return
	}

	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) DataItemPropGet(request []byte) (response []byte, err error) {
	var req nodeinterface.DataItemPropGetRequest
	var resp nodeinterface.DataItemPropGetResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp.Props, err = c.system.DataItemPropGet(req.ItemName)
	if err != nil {
		return
	}

	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) DataItemHistory(request []byte) (response []byte, err error) {
	var req nodeinterface.DataItemHistoryRequest
	var resp nodeinterface.DataItemHistoryResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	logger.Println("HttpServer DataItemHistory", req.Name, (req.DTEnd-req.DTBegin)/1000000)

	if req.DTEnd-req.DTBegin < 1 {
		err = errors.New("wrong time range (min)")
		return
	}

	if req.DTEnd-req.DTBegin > 2*365*24*3600*1000000 {
		err = errors.New("wrong time range (max)")
		return
	}

	resp.History, err = c.system.ReadHistory(req.Name, req.DTBegin, req.DTEnd)
	if err != nil {
		return
	}

	logger.Println("HttpServer DataItemHistory", req.Name, (req.DTEnd-req.DTBegin)/1000000, len(resp.History.Items))

	response, err = json.MarshalIndent(resp, "", " ")
	return
}

func (c *HttpServer) DataItemHistoryChart(request []byte) (response []byte, err error) {
	var req nodeinterface.DataItemHistoryChartRequest
	var resp nodeinterface.DataItemHistoryChartResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp.Items = make([]*nodeinterface.DataItemHistoryChartResponseDataItem, 0)

	itemNames := make([]string, 0)
	for _, reqItem := range req.Items {
		itemNames = append(itemNames, reqItem.Name)
	}
	itemValuesList := c.system.GetItemsValues(itemNames)
	itemValuesMap := make(map[string]common_interfaces.ItemStateInfo)
	for _, itemValue := range itemValuesList {
		itemValuesMap[itemValue.Name] = itemValue
	}

	for _, reqItem := range req.Items {
		itemHistoryError := ""
		validTimeRange := true

		if reqItem.GroupTimeRange < 1 {
			itemHistoryError = "wrong group_time_range"
			validTimeRange = false
		}

		if reqItem.DTEnd-reqItem.DTBegin < 1 {
			itemHistoryError = "wrong time range (min)"
			validTimeRange = false
		}

		if reqItem.DTEnd-reqItem.DTBegin > 2*365*24*3600*1000000 {
			itemHistoryError = "wrong time range (max)"
			validTimeRange = false
		}

		if reqItem.GroupTimeRange < 1 {
			itemHistoryError = ""
			validTimeRange = false
		} else {
			expectedItemsCount := (reqItem.DTEnd - reqItem.DTBegin) / reqItem.GroupTimeRange
			if expectedItemsCount > 10000 {
				itemHistoryError = "wrong time range (max items)"
				validTimeRange = false
			}
		}

		resultItems := make([]*nodeinterface.DataItemHistoryChartResponseDataItemValue, 0)

		if validTimeRange {

			var respItems *history.ReadResult

			respItems, err = c.system.ReadHistory(reqItem.Name, reqItem.DTBegin, reqItem.DTEnd)
			if err != nil {
				return
			}

			//logger.Println("ReadHistory: ", respItems.Items)

			rawValues := make([]*common_interfaces.ItemValue, 0)
			rawValuesGroupIndex := make([]int64, 0)

			for _, item := range respItems.Items {
				rawValues = append(rawValues, item)
				groupIndex := (item.DT - reqItem.DTBegin) / reqItem.GroupTimeRange
				rawValuesGroupIndex = append(rawValuesGroupIndex, groupIndex)
			}

			lastGroupIndex := int64(-1)
			var currentValueRange *nodeinterface.DataItemHistoryChartResponseDataItemValue

			for index := range rawValuesGroupIndex {
				r := rawValues[index]
				validValue := false

				if lastGroupIndex != rawValuesGroupIndex[index] {
					if currentValueRange != nil {
						resultItems = append(resultItems, currentValueRange)
						currentValueRange = nil
					}
					lastGroupIndex = rawValuesGroupIndex[index]
				}

				if currentValueRange == nil {
					currentValueRange = &nodeinterface.DataItemHistoryChartResponseDataItemValue{}
					currentValueRange.DatetimeFirst = r.DT - (r.DT % reqItem.GroupTimeRange)
					currentValueRange.DatetimeLast = r.DT - (r.DT % reqItem.GroupTimeRange) + reqItem.GroupTimeRange - 1
					currentValueRange.Qualities = make([]int64, 0)
					currentValueRange.MinValue = math.MaxFloat64
					currentValueRange.MaxValue = -math.MaxFloat64
					currentValueRange.AvgValue = 0
					currentValueRange.FirstValue = 0
					currentValueRange.LastValue = 0
				}

				if r.UOM != uom.ERROR && r.UOM != uom.STARTED && r.UOM != uom.STOPPED {
					valueAsString := strings.Trim(r.Value, " \r\n\t")
					valueAsFloat, err := strconv.ParseFloat(valueAsString, 64)

					if r.UOM != "" {
						currentValueRange.UOM = r.UOM
					}

					if err == nil {
						validValue = true

						if valueAsFloat < currentValueRange.MinValue {
							currentValueRange.MinValue = valueAsFloat
						}
						if valueAsFloat > currentValueRange.MaxValue {
							currentValueRange.MaxValue = valueAsFloat
						}
						currentValueRange.AvgValue += valueAsFloat
						currentValueRange.SumValue += valueAsFloat
						if currentValueRange.CountOfValues > 0 {
							currentValueRange.AvgValue /= 2
						}

						if currentValueRange.CountOfValues == 0 {
							currentValueRange.FirstValue = valueAsFloat
						}

						currentValueRange.LastValue = valueAsFloat

						currentValueRange.CountOfValues++
					}
				}

				if r.UOM != "error" && validValue {
					foundGood := false
					for _, q := range currentValueRange.Qualities {
						if q == 192 {
							foundGood = true
						}
					}
					if !foundGood {
						currentValueRange.Qualities = append(currentValueRange.Qualities, 192)
						currentValueRange.HasGood = true
					}
				} else {
					foundBad := false
					for _, q := range currentValueRange.Qualities {
						if q == 0 {
							foundBad = true
						}
					}
					if !foundBad {
						currentValueRange.Qualities = append(currentValueRange.Qualities, 0)
						currentValueRange.HasBad = true
					}
				}

			}
			if currentValueRange != nil {
				resultItems = append(resultItems, currentValueRange)
				currentValueRange = nil
			}
		}

		dataItem := &nodeinterface.DataItemHistoryChartResponseDataItem{}
		dataItem.Name = reqItem.Name
		dataItem.DTBegin = reqItem.DTBegin
		dataItem.DTEnd = reqItem.DTEnd
		dataItem.GroupTimeRange = reqItem.GroupTimeRange
		dataItem.Items = resultItems
		dataItem.HistoryError = itemHistoryError

		if itemValue, ok := itemValuesMap[dataItem.Name]; ok {
			dataItem.Value = &itemValue
		} else {
			dataItem.Value = &common_interfaces.ItemStateInfo{
				Id:    0,
				Name:  dataItem.Name,
				Value: "",
				DT:    0,
				UOM:   "",
			}
		}

		resp.Items = append(resp.Items, dataItem)
	}
	response, err = json.Marshal(resp)
	return
}
