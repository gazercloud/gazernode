package units_common

type UnitInfo struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	Type           string `json:"type"`
	TypeForDisplay string `json:"type_for_display"`
	Config         string `json:"config"`
	Enable         bool   `json:"enable"`
}
