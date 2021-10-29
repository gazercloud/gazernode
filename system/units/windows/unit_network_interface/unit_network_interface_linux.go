package unit_network_interface

import (
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/resources"
	"github.com/gazercloud/gazernode/system/units/units_common"
)

type UnitNetworkInterface struct {
	units_common.Unit

	interfaceName string
}

var Image []byte

func init() {
	Image = resources.R_files_sensors_unit_computer_network_interface_png
}

func New() common_interfaces.IUnit {
	var c UnitNetworkInterface
	return &c
}

func (c *UnitNetworkInterface) GetConfigMeta() string {
	meta := units_common.NewUnitConfigItem("", "", "", "", "", "", "")
	meta.Add("interface_name", "Network Interface", "-", "string", "", "", "network_interface")
	meta.Add("period", "Period, ms", "1000", "num", "0", "999999", "0")
	return meta.Marshal()
}

func (c *UnitNetworkInterface) InternalUnitStart() error {
	go c.Tick()
	return nil
}

func (c *UnitNetworkInterface) InternalUnitStop() {
}

func (c *UnitNetworkInterface) Tick() {
}
