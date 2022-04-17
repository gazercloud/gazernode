package units_system

import (
	"errors"
	"fmt"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/resources"
	"github.com/gazercloud/gazernode/system/protocols/nodeinterface"
	"github.com/gazercloud/gazernode/system/units/databases/unit_postgreesql"
	"github.com/gazercloud/gazernode/system/units/files/unit_filecontent"
	"github.com/gazercloud/gazernode/system/units/files/unit_filesize"
	"github.com/gazercloud/gazernode/system/units/gazer/unit_gazer_cloud"
	"github.com/gazercloud/gazernode/system/units/general/unit_general_cgi"
	"github.com/gazercloud/gazernode/system/units/general/unit_general_cgi_key_value"
	"github.com/gazercloud/gazernode/system/units/general/unit_hhgttg"
	"github.com/gazercloud/gazernode/system/units/general/unit_manual"
	"github.com/gazercloud/gazernode/system/units/general/unit_signal_generator"
	"github.com/gazercloud/gazernode/system/units/network/unit_http_json_requester"
	"github.com/gazercloud/gazernode/system/units/network/unit_ping"
	"github.com/gazercloud/gazernode/system/units/network/unit_ssl"
	"github.com/gazercloud/gazernode/system/units/network/unit_tcp_connect"
	unit_tcp_telnet_control "github.com/gazercloud/gazernode/system/units/network/unit_tcp_control"
	"github.com/gazercloud/gazernode/system/units/raspberry_pi/unit_raspberry_pi_cpu_temp"
	"github.com/gazercloud/gazernode/system/units/raspberry_pi/unit_raspberry_pi_gpio"
	unit_serial_port_key_value "github.com/gazercloud/gazernode/system/units/serial_port/serial_port_key_value"
	"github.com/gazercloud/gazernode/system/units/units_common"
	"github.com/gazercloud/gazernode/system/units/windows/unit_network"
	"github.com/gazercloud/gazernode/system/units/windows/unit_process"
	"github.com/gazercloud/gazernode/system/units/windows/unit_processes"
	"github.com/gazercloud/gazernode/system/units/windows/unit_storage"
	"github.com/gazercloud/gazernode/system/units/windows/unit_system_memory"
	"github.com/gazercloud/gazernode/utilities/logger"
	"runtime"
	"sync"
	"time"
)

type UnitsSystem struct {
	units        []common_interfaces.IUnit
	unitTypes    []*UnitType
	unitTypesMap map[string]*UnitType
	iDataStorage common_interfaces.IDataStorage
	mtx          sync.Mutex
}

var unitCategoriesIcons map[string][]byte
var unitCategoriesNames map[string]string

func init() {
	unitCategoriesIcons = make(map[string][]byte)
	unitCategoriesIcons["network"] = resources.R_files_sensors_category_network_png
	unitCategoriesIcons["computer"] = resources.R_files_sensors_category_computer_png
	unitCategoriesIcons["file"] = resources.R_files_sensors_category_file_png
	unitCategoriesIcons["general"] = resources.R_files_sensors_category_general_png

	unitCategoriesIcons["serial_port"] = resources.R_files_sensors_category_serial_port_png
	unitCategoriesIcons["raspberry_pi"] = resources.R_files_sensors_category_raspberry_pi_png
	unitCategoriesIcons["database"] = resources.R_files_sensors_category_database_png
	unitCategoriesIcons["gazer"] = resources.R_files_sensors_category_gazer_png
	unitCategoriesIcons[""] = resources.R_files_sensors_category_all_png

	unitCategoriesNames = make(map[string]string)
	unitCategoriesNames["network"] = "Network"
	unitCategoriesNames["computer"] = "Computer"
	unitCategoriesNames["file"] = "File"
	unitCategoriesNames["general"] = "General"
	unitCategoriesNames["serial_port"] = "Serial Port"
	unitCategoriesNames["raspberry_pi"] = "RaspberryPI"
	unitCategoriesNames["database"] = "Database"
	unitCategoriesNames["gazer"] = "Gazer"
	unitCategoriesNames[""] = "All"
}

