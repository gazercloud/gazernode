package unit_calculator

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/resources"
	"github.com/gazercloud/gazernode/system/units/units_common"
	"github.com/gazercloud/gazernode/utilities/uom"
	"strconv"
	"time"
)

type Item struct {
	Item    string `json:"item"`
	Formula string `json:"formula"`
}

type Config struct {
	Period float64 `json:"period"`
	Items  []Item  `json:"items"`
}

type UnitCalculator struct {
	units_common.Unit
	config Config
}

func New() common_interfaces.IUnit {
	var c UnitCalculator
	return &c
}

const (
	ItemNameStatus = "Status"
)

var Image []byte

func init() {
	Image = resources.R_files_sensors_unit_network_tcp_telnet_control_png
}

func (c *UnitCalculator) GetConfigMeta() string {
	meta := units_common.NewUnitConfigItem("", "", "", "", "", "", "")
	meta.Add("period", "Period, ms", "1000", "num", "0", "999999", "")
	t1 := meta.Add("items", "Items", "", "table", "", "", "")
	t1.Add("item", "Item Name", "item1", "string", "", "", "")
	t1.Add("formula", "Formula", "2+2", "text", "", "", "")
	return meta.Marshal()
}

func (c *UnitCalculator) InternalUnitStart() error {
	var err error

	err = json.Unmarshal([]byte(c.GetConfig()), &c.config)
	if err != nil {
		err = errors.New("config error")
		c.SetString(ItemNameStatus, err.Error(), uom.ERROR)
		return err
	}

	if c.config.Period < 100 {
		err = errors.New("wrong period")
		c.SetString(ItemNameStatus, err.Error(), uom.ERROR)
		return err
	}

	c.SetMainItem(ItemNameStatus)

	c.SetString(ItemNameStatus, "", "starting")

	/*for _, item := range c.config.Items {
		c.AddToWatch(item.Item)
	}*/

	go c.Tick()
	return nil
}

func (c *UnitCalculator) InternalUnitStop() {
	c.Stopping = true
}

func (c *UnitCalculator) Tick() {
	c.Started = true
	c.SetString(ItemNameStatus, "", "started")
	dtOperationTime := time.Time{}

	for _, item := range c.config.Items {
		c.SetString(item.Item, "", uom.STARTED)
	}

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
		dtOperationTime = time.Now()

		for index, _ := range c.config.Items {
			c.calcItem(&c.config.Items[index])
		}
	}

	for _, item := range c.config.Items {
		c.SetString(item.Item, "", uom.STOPPED)
	}

	c.SetString(ItemNameStatus, "", uom.STOPPED)
	c.Started = false
}

func (c *UnitCalculator) calcItem(item *Item) {
	result := ""

	expression, err := govaluate.NewEvaluableExpressionWithFunctions(item.Formula, c.functions())
	if err != nil {
		result = err.Error()
		return
	} else {
		var res interface{}
		parameters := make(map[string]interface{}, 8)
		res, err = expression.Evaluate(parameters)
		if err == nil {
			result = fmt.Sprint(res)
		} else {
			result = err.Error()
		}
	}

	c.SetString(item.Item, result, uom.NONE)
}

func (c *UnitCalculator) ItemChanged(itemName string, value common_interfaces.ItemValue) {
	if !c.Started {
		return
	}
	/*for _, item := range c.config.Items {
		if item.ItemFrom == itemName {
			c.IDataStorage().SetItem(item.ItemTo, value.Value, value.UOM, time.Now().UTC(), true)
		}
	}*/
}

func (c *UnitCalculator) functions() map[string]govaluate.ExpressionFunction {
	functionsResult := map[string]govaluate.ExpressionFunction{
		"item": func(args ...interface{}) (interface{}, error) {
			if len(args) < 1 {
				return nil, errors.New("no_item")
			}
			itemName, ok := args[0].(string)
			if !ok {
				return nil, errors.New("wrong_item_name_type")
			}
			val := c.GetValue(itemName)
			valFloat, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return nil, err
			}
			return valFloat, nil
		},
		"item_string": func(args ...interface{}) (interface{}, error) {
			if len(args) < 1 {
				return nil, errors.New("no_item")
			}
			itemName, ok := args[0].(string)
			if !ok {
				return nil, errors.New("wrong_item_name_type")
			}
			val := c.GetValue(itemName)
			return val, nil
		},
		"bool": func(args ...interface{}) (interface{}, error) {
			res := false
			if len(args) < 1 {
				return nil, errors.New("wrong_argument")
			}
			strVal := fmt.Sprint(args[0])

			// Direct
			if strVal == "true" || strVal == "1" || strVal == "True" {
				return true, nil
			}
			if strVal == "false" || strVal == "0" || strVal == "False" {
				return false, nil
			}

			// From Float/Int
			floatVal, err := strconv.ParseFloat(strVal, 64)
			if err != nil {
				return nil, errors.New("wrong_value")
			}
			intVal := int(floatVal)
			if intVal != 0 {
				res = true
			} else {
				res = false
			}

			return res, nil
		},
		"int": func(args ...interface{}) (interface{}, error) {
			var err error
			res := float64(0)
			if len(args) < 1 {
				return nil, errors.New("wrong_argument")
			}
			strVal := fmt.Sprint(args[0])

			// From Float/Int
			res, err = strconv.ParseFloat(strVal, 64)
			if err != nil {
				return nil, errors.New("wrong_value")
			}

			return int(res), nil
		},
		"float": func(args ...interface{}) (interface{}, error) {
			var err error
			res := float64(0)
			if len(args) < 1 {
				return nil, errors.New("wrong_argument")
			}
			strVal := fmt.Sprint(args[0])

			// From Float/Int
			res, err = strconv.ParseFloat(strVal, 64)
			if err != nil {
				return nil, errors.New("wrong_value")
			}

			return res, nil
		},
		"format": func(args ...interface{}) (interface{}, error) {
			var err error
			if len(args) < 2 {
				return nil, errors.New("wrong_arguments")
			}
			strVal := fmt.Sprint(args[0])
			strPrecision := fmt.Sprint(args[1])

			var val float64
			val, err = strconv.ParseFloat(strVal, 64)
			if err != nil {
				return nil, errors.New("wrong_value")
			}
			var precisionFloat float64
			precisionFloat, err = strconv.ParseFloat(strPrecision, 64)
			if err != nil {
				return nil, errors.New("wrong_value")
			}
			precision := int(precisionFloat)
			if precision < 0 || precisionFloat > 20 {
				return nil, errors.New("wrong_precision")
			}

			return strconv.FormatFloat(val, 'f', precision, 64), nil
		},
	}
	return functionsResult
}
