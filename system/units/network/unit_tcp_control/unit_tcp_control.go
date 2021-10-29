package unit_tcp_telnet_control

import (
	"encoding/json"
	"errors"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/resources"
	"github.com/gazercloud/gazernode/system/units/units_common"
	"net"
	"strings"
	"time"
)

type Item struct {
	Name   string `json:"item_name"`
	Format string `json:"format"`
}

type Config struct {
	Addr    string `json:"addr"`
	Timeout float64
	Period  float64
	Items   []Item `json:"items"`
}

type UnitTcpControl struct {
	units_common.Unit
	config Config
}

func New() common_interfaces.IUnit {
	var c UnitTcpControl
	return &c
}

const (
	ItemNameStatus = "Status"
)

var Image []byte

func init() {
	Image = resources.R_files_sensors_unit_network_tcp_telnet_control_png
}

func (c *UnitTcpControl) GetConfigMeta() string {
	meta := units_common.NewUnitConfigItem("", "", "", "", "", "", "")
	meta.Add("addr", "Address", "localhost:7777", "string", "", "", "")
	meta.Add("period", "Period, ms", "1000", "num", "0", "999999", "0")
	meta.Add("timeout", "Timeout, ms", "1000", "num", "0", "999999", "0")
	t1 := meta.Add("items", "Items", "", "table", "", "", "")
	t1.Add("item_name", "Item Name", "Unit/Item", "string", "", "", "data-item")
	t1.Add("format", "Format", "item=%VALUE%\\r\\n", "string", "", "", "")
	return meta.Marshal()
}

func (c *UnitTcpControl) InternalUnitStart() error {
	var err error

	err = json.Unmarshal([]byte(c.GetConfig()), &c.config)
	if err != nil {
		err = errors.New("config error")
		c.SetString(ItemNameStatus, err.Error(), "error")
		return err
	}

	if c.config.Addr == "" {
		err = errors.New("wrong address")
		c.SetString(ItemNameStatus, err.Error(), "error")
		return err
	}

	if c.config.Timeout < 100 {
		err = errors.New("wrong timeout (<100)")
		c.SetString(ItemNameStatus, err.Error(), "error")
		return err
	}
	if c.config.Timeout > 10000 {
		err = errors.New("wrong timeout (>10000)")
		c.SetString(ItemNameStatus, err.Error(), "error")
		return err
	}

	if c.config.Period < 100 {
		err = errors.New("wrong period")
		c.SetString(ItemNameStatus, err.Error(), "error")
		return err
	}

	if c.config.Period < c.config.Timeout {
		err = errors.New("wrong period (<timeout)")
		c.SetString(ItemNameStatus, err.Error(), "error")
		return err
	}
	if c.config.Period < 100 {
		err = errors.New("wrong period (<100)")
		c.SetString(ItemNameStatus, err.Error(), "error")
		return err
	}
	if c.config.Period > 60000 {
		err = errors.New("wrong period (>60000)")
		c.SetString(ItemNameStatus, err.Error(), "error")
		return err
	}

	c.SetMainItem(ItemNameStatus)

	c.SetString(ItemNameStatus, "", "")

	go c.Tick()
	return nil
}

func (c *UnitTcpControl) InternalUnitStop() {
	c.Stopping = true
}

func (c *UnitTcpControl) Tick() {
	var err error
	var conn net.Conn

	c.Started = true
	dtLastTime := time.Now().UTC()

	for !c.Stopping {
		for {
			if c.Stopping || time.Now().Sub(dtLastTime) > time.Duration(c.config.Period)*time.Millisecond {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
		if c.Stopping {
			c.SetString(ItemNameStatus, "stopped", "")
			break
		}
		dtLastTime = time.Now().UTC()

		timeBegin := time.Now()
		if conn == nil {
			conn, err = net.DialTimeout("tcp", c.config.Addr, time.Duration(c.config.Timeout)*time.Millisecond)
		}
		timeEnd := time.Now()
		duration := timeEnd.Sub(timeBegin)

		for _, item := range c.config.Items {
			val, err := c.GetItem(item.Name)
			if err == nil {
				format := item.Format
				format = strings.ReplaceAll(format, "\\r", "\r")
				format = strings.ReplaceAll(format, "\\n", "\n")
				format = strings.ReplaceAll(format, "%VALUE%", val.Value)
				format = strings.ReplaceAll(format, "%DT%", time.Unix(0, val.DT*1000).Format("2006-01-02 15:04:05.000"))
				format = strings.ReplaceAll(format, "%UOM%", val.UOM)
				dataStr := format
				if conn != nil {
					c.SetString("string", dataStr, "")
					_, err = conn.Write([]byte(dataStr))
					if err != nil {
						conn.Close()
						conn = nil
					}
				}
			}
		}

		if err != nil {
			c.SetString(ItemNameStatus, "timeout", "error")
		} else {
			if !c.Stopping {
				c.SetInt(ItemNameStatus, int(duration.Milliseconds()), "ms")
				c.SetError("")
			} else {
				c.SetError("")
			}
		}
	}

	if conn != nil {
		conn.Close()
	}

	c.SetString(ItemNameStatus, "", "stopped")
	c.Started = false
}