func New(iDataStorage common_interfaces.IDataStorage) *UnitsSystem {
	var c UnitsSystem
	c.unitTypes = make([]*UnitType, 0)
	c.unitTypesMap = make(map[string]*UnitType)
	c.iDataStorage = iDataStorage

	var unitType *UnitType

	unitType = c.RegisterUnit("network_ping", "network", "Ping", unit_ping.New, unit_ping.Image, "")
	unitType.Help = "https://gazer.cloud/unit-types/network/ping/"
	unitType = c.RegisterUnit("network_tcp_connect", "network", "TCP Connect", unit_tcp_connect.New, unit_tcp_connect.Image, "")
	unitType.Help = "https://gazer.cloud/unit-types/network/tcp-connect/"

	unitType = c.RegisterUnit("network_http_json_requester", "network", "JSON Requester", unit_http_json_requester.New, unit_http_json_requester.Image, "")
	unitType.Help = "https://gazer.cloud/unit-types/network/json-requester/"

	unitType = c.RegisterUnit("network_tcp_telnet_control", "network", "TCP Telnet Control", unit_tcp_telnet_control.New, unit_tcp_telnet_control.Image, "")
	unitType.Help = "https://gazer.cloud/unit-types/network/tcp-telnet-control/"

	unitType = c.RegisterUnit("network_ssl", "network", "SSL", unit_ssl.New, unit_ssl.Image, "")
	unitType.Help = "https://gazer.cloud/unit-types/network/ssl/"

	//unitType = c.RegisterUnit("network_http_json_items_server", "network", "HTTP Json Items Server", unit_http_json_items_server.New, unit_http_json_items_server.Image, "")
	//unitType = c.RegisterUnit("network_http_json_units_server", "network", "HTTP Json Units Server", unit_http_json_units_server.New, unit_http_json_units_server.Image, "")

	unitType = c.RegisterUnit("computer_memory", "computer", "Memory", unit_system_memory.New, unit_system_memory.Image, "")
	unitType.Help = "https://gazer.cloud/unit-types/computer/memory/"
	unitType = c.RegisterUnit("computer_process", "computer", "Process", unit_process.New, unit_process.Image, "")
	unitType.Help = "https://gazer.cloud/unit-types/computer/process/"
	unitType = c.RegisterUnit("computer_processes", "computer", "Processes", unit_processes.New, unit_processes.Image, "")
	unitType.Help = "https://gazer.cloud/unit-types/computer/processes/"
	unitType = c.RegisterUnit("computer_storage", "computer", "Storage", unit_storage.New, unit_storage.Image, "")
	unitType.Help = "https://gazer.cloud/unit-types/computer/storage/"
	unitType = c.RegisterUnit("computer_network", "computer", "Network", unit_network.New, unit_network.Image, "")
	unitType.Help = "https://gazer.cloud/unit-types/computer/network/"
	/*unitType = c.RegisterUnit("computer_named_pipe_server", "computer", "Named Pipe Server", unit_system_named_pipe_server.New, unit_system_named_pipe_server.Image, "")
	unitType.Help = "https://gazer.cloud/unit-types/computer/memory/"
	*/
	//unitType = c.RegisterUnit("computer_network_interface", "computer", "Network Interface", unit_network_interface.New, unit_network_interface.Image, "")
	//unitType.Help = "https://gazer.cloud/unit-types/computer/network-interface/"

	unitType = c.RegisterUnit("file_size", "file", "File Size", unit_filesize.New, unit_filesize.Image, "")
	unitType.Help = "https://gazer.cloud/unit-types/files/file-size/"
	unitType = c.RegisterUnit("file_content", "file", "File Content", unit_filecontent.New, unit_filecontent.Image, "")
	unitType.Help = "https://gazer.cloud/unit-types/files/file-content/"
	//unitType = c.RegisterUnit("file_csv_export", "file", "CSV Export", unit_csv_export.New, unit_csv_export.Image, "")

	unitType = c.RegisterUnit("general_cgi", "general", "Console", unit_general_cgi.New, unit_general_cgi.Image, "")
	unitType.Help = "https://gazer.cloud/unit-types/general/cgi/"
	unitType = c.RegisterUnit("general_cgi_key_value", "general", "Console Key=Value", unit_general_cgi_key_value.New, unit_general_cgi_key_value.Image, "")
	unitType.Help = "https://gazer.cloud/unit-types/general/cgi-key-value/"
	unitType = c.RegisterUnit("general_manual", "general", "Manual Items", unit_manual.New, unit_manual.Image, "")
	unitType.Help = "https://gazer.cloud/unit-types/general/manual-items/"
	unitType = c.RegisterUnit("general_hhgttg", "general", "HHGTTG", unit_hhgttg.New, unit_hhgttg.Image, "")
	unitType.Help = "https://gazer.cloud/unit-types/general/hhgttg/"
	unitType = c.RegisterUnit("general_signal_generator", "general", "Signal Generator", unit_signal_generator.New, unit_signal_generator.Image, "")
	unitType.Help = "https://gazer.cloud/unit-types/general/signal-generator/"

	//unitType = c.RegisterUnit("general_calculator", "general", "Calculator", unit_calculator.New, unit_calculator.Image, "")
	//unitType.Help = "https://gazer.cloud/unit-types/"

	unitType = c.RegisterUnit("serial_port_key_value", "serial_port", "Serial Port Key=Value", unit_serial_port_key_value.New, unit_serial_port_key_value.Image, "Key/value unit via Serial Port. Format: key=value<new_line>")
	unitType.Help = "https://gazer.cloud/unit-types/serial-port/serial-port-key-value/"

	if runtime.GOOS == "linux" {
		unitType = c.RegisterUnit("raspberry_pi_gpio", "raspberry_pi", "Raspberry PI GPIO", unit_raspberry_pi_gpio.New, unit_raspberry_pi_gpio.Image, "RaspberryPI GPIO")
		unitType.Help = ""
		unitType = c.RegisterUnit("raspberry_pi_cpu_temp", "raspberry_pi", "Raspberry PI CPU temperature", unit_raspberry_pi_cpu_temp.New, unit_raspberry_pi_cpu_temp.Image, "RaspberryPI CPU Temperature")
		unitType.Help = "https://gazer.cloud/unit-types/raspberrypi/cpu-temperature/"
	}

	unitType = c.RegisterUnit("database_postgresql", "database", "PostgreSQL", unit_postgreesql.New, unit_postgreesql.Image, "PostgreSQL database query execute")
	unitType.Help = "https://gazer.cloud/unit-types/databases/postgresql/"

	/*unitType = c.RegisterUnit("database_mysql", "database", "MySQL", unit_mysql.New, unit_mysql.Image, "MySQL database query execute")
		unitType.Help = `
	No description available
	`
	*/
	unitType = c.RegisterUnit("gazer_cloud", "gazer", "Gazer Cloud", unit_gazer_cloud.New, unit_gazer_cloud.Image, "Gazer Cloud Monitoring")
	unitType.Help = ""

	//unitType = c.RegisterUnit("industrial_modbus", "industrial", "Modbus TCP", unit_modbus.New, unit_modbus.Image, "Modbus TCP")

	return &c
}

