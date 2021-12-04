package gazer_client

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/history"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
	"io/fs"
	"io/ioutil"
)

func (c *GazerNodeClient) Write(itemName string, value string) error {
	var req nodeinterface.DataItemWriteRequest
	req.ItemName = itemName
	req.Value = value
	var call Call
	call.function = nodeinterface.FuncDataItemWrite
	call.request, _ = json.Marshal(req)
	call.client = c
	c.thCall(&call)
	err := call.err
	var resp nodeinterface.DataItemWriteResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return err
}

func (c *GazerNodeClient) GetItemsValues(items []string) ([]common_interfaces.ItemGetUnitItems, error) {
	var call Call
	var req nodeinterface.DataItemListRequest
	req.Items = items
	call.function = nodeinterface.FuncDataItemList
	call.request, _ = json.Marshal(req)
	call.client = c
	c.thCall(&call)
	err := call.err
	var resp nodeinterface.DataItemListResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return resp.Items, err
}

func (c *GazerNodeClient) GetAllItems() ([]common_interfaces.ItemGetUnitItems, error) {
	var call Call
	var req nodeinterface.DataItemListAllRequest
	call.function = nodeinterface.FuncDataItemListAll
	call.request, _ = json.Marshal(req)
	call.client = c
	c.thCall(&call)
	err := call.err
	var resp nodeinterface.DataItemListAllResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return resp.Items, err
}

func (c *GazerNodeClient) ReadHistory(name string, dtBegin int64, dtEnd int64) (*history.ReadResult, error) {
	var call Call
	var req nodeinterface.DataItemHistoryRequest
	req.Name = name
	req.DTBegin = dtBegin
	req.DTEnd = dtEnd
	call.function = nodeinterface.FuncDataItemHistory
	call.request, _ = json.MarshalIndent(req, "", " ")
	call.client = c
	c.thCall(&call)
	err := call.err
	var resp nodeinterface.DataItemHistoryResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return resp.History, err
}

func (c *GazerNodeClient) ReadHistoryChart(name string, dtBegin int64, dtEnd int64, groupTimeRange int64) (*nodeinterface.DataItemHistoryChartResponse, error) {
	var call Call
	var req nodeinterface.DataItemHistoryChartRequest
	req.Name = name
	req.DTBegin = dtBegin
	req.DTEnd = dtEnd
	req.GroupTimeRange = groupTimeRange
	req.OutFormat = "zip"
	call.function = nodeinterface.FuncDataItemHistoryChart
	call.request, _ = json.MarshalIndent(req, "", " ")
	call.client = c
	c.thCall(&call)
	var resp nodeinterface.DataItemHistoryChartResponse
	err := call.err
	if err == nil {
		type ZipOut struct {
			Data string `json:"data"`
		}
		var zipOut ZipOut
		err = json.Unmarshal([]byte(call.response), &zipOut)
		if err == nil {
			var data []byte
			data, err = base64.StdEncoding.DecodeString(zipOut.Data)
			if err == nil {
				buf := bytes.NewReader(data)
				var zipFile *zip.Reader
				zipFile, err = zip.NewReader(buf, buf.Size())
				var file fs.File
				file, err = zipFile.Open("data")
				if err == nil {
					var bs []byte
					bs, err = ioutil.ReadAll(file)
					if err == nil {
						err = json.Unmarshal(bs, &resp)
						if err == nil {
							//logger.Println("ok")
						}
					}
					_ = file.Close()
				}
			}
		}
	}
	return &resp, err
}

func (c *GazerNodeClient) DataItemRemove(items []string) error {
	var req nodeinterface.DataItemRemoveRequest
	req.Items = items
	var call Call
	call.function = nodeinterface.FuncDataItemRemove
	call.request, _ = json.Marshal(req)
	call.client = c
	c.thCall(&call)
	err := call.err
	var resp nodeinterface.DataItemRemoveResponse
	if err == nil {
		err = json.Unmarshal([]byte(call.response), &resp)
	}
	return err
}
