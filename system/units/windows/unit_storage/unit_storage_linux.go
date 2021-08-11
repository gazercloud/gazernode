package unit_storage

import (
	"encoding/json"
	"errors"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/resources"
	"github.com/gazercloud/gazernode/system/units/units_common"
	"golang.org/x/sys/unix"
	"time"
)

type UnitStorage struct {
	units_common.Unit

	disk     string
	periodMs int
}

var Image []byte

func init() {
	Image = resources.R_files_sensors_unit_computer_storage_png
}

func New() common_interfaces.IUnit {
	var c UnitStorage
	return &c
}

func (c *UnitStorage) InternalUnitStart() error {
	var err error
	c.SetMainItem("SpaceUsedPercents")

	type Config struct {
		Path   string  `json:"path"`
		Period float64 `json:"period"`
	}

	var config Config
	err = json.Unmarshal([]byte(c.GetConfig()), &config)
	if err != nil {
		err = errors.New("config error")
		c.SetString("Status", err.Error(), "error")
		return err
	}

	c.disk = config.Path
	if c.disk == "" {
		err = errors.New("wrong path")
		c.SetString("Status", err.Error(), "error")
		return err
	}

	c.periodMs = int(config.Period)
	if c.periodMs < 100 {
		err = errors.New("wrong period")
		c.SetString("Status", err.Error(), "error")
		return err
	}

	c.SetString("SpaceTotal", "", "")
	c.SetString("SpaceFree", "", "")
	c.SetString("SpaceUsed", "", "")

	c.SetString("BlocksTotal", "", "")
	c.SetString("BlocksFree", "", "")
	c.SetString("BlocksUsed", "", "")

	c.SetString("INodesTotal", "", "")
	c.SetString("INodesFree", "", "")
	c.SetString("INodesUsed", "", "")

	c.SetString("SpaceUsedPercents", "", "")
	c.SetString("BlocksUsedPercents", "", "")
	c.SetString("INodesUsedPercents", "", "")

	c.SetString("Status", "", "")

	go c.Tick()
	return nil
}

func (c *UnitStorage) InternalUnitStop() {
}

func (c *UnitStorage) GetConfigMeta() string {
	meta := units_common.NewUnitConfigItem("", "", "", "", "", "", "")
	meta.Add("path", "Mount Path", "/", "string", "", "", "")
	meta.Add("period", "Period, ms", "1000", "num", "0", "999999", "")
	return meta.Marshal()
}

func (c *UnitStorage) Tick() {
	var err error
	dtOperationTime := time.Time{}

	c.Started = true
	for !c.Stopping {
		for {
			if c.Stopping || time.Now().Sub(dtOperationTime) > time.Duration(c.periodMs)*time.Millisecond {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
		if c.Stopping {
			break
		}
		dtOperationTime = time.Now().UTC()

		var free, total uint64

		var stat unix.Statfs_t
		err = unix.Statfs(c.disk, &stat)
		free = uint64(stat.Bsize) * stat.Bfree
		total = uint64(stat.Bsize) * stat.Blocks

		if err != nil {
			c.SetString("Status", err.Error(), "error")

			c.SetString("SpaceTotal", "", "error")
			c.SetString("SpaceFree", "", "error")
			c.SetString("SpaceUsed", "", "error")

			c.SetString("BlocksTotal", "", "error")
			c.SetString("BlocksFree", "", "error")
			c.SetString("BlocksUsed", "", "error")

			c.SetString("INodesTotal", "", "error")
			c.SetString("INodesFree", "", "error")
			c.SetString("INodesUsed", "", "error")

			c.SetString("SpaceUsedPercents", "", "error")
			c.SetString("BlocksUsedPercents", "", "error")
			c.SetString("INodesUsedPercents", "", "error")

		} else {
			c.SetString("Status", "", "")

			c.SetUInt64("SpaceTotal", total/1024/1024, "MB")
			c.SetUInt64("SpaceFree", free/1024/1024, "MB")
			c.SetUInt64("SpaceUsed", (total-free)/1024/1024, "MB")

			c.SetUInt64("BlocksTotal", stat.Blocks, "")
			c.SetUInt64("BlocksFree", stat.Bfree, "")
			c.SetUInt64("BlocksUsed", stat.Blocks-stat.Bfree, "")

			c.SetUInt64("INodesTotal", stat.Files, "")
			c.SetUInt64("INodesFree", stat.Ffree, "")
			c.SetUInt64("INodesUsed", stat.Files-stat.Ffree, "")

			c.SetFloat64("SpaceUsedPercents", 100*float64(total-free)/float64(total), "%", 1)
			c.SetFloat64("BlocksUsedPercents", 100*float64(stat.Blocks-stat.Bfree)/float64(stat.Blocks), "%", 1)
			c.SetFloat64("INodesUsedPercents", 100*float64(stat.Files-stat.Ffree)/float64(stat.Files), "%", 1)
		}

		//summaryTotal := strconv.FormatFloat(float64(TotalSpace) / 1024 / 1024 / 1024 / 1024, 'f', 1, 64)
		//summaryUsed := strconv.FormatFloat(float64(UsedSpace) / 1024 / 1024 / 1024 / 1024, 'f', 1, 64)
	}

	c.Started = false
}
