package units_common

import "github.com/gazercloud/gazernode/common_interfaces"

type UnitInfo struct {
	Id             string                           `json:"id"`
	DisplayName    string                           `json:"name"`
	Type           string                           `json:"type"`
	TypeForDisplay string                           `json:"type_for_display"`
	Config         string                           `json:"config"`
	Enable         bool                             `json:"enable"`
	Properties     []common_interfaces.ItemProperty `json:"p"`
}
