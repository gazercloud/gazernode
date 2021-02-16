package units_system

import (
	"errors"
	"fmt"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/logger"
	"github.com/gazercloud/gazernode/system/units/files/unit_csv_export"
	"github.com/gazercloud/gazernode/system/units/files/unit_filecontent"
	"github.com/gazercloud/gazernode/system/units/files/unit_filesize"
	"github.com/gazercloud/gazernode/system/units/general/unit_general_cgi"
	"github.com/gazercloud/gazernode/system/units/general/unit_general_cgi_key_value"
	"github.com/gazercloud/gazernode/system/units/general/unit_hhgttg"
	"github.com/gazercloud/gazernode/system/units/general/unit_manual"
	"github.com/gazercloud/gazernode/system/units/general/unit_signal_generator"
	"github.com/gazercloud/gazernode/system/units/network/unit_http_json_items_server"
	"github.com/gazercloud/gazernode/system/units/network/unit_http_json_units_server"
	"github.com/gazercloud/gazernode/system/units/network/unit_ping"
	"github.com/gazercloud/gazernode/system/units/network/unit_tcp_connect"
	unit_tcp_telnet_control "github.com/gazercloud/gazernode/system/units/network/unit_tcp_control"
	unit_serial_port_key_value "github.com/gazercloud/gazernode/system/units/serial_port/serial_port_key_value"
	"github.com/gazercloud/gazernode/system/units/units_common"
	"github.com/gazercloud/gazernode/system/units/windows/unit_process"
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
	unitType = c.RegisterUnit("network_tcp_telnet_control", "network", "TCP Telnet Control", unit_tcp_telnet_control.New, unit_tcp_telnet_control.Image, "")
	unitType = c.RegisterUnit("network_http_json_items_server", "network", "HTTP Json Items Server", unit_http_json_items_server.New, unit_http_json_items_server.Image, "")
	unitType = c.RegisterUnit("network_http_json_units_server", "network", "HTTP Json Units Server", unit_http_json_units_server.New, unit_http_json_units_server.Image, "")

	unitType = c.RegisterUnit("windows_memory", "windows", "OS Memory", unit_system_memory.New, unit_system_memory.Image, "")
	unitType = c.RegisterUnit("windows_process", "windows", "OS Process", unit_process.New, unit_process.Image, "")

	unitType = c.RegisterUnit("file_size", "file", "File Size", unit_filesize.New, unit_filesize.Image, "")
	unitType = c.RegisterUnit("file_content", "file", "File Content", unit_filecontent.New, unit_filecontent.Image, "")
	unitType = c.RegisterUnit("file_csv_export", "file", "CSV Export", unit_csv_export.New, unit_csv_export.Image, "")

	unitType = c.RegisterUnit("general_cgi", "general", "Console", unit_general_cgi.New, unit_general_cgi.Image, "")
	unitType = c.RegisterUnit("general_cgi_key_value", "general", "Console Key=Value", unit_general_cgi_key_value.New, unit_general_cgi_key_value.Image, "")
	unitType = c.RegisterUnit("general_manual", "general", "Manual Items", unit_manual.New, unit_manual.Image, "")
	unitType = c.RegisterUnit("general_hhgttg", "general", "HHGTTG", unit_hhgttg.New, unit_hhgttg.Image, "")
	unitType = c.RegisterUnit("general_signal_generator", "general", "Signal Generator", unit_signal_generator.New, unit_signal_generator.Image, "")

	unitType = c.RegisterUnit("serial_port_key_value", "serial_port", "Serial Port Key=Value", unit_serial_port_key_value.New, unit_serial_port_key_value.Image, "Key/value unit via Serial Port. Format: key=value<new_line>")

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

func (c *UnitsSystem) UnitTypes() []common_interfaces.UnitTypeInfo {
	result := make([]common_interfaces.UnitTypeInfo, 0)
	for _, st := range c.unitTypes {
		var unitTypeInfo common_interfaces.UnitTypeInfo
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

func (c *UnitsSystem) UnitCategories() []common_interfaces.UnitCategoryInfo {
	result := make([]common_interfaces.UnitCategoryInfo, 0)
	addedCats := make(map[string]bool)
	for _, st := range c.unitTypes {
		if _, ok := addedCats[st.Category]; !ok {
			var unitCategoryInfo common_interfaces.UnitCategoryInfo
			unitCategoryInfo.Name = st.Category
			unitCategoryInfo.Image = st.Picture
			result = append(result, unitCategoryInfo)
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

func (c *UnitsSystem) AddUnit(unitType string, unitId string, name string, config string) error {
	nameIsExists := false
	c.mtx.Lock()
	for _, s := range c.units {
		if s.Name() == name {
			nameIsExists = true
		}
	}
	c.mtx.Unlock()

	if !nameIsExists {
		unit := c.MakeUnitByType(unitType, c.iDataStorage)
		if unit != nil {
			unit.SetId(unitId)
			unit.SetName(name)
			unit.SetType(unitType)
			unit.SetConfig(config)
			unit.SetIUnit(unit)
			c.units = append(c.units, unit)
		}
	} else {
		return errors.New("unit already exists")
	}
	return nil
}

func (c *UnitsSystem) GetUnitState(unitId string) (common_interfaces.UnitState, error) {
	var unit common_interfaces.IUnit
	c.mtx.Lock()
	for _, s := range c.units {
		if s.Id() == unitId {
			unit = s
		}
	}
	c.mtx.Unlock()

	if unit != nil {
		var unitState common_interfaces.UnitState
		unitState.Status = ""
		unitState.MainItem = unit.Name() + "/" + unit.MainItem()
		if unit.IsStarted() {
			unitState.Status = "started"
		} else {
			unitState.Status = "stopped"
		}
		return unitState, nil
	}
	return common_interfaces.UnitState{}, errors.New("no unit found")
}

func (c *UnitsSystem) ListOfUnits() []units_common.UnitInfo {
	result := make([]units_common.UnitInfo, 0)
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
