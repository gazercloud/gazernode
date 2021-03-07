package unit_csv_export

import (
	"encoding/json"
	"errors"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/resources"
	"github.com/gazercloud/gazernode/system/units/units_common"
	"os"
	"time"
)

type Item struct {
	Name     string `json:"item_name"`
	FileName string `json:"file_name"`
}

type Config struct {
	Directory string  `json:"directory"`
	Period    float64 `json:"period"`
	Separator string  `json:"separator"`
	Items     []Item  `json:"items"`
}

type UnitCsvExport struct {
	units_common.Unit
	fileName string
	periodMs int
	config   Config
}

func New() common_interfaces.IUnit {
	var c UnitCsvExport
	return &c
}

const (
	ItemNameStatus = "Status"
)

var Image []byte

func init() {
	Image = resources.R_files_sensors_sensor_file_csv_export_png
	//Image = uiresources.ResBin("icons/material/file/drawable-hdpi/outline_text_snippet_black_48dp.png")
}

func (c *UnitCsvExport) GetConfigMeta() string {
	meta := units_common.NewUnitConfigItem("", "", "", "", "", "", "")
	meta.Add("directory", "Export path", "", "string", "", "", "")
	meta.Add("period", "Period, ms", "1000", "num", "0", "999999", "")
	meta.Add("separator", "Separator", ",", "string", "", "", "")
	t1 := meta.Add("items", "Items", "", "table", "", "", "")
	t1.Add("item_name", "Item Name", "Unit/Item", "string", "", "", "data-items")
	t1.Add("file_name", "File Name", "item.csv", "string", "", "", "")
	return meta.Marshal()
}

func (c *UnitCsvExport) InternalUnitStart() error {
	var err error
	c.SetString(ItemNameStatus, "", "starting")
	c.SetMainItem(ItemNameStatus)

	c.config.Separator = ","

	err = json.Unmarshal([]byte(c.GetConfig()), &c.config)
	if err != nil {
		err = errors.New("config error")
		c.SetString(ItemNameStatus, err.Error(), "error")
		return err
	}

	if c.config.Directory == "" {
		err = errors.New("wrong directory")
		c.SetString(ItemNameStatus, err.Error(), "error")
		return err
	}

	if c.config.Period < 100 {
		err = errors.New("wrong period")
		c.SetString(ItemNameStatus, err.Error(), "error")
		return err
	}

	go c.Tick()

	c.SetString(ItemNameStatus, "", "started")
	return nil
}

func (c *UnitCsvExport) InternalUnitStop() {
}

func (c *UnitCsvExport) Tick() {
	c.Started = true
	dtOperationTime := time.Now().UTC()
	for !c.Stopping {
		for {
			if c.Stopping || time.Now().Sub(dtOperationTime) > time.Duration(c.config.Period)*time.Millisecond {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
		if c.Stopping {
			break
		}
		dtOperationTime = time.Now().UTC()

		for _, item := range c.config.Items {
			value, err := c.GetItem(item.Name)
			if err == nil {
				var f *os.File
				_ = os.MkdirAll(c.config.Directory, 0755)
				f, err = os.OpenFile(c.config.Directory+"/"+item.FileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 755)
				if err != nil {
					continue
				}
				sep := c.config.Separator

				line := ""
				line += time.Unix(0, value.DT*1000).Format("2006-01-02 15:04:05.000") + sep + value.Value + sep + value.UOM + "\r\n"
				_, err = f.Write([]byte(line))
				_ = f.Close()
			}
		}
	}
	c.SetString(ItemNameStatus, "", "stopped")
	c.Started = false
}
