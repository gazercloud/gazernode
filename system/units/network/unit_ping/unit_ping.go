package unit_ping

import (
	"encoding/json"
	"errors"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/resources"
	"github.com/gazercloud/gazernode/system/units/units_common"
	"github.com/gazercloud/gazernode/utilities"
	"github.com/gazercloud/gazernode/utilities/gazerping"
	"github.com/gazercloud/gazernode/utilities/uom"
	"math"
	"net"
	"runtime"
	"time"
)

type UnitPing struct {
	units_common.Unit

	addr      string
	timeoutMs int
	periodMs  int
	frameSize int
}

func New() common_interfaces.IUnit {
	var c UnitPing
	return &c
}

const (
	ItemNameAddress  = "Address"
	ItemNameTime     = "Time"
	ItemNameIP       = "IP"
	ItemNameDataSize = "DataSize"
)

var Image []byte

func init() {
	Image = resources.R_files_sensors_unit_network_ping_png
}

func (c *UnitPing) GetConfigMeta() string {
	meta := units_common.NewUnitConfigItem("", "", "", "", "", "", "")
	pAddr := meta.Add("addr", "Address", "localhost", "string", "", "", "")
	pAddr.ItemIsDisplayName = true
	meta.Add("period", "Period, ms", "1000", "num", "100", "60000", "0")
	meta.Add("timeout", "Timeout, ms", "1000", "num", "100", "10000", "0")
	meta.Add("frame_size", "Frame Size, bytes", "64", "num", "4", "1400", "0")
	return meta.Marshal()
}

func (c *UnitPing) InternalUnitStart() error {
	var err error
	c.SetMainItem(ItemNameTime)

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

	c.timeoutMs = int(math.Round(config.Timeout))
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

	c.periodMs = int(math.Round(config.Period))
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

	c.frameSize = int(math.Round(config.FrameSize))
	if c.frameSize < 1 {
		err = errors.New("wrong Frame Size (<1)")
		c.SetString(ItemNameTime, err.Error(), "error")
		return err
	}
	if c.frameSize > 1400 {
		err = errors.New("wrong FrameSize (>1400)")
		c.SetString(ItemNameTime, err.Error(), "error")
		return err
	}

	c.SetString(ItemNameAddress, c.addr, "-")
	//c.SetString(ItemNameTime, "", "")
	//c.SetString(ItemNameIP, "", "-")
	c.SetInt(ItemNameDataSize, c.frameSize, uom.BYTES)

	c.SetPropertyIfDoesntExist(ItemNameAddress, "color", "#AA0000")

	go c.Tick()
	return nil
}

func (c *UnitPing) InternalUnitStop() {
}

func (c *UnitPing) Tick() {
	var lastError string
	var lastIP string

	c.SetString(ItemNameTime, "", "started")

	c.Started = true
	var dtLastPingTime time.Time
	for !c.Stopping {
		for {
			if c.Stopping || time.Now().Sub(dtLastPingTime) > time.Duration(c.periodMs)*time.Millisecond {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
		if c.Stopping {
			break
		}

		if !utilities.IsRoot() && runtime.GOOS == "linux" {
			c.SetString(ItemNameTime, "no root", "error")
			c.SetError("ping.NewPinger: " + "no root")
			dtLastPingTime = time.Now().UTC()
			continue
		}

		//var timeoutMSec int32 = int32(c.timeoutMs)
		//var frameSize int32 = int32(c.frameSize)

		if c.addr == "" {
			c.SetError("ipaddress == ''")
			c.SetString(ItemNameTime, "wrong address", "error")
			continue
		}

		//logger.Println("PING 1 ", c.addr)

		var err error
		/*pingObject, err := ping.NewPinger(c.addr)
		if err != nil {
			c.SetString(ItemNameTime, err.Error(), "error")
			c.SetError("ping.NewPinger: " + err.Error())
			dtLastPingTime = time.Now().UTC()
			continue
		}

		if utilities.IsRoot() || runtime.GOOS == "windows" {
			pingObject.SetPrivileged(true)
		}

		pingObject.Count = 1
		pingObject.Size = int(frameSize)
		pingObject.Timeout = time.Duration(timeoutMSec) * time.Millisecond
		*/

		var pingTime int
		var peer net.Addr

		pingTime, peer, err = gazerping.Ping(c.addr, c.frameSize, c.timeoutMs)

		if err == nil {
			ip := peer.String()
			if ip != lastIP {
				lastIP = ip
				c.SetString(ItemNameIP, ip, "-")
			}
			if !c.Stopping {
				t := pingTime
				c.SetInt(ItemNameTime, t, uom.MS)
				if lastError != "" {
					c.SetError("")
				}
				lastError = ""
			}
		} else {
			if lastError != err.Error() {
				lastError = err.Error()
				lastIP = ""
				c.SetError(lastError)
				c.SetString(ItemNameIP, lastIP, "error")
			}
			c.SetString(ItemNameTime, lastError, "error")
		}

		dtLastPingTime = time.Now().UTC()
	}
	c.SetString(ItemNameTime, "", "stopped")
	c.Started = false
}
