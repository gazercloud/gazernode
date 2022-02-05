package system

import (
	"fmt"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/system/protocols/lookup"
	"github.com/gazercloud/gazernode/system/protocols/nodeinterface"
	"github.com/gazercloud/gazernode/system/units/windows/unit_process"
	"github.com/gazercloud/gazernode/utilities/logger"
	"go.bug.st/serial"
	"net"
	"strconv"
	"strings"
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

func (c *System) AddUnit(unitName string, unitType string, config string, fromCloud bool) (string, error) {
	logger.Println("System - AddUnit - ", unitName, unitType)
	unitId := ""
	c.mtx.Lock()
	maxUnitId := uint64(0)
	for _, u := range c.unitsSystem.Units() {
		uId := u.Id
		if len(uId) > 1 && uId[0] == 'u' {
			uIdInt, uIdParseError := strconv.ParseUint(uId[1:], 10, 64)
			if uIdParseError == nil {
				if uIdInt > maxUnitId {
					maxUnitId = uIdInt
				}
			}
		}
	}
	maxUnitId++
	unitId = "u" + strconv.FormatUint(maxUnitId, 10)
	c.mtx.Unlock()

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

	unitState.Items = make([]common_interfaces.ItemGetUnitItems, 0)
	for _, item := range c.items {
		if strings.HasPrefix(item.Name, unitState.UnitId+"/") {
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

func (c *System) GetUnitStateAll() (nodeinterface.UnitStateAllResponse, error) {
	var err error
	var result nodeinterface.UnitStateAllResponse

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
			//i.CloudChannels = c.publicChannels.GetChannelsWithItem(item.Name)
			//i.CloudChannelsNames = c.publicChannels.GetChannelsNamesWithItem(item.Name)
			items = append(items, i)
		}
	}
	c.mtx.Unlock()

	return items
}

func (c *System) RemoveItemsOfUnit(unitId string) error {
	c.mtx.Lock()
	itemsToRemove := make([]string, 0)
	for _, item := range c.items {
		if strings.HasPrefix(item.Name, unitId+"/") {
			itemsToRemove = append(itemsToRemove, item.Name)
		}
	}
	c.mtx.Unlock()

	_ = c.RemoveItems(itemsToRemove)

	return nil
}

func (c *System) GetItemsValues(reqItems []string) []common_interfaces.ItemStateInfo {
	var items []common_interfaces.ItemStateInfo
	items = make([]common_interfaces.ItemStateInfo, 0)

	c.mtx.Lock()
	for _, itemName := range reqItems {
		if item, ok := c.itemsByName[itemName]; ok {
			var i common_interfaces.ItemStateInfo
			i.Id = item.Id
			i.Name = item.Name
			i.Value = item.Value.Value
			i.DT = item.Value.DT
			i.UOM = item.Value.UOM

			{
				unitId := ""
				unitName := ""
				posOfSlash := strings.Index(i.Name, "/")
				if posOfSlash > 0 {
					var err error
					unitId = i.Name[:posOfSlash]
					unitName, err = c.unitsSystem.GetUnitDisplayName(unitId)
					if err != nil {
						unitName = ""
					} else {
						i.DisplayName = strings.Replace(i.Name, unitId+"/", unitName+"/", 1)
					}
				}
			}

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

func (c *System) UnitPropSet(unitId string, props []nodeinterface.PropItem) error {
	err := c.unitsSystem.UnitPropSet(unitId, props)
	c.SaveConfig()
	return err
}

func (c *System) UnitPropGet(unitId string) ([]nodeinterface.PropItem, error) {
	res, err := c.unitsSystem.UnitPropGet(unitId)
	c.SaveConfig()
	return res, err
}

func (c *System) Lookup(entity string) (lookup.Result, error) {
	var result lookup.Result
	result.Columns = make([]lookup.ResultColumn, 0)
	result.Rows = make([]lookup.ResultRow, 0)
	result.Entity = entity
	if entity == "serial-ports" {
		result.KeyColumn = "port"
		result.AddColumn("port", "Port", false)
		ports, err := serial.GetPortsList()
		if err == nil {
			for _, name := range ports {
				result.AddRow1(name)
			}
		}
	}
	if entity == "processes" {
		result.KeyColumn = "name"
		result.AddColumn("name", "Process Name", false)
		result.AddColumn("id", "Process Id", false)
		processes := unit_process.GetProcesses()
		for _, proc := range processes {
			result.AddRow2(proc.Name+"#"+fmt.Sprint(proc.Id), fmt.Sprint(proc.Id))
		}
	}
	if entity == "network_interface" {
		result.KeyColumn = "name"
		result.AddColumn("name", "Name", false)
		result.AddColumn("id", "Index", false)

		interfaces, err := net.Interfaces()
		if err == nil {
			for _, ni := range interfaces {
				result.AddRow2(ni.Name, fmt.Sprint(ni.Index))
			}
		}
	}
	if entity == "data-item" {
		result.KeyColumn = "id"
		result.AddColumn("id", "ID", false)
		result.AddColumn("name", "Name", true)
		result.AddColumn("display_name", "Name", false)
		c.mtx.Lock()
		for _, proc := range c.items {
			if strings.Contains(proc.Name, "/.service/") {
				continue
			}
			unitId := ""
			unitName := ""
			itemDisplayName := ""
			posOfSlash := strings.Index(proc.Name, "/")
			if posOfSlash > 0 {
				var err error
				unitId = proc.Name[:posOfSlash]
				unitName, err = c.unitsSystem.GetUnitDisplayName(unitId)
				if err != nil {
					unitName = ""
				} else {
					itemDisplayName = strings.Replace(proc.Name, unitId+"/", unitName+"/", 1)
				}
			}
			result.AddRow3(fmt.Sprint(proc.Id), proc.Name, itemDisplayName)
		}
		c.mtx.Unlock()
	}
	if entity == "serial-port-parity" {
		result.AddColumn("name", "Parity", false)
		result.AddRow1("none")
		result.AddRow1("odd")
		result.AddRow1("even")
		result.AddRow1("mark")
		result.AddRow1("space")
	}
	if entity == "serial-port-stop-bits" {
		result.AddColumn("name", "Stop Bits", false)
		result.AddRow1("1")
		result.AddRow1("1.5")
		result.AddRow1("2")
	}
	if entity == "gpio-mode" {
		result.AddColumn("name", "GPIO Mode", false)
		result.AddRow1("input")
		result.AddRow1("output")
	}
	if entity == "raspberry-pi-gpio" {
		result.AddColumn("name", "Index", false)
		result.AddColumn("full_name", "Full Name", false)
		result.AddColumn("desc", "Description", false)

		for i := 2; i <= 27; i++ {
			result.AddRow3(fmt.Sprint(i), "GPIO"+fmt.Sprint(i), "")
		}
	}
	if entity == "raspberry-pi-gpio-pull" {
		result.AddColumn("name", "Name", false)
		result.AddRow1("off")
		result.AddRow1("up")
		result.AddRow1("down")
	}
	return result, nil
}
