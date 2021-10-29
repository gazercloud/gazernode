package unit_process

import (
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/system/units/units_common"
)

type UnitSystemProcess struct {
	units_common.Unit

	processIdActive   bool
	processId         uint32
	processNameActive bool
	processName       string
	periodMs          int

	actualProcessName string
}

func New() common_interfaces.IUnit {
	var c UnitSystemProcess
	return &c
}

func (c *UnitSystemProcess) GetConfigMeta() string {
	meta := units_common.NewUnitConfigItem("", "", "", "", "", "", "")
	meta.Add("process_name", "Process Name", "notepad.exe", "string", "", "", "processes")
	meta.Add("period", "Period, ms", "1000", "num", "0", "999999", "0")
	return meta.Marshal()
}

type ProcessInfo struct {
	Name string
	Id   int
	Info string
}
