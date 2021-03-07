package unit_signal_generator

import (
	"encoding/json"
	"errors"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/resources"
	"github.com/gazercloud/gazernode/system/units/units_common"
	"math"
	"time"
)

type Item struct {
	Name         string  `json:"name"`
	SignalType   string  `json:"type"`
	SignalPeriod float64 `json:"period"`
	SignalMin    float64 `json:"min"`
	SignalMax    float64 `json:"max"`
}

type Config struct {
	Items []Item `json:"items"`
}

type UnitSignalGenerator struct {
	units_common.Unit
	periodMs int
	config   Config
}

func New() common_interfaces.IUnit {
	var c UnitSignalGenerator
	return &c
}

const (
	ItemNameStatus = "Status"
)

var Image []byte

func init() {
	Image = resources.R_files_sensors_sensor_general_sig_generator_png
}

func (c *UnitSignalGenerator) GetConfigMeta() string {
	meta := units_common.NewUnitConfigItem("", "", "", "", "", "", "")
	t1 := meta.Add("items", "Items", "", "table", "", "", "")
	t1.Add("name", "Item Name", "sin", "string", "", "", "")
	t1.Add("type", "Signal Type", "sin", "string", "", "", "")
	t1.Add("period", "Period, sec", "10", "num", "0", "99999", "")
	t1.Add("min", "Min Value", "0", "num", "-999999999", "999999999", "")
	t1.Add("max", "Max Value", "100", "num", "-999999999", "999999999", "")
	return meta.Marshal()
}

func (c *UnitSignalGenerator) InternalUnitStart() error {
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

func (c *UnitSignalGenerator) InternalUnitStop() {
}

func (c *UnitSignalGenerator) Tick() {
	c.Started = true

	type ItemState struct {
		Item
		Progress float64
	}

	items := make([]ItemState, 0)

	for _, item := range c.config.Items {
		items = append(items, ItemState{Item: item, Progress: 0})
	}

	dtLastActionTime := time.Now()

	for _, item := range items {
		c.SetString(item.Name, "", "started")
	}

	for !c.Stopping {
		for {
			if c.Stopping || time.Now().Sub(dtLastActionTime) > 100*time.Millisecond {
				time.Sleep(1 * time.Microsecond)
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
		if c.Stopping {
			break
		}
		dtLastActionTime = time.Now()

		tMS := time.Now().UTC().UnixNano() / 1000000

		for _, item := range items {
			if item.SignalPeriod > 0 {
				tItem := tMS % int64(item.SignalPeriod*1000)
				itemProgress := float64(tItem) / (item.SignalPeriod * 1000)
				value := math.Sin(itemProgress * 2 * math.Pi)
				delta := item.SignalMax - item.SignalMin
				value = (value+1)/2*delta + item.SignalMin

				c.SetFloat64(item.Name, value, "", 3)
			}
		}

	}

	for _, item := range items {
		c.SetString(item.Name, "", "stopped")
	}

	c.SetString(ItemNameStatus, "", "stopped")
	c.Started = false
}
