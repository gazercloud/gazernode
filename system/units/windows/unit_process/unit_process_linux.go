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

		var ru syscall.Rusage

		err = syscall.Getrusage(0, &ru)

		if err != nil {
			// Common
			c.SetInt64("Maxrss", ru.Maxrss, "")
			c.SetInt64("Ixrss", ru.Ixrss, "")
			c.SetInt64("Idrss", ru.Idrss, "")
			c.SetInt64("Isrss", ru.Isrss, "")
			c.SetInt64("Minflt", ru.Minflt, "")
			c.SetInt64("Majflt", ru.Majflt, "")
			c.SetInt64("Nswap", ru.Nswap, "")
			c.SetInt64("Inblock", ru.Inblock, "")
			c.SetInt64("Oublock", ru.Oublock, "")
			c.SetInt64("Msgsnd", ru.Msgsnd, "")
			c.SetInt64("Msgrcv", ru.Msgrcv, "")
			c.SetInt64("Nsignals", ru.Nsignals, "")
			c.SetInt64("Nvcsw", ru.Nvcsw, "")
			c.SetInt64("Nivcsw", ru.Nivcsw, "")
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
