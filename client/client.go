package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/history"
	"github.com/gazercloud/gazernode/logger"
	"github.com/gazercloud/gazernode/system/cloud"
	"github.com/gazercloud/gazernode/system/units/units_common"
	"github.com/gazercloud/gazerui/uievents"
	"github.com/gazercloud/gazerui/uiinterfaces"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Client struct {
	window   uiinterfaces.Window
	received []*Call
	mtx      sync.Mutex
	tm       *uievents.FormTimer
	client   *http.Client
	watcher  *ItemsWatcher
}

type Call struct {
	function   string
	request    []byte
	response   string
	onResponse func(call *Call)
	err        error
	client     *Client
}

func New(window uiinterfaces.Window) *Client {
	var c Client
	pc := &c

	tr := &http.Transport{}
	c.client = &http.Client{Transport: tr}
	c.client.Timeout = 3 * time.Second

	c.tm = window.NewTimer(100, pc.timer)
	c.tm.StartTimer()

	c.watcher = NewItemsWatcher(&c)

	return pc
}

func (c *Client) timer() {
	c.mtx.Lock()
	for _, call := range c.received {
		call.onResponse(call)
	}
	c.received = make([]*Call, 0)
	c.mtx.Unlock()
}

func (c *Client) GetItemValue(name string) common_interfaces.ItemValue {
	return c.watcher.Get(name)
}

func (c *Client) thCall(call *Call) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	{
		fw, _ := writer.CreateFormField("func")
		fw.Write([]byte(call.function))
	}
	{
		fw, _ := writer.CreateFormField("rj")
		if call.request == nil {
			fw.Write(make([]byte, 0))
		} else {
			fw.Write(call.request)
		}

	}
	writer.Close()

	AddStatSent(body.Len())

	addr := "127.0.0.1"
	if false {
		addr = "192.168.24.233"
	}

	response, err := c.client.Post("http://"+addr+":8084/api/request", writer.FormDataContentType(), &body)
	if err != nil {
		call.err = errors.New("no connection to local service")
		logger.Println(err)
	} else {
		content, _ := ioutil.ReadAll(response.Body)
		call.response = strings.TrimSpace(string(content))
		AddStatReceived(len(call.response))
		response.Body.Close()

		type ErrorContainer struct {
			Error string `json:"error"`
		}
		var errCont ErrorContainer
		json.Unmarshal([]byte(call.response), &errCont)
		if len(errCont.Error) > 0 {
			call.err = errors.New(errCont.Error)
		}
	}

	//client.CloseIdleConnections()

	call.client.mtx.Lock()
	call.client.received = append(call.client.received, call)
	call.client.mtx.Unlock()
}

