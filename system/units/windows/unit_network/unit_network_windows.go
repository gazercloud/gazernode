package unit_network

import (
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/resources"
	"github.com/gazercloud/gazernode/system/units/units_common"
	"github.com/kbinani/win"
	"net"
	"time"
)

type UnitNetwork struct {
	units_common.Unit
}

var Image []byte

func init() {
	Image = resources.R_files_sensors_unit_computer_network_png
}

func New() common_interfaces.IUnit {
	var c UnitNetwork
	return &c
}

func (c *UnitNetwork) InternalUnitStart() error {
	//c.SetString("TotalSpeed", "", "")
	c.SetMainItem("TotalSpeed")

	/*interfaces, err := net.Interfaces()
	if err == nil {
		for _, ni := range interfaces {
			var table win.MIB_IF_ROW2
			table.InterfaceIndex = win.NET_IFINDEX(ni.Index)
			win.GetIfEntry2(&table)
			if table.Type == 24 {
				continue
			}
			c.SetString(ni.Name+"/InSpeed", "", "")
			c.SetString(ni.Name+"/OutSpeed", "", "")
		}
		c.SetString("TotalInSpeed", "", "")
		c.SetString("TotalOutSpeed", "", "")
	} else {
		c.SetError("")
	}*/

	go c.Tick()
	return nil
}

func (c *UnitNetwork) InternalUnitStop() {
}

func (c *UnitNetwork) GetConfigMeta() string {
	meta := units_common.NewUnitConfigItem("", "", "", "", "", "", "")
	return meta.Marshal()
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
				var table win.MIB_IF_ROW2
				table.InterfaceIndex = win.NET_IFINDEX(ni.Index)
				win.GetIfEntry2(&table)

				if table.Type == 24 {
					continue
				}

				totalIn := uint64(table.InUcastPkts) + uint64(table.InNUcastPkts) + uint64(table.InDiscards)
				totalInBytes := uint64(table.InOctets)
				totalOut := uint64(table.OutUcastPkts) + uint64(table.OutNUcastPkts) + uint64(table.OutDiscards)
				totalOutBytes := uint64(table.OutOctets)

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
