package units_system

import (
	"errors"
	"fmt"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/logger"
	"github.com/gazercloud/gazernode/protocols/nodeinterface"
	"github.com/gazercloud/gazernode/resources"
	"github.com/gazercloud/gazernode/system/units/databases/unit_postgreesql"
	"github.com/gazercloud/gazernode/system/units/files/unit_filecontent"
	"github.com/gazercloud/gazernode/system/units/files/unit_filesize"
	"github.com/gazercloud/gazernode/system/units/general/unit_general_cgi"
	"github.com/gazercloud/gazernode/system/units/general/unit_general_cgi_key_value"
	"github.com/gazercloud/gazernode/system/units/general/unit_hhgttg"
	"github.com/gazercloud/gazernode/system/units/general/unit_manual"
	"github.com/gazercloud/gazernode/system/units/general/unit_signal_generator"
	"github.com/gazercloud/gazernode/system/units/network/unit_ping"
	"github.com/gazercloud/gazernode/system/units/network/unit_tcp_connect"
	unit_tcp_telnet_control "github.com/gazercloud/gazernode/system/units/network/unit_tcp_control"
	unit_serial_port_key_value "github.com/gazercloud/gazernode/system/units/serial_port/serial_port_key_value"
	"github.com/gazercloud/gazernode/system/units/units_common"
	"github.com/gazercloud/gazernode/system/units/windows/unit_network"
	"github.com/gazercloud/gazernode/system/units/windows/unit_network_interface"
	"github.com/gazercloud/gazernode/system/units/windows/unit_process"
	"github.com/gazercloud/gazernode/system/units/windows/unit_storage"
	"github.com/gazercloud/gazernode/system/units/windows/unit_system_memory"
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
	unitCategoriesIcons["network"] = resources.R_files_sensors_sensor_network_png
	unitCategoriesIcons["windows"] = resources.R_files_sensors_sensor_os_png
	unitCategoriesIcons["file"] = resources.R_files_sensors_sensor_files_png
	unitCategoriesIcons["general"] = resources.R_files_sensors_sensor_general_png
	unitCategoriesIcons["serial_port"] = resources.R_files_sensors_sensor_serial_port_png
	unitCategoriesIcons["databases"] = resources.R_files_sensors_sensor_network_png
	unitCategoriesIcons[""] = resources.R_files_sensors_sensor_all_png

	unitCategoriesNames = make(map[string]string)
	unitCategoriesNames["network"] = "Network"
	unitCategoriesNames["windows"] = "Windows"
	unitCategoriesNames["file"] = "Files"
	unitCategoriesNames["general"] = "General"
	unitCategoriesNames["serial_port"] = "Serial Port"
	unitCategoriesNames["databases"] = "Databases"
	unitCategoriesNames[""] = "All"
}

func New(iDataStorage common_interfaces.IDataStorage) *UnitsSystem {
	var c UnitsSystem
	c.unitTypes = make([]*UnitType, 0)
	c.unitTypesMap = make(map[string]*UnitType)
	c.iDataStorage = iDataStorage

	var unitType *UnitType

	unitType = c.RegisterUnit("network_ping", "network", "Ping", unit_ping.New, unit_ping.Image, "")
	unitType.Help = `The sensor sends ICMP-requests to the remote host and measures the response time. 
The measured time is written to the data item "Time".
It works like standard command "ping".
You can specify period between frames and timeout in milliseconds.
Timeout: 100-10000 ms. Default value is 1000
Period: 0-999999 ms. Default value is 1000
Frame size: 4-500 bytes. Default value is 64
`
	unitType = c.RegisterUnit("network_tcp_connect", "network", "TCP Connect", unit_tcp_connect.New, unit_tcp_connect.Image, "")
	unitType.Help = `The sensor tries to connect to the specified address and writes the result to the data item.

Address - Destination host
Timeout - Timeout in milliseconds to wait for each reply.
Period - The period between sensor activities

The result time is written to the data item "Time"
`

	unitType = c.RegisterUnit("network_tcp_telnet_control", "network", "TCP Telnet Control", unit_tcp_telnet_control.New, unit_tcp_telnet_control.Image, "")
	unitType.Help = `
No description available
`
	//unitType = c.RegisterUnit("network_http_json_items_server", "network", "HTTP Json Items Server", unit_http_json_items_server.New, unit_http_json_items_server.Image, "")
	//unitType = c.RegisterUnit("network_http_json_units_server", "network", "HTTP Json Units Server", unit_http_json_units_server.New, unit_http_json_units_server.Image, "")

	unitType = c.RegisterUnit("windows_memory", "windows", "OS Memory", unit_system_memory.New, unit_system_memory.Image, "")
	unitType.Help = `
No description available
`
	unitType = c.RegisterUnit("windows_process", "windows", "OS Process", unit_process.New, unit_process.Image, "")
	unitType.Help = `
The sensor periodically gets information about the process and writes it to the corresponding data items.
`
	unitType = c.RegisterUnit("windows_storage", "windows", "OS Storage", unit_storage.New, unit_storage.Image, "")
	unitType.Help = `
No description available
`
	unitType = c.RegisterUnit("windows_network", "windows", "OS Network", unit_network.New, unit_network.Image, "")
	unitType.Help = `
No description available
`
	unitType = c.RegisterUnit("windows_network_interface", "windows", "OS Network Interface", unit_network_interface.New, unit_network_interface.Image, "")
	unitType.Help = `
No description available
`

	unitType = c.RegisterUnit("file_size", "file", "File Size", unit_filesize.New, unit_filesize.Image, "")
	unitType.Help = `
The sensor writes the file size to the data item.
`
	unitType = c.RegisterUnit("file_content", "file", "File Content", unit_filecontent.New, unit_filecontent.Image, "")
	unitType.Help = `
The sensor reads the file contents and writes it to the data item. 
The maximum size to be read is 1 kilobyte.
`
	//unitType = c.RegisterUnit("file_csv_export", "file", "CSV Export", unit_csv_export.New, unit_csv_export.Image, "")

	unitType = c.RegisterUnit("general_cgi", "general", "Console", unit_general_cgi.New, unit_general_cgi.Image, "")
	unitType.Help = `
CGI is an interface for requesting information through the command line interface. 
The sensor reads the output stream of the external program. All external program output is written to the data item. 
In fact, this sensor allows you transfer a file content to the cloud in real time.
`
	unitType = c.RegisterUnit("general_cgi_key_value", "general", "Console Key=Value", unit_general_cgi_key_value.New, unit_general_cgi_key_value.Image, "")
	unitType.Help = `
The Sensor is similar to CGI-sensor, but it can parse data by sorting it out into data elements. 
The sensor requires each data item to be written on a separate line. 
The item name is placed before the equal sign(=), and the value is placed after it.
`
	unitType = c.RegisterUnit("general_manual", "general", "Manual Items", unit_manual.New, unit_manual.Image, "")
	unitType.Help = `
No description available
`
	unitType = c.RegisterUnit("general_hhgttg", "general", "HHGTTG", unit_hhgttg.New, unit_hhgttg.Image, "")
	unitType.Help = `
Ultimate Question of Life, the Universe, and Everything
`
	unitType = c.RegisterUnit("general_signal_generator", "general", "Signal Generator", unit_signal_generator.New, unit_signal_generator.Image, "")
	unitType.Help = `
No description available
`

	unitType = c.RegisterUnit("serial_port_key_value", "serial_port", "Serial Port Key=Value", unit_serial_port_key_value.New, unit_serial_port_key_value.Image, "Key/value unit via Serial Port. Format: key=value<new_line>")
	unitType.Help = `
No description available
`

	unitType = c.RegisterUnit("databases_postgresql", "databases", "PostgreSQL", unit_postgreesql.New, unit_postgreesql.Image, "PostgreSQL database query execute")
	unitType.Help = `
No description available
`

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
				startedUnits = append(startedUnits, fmt.Sprint(unit.Id(), " / ", unit.Type(), " / ", unit.Name()))
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

	return unit
}

