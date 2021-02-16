package unit_process

import (
	"errors"
	"fmt"
	"github.com/shirou/gopsutil/process"
	"time"
)

var Image []byte

func init() {
}

func (c *UnitSystemProcess) InternalUnitStart() error {
	//go c.Tick()
	return errors.New("not implemented")
}

func (c *UnitSystemProcess) InternalUnitStop() {
}

func (c *UnitSystemProcess) Tick() {
	c.Started = true

	for !c.Stopping {
		time.Sleep(1000 * time.Millisecond)
		found := false

		processes, _ := process.Processes()
		for _, pr := range processes {
			name, err := pr.Name()
			if err == nil {
				if name == "chrome.exe" {
					found = true
					break
				}
			}
		}

		if !found {
			c.SetString("error", "no process found", "")
		} else {
			c.SetString("error", "", "")
		}
	}

	fmt.Println("stopped")
	c.Started = false
}

func GetProcesses() []ProcessInfo {
	result := make([]ProcessInfo, 0)
	return result
}
