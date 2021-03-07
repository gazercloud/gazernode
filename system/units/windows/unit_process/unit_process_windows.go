package unit_process

import (
	"encoding/json"
	"errors"
	"github.com/gazercloud/gazernode/logger"
	"github.com/gazercloud/gazernode/resources"
	"golang.org/x/sys/windows"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

var Image []byte

func init() {
	Image = resources.R_files_sensors_sensor_windows_proc_png
}

func (c *UnitSystemProcess) InternalUnitStart() error {
	var err error
	type Config struct {
		ProcessName string  `json:"process_name"`
		Period      float64 `json:"period"`
	}

	c.SetMainItem("Main/Working Set Size")

	{
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
	}

	var config Config
	err = json.Unmarshal([]byte(c.GetConfig()), &config)
	if err != nil {
		logger.Println("ERROR[UnitSystemProcess]:", err)
		err = errors.New("config error")
		c.SetString("Common/ProcessID", err.Error(), "error")
		return err
	}

	c.processName = config.ProcessName
	if c.processName == "" {
		err = errors.New("wrong address")
		c.SetString("Common/ProcessID", err.Error(), "error")
		return err
	}

	c.periodMs = int(config.Period)
	if c.periodMs < 100 {
		err = errors.New("wrong period")
		c.SetString("Common/ProcessID", err.Error(), "error")
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
					if strings.ToLower(name) == strings.ToLower(c.processName) {
						processId = int32(entry.ProcessID)
						processJustFound = true
						break
					}
					err = windows.Process32Next(handle, &entry)
				}

				_ = windows.CloseHandle(handle)
			}
		}

		if processId >= 0 {
			hProcess, err := windows.OpenProcess(windows.STANDARD_RIGHTS_REQUIRED|windows.SYNCHRONIZE|windows.SPECIFIC_RIGHTS_ALL, false, uint32(processId))
			if err == nil {
				// Common
				c.SetString("Common/Name", c.processName, "")
				c.SetUInt32("Common/ProcessID", uint32(processId), "")

				{
					res, _ := GetProcessMemoryInfo(hProcess)
					c.SetUInt64("Main/Working Set Size", res.WorkingSetSize/1024, "KB")

					c.SetUInt32("Memory/Page Faults", res.PageFaultCount, "")
					c.SetUInt64("Memory/Peak Working SetSize", res.PeakWorkingSetSize/1024, "KB")
					c.SetUInt64("Memory/Private Usage", res.PrivateUsage/1024, "KB")
				}

				c.SetUInt32("Main/Thread Count", uint32(ProcessThreadsCount(uint32(processId))), "")
				c.SetUInt32("Main/Handle Count", uint32(GetProcessHandleCount(hProcess)), "")
				{
					cntGDI, cntUser, cntGDIPeak, cntUserPeak, _ := GetGuiResources(hProcess)
					c.SetInt64("Main/GDI Objects", cntGDI, "")
					c.SetInt64("Main/GDI Objects Peak", cntGDIPeak, "")
					c.SetInt64("Main/User Objects", cntUser, "")
					c.SetInt64("Main/User Objects Peak", cntUserPeak, "")
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

						c.SetInt64("CPU/Kernel Mode Time", kernelTimeMs, "ms")
						c.SetInt64("CPU/User Mode Time", userTimeMs, "ms")
						c.SetFloat64("CPU/Usage", usageCpu*100, "%", 1)
						c.SetFloat64("CPU/Usage Kernel", usageCpuKernel*100, "%", 1)
						c.SetFloat64("CPU/Usage User", usageCpuUser*100, "%", 1)

						{
							res, _ := GetProcessIoCounters(hProcess)
							c.SetUInt64("Operations/Read Operation Count", res.ReadOperationCount, "")
							c.SetUInt64("Operations/Read Transfer Count", res.ReadTransferCount, "")
							c.SetUInt64("Operations/Write Operation Count", res.WriteOperationCount, "")
							c.SetUInt64("Operations/Write Transfer Count", res.WriteTransferCount, "")
							c.SetUInt64("Operations/Other Operation Count", res.OtherOperationCount, "")
							c.SetUInt64("Operations/Other Transfer Count", res.OtherTransferCount, "")
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
			c.SetString("Common/Name", c.processName, "")
			c.SetString("Common/ProcessID", "not found", "error")

			// Main
			c.SetUInt64("Main/Working Set Size", 0, "error")
			c.SetUInt32("Main/Thread Count", 0, "error")
			c.SetUInt32("Main/Handle Count", 0, "error")
			c.SetInt64("Main/GDI Objects", 0, "error")
			c.SetInt64("Main/GDI Objects Peak", 0, "error")
			c.SetInt64("Main/User Objects", 0, "error")
			c.SetInt64("Main/User Objects Peak", 0, "error")

			// Operations
			c.SetUInt64("Operations/Read Operation Count", 0, "error")
			c.SetUInt64("Operations/Read Transfer Count", 0, "error")
			c.SetUInt64("Operations/Write Operation Count", 0, "error")
			c.SetUInt64("Operations/Write Transfer Count", 0, "error")
			c.SetUInt64("Operations/Other Operation Count", 0, "error")
			c.SetUInt64("Operations/Other Transfer Count", 0, "error")

			// CPU
			c.SetUInt64("CPU/Kernel Mode Time", 0, "error")
			c.SetUInt64("CPU/User Mode Time", 0, "error")
			c.SetInt64("CPU/Usage", 0, "error")
			c.SetInt64("CPU/Usage Kernel", 0, "error")
			c.SetInt64("CPU/Usage User", 0, "error")
			// Memory
			c.SetUInt32("Memory/Page Faults", 0, "error")
			c.SetUInt32("Memory/Peak Working SetSize", 0, "error")
			c.SetUInt64("Memory/Private Usage", 0, "error")
		}

		dtOperationTime = time.Now().UTC()
	}

	{
		// Common
		c.SetString("Common/Name", "", "stopped")
		c.SetString("Common/ProcessID", "", "stopped")

		// Main
		c.SetString("Main/Working Set Size", "", "stopped")
		c.SetString("Main/Thread Count", "", "stopped")
		c.SetString("Main/Handle Count", "", "stopped")
		c.SetString("Main/GDI Objects", "", "stopped")
		c.SetString("Main/GDI Objects Peak", "", "stopped")
		c.SetString("Main/User Objects", "", "stopped")
		c.SetString("Main/User Objects Peak", "", "stopped")

		// Operations
		c.SetString("Operations/Read Operation Count", "", "stopped")
		c.SetString("Operations/Read Transfer Count", "", "stopped")
		c.SetString("Operations/Write Operation Count", "", "stopped")
		c.SetString("Operations/Write Transfer Count", "", "stopped")
		c.SetString("Operations/Other Operation Count", "", "stopped")
		c.SetString("Operations/Other Transfer Count", "", "stopped")

		// CPU
		c.SetString("CPU/Kernel Mode Time", "", "stopped")
		c.SetString("CPU/User Mode Time", "", "stopped")
		c.SetString("CPU/Usage", "", "stopped")
		c.SetString("CPU/Usage Kernel", "", "stopped")
		c.SetString("CPU/Usage User", "", "stopped")
		// Memory
		c.SetString("Memory/Page Faults", "", "stopped")
		c.SetString("Memory/Peak Working SetSize", "", "stopped")
		c.SetString("Memory/Private Usage", "", "stopped")
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
