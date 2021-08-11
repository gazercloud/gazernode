package unit_hhgttg

import (
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/resources"
	"github.com/gazercloud/gazernode/system/units/units_common"
	"time"
)

type UnitHHGTTG struct {
	units_common.Unit
}

func New() common_interfaces.IUnit {
	var c UnitHHGTTG
	return &c
}

const (
	ItemNameValue = "Ultimate Question of Life, the Universe, and Everything"
)

var Image []byte

func init() {
	Image = resources.R_files_sensors_unit_general_hhgttg_png
}

func (c *UnitHHGTTG) GetConfigMeta() string {
	meta := units_common.NewUnitConfigItem("", "", "", "", "", "", "")
	return meta.Marshal()
}

func (c *UnitHHGTTG) InternalUnitStart() error {
	c.SetMainItem(ItemNameValue)

	c.SetString(ItemNameValue, "42", "")

	go c.Tick()
	return nil
}

func (c *UnitHHGTTG) InternalUnitStop() {
	c.Stopping = true
	c.SetString(ItemNameValue, "-42", "")
}

func (c *UnitHHGTTG) Tick() {
	c.Started = true
	for !c.Stopping {
		for {
			if c.Stopping {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
		if c.Stopping {
			break
		}
	}
	c.SetString(ItemNameValue, "", "-42")
	c.Started = false
}
