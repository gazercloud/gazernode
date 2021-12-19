package system

import (
	"fmt"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/system/protocols/lookup"
	nodeinterface2 "github.com/gazercloud/gazernode/system/protocols/nodeinterface"
	"github.com/gazercloud/gazernode/system/units/windows/unit_process"
	"github.com/gazercloud/gazernode/utilities/logger"
	"go.bug.st/serial"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

func SplitWithoutEmpty(req string, sep rune) []string {
	return strings.FieldsFunc(req, func(r rune) bool {
		return r == sep
	})
}

func (c *System) UnitTypes(category string, filter string, offset int, maxCount int) nodeinterface2.UnitTypeListResponse {
	unitTypes := c.unitsSystem.UnitTypes()

	var result nodeinterface2.UnitTypeListResponse
	result.TotalCount = len(unitTypes)
	result.Types = make([]nodeinterface2.UnitTypeListResponseItem, 0)
	filterParts := SplitWithoutEmpty(strings.ToLower(filter), ' ')

	for _, sType := range unitTypes {
		inFilter := 0
		sTypeString := sType.Type + sType.DisplayName + sType.Description + sType.Category
		sTypeString = strings.ToLower(sTypeString)
		for _, filterPart := range filterParts {
			if strings.Contains(sTypeString, filterPart) {
				inFilter++
			}
		}
		if inFilter == len(filterParts) && (category == "" || category == sType.Category) {
			if result.InFilterCount >= offset && len(result.Types) < maxCount {
				result.Types = append(result.Types, sType)
			}
			result.InFilterCount++
		}
	}

	return result
}

func (c *System) UnitCategories() nodeinterface2.UnitTypeCategoriesResponse {
	return c.unitsSystem.UnitCategories()
}

func (c *System) AddUnit(unitName string, unitType string, config string, fromCloud bool) (string, error) {
	logger.Println("System - AddUnit - ", unitName, unitType)
	unitId := strconv.FormatInt(time.Now().UnixNano(), 16) + "_" + strconv.FormatInt(int64(rand.Int()), 16)
	unit, err := c.unitsSystem.AddUnit(unitType, unitId, unitName, config, fromCloud)
	if err != nil {
		return "", err
	}
	err = unit.Start(c)
	if err != nil {
		return "", err
	}
	err = c.SaveConfig()
	if err != nil {
		return "", err
	}
	return unitId, err
}

func (c *System) GetUnitState(unitId string) (nodeinterface2.UnitStateResponse, error) {
	unitState, err := c.unitsSystem.GetUnitState(unitId)
	if err != nil {
		return nodeinterface2.UnitStateResponse{UnitId: unitId, UOM: "error"}, err
	}
	unitState.UnitId = unitId
	c.mtx.Lock()
	if item, ok := c.itemsByName[unitState.MainItem]; ok {
		unitState.Value = item.Value.Value
		unitState.UOM = item.Value.UOM
	}

	unitState.Items = make([]common_interfaces.ItemGetUnitItems, 0)
	for _, item := range c.items {
		if strings.HasPrefix(item.Name, unitState.UnitName+"/") {
			var i common_interfaces.ItemGetUnitItems
			i.Name = item.Name
			i.Value = item.Value
			i.Value.DT = item.Value.DT
			i.Value.UOM = item.Value.UOM
			unitState.Items = append(unitState.Items, i)
		}
	}

	c.mtx.Unlock()

	return unitState, err
}

func (c *System) GetUnitStateAll() (nodeinterface2.UnitStateAllResponse, error) {
	var err error
	var result nodeinterface2.UnitStateAllResponse

	result, err = c.unitsSystem.GetUnitStateAll()
	if err != nil {
		return result, err
	}

	c.mtx.Lock()
	for i := range result.Items {
		if item, ok := c.itemsByName[result.Items[i].MainItem]; ok {
			result.Items[i].Value = item.Value.Value
			result.Items[i].UOM = item.Value.UOM
		}
	}
	c.mtx.Unlock()

	return result, err
}

func (c *System) GetConfig(unitId string) (string, string, string, string, error) {
	return c.unitsSystem.GetConfig(unitId)
}

func (c *System) GetConfigByType(unitType string) (string, string, error) {
	return c.unitsSystem.GetConfigByType(unitType)
}

func (c *System) SetConfig(unitId string, name string, config string, fromCloud bool) error {
	err := c.unitsSystem.SetConfig(unitId, name, config, fromCloud)
	//logger.Println("system - SetConfig:", unitId, "name:", name, "error:", err)
	if err != nil {
		return err
	}
	err = c.SaveConfig()
	//logger.Println("system - SetConfig - save config:", unitId, "name:", name, "error:", err)
	return err
}

func (c *System) RemoveUnits(units []string) error {
	logger.Println("System RemoveUnits", units)
	err := c.unitsSystem.RemoveUnits(units)
	if err != nil {
		return err
	}
	err = c.SaveConfig()
	return err
}

func (c *System) StartUnits(ids []string) error {
	var err error
	for _, unit := range ids {
		_ = c.unitsSystem.StartUnit(unit)
	}
	err = c.SaveConfig()
	return err
}

func (c *System) StopUnits(ids []string) error {
	var err error
	for _, unit := range ids {
		_ = c.unitsSystem.StopUnit(unit)
	}
	err = c.SaveConfig()
	return err
}

func (c *System) ListOfUnits() nodeinterface2.UnitListResponse {
	return c.unitsSystem.ListOfUnits()
}

func (c *System) GetUnitValues(unitName string) []common_interfaces.ItemGetUnitItems {
	var items []common_interfaces.ItemGetUnitItems
	items = make([]common_interfaces.ItemGetUnitItems, 0)

	c.mtx.Lock()

	for _, item := range c.items {
		if strings.HasPrefix(item.Name, unitName+"/") {
			var i common_interfaces.ItemGetUnitItems
			i.Name = item.Name
			i.Value = item.Value
			i.Value.DT = item.Value.DT
			i.Value.UOM = item.Value.UOM
			//i.Value.Flags = item.Value.Flags
			//i.CloudChannels = c.publicChannels.GetChannelsWithItem(item.Name)
			//i.CloudChannelsNames = c.publicChannels.GetChannelsNamesWithItem(item.Name)
			items = append(items, i)
		}
	}
	c.mtx.Unlock()

	return items
}

func (c *System) RemoveItemsOfUnit(unitName string) error {
	c.mtx.Lock()
	itemsToRemove := make([]string, 0)
	for _, item := range c.items {
		if strings.HasPrefix(item.Name, unitName+"/") {
			itemsToRemove = append(itemsToRemove, item.Name)
		}
	}
	c.mtx.Unlock()

	_ = c.RemoveItems(itemsToRemove)

	return nil
}

func (c *System) GetItemsValues(reqItems []string) []common_interfaces.ItemGetUnitItems {
	var items []common_interfaces.ItemGetUnitItems
	items = make([]common_interfaces.ItemGetUnitItems, 0)

	c.mtx.Lock()
	for _, itemName := range reqItems {
		if item, ok := c.itemsByName[itemName]; ok {
			var i common_interfaces.ItemGetUnitItems
			i.Id = item.Id
			i.Name = item.Name
			i.Value = item.Value
			i.Value.DT = item.Value.DT
			i.Value.UOM = item.Value.UOM
			//i.Value.Flags = item.Value.Flags
			//i.CloudChannels = c.publicChannels.GetChannelsWithItem(item.Name)
			items = append(items, i)
		}
	}
	c.mtx.Unlock()

	return items
}

func (c *System) GetAllItems() []common_interfaces.ItemGetUnitItems {
	var items []common_interfaces.ItemGetUnitItems
	items = make([]common_interfaces.ItemGetUnitItems, 0)

	c.mtx.Lock()

	for _, item := range c.items {
		var i common_interfaces.ItemGetUnitItems
		i.Id = item.Id
		i.Name = item.Name
		i.Value = item.Value
		i.Value.DT = item.Value.DT
		i.Value.UOM = item.Value.UOM
		//i.Value.Flags = item.Value.Flags
		//i.CloudChannels = c.publicChannels.GetChannelsWithItem(item.Name)
		items = append(items, i)
	}
	c.mtx.Unlock()

	return items
}

func (c *System) Lookup(entity string) (lookup.Result, error) {
	var result lookup.Result
	result.Columns = make([]lookup.ResultColumn, 0)
	result.Rows = make([]lookup.ResultRow, 0)
	result.Entity = entity
	if entity == "serial-ports" {
		result.KeyColumn = "port"
		result.AddColumn("port", "Port")
		ports, err := serial.GetPortsList()
		if err == nil {
			for _, name := range ports {
				result.AddRow1(name)
			}
		}
	}
	if entity == "processes" {
		result.KeyColumn = "name"
		result.AddColumn("name", "Process Name")
		result.AddColumn("id", "Process Id")
		processes := unit_process.GetProcesses()
		for _, proc := range processes {
			result.AddRow2(proc.Name+"#"+fmt.Sprint(proc.Id), fmt.Sprint(proc.Id))
		}
	}
	if entity == "network_interface" {
		result.KeyColumn = "name"
		result.AddColumn("name", "Name")
		result.AddColumn("id", "Index")

		interfaces, err := net.Interfaces()
		if err == nil {
			for _, ni := range interfaces {
				result.AddRow2(ni.Name, fmt.Sprint(ni.Index))
			}
		}
	}
	if entity == "data-item" {
		result.KeyColumn = "id"
		result.AddColumn("id", "Data Item Id")
		result.AddColumn("name", "Data Item Name")
		c.mtx.Lock()
		for _, proc := range c.items {
			result.AddRow2(fmt.Sprint(proc.Id), proc.Name)
		}
		c.mtx.Unlock()
	}
	if entity == "serial-port-parity" {
		result.AddColumn("name", "Parity")
		result.AddRow1("none")
		result.AddRow1("odd")
		result.AddRow1("even")
		result.AddRow1("mark")
		result.AddRow1("space")
	}
	if entity == "serial-port-stop-bits" {
		result.AddColumn("name", "Stop Bits")
		result.AddRow1("1")
		result.AddRow1("1.5")
		result.AddRow1("2")
	}
	if entity == "gpio-mode" {
		result.AddColumn("name", "GPIO Mode")
		result.AddRow1("input")
		result.AddRow1("output")
	}
	if entity == "raspberry-pi-gpio" {
		result.AddColumn("name", "Index")
		result.AddColumn("full_name", "Full Name")
		result.AddColumn("desc", "Description")

		for i := 2; i <= 27; i++ {
			result.AddRow3(fmt.Sprint(i), "GPIO"+fmt.Sprint(i), "")
		}
	}
	if entity == "raspberry-pi-gpio-pull" {
		result.AddColumn("name", "Name")
		result.AddRow1("off")
		result.AddRow1("up")
		result.AddRow1("down")
	}
	return result, nil
}
