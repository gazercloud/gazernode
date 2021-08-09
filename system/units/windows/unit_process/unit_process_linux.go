package unit_process

import (
	"encoding/json"
	"errors"
	"github.com/gazercloud/gazernode/logger"
	"github.com/gazercloud/gazernode/resources"
	"strconv"
	"syscall"
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

	c.SetMainItem("Main/Working Set Size")

	{
		// Common
		c.SetString("Common/Name", "", "")
		c.SetString("Common/ProcessID", "", "")
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

	//processId := int32(-1)

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

		var ru syscall.Rusage

		var err error

		err = syscall.Getrusage(0, &ru)

		if err == nil {
			// Common
			c.SetInt64("Maxrss", int64(ru.Maxrss), "")
			c.SetInt64("Ixrss", int64(ru.Ixrss), "")
			c.SetInt64("Idrss", int64(ru.Idrss), "")
			c.SetInt64("Isrss", int64(ru.Isrss), "")
			c.SetInt64("Minflt", int64(ru.Minflt), "")
			c.SetInt64("Majflt", int64(ru.Majflt), "")
			c.SetInt64("Nswap", int64(ru.Nswap), "")
			c.SetInt64("Inblock", int64(ru.Inblock), "")
			c.SetInt64("Oublock", int64(ru.Oublock), "")
			c.SetInt64("Msgsnd", int64(ru.Msgsnd), "")
			c.SetInt64("Msgrcv", int64(ru.Msgrcv), "")
			c.SetInt64("Nsignals", int64(ru.Nsignals), "")
			c.SetInt64("Nvcsw", int64(ru.Nvcsw), "")
			c.SetInt64("Nivcsw", int64(ru.Nivcsw), "")
			c.SetString("Status", "ok", "")
		} else {
			c.SetString("Status", err.Error(), "")
		}

		dtOperationTime = time.Now().UTC()
	}

	{
		// Common
		c.SetString("Common/Name", "", "stopped")
		c.SetString("Common/ProcessID", "", "stopped")

	}

	logger.Println("UNIT <Process Windows> stopped:", c.Id())
	c.Started = false
}

func GetProcesses() []ProcessInfo {
	result := make([]ProcessInfo, 0)

	return result
}
