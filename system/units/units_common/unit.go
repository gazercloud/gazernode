package units_common

import (
	"errors"
	"fmt"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/utilities/uom"
	"strconv"
	"sync"
	"time"
)

type Unit struct {
	mtx sync.Mutex

	unitId          string
	unitType        string
	unitDisplayName string
	config          string
	iUnit           common_interfaces.IUnit
	iDataStorage    common_interfaces.IDataStorage
	lastError       string
	lastErrorDT     time.Time
	lastInfo        string
	lastInfoDT      time.Time

	Properties map[string]common_interfaces.ItemProperty

	lastLogDT time.Time

	Started  bool
	Stopping bool

	watchItems map[string]bool
}

func (c *Unit) Init() {
	c.Properties = make(map[string]common_interfaces.ItemProperty)
}

func (c *Unit) PropSetIfNotExists(name string, value string) {
	c.mtx.Lock()
	if _, ok := c.Properties[name]; !ok {
		c.Properties[name] = common_interfaces.ItemProperty{
			Name:  name,
			Value: value,
		}
	}
	c.mtx.Unlock()
}

func (c *Unit) Prop(name string) string {
	result := ""
	c.mtx.Lock()
	if prop, ok := c.Properties[name]; ok {
		result = prop.Value
	}
	c.mtx.Unlock()
	return result
}

func (c *Unit) PropSet(props []common_interfaces.ItemProperty) {
	c.mtx.Lock()
	for _, prop := range props {
		c.Properties[prop.Name] = prop
	}
	c.mtx.Unlock()
}

func (c *Unit) PropGet() []common_interfaces.ItemProperty {
	result := make([]common_interfaces.ItemProperty, 0)
	c.mtx.Lock()
	for _, prop := range c.Properties {
		result = append(result, prop)
	}
	c.mtx.Unlock()
	return result
}

func (c *Unit) Id() string {
	return c.unitId
}

func (c *Unit) SetId(id string) {
	c.unitId = id
}

func (c *Unit) SetIUnit(iUnit common_interfaces.IUnit) {
	c.iUnit = iUnit
}

func (c *Unit) SetMainItem(mainItem string) {
	c.PropSetIfNotExists("main_item", c.Id()+"/"+mainItem)
}

func (c *Unit) MainItem() string {
	return c.Prop("main_item")
}

func (c *Unit) Type() string {
	return c.unitType
}

func (c *Unit) SetType(unitType string) {
	c.unitType = unitType
}

func (c *Unit) DisplayName() string {
	return c.unitDisplayName
}

func (c *Unit) SetDisplayName(unitDisplayName string) {
	c.unitDisplayName = unitDisplayName
}

func (c *Unit) SetConfig(config string) {
	c.config = config
}

func (c *Unit) GetConfig() string {
	return c.config
}

func (c *Unit) GetConfigMeta() string {
	return ""
}

func (c *Unit) InternalInitItems() {
	c.SetStringForAll("", uom.STARTED)
}

func (c *Unit) InternalDeInitItems() {
	c.SetStringForAll("", uom.STOPPED)
}

func (c *Unit) Start(iDataStorage common_interfaces.IDataStorage) error {
	var err error
	c.watchItems = make(map[string]bool)
	c.iDataStorage = iDataStorage
	if c.Started {
		return errors.New("already started")
	}
	c.LogInfo("")
	c.LogInfo("starting ...")
	c.SetStringService("name", c.DisplayName(), "")
	c.SetError("")
	c.SetStringService("status", "started", "")
	c.SetStringService("Unit", c.Type(), "")

	c.Stopping = false

	c.iUnit.InternalInitItems()
	err = c.iUnit.InternalUnitStart()

	if err != nil {
		c.SetError(err.Error())
		c.LogError(err.Error())
	} else {
		c.LogInfo("started")
	}

	return err
}

func (c *Unit) Stop() {
	if !c.Started {
		return
	}
	c.LogInfo("stopping ...")

	for itemWatched, _ := range c.watchItems {
		c.iDataStorage.RemoveFromWatch(c.Id(), itemWatched)
	}

	c.SetStringService("status", "stopping", "")
	c.Stopping = true
	for c.Started {
		time.Sleep(100 * time.Millisecond)
	}
	c.SetStringService("status", "stopped", "")
	c.iUnit.InternalDeInitItems()
	c.LogInfo("stopped")
}

func (c *Unit) IsStarted() bool {
	return c.Started
}

const (
	UnitServicePrefix = ".service/"
	ItemNameError     = "error"
)

func (c *Unit) IDataStorage() common_interfaces.IDataStorage {
	return c.iDataStorage
}

func (c *Unit) SetStringService(name string, value string, UOM string) {
	fullName := c.Id() + "/" + UnitServicePrefix + name
	c.iDataStorage.SetItemByName(fullName, value, UOM, time.Now().UTC(), false)
}

