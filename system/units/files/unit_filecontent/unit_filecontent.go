package unit_filecontent

import (
	"encoding/json"
	"errors"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/resources"
	"github.com/gazercloud/gazernode/system/units/units_common"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

type UnitFileContent struct {
	units_common.Unit
	fileName  string
	periodMs  int
	trim      bool
	parse     bool
	scale     float64
	offset    float64
	uom       string
	precision int
}

func New() common_interfaces.IUnit {
	var c UnitFileContent
	return &c
}

const (
	ItemNameContent = "Content"
)

var Image []byte

func init() {
	Image = resources.R_files_sensors_unit_file_file_content_png
}

func (c *UnitFileContent) GetConfigMeta() string {
	meta := units_common.NewUnitConfigItem("", "", "", "", "", "", "")
	meta.Add("file_name", "File Name", "file.txt", "string", "", "", "")
	meta.Add("period", "Period, ms", "1000", "num", "0", "999999", "0")
	meta.Add("trim", "Trim", "true", "bool", "", "", "")
	meta.Add("parse", "Parse", "true", "bool", "", "", "")
	meta.Add("scale", "Scale", "1", "num", "-999999999", "99999999", "6")
	meta.Add("offset", "Offset", "0", "num", "-999999999", "99999999", "6")
	meta.Add("uom", "UOM", "", "string", "", "", "")
	meta.Add("precision", "Precision", "3", "num", "0", "99", "")
	return meta.Marshal()
}

func (c *UnitFileContent) InternalUnitStart() error {
	var err error
	c.SetString(ItemNameContent, "", "")
	c.SetMainItem(ItemNameContent)

	type Config struct {
		FileName   string  `json:"file_name"`
		Period     float64 `json:"period"`
		Trim       bool    `json:"trim"`
		ParseFloat bool    `json:"parse"`
		Scale      float64 `json:"scale"`
		Offset     float64 `json:"offset"`
		UOM        string  `json:"uom"`
		Precision  float64 `json:"precision"`
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

	c.precision = int(config.Precision)
	if c.precision < 0 || c.precision > 100 {
		err = errors.New("wrong precision")
		c.SetString(ItemNameContent, err.Error(), "error")
		return err
	}

	c.trim = config.Trim
	c.parse = config.ParseFloat
	c.scale = config.Scale
	c.offset = config.Offset
	c.uom = config.UOM

	go c.Tick()
	return nil
}

func (c *UnitFileContent) InternalUnitStop() {
}

func (c *UnitFileContent) Tick() {
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
		contentStr := string(content)
		if c.trim {
			contentStr = strings.Trim(contentStr, " \n\r\t")
		}

		var contentFloat float64

		if c.parse {
			contentFloat, err = strconv.ParseFloat(contentStr, 64)
			if err == nil {
				contentFloat = contentFloat*c.scale + c.offset
				contentStr = strconv.FormatFloat(contentFloat, 'f', c.precision, 64)
				if strings.Index(contentStr, ".") >= 0 {
					contentStr = strings.TrimRight(contentStr, "0")
				}
			}
		}

		if err == nil {
			c.SetString(ItemNameContent, contentStr, c.uom)
			c.SetError("")
		} else {
			c.SetString(ItemNameContent, string(content), "error")
			c.SetError(err.Error())
		}
	}
	c.SetString(ItemNameContent, "", "stopped")
	c.Started = false
}
