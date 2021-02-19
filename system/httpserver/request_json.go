package httpserver

import (
	"errors"
	"github.com/gazercloud/gazernode/logger"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
)

func (c *HttpServer) requestJson(function string, requestText []byte) ([]byte, error) {
	var err error
	var result []byte
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
		result, err = c.UnitAdd(requestText)
	case nodeinterface.FuncUnitRemove:
		result, err = c.UnitRemove(requestText)
	case nodeinterface.FuncUnitState:
		result, err = c.UnitState(requestText)
	case nodeinterface.FuncUnitItemsValues:
		result, err = c.UnitItemsValues(requestText)
	case nodeinterface.FuncUnitList:
		result, err = c.UnitList(requestText)
	case nodeinterface.FuncUnitStart:
		result, err = c.UnitStart(requestText)
	case nodeinterface.FuncUnitStop:
		result, err = c.UnitStop(requestText)
	case nodeinterface.FuncUnitSetConfig:
		result, err = c.UnitSetConfig(requestText)
	case nodeinterface.FuncUnitGetConfig:
		result, err = c.UnitGetConfig(requestText)

		// *** Service ***
	case nodeinterface.FuncServiceLookup:
		result, err = c.ServiceLookup(requestText)
	case nodeinterface.FuncServiceStatistics:
		result, err = c.ServiceStatistics(requestText)

		// *** Resource ***
	case nodeinterface.FuncResourceAdd:
		result, err = c.ResourceAdd(requestText)
	case nodeinterface.FuncResourceSet:
		result, err = c.ResourceSet(requestText)
	case nodeinterface.FuncResourceGet:
		result, err = c.ResourceGet(requestText)
	case nodeinterface.FuncResourceRemove:
		result, err = c.ResourceRemove(requestText)
	case nodeinterface.FuncResourceRename:
		result, err = c.ResourceRename(requestText)
	case nodeinterface.FuncResourceList:
		result, err = c.ResourceList(requestText)

		// *** Public Channel ***
	case nodeinterface.FuncPublicChannelList:
		result, err = c.PublicChannelList(requestText)
	case nodeinterface.FuncPublicChannelAdd:
		result, err = c.PublicChannelAdd(requestText)
	case nodeinterface.FuncPublicChannelSetName:
		result, err = c.PublicChannelSetName(requestText)
	case nodeinterface.FuncPublicChannelRemove:
		result, err = c.PublicChannelRemove(requestText)
	case nodeinterface.FuncPublicChannelItemAdd:
		result, err = c.PublicChannelItemAdd(requestText)
	case nodeinterface.FuncPublicChannelItemRemove:
		result, err = c.PublicChannelItemRemove(requestText)
	case nodeinterface.FuncPublicChannelItemsState:
		result, err = c.PublicChannelItemsState(requestText)

		// *** Data Item ***
	case nodeinterface.FuncDataItemList:
		result, err = c.DataItemList(requestText)
	case nodeinterface.FuncDataItemListAll:
		result, err = c.DataItemListAll(requestText)
	case nodeinterface.FuncDataItemWrite:
		result, err = c.DataItemWrite(requestText)
	case nodeinterface.FuncDataItemHistory:
		result, err = c.DataItemHistory(requestText)

	default:
		err = errors.New("function not supported")
	}

	if err == nil {
		return result, nil
	}

	logger.Println("Function execution error: ", err, "\r\n", requestText)
	return nil, err
}

var TempValue int

func init() {
	TempValue = 5
}
