package unit_network_interface

import (
	"encoding/json"
	"errors"
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/resources"
	"github.com/gazercloud/gazernode/system/units/units_common"
	"github.com/gazercloud/gazernode/utilities/logger"
	"github.com/kbinani/win"
	"net"
	"time"
)

type UnitNetworkInterface struct {
	units_common.Unit

	interfaceName string
}

var Image []byte

func init() {
	Image = resources.R_files_sensors_unit_computer_network_interface_png
}

func New() common_interfaces.IUnit {
	var c UnitNetworkInterface
	return &c
}

func (c *UnitNetworkInterface) GetConfigMeta() string {
	meta := units_common.NewUnitConfigItem("", "", "", "", "", "", "")
	meta.Add("interface_name", "Network Interface", "-", "string", "", "", "network_interface")
	meta.Add("period", "Period, ms", "1000", "num", "0", "999999", "0")
	return meta.Marshal()
}

type NetworkInterface struct {
	Name string
	Id   int
}

func (c *UnitNetworkInterface) InternalUnitStart() error {
	var err error
	type Config struct {
		InterfaceName string  `json:"interface_name"`
		Period        float64 `json:"period"`
	}

	var config Config
	err = json.Unmarshal([]byte(c.GetConfig()), &config)
	if err != nil {
		logger.Println("ERROR[UnitNetworkInterface]:", err)
		err = errors.New("config error")
		c.SetString("TotalSpeed", err.Error(), "error")
		return err
	}

	c.interfaceName = config.InterfaceName

	c.SetString("TotalSpeed", "", "")
	c.SetMainItem("TotalSpeed")

	c.SetString("Status", "", "")
	c.SetString("InOctets", "", "")
	c.SetString("InUcastPkts", "", "")
	c.SetString("InNUcastPkts", "", "")
	c.SetString("InDiscards", "", "")
	c.SetString("InErrors", "", "")

	c.SetString("OutOctets", "", "")
	c.SetString("OutUcastPkts", "", "")
	c.SetString("OutNUcastPkts", "", "")
	c.SetString("OutDiscards", "", "")
	c.SetString("OutErrors", "", "")

	c.SetString("InSpeedFrames", "", "")
	c.SetString("InSpeed", "", "")
	c.SetString("OutSpeedFrames", "", "")
	c.SetString("OutSpeed", "", "")

	go c.Tick()
	return nil
}

func (c *UnitNetworkInterface) InternalUnitStop() {
}

func (c *UnitNetworkInterface) Tick() {
	var err error
	c.Started = true

	type LastCounters struct {
		DT            time.Time
		TotalIn       uint64
		TotalOut      uint64
		TotalInBytes  uint64
		TotalOutBytes uint64
	}

	lastCounters := make(map[int]LastCounters)

	for !c.Stopping {
		for i := 0; i < 10; i++ {
			if c.Stopping {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}

		totalSpeed := 0.0
		totalInSpeed := 0.0
		totalOutSpeed := 0.0
		var interfaces []net.Interface
		interfaces, err = net.Interfaces()
		if err == nil {
			for _, ni := range interfaces {
				if ni.Name != c.interfaceName {
					continue
				}
				var table win.MIB_IF_ROW2
				table.InterfaceIndex = win.NET_IFINDEX(ni.Index)
				win.GetIfEntry2(&table)

				c.SetInt("Index", ni.Index, "")
				c.SetString("HardwareAddress", ni.HardwareAddr.String(), "")

				if table.OperStatus == 1 {
					c.SetUInt64("Status", uint64(table.OperStatus), "")
					c.SetUInt64("InOctets", uint64(table.InOctets), "bytes")
					c.SetUInt64("InUcastPkts", uint64(table.InUcastPkts), "")
					c.SetUInt64("InNUcastPkts", uint64(table.InNUcastPkts), "")
					c.SetUInt64("InDiscards", uint64(table.InDiscards), "")
					c.SetUInt64("InErrors", uint64(table.InErrors), "bytes")

					c.SetUInt64("OutOctets", uint64(table.OutOctets), "bytes")
					c.SetUInt64("OutUcastPkts", uint64(table.OutUcastPkts), "")
					c.SetUInt64("OutNUcastPkts", uint64(table.OutNUcastPkts), "")
					c.SetUInt64("OutDiscards", uint64(table.OutDiscards), "")
					c.SetUInt64("OutErrors", uint64(table.OutErrors), "bytes")

					totalIn := uint64(table.InUcastPkts) + uint64(table.InNUcastPkts) + uint64(table.InDiscards)
					totalInBytes := uint64(table.InOctets)
					totalOut := uint64(table.OutUcastPkts) + uint64(table.OutNUcastPkts) + uint64(table.OutDiscards)
					totalOutBytes := uint64(table.OutOctets)

					nowTime := time.Now()

					if cs, ok := lastCounters[ni.Index]; ok {
						seconds := nowTime.Sub(cs.DT).Seconds()
						if seconds > 0.001 {
							c.SetFloat64("InSpeedFrames", float64(totalIn-cs.TotalIn)/seconds, "", 1)
							c.SetFloat64("InSpeed", float64(totalInBytes-cs.TotalInBytes)/seconds/1024.0, "KB/sec", 1)
							c.SetFloat64("OutSpeedFrames", float64(totalOut-cs.TotalOut)/seconds, "", 1)
							c.SetFloat64("OutSpeed", float64(totalOutBytes-cs.TotalOutBytes)/seconds/1024.0, "KB/sec", 1)

							totalInSpeed += float64(totalInBytes-cs.TotalInBytes) / seconds / 1024.0
							totalOutSpeed += float64(totalOutBytes-cs.TotalOutBytes) / seconds / 1024.0
						}
					}

					lastCounters[ni.Index] = LastCounters{
						DT:            nowTime,
						TotalIn:       totalIn,
						TotalOut:      totalOut,
						TotalInBytes:  totalInBytes,
						TotalOutBytes: totalOutBytes,
					}
				} else {
					delete(lastCounters, ni.Index)
					c.SetString("Status", "", "error")
					c.SetString("InOctets", "", "error")
					c.SetString("InUcastPkts", "", "error")
					c.SetString("InNUcastPkts", "", "error")
					c.SetString("InDiscards", "", "error")
					c.SetString("InErrors", "", "error")

					c.SetString("OutOctets", "", "error")
					c.SetString("OutUcastPkts", "", "error")
					c.SetString("OutNUcastPkts", "", "error")
					c.SetString("OutDiscards", "", "error")
					c.SetString("OutErrors", "", "error")

					c.SetString("InSpeedFrames", "", "error")
					c.SetString("OutSpeedFrames", "", "error")

					c.SetString("InSpeed", "", "error")
					c.SetString("OutSpeed", "", "error")
				}
			}

			totalSpeed = totalInSpeed + totalOutSpeed
			c.SetFloat64("TotalSpeed", totalSpeed, "KB/sec", 1)
		} else {
			c.SetError(err.Error())
		}
	}

	c.Started = false
}
