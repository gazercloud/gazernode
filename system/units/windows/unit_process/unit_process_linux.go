package unit_process

import (
	"encoding/json"
	"errors"
	"github.com/gazercloud/gazernode/logger"
	"github.com/gazercloud/gazernode/resources"
	"github.com/prometheus/procfs"
	"strconv"
	"strings"
	"time"
)

var Image []byte

func init() {
	Image = resources.R_files_sensors_category_network_png
}

func (c *UnitSystemProcess) InternalUnitStart() error {
	var err error
	type Config struct {
		ProcessName string  `json:"process_name"`
		Period      float64 `json:"period"`
	}

	c.SetMainItem("ResidentMemory")

	{
		c.SetString("ResidentMemory", "", "stopped")
		c.SetString("VirtualMemory", "", "stopped")
		c.SetString("CPUTime", "", "stopped")
		c.SetString("FileDescriptors", "", "stopped")
	}

	var config Config
	err = json.Unmarshal([]byte(c.GetConfig()), &config)
	if err != nil {
		logger.Println("ERROR[UnitSystemProcess]:", err)
		err = errors.New("config error")
		c.SetString("Common/ProcessID", err.Error(), "error")
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

	processId := int(-1)
	var proc procfs.Proc

	//lastKernelTimeMs := int64(0)
	//lastUserTimeMs := int64(0)
	//lastReadProcessTimes := time.Now().UTC()

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

			matchId := false
			matchName := false

			allProcesses, err := procfs.AllProcs()
			if err != nil {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			for _, p := range allProcesses {
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
				}
			}
		}

		if processId == -1 {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		pStat, err := proc.Stat()
		//pStat.CSTime
		if err == nil {
			c.SetFloat64("ResidentMemory", float64(pStat.ResidentMemory()), "", 0)
			c.SetFloat64("CPUTime", float64(pStat.CPUTime()), "", 3)
			c.SetFloat64("VirtualMemory", float64(pStat.VirtualMemory()), "", 0)
		} else {
			c.SetString("ResidentMemory", "", "error")
			c.SetString("CPUTime", "", "error")
			c.SetString("VirtualMemory", "", "error")
			processId = -1
		}

		fdInfo, err := proc.FileDescriptorsInfo()
		if err == nil {
			c.SetInt("FileDescriptors", fdInfo.Len(), "")
		} else {
			c.SetString("FileDescriptors", "", "error")
			processId = -1
		}

		dtOperationTime = time.Now().UTC()
	}

	{
		c.SetString("ResidentMemory", "", "stopped")
		c.SetString("VirtualMemory", "", "stopped")
		c.SetString("CPUTime", "", "stopped")
		c.SetString("FileDescriptors", "", "stopped")
	}

	logger.Println("UNIT <Process Windows> stopped:", c.Id())
	c.Started = false
}

func GetProcesses() []ProcessInfo {
	result := make([]ProcessInfo, 0)

	return result
}