func (c *UnitsSystem) AddUnit(unitType string, unitId string, name string, config string) (common_interfaces.IUnit, error) {
	var unit common_interfaces.IUnit
	nameIsExists := false
	c.mtx.Lock()
	for _, s := range c.units {
		if s.Name() == name {
			nameIsExists = true
		}
	}
	c.mtx.Unlock()

	if !nameIsExists {
		unit = c.MakeUnitByType(unitType, c.iDataStorage)
		if unit != nil {
			unit.SetId(unitId)
			unit.SetName(name)
			unit.SetType(unitType)
			unit.SetConfig(config)
			unit.SetIUnit(unit)
			c.units = append(c.units, unit)
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
		unitState.UnitName = unit.Name()
		unitState.MainItem = unit.Name() + "/" + unit.MainItem()
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
		unitState.UnitName = unit.Name()
		unitState.MainItem = unit.Name() + "/" + unit.MainItem()
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
		sens.Name = s.Name()
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
		sens.Name = s.Name()
		sens.Enable = s.IsStarted()
		sens.TypeForDisplay = c.UnitTypeForDisplayByType(s.Type())
		sens.Config = s.GetConfig()
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
	c.mtx.Lock()
	defer c.mtx.Unlock()

	for _, unitToRemove := range units {
		for index, s := range c.units {
			if s.Id() == unitToRemove {
				s.Stop()
				c.units = append(c.units[:index], c.units[index+1:]...)
				_ = c.iDataStorage.RemoveItemsOfUnit(s.Name())
				break
			}
		}
	}

	return nil
}

func (c *UnitsSystem) GetUnitName(unitId string) (string, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	for _, s := range c.units {
		if s.Id() == unitId {
			return s.Name(), nil
		}
	}
	return "", errors.New("no unit found")
}

func (c *UnitsSystem) GetConfig(unitId string) (string, string, string, string, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	for _, s := range c.units {
		if s.Id() == unitId {
			return s.Name(), s.GetConfig(), s.GetConfigMeta(), s.Type(), nil
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

func (c *UnitsSystem) SetConfig(unitId string, name string, config string) error {
	var unit common_interfaces.IUnit

	c.mtx.Lock()
	for _, s := range c.units {
		if s.Id() == unitId {
			unit = s
		}
	}
	c.mtx.Unlock()

	if unit != nil {
		unit.Stop()
		oldName := unit.Name()

		if oldName != name {

			nameIsExists := false
			c.mtx.Lock()
			for _, s := range c.units {
				if s.Name() == name {
					nameIsExists = true
				}
			}
			c.mtx.Unlock()

			if !nameIsExists {
				unit.SetName(name)
				c.iDataStorage.RenameItems(oldName+"/", name+"/")
			}
		}

		unit.SetConfig(config)

		unit.Start(c.iDataStorage)
	}

	return nil
}
