package unit_raspberry_pi_gpio

import (
	"encoding/json"
	"errors"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/resources"
	"github.com/gazercloud/gazernode/system/units/units_common"
	"github.com/gazercloud/gazernode/utilities/uom"
	"github.com/stianeikeland/go-rpio/v4"
	"strconv"
	"time"
)

// "github.com/stianeikeland/go-rpio/v4"

type UnitRaspberryPiGPIO struct {
	units_common.Unit
	periodMs int
	config   Config
}

type ConfigItem struct {
	Name    string `json:"name"`
	Index   string `json:"index"`
	Mode    string `json:"mode"`
	Default string `json:"default"`
	Pull    string `json:"pull"`
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
	t1.Add("index", "GPIO#", "0", "string", "", "", "raspberry-pi-gpio")
	t1.Add("mode", "Mode", "input", "string", "", "", "gpio-mode")
	t1.Add("default", "Default", "0", "string", "", "", "")
	t1.Add("pull", "Pull", "off", "string", "", "", "raspberry-pi-gpio-pull")
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
	if err := rpio.Open(); err != nil {
		c.SetString(ItemNameResult, err.Error(), uom.STARTED)
		c.Started = false
		return
	}
	defer rpio.Close()

	for _, item := range c.config.Pins {
		indexOfPin, err := strconv.ParseInt(item.Index, 10, 64)
		indexOfPinInt := int(indexOfPin)
		if err == nil && indexOfPinInt >= 2 && indexOfPinInt <= 27 {
			if item.Mode == "input" {
				pin := rpio.Pin(indexOfPinInt)
				pin.Input()
				if item.Pull == "off" {
					pin.PullOff()
				}
				if item.Pull == "up" {
					pin.PullUp()
				}
				if item.Pull == "down" {
					pin.PullDown()
				}
			}
			if item.Mode == "output" {
				pin := rpio.Pin(indexOfPinInt)
				pin.Output()
				if item.Default == "1" {
					pin.High()
				} else {
					pin.Low()
				}
			}
			c.SetString(item.Name, "", uom.STARTED)
		} else {
			c.SetString(item.Name, "wrong pin index", "error")
		}
	}

	c.Started = true
	dtOperationTime := time.Now().UTC()

	c.SetString(ItemNameResult, "", uom.STARTED)

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

		for _, item := range c.config.Pins {
			indexOfPin, err := strconv.ParseInt(item.Index, 10, 64)
			indexOfPinInt := int(indexOfPin)
			if err == nil && indexOfPinInt >= 2 && indexOfPinInt <= 27 {
				if item.Mode == "input" {
					pin := rpio.Pin(indexOfPinInt)
					if rpio.ReadPin(pin) == rpio.High {
						c.SetString(item.Name, "1", "")
					} else {
						c.SetString(item.Name, "0", "")
					}
				}
				if item.Mode == "output" {
					pin := rpio.Pin(indexOfPinInt)
					st, err := c.IDataStorage().GetItem(c.Name() + "/" + item.Name)
					if err == nil {
						if st.Value.Value == "1" {
							pin.High()
						} else {
							pin.Low()
						}
					}
				}
			}
		}

	}

	for _, item := range c.config.Pins {
		indexOfPin, err := strconv.ParseInt(item.Index, 10, 64)
		indexOfPinInt := int(indexOfPin)
		if err == nil && indexOfPinInt >= 2 && indexOfPinInt <= 27 {
			pin := rpio.Pin(indexOfPinInt)
			pin.Input()
		}
		c.SetString(item.Name, "", uom.STOPPED)
	}

	c.SetString(ItemNameResult, "", uom.STOPPED)
	c.Started = false
}