func (c *Unit) LogInfo(value string) {
	dt := time.Now().UTC()
	if dt.Sub(c.lastLogDT) < 1*time.Microsecond {
		dt = dt.Add(1 * time.Microsecond)
	}
	c.lastLogDT = dt
	if c.lastInfo != value || time.Now().UTC().Sub(c.lastInfoDT) > 5*time.Second {
		fullName := c.Id() + "/" + UnitServicePrefix + "log"
		c.iDataStorage.SetItemByName(fullName, value, "", dt, false)
		c.lastInfoDT = time.Now().UTC()
	}
	c.lastInfo = value
	time.Sleep(1 * time.Microsecond)
}

func (c *Unit) LogError(value string) {
	dt := time.Now().UTC()
	if dt.Sub(c.lastLogDT) < 1*time.Microsecond {
		dt = dt.Add(1 * time.Microsecond)
	}
	c.lastLogDT = dt

	if c.lastError != value || time.Now().UTC().Sub(c.lastErrorDT) > 5*time.Second {
		fullName := c.Id() + "/" + UnitServicePrefix + "log"
		c.iDataStorage.SetItemByName(fullName, value, "error", dt, false)
		c.lastErrorDT = time.Now().UTC()
	}
	c.lastError = value
	time.Sleep(1 * time.Microsecond)
}

func (c *Unit) SetError(value string) {
	fullName := c.Id() + "/" + UnitServicePrefix + ItemNameError
	c.iDataStorage.SetItemByName(fullName, value, "", time.Now().UTC(), false)
}

func (c *Unit) SetStringForAll(value string, UOM string) {
	fullName := c.Id()
	c.iDataStorage.SetAllItemsByUnitName(fullName, value, UOM, time.Now().UTC(), false)
}

func (c *Unit) SetString(name string, value string, UOM string) {
	fullName := c.Id()
	if len(name) > 0 {
		fullName = c.Id() + "/" + name
	}
	c.iDataStorage.SetItemByName(fullName, value, UOM, time.Now().UTC(), false)
}

func (c *Unit) SetPropertyIfDoesntExist(itemName string, propName string, propValue string) {
	c.iDataStorage.SetPropertyIfDoesntExist(c.Id()+"/"+itemName, propName, propValue)
}

func (c *Unit) TouchItem(name string) {
	fullName := c.Id()
	if len(name) > 0 {
		fullName = c.Id() + "/" + name
	}
	c.iDataStorage.TouchItem(fullName)
}

func (c *Unit) SetInt(name string, value int, UOM string) {
	c.SetString(name, strconv.Itoa(value), UOM)
}

func (c *Unit) SetInt64(name string, value int64, UOM string) {
	c.SetString(name, fmt.Sprint(value), UOM)
}

func (c *Unit) SetUInt64(name string, value uint64, UOM string) {
	c.SetString(name, fmt.Sprint(value), UOM)
}

func (c *Unit) SetInt32(name string, value int32, UOM string) {
	c.SetString(name, fmt.Sprint(value), UOM)
}

func (c *Unit) SetUInt32(name string, value uint32, UOM string) {
	c.SetString(name, fmt.Sprint(value), UOM)
}

func (c *Unit) SetInt16(name string, value int16, UOM string) {
	c.SetString(name, fmt.Sprint(value), UOM)
}

func (c *Unit) SetUInt16(name string, value uint16, UOM string) {
	c.SetString(name, fmt.Sprint(value), UOM)
}

func (c *Unit) SetFloat64(name string, value float64, UOM string, precision int) {
	c.SetString(name, strconv.FormatFloat(value, 'f', precision, 64), UOM)
}

func (c *Unit) GetValue(name string) string {
	item, err := c.iDataStorage.GetItem(name)
	if err != nil {
		return ""
	}
	return item.Value.Value
}

func (c *Unit) GetItem(name string) (common_interfaces.ItemValue, error) {
	item, err := c.iDataStorage.GetItem(name)
	if err != nil {
		return common_interfaces.ItemValue{}, err
	}
	return item.Value, nil
}

func (c *Unit) GetItemsOfUnit(unitId string) ([]common_interfaces.ItemGetUnitItems, error) {
	return c.iDataStorage.GetUnitValues(unitId), nil
}

func (c *Unit) AddToWatch(itemName string) {
	c.iDataStorage.AddToWatch(c.Id(), itemName)
	c.watchItems[itemName] = true
}

func (c *Unit) RemoveFromWatch(itemName string) {
	c.iDataStorage.RemoveFromWatch(c.Id(), itemName)
	delete(c.watchItems, itemName)
}

func (c *Unit) ItemChanged(itemName string, value common_interfaces.ItemValue) {
}
