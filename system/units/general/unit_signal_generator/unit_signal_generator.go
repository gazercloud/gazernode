package unit_signal_generator

import (
	"encoding/json"
	"errors"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/resources"
	"github.com/gazercloud/gazernode/system/units/units_common"
	"math"
	"math/rand"
	"time"
)

type Item struct {
	Name         string  `json:"name"`
	SignalType   string  `json:"type"`
	SignalPeriod float64 `json:"period"`
	SignalMin    float64 `json:"min"`
	SignalMax    float64 `json:"max"`
	Precision    float64 `json:"precision"`
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
	Image = resources.R_files_sensors_unit_general_signal_generator_png
}

func (c *UnitSignalGenerator) GetConfigMeta() string {
	meta := units_common.NewUnitConfigItem("", "", "", "", "", "", "")
	t1 := meta.Add("items", "Items", "", "table", "", "", "")
	t1.Add("name", "Item Name", "sin", "string", "", "", "")
	t1.Add("type", "Signal Type", "sin", "string", "", "", "")
	t1.Add("period", "Period, sec", "10", "num", "0", "99999", "")
	t1.Add("min", "Min Value", "0", "num", "-999999999", "999999999", "")
	t1.Add("max", "Max Value", "100", "num", "-999999999", "999999999", "")
	t1.Add("precision", "Precision", "3", "num", "0", "99", "")
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

	for _, item := range c.config.Items {
		if item.Precision < 0 || item.Precision > 100 {
			err = errors.New("wrong precision")
			c.SetString(ItemNameStatus, err.Error(), "error")
			return err
		}
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

		itemNames := make([]string, 0)
		itemNamesMap := make(map[string]bool)
		for _, item := range items {
			if _, alreadyAdded := itemNamesMap[item.Name]; !alreadyAdded {
				itemNamesMap[item.Name] = true
				itemNames = append(itemNames, item.Name)
			}
		}

		for _, itemName := range itemNames {
			value := 0.0
			maxPrecision := 0
			for _, item := range items {
				if item.Name == itemName && item.SignalPeriod > 0 {
					tItem := tMS % int64(item.SignalPeriod*1000)
					itemProgress := float64(tItem) / (item.SignalPeriod * 1000)
					if item.SignalType == "sin" {
						value += c.genSin(item.SignalMin, item.SignalMax, itemProgress)
					}
					if item.SignalType == "sinsin" {
						value += c.genSinSin(item.SignalMin, item.SignalMax, itemProgress)
					}
					if item.SignalType == "tan" {
						value += c.genTg(item.SignalMin, item.SignalMax, itemProgress)
					}
					if item.SignalType == "meander" {
						value += c.genMeander(item.SignalMin, item.SignalMax, itemProgress)
					}
					if item.SignalType == "noise" {
						value += c.genNoise(item.SignalMin, item.SignalMax, itemProgress)
					}

				}

				if int(item.Precision) > maxPrecision {
					maxPrecision = int(item.Precision)
				}
			}
			c.SetFloat64(itemName, value, "", maxPrecision)
		}

	}

	for _, item := range items {
		c.SetString(item.Name, "", "stopped")
	}

	c.SetString(ItemNameStatus, "", "stopped")
	c.Started = false
}

func (c *UnitSignalGenerator) genSin(min float64, max float64, progress01 float64) float64 {
	value := math.Sin(progress01 * 2 * math.Pi)
	delta := max - min
	value = (value+1)/2*delta + min
	return value
}

func (c *UnitSignalGenerator) genSinSin(min float64, max float64, progress01 float64) float64 {
	prSmall := float64(int(math.Round(progress01*100))%20) / 20.0
	value := math.Sin(progress01 * 2 * math.Pi)
	value += math.Sin(prSmall * 2 * math.Pi)
	delta := max - min
	value = (value+1)/2*delta + min
	return value
}

func (c *UnitSignalGenerator) genTg(min float64, max float64, progress01 float64) float64 {
	value := math.Tan(progress01 * 2 * math.Pi)
	delta := max - min
	value = (value+1)/2*delta + min
	return value
}

func (c *UnitSignalGenerator) genMeander(min float64, max float64, progress01 float64) float64 {
	value := 0.0
	if progress01 >= 0.5 {
		value = 1
	}
	delta := max - min
	value = value*delta + min
	return value
}

func (c *UnitSignalGenerator) genNoise(min float64, max float64, progress01 float64) float64 {
	delta := max - min
	value := rand.Float64()*delta + min
	return value
}
