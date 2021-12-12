package unit_raspberry_pi_gpio

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/resources"
	"github.com/gazercloud/gazernode/system/units/units_common"
	"github.com/stianeikeland/go-rpio"
	"os"
	"time"
)

type UnitRaspberryPiGPIO struct {
	units_common.Unit
	periodMs int
	config   Config
}

type ConfigItem struct {
	Name string `json:"name"`
	Mode string `json:"mode"`
}

type Config struct {
	Period float64       `json:"period"`
	Pins   []*ConfigItem `json:"pins"`
}

func New() common_interfaces.IUnit {
	var c UnitRaspberryPiGPIO
	return &c
}

const (
	ItemNameResult = "Result"
)

var Image []byte

func init() {
	Image = resources.R_files_sensors_unit_raspberry_pi_gpio_png
}

func (c *UnitRaspberryPiGPIO) GetConfigMeta() string {
	meta := units_common.NewUnitConfigItem("", "", "", "", "", "", "")
	meta.Add("period", "Period, ms", "1000", "num", "0", "999999", "")
	t1 := meta.Add("pins", "Pins", "", "table", "", "", "")
	t1.Add("name", "Name", "pin_name", "string", "", "", "")
	t1.Add("mode", "Mode", "input", "string", "", "", "gpio-mode")
	return meta.Marshal()
}

func (c *UnitRaspberryPiGPIO) InternalUnitStart() error {
	var err error
	c.SetString(ItemNameResult, "", "")
	c.SetMainItem(ItemNameResult)

	err = json.Unmarshal([]byte(c.GetConfig()), &c.config)
	if err != nil {
		err = errors.New("config error")
		c.SetString(ItemNameResult, err.Error(), "error")
		return err
	}

	c.periodMs = int(c.config.Period)
	if c.periodMs < 100 {
		err = errors.New("wrong period")
		c.SetString(ItemNameResult, err.Error(), "error")
		return err
	}

	go c.Tick()
	return nil
}

func (c *UnitRaspberryPiGPIO) InternalUnitStop() {
}

func (c *UnitRaspberryPiGPIO) Tick() {

	c.SetInt("count of pins", len(c.config.Pins), "init")
	for _, item := range c.config.Pins {
		c.SetString(item.Name, item.Mode, "init")
	}

	c.Started = true
	dtOperationTime := time.Now().UTC()

	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer rpio.Close()

	pin := rpio.Pin(17)
	pin.Output()

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
		pin.Toggle()
	}
	c.SetString(ItemNameResult, "", "stopped")
	c.Started = false
}
