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

	disk     string
	periodMs int
}

var Image []byte

func init() {
	Image = resources.R_files_sensors_sensor_windows_ram_png
}

func New() common_interfaces.IUnit {
	var c UnitStorage
	return &c
}

func (c *UnitStorage) InternalUnitStart() error {
	var err error
	c.SetString("UsedPercents", "", "")
	c.SetMainItem("UsedPercents")

	type Config struct {
		Path   string  `json:"path"`
		Period float64 `json:"period"`
	}

	var config Config
	err = json.Unmarshal([]byte(c.GetConfig()), &config)
	if err != nil {
		err = errors.New("config error")
		c.SetString("UsedPercents", err.Error(), "error")
		return err
	}

	c.drive = config.Path
	if c.fileName == "" {
		err = errors.New("wrong path")
		c.SetString("UsedPercents", err.Error(), "error")
		return err
	}

	c.periodMs = int(config.Period)
	if c.periodMs < 100 {
		err = errors.New("wrong period")
		c.SetString("UsedPercents", err.Error(), "error")
		return err
	}

	diskName := strings.ReplaceAll(c.disk, "/", "_")
	c.SetString(diskName+"/Total", "", "")
	c.SetString(diskName+"/Free", "", "")
	c.SetString(diskName+"/Used", "", "")
	c.SetString(diskName+"/Utilization", "", "")

	go c.Tick()
	return nil
}

func (c *UnitStorage) InternalUnitStop() {
}

func (c *UnitStorage) GetConfigMeta() string {
	meta := units_common.NewUnitConfigItem("", "", "", "", "", "", "")
	meta.Add("path", "Path", "/", "string", "", "", "")
	meta.Add("period", "Period, ms", "1000", "num", "0", "999999", "")
	return meta.Marshal()
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

		var TotalSpace uint64
		var UsedSpace uint64

		diskName := strings.ReplaceAll(c.disk, "/", "_")
		var free, total uint64

		var stat unix.Statfs_t
		err = unix.Statfs(c.disk, &stat)
		free = uint64(stat.Bsize) * stat.Bfree
		total = uint64(stat.Bsize) * stat.Blocks

		if err != nil {
			c.SetString(diskName+"/Total", "", "error")
			//c.SetString(disk+"/Available", "", "error")
			c.SetString(diskName+"/Free", "", "error")
			c.SetString(diskName+"/Used", "", "error")
			c.SetString(diskName+"/Utilization", "", "error")
		} else {
			c.SetUInt64(diskName+"/Total", total/1024/1024, "MB")
			//c.SetUInt64(disk+"/Available", avail / 1024 / 1024, "MB")
			c.SetUInt64(diskName+"/Free", free/1024/1024, "MB")
			c.SetUInt64(diskName+"/Used", (total-free)/1024/1024, "MB")
			c.SetFloat64(diskName+"/Utilization", 100*float64(total-free)/float64(total), "%", 1)

			TotalSpace += total
			UsedSpace += total - free
		}

		//summaryTotal := strconv.FormatFloat(float64(TotalSpace) / 1024 / 1024 / 1024 / 1024, 'f', 1, 64)
		//summaryUsed := strconv.FormatFloat(float64(UsedSpace) / 1024 / 1024 / 1024 / 1024, 'f', 1, 64)
		summaryUtilization := strconv.FormatFloat(100*float64(UsedSpace)/float64(TotalSpace), 'f', 1, 64)

		summary := summaryUtilization

		c.SetString("UsedPercents", summary, "%")
	}

	c.Started = false
}