func (c *Client) UnitTypes(category string, filter string, offset int, maxCount int, f func(common_interfaces.UnitTypes, error)) {
	var call Call
	type Request struct {
		Category string `json:"category"`
		Filter   string `json:"filter"`
		Offset   int    `json:"offset"`
		MaxCount int    `json:"max_count"`
	}
	var req Request
	req.Category = category
	req.Filter = filter
	req.Offset = offset
	req.MaxCount = maxCount

	call.function = "unit_types"
	call.request, _ = json.MarshalIndent(req, "", " ")
	call.onResponse = func(call *Call) {
		var result common_interfaces.UnitTypes
		result.Types = make([]common_interfaces.UnitTypeInfo, 0)
		err := json.Unmarshal([]byte(call.response), &result)
		if f != nil {
			f(result, err)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *Client) UnitCategories(f func([]common_interfaces.UnitCategoryInfo, error)) {
	var call Call
	type Request struct {
	}
	var req Request

	call.function = "unit_categories"
	call.request, _ = json.MarshalIndent(req, "", " ")
	call.onResponse = func(call *Call) {
		var result []common_interfaces.UnitCategoryInfo
		result = make([]common_interfaces.UnitCategoryInfo, 0)
		err := json.Unmarshal([]byte(call.response), &result)
		if f != nil {
			f(result, err)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *Client) ListOfUnits(f func([]units_common.UnitInfo, error)) {
	var call Call
	call.function = "units"
	call.request = nil
	call.onResponse = func(call *Call) {
		var result []units_common.UnitInfo
		result = make([]units_common.UnitInfo, 0)
		err := json.Unmarshal([]byte(call.response), &result)
		if f != nil {
			f(result, err)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *Client) Lookup(entity string, f func(*units_common.LookupResult, error)) {

	type Request struct {
		Entity     string `json:"entity"`
		Parameters string `json:"parameters"`
	}
	var req Request
	req.Entity = entity

	var call Call
	call.function = "lookup"
	call.request, _ = json.MarshalIndent(req, "", " ")
	call.onResponse = func(call *Call) {
		result := units_common.NewLookupResult()
		err := json.Unmarshal([]byte(call.response), result)
		if f != nil {
			f(result, err)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *Client) Write(itemName string, value string, f func(error)) {

	type Request struct {
		ItemName string `json:"item_name"`
		Value    string `json:"value"`
	}
	var req Request
	req.ItemName = itemName
	req.Value = value

	var call Call
	call.function = "write"
	call.request, _ = json.MarshalIndent(req, "", " ")
	call.onResponse = func(call *Call) {
		if f != nil {
			f(call.err)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *Client) AddUnit(unitType string, unitName string, f func(string, error)) {
	var call Call
	type Request struct {
		UnitType string `json:"type"`
		UnitName string `json:"name"`
	}
	var req Request
	req.UnitType = unitType
	req.UnitName = unitName

	type Response struct {
		Error  string `json:"error"`
		UnitId string `json:"unit_id"`
	}

	call.function = "add_unit"
	call.request, _ = json.MarshalIndent(req, "", " ")
	call.onResponse = func(call *Call) {
		var result Response
		err := json.Unmarshal([]byte(call.response), &result)
		if f != nil && err == nil {
			if result.Error != "" {
				err = errors.New(result.Error)
			}
			f(result.UnitId, err)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *Client) RemoveUnit(units []string, f func(error)) {
	var call Call
	type Request struct {
		Units []string `json:"ids"`
	}
	var req Request
	req.Units = units
	call.function = "remove_unit"
	call.request, _ = json.MarshalIndent(req, "", " ")
	call.onResponse = func(call *Call) {
		if f != nil {
			f(nil)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *Client) GetUnitState(unitId string, f func(common_interfaces.UnitState, error)) {
	var call Call
	type Request struct {
		UnitId string `json:"id"`
	}
	var req Request
	req.UnitId = unitId

	call.function = "unit_state"
	call.request, _ = json.MarshalIndent(req, "", " ")
	call.onResponse = func(call *Call) {
		var result common_interfaces.UnitState
		err := json.Unmarshal([]byte(call.response), &result)
		if f != nil {
			f(result, err)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *Client) StartUnits(ids []string, f func(error)) {
	var call Call
	type Request struct {
		Ids []string `json:"ids"`
	}
	var req Request
	req.Ids = ids
	call.function = "start_units"
	call.request, _ = json.MarshalIndent(req, "", " ")
	call.onResponse = func(call *Call) {
		if f != nil {
			f(nil)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *Client) StopUnits(ids []string, f func(error)) {
	var call Call
	type Request struct {
		Ids []string `json:"ids"`
	}
	var req Request
	req.Ids = ids
	call.function = "stop_units"
	call.request, _ = json.MarshalIndent(req, "", " ")
	call.onResponse = func(call *Call) {
		if f != nil {
			f(nil)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *Client) GetUnitValues(unitName string, f func([]common_interfaces.ItemGetUnitItems, error)) {
	var call Call
	type Request struct {
		UnitName string `json:"unit_name"`
	}
	var req Request
	req.UnitName = unitName
	call.function = "unit_values"
	call.request, _ = json.MarshalIndent(req, "", " ")
	call.onResponse = func(call *Call) {
		var result []common_interfaces.ItemGetUnitItems
		result = make([]common_interfaces.ItemGetUnitItems, 0)
		err := json.Unmarshal([]byte(call.response), &result)
		if f != nil {
			f(result, err)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *Client) GetItemsValues(items []string, f func([]common_interfaces.ItemGetUnitItems, error)) {
	var call Call
	type Request struct {
		Items []string `json:"items"`
	}
	var req Request
	req.Items = items
	call.function = "items"
	call.request, _ = json.MarshalIndent(req, "", " ")
	call.onResponse = func(call *Call) {
		var result []common_interfaces.ItemGetUnitItems
		result = make([]common_interfaces.ItemGetUnitItems, 0)
		err := json.Unmarshal([]byte(call.response), &result)
		if f != nil {
			f(result, err)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *Client) GetAllItems(f func([]common_interfaces.ItemGetUnitItems, error)) {
	var call Call
	call.function = "all_items"
	call.request = nil
	call.onResponse = func(call *Call) {
		var result []common_interfaces.ItemGetUnitItems
		result = make([]common_interfaces.ItemGetUnitItems, 0)
		err := json.Unmarshal([]byte(call.response), &result)
		if f != nil {
			f(result, err)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *Client) GetCloudChannelValues(channelId string, f func([]common_interfaces.Item, error)) {
	var call Call
	type Request struct {
		ChannelId string `json:"id"`
	}
	var req Request
	req.ChannelId = channelId
	call.function = "cloud_channel_values"
	call.request, _ = json.MarshalIndent(req, "", " ")
	call.onResponse = func(call *Call) {
		var result []common_interfaces.Item
		result = make([]common_interfaces.Item, 0)
		err := json.Unmarshal([]byte(call.response), &result)
		if f != nil {
			f(result, err)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *Client) SetUnitConfig(unitId string, unitName string, unitConfig string, f func(error)) {
	var call Call
	type Request struct {
		UnitId     string `json:"id"`
		UnitName   string `json:"name"`
		UnitConfig string `json:"config"`
	}
	var req Request
	req.UnitId = unitId
	req.UnitName = unitName
	req.UnitConfig = unitConfig
	call.function = "set_unit_config"
	call.request, _ = json.MarshalIndent(req, "", " ")
	call.onResponse = func(call *Call) {
		if f != nil {
			f(nil)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *Client) GetUnitConfig(unitId string, f func(string, string, string, string, error)) {
	var call Call
	type Request struct {
		UnitId string `json:"id"`
	}
	var req Request
	req.UnitId = unitId
	call.function = "unit_config"
	call.request, _ = json.MarshalIndent(req, "", " ")
	call.onResponse = func(call *Call) {
		type Response struct {
			UnitId         string `json:"id"`
			UnitName       string `json:"name"`
			UnitType       string `json:"type"`
			UnitConfig     string `json:"config"`
			UnitConfigMeta string `json:"config_meta"`
		}
		var resp Response
		err := json.Unmarshal([]byte(call.response), &resp)
		if f != nil {
			f(resp.UnitName, resp.UnitConfig, resp.UnitConfigMeta, resp.UnitType, err)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *Client) GetUnitConfigByType(unitType string, f func(string, string, error)) {
	var call Call
	type Request struct {
		UnitType string `json:"type"`
	}
	var req Request
	req.UnitType = unitType
	call.function = "unit_config_by_type"
	call.request, _ = json.MarshalIndent(req, "", " ")
	call.onResponse = func(call *Call) {
		type Response struct {
			UnitId         string `json:"id"`
			UnitName       string `json:"name"`
			UnitConfigMeta string `json:"config_meta"`
		}
		var resp Response
		err := json.Unmarshal([]byte(call.response), &resp)
		if f != nil {
			f(resp.UnitName, resp.UnitConfigMeta, err)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *Client) AddCloudChannel(channelName string, f func(error)) {
	var call Call
	type Request struct {
		ChannelName string `json:"name"`
	}
	var req Request
	req.ChannelName = channelName

	call.function = "add_cloud_channel"
	call.request, _ = json.MarshalIndent(req, "", " ")
	call.onResponse = func(call *Call) {
		if f != nil {
			f(call.err)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *Client) EditCloudChannel(channelId string, channelName string, f func(error)) {
	var call Call
	type Request struct {
		ChannelId   string `json:"id"`
		ChannelName string `json:"name"`
	}
	var req Request
	req.ChannelId = channelId
	req.ChannelName = channelName

	call.function = "edit_cloud_channel"
	call.request, _ = json.MarshalIndent(req, "", " ")
	call.onResponse = func(call *Call) {
		if f != nil {
			f(nil)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *Client) RemoveCloudChannel(channelId string, f func(error)) {
	var call Call
	type Request struct {
		ChannelId string `json:"id"`
	}
	var req Request
	req.ChannelId = channelId

	call.function = "remove_cloud_channel"
	call.request, _ = json.MarshalIndent(req, "", " ")
	call.onResponse = func(call *Call) {
		if f != nil {
			f(nil)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *Client) GetCloudChannels(f func([]cloud.ChannelInfo, error)) {
	var call Call
	type Request struct {
	}
	var req Request
	type Response struct {
		Channels []cloud.ChannelInfo
	}

	call.function = "cloud_channels"
	call.request, _ = json.MarshalIndent(req, "", " ")
	call.onResponse = func(call *Call) {
		var resp Response
		json.Unmarshal([]byte(call.response), &resp)
		if f != nil {
			f(resp.Channels, nil)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *Client) CloudAddItems(channels []string, items []string, f func(error)) {
	var call Call
	type Request struct {
		Channels []string `json:"ids"`
		Items    []string `json:"items"`
	}
	var req Request
	req.Channels = channels
	req.Items = items

	call.function = "cloud_add_items"
	call.request, _ = json.MarshalIndent(req, "", " ")
	call.onResponse = func(call *Call) {
		if f != nil {
			f(nil)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *Client) CloudRemoveItems(channels []string, items []string, f func(error)) {
	var call Call
	type Request struct {
		Channels []string `json:"ids"`
		Items    []string `json:"items"`
	}
	var req Request
	req.Channels = channels
	req.Items = items

	call.function = "cloud_remove_items"
	call.request, _ = json.MarshalIndent(req, "", " ")
	call.onResponse = func(call *Call) {
		if f != nil {
			f(nil)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *Client) ReadHistory(name string, dtBegin int64, dtEnd int64, f func(*history.ReadResult, error)) {
	var call Call
	type Request struct {
		Name    string `json:"name"`
		DTBegin int64  `json:"dt_begin"`
		DTEnd   int64  `json:"dt_end"`
	}
	var req Request
	req.Name = name
	req.DTBegin = dtBegin
	req.DTEnd = dtEnd
	call.function = "history"
	call.request, _ = json.MarshalIndent(req, "", " ")
	call.onResponse = func(call *Call) {
		var result history.ReadResult
		err := json.Unmarshal([]byte(call.response), &result)
		if f != nil {
			f(&result, err)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *Client) GetStatistics(f func(common_interfaces.Statistics, error)) {
	var call Call
	type Request struct {
	}
	var req Request
	type Response struct {
		Stat common_interfaces.Statistics
	}

	call.function = "statistics"
	call.request, _ = json.MarshalIndent(req, "", " ")
	call.onResponse = func(call *Call) {
		var resp Response
		json.Unmarshal([]byte(call.response), &resp)
		if f != nil {
			f(resp.Stat, call.err)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *Client) ResAdd(name string, tp string, content []byte, f func(string, error)) {
	var call Call
	type Request struct {
		Name    string `json:"name"`
		Type    string `json:"type"`
		Content []byte `json:"content"`
	}
	var req Request
	req.Name = name
	req.Type = tp
	req.Content = content

	call.function = "res_add"
	call.request, _ = json.MarshalIndent(req, "", " ")
	call.onResponse = func(call *Call) {
		type Response struct {
			Id string
		}
		var resp Response
		_ = json.Unmarshal([]byte(call.response), &resp)

		if f != nil {
			f(resp.Id, nil)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *Client) ResSet(id string, thumbnail []byte, content []byte, f func(error)) {
	var call Call
	type Request struct {
		Id        string `json:"id"`
		Thumbnail []byte `json:"thumbnail"`
		Content   []byte `json:"content"`
	}
	var req Request
	req.Id = id
	req.Thumbnail = thumbnail
	req.Content = content

	call.function = "res_set"
	call.request, _ = json.MarshalIndent(req, "", " ")
	call.onResponse = func(call *Call) {
		if f != nil {
			f(call.err)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *Client) ResGet(id string, f func(*common_interfaces.ResourcesItem, error)) {
	var call Call
	type Request struct {
		Id string `json:"id"`
	}
	var req Request
	req.Id = id

	call.function = "res_get"
	call.request, _ = json.MarshalIndent(req, "", " ")
	call.onResponse = func(call *Call) {
		if f != nil {
			type Response struct {
				Item *common_interfaces.ResourcesItem
			}
			var resp Response
			_ = json.Unmarshal([]byte(call.response), &resp)
			f(resp.Item, call.err)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *Client) ResList(tp string, filter string, offset int, maxCount int, f func(common_interfaces.ResourcesInfo, error)) {
	var call Call
	type Request struct {
		Type     string `json:"type"`
		Filter   string `json:"filter"`
		Offset   int    `json:"offset"`
		MaxCount int    `json:"max_count"`
	}
	var req Request
	req.Type = tp
	req.Filter = filter
	req.Offset = offset
	req.MaxCount = maxCount

	call.function = "res_list"
	call.request, _ = json.MarshalIndent(req, "", " ")
	call.onResponse = func(call *Call) {
		if f != nil {
			type Response struct {
				Items common_interfaces.ResourcesInfo
			}
			var resp Response

			_ = json.Unmarshal([]byte(call.response), &resp)
			f(resp.Items, call.err)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *Client) ResRemove(id string, f func(error)) {
	var call Call
	type Request struct {
		Id string `json:"id"`
	}
	var req Request
	req.Id = id

	call.function = "res_remove"
	call.request, _ = json.MarshalIndent(req, "", " ")
	call.onResponse = func(call *Call) {
		if f != nil {
			type Response struct {
			}
			var resp Response
			_ = json.Unmarshal([]byte(call.response), &resp)
			f(call.err)
		}
	}
	call.client = c
	go c.thCall(&call)
}

func (c *Client) ResRename(id string, name string, f func(error)) {
	var call Call
	type Request struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}
	var req Request
	req.Id = id
	req.Name = name

	call.function = "res_rename"
	call.request, _ = json.MarshalIndent(req, "", " ")
	call.onResponse = func(call *Call) {
		if f != nil {
			type Response struct {
			}
			var resp Response
			_ = json.Unmarshal([]byte(call.response), &resp)
			f(call.err)
		}
	}
	call.client = c
	go c.thCall(&call)
}
