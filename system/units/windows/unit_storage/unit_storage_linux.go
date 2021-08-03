package unit_storage

import (
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/resources"
	"github.com/gazercloud/gazernode/system/units/units_common"
	"golang.org/x/sys/unix"
	"strconv"
	"strings"
	"time"
)

type UnitStorage struct {
	units_common.Unit
}

var Image []byte

func init() {
	Image = resources.R_files_sensors_sensor_windows_ram_png
}

func New() common_interfaces.IUnit {
	var c UnitStorage
	return &c
}

func (c *UnitStorage) drives() []string {
	drives := make([]string, 0)
	drives = append(drives, "/")
	return drives
}

func (c *UnitStorage) InternalUnitStart() error {
	drives := c.drives()
	c.SetString("UsedPercents", "", "")
	c.SetMainItem("UsedPercents")

	for _, disk := range drives {
		c.SetString(disk+"/Total", "", "")
		//c.SetString(disk+"/Available", "", "")
		c.SetString(disk+"/Free", "", "")
		c.SetString(disk+"/Used", "", "")
		c.SetString(disk+"/Utilization", "", "")
	}

	go c.Tick()
	return nil
}

func (c *UnitStorage) InternalUnitStop() {
}

func (c *UnitStorage) GetConfigMeta() string {
	return ""
}

func (c *UnitStorage) Tick() {
	var err error
	c.Started = true
	for !c.Stopping {
		for i := 0; i < 10; i++ {
			if c.Stopping {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}

		drives := c.drives()

		var TotalSpace int64
		var UsedSpace int64

		for _, disk := range drives {
			diskName := strings.ReplaceAll("/", "_")
			var free, total, avail int64

			var stat unix.Statfs_t
			err = unix.Statfs(disk, &stat)
			free = stat.Bsize * stat.Bfree
			total = stat.Bsize * stat.Blocks
			avail = stat.Bsize * stat.Bfree

			if err != nil {
				c.SetString(diskName+"/Total", "", "error")
				//c.SetString(disk+"/Available", "", "error")
				c.SetString(diskName+"/Free", "", "error")
				c.SetString(diskName+"/Used", "", "error")
				c.SetString(diskName+"/Utilization", "", "error")
			} else {
				c.SetInt64(diskName+"/Total", total/1024/1024, "MB")
				//c.SetUInt64(disk+"/Available", avail / 1024 / 1024, "MB")
				c.SetInt64(diskName+"/Free", free/1024/1024, "MB")
				c.SetInt64(diskName+"/Used", (total-free)/1024/1024, "MB")
				c.SetFloat64(diskName+"/Utilization", 100*float64(total-free)/float64(total), "%", 1)

				TotalSpace += total
				UsedSpace += total - free
			}
		}

		//summaryTotal := strconv.FormatFloat(float64(TotalSpace) / 1024 / 1024 / 1024 / 1024, 'f', 1, 64)
		//summaryUsed := strconv.FormatFloat(float64(UsedSpace) / 1024 / 1024 / 1024 / 1024, 'f', 1, 64)
		summaryUtilization := strconv.FormatFloat(100*float64(UsedSpace)/float64(TotalSpace), 'f', 1, 64)

		summary := summaryUtilization

		c.SetString("UsedPercents", summary, "%")
	}

	c.Started = false
}
