package unit_gazer_cloud

import (
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/resources"
	"github.com/gazercloud/gazernode/system/units/units_common"
	"time"
)

type UnitGazerCloud struct {
	units_common.Unit
}

var Image []byte

func init() {
	Image = resources.R_files_sensors_category_gazer_png
}

func New() common_interfaces.IUnit {
	var c UnitGazerCloud
	return &c
}

func (c *UnitGazerCloud) InternalUnitStart() error {
	c.SetMainItem("CallsPerSecond")

	c.SetString("CallsPerSecond", "", "")
	c.SetString("InTraffic", "", "")
	c.SetString("OutTraffic", "", "")

	go c.Tick()
	return nil
}

func (c *UnitGazerCloud) InternalUnitStop() {
}

func (c *UnitGazerCloud) GetConfigMeta() string {
	meta := units_common.NewUnitConfigItem("", "", "", "", "", "", "")
	return meta.Marshal()
}

func (c *UnitGazerCloud) Tick() {
	c.Started = true
	for !c.Stopping {
		for i := 0; i < 10; i++ {
			if c.Stopping {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}

		statGazerCloud := c.IDataStorage().StatGazerCloud()

		c.SetFloat64("CallsPerSecond", statGazerCloud.CallsPerSecond, "", 1)
		c.SetFloat64("InTraffic", statGazerCloud.ReceiveSpeed/1024, "KB/sec", 1)
		c.SetFloat64("OutTraffic", statGazerCloud.SendSpeed/1024, "KB/sec", 1)
	}

	time.Sleep(1 * time.Millisecond)
	c.SetString("CallsPerSecond", "", "stopped")
	c.SetString("InTraffic", "", "")
	c.SetString("OutTraffic", "", "")

	c.Started = false
}
