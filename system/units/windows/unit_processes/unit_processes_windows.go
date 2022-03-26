package unit_processes

import (
	"encoding/json"
	"errors"
	"github.com/gazercloud/gazernode/resources"
	"github.com/gazercloud/gazernode/utilities/logger"
	"github.com/gazercloud/gazernode/utilities/uom"
	"golang.org/x/sys/windows"
	"syscall"
	"time"
	"unsafe"
)

var Image []byte

func init() {
	Image = resources.R_files_sensors_unit_comruter_process_png
}

func (c *UnitSystemProcesses) InternalUnitStart() error {
	var err error
	type Config struct {
		Period float64 `json:"period"`
	}

	c.SetMainItem("Count")

	var config Config
	err = json.Unmarshal([]byte(c.GetConfig()), &config)
	if err != nil {
		logger.Println("ERROR[UnitSystemProcess]:", err)
		err = errors.New("config error")
		c.SetString("Count", err.Error(), uom.ERROR)
		return err
	}

	c.periodMs = int(config.Period)
	if c.periodMs < 100 {
		err = errors.New("wrong period")
		c.SetString("Count", err.Error(), uom.ERROR)
		return err
	}

	go c.Tick()
	return nil
}

func (c *UnitSystemProcesses) InternalUnitStop() {
}

func (c *UnitSystemProcesses) Tick() {
	c.Started = true
	logger.Println("UNIT <Processes Windows> started:", c.Id())
	c.SetStringForAll("", uom.STARTED)

	dtOperationTime := time.Now().UTC()

	for !c.Stopping {
		for {
			if c.Stopping || time.Now().Sub(dtOperationTime) > time.Duration(c.periodMs)*time.Millisecond {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
		if c.Stopping {
			break
		}

		count := float64(0)

		handle, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
		if err == nil {
			var entry windows.ProcessEntry32
			entry.Size = uint32(unsafe.Sizeof(entry))
			err = windows.Process32First(handle, &entry)
			for err == nil {
				nameSize := 0
				for i := 0; i < 260; i++ {
					if entry.ExeFile[nameSize] == 0 {
						break
					}
					nameSize++
				}

				id := entry.ProcessID
				_ = id
				name := syscall.UTF16ToString(entry.ExeFile[:nameSize])
				_ = name

				count += 1

				err = windows.Process32Next(handle, &entry)
			}

			_ = windows.CloseHandle(handle)
		}

		c.SetFloat64("Count", count, uom.NONE, 0)

		dtOperationTime = time.Now().UTC()
	}

	logger.Println("UNIT <Processes Windows> stopped:", c.Id())
	c.Started = false
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

type PROCESS_MEMORY_COUNTERS_EX struct {
	CB                         uint32
	PageFaultCount             uint32
	PeakWorkingSetSize         uint64
	WorkingSetSize             uint64
	QuotaPeakPagedPoolUsage    uint64
	QuotaPagedPoolUsage        uint64
	QuotaPeakNonPagedPoolUsage uint64
	QuotaNonPagedPoolUsage     uint64
	PagefileUsage              uint64
	PeakPagefileUsage          uint64
	PrivateUsage               uint64
}

func errnoErr(e syscall.Errno) error {
	switch e {
	case 0:
		return nil
	}
	return e
}

func GetProcessMemoryInfo(handle windows.Handle) (pc PROCESS_MEMORY_COUNTERS_EX, err error) {
	var res PROCESS_MEMORY_COUNTERS_EX
	r1, _, e1 := syscall.Syscall(getProcessMemoryInfo.Addr(), 3, uintptr(handle), uintptr(unsafe.Pointer(&res)), uintptr(80))
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return res, err
}

func ProcessThreadsCount(processId uint32) int {
	count := 0
	handle, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPTHREAD, 0)
	if err == nil {
		var entry windows.ThreadEntry32
		entry.Size = uint32(unsafe.Sizeof(entry))
		err := windows.Thread32First(handle, &entry)
		for err == nil {
			if entry.OwnerProcessID == processId {
				count++
			}
			entry.Size = uint32(unsafe.Sizeof(entry))
			err = windows.Thread32Next(handle, &entry)
		}

		_ = windows.CloseHandle(handle)
	}

	return count
}

func GetProcessHandleCount(handle windows.Handle) int {
	var res uint32
	syscall.Syscall(getProcessHandleCount.Addr(), 2, uintptr(handle), uintptr(unsafe.Pointer(&res)), 0)
	return int(res)
}

func GetProcessIoCounters(handle windows.Handle) (cnt windows.IO_COUNTERS, err error) {
	var res windows.IO_COUNTERS
	r1, _, e1 := syscall.Syscall(getProcessIoCounters.Addr(), 2, uintptr(handle), uintptr(unsafe.Pointer(&res)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return res, err
}

func GetGuiResources(handle windows.Handle) (cntGDI int64, cntUser int64, cntGDIPeak int64, cntUserPeak int64, err error) {
	var flags uint32
	flags = 0
	r1, _, _ := syscall.Syscall(getGuiResources.Addr(), 2, uintptr(handle), uintptr(flags), 0)
	cntGDI = int64(r1)
	flags = 1
	r2, _, _ := syscall.Syscall(getGuiResources.Addr(), 2, uintptr(handle), uintptr(flags), 0)
	cntUser = int64(r2)
	flags = 2
	r3, _, _ := syscall.Syscall(getGuiResources.Addr(), 2, uintptr(handle), uintptr(flags), 0)
	cntGDIPeak = int64(r3)
	flags = 4
	r4, _, _ := syscall.Syscall(getGuiResources.Addr(), 2, uintptr(handle), uintptr(flags), 0)
	cntUserPeak = int64(r4)
	return
}

func GetProcesses() []ProcessInfo {
	result := make([]ProcessInfo, 0)
	handle, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err == nil {
		var entry windows.ProcessEntry32
		entry.Size = uint32(unsafe.Sizeof(entry))
		err = windows.Process32First(handle, &entry)
		for err == nil {
			nameSize := 0
			for i := 0; i < 260; i++ {
				if entry.ExeFile[nameSize] == 0 {
					break
				}
				nameSize++
			}
			name := syscall.UTF16ToString(entry.ExeFile[:nameSize])

			var pi ProcessInfo
			pi.Id = int(entry.ProcessID)
			pi.Name = name
			result = append(result, pi)

			err = windows.Process32Next(handle, &entry)
		}

		_ = windows.CloseHandle(handle)
	}

	return result
}
