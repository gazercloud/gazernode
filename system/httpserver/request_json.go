package httpserver

import (
	"encoding/json"
	"errors"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/history"
	"github.com/gazercloud/gazernode/logger"
	"github.com/gazercloud/gazernode/system/cloud"
	"github.com/gazercloud/gazernode/system/units/units_common"
	"time"
)

func (c *HttpServer) requestJson(function string, requestText string) (string, error) {
	var err error
	var result string
	switch function {
	case "unit_types":
		{
			type Request struct {
				Category string `json:"category"`
				Filter   string `json:"filter"`
				Offset   int    `json:"offset"`
				MaxCount int    `json:"max_count"`
			}
			var req Request
			var unitTypes common_interfaces.UnitTypes
			unitTypes.Types = make([]common_interfaces.UnitTypeInfo, 0)
			err = json.Unmarshal([]byte(requestText), &req)
			if err == nil {
				unitTypes = c.system.UnitTypes(req.Category, req.Filter, req.Offset, req.MaxCount)
			}

			var resultBytes []byte
			resultBytes, err = json.MarshalIndent(unitTypes, "", " ")
			if err == nil {
				result = string(resultBytes)
			}
		}
	case "unit_categories":
		{
			type Request struct {
			}
			var req Request
			var unitCategories []common_interfaces.UnitCategoryInfo
			unitCategories = make([]common_interfaces.UnitCategoryInfo, 0)
			err = json.Unmarshal([]byte(requestText), &req)
			if err == nil {
				unitCategories = c.system.UnitCategories()
			}

			var resultBytes []byte
			resultBytes, err = json.MarshalIndent(unitCategories, "", " ")
			if err == nil {
				result = string(resultBytes)
			}
		}
	case "units":
		{
			var resultBytes []byte
			var units []units_common.UnitInfo
			units = c.system.ListOfUnits()
			resultBytes, err = json.MarshalIndent(units, "", " ")
			if err == nil {
				result = string(resultBytes)
			}
		}
	case "lookup":
		{
			type Request struct {
				Entity     string `json:"entity"`
				Parameters string `json:"parameters"`
			}
			var req Request
			err = json.Unmarshal([]byte(requestText), &req)
			if err == nil {
				var resultBytes []byte
				var lookupResult *units_common.LookupResult
				lookupResult = c.system.Lookup(req.Entity)
				resultBytes, err = json.MarshalIndent(lookupResult, "", " ")
				if err == nil {
					result = string(resultBytes)
				}
			}
		}
	case "write":
		{
			type Request struct {
				ItemName string `json:"item_name"`
				Value    string `json:"value"`
			}

			var req Request
			err = json.Unmarshal([]byte(requestText), &req)
			if err == nil {
				err = c.system.SetItem(req.ItemName, req.Value, "", time.Now().UTC(), "")
			}
		}
	case "start_units":
		{
			type Request struct {
				Ids []string `json:"ids"`
			}
			var req Request
			err = json.Unmarshal([]byte(requestText), &req)
			if err == nil {
				err = c.system.StartUnits(req.Ids)
			}
		}
	case "stop_units":
		{
			type Request struct {
				Ids []string `json:"ids"`
			}
			var req Request
			err = json.Unmarshal([]byte(requestText), &req)
			if err == nil {
				err = c.system.StopUnits(req.Ids)
			}
		}
	case "set_unit_config":
		{
			type Request struct {
				UnitId     string `json:"id"`
				UnitName   string `json:"name"`
				UnitConfig string `json:"config"`
			}
			var req Request
			err = json.Unmarshal([]byte(requestText), &req)
			if err == nil {
				err = c.system.SetConfig(req.UnitId, req.UnitName, req.UnitConfig)
			}
		}
	case "unit_config":
		{
			type Request struct {
				UnitId string `json:"id"`
			}
			type Response struct {
				UnitId         string `json:"id"`
				UnitName       string `json:"name"`
				UnitType       string `json:"type"`
				UnitConfig     string `json:"config"`
				UnitConfigMeta string `json:"config_meta"`
			}
			var req Request
			var resp Response
			var responseBytes []byte
			err = json.Unmarshal([]byte(requestText), &req)
			if err == nil {
				resp.UnitName, resp.UnitConfig, resp.UnitConfigMeta, resp.UnitType, err = c.system.GetConfig(req.UnitId)
			}
			responseBytes, err = json.MarshalIndent(resp, "", " ")
			if err == nil {
				result = string(responseBytes)
			}
		}
	case "unit_config_by_type":
		{
			type Request struct {
				UnitType string `json:"type"`
			}
			type Response struct {
				UnitId         string `json:"id"`
				UnitName       string `json:"name"`
				UnitConfigMeta string `json:"config_meta"`
			}
			var req Request
			var resp Response
			var responseBytes []byte
			err = json.Unmarshal([]byte(requestText), &req)
			if err == nil {
				resp.UnitName, resp.UnitConfigMeta, err = c.system.GetConfigByType(req.UnitType)
			}
			responseBytes, err = json.MarshalIndent(resp, "", " ")
			if err == nil {
				result = string(responseBytes)
			}
		}
	case "remove_unit":
		{
			type Request struct {
				Units []string `json:"ids"`
			}
			var req Request
			err = json.Unmarshal([]byte(requestText), &req)
			if err == nil {
				err = c.system.RemoveUnits(req.Units)
			}
		}
	case "add_unit":
		{
			type Request struct {
				UnitType string `json:"type"`
				UnitName string `json:"name"`
			}
			var req Request
			err = json.Unmarshal([]byte(requestText), &req)
			if err == nil {
				var unitId string
				unitId, err = c.system.AddUnit(req.UnitName, req.UnitType)
				if err == nil {

					type Response struct {
						UnitId string `json:"unit_id"`
					}

					var resp Response
					resp.UnitId = unitId

					var resultBytes []byte
					resultBytes, err = json.MarshalIndent(resp, "", " ")
					if err == nil {
						result = string(resultBytes)
					}
				}

			}
		}
	case "unit_state":
		{
			type Request struct {
				UnitId string `json:"id"`
			}
			var req Request
			err = json.Unmarshal([]byte(requestText), &req)
			if err == nil {
				var unitState common_interfaces.UnitState
				unitState, err = c.system.GetUnitState(req.UnitId)

				var resultBytes []byte
				resultBytes, err = json.MarshalIndent(unitState, "", " ")
				if err == nil {
					result = string(resultBytes)
				}

			}
		}
	case "unit_values":
		{
			type Request struct {
				UnitName string `json:"unit_name"`
			}
			var resultBytes []byte
			var unitValues []common_interfaces.ItemGetUnitItems
			var req Request
			err = json.Unmarshal([]byte(requestText), &req)
			if err == nil {
				unitValues = c.system.GetUnitValues(req.UnitName)
			}

			resultBytes, err = json.MarshalIndent(unitValues, "", " ")
			if err == nil {
				result = string(resultBytes)
			}

		}
	case "items":
		{
			type Request struct {
				Items []string `json:"items"`
			}
			var resultBytes []byte
			var unitValues []common_interfaces.ItemGetUnitItems
			var req Request
			err = json.Unmarshal([]byte(requestText), &req)
			if err == nil {
				unitValues = c.system.GetItemsValues(req.Items)
			}

			resultBytes, err = json.MarshalIndent(unitValues, "", " ")
			if err == nil {
				result = string(resultBytes)
			}

		}
	case "all_items":
		{
			var resultBytes []byte
			var unitValues []common_interfaces.ItemGetUnitItems
			unitValues = c.system.GetAllItems()

			resultBytes, err = json.MarshalIndent(unitValues, "", " ")
			if err == nil {
				result = string(resultBytes)
			}

		}
	case "cloud_channel_values":
		{
			type Request struct {
				ChannelId string `json:"id"`
			}
			var resultBytes []byte
			var unitValues []common_interfaces.Item
			var req Request
			err = json.Unmarshal([]byte(requestText), &req)
			if err == nil {
				unitValues, err = c.system.GetCloudChannelValues(req.ChannelId)
			}

			resultBytes, err = json.MarshalIndent(unitValues, "", " ")
			if err == nil {
				result = string(resultBytes)
			}

		}
	case "add_cloud_channel":
		{
			type Request struct {
				ChannelName string `json:"name"`
			}
			var resultBytes []byte
			var req Request
			err = json.Unmarshal([]byte(requestText), &req)
			if err == nil {
				err = c.system.AddCloudChannel(req.ChannelName)
			}

			resultBytes, err = json.MarshalIndent("", "", " ")
			if err == nil {
				result = string(resultBytes)
			}

		}
	case "edit_cloud_channel":
		{
			type Request struct {
				ChannelId   string `json:"id"`
				ChannelName string `json:"name"`
			}
			var resultBytes []byte
			var req Request
			err = json.Unmarshal([]byte(requestText), &req)
			if err == nil {
				err = c.system.EditCloudChannel(req.ChannelId, req.ChannelName)
			}

			resultBytes, err = json.MarshalIndent("", "", " ")
			if err == nil {
				result = string(resultBytes)
			}

		}
	case "remove_cloud_channel":
		{
			type Request struct {
				ChannelId string `json:"id"`
			}
			var resultBytes []byte
			var req Request
			err = json.Unmarshal([]byte(requestText), &req)
			if err == nil {
				err = c.system.RemoveCloudChannel(req.ChannelId)
			}

			resultBytes, err = json.MarshalIndent("", "", " ")
			if err == nil {
				result = string(resultBytes)
			}

		}
	case "cloud_add_items":
		{
			type Request struct {
				Channels []string `json:"ids"`
				Items    []string `json:"items"`
			}
			var resultBytes []byte
			var req Request
			err = json.Unmarshal([]byte(requestText), &req)
			if err == nil {
				err = c.system.CloudAddItems(req.Channels, req.Items)
			}

			resultBytes, err = json.MarshalIndent("", "", " ")
			if err == nil {
				result = string(resultBytes)
			}

		}
	case "cloud_remove_items":
		{
			type Request struct {
				Channels []string `json:"ids"`
				Items    []string `json:"items"`
			}
			var resultBytes []byte
			var req Request
			err = json.Unmarshal([]byte(requestText), &req)
			if err == nil {
				err = c.system.CloudRemoveItems(req.Channels, req.Items)
			}

			resultBytes, err = json.MarshalIndent("", "", " ")
			if err == nil {
				result = string(resultBytes)
			}

		}
	case "cloud_remove_all_items":
		{
			type Request struct {
				ChannelId string `json:"id"`
			}
			var resultBytes []byte
			var req Request
			err = json.Unmarshal([]byte(requestText), &req)
			if err == nil {
				err = c.system.CloudRemoveAllItems(req.ChannelId)
			}

			resultBytes, err = json.MarshalIndent("", "", " ")
			if err == nil {
				result = string(resultBytes)
			}

		}
	case "cloud_channels":
		{
			type Request struct {
			}
			var resultBytes []byte
			type Response struct {
				Channels []cloud.ChannelInfo
			}
			var req Request
			var resp Response
			err = json.Unmarshal([]byte(requestText), &req)
			if err == nil {
				resp.Channels, err = c.system.GetCloudChannels()
			}

			resultBytes, err = json.MarshalIndent(resp, "", " ")
			if err == nil {
				result = string(resultBytes)
			}
		}
	case "history":
		{
			type Request struct {
				Name    string `json:"name"`
				DTBegin int64  `json:"dt_begin"`
				DTEnd   int64  `json:"dt_end"`
			}

			var resultBytes []byte
			var readResult *history.ReadResult
			var req Request
			err = json.Unmarshal([]byte(requestText), &req)
			if err == nil {
				readResult, err = c.system.ReadHistory(req.Name, req.DTBegin, req.DTEnd)
			}

			resultBytes, err = json.MarshalIndent(readResult, "", " ")
			if err == nil {
				result = string(resultBytes)
			}
		}
	case "statistics":
		{
			type Request struct {
			}
			var resultBytes []byte
			type Response struct {
				Stat common_interfaces.Statistics
			}
			var req Request
			var resp Response
			err = json.Unmarshal([]byte(requestText), &req)
			if err == nil {
				resp.Stat, err = c.system.GetStatistics()
			}

			resultBytes, err = json.MarshalIndent(resp, "", " ")
			if err == nil {
				result = string(resultBytes)
			}
		}
	case "res_add":
		{
			type Request struct {
				Name    string
				Type    string
				Content []byte
			}
			var req Request
			type Response struct {
				Id string
			}
			var resp Response

			err = json.Unmarshal([]byte(requestText), &req)
			if err == nil {
				resp.Id, err = c.system.ResAdd(req.Name, req.Type, req.Content)
			}

			var resultBytes []byte
			resultBytes, err = json.MarshalIndent(resp, "", " ")
			if err == nil {
				result = string(resultBytes)
			}
		}
	case "res_set":
		{
			type Request struct {
				Id        string `json:"id"`
				Thumbnail []byte `json:"thumbnail"`
				Content   []byte `json:"content"`
			}
			var req Request
			type Response struct {
			}
			var resp Response

			err = json.Unmarshal([]byte(requestText), &req)
			if err == nil {
				err = c.system.ResSet(req.Id, req.Thumbnail, req.Content)
			}

			var resultBytes []byte
			resultBytes, err = json.MarshalIndent(resp, "", " ")
			if err == nil {
				result = string(resultBytes)
			}
		}
	case "res_get":
		{
			type Request struct {
				Id string
			}
			var req Request
			type Response struct {
				Item *common_interfaces.ResourcesItem
			}
			var resp Response

			err = json.Unmarshal([]byte(requestText), &req)
			if err == nil {
				resp.Item, err = c.system.ResGet(req.Id)
				if err == nil {
					var resultBytes []byte
					resultBytes, err = json.MarshalIndent(resp, "", " ")
					result = string(resultBytes)
				}
			}

		}
	case "res_remove":
		{
			type Request struct {
				Id string `json:"id"`
			}
			var req Request
			type Response struct {
			}
			var resp Response

			err = json.Unmarshal([]byte(requestText), &req)
			if err == nil {
				err = c.system.ResRemove(req.Id)
			}

			var resultBytes []byte
			resultBytes, err = json.MarshalIndent(resp, "", " ")
			if err == nil {
				result = string(resultBytes)
			}
		}
	case "res_rename":
		{
			type Request struct {
				Id   string `json:"id"`
				Name string `json:"name"`
			}
			var req Request
			type Response struct {
			}
			var resp Response

			err = json.Unmarshal([]byte(requestText), &req)
			if err == nil {
				err = c.system.ResRename(req.Id, req.Name)
			}

			var resultBytes []byte
			resultBytes, err = json.MarshalIndent(resp, "", " ")
			if err == nil {
				result = string(resultBytes)
			}
		}
	case "res_list":
		{
			type Request struct {
				Type     string `json:"type"`
				Filter   string `json:"filter"`
				Offset   int    `json:"offset"`
				MaxCount int    `json:"max_count"`
			}
			var req Request
			type Response struct {
				Items common_interfaces.ResourcesInfo
			}
			var resp Response

			err = json.Unmarshal([]byte(requestText), &req)
			if err == nil {
				resp.Items = c.system.ResList(req.Type, req.Filter, req.Offset, req.MaxCount)
			}

			var resultBytes []byte
			resultBytes, err = json.MarshalIndent(resp, "", " ")
			if err == nil {
				result = string(resultBytes)
			}
		}

	default:
		err = errors.New("function not supported")
	}

	if err == nil {
		return result, nil
	}

	logger.Println("Function execution error: ", err, "\r\n", requestText)
	return "", err
}

var TempValue int

func init() {
	TempValue = 5
}
