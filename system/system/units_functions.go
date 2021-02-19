package system

import (
	"fmt"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/protocols/lookup"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
	"github.com/gazercloud/gazernode/system/units/windows/unit_process"
	"go.bug.st/serial"
	"strings"
	"time"
)

func SplitWithoutEmpty(req string, sep rune) []string {
	return strings.FieldsFunc(req, func(r rune) bool {
		return r == sep
	})
}

func (c *System) UnitTypes(category string, filter string, offset int, maxCount int) nodeinterface.UnitTypeListResponse {
	unitTypes := c.unitsSystem.UnitTypes()

	var result nodeinterface.UnitTypeListResponse
	result.TotalCount = len(unitTypes)
	result.Types = make([]nodeinterface.UnitTypeListResponseItem, 0)
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

func (c *System) UnitCategories() nodeinterface.UnitTypeCategoriesResponse {
	return c.unitsSystem.UnitCategories()
}

func (c *System) AddUnit(unitName string, unitType string) (string, error) {
	unitId := time.Now().Format("2006-01-02 15-04-05")
	err := c.unitsSystem.AddUnit(unitType, unitId, unitName, "")
	if err != nil {
		return "", err
	}
	err = c.SaveConfig()
	if err != nil {
		return "", err
	}
	return unitId, err
}

func (c *System) GetUnitState(unitId string) (nodeinterface.UnitStateResponse, error) {
	unitState, err := c.unitsSystem.GetUnitState(unitId)
	if err != nil {
		return nodeinterface.UnitStateResponse{UnitId: unitId, UOM: "error"}, err
	}
	unitState.UnitId = unitId
	c.mtx.Lock()
	if item, ok := c.itemsByName[unitState.MainItem]; ok {
		unitState.Value = item.Value.Value
		unitState.UOM = item.Value.UOM
	}
	c.mtx.Unlock()

	return unitState, err
}

func (c *System) GetConfig(unitId string) (string, string, string, string, error) {
	return c.unitsSystem.GetConfig(unitId)
}

func (c *System) GetConfigByType(unitType string) (string, string, error) {
	return c.unitsSystem.GetConfigByType(unitType)
}

func (c *System) SetConfig(unitId string, name string, config string) error {
	err := c.unitsSystem.SetConfig(unitId, name, config)
	if err != nil {
		return err
	}
	err = c.SaveConfig()
	return err
}

func (c *System) RemoveUnits(units []string) error {
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

func (c *System) ListOfUnits() nodeinterface.UnitListResponse {
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
			i.CloudChannels = c.cloud.GetChannelsWithItem(item.Name)
			i.CloudChannelsNames = c.cloud.GetChannelsNamesWithItem(item.Name)
			items = append(items, i)
		}
	}
	c.mtx.Unlock()

	return items
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
			i.CloudChannels = c.cloud.GetChannelsWithItem(item.Name)
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
		i.CloudChannels = c.cloud.GetChannelsWithItem(item.Name)
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
			result.AddRow2(proc.Name, fmt.Sprint(proc.Id))
		}
	}
	if entity == "data-items" {
		result.KeyColumn = "name"
		result.AddColumn("name", "Data Item Name")
		c.mtx.Lock()
		for _, proc := range c.items {
			result.AddRow1(proc.Name)
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
	return result, nil
}
