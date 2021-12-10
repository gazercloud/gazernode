package unit_modbus

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/system/units/units_common"
	"github.com/gazercloud/gazernode/utilities/logger"
	"net"
	"sync"
	"time"
)

type Item struct {
	Name string  `json:"item_name"`
	Addr float64 `json:"addr"`
	Type float64 `json:"type"`
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

	mtxConn sync.Mutex
	conn    net.Conn
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
	meta.Add("period", "Period, ms", "1000", "num", "0", "999999", "")
	meta.Add("timeout", "Timeout, ms", "1000", "num", "0", "999999", "")
	t1 := meta.Add("items", "Items", "", "table", "", "", "")
	t1.Add("item_name", "Item Name", "Unit/Item", "string", "", "", "data-item")
	t1.Add("addr", "Address", "0", "num", "0", "65535", "")
	t1.Add("type", "Type", "0", "num", "0", "65535", "")
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
	go c.readThread()
	return nil
}

func (c *UnitModbus) InternalUnitStop() {
	c.Stopping = true
}

func (c *UnitModbus) WriteCoil(conn net.Conn, addr uint16, on bool) error {
	var err error
	data := []byte{0x00, 0x00, 0x00, 0x00, 0x0, 0x6, 1, 0x05, 0x00, 0x00, 0xFF, 0x00}
	if on {
		data = []byte{0x00, 0x00, 0x00, 0x00, 0x0, 0x6, 1, 0x05, 0x00, 0x00, 0xFF, 0x00}
	} else {
		data = []byte{0x00, 0x00, 0x00, 0x00, 0x0, 0x6, 1, 0x05, 0x00, 0x00, 0x00, 0x00}
	}
	binary.BigEndian.PutUint16(data[8:], addr)
	var n int
	for n < len(data) {
		var writeResult int
		writeResult, err = conn.Write(data[n:])
		if err != nil {
			break
		}
		n += writeResult
	}
	return err
}

func (c *UnitModbus) ReadCoil(conn net.Conn, addr uint16) error {
	var err error
	data := []byte{0x00, 0x00, 0x00, 0x00, 0x0, 0x6, 1, 0x02, 0x00, 0x00, 0x00, 0x08}
	//binary.BigEndian.PutUint16(data[8:], addr)
	//binary.BigEndian.PutUint16(data[10:], addr)
	var n int
	for n < len(data) {
		var writeResult int
		writeResult, err = conn.Write(data[n:])
		if err != nil {
			break
		}
		n += writeResult
	}
	return err
}

func (c *UnitModbus) readThread() {
	var err error
	inputBuffer := make([]byte, 0)
	logger.Println("modbus - start")

	for !c.Stopping {
		c.mtxConn.Lock()
		connection := c.conn
		c.mtxConn.Unlock()
		if connection == nil {
			time.Sleep(10 * time.Millisecond)
			continue
		}
		buffer := make([]byte, 1024)
		var n int
		//logger.Println("modbus - reading ...")
		n, err = connection.Read(buffer)
		if err != nil {
			c.disconnect()
			continue
		}
		//logger.Println("modbus - read", n)
		if n > 0 {
			inputBuffer = append(inputBuffer, buffer[:n]...)
		}

		// parsing
		for len(inputBuffer) > 6 {
			lenOfFrame := int(binary.BigEndian.Uint16(inputBuffer[4:]))
			if len(inputBuffer) >= lenOfFrame+6 {
				//fmt.Println("MODBUS: ", inputBuffer)
				fmt.Println("frame: ", inputBuffer[0:lenOfFrame+6])
				if inputBuffer[2] == 0 && inputBuffer[3] == 0 {
					if inputBuffer[7] == 0x02 {
						fmt.Println("INPUT MODBUS!", inputBuffer[8:lenOfFrame+6])
						if inputBuffer[9] == 0 {
							c.SetString("input", "0", "")
						} else {
							c.SetString("input", "1", "")
						}

					}
				}
				inputBuffer = inputBuffer[lenOfFrame+6:]
			} else {
				break
			}
		}
	}

	logger.Println("modbus - stop")
}

func (c *UnitModbus) connect() error {
	var err error
	c.mtxConn.Lock()
	c.conn, err = net.DialTimeout("tcp", c.config.Addr, time.Duration(c.config.Timeout)*time.Millisecond)
	if err != nil {
		c.SetString(ItemNameStatus, err.Error(), "error")
		c.conn = nil
	}
	c.mtxConn.Unlock()
	return err
}

func (c *UnitModbus) disconnect() {
	c.mtxConn.Lock()
	if c.conn != nil {
		_ = c.conn.Close()
		c.conn = nil
	}
	c.mtxConn.Unlock()
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

		if c.conn == nil {
			if c.connect() != nil {
				continue
			}
		}

		timeBegin := time.Now()

		for _, item := range c.config.Items {
			if item.Type == 0 {
				var val common_interfaces.ItemValue
				val, err = c.GetItem(item.Name)
				if err == nil {
					err = c.WriteCoil(c.conn, uint16(item.Addr), val.Value == "1")
				}
			}
		}

		for _, item := range c.config.Items {
			if item.Type == 1 {
				if err == nil {
					err = c.ReadCoil(c.conn, uint16(item.Addr))
				}
			}
		}

		timeEnd := time.Now()
		duration := timeEnd.Sub(timeBegin)

		if err != nil {
			c.SetString(ItemNameStatus, err.Error(), "error")
			c.disconnect()
		} else {
			if !c.Stopping {
				c.SetInt(ItemNameStatus, int(duration.Milliseconds()), "ms")
				c.SetError("")
			} else {
				c.SetError("")
			}
		}
	}

	c.disconnect()
	c.SetString(ItemNameStatus, "", "stopped")
	c.Started = false
}
