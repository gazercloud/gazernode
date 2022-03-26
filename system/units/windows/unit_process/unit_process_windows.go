package unit_process

import (
	"encoding/json"
	"errors"
	"github.com/gazercloud/gazernode/resources"
	"github.com/gazercloud/gazernode/utilities/logger"
	"github.com/gazercloud/gazernode/utilities/uom"
	"golang.org/x/sys/windows"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

var Image []byte

func init() {
	Image = resources.R_files_sensors_unit_comruter_process_png
}

func (c *UnitSystemProcess) InternalUnitStart() error {
	var err error
	type Config struct {
		ProcessName string  `json:"process_name"`
		Period      float64 `json:"period"`
	}

	c.SetMainItem("Main/Working Set Size")

	/*{
		// Common
		c.SetString("Common/Name", "", "")
		c.SetString("Common/ProcessID", "", "")

		// Main
		c.SetString("Main/Working Set Size", "", "")
		c.SetString("Main/Thread Count", "", "")
		c.SetString("Main/Handle Count", "", "")
		c.SetString("Main/GDI Objects", "", "")
		c.SetString("Main/GDI Objects Peak", "", "")
		c.SetString("Main/User Objects", "", "")
		c.SetString("Main/User Objects Peak", "", "")

		// Operations
		c.SetString("Operations/Read Operation Count", "", "")
		c.SetString("Operations/Read Transfer Count", "", "")
		c.SetString("Operations/Write Operation Count", "", "")
		c.SetString("Operations/Write Transfer Count", "", "")
		c.SetString("Operations/Other Operation Count", "", "")
		c.SetString("Operations/Other Transfer Count", "", "")

		// CPU
		c.SetString("CPU/Kernel Mode Time", "", "")
		c.SetString("CPU/User Mode Time", "", "")
		c.SetString("CPU/Usage", "", "")
		c.SetString("CPU/Usage Kernel", "", "")
		c.SetString("CPU/Usage User", "", "")

		// Memory
		c.SetString("Memory/Page Faults", "", "")
		c.SetString("Memory/Peak Working SetSize", "", "")
		c.SetString("Memory/Private Usage", "", "")
	}*/

	var config Config
	err = json.Unmarshal([]byte(c.GetConfig()), &config)
	if err != nil {
		logger.Println("ERROR[UnitSystemProcess]:", err)
		err = errors.New("config error")
		c.SetString("Common/ProcessID", err.Error(), uom.ERROR)
		return err
	}

	runes := []rune(config.ProcessName)
	indexOfSlash := -1
	for i, r := range runes {
		if r == '#' {
			indexOfSlash = i
		}
	}

	if indexOfSlash > -1 {
		var prId uint64
		prId, err = strconv.ParseUint(string(runes[indexOfSlash+1:]), 10, 32)
		if err == nil {
			c.processId = uint32(prId)
			c.processIdActive = true
		} else {
			c.processId = 0
			c.processIdActive = false
		}

		c.processNameActive = false
		c.processName = ""

		if indexOfSlash > 0 {
			c.processNameActive = true
			c.processName = string(runes[:indexOfSlash])
		}
	} else {
		c.processIdActive = false
		c.processId = 0
		c.processNameActive = true
		c.processName = config.ProcessName
	}

	if !c.processIdActive && !c.processNameActive {
		err = errors.New("wrong filter")
		c.SetString("Common/ProcessID", err.Error(), uom.ERROR)
		return err
	}

	c.periodMs = int(config.Period)
	if c.periodMs < 100 {
		err = errors.New("wrong period")
		c.SetString("Common/ProcessID", err.Error(), uom.ERROR)
		return err
	}

	go c.Tick()
	return nil
}

func (c *UnitSystemProcess) InternalUnitStop() {
}

func (c *UnitSystemProcess) Tick() {
	c.Started = true
	logger.Println("UNIT <Process Windows> started:", c.Id())

	dtOperationTime := time.Now().UTC()

	processId := int32(-1)

	lastKernelTimeMs := int64(0)
	lastUserTimeMs := int64(0)
	lastReadProcessTimes := time.Now().UTC()

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

		processJustFound := false

		if processId < 0 {
			if !c.processIdActive && !c.processNameActive {
				time.Sleep(100 * time.Millisecond)
				continue
			}

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
					name := syscall.UTF16ToString(entry.ExeFile[:nameSize])

					// Filtering
					matchId := false
					matchName := false
					if c.processIdActive {
						if id == c.processId {
							matchId = true
						}
					} else {
						matchId = true
					}
					if c.processNameActive {
						if strings.ToLower(name) == strings.ToLower(c.processName) {
							matchName = true
						}
					} else {
						matchName = true
					}
					if matchId && matchName {
						processId = int32(entry.ProcessID)
						processJustFound = true
						c.actualProcessName = name
						break
					}
					// /////////////////

					err = windows.Process32Next(handle, &entry)
				}

				_ = windows.CloseHandle(handle)
			}
		}

		if processId >= 0 {
			hProcess, err := windows.OpenProcess(windows.STANDARD_RIGHTS_REQUIRED|windows.SYNCHRONIZE|windows.SPECIFIC_RIGHTS_ALL, false, uint32(processId))
			if err == nil {
				// Common
				c.SetString("Common/Name", c.actualProcessName, uom.NONE)
				c.SetUInt32("Common/ProcessID", uint32(processId), uom.NONE)

				{
					res, _ := GetProcessMemoryInfo(hProcess)
					c.SetUInt64("Main/Working Set Size", res.WorkingSetSize/1024, uom.KB)

					c.SetUInt32("Memory/Page Faults", res.PageFaultCount, uom.NONE)
					c.SetUInt64("Memory/Peak Working SetSize", res.PeakWorkingSetSize/1024, uom.KB)
					c.SetUInt64("Memory/Private Usage", res.PrivateUsage/1024, uom.KB)
				}

				c.SetUInt32("Main/Thread Count", uint32(ProcessThreadsCount(uint32(processId))), uom.NONE)
				c.SetUInt32("Main/Handle Count", uint32(GetProcessHandleCount(hProcess)), uom.NONE)
				{
					cntGDI, cntUser, cntGDIPeak, cntUserPeak, _ := GetGuiResources(hProcess)
					c.SetInt64("Main/GDI Objects", cntGDI, uom.NONE)
					c.SetInt64("Main/GDI Objects Peak", cntGDIPeak, uom.NONE)
					c.SetInt64("Main/User Objects", cntUser, uom.NONE)
					c.SetInt64("Main/User Objects Peak", cntUserPeak, uom.NONE)
				}

				{
					var ftStart windows.Filetime
					var ftEnd windows.Filetime
					var ftKernel windows.Filetime
					var ftUser windows.Filetime
					err = windows.GetProcessTimes(hProcess, &ftStart, &ftEnd, &ftKernel, &ftUser)
					if err == nil {
						kernelTimeMs := (int64(ftKernel.HighDateTime)<<32 + int64(ftKernel.LowDateTime)) / 100000
						userTimeMs := (int64(ftUser.HighDateTime)<<32 + int64(ftUser.LowDateTime)) / 100000
						deltaKernel := kernelTimeMs - lastKernelTimeMs
						deltaUser := userTimeMs - lastUserTimeMs
						duringMs := time.Now().UTC().Sub(lastReadProcessTimes).Milliseconds()

						usageCpuKernel := float64(0)
						usageCpuUser := float64(0)
						usageCpu := float64(0)

						if duringMs > 0 {
							usageCpuKernel = float64(deltaKernel) / float64(duringMs)
							usageCpuUser = float64(deltaUser) / float64(duringMs)
							usageCpu = float64(deltaKernel+deltaUser) / float64(duringMs)
						}

						lastReadProcessTimes = time.Now().UTC()
						lastKernelTimeMs = kernelTimeMs
						lastUserTimeMs = userTimeMs

						if processJustFound {
							usageCpuKernel = 0
							usageCpuUser = 0
							usageCpu = 0
						}

						c.SetInt64("CPU/Kernel Mode Time", kernelTimeMs, uom.MS)
						c.SetInt64("CPU/User Mode Time", userTimeMs, uom.MS)
						c.SetFloat64("CPU/Usage", usageCpu*100, uom.PERCENTS, 1)
						c.SetFloat64("CPU/Usage Kernel", usageCpuKernel*100, uom.PERCENTS, 1)
						c.SetFloat64("CPU/Usage User", usageCpuUser*100, uom.PERCENTS, 1)

						{
							res, _ := GetProcessIoCounters(hProcess)
							c.SetUInt64("Operations/Read Operation Count", res.ReadOperationCount, uom.NONE)
							c.SetUInt64("Operations/Read Transfer Count", res.ReadTransferCount, uom.NONE)
							c.SetUInt64("Operations/Write Operation Count", res.WriteOperationCount, uom.NONE)
							c.SetUInt64("Operations/Write Transfer Count", res.WriteTransferCount, uom.NONE)
							c.SetUInt64("Operations/Other Operation Count", res.OtherOperationCount, uom.NONE)
							c.SetUInt64("Operations/Other Transfer Count", res.OtherTransferCount, uom.NONE)
						}

					}
				}

				//GetGuiResources function (winuser.h)

				_ = windows.CloseHandle(hProcess)
			} else {
				processId = -1
			}
		}

		if processId < 0 {
			// Common
			c.SetString("Common/Name", c.processName, uom.NONE)
			c.SetString("Common/ProcessID", "not found", uom.ERROR)

			// Main
			c.SetString("Main/Working Set Size", "", uom.ERROR)
			c.SetString("Main/Thread Count", "", uom.ERROR)
			c.SetString("Main/Handle Count", "", uom.ERROR)
			c.SetString("Main/GDI Objects", "", uom.ERROR)
			c.SetString("Main/GDI Objects Peak", "", uom.ERROR)
			c.SetString("Main/User Objects", "", uom.ERROR)
			c.SetString("Main/User Objects Peak", "", uom.ERROR)

			// Operations
			c.SetString("Operations/Read Operation Count", "", uom.ERROR)
			c.SetString("Operations/Read Transfer Count", "", uom.ERROR)
			c.SetString("Operations/Write Operation Count", "", uom.ERROR)
			c.SetString("Operations/Write Transfer Count", "", uom.ERROR)
			c.SetString("Operations/Other Operation Count", "", uom.ERROR)
			c.SetString("Operations/Other Transfer Count", "", uom.ERROR)

			// CPU
			c.SetString("CPU/Kernel Mode Time", "", uom.ERROR)
			c.SetString("CPU/User Mode Time", "", uom.ERROR)
			c.SetString("CPU/Usage", "", uom.ERROR)
			c.SetString("CPU/Usage Kernel", "", uom.ERROR)
			c.SetString("CPU/Usage User", "", uom.ERROR)
			// Memory
			c.SetString("Memory/Page Faults", "", uom.ERROR)
			c.SetString("Memory/Peak Working SetSize", "", uom.ERROR)
			c.SetString("Memory/Private Usage", "", uom.ERROR)
		}

		dtOperationTime = time.Now().UTC()
	}

	{
		// Common
		c.SetString("Common/Name", "", uom.STOPPED)
		c.SetString("Common/ProcessID", "", uom.STOPPED)

		// Main
		c.SetString("Main/Working Set Size", "", uom.STOPPED)
		c.SetString("Main/Thread Count", "", uom.STOPPED)
		c.SetString("Main/Handle Count", "", uom.STOPPED)
		c.SetString("Main/GDI Objects", "", uom.STOPPED)
		c.SetString("Main/GDI Objects Peak", "", uom.STOPPED)
		c.SetString("Main/User Objects", "", uom.STOPPED)
		c.SetString("Main/User Objects Peak", "", uom.STOPPED)

		// Operations
		c.SetString("Operations/Read Operation Count", "", uom.STOPPED)
		c.SetString("Operations/Read Transfer Count", "", uom.STOPPED)
		c.SetString("Operations/Write Operation Count", "", uom.STOPPED)
		c.SetString("Operations/Write Transfer Count", "", uom.STOPPED)
		c.SetString("Operations/Other Operation Count", "", uom.STOPPED)
		c.SetString("Operations/Other Transfer Count", "", uom.STOPPED)

		// CPU
		c.SetString("CPU/Kernel Mode Time", "", uom.STOPPED)
		c.SetString("CPU/User Mode Time", "", uom.STOPPED)
		c.SetString("CPU/Usage", "", uom.STOPPED)
		c.SetString("CPU/Usage Kernel", "", uom.STOPPED)
		c.SetString("CPU/Usage User", "", uom.STOPPED)
		// Memory
		c.SetString("Memory/Page Faults", "", uom.STOPPED)
		c.SetString("Memory/Peak Working SetSize", "", uom.STOPPED)
		c.SetString("Memory/Private Usage", "", uom.STOPPED)
	}

	logger.Println("UNIT <Process Windows> stopped:", c.Id())
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
