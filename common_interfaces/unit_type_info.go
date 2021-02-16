package common_interfaces

type UnitTypeInfo struct {
	Type        string `json:"type"`
	Category    string `json:"category"`
	DisplayName string `json:"display_name"`
	Help        string `json:"help"`
	Description string `json:"description"`
	Image       []byte `json:"image"`
}
