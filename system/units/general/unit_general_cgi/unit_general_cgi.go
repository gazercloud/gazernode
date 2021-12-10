package unit_general_cgi

import (
	"encoding/json"
	"errors"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/resources"
	"github.com/gazercloud/gazernode/system/units/units_common"
	"github.com/gazercloud/gazernode/utilities/logger"
	"os/exec"
	"time"
)

type UnitGeneralCGI struct {
	units_common.Unit
	command      string
	arguments    string
	maxValueSize int
	periodMs     int
	showError    bool
	CutSpecChars bool
}

func New() common_interfaces.IUnit {
	var c UnitGeneralCGI
	return &c
}

const (
	ItemNameResult = "Result"
)

var Image []byte

func init() {
	Image = resources.R_files_sensors_unit_general_console_png
}

func (c *UnitGeneralCGI) GetConfigMeta() string {
	meta := units_common.NewUnitConfigItem("", "", "", "", "", "", "")
	meta.Add("command", "Command", "", "string", "", "", "")
	meta.Add("arguments", "Arguments", "", "string", "", "", "")
	meta.Add("period", "Period, ms", "1000", "num", "0", "999999", "0")
	meta.Add("max_value_size", "Max Size Of Result", "100", "num", "1", "999", "0")
	meta.Add("show_error", "Show error in result", "true", "bool", "", "", "")
	meta.Add("cut_new_lines", "Cut special chars", "true", "bool", "", "", "")
	return meta.Marshal()
}

func (c *UnitGeneralCGI) InternalUnitStart() error {
	var err error
	c.SetString(ItemNameResult, "", "")
	c.SetMainItem(ItemNameResult)

	type Config struct {
		Command      string  `json:"command"`
		Arguments    string  `json:"arguments"`
		Period       float64 `json:"period"`
		ShowError    bool    `json:"show_error"`
		CutNewLines  bool    `json:"cut_new_lines"`
		MaxValueSize float64 `json:"max_value_size"`
	}

	var config Config
	config.CutNewLines = true
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

	c.maxValueSize = int(config.MaxValueSize)
	c.showError = config.ShowError
	c.CutSpecChars = config.CutNewLines

	go c.Tick()
	return nil
}

func (c *UnitGeneralCGI) InternalUnitStop() {
}

func (c *UnitGeneralCGI) Tick() {
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
			out = out[:1024]
		}

		result := string(out)

		resultRunes := []rune(result)

		if len(resultRunes) > c.maxValueSize {
			result = string(resultRunes[:c.maxValueSize])
		}

		if c.CutSpecChars {
			resString := make([]rune, 0)
			for _, r := range result {
				if r > 31 {
					resString = append(resString, r)
				}
			}
			result = string(resString)
		}

		if err == nil {
			c.SetString(ItemNameResult, result, "")
			c.SetError("")
		} else {
			if c.showError {
				c.SetString(ItemNameResult, result, "error")
			} else {
				c.SetString(ItemNameResult, err.Error(), "error")
			}
			c.SetError(err.Error())
		}
	}
	c.SetString(ItemNameResult, "", "stopped")
	c.Started = false
}
