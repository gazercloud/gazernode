package system

import (
	"errors"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/history"
	"github.com/gazercloud/gazernode/product/productinfo"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
	"strings"
	"time"
)

func (c *System) SetItem(name string, value string, UOM string, dt time.Time, flags string) error {
	var item *common_interfaces.Item
	fullName := name
	c.mtx.Lock()
	if i, ok := c.itemsByName[fullName]; ok {
		item = i
	} else {
		item = common_interfaces.NewItem()
		item.Id = c.nextItemId
		item.Name = fullName
		c.itemsByName[item.Name] = item
		c.items = append(c.items, item)
		c.nextItemId++
	}
	item.Value.Value = value
	item.Value.DT = dt.UnixNano() / 1000
	item.Value.UOM = UOM
	//item.Value.Flags = flags
	c.mtx.Unlock()
	c.history.Write(item.Id, item.Value)
	return nil
}

func (c *System) TouchItem(name string) error {
	var item *common_interfaces.Item
	fullName := name
	c.mtx.Lock()
	if _, ok := c.itemsByName[fullName]; !ok {
		item = common_interfaces.NewItem()
		item.Id = c.nextItemId
		item.Name = fullName
		c.itemsByName[item.Name] = item
		c.items = append(c.items, item)
		c.nextItemId++
	}
	c.mtx.Unlock()
	return nil
}

func (c *System) GetItem(name string) (common_interfaces.Item, error) {
	var item common_interfaces.Item
	var found bool
	c.mtx.Lock()
	if i, ok := c.itemsByName[name]; ok {
		item = *i
		found = true
	}
	c.mtx.Unlock()

	if !found {
		return item, errors.New("no item found")
	}

	return item, nil
}

func (c *System) RemoveItems(itemsNames []string) error {
	var err error

	c.mtx.Lock()
	newItems := make([]*common_interfaces.Item, 0)
	itemsForRemove := make([]*common_interfaces.Item, 0)

	itemsNamesMap := make(map[string]bool)
	for _, itemName := range itemsNames {
		itemsNamesMap[itemName] = true
	}

	for _, item := range c.items {
		if _, ok := itemsNamesMap[item.Name]; ok {
			itemsForRemove = append(itemsForRemove, item)
		} else {
			newItems = append(newItems, item)
		}
	}
	c.items = newItems

	for _, item := range itemsForRemove {
		delete(c.itemsByName, item.Name)
		delete(c.itemsById, item.Id)
		c.history.RemoveItem(item.Id)
	}
	c.mtx.Unlock()

	c.publicChannels.RemoveItems(nil, itemsNames)

	err = c.SaveConfig()
	return err
}

func (c *System) GetItems() []common_interfaces.Item {
	var items []common_interfaces.Item
	c.mtx.Lock()
	items = make([]common_interfaces.Item, len(c.items))
	for index, item := range c.items {
		items[index] = *item
	}
	c.mtx.Unlock()
	return items
}

func (c *System) RenameItems(oldPrefix string, newPrefix string) {
	c.mtx.Lock()
	for _, item := range c.items {
		if strings.HasPrefix(item.Name, oldPrefix) {
			delete(c.itemsByName, item.Name)
			item.Name = strings.Replace(item.Name, oldPrefix, newPrefix, 1)
			c.itemsByName[item.Name] = item
		}
	}
	c.mtx.Unlock()

	c.publicChannels.RenameItems(oldPrefix, newPrefix)
}

func (c *System) ReadHistory(name string, dtBegin int64, dtEnd int64) (*history.ReadResult, error) {
	c.mtx.Lock()
	item, ok := c.itemsByName[name]
	c.mtx.Unlock()
	if ok {
		return c.history.Read(item.Id, dtBegin, dtEnd), nil
	}

	var result history.ReadResult
	return &result, errors.New("no item found")
}

func (c *System) GetStatistics() (common_interfaces.Statistics, error) {
	var res common_interfaces.Statistics
	res.CloudSentBytes = c.publicChannels.SentBytes()
	res.CloudReceivedBytes = c.publicChannels.ReceivedBytes()
	res.ApiCalls = c.apiCallsCount
	return res, nil
}

