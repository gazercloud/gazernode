package units_system

import (
	"github.com/gazercloud/gazernode/common_interfaces"
)

type UnitType struct {
	TypeCode    string `json:"type_code"`
	Category    string `json:"category"`
	DisplayName string `json:"display_name"`
	Help        string `json:"help"`
	Description string `json:"description"`
	Picture     []byte `json:"picture"`
	ConfigMeta  string `json:"config_meta"`

	Constructor func() common_interfaces.IUnit
}
