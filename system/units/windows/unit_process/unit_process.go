package unit_process

import (
	"github.com/gazercloud/gazernode/common_interfaces"
	"github.com/gazercloud/gazernode/system/units/units_common"
)

type UnitSystemProcess struct {
	units_common.Unit

	processName string
	periodMs    int
}

func New() common_interfaces.IUnit {
	var c UnitSystemProcess
	return &c
}

func (c *UnitSystemProcess) GetConfigMeta() string {
	meta := units_common.NewUnitConfigItem("", "", "", "", "", "", "")
	meta.Add("process_name", "Process Name", "notepad.exe", "string", "", "", "processes")
	meta.Add("period", "Period, ms", "1000", "num", "0", "999999", "")
	return meta.Marshal()
}

type ProcessInfo struct {
	Name string
	Id   int
	Info string
}