func (c *UnitsSystem) RegisterUnit(typeName string, category string, displayName string, constructor func() common_interfaces.IUnit, imgBytes []byte, description string) *UnitType {
	var sType UnitType
	sType.TypeCode = typeName
	sType.Category = category
	sType.DisplayName = displayName
	sType.Constructor = constructor
	sType.Picture = imgBytes
	sType.Description = description
	c.unitTypes = append(c.unitTypes, &sType)
	c.unitTypesMap[typeName] = &sType
	return &sType
}

func (c *UnitsSystem) UnitTypes() []nodeinterface.UnitTypeListResponseItem {
	result := make([]nodeinterface.UnitTypeListResponseItem, 0)
	for _, st := range c.unitTypes {
		var unitTypeInfo nodeinterface.UnitTypeListResponseItem
		unitTypeInfo.Type = st.TypeCode
		unitTypeInfo.Category = st.Category
		unitTypeInfo.DisplayName = st.DisplayName
		unitTypeInfo.Help = st.Help
		unitTypeInfo.Description = st.Description
		unitTypeInfo.Image = st.Picture
		result = append(result, unitTypeInfo)
	}
	return result
}

func (c *UnitsSystem) UnitCategories() nodeinterface.UnitTypeCategoriesResponse {
	var result nodeinterface.UnitTypeCategoriesResponse
	result.Items = make([]nodeinterface.UnitTypeCategoriesResponseItem, 0)
	addedCats := make(map[string]bool)

	catAllName := ""
	var unitCategoryInfoAll nodeinterface.UnitTypeCategoriesResponseItem
	unitCategoryInfoAll.Name = catAllName
	unitCategoryInfoAll.DisplayName = "All"
	if imgBytes, ok := unitCategoriesIcons[catAllName]; ok {
		unitCategoryInfoAll.Image = imgBytes
	}
	result.Items = append(result.Items, unitCategoryInfoAll)
	addedCats[catAllName] = true

	for _, st := range c.unitTypes {
		if _, ok := addedCats[st.Category]; !ok {
			var unitCategoryInfo nodeinterface.UnitTypeCategoriesResponseItem
			unitCategoryInfo.Name = st.Category
			if catName, ok := unitCategoriesNames[st.Category]; ok {
				unitCategoryInfo.DisplayName = catName
			} else {
				unitCategoryInfo.DisplayName = st.Category
			}
			if imgBytes, ok := unitCategoriesIcons[st.Category]; ok {
				unitCategoryInfo.Image = imgBytes
			} else {
				unitCategoryInfo.Image = st.Picture
			}
			result.Items = append(result.Items, unitCategoryInfo)
			addedCats[st.Category] = true
		}
	}
	return result
}

