package system

import (
	"errors"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/product/productinfo"
	"github.com/gazercloud/gazernode/system/history"
	nodeinterface2 "github.com/gazercloud/gazernode/system/protocols/nodeinterface"
	"strings"
	"time"
)

func (c *System) SetItemByName(name string, value string, UOM string, dt time.Time, external bool) error {
	var item *common_interfaces.Item

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
	return c.SetItem(item.Id, itemValue)
}

func (c *System) SetItem(itemId uint64, value common_interfaces.ItemValue) error {
	var item *common_interfaces.Item
	c.mtx.Lock()
	if i, ok := c.itemsById[itemId]; ok {
		item = i
		item.Value = value
	}
	c.mtx.Unlock()
	if item == nil {
		return errors.New("item not found")
	}
	c.history.Write(item.Id, value)
	c.processTranslators(item)
	return nil
}

func (c *System) SetItemSource(itemName string, source string) error {
	var itemSrc *common_interfaces.Item
	var itemDest *common_interfaces.Item
	var okSrc, okDest bool

	c.mtx.Lock()
	defer c.mtx.Unlock()

	itemSrc, okSrc = c.itemsByName[source]
	itemDest, okDest = c.itemsByName[itemName]

	if source != "" && !okSrc {
		return errors.New("source item not found")
	}
	if !okDest {
		return errors.New("dest item not found")
	}

	translatorsForDelete := make([]uint64, 0)
	for _, tr := range c.translators {
		if tr.Dest == itemDest.Id {
			translatorsForDelete = append(translatorsForDelete, tr.Id)
		}
	}

	for _, translatorForDelete := range translatorsForDelete {
		delete(c.translators, translatorForDelete)
	}

	if okSrc {
		_, err := c.addTranslator(Translator{
			Id:   0,
			Src:  itemSrc.Id,
			Dest: itemDest.Id,
		})

		return err
	}
	return nil
}

func (c *System) addTranslator(translator Translator) (uint64, error) {
	maxTranslatorId := uint64(0)
	for _, tr := range c.translators {
		if tr.Id > maxTranslatorId {
			maxTranslatorId = tr.Id
		}
	}
	if maxTranslatorId == ^uint64(0) {
		return 0, errors.New("no more translator's identifiers")
	}
	nextTranslatorId := maxTranslatorId + 1
	translator.Id = nextTranslatorId
	c.translators[nextTranslatorId] = &translator
	return nextTranslatorId, nil
}

func (c *System) processTranslators(src *common_interfaces.Item) {
	involvedTranslators := make([]*Translator, 0)
	c.mtx.Lock()
	for _, tr := range c.translators {
		if tr.Src == src.Id {
			involvedTranslators = append(involvedTranslators, tr)
		}
	}
	c.mtx.Unlock()
	for _, translator := range involvedTranslators {
		c.executeTranslator(translator)
	}
}

func (c *System) executeTranslator(translator *Translator) error {
	c.mtx.Lock()
	itemSrc, itemSrcExists := c.itemsById[translator.Src]
	itemDest, itemDestExists := c.itemsById[translator.Dest]
	c.mtx.Unlock()
	if !itemSrcExists {
		return errors.New("src item not found")
	}
	if !itemDestExists {
		return errors.New("dest item not found")
	}

	return c.SetItem(itemDest.Id, itemSrc.Value)
}

type ItemWatcher struct {
	UnitIDs map[string]bool
}

type Translator struct {
	Id   uint64
	Src  uint64
	Dest uint64
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

func (c *System) GetApi() (nodeinterface2.ServiceApiResponse, error) {
	var res nodeinterface2.ServiceApiResponse
	res.Product = productinfo.Name()
	res.Version = productinfo.Version()
	res.BuildTime = productinfo.BuildTime()
	res.SupportedFunctions = nodeinterface2.ApiFunctions()

	return res, nil
}

func (c *System) SetNodeName(name string) error {
	c.nodeName = name
	return c.SaveConfig()
}

func (c *System) NodeName() string {
	return c.nodeName
}

func (c *System) GetInfo() (nodeinterface2.ServiceInfoResponse, error) {
	var res nodeinterface2.ServiceInfoResponse
	res.NodeName = c.NodeName()
	res.Version = productinfo.Version()
	res.BuildTime = productinfo.BuildTime()
	return res, nil
}
