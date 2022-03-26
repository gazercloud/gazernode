package system

import (
	"errors"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/product/productinfo"
	"github.com/gazercloud/gazernode/system/history"
	"github.com/gazercloud/gazernode/system/protocols/nodeinterface"
	"github.com/gazercloud/gazernode/utilities/logger"
	"strconv"
	"strings"
	"time"
)

func (c *System) SetItemByName(name string, value string, UOM string, dt time.Time, external bool) error {
	var item *common_interfaces.Item
	if name == "" {
		return nil
	}

	c.mtx.Lock()
	if i, ok := c.itemsByName[name]; ok {
		item = i
	} else {
		item = common_interfaces.NewItem()
		item.Id = c.nextItemId
		item.Name = name
		c.itemsByName[item.Name] = item
		c.itemsById[item.Id] = item
		c.items = append(c.items, item)
		c.nextItemId++
	}
	c.mtx.Unlock()

	var itemValue common_interfaces.ItemValue
	itemValue.Value = value
	itemValue.DT = dt.UnixNano() / 1000
	itemValue.UOM = UOM
	return c.SetItem(item.Id, itemValue, 0, external)
}

func (c *System) SetAllItemsByUnitName(unitName string, value string, UOM string, dt time.Time, external bool) error {
	items := make([]*common_interfaces.Item, 0)
	if unitName == "" {
		return nil
	}

	c.mtx.Lock()
	for _, i := range c.items {
		if strings.HasPrefix(i.Name, unitName+"/") {
			items = append(items, i)
		}
	}
	c.mtx.Unlock()

	for _, i := range items {
		var itemValue common_interfaces.ItemValue
		itemValue.Value = value
		itemValue.DT = dt.UnixNano() / 1000
		itemValue.UOM = UOM
		c.SetItem(i.Id, itemValue, 0, external)
	}

	return nil
}

func (c *System) SetItem(itemId uint64, value common_interfaces.ItemValue, counter int, external bool) error {
	var item *common_interfaces.Item
	var watchersUnits []string

	counter++
	if counter > 10 {
		return errors.New("recursion detected")
	}
	c.mtx.Lock()
	if i, ok := c.itemsById[itemId]; ok {
		item = i
		value.Value = item.PostprocessingValue(value.Value)
		item.Value = value
	}
	c.mtx.Unlock()
	if item == nil {
		logger.Println("set item error: ", itemId, "=", value.Value)
		return errors.New("item not found")
	}

	if watcher, ok := c.itemWatchers[item.Name]; ok {
		watchersUnits = make([]string, 0)
		for watcherUnitId, _ := range watcher.UnitIDs {
			watchersUnits = append(watchersUnits, watcherUnitId)
		}
	}

	c.history.Write(item.Id, value)

	if external {
		for _, unitId := range watchersUnits {
			logger.Println("SendToWatcher", unitId, value.Value)
			c.unitsSystem.SendToWatcher(unitId, item.Name, item.Value)
		}
	}

	for _, itemDest := range item.TranslateToItems {
		c.SetItem(itemDest.Id, value, counter, true)
	}
	return nil
}

func (c *System) DataItemPropSet(itemName string, props []nodeinterface.PropItem) error {
	c.mtx.Lock()
	if item, ok := c.itemsByName[itemName]; ok {
		for _, prop := range props {
			item.Properties[prop.PropName] = &common_interfaces.ItemProperty{
				Name:  prop.PropName,
				Value: prop.PropValue,
			}
		}
		c.applyItemsProperties()
	} else {
		c.mtx.Unlock()
		return errors.New("item not found")
	}
	c.mtx.Unlock()
	c.SaveConfig()
	return nil
}

func (c *System) applyItemsProperties() {
	for _, item := range c.items {
		for _, prop := range item.Properties {
			if prop.Name == "source" {
				if prop.Value != "" {
					srcItemId, err := strconv.ParseUint(prop.Value, 10, 64)
					if err == nil {
						for _, itemToClearUp := range c.items {
							delete(itemToClearUp.TranslateToItems, item.Id)
						}
						if srcItem, ok := c.itemsById[srcItemId]; ok {
							srcItem.TranslateToItems[item.Id] = item
						}
					}
				} else {
					for _, itemToClearUp := range c.items {
						delete(itemToClearUp.TranslateToItems, item.Id)
					}
				}
			}

			if prop.Name == "tune_trim" {
				item.PostprocessingTrim = prop.Value == "1"
			}
			if prop.Name == "tune_on" {
				item.PostprocessingAdjust = prop.Value == "1"
			}
			if prop.Name == "tune_scale" {
				item.PostprocessingScale, _ = strconv.ParseFloat(prop.Value, 64)
			}
			if prop.Name == "tune_offset" {
				item.PostprocessingOffset, _ = strconv.ParseFloat(prop.Value, 64)
			}
			if prop.Name == "tune_precision" {
				precision, _ := strconv.ParseInt(prop.Value, 10, 64)
				item.PostprocessingPrecision = int(precision)
			}
		}
	}
}

