package unit_system_memory

import (
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/resources"
	"github.com/gazercloud/gazernode/system/units/units_common"
	"github.com/gazercloud/gazernode/utilities/uom"
	"github.com/shirou/gopsutil/mem"
	"time"
)

type UnitSystemMemory struct {
	units_common.Unit
}

var Image []byte

func init() {
	Image = resources.R_files_sensors_unit_computer_memory_png
}

func New() common_interfaces.IUnit {
	var c UnitSystemMemory
	return &c
}

func (c *UnitSystemMemory) InternalUnitStart() error {
	c.SetMainItem("UsedPercent")

	c.SetString("Total", "", "")
	c.SetString("Available", "", "")
	c.SetString("Used", "", "")

	c.SetString("UsedPercent", "", "")

	go c.Tick()
	return nil
}

func (c *UnitSystemMemory) InternalUnitStop() {
}

func (c *UnitSystemMemory) GetConfigMeta() string {
	meta := units_common.NewUnitConfigItem("", "", "", "", "", "", "")
	return meta.Marshal()
}

func (c *UnitSystemMemory) Tick() {
	c.Started = true
	for !c.Stopping {
		for i := 0; i < 10; i++ {
			if c.Stopping {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}

		v, _ := mem.VirtualMemory()

		percents := (float64(v.Used) / float64(v.Total)) * 100.0

		// Common
		c.SetUInt64("Total", v.Total/1048576, uom.MB)
		c.SetUInt64("Available", v.Available/1048576, uom.MB)
		c.SetUInt64("Used", v.Used/1048576, uom.MB)
		c.SetFloat64("UsedPercent", percents, "%", 1)
	}

	time.Sleep(1 * time.Millisecond)
	c.SetString("Total", "", "stopped")
	c.SetString("Available", "", "stopped")
	c.SetString("Used", "", "stopped")
	c.SetString("UsedPercent", "", "stopped")

	c.Started = false
}