func (c *System) GetApi() (nodeinterface.ServiceApiResponse, error) {
	var res nodeinterface.ServiceApiResponse
	res.Product = productinfo.Name()
	res.Version = productinfo.Version()
	res.BuildTime = productinfo.BuildTime()

	res.SupportedFunctions = append(res.SupportedFunctions, nodeinterface.FuncUnitTypeList)
	res.SupportedFunctions = append(res.SupportedFunctions, nodeinterface.FuncUnitTypeCategories)
	res.SupportedFunctions = append(res.SupportedFunctions, nodeinterface.FuncUnitTypeConfigMeta)

	res.SupportedFunctions = append(res.SupportedFunctions, nodeinterface.FuncUnitAdd)
	res.SupportedFunctions = append(res.SupportedFunctions, nodeinterface.FuncUnitRemove)
	res.SupportedFunctions = append(res.SupportedFunctions, nodeinterface.FuncUnitState)
	res.SupportedFunctions = append(res.SupportedFunctions, nodeinterface.FuncUnitItemsValues)
	res.SupportedFunctions = append(res.SupportedFunctions, nodeinterface.FuncUnitList)
	res.SupportedFunctions = append(res.SupportedFunctions, nodeinterface.FuncUnitStart)
	res.SupportedFunctions = append(res.SupportedFunctions, nodeinterface.FuncUnitStop)
	res.SupportedFunctions = append(res.SupportedFunctions, nodeinterface.FuncUnitSetConfig)
	res.SupportedFunctions = append(res.SupportedFunctions, nodeinterface.FuncUnitGetConfig)

	res.SupportedFunctions = append(res.SupportedFunctions, nodeinterface.FuncDataItemList)
	res.SupportedFunctions = append(res.SupportedFunctions, nodeinterface.FuncDataItemListAll)
	res.SupportedFunctions = append(res.SupportedFunctions, nodeinterface.FuncDataItemWrite)
	res.SupportedFunctions = append(res.SupportedFunctions, nodeinterface.FuncDataItemHistory)

	res.SupportedFunctions = append(res.SupportedFunctions, nodeinterface.FuncPublicChannelList)
	res.SupportedFunctions = append(res.SupportedFunctions, nodeinterface.FuncPublicChannelAdd)
	res.SupportedFunctions = append(res.SupportedFunctions, nodeinterface.FuncPublicChannelSetName)
	res.SupportedFunctions = append(res.SupportedFunctions, nodeinterface.FuncPublicChannelRemove)
	res.SupportedFunctions = append(res.SupportedFunctions, nodeinterface.FuncPublicChannelItemAdd)
	res.SupportedFunctions = append(res.SupportedFunctions, nodeinterface.FuncPublicChannelItemRemove)
	res.SupportedFunctions = append(res.SupportedFunctions, nodeinterface.FuncPublicChannelItemsState)

	res.SupportedFunctions = append(res.SupportedFunctions, nodeinterface.FuncServiceLookup)
	res.SupportedFunctions = append(res.SupportedFunctions, nodeinterface.FuncServiceStatistics)
	res.SupportedFunctions = append(res.SupportedFunctions, nodeinterface.FuncServiceApi)

	res.SupportedFunctions = append(res.SupportedFunctions, nodeinterface.FuncResourceAdd)
	res.SupportedFunctions = append(res.SupportedFunctions, nodeinterface.FuncResourceSet)
	res.SupportedFunctions = append(res.SupportedFunctions, nodeinterface.FuncResourceGet)
	res.SupportedFunctions = append(res.SupportedFunctions, nodeinterface.FuncResourceRemove)
	res.SupportedFunctions = append(res.SupportedFunctions, nodeinterface.FuncResourceRename)
	res.SupportedFunctions = append(res.SupportedFunctions, nodeinterface.FuncResourceList)

	res.SupportedFunctions = append(res.SupportedFunctions, nodeinterface.FuncSessionOpen)
	res.SupportedFunctions = append(res.SupportedFunctions, nodeinterface.FuncSessionActivate)
	res.SupportedFunctions = append(res.SupportedFunctions, nodeinterface.FuncSessionRemove)
	res.SupportedFunctions = append(res.SupportedFunctions, nodeinterface.FuncSessionList)

	res.SupportedFunctions = append(res.SupportedFunctions, nodeinterface.FuncUserList)
	res.SupportedFunctions = append(res.SupportedFunctions, nodeinterface.FuncUserAdd)
	res.SupportedFunctions = append(res.SupportedFunctions, nodeinterface.FuncUserSetPassword)
	res.SupportedFunctions = append(res.SupportedFunctions, nodeinterface.FuncUserRemove)
	return res, nil
}
