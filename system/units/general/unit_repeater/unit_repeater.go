package unit_repeater

import (
	"encoding/json"
	"errors"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/resources"
	"github.com/gazercloud/gazernode/system/units/units_common"
	"time"
)

type Item struct {
	ItemFrom string `json:"item_from"`
	ItemTo   string `json:"item_to"`
}

type Config struct {
	Items []Item `json:"items"`
}

type UnitRepeater struct {
	units_common.Unit
	config Config
}

func New() common_interfaces.IUnit {
	var c UnitRepeater
	return &c
}

const (
	ItemNameStatus = "Status"
)

var Image []byte

func init() {
	Image = resources.R_files_sensors_unit_network_tcp_telnet_control_png
}

func (c *UnitRepeater) GetConfigMeta() string {
	meta := units_common.NewUnitConfigItem("", "", "", "", "", "", "")
	t1 := meta.Add("items", "Items", "", "table", "", "", "")
	t1.Add("item_from", "Source Item", "Unit/Item", "string", "", "", "data-item")
	t1.Add("item_to", "Target Item", "Unit/Item", "string", "", "", "data-item")
	return meta.Marshal()
}

func (c *UnitRepeater) InternalUnitStart() error {
	var err error

	err = json.Unmarshal([]byte(c.GetConfig()), &c.config)
	if err != nil {
		err = errors.New("config error")
		c.SetString(ItemNameStatus, err.Error(), "error")
		return err
	}

	c.SetMainItem(ItemNameStatus)

	c.SetString(ItemNameStatus, "", "starting")

	for _, item := range c.config.Items {
		c.AddToWatch(item.ItemFrom)
	}

	go c.Tick()
	return nil
}

func (c *UnitRepeater) InternalUnitStop() {
	c.Stopping = true
}

func (c *UnitRepeater) Tick() {
	c.Started = true
	c.SetString(ItemNameStatus, "", "started")
	for !c.Stopping {
		time.Sleep(100 * time.Millisecond)
		if c.Stopping {
			break
		}
	}
	c.SetString(ItemNameStatus, "", "stopped")
	c.Started = false
}