func (c *UnitsSystem) UnitTypeForDisplayByType(t string) string {
	if res, ok := c.unitTypesMap[t]; ok {
		return res.DisplayName
	}
	return t
}

func (c *UnitsSystem) Start() {
	for _, unit := range c.units {
		unit.Start(c.iDataStorage)
	}
}

func (c *UnitsSystem) Stop() {
	logger.Println("UNITS_SYSTEM stopping begin")
	for _, unit := range c.units {
		go unit.Stop()
	}

	time.Sleep(100 * time.Millisecond) // Wait for units
	startedUnits := make([]string, 0)

	regularQuit := false
	for i := 0; i < 10; i++ {
		startedUnits = make([]string, 0)
		for _, unit := range c.units {
			if unit.IsStarted() {
				startedUnits = append(startedUnits, fmt.Sprint(unit.Id(), " / ", unit.Type(), " / ", unit.DisplayName()))
			}
		}
		if len(startedUnits) == 0 {
			regularQuit = true
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	if !regularQuit {
		logger.Println("Units stopping - timeout")
		for _, startedUnit := range startedUnits {
			logger.Println("Started: ", startedUnit)
		}
	}
	logger.Println("UNITS_SYSTEM stopping end")
}

func (c *UnitsSystem) MakeUnitByType(unitType string, dataStorage common_interfaces.IDataStorage) common_interfaces.IUnit {
	var unit common_interfaces.IUnit

	for _, st := range c.unitTypes {
		if st.TypeCode == unitType {
			unit = st.Constructor()
			break
		}
	}

	if unit != nil {
		unit.Init()
	}

	return unit
}

func (c *UnitsSystem) AddUnit(unitType string, unitId string, displayName string, config string, fromCloud bool) (common_interfaces.IUnit, error) {
	var unit common_interfaces.IUnit
	nameIsExists := false
	c.mtx.Lock()
	for _, s := range c.units {
		if s.DisplayName() == displayName {
			nameIsExists = true
		}
	}
	c.mtx.Unlock()

	if fromCloud {
		if unitType == "general_cgi" || unitType == "general_cgi_key_value" {
			return nil, errors.New("cannot edit a cgi-unit via the Cloud")
		}
	}

	if !nameIsExists {
		unit = c.MakeUnitByType(unitType, c.iDataStorage)
		if unit != nil {
			unit.SetId(unitId)
			unit.SetDisplayName(displayName)
			unit.SetType(unitType)
			unit.SetConfig(config)
			unit.SetIUnit(unit)
			c.units = append(c.units, unit)
		} else {
			return nil, errors.New("cannot create unit")
		}
	} else {
		return nil, errors.New("unit already exists")
	}
	return unit, nil
}

func (c *UnitsSystem) GetUnitState(unitId string) (nodeinterface.UnitStateResponse, error) {
	var unit common_interfaces.IUnit
	c.mtx.Lock()
	for _, s := range c.units {
		if s.Id() == unitId {
			unit = s
		}
	}
	c.mtx.Unlock()

	if unit != nil {
		var unitState nodeinterface.UnitStateResponse
		unitState.Status = ""
		unitState.UnitDisplayName = unit.DisplayName()
		unitState.MainItem = unit.MainItem()
		unitState.Type = unit.Type()
		unitState.TypeName = c.UnitTypeForDisplayByType(unit.Type())
		if unit.IsStarted() {
			unitState.Status = "started"
		} else {
			unitState.Status = "stopped"
		}

		return unitState, nil
	}
	return nodeinterface.UnitStateResponse{}, errors.New("no unit found")
}

func (c *UnitsSystem) GetUnitStateAll() (nodeinterface.UnitStateAllResponse, error) {
	var result nodeinterface.UnitStateAllResponse
	result.Items = make([]nodeinterface.UnitStateAllResponseItem, 0)

	c.mtx.Lock()
	for _, unit := range c.units {
		var unitState nodeinterface.UnitStateAllResponseItem
		unitState.Status = ""
		unitState.UnitId = unit.Id()
		unitState.UnitDisplayName = unit.DisplayName()
		unitState.Type = unit.Type()
		unitState.TypeName = c.UnitTypeForDisplayByType(unit.Type())
		unitState.MainItem = unit.MainItem()
		if unit.IsStarted() {
			unitState.Status = "started"
		} else {
			unitState.Status = "stopped"
		}
		result.Items = append(result.Items, unitState)
	}
	c.mtx.Unlock()

	return result, nil
}

func (c *UnitsSystem) ListOfUnits() nodeinterface.UnitListResponse {
	var result nodeinterface.UnitListResponse
	result.Items = make([]nodeinterface.UnitListResponseItem, 0)
	c.mtx.Lock()
	defer c.mtx.Unlock()
	for _, s := range c.units {
		var sens nodeinterface.UnitListResponseItem
		sens.Id = s.Id()
		sens.Type = s.Type()
		sens.DisplayName = s.DisplayName()
		sens.Enable = s.IsStarted()
		sens.TypeForDisplay = c.UnitTypeForDisplayByType(s.Type())
		sens.Config = s.GetConfig()
		result.Items = append(result.Items, sens)
	}
	return result
}

func (c *UnitsSystem) Units() []units_common.UnitInfo {
	var result []units_common.UnitInfo
	result = make([]units_common.UnitInfo, 0)
	c.mtx.Lock()
	defer c.mtx.Unlock()
	for _, s := range c.units {
		var sens units_common.UnitInfo
		sens.Id = s.Id()
		sens.Type = s.Type()
		sens.DisplayName = s.DisplayName()
		sens.Enable = s.IsStarted()
		sens.TypeForDisplay = c.UnitTypeForDisplayByType(s.Type())
		sens.Config = s.GetConfig()
		sens.Properties = s.PropGet()
		result = append(result, sens)
	}
	return result
}

func (c *UnitsSystem) StartUnit(unitId string) error {
	for _, s := range c.units {
		if s.Id() == unitId {
			s.Start(c.iDataStorage)
		}
	}
	return nil
}

func (c *UnitsSystem) StopUnit(unitId string) error {
	for _, s := range c.units {
		if s.Id() == unitId {
			s.Stop()
		}
	}
	return nil
}

func (c *UnitsSystem) RemoveUnits(units []string) error {
	logger.Println("UnitsSystem RemoveUnits", units)
	c.mtx.Lock()

	var deletedUnit common_interfaces.IUnit
	var unitIndex int
	idsOfDeletedUnits := make([]string, 0)

	for _, unitToRemove := range units {
		for unitIndex, deletedUnit = range c.units {
			if deletedUnit.Id() == unitToRemove {
				logger.Println("UnitsSystem RemoveUnits unit", deletedUnit.Id())
				idsOfDeletedUnits = append(idsOfDeletedUnits, deletedUnit.Id())
				deletedUnit.Stop()
				c.units = append(c.units[:unitIndex], c.units[unitIndex+1:]...)
				break
			}
		}
	}

	c.mtx.Unlock()

	for _, idOfDeletedUnit := range idsOfDeletedUnits {
		_ = c.iDataStorage.RemoveItemsOfUnit(idOfDeletedUnit)
	}

	return nil
}

func (c *UnitsSystem) GetUnitDisplayName(unitId string) (string, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	for _, s := range c.units {
		if s.Id() == unitId {
			return s.DisplayName(), nil
		}
	}
	return "", errors.New("no unit found")
}

func (c *UnitsSystem) GetConfig(unitId string) (string, string, string, string, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	for _, s := range c.units {
		if s.Id() == unitId {
			return s.DisplayName(), s.GetConfig(), s.GetConfigMeta(), s.Type(), nil
		}
	}
	return "", "", "", "", errors.New("no unit found")
}

func (c *UnitsSystem) GetConfigByType(unitType string) (string, string, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	for _, st := range c.unitTypes {
		if st.TypeCode == unitType {
			sens := c.MakeUnitByType(st.TypeCode, nil)
			if sens != nil {
				return st.DisplayName, sens.GetConfigMeta(), nil
			} else {
				return "", "", errors.New("no unit type found")
			}
		}
	}

	return "", "", errors.New("no unit type found")
}

func (c *UnitsSystem) SetConfig(unitId string, name string, config string, fromCloud bool) error {
	var unit common_interfaces.IUnit

	c.mtx.Lock()
	for _, s := range c.units {
		if s.Id() == unitId {
			unit = s
		}
	}
	c.mtx.Unlock()

	if unit != nil {
		if fromCloud {
			if unit.Type() == "general_cgi" || unit.Type() == "general_cgi_key_value" {
				return errors.New("cannot edit a cgi-unit via the Cloud")
			}
		}
	}

	if unit != nil {

		unit.Stop()
		oldName := unit.DisplayName()

		if oldName != name {

			nameIsExists := false
			c.mtx.Lock()
			for _, s := range c.units {
				if s.DisplayName() == name {
					nameIsExists = true
				}
			}
			c.mtx.Unlock()

			if !nameIsExists {
				unit.SetDisplayName(name)
				//c.iDataStorage.RenameItems(oldName+"/", name+"/")
			}
		}

		unit.SetConfig(config)

		unit.Start(c.iDataStorage)
	}

	return nil
}

func (c *UnitsSystem) SendToWatcher(unitId string, itemName string, value common_interfaces.ItemValue) {
	var targetUnit common_interfaces.IUnit

	c.mtx.Lock()
	for _, unit := range c.units {
		if unit.Id() == unitId {
			targetUnit = unit
			break
		}
	}
	c.mtx.Unlock()

	if targetUnit != nil {
		targetUnit.ItemChanged(itemName, value)
	}

}

func (c *UnitsSystem) UnitPropSet(unitId string, props []nodeinterface.PropItem) error {
	var err error
	var targetUnit common_interfaces.IUnit
	c.mtx.Lock()
	for _, unit := range c.units {
		if unit.Id() == unitId {
			targetUnit = unit
			break
		}
	}
	if targetUnit != nil {
		properties := make([]common_interfaces.ItemProperty, 0)
		for _, prop := range props {
			properties = append(properties, common_interfaces.ItemProperty{
				Name:  prop.PropName,
				Value: prop.PropValue,
			})
		}
		targetUnit.PropSet(properties)
	} else {
		err = errors.New("no unit found")
	}
	c.mtx.Unlock()
	return err
}

func (c *UnitsSystem) UnitPropGet(unitId string) ([]nodeinterface.PropItem, error) {
	var err error
	result := make([]nodeinterface.PropItem, 0)
	var targetUnit common_interfaces.IUnit
	c.mtx.Lock()
	for _, unit := range c.units {
		if unit.Id() == unitId {
			targetUnit = unit
			break
		}
	}
	if targetUnit != nil {
		props := targetUnit.PropGet()
		for _, prop := range props {
			result = append(result, nodeinterface.PropItem{
				PropName:  prop.Name,
				PropValue: prop.Value,
			})
		}
	} else {
		err = errors.New("no unit found")
	}
	c.mtx.Unlock()
	return result, err
}
