package unit_filecontent

import (
	"encoding/json"
	"errors"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/resources"
	"github.com/gazercloud/gazernode/system/units/units_common"
	"io/ioutil"
	"time"
)

type UnitParamInLog struct {
	units_common.Unit
	fileName string
	periodMs int
}

func New() common_interfaces.IUnit {
	var c UnitParamInLog
	return &c
}

const (
	ItemNameContent = "Content"
)

var Image []byte

func init() {
	Image = resources.R_files_sensors_sensor_file_content_png
}

func (c *UnitParamInLog) GetConfigMeta() string {
	meta := units_common.NewUnitConfigItem("", "", "", "", "", "", "")
	meta.Add("file_name", "File Name", "file.txt", "string", "", "", "")
	meta.Add("period", "Period, ms", "1000", "num", "0", "999999", "")
	return meta.Marshal()
}

func (c *UnitParamInLog) InternalUnitStart() error {
	var err error
	c.SetString(ItemNameContent, "", "")
	c.SetMainItem(ItemNameContent)

	type Config struct {
		FileName string  `json:"file_name"`
		Period   float64 `json:"period"`
	}

	var config Config
	err = json.Unmarshal([]byte(c.GetConfig()), &config)
	if err != nil {
		err = errors.New("config error")
		c.SetString(ItemNameContent, err.Error(), "error")
		return err
	}

	c.fileName = config.FileName
	if c.fileName == "" {
		err = errors.New("wrong file")
		c.SetString(ItemNameContent, err.Error(), "error")
		return err
	}

	c.periodMs = int(config.Period)
	if c.periodMs < 100 {
		err = errors.New("wrong period")
		c.SetString(ItemNameContent, err.Error(), "error")
		return err
	}

	go c.Tick()
	return nil
}

func (c *UnitParamInLog) InternalUnitStop() {
}

func (c *UnitParamInLog) Tick() {
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

		content, err := ioutil.ReadFile(c.fileName)

		if len(content) > 1024 {
			err = errors.New("too much data")
			content = content[:1024]
		}

		if err == nil {
			c.SetString(ItemNameContent, string(content), "")
			c.SetError("")
		} else {
			c.SetString(ItemNameContent, string(content), "error")
			c.SetError(err.Error())
		}
	}
	c.SetString(ItemNameContent, "", "stopped")
	c.Started = false
}
