package unit_network

import (
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/resources"
	"github.com/gazercloud/gazernode/system/units/units_common"
	"io/ioutil"
	"net"
	"strconv"
	"time"
)

type UnitNetwork struct {
	units_common.Unit
}

var Image []byte

func init() {
	Image = resources.R_files_sensors_sensor_windows_ram_png
}

func New() common_interfaces.IUnit {
	var c UnitNetwork
	return &c
}

func (c *UnitNetwork) InternalUnitStart() error {
	c.SetString("TotalSpeed", "", "")
	c.SetMainItem("TotalSpeed")

	interfaces, err := net.Interfaces()
	if err == nil {
		for _, ni := range interfaces {
			c.SetString(ni.Name+"/InSpeed", "", "")
			c.SetString(ni.Name+"/OutSpeed", "", "")
		}
		c.SetString("TotalInSpeed", "", "")
		c.SetString("TotalOutSpeed", "", "")
	} else {
		c.SetError("")
	}

	go c.Tick()
	return nil
}

func (c *UnitNetwork) InternalUnitStop() {
}

func (c *UnitNetwork) GetConfigMeta() string {
	return ""
}

func (c *UnitNetwork) Tick() {
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
				rxPackets := 0
				rxBytes := 0
				txPackets := 0
				txBytes := 0

				rxPacketsStr, errParamRxPackets := ioutil.ReadFile("/sys/class/net/" + ni.Name + "/statistics/rx_packets")
				if errParamRxPackets == nil {
					rxPackets, errParamRxPackets = strconv.ParseInt(rxPacketsStr, 10, 64)
				}

				rxBytesStr, errParamRxBytes := ioutil.ReadFile("/sys/class/net/" + ni.Name + "/statistics/rx_bytes")
				if errParamRxBytes == nil {
					rxBytes, errParamRxBytes = strconv.ParseInt(rxBytesStr, 10, 64)
				}

				txPacketsStr, errParamTxPackets := ioutil.ReadFile("/sys/class/net/" + ni.Name + "/statistics/tx_packets")
				if errParamTxPackets == nil {
					txPackets, errParamTxPackets = strconv.ParseInt(txPacketsStr, 10, 64)
				}

				txBytesStr, errParamTxBytes := ioutil.ReadFile("/sys/class/net/" + ni.Name + "/statistics/tx_bytes")
				if errParamTxBytes == nil {
					txBytes, errParamTxBytes = strconv.ParseInt(txBytesStr, 10, 64)
				}

				totalIn := uint64(rxPackets)
				totalInBytes := uint64(rxBytes)
				totalOut := uint64(txPackets)
				totalOutBytes := uint64(txBytes)

				nowTime := time.Now()

				if table.OperStatus == 1 {
					if cs, ok := lastCounters[ni.Index]; ok {
						seconds := nowTime.Sub(cs.DT).Seconds()
						if seconds > 0.001 {
							c.SetFloat64(ni.Name+"/InSpeed", float64(totalInBytes-cs.TotalInBytes)/seconds/1024.0, "KB/sec", 1)
							c.SetFloat64(ni.Name+"/OutSpeed", float64(totalOutBytes-cs.TotalOutBytes)/seconds/1024.0, "KB/sec", 1)
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
					c.SetString(ni.Name+"/InSpeed", "", "error")
					c.SetString(ni.Name+"/OutSpeed", "", "error")
				}

			}

			totalSpeed = totalInSpeed + totalOutSpeed
			c.SetFloat64("TotalInSpeed", totalInSpeed, "KB/sec", 1)
			c.SetFloat64("TotalOutSpeed", totalOutSpeed, "KB/sec", 1)
			c.SetFloat64("TotalSpeed", totalSpeed, "KB/sec", 1)
		} else {
			c.SetError(err.Error())
		}
	}

	c.Started = false
}
