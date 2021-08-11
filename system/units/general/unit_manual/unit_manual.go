package unit_manual

import (
	"encoding/json"
	"errors"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/resources"
	"github.com/gazercloud/gazernode/system/units/units_common"
	"time"
)

type Item struct {
	Name      string `json:"item_name"`
	InitValue string `json:"init_value"`
}

type Config struct {
	Items []Item `json:"items"`
}

type UnitManual struct {
	units_common.Unit
	fileName string
	periodMs int
	config   Config
}

func New() common_interfaces.IUnit {
	var c UnitManual
	return &c
}

const (
	ItemNameStatus = "Status"
)

var Image []byte

func init() {
	Image = resources.R_files_sensors_unit_general_manual_items_png
}

func (c *UnitManual) GetConfigMeta() string {
	meta := units_common.NewUnitConfigItem("", "", "", "", "", "", "")
	t1 := meta.Add("items", "Items", "", "table", "", "", "")
	t1.Add("item_name", "Item Name", "item1", "string", "", "", "")
	t1.Add("init_value", "Init Value", "42", "string", "", "", "")
	return meta.Marshal()
}

func (c *UnitManual) InternalUnitStart() error {
	var err error
	c.SetString(ItemNameStatus, "", "starting")
	c.SetMainItem(ItemNameStatus)

	err = json.Unmarshal([]byte(c.GetConfig()), &c.config)
	if err != nil {
		err = errors.New("config error")
		c.SetString(ItemNameStatus, err.Error(), "error")
		return err
	}

	go c.Tick()

	c.SetString(ItemNameStatus, "", "started")
	return nil
}

func (c *UnitManual) InternalUnitStop() {
}

func (c *UnitManual) Tick() {
	c.Started = true

	for _, item := range c.config.Items {
		c.TouchItem(item.Name)
	}

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

	c.SetString(ItemNameStatus, "", "stopped")
	c.Started = false
}
