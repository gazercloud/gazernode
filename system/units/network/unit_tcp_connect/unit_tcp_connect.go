package unit_tcp_connect

import (
	"encoding/json"
	"errors"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/resources"
	"github.com/gazercloud/gazernode/system/units/units_common"
	"net"
	"time"
)

type UnitTcpConnect struct {
	units_common.Unit
	addr      string
	timeoutMs int
	periodMs  int
}

func New() common_interfaces.IUnit {
	var c UnitTcpConnect
	return &c
}

const (
	ItemNameTime = "Time"
	ItemNameAddr = "Address"
)

var Image []byte

func init() {
	Image = resources.R_files_sensors_sensor_network_tcp_connect_png
}

func (c *UnitTcpConnect) GetConfigMeta() string {
	meta := units_common.NewUnitConfigItem("", "", "", "", "", "", "")
	meta.Add("addr", "Address", "localhost:445", "string", "", "", "")
	meta.Add("period", "Period, ms", "1000", "num", "0", "999999", "")
	meta.Add("timeout", "Timeout, ms", "1000", "num", "0", "999999", "")
	return meta.Marshal()
}

func (c *UnitTcpConnect) InternalUnitStart() error {
	var err error
	c.addr = "r002.gazer.cloud:80"

	type Config struct {
		Addr      string  `json:"addr"`
		Timeout   float64 `json:"timeout"`
		Period    float64 `json:"period"`
		FrameSize float64 `json:"frame_size"`
	}

	var config Config
	err = json.Unmarshal([]byte(c.GetConfig()), &config)
	if err != nil {
		err = errors.New("config error")
		c.SetString(ItemNameTime, err.Error(), "error")
		return err
	}

	c.addr = config.Addr
	if c.addr == "" {
		err = errors.New("wrong address")
		c.SetString(ItemNameTime, err.Error(), "error")
		return err
	}

	c.timeoutMs = int(config.Timeout)
	if c.timeoutMs < 100 {
		err = errors.New("wrong timeout (<100)")
		c.SetString(ItemNameTime, err.Error(), "error")
		return err
	}
	if c.timeoutMs > 10000 {
		err = errors.New("wrong timeout (>10000)")
		c.SetString(ItemNameTime, err.Error(), "error")
		return err
	}

	c.periodMs = int(config.Period)
	if c.periodMs < 100 {
		err = errors.New("wrong period")
		c.SetString(ItemNameTime, err.Error(), "error")
		return err
	}

	if c.periodMs < c.timeoutMs {
		err = errors.New("wrong period (<timeout)")
		c.SetString(ItemNameTime, err.Error(), "error")
		return err
	}
	if c.periodMs < 100 {
		err = errors.New("wrong period (<100)")
		c.SetString(ItemNameTime, err.Error(), "error")
		return err
	}
	if c.periodMs > 60000 {
		err = errors.New("wrong period (>60000)")
		c.SetString(ItemNameTime, err.Error(), "error")
		return err
	}

	c.SetMainItem(ItemNameTime)

	c.SetStringService(ItemNameAddr, c.addr, "")
	c.SetString(ItemNameTime, "", "")

	go c.Tick()
	return nil
}

func (c *UnitTcpConnect) InternalUnitStop() {
}

func (c *UnitTcpConnect) Tick() {
	c.Started = true
	dtLastTime := time.Now().UTC()

	for !c.Stopping {
		for {
			if c.Stopping || time.Now().Sub(dtLastTime) > time.Duration(c.periodMs)*time.Millisecond {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
		if c.Stopping {
			c.SetString(ItemNameTime, "stopped", "")
			break
		}
		dtLastTime = time.Now().UTC()

		var err error
		var conn net.Conn

		timeBegin := time.Now()
		conn, err = net.DialTimeout("tcp", c.addr, time.Duration(c.timeoutMs)*time.Millisecond)
		timeEnd := time.Now()
		duration := timeEnd.Sub(timeBegin)

		if conn != nil {
			conn.Close()
		}

		if err != nil {
			c.SetString(ItemNameTime, "timeout", "error")
		} else {
			if !c.Stopping {
				c.SetInt(ItemNameTime, int(duration.Milliseconds()), "ms")
				c.SetError("")
			} else {
				c.SetError("")
			}
		}
	}
	c.SetString(ItemNameTime, "", "stopped")
	c.Started = false
}
