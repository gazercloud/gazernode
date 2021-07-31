package httpserver

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/history"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
	"io"
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

	err = c.system.SetItem(req.ItemName, req.Value, "", time.Now().UTC(), "")
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

func (c *HttpServer) DataItemHistory(request []byte) (response []byte, err error) {
	var req nodeinterface.DataItemHistoryRequest
	var resp nodeinterface.DataItemHistoryResponse
	err = json.Unmarshal(request, &req)
	if err != nil {
		return
	}

	resp.History, err = c.system.ReadHistory(req.Name, req.DTBegin, req.DTEnd)
	if err != nil {
		return
	}

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

	var respItems *history.ReadResult

	respItems, err = c.system.ReadHistory(req.Name, req.DTBegin, req.DTEnd)
	if err != nil {
		return
	}

	resultItems := make([]*nodeinterface.DataItemHistoryChartResponseItem, 0)
	rawValues := make([]*common_interfaces.ItemValue, 0)
	rawValuesGroupIndex := make([]int64, 0)

	for _, item := range respItems.Items {
		rawValues = append(rawValues, item)
		groupIndex := (item.DT - req.DTBegin) / req.GroupTimeRange
		rawValuesGroupIndex = append(rawValuesGroupIndex, groupIndex)
	}

	lastGroupIndex := int64(-1)
	var currentValueRange *nodeinterface.DataItemHistoryChartResponseItem

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
			currentValueRange = &nodeinterface.DataItemHistoryChartResponseItem{}
			currentValueRange.DatetimeFirst = r.DT - (r.DT % req.GroupTimeRange)
			currentValueRange.DatetimeLast = r.DT - (r.DT % req.GroupTimeRange) + req.GroupTimeRange - 1
			currentValueRange.Qualities = make([]int64, 0)
			currentValueRange.MinValue = 1000000000000
			currentValueRange.MaxValue = -1000000000000
			currentValueRange.AvgValue = 0
			currentValueRange.FirstValue = 0
			currentValueRange.LastValue = 0
		}

		if r.UOM != "error" {
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
	resp.Items = resultItems

	response, err = json.Marshal(resp)

	if req.OutFormat == "zip" {
		buf := new(bytes.Buffer)
		zipWriter := zip.NewWriter(buf)

		// Add some files to the archive.
		var zipFile io.Writer
		zipFile, err = zipWriter.Create("data")
		if err == nil {
			_, err = zipFile.Write([]byte(response))
		}

		// Make sure to check the error on Close.
		err = zipWriter.Close()
		type ZipOut struct {
			Data string `json:"data"`
		}

		sEnc := base64.StdEncoding.EncodeToString([]byte(buf.Bytes()))

		var zipData ZipOut
		zipData.Data = sEnc
		response, err = json.Marshal(zipData)
	}

	/*logger.Println("DataItemHistoryChart REQUEST dt(sec):", (req.DTEnd - req.DTBegin) / 1000000, "range:", req.GroupTimeRange,
	"itemsCount:", len(respItems.Items), "resCount", len(resultItems), "bytes:", len(response))*/

	return
}