func (c *System) DataItemPropGet(itemName string) ([]nodeinterface.PropItem, error) {
	result := make([]nodeinterface.PropItem, 0)

	c.mtx.Lock()
	if item, ok := c.itemsByName[itemName]; ok {
		for _, prop := range item.Properties {
			result = append(result, nodeinterface.PropItem{
				PropName:  prop.Name,
				PropValue: prop.Value,
			})

			if prop.Name == "source" {
				if prop.Value != "" {
					sourceItemId, errParseSourceItemId := strconv.ParseUint(prop.Value, 10, 64)
					if errParseSourceItemId == nil {
						if itemSource, ok := c.itemsById[sourceItemId]; ok {
							result = append(result, nodeinterface.PropItem{
								PropName:  "#source_item_name",
								PropValue: itemSource.Name,
							})
						}
					}
				}
			}
		}
	} else {
		c.mtx.Unlock()
		return nil, errors.New("item not found")
	}
	c.mtx.Unlock()
	return result, nil
}

type ItemWatcher struct {
	UnitIDs map[string]bool
}

func (c *System) AddToWatch(unitId string, itemName string) {
	c.mtx.Lock()
	watcher, ok := c.itemWatchers[itemName]
	if !ok {
		watcher = &ItemWatcher{
			UnitIDs: make(map[string]bool),
		}
		c.itemWatchers[itemName] = watcher
	}
	watcher.UnitIDs[unitId] = true
	c.mtx.Unlock()
}

func (c *System) RemoveFromWatch(unitId string, itemName string) {
	c.mtx.Lock()
	watcher, ok := c.itemWatchers[itemName]
	if ok {
		delete(watcher.UnitIDs, unitId)
	}
	if len(watcher.UnitIDs) == 0 {
		delete(c.itemWatchers, itemName)
	}
	c.mtx.Unlock()
}

func (c *System) SetProperty(itemName string, propName string, propValue string) {
	item, err := c.TouchItem(itemName)
	if err == nil {
		c.mtx.Lock()
		item.SetProperty(propName, propValue)
		c.mtx.Unlock()
	}
}

func (c *System) SetPropertyIfDoesntExist(itemName string, propName string, propValue string) {
	item, err := c.TouchItem(itemName)
	if err == nil {
		c.mtx.Lock()
		item.SetPropertyIfDoesntExist(propName, propValue)
		c.mtx.Unlock()
	}
}

func (c *System) TouchItem(name string) (*common_interfaces.Item, error) {
	var item *common_interfaces.Item
	fullName := name
	c.mtx.Lock()
	var ok bool
	if item, ok = c.itemsByName[fullName]; !ok {
		item = common_interfaces.NewItem()
		item.Id = c.nextItemId
		item.Name = fullName
		c.itemsByName[item.Name] = item
		c.itemsById[item.Id] = item
		c.items = append(c.items, item)
		c.nextItemId++
	}
	c.mtx.Unlock()
	return item, nil
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

	//c.publicChannels.RemoveItems(nil, itemsNames)

	err = c.SaveConfig()

	if len(itemsForRemove) == 0 {
		return errors.New("no items found")
	}

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

	//c.publicChannels.RenameItems(oldPrefix, newPrefix)
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
	//res.CloudSentBytes = c.publicChannels.SentBytes()
	//res.CloudReceivedBytes = c.publicChannels.ReceivedBytes()
	res.ApiCalls = c.apiCallsCount
	return res, nil
}

func (c *System) GetApi() (nodeinterface.ServiceApiResponse, error) {
	var res nodeinterface.ServiceApiResponse
	res.Product = productinfo.Name()
	res.Version = productinfo.Version()
	res.BuildTime = productinfo.BuildTime()
	res.SupportedFunctions = nodeinterface.ApiFunctions()

	return res, nil
}

func (c *System) SetNodeName(name string) error {
	c.nodeName = name
	return c.SaveConfig()
}

func (c *System) NodeName() string {
	return c.nodeName
}

func (c *System) GetInfo() (nodeinterface.ServiceInfoResponse, error) {
	var res nodeinterface.ServiceInfoResponse
	res.NodeName = c.NodeName()
	res.Version = productinfo.Version()
	res.BuildTime = productinfo.BuildTime()
	return res, nil
}
