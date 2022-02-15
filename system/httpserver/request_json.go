package httpserver

import (
	"errors"
	"github.com/gazercloud/gazernode/system/protocols/nodeinterface"
	"github.com/gazercloud/gazernode/utilities/logger"
)

func (c *HttpServer) RequestJson(function string, requestText []byte, host string, fromCloud bool) ([]byte, error) {
	var err error
	var result []byte

	c.system.RegApiCall()

	switch function {

	// *** UnitType ***
	case nodeinterface.FuncUnitTypeList:
		result, err = c.UnitTypeList(requestText)
	case nodeinterface.FuncUnitTypeCategories:
		result, err = c.UnitTypeCategories(requestText)
	case nodeinterface.FuncUnitTypeConfigMeta:
		result, err = c.UnitTypeConfigMeta(requestText)

		// *** Unit ***
	case nodeinterface.FuncUnitAdd:
		result, err = c.UnitAdd(requestText, fromCloud)
	case nodeinterface.FuncUnitRemove:
		result, err = c.UnitRemove(requestText)
	case nodeinterface.FuncUnitState:
		result, err = c.UnitState(requestText)
	case nodeinterface.FuncUnitStateAll:
		result, err = c.UnitStateAll(requestText)
	case nodeinterface.FuncUnitItemsValues:
		result, err = c.UnitItemsValues(requestText)
	case nodeinterface.FuncUnitList:
		result, err = c.UnitList(requestText)
	case nodeinterface.FuncUnitStart:
		result, err = c.UnitStart(requestText)
	case nodeinterface.FuncUnitStop:
		result, err = c.UnitStop(requestText)
	case nodeinterface.FuncUnitSetConfig:
		result, err = c.UnitSetConfig(requestText, fromCloud)
	case nodeinterface.FuncUnitGetConfig:
		result, err = c.UnitGetConfig(requestText)
	case nodeinterface.FuncUnitPropSet:
		result, err = c.UnitPropSet(requestText)
	case nodeinterface.FuncUnitPropGet:
		result, err = c.UnitPropGet(requestText)

		// *** Service ***
	case nodeinterface.FuncServiceLookup:
		result, err = c.ServiceLookup(requestText)
	case nodeinterface.FuncServiceStatistics:
		result, err = c.ServiceStatistics(requestText)
	case nodeinterface.FuncServiceApi:
		result, err = c.ServiceApi(requestText)
	case nodeinterface.FuncServiceSetNodeName:
		result, err = c.ServiceSetNodeName(requestText)
	case nodeinterface.FuncServiceNodeName:
		result, err = c.ServiceNodeName(requestText)
	case nodeinterface.FuncServiceInfo:
		result, err = c.ServiceInfo(requestText)

		// *** Resource ***
	case nodeinterface.FuncResourceAdd:
		result, err = c.ResourceAdd(requestText)
	case nodeinterface.FuncResourceSet:
		result, err = c.ResourceSet(requestText)
	case nodeinterface.FuncResourceGet:
		result, err = c.ResourceGet(requestText)
	case nodeinterface.FuncResourceGetThumbnail:
		result, err = c.ResourceGetThumbnail(requestText)
	case nodeinterface.FuncResourceRemove:
		result, err = c.ResourceRemove(requestText)
	case nodeinterface.FuncResourceList:
		result, err = c.ResourceList(requestText)
	case nodeinterface.FuncResourcePropSet:
		result, err = c.ResourcePropSet(requestText)
	case nodeinterface.FuncResourcePropGet:
		result, err = c.ResourcePropGet(requestText)

		// *** Cloud ***
	case nodeinterface.FuncCloudLogin:
		result, err = c.CloudLogin(requestText)
	case nodeinterface.FuncCloudLogout:
		result, err = c.CloudLogout(requestText)
	case nodeinterface.FuncCloudState:
		result, err = c.CloudState(requestText)
	case nodeinterface.FuncCloudNodes:
		result, err = c.CloudNodes(requestText)
	case nodeinterface.FuncCloudAddNode:
		result, err = c.CloudAddNode(requestText)
	case nodeinterface.FuncCloudUpdateNode:
		result, err = c.CloudUpdateNode(requestText)
	case nodeinterface.FuncCloudRemoveNode:
		result, err = c.CloudRemoveNode(requestText)
	case nodeinterface.FuncCloudGetSettings:
		result, err = c.CloudGetSettings(requestText)
	case nodeinterface.FuncCloudSetSettings:
		result, err = c.CloudSetSettings(requestText)
	case nodeinterface.FuncCloudAccountInfo:
		result, err = c.CloudAccountInfo(requestText)
	case nodeinterface.FuncCloudSetCurrentNodeId:
		result, err = c.CloudSetCurrentNodeId(requestText)
	case nodeinterface.FuncCloudGetSettingsProfiles:
		result, err = c.CloudGetSettingsProfiles(requestText)

	// *** Data Item ***
	case nodeinterface.FuncDataItemList:
		result, err = c.DataItemList(requestText)
	case nodeinterface.FuncDataItemListAll:
		result, err = c.DataItemListAll(requestText)
	case nodeinterface.FuncDataItemWrite:
		result, err = c.DataItemWrite(requestText)
	case nodeinterface.FuncDataItemHistory:
		result, err = c.DataItemHistory(requestText)
	case nodeinterface.FuncDataItemHistoryChart:
		result, err = c.DataItemHistoryChart(requestText)
	case nodeinterface.FuncDataItemRemove:
		result, err = c.DataItemRemove(requestText)
	case nodeinterface.FuncDataItemPropSet:
		result, err = c.DataItemPropSet(requestText)
	case nodeinterface.FuncDataItemPropGet:
		result, err = c.DataItemPropGet(requestText)

		// *** Data Item ***
	case nodeinterface.FuncSessionOpen:
		result, err = c.SessionOpen(requestText, host)
	case nodeinterface.FuncSessionActivate:
		result, err = c.SessionActivate(requestText)
	case nodeinterface.FuncSessionRemove:
		result, err = c.SessionRemove(requestText)
	case nodeinterface.FuncSessionList:
		result, err = c.SessionList(requestText)

	case nodeinterface.FuncUserList:
		result, err = c.UserList(requestText)
	case nodeinterface.FuncUserAdd:
		result, err = c.UserAdd(requestText)
	case nodeinterface.FuncUserSetPassword:
		result, err = c.UserSetPassword(requestText)
	case nodeinterface.FuncUserRemove:
		result, err = c.UserRemove(requestText)
	case nodeinterface.FuncUserPropSet:
		result, err = c.UserPropSet(requestText)
	case nodeinterface.FuncUserPropGet:
		result, err = c.UserPropGet(requestText)

	default:
		err = errors.New("function not supported")
	}

	if err == nil {
		return result, nil
	}

	logger.Println("Function execution error: ", err, "\r\n", function, string(requestText))
	return nil, err
}

var TempValue int

func init() {
	TempValue = 5
}
