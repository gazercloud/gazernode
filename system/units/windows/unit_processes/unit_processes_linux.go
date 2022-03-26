package unit_processes

import (
	"encoding/json"
	"errors"
	"github.com/gazercloud/gazernode/resources"
	"github.com/gazercloud/gazernode/utilities/logger"
	"github.com/prometheus/procfs"
	"strconv"
	"strings"
	"time"
)

var Image []byte

func init() {
	Image = resources.R_files_sensors_unit_comruter_process_png
}

func (c *UnitSystemProcesses) InternalUnitStart() error {
	var err error
	type Config struct {
		ProcessName string  `json:"process_name"`
		Period      float64 `json:"period"`
	}

	c.SetMainItem("ResidentMemory")

	{
		c.SetString("PID", "", "")
		c.SetString("ResidentMemory", "", "")
		c.SetString("VirtualMemory", "", "")
		c.SetString("CPU", "", "")
		c.SetString("FileDescriptors", "", "")
		c.SetString("Status", "", "")
		c.SetString("Command", "", "")
		c.SetString("Executable", "", "")
	}

	var config Config
	err = json.Unmarshal([]byte(c.GetConfig()), &config)
	if err != nil {
		logger.Println("ERROR[UnitSystemProcess]:", err)
		err = errors.New("config error")
		c.SetString("Status", err.Error(), "error")
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
		c.SetString("Status", err.Error(), "error")
		return err
	}

	c.periodMs = int(config.Period)
	if c.periodMs < 100 {
		err = errors.New("wrong period")
		c.SetString("Status", err.Error(), "error")
		return err
	}

	go c.Tick()
	return nil
}

func (c *UnitSystemProcesses) InternalUnitStop() {
}

func (c *UnitSystemProcesses) Tick() {
	c.Started = true
	logger.Println("UNIT <Process Windows> started:", c.Id())

	dtOperationTime := time.Now().UTC()

	processId := int(-1)
	var proc procfs.Proc

	lastCpuValid := false
	lastCpuValue := float64(0)
	lastCpuTime := time.Now()

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

		if processId == -1 {

			var err error

			allProcesses, err := procfs.AllProcs()
			if err != nil {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			for _, p := range allProcesses {
				matchId := false
				matchName := false

				if c.processIdActive {
					if int(c.processId) == p.PID {
						matchId = true
					}
				} else {
					matchId = true
				}

				if c.processNameActive {
					if comm, err := p.Comm(); err == nil && strings.Contains(comm, c.processName) {
						matchName = true
					}
				} else {
					matchName = true
				}

				if matchId && matchName {
					processId = p.PID
					proc = p

					comm, err := p.Comm()
					if err == nil {
						c.SetString("Command", comm, "")
					}
					exe, err := p.Executable()
					if err == nil {
						c.SetString("Executable", exe, "")
					}
					c.SetFloat64("PID", float64(p.PID), "", 0)
				}
			}
		}

		if processId == -1 {
			time.Sleep(100 * time.Millisecond)
			{
				c.SetString("Status", "no process found", "error")
				c.SetString("PID", "", "error")
				c.SetString("Command", "", "error")
				c.SetString("Executable", "", "error")
				c.SetString("ResidentMemory", "", "error")
				c.SetString("VirtualMemory", "", "error")
				c.SetString("CPU", "", "error")
				c.SetString("FileDescriptors", "", "error")
			}
			continue
		}

		pStat, err := proc.Stat()
		if err == nil {
			tNow := time.Now()
			dur := tNow.Sub(lastCpuTime).Seconds()
			cpuTime := pStat.CPUTime()

			if lastCpuValid {
				value := cpuTime - lastCpuValue
				if dur > 0.0000001 {
					usage := value / dur
					c.SetFloat64("CPU", usage*100, "%", 2)
				}
			}

			lastCpuTime = tNow
			lastCpuValue = cpuTime
			lastCpuValid = true

			c.SetFloat64("ResidentMemory", float64(pStat.ResidentMemory()/1024), "KB", 0)
			c.SetFloat64("VirtualMemory", float64(pStat.VirtualMemory()/1024), "KB", 0)
			c.SetString("Status", pStat.State, "")
		} else {
			c.SetString("ResidentMemory", "", "error")
			c.SetString("CPUTime", "", "error")
			c.SetString("VirtualMemory", "", "error")
			c.SetString("Status", err.Error(), "error")
			c.SetString("Command", "", "error")
			c.SetString("Executable", "", "error")
			lastCpuValid = false
			processId = -1
		}

		fdInfo, err := proc.FileDescriptorsInfo()
		if err == nil {
			c.SetInt("FileDescriptors", fdInfo.Len(), "")
		} else {
			c.SetString("FileDescriptors", "", "error")
			processId = -1
			lastCpuValid = false
		}

		dtOperationTime = time.Now().UTC()
	}

	{
		c.SetString("PID", "", "stopped")
		c.SetString("ResidentMemory", "", "stopped")
		c.SetString("VirtualMemory", "", "stopped")
		c.SetString("CPU", "", "stopped")
		c.SetString("FileDescriptors", "", "stopped")
		c.SetString("Status", "", "stopped")
		c.SetString("Command", "", "stopped")
		c.SetString("Executable", "", "stopped")
	}

	logger.Println("UNIT <Process Windows> stopped:", c.Id())
	c.Started = false
}

func GetProcesses() []ProcessInfo {
	result := make([]ProcessInfo, 0)

	allProcesses, err := procfs.AllProcs()
	if err != nil {
		return result
	}

	for _, p := range allProcesses {
		var proc ProcessInfo
		proc.Id = p.PID
		proc.Name, _ = p.Comm()
		proc.Info, _ = p.Executable()
		result = append(result, proc)
	}

	return result
}
