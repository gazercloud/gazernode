package units_common

import (
	"errors"
	"fmt"
	"github.com/gazercloud/gazernode/common_interfaces"
	"strconv"
	"time"
)

type Unit struct {
	unitId       string
	unitType     string
	unitName     string
	config       string
	iUnit        common_interfaces.IUnit
	iDataStorage common_interfaces.IDataStorage
	mainItem     string
	lastError    string
	lastErrorDT  time.Time
	lastInfo     string
	lastInfoDT   time.Time

	lastLogDT time.Time

	Started  bool
	Stopping bool

	watchItems map[string]bool
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
	c.mainItem = mainItem
}

func (c *Unit) MainItem() string {
	return c.mainItem
}

func (c *Unit) Type() string {
	return c.unitType
}

func (c *Unit) SetType(unitType string) {
	c.unitType = unitType
}

func (c *Unit) Name() string {
	return c.unitName
}

func (c *Unit) SetName(unitName string) {
	c.unitName = unitName
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

func (c *Unit) Start(iDataStorage common_interfaces.IDataStorage) error {
	var err error
	c.watchItems = make(map[string]bool)
	c.iDataStorage = iDataStorage
	if c.Started {
		return errors.New("already started")
	}
	c.LogInfo("")
	c.LogInfo("starting ...")
	c.SetStringService("name", c.Name(), "")
	c.SetError("")
	c.SetStringService("status", "started", "")
	c.SetStringService("Unit", c.Type(), "")

	c.Stopping = false
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
	c.iUnit.InternalUnitStop()
	for c.Started {
		time.Sleep(100 * time.Millisecond)
	}
	c.SetStringService("status", "stopped", "")
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
	fullName := c.Name() + "/" + UnitServicePrefix + name
	c.iDataStorage.SetItemByName(fullName, value, UOM, time.Now().UTC(), false)
}

func (c *Unit) LogInfo(value string) {
	dt := time.Now().UTC()
	if dt.Sub(c.lastLogDT) < 1*time.Microsecond {
		dt = dt.Add(1 * time.Microsecond)
	}
	c.lastLogDT = dt
	if c.lastInfo != value || time.Now().UTC().Sub(c.lastInfoDT) > 5*time.Second {
		fullName := c.Name() + "/" + UnitServicePrefix + "log"
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
		fullName := c.Name() + "/" + UnitServicePrefix + "log"
		c.iDataStorage.SetItemByName(fullName, value, "error", dt, false)
		c.lastErrorDT = time.Now().UTC()
	}
	c.lastError = value
	time.Sleep(1 * time.Microsecond)
}

func (c *Unit) SetError(value string) {
	fullName := c.Name() + "/" + UnitServicePrefix + ItemNameError
	c.iDataStorage.SetItemByName(fullName, value, "", time.Now().UTC(), false)
}

func (c *Unit) SetString(name string, value string, UOM string) {
	fullName := c.Name()
	if len(name) > 0 {
		fullName = c.Name() + "/" + name
	}
	c.iDataStorage.SetItemByName(fullName, value, UOM, time.Now().UTC(), false)
}

func (c *Unit) SetPropertyIfDoesntExist(itemName string, propName string, propValue string) {
	c.iDataStorage.SetPropertyIfDoesntExist(itemName, propName, propValue)
}

func (c *Unit) TouchItem(name string) {
	fullName := c.Name()
	if len(name) > 0 {
		fullName = c.Name() + "/" + name
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
