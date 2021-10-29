package unit_tcp_connect

import (
	"encoding/json"
	"errors"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/resources"
	"github.com/gazercloud/gazernode/system/units/units_common"
	"github.com/gazercloud/gazernode/utilities/uom"
	"net"
	"strings"
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
	ItemNameAddress = "Address"
	ItemNameTime    = "Time"
	ItemNameIP      = "IP"
)

var Image []byte

func init() {
	Image = resources.R_files_sensors_unit_network_tcp_connect_png
}

func (c *UnitTcpConnect) GetConfigMeta() string {
	meta := units_common.NewUnitConfigItem("", "", "", "", "", "", "")
	meta.Add("addr", "Address", "localhost:445", "string", "", "", "")
	meta.Add("period", "Period, ms", "1000", "num", "0", "999999", "0")
	meta.Add("timeout", "Timeout, ms", "1000", "num", "0", "999999", "0")
	return meta.Marshal()
}

func (c *UnitTcpConnect) InternalUnitStart() error {
	var err error

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
		c.SetString(ItemNameTime, err.Error(), uom.ERROR)
		return err
	}

	c.addr = config.Addr
	if c.addr == "" {
		err = errors.New("wrong address")
		c.SetString(ItemNameTime, err.Error(), uom.ERROR)
		return err
	}

	c.timeoutMs = int(config.Timeout)
	if c.timeoutMs < 100 {
		err = errors.New("wrong timeout (<100)")
		c.SetString(ItemNameTime, err.Error(), uom.ERROR)
		return err
	}
	if c.timeoutMs > 10000 {
		err = errors.New("wrong timeout (>10000)")
		c.SetString(ItemNameTime, err.Error(), uom.ERROR)
		return err
	}

	c.periodMs = int(config.Period)
	if c.periodMs < 100 {
		err = errors.New("wrong period")
		c.SetString(ItemNameTime, err.Error(), uom.ERROR)
		return err
	}

	if c.periodMs < c.timeoutMs {
		err = errors.New("wrong period (<timeout)")
		c.SetString(ItemNameTime, err.Error(), uom.ERROR)
		return err
	}
	if c.periodMs < 100 {
		err = errors.New("wrong period (<100)")
		c.SetString(ItemNameTime, err.Error(), uom.ERROR)
		return err
	}
	if c.periodMs > 60000 {
		err = errors.New("wrong period (>60000)")
		c.SetString(ItemNameTime, err.Error(), uom.ERROR)
		return err
	}

	c.SetMainItem(ItemNameTime)

	c.SetString(ItemNameTime, "", "")
	c.SetString(ItemNameAddress, c.addr, uom.EVENT)
	c.SetString(ItemNameIP, "", uom.EVENT)
	go c.Tick()
	return nil
}

func (c *UnitTcpConnect) InternalUnitStop() {
}

func (c *UnitTcpConnect) Tick() {
	var lastError string
	var lastIP string

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
			c.SetString(ItemNameTime, "", uom.STOPPED)
			break
		}
		dtLastTime = time.Now().UTC()

		var err error
		var conn net.Conn
		var resolvedAddr *net.IPAddr

		posOfColon := strings.Index(c.addr, ":")
		cleanAddr := c.addr
		if posOfColon > -1 {
			cleanAddr = c.addr[:posOfColon]
		}

		resolvedAddr, err = net.ResolveIPAddr("", cleanAddr)
		if err == nil {
			ip := resolvedAddr.IP.String()
			if ip != lastIP {
				lastIP = ip
				c.SetString(ItemNameIP, ip, uom.EVENT)
			}

			timeBegin := time.Now()
			conn, err = net.DialTimeout("tcp", c.addr, time.Duration(c.timeoutMs)*time.Millisecond)
			timeEnd := time.Now()
			duration := timeEnd.Sub(timeBegin)
			if err == nil {
				c.SetInt(ItemNameTime, int(duration.Milliseconds()), uom.MS)
			}

			if conn != nil {
				conn.Close()
			}
		}

		if err != nil {
			if lastError != err.Error() {
				lastError = err.Error()
				lastIP = ""
				c.SetError(lastError)
				c.SetString(ItemNameIP, lastIP, uom.ERROR)
			}
			c.SetString(ItemNameTime, lastError, uom.ERROR)
		} else {
			if !c.Stopping {
				if lastError != "" {
					c.SetError("")
				}
				lastError = ""
			}
		}
	}
	c.SetString(ItemNameTime, "", uom.STOPPED)
	c.Started = false
}
