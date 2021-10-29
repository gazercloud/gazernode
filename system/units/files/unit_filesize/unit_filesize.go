package unit_filesize

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/resources"
	"github.com/gazercloud/gazernode/system/units/units_common"
	"os"
	"time"
)

type UnitFileSize struct {
	units_common.Unit
	fileName string
	periodMs int
}

func New() common_interfaces.IUnit {
	var c UnitFileSize
	return &c
}

const (
	ItemNameSize   = "FileSize"
	ItemNameResult = "Result"
)

var Image []byte

func init() {
	Image = resources.R_files_sensors_unit_file_file_size_png
}

func (c *UnitFileSize) GetConfigMeta() string {
	meta := units_common.NewUnitConfigItem("", "", "", "", "", "", "")
	meta.Add("file_name", "File Name", "file.txt", "string", "", "", "")
	meta.Add("period", "Period, ms", "1000", "num", "0", "999999", "0")
	return meta.Marshal()
}

func (c *UnitFileSize) InternalUnitStart() error {
	var err error
	c.SetString(ItemNameSize, "", "")
	c.SetMainItem(ItemNameSize)

	type Config struct {
		FileName string  `json:"file_name"`
		Period   float64 `json:"period"`
	}

	var config Config
	err = json.Unmarshal([]byte(c.GetConfig()), &config)
	if err != nil {
		err = errors.New("config error")
		c.SetString(ItemNameSize, err.Error(), "error")
		return err
	}

	c.fileName = config.FileName
	if c.fileName == "" {
		err = errors.New("wrong file")
		c.SetString(ItemNameSize, err.Error(), "error")
		return err
	}

	c.periodMs = int(config.Period)
	if c.periodMs < 100 {
		err = errors.New("wrong period")
		c.SetString(ItemNameSize, err.Error(), "error")
		return err
	}

	go c.Tick()
	return nil
}

func (c *UnitFileSize) InternalUnitStop() {
}

func (c *UnitFileSize) Tick() {
	c.Started = true
	dtOperationTime := time.Now().UTC()
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

		stat, err := os.Stat(c.fileName)
		if err == nil {
			c.SetString(ItemNameSize, fmt.Sprint(stat.Size()), "bytes")
			c.SetError("")
		} else {
			c.SetString(ItemNameSize, "", "")
			c.SetError(err.Error())
		}
	}
	c.SetString(ItemNameSize, "", "stopped")
	c.Started = false
}
