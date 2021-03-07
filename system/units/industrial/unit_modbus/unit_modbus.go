package unit_modbus

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/system/units/units_common"
	"net"
	"time"
)

type Item struct {
	Name string  `json:"item_name"`
	Addr float64 `json:"addr"`
}

type Config struct {
	Addr    string `json:"addr"`
	Timeout float64
	Period  float64
	Items   []Item `json:"items"`
}

type UnitModbus struct {
	units_common.Unit
	config Config
}

func New() common_interfaces.IUnit {
	var c UnitModbus
	return &c
}

const (
	ItemNameStatus = "Status"
)

var Image []byte

func init() {
}

func (c *UnitModbus) GetConfigMeta() string {
	meta := units_common.NewUnitConfigItem("", "", "", "", "", "", "")
	meta.Add("addr", "Address", "localhost:502", "string", "", "", "")
	tChannels := meta.Add("channels", "Channels", "", "table", "", "", "")
	tChannels.Add("period", "Period, ms", "1000", "num", "0", "999999", "")
	tChannels.Add("timeout", "Timeout, ms", "1000", "num", "0", "999999", "")
	t1 := tChannels.Add("items", "Items", "", "table", "", "", "")
	t1.Add("item_name", "Item Name", "Unit/Item", "string", "", "", "data-item")
	t1.Add("addr", "Address", "0", "num", "0", "65535", "")
	return meta.Marshal()
}

func (c *UnitModbus) InternalUnitStart() error {
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

func (c *UnitModbus) InternalUnitStop() {
	c.Stopping = true
}

func (c *UnitModbus) WriteCoil(addr uint16, on bool) error {
	var err error
	var conn net.Conn
	conn, err = net.DialTimeout("tcp", c.config.Addr, time.Duration(c.config.Timeout)*time.Millisecond)

	data := []byte{0x00, 0x00, 0x00, 0x00, 0x0, 0x6, 1, 0x05, 0x00, 0x00, 0xFF, 0x00}
	if on {
		data = []byte{0x00, 0x00, 0x00, 0x00, 0x0, 0x6, 1, 0x05, 0x00, 0x00, 0xFF, 0x00}
	} else {
		data = []byte{0x00, 0x00, 0x00, 0x00, 0x0, 0x6, 1, 0x05, 0x00, 0x00, 0x00, 0x00}
	}

	binary.BigEndian.PutUint16(data[8:], addr)

	if conn != nil {
		_, err = conn.Write(data)
		conn.Close()
	}
	conn = nil
	return err
}

func (c *UnitModbus) Tick() {
	var err error

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
		timeEnd := time.Now()
		duration := timeEnd.Sub(timeBegin)

		for _, item := range c.config.Items {
			val, err := c.GetItem(item.Name)
			if err == nil {
				c.WriteCoil(uint16(item.Addr), val.Value == "1")
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

	c.SetString(ItemNameStatus, "", "stopped")
	c.Started = false
}
