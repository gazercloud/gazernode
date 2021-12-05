package unit_general_cgi_key_value

import (
	"encoding/json"
	"errors"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/resources"
	"github.com/gazercloud/gazernode/system/units/units_common"
	"github.com/gazercloud/gazernode/utilities/logger"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type ConfigItem struct {
	Name   string  `json:"name"`
	UOM    string  `json:"uom"`
	Scale  float64 `json:"scale"`
	Offset float64 `json:"offset"`
}

type UnitGeneralCGIKeyValue struct {
	units_common.Unit
	command   string
	arguments string
	periodMs  int

	showError bool

	receiveAll bool
	items      map[string]*ConfigItem

	receivedVariables map[string]string
}

func New() common_interfaces.IUnit {
	var c UnitGeneralCGIKeyValue
	return &c
}

const (
	ItemNameResult = "status"
)

var Image []byte

func init() {
	Image = resources.R_files_sensors_unit_general_console_key_value_png
}

func (c *UnitGeneralCGIKeyValue) GetConfigMeta() string {
	meta := units_common.NewUnitConfigItem("", "", "", "", "", "", "")
	meta.Add("command", "Command", "", "string", "", "", "")
	meta.Add("arguments", "Arguments", "", "string", "", "", "")
	meta.Add("period", "Period, ms", "1000", "num", "0", "999999", "0")
	meta.Add("show_error", "Show error in result", "true", "bool", "", "", "")
	meta.Add("receive_all", "Receive All", "true", "bool", "", "", "")
	t1 := meta.Add("items", "Elements", "", "table", "", "", "")
	t1.Add("name", "ID", "item_name", "string", "", "", "")
	t1.Add("uom", "UOM", "V", "string", "", "", "")
	t1.Add("scale", "Scale", "1", "num", "-999999999", "999999999", "3")
	t1.Add("offset", "Offset", "0", "num", "-999999999", "999999999", "3")
	return meta.Marshal()
}

func (c *UnitGeneralCGIKeyValue) InternalUnitStart() error {
	var err error
	c.SetString(ItemNameResult, "", "")
	c.SetMainItem(ItemNameResult)

	type Config struct {
		Command    string        `json:"command"`
		Arguments  string        `json:"arguments"`
		Period     float64       `json:"period"`
		ShowError  bool          `json:"show_error"`
		ReceiveAll bool          `json:"receive_all"`
		Items      []*ConfigItem `json:"items"`
	}

	var config Config
	err = json.Unmarshal([]byte(c.GetConfig()), &config)
	if err != nil {
		err = errors.New("config error")
		c.SetString(ItemNameResult, err.Error(), "error")
		return err
	}

	c.command = config.Command
	if c.command == "" {
		err = errors.New("wrong command")
		c.SetString(ItemNameResult, err.Error(), "error")
		return err
	}

	c.arguments = config.Arguments

	c.periodMs = int(config.Period)
	if c.periodMs < 100 {
		err = errors.New("wrong period")
		c.SetString(ItemNameResult, err.Error(), "error")
		return err
	}

	c.showError = config.ShowError

	c.receiveAll = config.ReceiveAll

	c.items = make(map[string]*ConfigItem)

	for _, item := range config.Items {
		c.items[item.Name] = item
	}

	c.receivedVariables = make(map[string]string)

	go c.Tick()
	return nil
}

func (c *UnitGeneralCGIKeyValue) InternalUnitStop() {
}

func (c *UnitGeneralCGIKeyValue) Tick() {
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

		cmd := exec.Command(c.command, c.arguments)
		out, err := cmd.CombinedOutput()
		if err != nil {
			logger.Println("cmd.Run() failed with error:", err)
		}

		if len(out) > 1024 {
			err = errors.New("too much data")
			out = []byte(err.Error())
		}

		if err != nil {
			if c.showError {
				c.SetString(ItemNameResult, string(out), "error")
			} else {
				c.SetString(ItemNameResult, err.Error(), "error")
			}
		} else {
			lines := strings.Split(string(out), "\n")
			for _, currentLine := range lines {

				if len(currentLine) > 0 {
					parts := strings.Split(currentLine, "=")
					if len(parts) > 1 {
						if len(parts[0]) > 0 {
							key := parts[0]
							value := parts[1]

							value = strings.ReplaceAll(value, "\r", "")
							value = strings.ReplaceAll(value, "\n", "")

							if item, ok := c.items[key]; ok {
								finalValue := value
								valueAsFloat, err := strconv.ParseFloat(value, 64)
								if err == nil {
									valueAsFloat = valueAsFloat*item.Scale + item.Offset
									finalValue = strconv.FormatFloat(valueAsFloat, 'f', -1, 64)
								}
								c.receivedVariables[key] = finalValue
								c.SetString(key, finalValue, item.UOM)
							} else {
								if c.receiveAll {
									finalValue := value
									c.receivedVariables[key] = finalValue
									c.SetString(key, finalValue, "")
								}
							}

						}
					}

					c.SetString("status", "ok", "")
				}
			}
		}
	}
	c.SetString(ItemNameResult, "", "stopped")
	c.Started = false
}
