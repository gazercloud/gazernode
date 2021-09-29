package unit_storage

import (
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/resources"
	"github.com/gazercloud/gazernode/system/units/units_common"
	"github.com/gazercloud/gazernode/utilities/uom"
	"golang.org/x/sys/windows"
	"strconv"
	"syscall"
	"time"
)

type UnitStorage struct {
	units_common.Unit
}

var Image []byte

func init() {
	Image = resources.R_files_sensors_unit_computer_storage_png
}

var (
	modpsapi    = windows.NewLazySystemDLL("psapi.dll")
	modkernel32 = windows.NewLazySystemDLL("kernel32.dll")
	moduser32   = windows.NewLazySystemDLL("user32.dll")
)

var getProcessMemoryInfo = modpsapi.NewProc("GetProcessMemoryInfo")
var getProcessHandleCount = modkernel32.NewProc("GetProcessHandleCount")
var getProcessIoCounters = modkernel32.NewProc("GetProcessIoCounters")
var getGuiResources = moduser32.NewProc("GetGuiResources")

var getLogicalDrivesHandle = modkernel32.NewProc("GetLogicalDrives")

func New() common_interfaces.IUnit {
	var c UnitStorage
	return &c
}

func (c *UnitStorage) bitsToDrives(bitMap uint32) (drives []string) {
	availableDrives := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

	for i := range availableDrives {
		if bitMap&1 == 1 {
			drives = append(drives, availableDrives[i])
		}
		bitMap >>= 1
	}

	return
}

func (c *UnitStorage) drives() []string {
	drives := make([]string, 0)
	drivesBits, err := windows.GetLogicalDrives()
	if err == nil {
		drives = c.bitsToDrives(drivesBits)
	}
	return drives
}

func (c *UnitStorage) InternalUnitStart() error {
	drives := c.drives()
	c.SetString("UsedPercents", "", uom.STARTED)
	c.SetMainItem("UsedPercents")

	for _, disk := range drives {
		c.SetString(disk+"/Total", "", uom.STARTED)
		//c.SetString(disk+"/Available", "", "")
		c.SetString(disk+"/Free", "", uom.STARTED)
		c.SetString(disk+"/Used", "", uom.STARTED)
		c.SetString(disk+"/Utilization", "", uom.STARTED)
	}

	go c.Tick()
	return nil
}

func (c *UnitStorage) InternalUnitStop() {
}

func (c *UnitStorage) GetConfigMeta() string {
	meta := units_common.NewUnitConfigItem("", "", "", "", "", "", "")
	return meta.Marshal()
}

func (c *UnitStorage) Tick() {
	var err error
	c.Started = true

	envolvedItems := make(map[string]string)

	for !c.Stopping {
		for i := 0; i < 10; i++ {
			if c.Stopping {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}

		drives := c.drives()

		var TotalSpace uint64
		var UsedSpace uint64

		for _, disk := range drives {
			var free, total, avail uint64
			var diskName *uint16
			diskName, err = syscall.UTF16PtrFromString(disk + ":\\")

			err = windows.GetDiskFreeSpaceEx(
				diskName,
				&free,
				&total,
				&avail,
			)

			envolvedItems[disk+"/Total"] = ""
			envolvedItems[disk+"/Free"] = ""
			envolvedItems[disk+"/Used"] = ""
			envolvedItems[disk+"/Utilization"] = ""

			if err != nil {
				c.SetString(disk+"/Total", "", "error")
				c.SetString(disk+"/Free", "", "error")
				c.SetString(disk+"/Used", "", "error")
				c.SetString(disk+"/Utilization", "", "error")
			} else {
				c.SetUInt64(disk+"/Total", total/1024/1024, "MB")
				//c.SetUInt64(disk+"/Available", avail / 1024 / 1024, "MB")
				c.SetUInt64(disk+"/Free", free/1024/1024, "MB")
				c.SetUInt64(disk+"/Used", (total-free)/1024/1024, "MB")
				c.SetFloat64(disk+"/Utilization", 100*float64(total-free)/float64(total), "%", 1)

				TotalSpace += total
				UsedSpace += total - free
			}
		}

		//summaryTotal := strconv.FormatFloat(float64(TotalSpace) / 1024 / 1024 / 1024 / 1024, 'f', 1, 64)
		//summaryUsed := strconv.FormatFloat(float64(UsedSpace) / 1024 / 1024 / 1024 / 1024, 'f', 1, 64)
		summaryUtilization := strconv.FormatFloat(100*float64(UsedSpace)/float64(TotalSpace), 'f', 1, 64)

		summary := summaryUtilization

		envolvedItems["UsedPercents"] = ""
		c.SetString("UsedPercents", summary, "%")
	}

	for envItem, _ := range envolvedItems {
		c.SetString(envItem, "", uom.STOPPED)
	}

	c.Started = false
}
